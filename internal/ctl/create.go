// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package ctl

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"strings"

	"github.com/awslabs/diagram-as-code/internal/cache"
	"github.com/awslabs/diagram-as-code/internal/definition"
	"github.com/awslabs/diagram-as-code/internal/types"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/maps"
)

func stringToColor(c string) color.RGBA {
	var r, g, b, a uint8
	_, err := fmt.Sscanf(c, "rgba(%d,%d,%d,%d)", &r, &g, &b, &a)
	if err != nil {
		log.Fatal(err)
	}
	return color.RGBA{r, g, b, a}
}

type TemplateStruct struct {
	Diagram `yaml:"Diagram"`
}

type Diagram struct {
	DefinitionFiles []DefinitionFile    `yaml:"DefinitionFiles"`
	Resources       map[string]Resource `yaml:"Resources"`
	Links           []Link              `yaml:"Links"`
}

type DefinitionFile struct {
	Type      string                         `yaml:"Type"` // URL,LocalFile,Embed
	Url       string                         `yaml:"Url"`
	LocalFile string                         `yaml:"LocalFile"`
	Embed     definition.DefinitionStructure `yaml:"Embed"`
}

type Resource struct {
	Type           string            `yaml:"Type"`
	Icon           string            `yaml:"Icon"`
	IconFill       *ResourceIconFill `yaml:"IconFill"`
	Direction      string            `yaml:"Direction"`
	Preset         string            `yaml:"Preset"`
	Align          string            `yaml:"Align"`
	HeaderAlign    string            `yaml:"HeaderAlign"`
	FillColor      string            `yaml:"FillColor"`
	Title          string            `yaml:"Title"`
	TitleColor     string            `yaml:"TitleColor"`
	Font           string            `yaml:"Font"`
	Children       []string          `yaml:"Children"`
	BorderColor    string            `yaml:"BorderColor"`
	BorderChildren []BorderChild     `yaml:"BorderChildren"`
}

type ResourceIconFill struct {
	Type  *string `yaml:"Type"`
	Color *string `yaml:"Color"`
}

type BorderChild struct {
	Position string `yaml:"Position"`
	Resource string `yaml:"Resource"`
}

type Link struct {
	Source          string          `yaml:"Source"`
	SourcePosition  string          `yaml:"SourcePosition"`
	SourceArrowHead types.ArrowHead `yaml:"SourceArrowHead"`
	Target          string          `yaml:"Target"`
	TargetPosition  string          `yaml:"TargetPosition"`
	TargetArrowHead types.ArrowHead `yaml:"TargetArrowHead"`
	Type            string          `yaml:"Type"`
	LineWidth       int             `yaml:"LineWidth"`
	LineColor       string          `yaml:"LineColor"`
	LineStyle       string          `yaml:"LineStyle"`
	Labels          LinkLabels      `yaml:"Labels"`
}

type LinkLabels struct {
	SourceRight *LinkLabel `yaml:"SourceRight"`
	SourceLeft  *LinkLabel `yaml:"SourceLeft"`
	TargetRight *LinkLabel `yaml:"TargetRight"`
	TargetLeft  *LinkLabel `yaml:"TargetLeft"`
}

type LinkLabel struct {
	Type  *string `yaml:"Type"`
	Title string  `yaml:"Title"`
	Color *string `yaml:"Color"`
	Font  *string `yaml:"Font"`
}

type CreateOptions struct {
	IsGoTemplate    bool
	OverrideDefFile string
}

