// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package ctl

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"strings"

	"github.com/awslabs/diagram-as-code/internal/cache"
	"github.com/awslabs/diagram-as-code/internal/definition"
	"github.com/awslabs/diagram-as-code/internal/types"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/maps"
	"golang.org/x/image/draw"
)

func stringToColor(c string) (color.RGBA, error) {
	var r, g, b, a uint8
	_, err := fmt.Sscanf(c, "rgba(%d,%d,%d,%d)", &r, &g, &b, &a)
	if err != nil {
		return color.RGBA{}, fmt.Errorf("failed to parse color string '%s': %w", c, err)
	}
	return color.RGBA{r, g, b, a}, nil
}

// OverwriteMode defines how to handle existing output files
type OverwriteMode int

const (
	// Ask shows confirmation prompt when output file exists (CLI default)
	Ask OverwriteMode = iota
	// Force overwrites without confirmation (CLI with --force)
	Force
	// NoOverwrite refuses to overwrite and returns error (MCP server default)
	NoOverwrite
)

// CheckOutputFileOverwrite checks if output file exists and handles according to mode
func CheckOutputFileOverwrite(outputFile string, mode OverwriteMode) error {
	// Check if file exists
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		// File doesn't exist, proceed
		return nil
	} else if err != nil {
		// Some other error occurred
		return fmt.Errorf("failed to check output file: %w", err)
	}

	// File exists, handle according to mode
	switch mode {
	case Force:
		// Force mode: proceed without confirmation
		return nil
	case NoOverwrite:
		// NoOverwrite mode: return error
		return fmt.Errorf("output file '%s' already exists", outputFile)
	case Ask:
		// Ask mode: show confirmation prompt
		return askOverwriteConfirmation(outputFile)
	default:
		return fmt.Errorf("unknown overwrite mode: %d", mode)
	}
}

// askOverwriteConfirmation shows interactive confirmation prompt
func askOverwriteConfirmation(outputFile string) error {
	fmt.Printf("File '%s' already exists. Overwrite? [y/N]: ", outputFile)

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read user input: %w", err)
	}

	response = strings.TrimSpace(strings.ToLower(response))
	if response == "y" || response == "yes" {
		return nil
	}

	return fmt.Errorf("operation cancelled by user")
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
	Options        *ResourceOptions  `yaml:"Options"`
}

type ResourceOptions struct {
	GroupingOffset *bool `yaml:"GroupingOffset"`
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
	OverwriteMode   OverwriteMode
	OverrideFont    string
	Width           int
	Height          int
}

func createDiagram(resources map[string]*types.Resource, outputfile *string, opts *CreateOptions) error {

	// Check for file overwrite before processing
	if err := CheckOutputFileOverwrite(*outputfile, opts.OverwriteMode); err != nil {
		return err
	}

	// Override font if specified
	if opts.OverrideFont != "" {
		for _, resource := range resources {
			resource.SetLabel(nil, nil, &opts.OverrideFont)
		}
	}

	log.Info("--- Draw diagram ---")
	canvas, exists := resources["Canvas"]
	if !exists {
		return fmt.Errorf("Canvas resource not found")
	}
	err := canvas.Scale(nil, nil)
	if err != nil {
		return fmt.Errorf("error scaling diagram: %w", err)
	}
	if err := canvas.ZeroAdjust(); err != nil {
		return fmt.Errorf("error adjusting diagram: %w", err)
	}

	// Resolve auto-positions after layout is complete
	for _, resource := range resources {
		for _, link := range resource.GetLinks() {
			link.ResolveAutoPositions()
		}
	}

	img, err := canvas.Draw(nil, nil)
	if err != nil {
		return fmt.Errorf("error drawing diagram: %w", err)
	}

	// Resize the image if width or height is specified
	if opts != nil && (opts.Width > 0 || opts.Height > 0) {
		log.Infof("Resizing image to width: %d, height: %d", opts.Width, opts.Height)
		resizedImg := resizeImage(img, opts.Width, opts.Height)
		img = resizedImg
	}

	log.Infof("Save %s\n", *outputfile)
	fmt.Printf("[Completed] AWS infrastructure diagram generated: %s\n", *outputfile)
	f, err := os.OpenFile(*outputfile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("error opening output file: %w", err)
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			log.Warnf("Failed to close output file: %v", closeErr)
		}
	}()
	if err := png.Encode(f, img); err != nil {
		return fmt.Errorf("error encoding PNG: %w", err)
	}
	return nil
}

