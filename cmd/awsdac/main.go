// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/awslabs/diagram-as-code/internal/cache"
	"github.com/awslabs/diagram-as-code/internal/definition"
	"github.com/awslabs/diagram-as-code/internal/types"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/maps"
	"gopkg.in/yaml.v3"
)

func stringToColor(c string) color.RGBA {
	var r, g, b, a uint8
	fmt.Sscanf(c, "rgba(%d,%d,%d,%d)", &r, &g, &b, &a)
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
	Type      string   `yaml:"Type"`
	Icon      string   `yaml:"Icon"`
	Direction string   `yaml:"Direction"`
	Preset    string   `yaml:"Preset"`
	Align     string   `yaml:"Align"`
	FillColor string   `yaml:"FillColor"`
	Title     string   `yaml:"Title"`
	Children  []string `yaml:"Children"`
}

type Link struct {
	Source         string `yaml:"Source"`
	SourcePosition string `yaml:"SourcePosition"`
	Target         string `yaml:"Target"`
	TargetPosition string `yaml:"TargetPosition"`
	LineWidth      int    `yaml:"LineWidth"`
}

func main() {
	debug := flag.Bool("v", false, "Debugging outputs")
	outputfile := flag.String("o", "output.png", "Output file")
	flag.Parse()
	inputfile := flag.Arg(0)

	if *debug {
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}

	log.Infof("input file: %s\n", inputfile)
	data, err := os.ReadFile(inputfile)
	if err != nil {
		log.Fatal(err)
	}

	var b TemplateStruct

	err = yaml.Unmarshal([]byte(data), &b)
	if err != nil {
		log.Fatal(err)
	}

	// Load definition files
	var ds definition.DefinitionStructure
	for _, v := range b.Diagram.DefinitionFiles {
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

	var resources map[string]types.Node = make(map[string]types.Node)
	resources["Canvas"] = new(types.Group).Init()

	log.Info("Add resources")
	for k, v := range b.Resources {
		title := v.Title
		log.Infof("Load Resource: %s (%s)\n", k, v.Type)
		switch v.Type {
		case "AWS::Diagram::Canvas":
			resources[k].SetBorderColor(color.RGBA{0, 0, 0, 0})
			resources[k].SetFillColor(color.RGBA{255, 255, 255, 255})
		case "AWS::Diagram::Group":
			resources[k] = new(types.Group).Init()
		case "AWS::Diagram::Resource":
			resources[k] = new(types.Resource).Init()
		case "AWS::Diagram::VerticalStack":
			resources[k] = new(types.VerticalStack).Init()
		case "AWS::Diagram::HorizontalStack":
			resources[k] = new(types.HorizontalStack).Init()
		default:
			def, ok := ds.Definitions[v.Type]
			if !ok {
				log.Fatalf("Unknown type: %s\n", v.Type)
			}
			if def.Type == "Resource" {
				resources[k] = new(types.Resource).Init()
			} else if def.Type == "Group" {
				resources[k] = new(types.Group).Init()
			}
			if fill := def.Fill; fill != nil {
				resources[k].SetFillColor(stringToColor(fill.Color))
			}
			if border := def.Border; border != nil {
				resources[k].SetBorderColor(stringToColor(border.Color))
			}
			if label := def.Label; label != nil {
				resources[k].SetLabel(label.Title, stringToColor(label.Color))
			}
			if icon := def.Icon; icon != nil {
				if def.CacheFilePath == "" {
					break
				}
				resources[k].LoadIcon(def.CacheFilePath)
			}
		}
		switch v.Preset {
		case "BlankGroup":
			resources[k].SetIconBounds(image.Rect(0, 0, 64, 64))
		case "":
			break
		default:
			def, ok := ds.Definitions[v.Preset]
			if !ok {
				log.Fatalf("Unknown preset: %s\n", v.Preset)
			}
			if fill := def.Fill; fill != nil {
				resources[k].SetFillColor(stringToColor(fill.Color))
			}
			if border := def.Border; border != nil {
				resources[k].SetBorderColor(stringToColor(border.Color))
			}
			if label := def.Label; label != nil {
				resources[k].SetLabel(label.Title, stringToColor(label.Color))
			}
			if icon := def.Icon; icon != nil {
				if def.CacheFilePath == "" {
					continue
				}
				resources[k].LoadIcon(def.CacheFilePath)
			}
		}
		if v.Icon != "" {
			resources[k].LoadIcon(v.Icon)
		}
		if v.Title != "" {
			resources[k].SetLabel(title, color.RGBA{0, 0, 0, 255})
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
	}
	log.Info("Add children")
	for k, v := range b.Resources {
		for _, child := range v.Children {
			_, ok := resources[child]
			if !ok {
				log.Warnf("Not found resource %s", child)
				continue
			}
			log.Infof("Add child(%s) on %s", child, k)
			resources[k].AddChild(resources[child])
		}
	}

	log.Info("Add links")
	for _, v := range b.Links {
		_, ok := resources[v.Source]
		if !ok {
			log.Warnf("Not found resource %s", v.Source)
			continue
		}
		source := resources[v.Source]

		_, ok = resources[v.Target]
		if !ok {
			log.Warnf("Not found resource %s", v.Target)
			continue
		}
		target := resources[v.Target]

		log.Infof("Add link(%s-%s)", v.Source, v.Target)
		lineWidth := v.LineWidth
		if lineWidth == 0 {
			lineWidth = 2
		}
		link := new(types.Link).Init(&source, v.SourcePosition, &target, v.TargetPosition, lineWidth)
		resources[v.Source].AddLink(link)
		resources[v.Target].AddLink(link)
	}

	log.Info("Drawing")
	resources["Canvas"].Scale()
	resources["Canvas"].ZeroAdjust()
	img := resources["Canvas"].Draw(nil)

	log.Infof("Save %s\n", *outputfile)
	f, _ := os.OpenFile(*outputfile, os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	png.Encode(f, img)
}
