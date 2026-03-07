#!/usr/bin/env python3
"""
dac-to-drawio.py
Converts DAC files (diagram-as-code YAML) to draw.io format (.drawio / mxGraph XML).

Layout strategy:
  - All elements are positioned using ABSOLUTE coordinates (parent="1")
  - Groups are rendered as background rectangles (z-order: groups first, icons after)
  - This keeps the visual layout aligned with awsdac PNG output
  - This avoids error accumulation from nested relative coordinates

Usage:
    python3 tools/dac-to-drawio.py <input.yaml> [-o output.drawio]
"""

import yaml
import sys
import os
import argparse
import xml.etree.ElementTree as ET
from xml.dom import minidom
from collections import deque

# --- Layout constants ---
PADDING   = 36   # Inner padding inside groups (between border and children)
GAP       = 16   # Gap between sibling elements
HEADER_H  = 48   # Group header height (space for icon + label)
ICON_W    = 78   # Default width for leaf icons
ICON_H    = 78   # Default height for leaf icons
LABEL_H   = 24   # Extra height for icon label below the icon

# --- draw.io styles by DAC resource type ---
# is_group=True  -> rendered as a background rectangle (AWS group)
# is_group=False -> rendered as a resource icon
STYLES: dict = {
    # ── AWS resource icons ───────────────────────────────────────────────────
    "AWS::EC2::Instance": dict(
        style=("shape=mxgraph.aws4.resourceIcon;resIcon=mxgraph.aws4.ec2;"
               "fillColor=#ED7100;gradientColor=none;strokeColor=none;"
               "fontFamily=Helvetica;fontSize=11;"
               "verticalLabelPosition=bottom;verticalAlign=top;align=center;"),
        w=ICON_W, h=ICON_H, is_group=False,
    ),
    "AWS::EC2::InternetGateway": dict(
        style=("shape=mxgraph.aws4.resourceIcon;resIcon=mxgraph.aws4.internet_gateway;"
               "fillColor=#8C4FFF;gradientColor=none;strokeColor=none;"
               "fontFamily=Helvetica;fontSize=11;"
               "verticalLabelPosition=bottom;verticalAlign=top;align=center;"),
        w=ICON_W, h=ICON_H, is_group=False,
    ),
    "AWS::ElasticLoadBalancingV2::LoadBalancer": dict(
        style=("shape=mxgraph.aws4.resourceIcon;resIcon=mxgraph.aws4.application_load_balancer;"
               "fillColor=#8C4FFF;gradientColor=none;strokeColor=none;"
               "fontFamily=Helvetica;fontSize=11;"
               "verticalLabelPosition=bottom;verticalAlign=top;align=center;"),
        w=ICON_W, h=ICON_H, is_group=False,
    ),
    "AWS::RDS::DBInstance": dict(
        style=("shape=mxgraph.aws4.resourceIcon;resIcon=mxgraph.aws4.rds;"
               "fillColor=#C925D1;gradientColor=none;strokeColor=none;"
               "fontFamily=Helvetica;fontSize=11;"
               "verticalLabelPosition=bottom;verticalAlign=top;align=center;"),
        w=ICON_W, h=ICON_H, is_group=False,
    ),
    "AWS::Lambda::Function": dict(
        style=("shape=mxgraph.aws4.resourceIcon;resIcon=mxgraph.aws4.lambda;"
               "fillColor=#ED7100;gradientColor=none;strokeColor=none;"
               "fontFamily=Helvetica;fontSize=11;"
               "verticalLabelPosition=bottom;verticalAlign=top;align=center;"),
        w=ICON_W, h=ICON_H, is_group=False,
    ),
    "AWS::S3::Bucket": dict(
        style=("shape=mxgraph.aws4.resourceIcon;resIcon=mxgraph.aws4.s3;"
               "fillColor=#3F8624;gradientColor=none;strokeColor=none;"
               "fontFamily=Helvetica;fontSize=11;"
               "verticalLabelPosition=bottom;verticalAlign=top;align=center;"),
        w=ICON_W, h=ICON_H, is_group=False,
    ),
    "AWS::CloudFront::Distribution": dict(
        style=("shape=mxgraph.aws4.resourceIcon;resIcon=mxgraph.aws4.cloudfront;"
               "fillColor=#8C4FFF;gradientColor=none;strokeColor=none;"
               "fontFamily=Helvetica;fontSize=11;"
               "verticalLabelPosition=bottom;verticalAlign=top;align=center;"),
        w=ICON_W, h=ICON_H, is_group=False,
    ),
    "AWS::CloudFront": dict(  # service-level fallback
        style=("shape=mxgraph.aws4.resourceIcon;resIcon=mxgraph.aws4.cloudfront;"
               "fillColor=#8C4FFF;gradientColor=none;strokeColor=none;"
               "fontFamily=Helvetica;fontSize=11;"
               "verticalLabelPosition=bottom;verticalAlign=top;align=center;"),
        w=ICON_W, h=ICON_H, is_group=False,
    ),
    "AWS::ECS::Cluster": dict(
        style=("shape=mxgraph.aws4.resourceIcon;resIcon=mxgraph.aws4.ecs;"
               "fillColor=#ED7100;gradientColor=none;strokeColor=none;"
               "fontFamily=Helvetica;fontSize=11;"
               "verticalLabelPosition=bottom;verticalAlign=top;align=center;"),
        w=ICON_W, h=ICON_H, is_group=False,
    ),
    "AWS::ApiGateway::RestApi": dict(
        style=("shape=mxgraph.aws4.resourceIcon;resIcon=mxgraph.aws4.api_gateway;"
               "fillColor=#8C4FFF;gradientColor=none;strokeColor=none;"
               "fontFamily=Helvetica;fontSize=11;"
               "verticalLabelPosition=bottom;verticalAlign=top;align=center;"),
        w=ICON_W, h=ICON_H, is_group=False,
    ),
    "AWS::ApiGateway": dict(  # fallback
        style=("shape=mxgraph.aws4.resourceIcon;resIcon=mxgraph.aws4.api_gateway;"
               "fillColor=#8C4FFF;gradientColor=none;strokeColor=none;"
               "fontFamily=Helvetica;fontSize=11;"
               "verticalLabelPosition=bottom;verticalAlign=top;align=center;"),
        w=ICON_W, h=ICON_H, is_group=False,
    ),
    "AWS::Diagram::Resource": dict(
        style=("shape=mxgraph.aws4.user;"
               "fillColor=#000000;gradientColor=none;strokeColor=none;"
               "fontFamily=Helvetica;fontSize=11;"
               "verticalLabelPosition=bottom;verticalAlign=top;align=center;"),
        w=58, h=58, is_group=False,
    ),
    # ── AWS groups (rendered as background rectangles) ───────────────────────
    "AWS::EC2::VPC": dict(
        style=("shape=mxgraph.aws4.group;grIcon=mxgraph.aws4.group_vpc;"
               "grStroke=0;fillColor=#E5F5F8;strokeColor=#147EBA;"
               "fontFamily=Helvetica;fontSize=12;fontStyle=1;"
               "verticalLabelPosition=top;verticalAlign=bottom;align=left;"),
        w=400, h=300, is_group=True,
    ),
    "AWS::EC2::Subnet": dict(
        style=("shape=mxgraph.aws4.group;grIcon=mxgraph.aws4.group_public_subnet;"
               "grStroke=0;fillColor=#E9F3E6;strokeColor=#67AB9F;"
               "fontFamily=Helvetica;fontSize=11;"
               "verticalLabelPosition=top;verticalAlign=bottom;align=left;"),
        w=160, h=160, is_group=True,
    ),
    "AWS::Diagram::Cloud": dict(
        style=("shape=mxgraph.aws4.group;grIcon=mxgraph.aws4.group_aws_cloud_alt;"
               "grStroke=0;fillColor=#FFFFFF;strokeColor=#232F3E;"
               "fontFamily=Helvetica;fontSize=13;fontStyle=1;"
               "verticalLabelPosition=top;verticalAlign=bottom;align=left;"),
        w=500, h=400, is_group=True,
    ),
    "AWS::Diagram::Canvas": dict(
        style="fillColor=none;strokeColor=none;",
        w=900, h=700, is_group=True,
    ),
    "AWS::Diagram::HorizontalStack": dict(
        style="fillColor=none;strokeColor=none;",
        w=200, h=100, is_group=True,
    ),
    "AWS::Diagram::VerticalStack": dict(
        style="fillColor=none;strokeColor=none;",
        w=120, h=200, is_group=True,
    ),
    # ── Unknown type fallback ────────────────────────────────────────────────
    "_default": dict(
        style=("rounded=1;whiteSpace=wrap;html=1;"
               "fillColor=#dae8fc;strokeColor=#6c8ebf;"
               "fontFamily=Helvetica;fontSize=11;"),
        w=120, h=60, is_group=False,
    ),
}

