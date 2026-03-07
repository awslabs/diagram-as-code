// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package ctl

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/awslabs/diagram-as-code/internal/definition"
	"github.com/awslabs/diagram-as-code/internal/types"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// drawioScale converts awsdac pixel coordinates to draw.io units.
// The rendering engine uses pixels; 0.5 produces comfortable on-screen sizing.
const drawioScale = 0.5

// --- DAC type to draw.io style mapping ---

type drawioNode struct {
	style   string
	isGroup bool // true = background rectangle (group), false = resource icon
}

// stencilStyle returns the style for an icon using native draw.io stencils.
// shapeName is the exact aws4 stencil name (e.g. "mxgraph.aws4.lambda").
// fillColor is the AWS category color (e.g. "#ED7100" for compute).
func stencilStyle(shapeName, fillColor string) string {
	return fmt.Sprintf(
		"outlineConnect=0;fontColor=#232F3E;gradientColor=none;"+
			"strokeColor=none;fillColor=%s;"+
			"labelBackgroundColor=#ffffff;align=center;html=1;"+
			"fontSize=11;fontStyle=0;aspect=fixed;pointerEvents=1;"+
			"shape=%s;verticalLabelPosition=bottom;verticalAlign=top;",
		fillColor, shapeName,
	)
}

// groupStyle returns the style for an AWS4 group with a visible border.
func groupStyle(grIcon, fillColor, strokeColor string) string {
	return fmt.Sprintf(
		"shape=mxgraph.aws4.group;grIcon=%s;grStroke=0;"+
			"fillColor=%s;strokeColor=%s;strokeWidth=2;"+
			"fontFamily=Helvetica;fontSize=12;"+
			"verticalLabelPosition=top;verticalAlign=bottom;align=left;",
		grIcon, fillColor, strokeColor,
	)
}

// Official AWS category colors.
const (
	awsColorCompute    = "#ED7100" // Compute (EC2, Lambda, ECS, Fargate)
	awsColorStorage    = "#3F8624" // Storage (S3, EBS)
	awsColorDatabase   = "#C925D1" // Database (RDS, DynamoDB, ElastiCache)
	awsColorNetworking = "#8C4FFF" // Networking & CDN (VPC, ALB, CloudFront, API GW)
	awsColorSecurity   = "#DD344C" // Security (IAM, etc.)
	awsColorGeneral    = "#232F3E" // General / Management
)