// resizeImage resizes the image while maintaining aspect ratio
func resizeImage(src *image.RGBA, width, height int) *image.RGBA {
	// Get original dimensions
	bounds := src.Bounds()
	srcWidth := bounds.Dx()
	srcHeight := bounds.Dy()

	// If neither width nor height is specified, return the original image
	if width == 0 && height == 0 {
		return src
	}

	// Calculate new dimensions while maintaining aspect ratio
	var ratio float64
	if width > 0 && height > 0 {
		// Both width and height specified, fit within these constraints
		widthRatio := float64(width) / float64(srcWidth)
		heightRatio := float64(height) / float64(srcHeight)

		// Use the smaller ratio to ensure the image fits within the specified dimensions
		ratio = math.Min(widthRatio, heightRatio)
	} else if width > 0 {
		// Only width specified
		ratio = float64(width) / float64(srcWidth)
	} else {
		// Only height specified
		ratio = float64(height) / float64(srcHeight)
	}

	newWidth := int(float64(srcWidth) * ratio)
	newHeight := int(float64(srcHeight) * ratio)

	// Create a new RGBA image with the calculated dimensions
	dst := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	// Resize the image using CatmullRom algorithm for better quality
	draw.CatmullRom.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Over, nil)

	return dst
}

func loadDefinitionFiles(template *TemplateStruct, ds *definition.DefinitionStructure) error {

	// Load definition files
	for _, v := range template.DefinitionFiles {
		switch v.Type {
		case "URL":
			log.Infof("Fetch definition file from URL: %s\n", v.Url)
			cacheFilePath, err := cache.FetchFile(v.Url)
			if err != nil {
				return fmt.Errorf("failed to fetch definition file from URL %s: %w", v.Url, err)
			}
			log.Infof("Read definition file from cache file: %s\n", cacheFilePath)
			err = ds.LoadDefinitions(cacheFilePath)
			if err != nil {
				return fmt.Errorf("failed to load definitions from cache file %s: %w", cacheFilePath, err)
			}
		case "LocalFile":
			log.Infof("Read definition file from path: %s\n", v.LocalFile)
			err := ds.LoadDefinitions(v.LocalFile)
			if err != nil {
				return fmt.Errorf("failed to load definitions from local file %s: %w", v.LocalFile, err)
			}
		case "Embed":
			log.Info("Read embedded definitions")
			maps.Copy(ds.Definitions, v.Embed.Definitions)
		}
	}
	return nil

}