# Total visual icon height (includes label below).
def icon_visual_h(h: int) -> int:
    return h + LABEL_H


def get_style(rtype: str) -> dict:
    return STYLES.get(rtype, STYLES["_default"])


# --- Layout engine ---

def is_horizontal_dir(rtype: str, direction: str) -> bool:
    if rtype == "AWS::Diagram::HorizontalStack":
        return True
    if rtype == "AWS::Diagram::VerticalStack":
        return False
    return direction != "vertical"


def compute_size(name: str, resources: dict, _cache: dict = None) -> tuple:
    """
    Returns (w, h) for the resource bounding box including all children.
    Uses cache to avoid recomputing large graphs.
    """
    if _cache is None:
        _cache = {}
    if name in _cache:
        return _cache[name]

    res      = resources.get(name, {})
    rtype    = res.get("Type", "")
    children = [c for c in res.get("Children", []) if c in resources]
    direction = res.get("Direction", "horizontal")
    si        = get_style(rtype)
    horiz     = is_horizontal_dir(rtype, direction)

    if not children:
        # For leaf icons, visual height includes the label below.
        if not si["is_group"]:
            result = (si["w"], icon_visual_h(si["h"]))
        else:
            result = (si["w"], si["h"])
        _cache[name] = result
        return result

    child_sizes = [compute_size(c, resources, _cache) for c in children]

    if horiz:
        content_w = sum(cw for cw, _ in child_sizes) + GAP * (len(child_sizes) - 1)
        content_h = max(ch for _, ch in child_sizes)
    else:
        content_w = max(cw for cw, _ in child_sizes)
        content_h = sum(ch for _, ch in child_sizes) + GAP * (len(child_sizes) - 1)

    total_w = content_w + PADDING * 2
    total_h = content_h + PADDING * 2 + HEADER_H

    result = (max(total_w, si["w"]), max(total_h, si["h"]))
    _cache[name] = result
    return result