func createDiagram(resources map[string]*types.Resource, outputfile *string) {

	log.Info("--- Draw diagram ---")
	err := resources["Canvas"].Scale(nil, nil)
	if err != nil {
		log.Fatalf("Error scaling diagram: %v", err)
	}
	resources["Canvas"].ZeroAdjust()
	img := resources["Canvas"].Draw(nil, nil)

	log.Infof("Save %s\n", *outputfile)
	fmt.Printf("[Completed] AWS infrastructure diagram generated: %s\n", *outputfile)
	f, _ := os.OpenFile(*outputfile, os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	png.Encode(f, img)
}

func loadDefinitionFiles(template *TemplateStruct, ds *definition.DefinitionStructure) {

	// Load definition files
	for _, v := range template.Diagram.DefinitionFiles {
		switch v.Type {
		case "URL":
			log.Infof("Fetch definition file from URL: %s\n", v.Url)
			cacheFilePath, err := cache.FetchFile(v.Url)
			if err != nil {
				log.Fatal(err)
			}
			log.Infof("Read definition file from cache file: %s\n", cacheFilePath)
			err = ds.LoadDefinitions(cacheFilePath)
			if err != nil {
				log.Fatal(err)
			}
		case "LocalFile":
			log.Infof("Read definition file from path: %s\n", v.LocalFile)
			err := ds.LoadDefinitions(v.LocalFile)
			if err != nil {
				log.Fatal(err)
			}
		case "Embed":
			log.Info("Read embedded definitions")
			maps.Copy(ds.Definitions, v.Embed.Definitions)
		}
	}

}

func loadResources(template *TemplateStruct, ds definition.DefinitionStructure, resources map[string]*types.Resource) {

	resources["Canvas"] = new(types.Resource).Init()

	for k, v := range template.Resources {
		// Override order: Definition{Resource Type -> Preset} -> Template
		log.Infof("Load Resource: %s (%s)\n", k, v.Type)
		switch v.Type {
		case "":
			log.Warnf("%s does not have Type field. Skipping this resource.", k)
			continue
		case "AWS::Diagram::Canvas":
			resources[k].SetBorderColor(color.RGBA{0, 0, 0, 0})
			resources[k].SetFillColor(color.RGBA{255, 255, 255, 255})
		case "AWS::Diagram::Resource":
			resources[k] = new(types.Resource).Init()
		case "AWS::Diagram::VerticalStack":
			resources[k] = new(types.VerticalStack).Init()
		case "AWS::Diagram::HorizontalStack":
			resources[k] = new(types.HorizontalStack).Init()
		default:
			def, ok := ds.Definitions[v.Type]
			if !ok {
				newType := fallbackToServiceIcon(v.Type)
				_, check := ds.Definitions[newType]
				if !check {
					log.Warnf("Type %s is not defined in the DAC definition file. It cannot be fall backed to service icon. Ignore this type.\n", v.Type)
					continue
				}
				log.Warnf("Type %s is not defined in the DAC definition file. It's fall backed to its service icon (Type %s).\n", v.Type, newType)
				def = ds.Definitions[newType]

				// Change the title to indicate the original resource type for fallback icons.
				if v.Title == "" {
					resources[k].SetLabel(&v.Type, nil, nil)
				}
			}
			if def.Type == "Resource" {
				resources[k] = new(types.Resource).Init()
			} else if def.Type == "Group" {
				resources[k] = new(types.Resource).Init()
			}
			if fill := def.Fill; fill != nil {
				resources[k].SetFillColor(stringToColor(fill.Color))
			}
			if border := def.Border; border != nil {
				resources[k].SetBorderColor(stringToColor(border.Color))
				switch border.Type {
				case "straight":
					resources[k].SetBorderType(types.BORDER_TYPE_STRAIGHT)
				case "dashed":
					resources[k].SetBorderType(types.BORDER_TYPE_DASHED)
				default:
					resources[k].SetBorderType(types.BORDER_TYPE_STRAIGHT)
				}
			}
			if label := def.Label; label != nil {
				if label.Title != "" {
					resources[k].SetLabel(&label.Title, nil, nil)
				}
				if label.Color != "" {
					c := stringToColor(label.Color)
					resources[k].SetLabel(nil, &c, nil)
				}
				if label.Font != "" {
					resources[k].SetLabel(nil, nil, &label.Font)
				}
			}
			if headerAlign := def.HeaderAlign; headerAlign != "" {
				resources[k].SetHeaderAlign(headerAlign)
			}
			if icon := def.Icon; icon != nil {
				if def.CacheFilePath == "" {
					break
				}
				err := resources[k].LoadIcon(def.CacheFilePath)
				if err != nil {
					panic(err)
				}
			}
		}

		switch v.Preset {
		case "BlankGroup":
			resources[k].SetIconBounds(image.Rect(0, 0, 64, 64))
			resources[k].SetBorderColor(color.RGBA{0, 0, 0, 0})
		case "":
		default:
			def, ok := ds.Definitions[v.Preset]
			if !ok {
				log.Warnf("Unknown preset %s on %s\n", v.Preset, v.Type)
			}
			if fill := def.Fill; fill != nil {
				resources[k].SetFillColor(stringToColor(fill.Color))
			}
			if border := def.Border; border != nil {
				resources[k].SetBorderColor(stringToColor(border.Color))
				switch border.Type {
				case "straight":
					resources[k].SetBorderType(types.BORDER_TYPE_STRAIGHT)
				case "dashed":
					resources[k].SetBorderType(types.BORDER_TYPE_DASHED)
				default:
					resources[k].SetBorderType(types.BORDER_TYPE_STRAIGHT)
				}
			}
			if label := def.Label; label != nil {
				if label.Title != "" {
					resources[k].SetLabel(&label.Title, nil, nil)
				}
				if label.Color != "" {
					c := stringToColor(label.Color)
					resources[k].SetLabel(nil, &c, nil)
				}
				if label.Font != "" {
					resources[k].SetLabel(nil, nil, &label.Font)
				}
			}
			if headerAlign := def.HeaderAlign; headerAlign != "" {
				resources[k].SetHeaderAlign(headerAlign)
			}
			if icon := def.Icon; icon != nil {
				if def.CacheFilePath != "" {
					err := resources[k].LoadIcon(def.CacheFilePath)
					if err != nil {
						panic(err)
					}
				}
			}
		}
		if v.Icon != "" {
			err := resources[k].LoadIcon(v.Icon)
			if err != nil {
				panic(err)
			}
		}
		if v.IconFill != nil {
			switch *v.IconFill.Type {
			case "none":
				resources[k].SetIconFill(types.ICON_FILL_TYPE_NONE, nil)
			case "rect":
				if v.IconFill.Color != nil {
					c := stringToColor(*v.IconFill.Color)
					resources[k].SetIconFill(types.ICON_FILL_TYPE_RECT, &c)
				} else {
					resources[k].SetIconFill(types.ICON_FILL_TYPE_RECT, nil)
				}
			default:
				resources[k].SetIconFill(types.ICON_FILL_TYPE_NONE, nil)
			}
		}
		if v.Title != "" {
			resources[k].SetLabel(&v.Title, nil, nil)
		}
		if v.TitleColor != "" {
			c := stringToColor(v.TitleColor)
			resources[k].SetLabel(nil, &c, nil)
		}
		if v.HeaderAlign != "" {
			resources[k].SetHeaderAlign(v.HeaderAlign)
		}
		if v.Font != "" {
			resources[k].SetLabel(nil, nil, &v.Font)
		}
		if v.Align != "" {
			resources[k].SetAlign(v.Align)
		}
		if v.Direction != "" {
			resources[k].SetDirection(v.Direction)
		}
		if v.FillColor != "" {
			resources[k].SetFillColor(stringToColor(v.FillColor))
		}
		if v.BorderColor != "" {
			resources[k].SetBorderColor(stringToColor(v.BorderColor))
		}
	}

}

func fallbackToServiceIcon(inputType string) string {

	parts := strings.SplitN(inputType, "::", 3)
	possibleServiceType := strings.Join(parts[:2], "::")

	return possibleServiceType
}

func associateChildren(template *TemplateStruct, resources map[string]*types.Resource) {

	for logicalId, v := range template.Resources {
		_, ok := resources[logicalId]
		if !ok {
			log.Fatalf("Unknown resource %s\n", logicalId)
		}
		for _, child := range v.Children {
			_, ok := resources[child]
			if ok {
				log.Infof("Add child(%s) on %s", child, logicalId)
				resources[logicalId].AddChild(resources[child])
			} else {
				log.Warnf("Child `%s` was not found, ignoring it.", child)
			}
		}
		for _, borderChild := range v.BorderChildren {
			_, ok := resources[borderChild.Resource]
			if !ok {
				log.Warnf("Child `%s` was not found, ignoring it.", borderChild.Resource)
				continue
			}
			log.Infof("Add BorderChild(%s) on %s", borderChild.Resource, logicalId)

			position, err := types.ConvertWindrose(borderChild.Position)
			if err != nil {
				panic(err)
			}
			bc := types.BorderChild{
				Position: position,
				Resource: resources[borderChild.Resource],
			}
			resources[logicalId].AddBorderChild(&bc)
		}
	}
}

func convertLabel(label *LinkLabel) *types.LinkLabel {
	r := &types.LinkLabel{}
	if label.Type != nil {
		switch *label.Type {
		case "horizontal":
			r.Type = types.LINK_LABEL_TYPE_HORIZONTAL
		default:
			r.Type = types.LINK_LABEL_TYPE_HORIZONTAL
		}
	} else {
		r.Type = types.LINK_LABEL_TYPE_HORIZONTAL
	}
	r.Title = label.Title
	if label.Color != nil {
		c := stringToColor(*label.Color)
		r.Color = &c
	}
	if label.Font != nil {
		r.Font = *label.Font
	}
	return r
}

func loadLinks(template *TemplateStruct, resources map[string]*types.Resource) {

	for _, v := range template.Links {
		_, ok := resources[v.Source]
		if !ok {
			log.Warnf("Not found Source esource %s", v.Source)
			continue
		}
		source := resources[v.Source]

		_, ok = resources[v.Target]
		if !ok {
			log.Warnf("Not found Target resource %s", v.Target)
			continue
		}
		target := resources[v.Target]

		log.Infof("Add link(%s-%s)", v.Source, v.Target)
		lineWidth := v.LineWidth
		if lineWidth == 0 {
			lineWidth = 2
		}

		lineColor := color.RGBA{0, 0, 0, 255}
		if v.LineColor != "" {
			lineColor = stringToColor(v.LineColor)
		}

		sourcePosition, err := types.ConvertWindrose(v.SourcePosition)
		if err != nil {
			panic(err)
		}
		targetPosition, err := types.ConvertWindrose(v.TargetPosition)
		if err != nil {
			panic(err)
		}
		link := new(types.Link).Init(source, sourcePosition, v.SourceArrowHead, target, targetPosition, v.TargetArrowHead, lineWidth, lineColor)
		link.SetType(v.Type)
		link.SetLineStyle(v.LineStyle)
		if v.Labels.SourceRight != nil {
			link.Labels.SourceRight = convertLabel(v.Labels.SourceRight)
		}
		if v.Labels.SourceLeft != nil {
			link.Labels.SourceLeft = convertLabel(v.Labels.SourceLeft)
		}
		if v.Labels.TargetRight != nil {
			link.Labels.TargetRight = convertLabel(v.Labels.TargetRight)
		}
		if v.Labels.TargetLeft != nil {
			link.Labels.TargetLeft = convertLabel(v.Labels.TargetLeft)
		}
		resources[v.Source].AddLink(link)
		resources[v.Target].AddLink(link)
	}

}

func IsURL(str string) bool {
	if strings.HasPrefix(str, "http://") || strings.HasPrefix(str, "https://") {
		return true
	}
	return false
}