var dacTypeStyles = map[string]drawioNode{
	// ── Compute ────────────────────────────────────────────────────────────
	"AWS::EC2::Instance": {style: stencilStyle("mxgraph.aws4.ec2", awsColorCompute)},
	"AWS::Lambda::Function": {style: stencilStyle("mxgraph.aws4.lambda", awsColorCompute)},
	"AWS::ECS::Cluster":     {style: stencilStyle("mxgraph.aws4.ecs", awsColorCompute)},
	"AWS::ECS::Service":     {style: stencilStyle("mxgraph.aws4.ecs_service", awsColorCompute)},
	"AWS::ECS::TaskDefinition": {style: stencilStyle("mxgraph.aws4.ecs_task", awsColorCompute)},

	// ── Storage ────────────────────────────────────────────────────────────
	"AWS::S3::Bucket": {style: stencilStyle("mxgraph.aws4.s3", awsColorStorage)},

	// ── Database ───────────────────────────────────────────────────────────
	"AWS::RDS::DBInstance":  {style: stencilStyle("mxgraph.aws4.rds", awsColorDatabase)},
	"AWS::RDS::DBCluster":   {style: stencilStyle("mxgraph.aws4.rds", awsColorDatabase)},
	"AWS::DynamoDB::Table":  {style: stencilStyle("mxgraph.aws4.dynamodb", awsColorDatabase)},
	"AWS::ElastiCache::CacheCluster": {style: stencilStyle("mxgraph.aws4.elasticache", awsColorDatabase)},

	// ── Networking & CDN ───────────────────────────────────────────────────
	"AWS::EC2::InternetGateway":              {style: stencilStyle("mxgraph.aws4.internet_gateway", awsColorNetworking)},
	"AWS::ElasticLoadBalancingV2::LoadBalancer": {style: stencilStyle("mxgraph.aws4.application_load_balancer", awsColorNetworking)},
	"AWS::ElasticLoadBalancing::LoadBalancer":   {style: stencilStyle("mxgraph.aws4.classic_load_balancer", awsColorNetworking)},
	"AWS::CloudFront::Distribution":          {style: stencilStyle("mxgraph.aws4.cloudfront", awsColorNetworking)},
	"AWS::CloudFront":                        {style: stencilStyle("mxgraph.aws4.cloudfront", awsColorNetworking)},
	"AWS::ApiGateway::RestApi":               {style: stencilStyle("mxgraph.aws4.api_gateway", awsColorNetworking)},
	"AWS::ApiGateway":                        {style: stencilStyle("mxgraph.aws4.api_gateway", awsColorNetworking)},
	"AWS::EC2::NatGateway":                   {style: stencilStyle("mxgraph.aws4.nat_gateway", awsColorNetworking)},
	"AWS::Route53::HostedZone":               {style: stencilStyle("mxgraph.aws4.route_53", awsColorNetworking)},
	"AWS::SNS::Topic":                        {style: stencilStyle("mxgraph.aws4.sns", "#E7157B")},
	"AWS::SQS::Queue":                        {style: stencilStyle("mxgraph.aws4.sqs", "#E7157B")},

	// ── Grupos / containers ────────────────────────────────────────────────
	"AWS::EC2::VPC": {
		style:   groupStyle("mxgraph.aws4.group_vpc", "#E5F5F8", "#147EBA"),
		isGroup: true,
	},
	"AWS::EC2::Subnet": {
		style:   groupStyle("mxgraph.aws4.group_public_subnet", "#E9F3E6", "#67AB9F"),
		isGroup: true,
	},
	"AWS::Diagram::Cloud": {
		style:   groupStyle("mxgraph.aws4.group_aws_cloud_alt", "#FFFFFF", "#232F3E"),
		isGroup: true,
	},
	"AWS::Diagram::Canvas": {
		style:   "fillColor=none;strokeColor=none;",
		isGroup: true,
	},
	"AWS::Diagram::HorizontalStack": {
		style:   "fillColor=none;strokeColor=none;",
		isGroup: true,
	},
	"AWS::Diagram::VerticalStack": {
		style:   "fillColor=none;strokeColor=none;",
		isGroup: true,
	},
	// ── Generic resource (User, etc.) ──────────────────────────────────────
	"AWS::Diagram::Resource": {
		style: stencilStyle("mxgraph.aws4.user", awsColorGeneral),
	},
}

func getDrawioNode(rtype string) drawioNode {
	if n, ok := dacTypeStyles[rtype]; ok {
		return n
	}
	// Fallback: try service-level icon (e.g. AWS::ApiGateway::X -> AWS::ApiGateway).
	parts := strings.SplitN(rtype, "::", 3)
	if len(parts) >= 2 {
		svcType := strings.Join(parts[:2], "::")
		if n, ok := dacTypeStyles[svcType]; ok {
			return n
		}
	}
	return drawioNode{
		style:   "rounded=1;whiteSpace=wrap;html=1;fillColor=#dae8fc;strokeColor=#6c8ebf;fontFamily=Helvetica;",
		isGroup: false,
	}
}

// --- XML generation ---

// rgbaToHex converts "rgba(r,g,b,a)" to "#RRGGBB" for draw.io.
func rgbaToHex(rgba string) string {
	var r, g, b, a uint8
	if _, err := fmt.Sscanf(rgba, "rgba(%d,%d,%d,%d)", &r, &g, &b, &a); err == nil {
		return fmt.Sprintf("#%02X%02X%02X", r, g, b)
	}
	return rgba // Keep original value if parsing fails.
}

func px(v int) string {
	return fmt.Sprintf("%.0f", float64(v)*drawioScale)
}

// escapeXML applies basic escaping for XML attribute values.
func escapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "\n", "&#xa;")
	return s
}