def compute_positions(name: str, resources: dict,
                      ox: int, oy: int,
                      size_cache: dict) -> dict:
    """
    Computes ABSOLUTE positions for all nodes recursively.
    Returns {name: (abs_x, abs_y, w, h)}.
    """
    w, h = size_cache[name]
    positions = {name: (ox, oy, w, h)}

    res      = resources.get(name, {})
    rtype    = res.get("Type", "")
    children = [c for c in res.get("Children", []) if c in resources]
    direction = res.get("Direction", "horizontal")
    si        = get_style(rtype)
    horiz     = is_horizontal_dir(rtype, direction)

    if not children:
        return positions

    # Inner area (excluding padding and header)
    cursor_x = ox + PADDING
    cursor_y = oy + HEADER_H + PADDING

    for child_name in children:
        cw, ch = size_cache[child_name]
        child_positions = compute_positions(child_name, resources,
                                            cursor_x, cursor_y, size_cache)
        positions.update(child_positions)
        if horiz:
            cursor_x += cw + GAP
        else:
            cursor_y += ch + GAP

    return positions


# --- draw.io XML generation ---

def bfs_nodes(resources: dict) -> list:
    """BFS from Canvas to order nodes (parents before children)."""
    order, visited = [], set()
    queue = deque(["Canvas"])
    while queue:
        n = queue.popleft()
        if n in visited or n not in resources:
            continue
        visited.add(n)
        order.append(n)
        for child in resources[n].get("Children", []):
            queue.append(child)
    # Nodes outside the main tree (e.g. isolated BorderChildren)
    for name in resources:
        if name not in visited:
            order.append(name)
    return order


