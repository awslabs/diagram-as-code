#!/usr/bin/env python3
"""
dac-to-drawio.py
Converte arquivos DAC (diagram-as-code YAML) para o formato draw.io (.drawio / XML mxGraph).

Estratégia de layout:
  - Todos os elementos são posicionados com coordenadas ABSOLUTAS (parent="1")
  - Grupos ficam como retângulos de fundo (z-order: grupos primeiro, ícones depois)
  - Isso garante que o layout fique visualmente idêntico ao PNG gerado pelo awsdac
  - Evita o problema de acumulação de erros com coordenadas relativas aninhadas

Uso:
    python3 tools/dac-to-drawio.py <arquivo.yaml> [-o saida.drawio]
"""

import yaml
import sys
import os
import argparse
import xml.etree.ElementTree as ET
from xml.dom import minidom
from collections import deque

# ─── Constantes de layout ────────────────────────────────────────────────────
PADDING   = 36   # Espaço interno em grupos (entre borda e filhos)
GAP       = 16   # Espaço entre irmãos
HEADER_H  = 48   # Altura do cabeçalho do grupo (espaço para ícone + rótulo)
ICON_W    = 78   # Largura padrão de ícones folha
ICON_H    = 78   # Altura padrão de ícones folha
LABEL_H   = 24   # Altura extra para rótulo abaixo do ícone

# ─── Estilos draw.io por tipo de recurso DAC ─────────────────────────────────
# is_group=True → renderizado como retângulo de fundo (grupo AWS)
# is_group=False → renderizado como ícone de recurso
STYLES: dict = {
    # ── Ícones de recursos AWS ────────────────────────────────────────────────
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
    "AWS::CloudFront": dict(  # fallback de serviço
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
    # ── Grupos AWS (renderizados como retângulos de fundo) ────────────────────
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
    # ── Tipo desconhecido ─────────────────────────────────────────────────────
    "_default": dict(
        style=("rounded=1;whiteSpace=wrap;html=1;"
               "fillColor=#dae8fc;strokeColor=#6c8ebf;"
               "fontFamily=Helvetica;fontSize=11;"),
        w=120, h=60, is_group=False,
    ),
}

# Tamanho visual total de um ícone (inclui rótulo abaixo)
def icon_visual_h(h: int) -> int:
    return h + LABEL_H


def get_style(rtype: str) -> dict:
    return STYLES.get(rtype, STYLES["_default"])


# ─── Motor de layout ──────────────────────────────────────────────────────────

def is_horizontal_dir(rtype: str, direction: str) -> bool:
    if rtype == "AWS::Diagram::HorizontalStack":
        return True
    if rtype == "AWS::Diagram::VerticalStack":
        return False
    return direction != "vertical"


def compute_size(name: str, resources: dict, _cache: dict = None) -> tuple:
    """
    Retorna (w, h) do bounding box de um recurso incluindo todos os filhos.
    Usa cache para evitar recálculo em grafos grandes.
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
        # Ícones folha: o "h" visual inclui o rótulo abaixo
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
    Calcula posições ABSOLUTAS de todos os nós recursivamente.
    Retorna {name: (abs_x, abs_y, w, h)}.
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

    # Área interior (excluindo padding e cabeçalho)
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


# ─── Geração do XML draw.io ───────────────────────────────────────────────────

def bfs_nodes(resources: dict) -> list:
    """BFS a partir do Canvas para ordenar nós (pais antes de filhos)."""
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
    # Nós fora da árvore (ex: BorderChildren isolados)
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

    # ── 1. Calcular tamanhos ──────────────────────────────────────────────────
    size_cache: dict = {}
    for name in resources:
        compute_size(name, resources, size_cache)

    # ── 2. Calcular posições absolutas ────────────────────────────────────────
    canvas_w, canvas_h = size_cache.get("Canvas", (900, 700))
    positions = compute_positions("Canvas", resources, 20, 20, size_cache)

    # ── 3. Separar grupos de ícones (grupos primeiro no z-order) ─────────────
    node_order = bfs_nodes(resources)
    groups = [n for n in node_order if get_style(resources.get(n, {}).get("Type", ""))["is_group"]]
    icons  = [n for n in node_order if not get_style(resources.get(n, {}).get("Type", ""))["is_group"]]
    draw_order = groups + icons  # grupos no fundo, ícones na frente

    # ── 4. Montar XML ─────────────────────────────────────────────────────────
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
        # Título: substitui quebras de linha por entidade XML
        title = res.get("Title", name).replace("\n", "&#xa;")
        x, y, w, h = positions[name]
        si = get_style(rtype)

        # ── Todos os nós são filhos diretos do layer "1" (posição absoluta) ──
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

    # ── 5. Arestas ────────────────────────────────────────────────────────────
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

        # Seta na ponta do destino
        has_target_arrow = bool(link.get("TargetArrowHead"))
        if not has_target_arrow:
            edge_style += "endArrow=none;"

        line_color = link.get("LineColor", "")
        if line_color:
            # rgba(r,g,b,a) → #rrggbb para draw.io
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

    # ── 6. Serializar ─────────────────────────────────────────────────────────
    raw    = ET.tostring(model, encoding="unicode")
    pretty = minidom.parseString(raw).toprettyxml(indent="  ")
    lines  = pretty.splitlines()
    if lines and lines[0].startswith("<?xml"):
        lines = lines[1:]
    return "\n".join(lines)


# ─── Entrada ─────────────────────────────────────────────────────────────────

def main():
    parser = argparse.ArgumentParser(
        description="Converte arquivos DAC YAML para draw.io (.drawio)")
    parser.add_argument("input",  help="Arquivo DAC YAML de entrada")
    parser.add_argument("-o", "--output",
                        help="Arquivo de saída .drawio (padrão: mesmo nome do input)")
    args = parser.parse_args()

    if not os.path.exists(args.input):
        print(f"Erro: arquivo '{args.input}' não encontrado.", file=sys.stderr)
        sys.exit(1)

    output = args.output or os.path.splitext(args.input)[0] + ".drawio"
    xml_content = build_drawio(args.input)

    with open(output, "w", encoding="utf-8") as f:
        f.write(xml_content)
    print(f"[OK] draw.io gerado: {output}")


if __name__ == "__main__":
    main()