func loadResources(template *TemplateStruct, ds definition.DefinitionStructure, resources map[string]*types.Resource) error {

	resources["Canvas"] = new(types.Resource).Init()

	for k, v := range template.Resources {
		// Override order: Definition{Resource Type -> Preset} -> Template
		log.Infof("Load Resource: %s (%s)\n", k, v.Type)
		switch v.Type {
		case "":
			log.Warnf("%s does not have Type field. Skipping this resource.", k)
			continue
		case "AWS::Diagram::Canvas":
			resource, exists := resources[k]
			if !exists {
				return fmt.Errorf("Canvas resource %s not found in resources map", k)
			}
			resource.SetBorderColor(color.RGBA{0, 0, 0, 0})
			resource.SetFillColor(color.RGBA{255, 255, 255, 255})
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
				fallbackDef, check := ds.Definitions[newType]
				if !check || fallbackDef == nil {
					log.Warnf("Type %s is not defined in the DAC definition file. It cannot be fall backed to service icon. Ignore this type.\n", v.Type)
					continue
				}
				log.Warnf("Type %s is not defined in the DAC definition file. It's fall backed to its service icon (Type %s).\n", v.Type, newType)
				def = fallbackDef
			}
			if def == nil {
				log.Warnf("Definition for %s is nil. Skip this resource.\n", v.Type)
				continue
			}
			switch def.Type {
			case "Resource":
				resources[k] = new(types.Resource).Init()
			case "Group":
				resources[k] = new(types.Resource).Init()
			}
			if fill := def.Fill; fill != nil {
				fillColor, err := stringToColor(fill.Color)
				if err != nil {
					return fmt.Errorf("failed to parse fill color for resource %s: %w", k, err)
				}
				resource, exists := resources[k]
				if !exists {
					return fmt.Errorf("resource %s not found for fill color", k)
				}
				resource.SetFillColor(fillColor)
			}
			if border := def.Border; border != nil {
				borderColor, err := stringToColor(border.Color)
				if err != nil {
					return fmt.Errorf("failed to parse border color for resource %s: %w", k, err)
				}
				resource, exists := resources[k]
				if !exists {
					return fmt.Errorf("resource %s not found when setting border", k)
				}
				resource.SetBorderColor(borderColor)
				switch border.Type {
				case "straight":
					resource.SetBorderType(types.BORDER_TYPE_STRAIGHT)
				case "dashed":
					resource.SetBorderType(types.BORDER_TYPE_DASHED)
				default:
					resource.SetBorderType(types.BORDER_TYPE_STRAIGHT)
				}
			}
			if label := def.Label; label != nil {
				resource, exists := resources[k]
				if !exists {
					return fmt.Errorf("resource %s not found for label", k)
				}
				if label.Title != "" {
					resource.SetLabel(&label.Title, nil, nil)
				}
				if label.Color != "" {
					c, err := stringToColor(label.Color)
					if err != nil {
						return fmt.Errorf("failed to parse label color for resource %s: %w", k, err)
					}
					resource, exists := resources[k]
					if !exists {
						return fmt.Errorf("resource %s not found when setting label color", k)
					}
					resource.SetLabel(nil, &c, nil)
				}
				if label.Font != "" {
					resource, exists := resources[k]
					if !exists {
						return fmt.Errorf("resource %s not found when setting label font", k)
					}
					resource.SetLabel(nil, nil, &label.Font)
				}
			}
			if headerAlign := def.HeaderAlign; headerAlign != "" {
				resource, exists := resources[k]
				if !exists {
					return fmt.Errorf("resource %s not found when setting header align", k)
				}
				resource.SetHeaderAlign(headerAlign)
			}
			if icon := def.Icon; icon != nil {
				if def.CacheFilePath == "" {
					break
				}
				resource, exists := resources[k]
				if !exists {
					return fmt.Errorf("resource %s not found when loading icon", k)
				}
				err := resource.LoadIcon(def.CacheFilePath)
				if err != nil {
					return fmt.Errorf("failed to load icon from cache file path: %w", err)
				}
			}
		}

		switch v.Preset {
		case "BlankGroup":
			resource, exists := resources[k]
			if !exists {
				return fmt.Errorf("resource %s not found for BlankGroup preset", k)
			}
			resource.SetIconBounds(image.Rect(0, 0, 64, 64))
			resource.SetBorderColor(color.RGBA{0, 0, 0, 0})
		case "":
		default:
			def, ok := ds.Definitions[v.Preset]
			if !ok {
				log.Warnf("Unknown preset %s on %s\n", v.Preset, v.Type)
			} else {
				resource, exists := resources[k]
				if !exists {
					return fmt.Errorf("resource %s not found for preset configuration", k)
				}
				if fill := def.Fill; fill != nil {
					fillColor, err := stringToColor(fill.Color)
					if err != nil {
						return fmt.Errorf("failed to parse fill color for resource %s: %w", k, err)
					}
					resource.SetFillColor(fillColor)
				}
				if border := def.Border; border != nil {
					borderColor, err := stringToColor(border.Color)
					if err != nil {
						return fmt.Errorf("failed to parse border color for resource %s: %w", k, err)
					}
					resource.SetBorderColor(borderColor)
					switch border.Type {
					case "straight":
						resource.SetBorderType(types.BORDER_TYPE_STRAIGHT)
					case "dashed":
						resource.SetBorderType(types.BORDER_TYPE_DASHED)
					default:
						resource.SetBorderType(types.BORDER_TYPE_STRAIGHT)
					}
				}
				if label := def.Label; label != nil {
					if label.Title != "" {
						resource.SetLabel(&label.Title, nil, nil)
					}
					if label.Color != "" {
						c, err := stringToColor(label.Color)
						if err != nil {
							return fmt.Errorf("failed to parse label color for resource %s: %w", k, err)
						}
						resource.SetLabel(nil, &c, nil)
					}
					if label.Font != "" {
						resource.SetLabel(nil, nil, &label.Font)
					}
				}
				if headerAlign := def.HeaderAlign; headerAlign != "" {
					resource.SetHeaderAlign(headerAlign)
				}
				if icon := def.Icon; icon != nil {
					if def.CacheFilePath != "" {
						err := resource.LoadIcon(def.CacheFilePath)
						if err != nil {
							return fmt.Errorf("failed to load icon from cache file path: %w", err)
						}
					}
				}
			}
		}
		if v.Icon != "" {
			resource, exists := resources[k]
			if !exists {
				return fmt.Errorf("resource %s not found for icon loading", k)
			}
			err := resource.LoadIcon(v.Icon)
			if err != nil {
				return fmt.Errorf("failed to load icon from file: %w", err)
			}
		}
		if v.IconFill != nil {
			resource, exists := resources[k]
			if !exists {
				return fmt.Errorf("resource %s not found for icon fill", k)
			}
			switch *v.IconFill.Type {
			case "none":
				resource.SetIconFill(types.ICON_FILL_TYPE_NONE, nil)
			case "rect":
				if v.IconFill.Color != nil {
					c, err := stringToColor(*v.IconFill.Color)
					if err != nil {
						return fmt.Errorf("failed to parse icon fill color for resource %s: %w", k, err)
					}
					resource.SetIconFill(types.ICON_FILL_TYPE_RECT, &c)
				} else {
					resource.SetIconFill(types.ICON_FILL_TYPE_RECT, nil)
				}
			default:
				resource.SetIconFill(types.ICON_FILL_TYPE_NONE, nil)
			}
		}
		if v.Title != "" {
			resource, exists := resources[k]
			if !exists {
				return fmt.Errorf("resource %s not found for title", k)
			}
			resource.SetLabel(&v.Title, nil, nil)
		}
		if v.TitleColor != "" {
			c, err := stringToColor(v.TitleColor)
			if err != nil {
				return fmt.Errorf("failed to parse title color for resource %s: %w", k, err)
			}
			resource, exists := resources[k]
			if !exists {
				return fmt.Errorf("resource %s not found for title color", k)
			}
			resource.SetLabel(nil, &c, nil)
		}
		if v.HeaderAlign != "" {
			resource, exists := resources[k]
			if !exists {
				return fmt.Errorf("resource %s not found for header align", k)
			}
			resource.SetHeaderAlign(v.HeaderAlign)
		}
		if v.Font != "" {
			resource, exists := resources[k]
			if !exists {
				return fmt.Errorf("resource %s not found for font", k)
			}
			resource.SetLabel(nil, nil, &v.Font)
		}
		if v.Align != "" {
			resource, exists := resources[k]
			if !exists {
				return fmt.Errorf("resource %s not found for align", k)
			}
			resource.SetAlign(v.Align)
		}
		if v.Direction != "" {
			resource, exists := resources[k]
			if !exists {
				return fmt.Errorf("resource %s not found for direction", k)
			}
			resource.SetDirection(v.Direction)
		}
		if v.FillColor != "" {
			fillColor, err := stringToColor(v.FillColor)
			if err != nil {
				return fmt.Errorf("failed to parse fill color for resource %s: %w", k, err)
			}
			resource, exists := resources[k]
			if !exists {
				return fmt.Errorf("resource %s not found for fill color", k)
			}
			resource.SetFillColor(fillColor)
		}
		if v.BorderColor != "" {
			borderColor, err := stringToColor(v.BorderColor)
			if err != nil {
				return fmt.Errorf("failed to parse border color for resource %s: %w", k, err)
			}
			resource, exists := resources[k]
			if !exists {
				return fmt.Errorf("resource %s not found for border color", k)
			}
			resource.SetBorderColor(borderColor)
		}

		// Process Options
		if v.Options != nil {
			resource, exists := resources[k]
			if !exists {
				return fmt.Errorf("resource %s not found for options", k)
			}
			if v.Options.GroupingOffset != nil {
				resource.SetGroupingOffset(*v.Options.GroupingOffset)
			}
		}
	}

	return nil
}