func exportToDrawio(
	template *TemplateStruct,
	resources map[string]*types.Resource,
	outputfile string,
) error {
	var buf bytes.Buffer

	// ── BFS insertion order (parents before children) ─────────────────────
	order := []string{}
	visited := map[string]bool{}
	queue := []string{"Canvas"}
	for len(queue) > 0 {
		n := queue[0]
		queue = queue[1:]
		if visited[n] || resources[n] == nil {
			continue
		}
		visited[n] = true
		order = append(order, n)
		if res, ok := template.Resources[n]; ok {
			for _, child := range res.Children {
				queue = append(queue, child)
			}
		}
	}
	// Resources outside the tree (e.g. isolated BorderChildren)
	for name := range template.Resources {
		if !visited[name] && resources[name] != nil {
			order = append(order, name)
		}
	}

	// ── Split groups and icons for z-order (groups first) ─────────────────
	var groups, icons []string
	for _, name := range order {
		res := template.Resources[name]
		node := getDrawioNode(res.Type)
		hasChildren := len(res.Children) > 0
		if node.isGroup || hasChildren {
			groups = append(groups, name)
		} else {
			icons = append(icons, name)
		}
	}
	drawOrder := append(groups, icons...)

	// ── ID map ────────────────────────────────────────────────────────────
	nameToID := map[string]string{}
	cellID := 2
	for _, name := range drawOrder {
		nameToID[name] = fmt.Sprintf("%d", cellID)
		cellID++
	}

	// ── XML header ─────────────────────────────────────────────────────────
	buf.WriteString(`<mxGraphModel dx="1422" dy="762" grid="1" gridSize="10" guides="1" tooltips="1" connect="1" arrows="1" fold="1" page="1" pageScale="1" pageWidth="1654" pageHeight="1169" math="0" shadow="0">`)
	buf.WriteString("\n  <root>\n")
	buf.WriteString("    <mxCell id=\"0\"/>\n")
	buf.WriteString("    <mxCell id=\"1\" parent=\"0\"/>\n")

	// ── Nodes ──────────────────────────────────────────────────────────────
	for _, name := range drawOrder {
		res := template.Resources[name]
		runtime := resources[name]
		if runtime == nil {
			continue
		}

		bounds := runtime.GetBindings()
		if bounds.Empty() {
			continue
		}

		node := getDrawioNode(res.Type)
		hasChildren := len(res.Children) > 0

		label := res.Title
		if label == "" {
			label = name
		}

		// Private subnet preset: swap group icon.
		if res.Type == "AWS::EC2::Subnet" && res.Preset == "PrivateSubnet" {
			node.style = groupStyle("mxgraph.aws4.group_private_subnet", "#F4E6FA", "#AD688E")
		}

		style := node.style
		value := label

		// For leaf icons: use official AWS SVG from the Asset Package.
		if !node.isGroup && !hasChildren {
			if dataURI := GetAWSIconDataURI(res.Type); dataURI != "" {
				style = fmt.Sprintf(
					"shape=image;html=1;aspect=fixed;"+
						"verticalLabelPosition=bottom;verticalAlign=top;align=center;"+
						"fontSize=11;fontColor=#232F3E;labelBackgroundColor=#ffffff;"+
						"image=%s;",
					dataURI,
				)
			}
		}

		// Resources with children use YAML fill/border when present, otherwise defaults.
		if !node.isGroup && hasChildren {
			fillColor := "#dae8fc"
			strokeColor := "#6c8ebf"
			tmplRes := template.Resources[name]
			if tmplRes.FillColor != "" {
				fillColor = rgbaToHex(tmplRes.FillColor)
			}
			if tmplRes.BorderColor != "" {
				strokeColor = rgbaToHex(tmplRes.BorderColor)
			}
			style = fmt.Sprintf(
				"rounded=1;whiteSpace=wrap;html=1;"+
					"fillColor=%s;strokeColor=%s;strokeWidth=2;"+
					"fontFamily=Helvetica;fontSize=12;"+
					"verticalAlign=top;align=left;",
				fillColor, strokeColor,
			)
		}

		x := px(bounds.Min.X)
		y := px(bounds.Min.Y)
		w := px(bounds.Dx())
		h := px(bounds.Dy())

		id := nameToID[name]

		fmt.Fprintf(&buf,
			"    <mxCell id=\"%s\" value=\"%s\" style=\"%s\" vertex=\"1\" parent=\"1\">\n"+
				"      <mxGeometry x=\"%s\" y=\"%s\" width=\"%s\" height=\"%s\" as=\"geometry\"/>\n"+
				"    </mxCell>\n",
			id, escapeXML(value), escapeXML(style),
			x, y, w, h,
		)
	}

	// ── Edges ──────────────────────────────────────────────────────────────
	for _, link := range template.Links {
		srcID, srcOK := nameToID[link.Source]
		tgtID, tgtOK := nameToID[link.Target]
		if !srcOK || !tgtOK {
			log.Warnf("drawio: link %s→%s: resource not found, skipping", link.Source, link.Target)
			continue
		}

		edgeStyle := "edgeStyle=orthogonalEdgeStyle;rounded=0;orthogonalLoop=1;jettySize=auto;html=1;"

		if link.LineStyle == "dashed" {
			edgeStyle += "dashed=1;"
		}
		if link.LineColor != "" {
			// rgba(r,g,b,a) -> #RRGGBB
			var r, g, b, a uint8
			if _, err := fmt.Sscanf(link.LineColor, "rgba(%d,%d,%d,%d)", &r, &g, &b, &a); err == nil {
				edgeStyle += fmt.Sprintf("strokeColor=#%02X%02X%02X;", r, g, b)
			}
		}
		if link.TargetArrowHead.Type == "" {
			edgeStyle += "endArrow=none;"
		} else {
			edgeStyle += "endArrow=open;endFill=0;"
		}
		if link.SourceArrowHead.Type != "" {
			edgeStyle += "startArrow=open;startFill=0;"
		}

		fmt.Fprintf(&buf,
			"    <mxCell id=\"%d\" value=\"\" style=\"%s\" edge=\"1\" source=\"%s\" target=\"%s\" parent=\"1\">\n"+
				"      <mxGeometry relative=\"1\" as=\"geometry\"/>\n"+
				"    </mxCell>\n",
			cellID, escapeXML(edgeStyle), srcID, tgtID,
		)
		cellID++
	}

	buf.WriteString("  </root>\n</mxGraphModel>\n")

	return os.WriteFile(outputfile, buf.Bytes(), 0600)
}