def build_drawio(dac_file: str) -> str:
    with open(dac_file, encoding="utf-8") as f:
        doc = yaml.safe_load(f)

    diagram   = doc.get("Diagram", {})
    resources = diagram.get("Resources", {})
    links     = diagram.get("Links", [])

    # ── 1. Compute sizes ─────────────────────────────────────────────────────
    size_cache: dict = {}
    for name in resources:
        compute_size(name, resources, size_cache)

    # ── 2. Compute absolute positions ────────────────────────────────────────
    canvas_w, canvas_h = size_cache.get("Canvas", (900, 700))
    positions = compute_positions("Canvas", resources, 20, 20, size_cache)

    # ── 3. Split groups and icons (groups first in z-order) ─────────────────
    node_order = bfs_nodes(resources)
    groups = [n for n in node_order if get_style(resources.get(n, {}).get("Type", ""))["is_group"]]
    icons  = [n for n in node_order if not get_style(resources.get(n, {}).get("Type", ""))["is_group"]]
    draw_order = groups + icons  # groups in the back, icons in front

    # ── 4. Build XML ─────────────────────────────────────────────────────────
    model = ET.Element("mxGraphModel",
                       dx="1422", dy="762",
                       grid="1", gridSize="10",
                       guides="1", tooltips="1",
                       connect="1", arrows="1",
                       fold="1", page="1",
                       pageScale="1",
                       pageWidth="1654",
                       pageHeight="1169",
                       math="0", shadow="0")
    xml_root = ET.SubElement(model, "root")
    ET.SubElement(xml_root, "mxCell", id="0")
    ET.SubElement(xml_root, "mxCell", id="1", parent="0")

    name_to_id: dict = {}
    cell_id = 2

    for name in draw_order:
        if name not in positions:
            continue

        res   = resources.get(name, {})
        rtype = res.get("Type", "")
        # Title: replace line breaks with XML entity.
        title = res.get("Title", name).replace("\n", "&#xa;")
        x, y, w, h = positions[name]
        si = get_style(rtype)

        # ── All nodes are direct children of layer "1" (absolute positions) ─
        attrs = dict(
            id=str(cell_id),
            value=title,
            style=si["style"],
            vertex="1",
            parent="1",
        )
        cell = ET.SubElement(xml_root, "mxCell", **attrs)
        ET.SubElement(cell, "mxGeometry",
                      x=str(x), y=str(y),
                      width=str(w), height=str(h),
                      **{"as": "geometry"})

        name_to_id[name] = str(cell_id)
        cell_id += 1

    # ── 5. Edges ─────────────────────────────────────────────────────────────
    for link in links:
        src_id = name_to_id.get(link.get("Source", ""))
        tgt_id = name_to_id.get(link.get("Target", ""))
        if not src_id or not tgt_id:
            continue

        link_type = link.get("Type", "")
        if link_type == "orthogonal":
            edge_style = ("edgeStyle=orthogonalEdgeStyle;rounded=0;"
                          "orthogonalLoop=1;jettySize=auto;html=1;")
        else:
            edge_style = ("edgeStyle=orthogonalEdgeStyle;rounded=0;"
                          "orthogonalLoop=1;jettySize=auto;html=1;"
                          "exitX=0.5;exitY=1;exitDx=0;exitDy=0;"
                          "entryX=0.5;entryY=0;entryDx=0;entryDy=0;")

        # Arrowhead on the target side
        has_target_arrow = bool(link.get("TargetArrowHead"))
        if not has_target_arrow:
            edge_style += "endArrow=none;"

        line_color = link.get("LineColor", "")
        if line_color:
            # rgba(r,g,b,a) -> #rrggbb for draw.io
            try:
                parts = line_color.replace("rgba(", "").replace(")", "").split(",")
                r, g, b = int(parts[0]), int(parts[1]), int(parts[2])
                edge_style += f"strokeColor=#{r:02X}{g:02X}{b:02X};"
            except Exception:
                pass

        line_style = link.get("LineStyle", "")
        if line_style == "dashed":
            edge_style += "dashed=1;"

        edge = ET.SubElement(xml_root, "mxCell",
                             id=str(cell_id),
                             value="",
                             style=edge_style,
                             edge="1",
                             source=src_id,
                             target=tgt_id,
                             parent="1")
        ET.SubElement(edge, "mxGeometry", relative="1", **{"as": "geometry"})
        cell_id += 1

    # ── 6. Serialize ─────────────────────────────────────────────────────────
    raw    = ET.tostring(model, encoding="unicode")
    pretty = minidom.parseString(raw).toprettyxml(indent="  ")
    lines  = pretty.splitlines()
    if lines and lines[0].startswith("<?xml"):
        lines = lines[1:]
    return "\n".join(lines)


# --- CLI entrypoint ---

def main():
    parser = argparse.ArgumentParser(
        description="Convert DAC YAML files to draw.io (.drawio)")
    parser.add_argument("input",  help="Input DAC YAML file")
    parser.add_argument("-o", "--output",
                        help="Output .drawio file (default: input name with .drawio)")
    args = parser.parse_args()

    if not os.path.exists(args.input):
        print(f"Error: file '{args.input}' was not found.", file=sys.stderr)
        sys.exit(1)

    output = args.output or os.path.splitext(args.input)[0] + ".drawio"
    xml_content = build_drawio(args.input)

    with open(output, "w", encoding="utf-8") as f:
        f.write(xml_content)
    print(f"[OK] draw.io generated: {output}")


if __name__ == "__main__":
    main()