func fallbackToServiceIcon(inputType string) string {

	parts := strings.SplitN(inputType, "::", 3)
	possibleServiceType := strings.Join(parts[:2], "::")

	return possibleServiceType
}

func associateChildren(template *TemplateStruct, resources map[string]*types.Resource) error {

	for logicalId, v := range template.Resources {
		resource, ok := resources[logicalId]
		if !ok {
			return fmt.Errorf("unknown resource %s", logicalId)
		}
		for _, child := range v.Children {
			childResource, ok := resources[child]
			if ok {
				log.Infof("Add child(%s) on %s", child, logicalId)
				if err := resource.AddChild(childResource); err != nil {
					return fmt.Errorf("failed to add child %s to %s: %w", child, logicalId, err)
				}
			} else {
				log.Warnf("Child `%s` was not found, ignoring it.", child)
			}
		}
		for _, borderChild := range v.BorderChildren {
			borderChildResource, ok := resources[borderChild.Resource]
			if !ok {
				log.Warnf("Child `%s` was not found, ignoring it.", borderChild.Resource)
				continue
			}
			log.Infof("Add BorderChild(%s) on %s", borderChild.Resource, logicalId)

			position, err := types.ConvertWindrose(borderChild.Position)
			if err != nil {
				return fmt.Errorf("failed to convert windrose position: %w", err)
			}
			bc := types.BorderChild{
				Position: position,
				Resource: borderChildResource,
			}
			if err := resource.AddBorderChild(&bc); err != nil {
				return fmt.Errorf("failed to add border child: %w", err)
			}
		}
	}
	return nil
}