// --- Public entry point ---

// CreateDrawioFromDacFile reads a DAC YAML file and produces a .drawio file
// with the same pixel-accurate layout as PNG output (same awsdac layout engine).
func CreateDrawioFromDacFile(inputfile string, outputfile *string, opts *CreateOptions) error {
	log.Infof("drawio: input file: %s", inputfile)

	data, err := getTemplate(inputfile)
	if err != nil {
		return fmt.Errorf("failed to get template: %w", err)
	}

	var processedData []byte
	if opts.IsGoTemplate {
		processedData, err = processTemplate(data)
		if err != nil {
			return fmt.Errorf("failed to process template: %w", err)
		}
	} else {
		processedData = data
	}

	var template TemplateStruct
	dec := yaml.NewDecoder(bytes.NewReader(processedData))
	dec.KnownFields(true)
	if err := dec.Decode(&template); err != nil {
		return fmt.Errorf("failed to decode YAML: %w", err)
	}

	var ds definition.DefinitionStructure
	if opts.OverrideDefFile != "" {
		var overrideDefTemplate TemplateStruct
		defFile := DefinitionFile{}
		if IsURL(opts.OverrideDefFile) {
			defFile.Type = "URL"
			defFile.Url = opts.OverrideDefFile
		} else {
			defFile.Type = "LocalFile"
			defFile.LocalFile = opts.OverrideDefFile
		}
		overrideDefTemplate.DefinitionFiles = append(overrideDefTemplate.DefinitionFiles, defFile)
		if err := loadDefinitionFiles(&overrideDefTemplate, &ds, true); err != nil {
			return fmt.Errorf("failed to load override definition files: %w", err)
		}
	} else {
		if err := loadDefinitionFiles(&template, &ds, opts.AllowUntrustedDefinitions); err != nil {
			return fmt.Errorf("failed to load definition files: %w", err)
		}
	}

	resources := make(map[string]*types.Resource)
	if err := loadResources(&template, ds, resources); err != nil {
		return fmt.Errorf("failed to load resources: %w", err)
	}
	if err := associateChildren(&template, resources); err != nil {
		return fmt.Errorf("failed to associate children: %w", err)
	}
	if err := loadLinks(&template, resources); err != nil {
		return fmt.Errorf("failed to load links: %w", err)
	}

	// Reorder children (same behavior as PNG generation).
	canvas, exists := resources["Canvas"]
	if !exists {
		return fmt.Errorf("Canvas resource not found")
	}
	var allLinks []*types.Link
	for _, r := range resources {
		allLinks = append(allLinks, r.GetLinks()...)
	}
	types.ReorderChildrenByLinks(canvas, allLinks)

	// ── Layout: use the same engine as PNG generation ─────────────────────
	if err := canvas.Scale(nil, nil); err != nil {
		return fmt.Errorf("error scaling diagram: %w", err)
	}
	if err := canvas.ZeroAdjust(); err != nil {
		return fmt.Errorf("error adjusting diagram: %w", err)
	}

	// Resolve auto-positions (required for links, but does not change bounds).
	for _, r := range resources {
		for _, link := range r.GetLinks() {
			if err := link.ResolveAutoPositions(); err != nil {
				return fmt.Errorf("failed to resolve auto-positions: %w", err)
			}
		}
	}

	// ── Export draw.io ─────────────────────────────────────────────────────
	if err := exportToDrawio(&template, resources, *outputfile); err != nil {
		return fmt.Errorf("failed to export drawio: %w", err)
	}

	return nil
}