func convertLabel(label *LinkLabel) (*types.LinkLabel, error) {
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
		c, err := stringToColor(*label.Color)
		if err != nil {
			return nil, fmt.Errorf("failed to parse label color: %w", err)
		}
		r.Color = &c
	}
	if label.Font != nil {
		r.Font = *label.Font
	}
	return r, nil
}

func loadLinks(template *TemplateStruct, resources map[string]*types.Resource) error {

	for _, v := range template.Links {
		sourceResource, ok := resources[v.Source]
		if !ok {
			log.Warnf("Not found Source esource %s", v.Source)
			continue
		}
		source := sourceResource

		targetResource, ok := resources[v.Target]
		if !ok {
			log.Warnf("Not found Target resource %s", v.Target)
			continue
		}
		target := targetResource

		log.Infof("Add link(%s-%s)", v.Source, v.Target)
		lineWidth := v.LineWidth
		if lineWidth == 0 {
			lineWidth = 2
		}

		lineColor := color.RGBA{0, 0, 0, 255}
		if v.LineColor != "" {
			var err error
			lineColor, err = stringToColor(v.LineColor)
			if err != nil {
				return fmt.Errorf("failed to parse line color: %w", err)
			}
		}

		// Convert positions (empty string and "auto" both become WINDROSE_AUTO)
		sourcePosition, err := types.ConvertWindrose(v.SourcePosition)
		if err != nil {
			return fmt.Errorf("failed to convert source windrose position: %w", err)
		}
		targetPosition, err := types.ConvertWindrose(v.TargetPosition)
		if err != nil {
			return fmt.Errorf("failed to convert target windrose position: %w", err)
		}

		link := new(types.Link).Init(source, sourcePosition, v.SourceArrowHead, target, targetPosition, v.TargetArrowHead, lineWidth, lineColor)
		link.SetType(v.Type)
		link.SetLineStyle(v.LineStyle)
		if v.Labels.SourceRight != nil {
			label, err := convertLabel(v.Labels.SourceRight)
			if err != nil {
				return fmt.Errorf("failed to convert source right label: %w", err)
			}
			link.Labels.SourceRight = label
		}
		if v.Labels.SourceLeft != nil {
			label, err := convertLabel(v.Labels.SourceLeft)
			if err != nil {
				return fmt.Errorf("failed to convert source left label: %w", err)
			}
			link.Labels.SourceLeft = label
		}
		if v.Labels.TargetRight != nil {
			label, err := convertLabel(v.Labels.TargetRight)
			if err != nil {
				return fmt.Errorf("failed to convert target right label: %w", err)
			}
			link.Labels.TargetRight = label
		}
		if v.Labels.TargetLeft != nil {
			label, err := convertLabel(v.Labels.TargetLeft)
			if err != nil {
				return fmt.Errorf("failed to convert target left label: %w", err)
			}
			link.Labels.TargetLeft = label
		}
		source.AddLink(link)
		target.AddLink(link)
	}
	return nil
}

func IsURL(str string) bool {
	if strings.HasPrefix(str, "http://") || strings.HasPrefix(str, "https://") {
		return true
	}
	return false
}
