// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
// Reference: https://github.com/aws-cloudformation/rain/blob/main/cft/graph/util.go

package ctl

import (
	"fmt"
	"image/color"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/aws-cloudformation/rain/cft"
	"github.com/aws-cloudformation/rain/cft/parse"
	"github.com/awslabs/diagram-as-code/internal/definition"
	"github.com/awslabs/diagram-as-code/internal/types"
	log "github.com/sirupsen/logrus"
)

var template = TemplateStruct{
	Diagram: Diagram{
		DefinitionFiles: []DefinitionFile{
			{
				Type: "URL",
				Url:  "https://raw.githubusercontent.com/awslabs/diagram-as-code/main/definitions/definition-for-aws-icons-light.yaml",
			},
		},
		Resources: map[string]Resource{
			"Canvas": {
				Type: "AWS::Diagram::Canvas",
				Children: []string{
					"AWSCloud",
				},
			},
			"AWSCloud": {
				Type:     "AWS::Diagram::Cloud",
				Preset:   "AWSCloudNoLogo",
				Align:    "center",
				Children: []string{},
			},
		},
		Links: []Link{},
	},
}

func CreateDiagramFromCFnTemplate(inputfile string, outputfile *string, generateDacFile bool, opts *CreateOptions) error {

	log.Infof("input file path: %s\n", inputfile)

	var cfn_template cft.Template

	if IsURL(inputfile) {
		// URL from remote
		resp, err := http.Get(inputfile)
		if err != nil {
			return fmt.Errorf("failed to get URL: %w", err)
		}
		defer func() {
			if closeErr := resp.Body.Close(); closeErr != nil {
				log.Warnf("Failed to close response body: %v", closeErr)
			}
		}()

		cfn_template, err = parse.Reader(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to parse CloudFormation template from URL: %w", err)
		}

	} else {
		// Local file
		var err error
		cfn_template, err = parse.File(inputfile)
		if err != nil {
			return fmt.Errorf("failed to parse CloudFormation template file: %w", err)
		}
	}

	var ds definition.DefinitionStructure
	resources := make(map[string]*types.Resource)

	log.Info("--- Load DefinitionFiles section ---")
	if opts.OverrideDefFile != "" {
		var overrideDefTemplate TemplateStruct
		if IsURL(opts.OverrideDefFile) {
			log.Infof("As given overrideDefFile, use %s as URL instead of %v", opts.OverrideDefFile, &template.DefinitionFiles)
			var defFile = DefinitionFile{
				Type: "URL",
				Url:  opts.OverrideDefFile,
			}
			overrideDefTemplate.DefinitionFiles = append(overrideDefTemplate.DefinitionFiles, defFile)
		} else {
			log.Infof("As given overrideDefFile, use %s as LocalFile instead of %v", opts.OverrideDefFile, &template.DefinitionFiles)
			var defFile = DefinitionFile{
				Type:      "LocalFile",
				LocalFile: opts.OverrideDefFile,
			}
			overrideDefTemplate.DefinitionFiles = append(overrideDefTemplate.DefinitionFiles, defFile)
		}
		if err := loadDefinitionFiles(&overrideDefTemplate, &ds); err != nil {
			return fmt.Errorf("failed to load override definition files: %w", err)
		}
		log.Infof("overrideDefTemplate: %+v", overrideDefTemplate)
	} else {
		if err := loadDefinitionFiles(&template, &ds); err != nil {
			return fmt.Errorf("failed to load definition files: %w", err)
		}
	}

	log.Info("--- Convert CloudFormation template to diagram structures ---")
	convertTemplate(cfn_template, &template, ds)

	log.Info("--- Ensuring a single parent for resources with multiple parents ---")
	ensureSingleParent(&template)

	log.Info("--- Load Resources section ---")
	if err := loadResources(&template, ds, resources); err != nil {
		return fmt.Errorf("failed to load resources: %w", err)
	}

	log.Info("--- Associate children with parent resources ---")
	associateCFnChildren(&template, ds, resources)

	if generateDacFile {
		log.Info("--- Generate dac file from CloudFormation template ---")
		go generateDacFileFromCFnTemplate(&template, *outputfile)
	}

	if err := createDiagram(resources, outputfile, opts); err != nil {
		return fmt.Errorf("failed to create diagram: %w", err)
	}
	return nil
}

func convertTemplate(cfn_template cft.Template, template *TemplateStruct, ds definition.DefinitionStructure) {

	resources_cfn_template := cfn_template.Map()["Resources"]

	if resourcesMap, ok := resources_cfn_template.(map[string]interface{}); ok {

		//Initialized with all logical IDs written in the template
		for logicalId, res := range resourcesMap {
			resource := res.(map[string]interface{})
			typeValue, _ := resource["Type"].(string)

			if _, ok := template.Resources[logicalId]; !ok {
				template.Resources[logicalId] = Resource{
					Type: typeValue,
				}
			}
		}

		//Check dependencies between resources
		for logicalId, res := range resourcesMap {
			resource := res.(map[string]interface{})

			var findParent bool

			//In CloudFormation templates, parameter names and resources are often related.
			//However, a parameter is not a "parent resource" of its resource.
			for _, related := range findRefs(resource, logicalId) {

				related = strings.Split(related, ".")[0]
				related_resource_type := template.Diagram.Resources[related].Type

				//related_resource_type does not have "Type". This means it may be a Parameter value
				if related_resource_type == "" {
					log.Infof("%s does not have \"Type\".", related)
					continue
				}

				def, ok := ds.Definitions[related_resource_type]
				if !ok {
					log.Infof("%s is not defined in the definition file.", related_resource_type)
					continue
				}

				//related_resource_type can not have children resources due to the restrict of definition file.
				if def == nil || !def.CFn.HasChildren {
					log.Infof("%s cannot have children resource.", related)
					continue
				}

				//Find parent
				findParent = true
				parent_logicalId := related
				parent_resources := template.Resources[parent_logicalId]
				parent_resources.Children = append(parent_resources.Children, logicalId)
				template.Resources[parent_logicalId] = parent_resources
			}

			//If there is no parent resource, consider "AWSCloud" as the parent
			if !findParent {
				parents := template.Resources["AWSCloud"]
				parents.Children = append(parents.Children, logicalId)
				template.Resources["AWSCloud"] = parents
			}
		}
	}
}

func ensureSingleParent(template *TemplateStruct) {
	for logicalId, resource := range template.Resources {

		if logicalId == "Canvas" || logicalId == "AWSCloud" {
			continue
		}

		if len(resource.Children) > 1 {

			for _, childID := range resource.Children {
				child, ok := template.Resources[childID]
				if !ok {
					continue
				}
				if len(child.Children) > 0 {
					grandchildrenIds := make([]string, 0)

					for _, grandchildID := range child.Children {
						if contains(resource.Children, grandchildID) {
							grandchildrenIds = append(grandchildrenIds, grandchildID)
							log.Infof("Found grandchild %s in resource %s", grandchildID, logicalId)
						}
					}

					newChildren := make([]string, 0, len(child.Children))
					for _, id := range resource.Children {
						if !contains(grandchildrenIds, id) {
							newChildren = append(newChildren, id)
						}
					}

					grandparent_resources := template.Resources[logicalId]
					grandparent_resources.Children = newChildren
					template.Resources[logicalId] = grandparent_resources
					resource.Children = newChildren
					log.Infof("Updated resource %s children: %v", logicalId, newChildren)
				}
			}
		}
	}
}

func associateCFnChildren(template *TemplateStruct, ds definition.DefinitionStructure, resources map[string]*types.Resource) {

	for logicalId, resource := range template.Resources {

		def, ok := ds.Definitions[resource.Type]

		if resource.Type == "" || !ok {
			log.Infof("%s is not defined in CloudFormation template or definition file. Skip process", logicalId)
			continue
		}

		if def == nil || !def.CFn.HasChildren {
			log.Infof("%s cannot have children resource.", logicalId)
			continue
		}

		if _, ok = resources[logicalId]; !ok {
			log.Infof("%s is not defined as a resource. Skip process", logicalId)
			continue
		}

		newChildren := make([]string, 0)

		for _, child := range resource.Children {
			_, ok := resources[child]
			if !ok {
				log.Infof("%s does not have parent resource", child)
				continue
			}
			log.Infof("Add child(%s) on %s", child, logicalId)

			if err := resources[logicalId].AddChild(resources[child]); err != nil {
				log.Warnf("Failed to add child %s to %s: %v", child, logicalId, err)
				continue
			}

			if def == nil || def.Border == nil {
				resources[logicalId].SetBorderColor(color.RGBA{0, 0, 0, 255})
				resources[logicalId].SetFillColor(color.RGBA{0, 0, 0, 0})
			}

			// Update yaml template for providing
			newChildren = append(newChildren, child)
			template_resource := template.Resources[logicalId]
			template_resource.Children = newChildren
			template.Resources[logicalId] = template_resource

		}
	}
}

func generateDacFileFromCFnTemplate(template *TemplateStruct, outputfile string) {

	yamlData, err := yaml.Marshal(template)
	if err != nil {
		fmt.Printf("Error while marshaling data: %v", err)
		return
	}

	outputBase := strings.TrimSuffix(outputfile, filepath.Ext(outputfile))
	outputYamlFile := outputBase + ".yaml"

	err = os.WriteFile(outputYamlFile, yamlData, 0644)
	if err != nil {
		fmt.Printf("Error while writing to file: %v", err)
		return
	}

	fmt.Printf("[Completed] dac (diagram-as-code) data written to %s\n", outputYamlFile)
}

func findRefs(t map[string]interface{}, fromName string) []string {
	refs := make([]string, 0)
	var subRe = regexp.MustCompile(`\$\{([^!].+?)\}`)

	for key, value := range t {
		switch key {
		case "DependsOn":
			switch v := value.(type) {
			case string:
				refs = append(refs, v)
			case []interface{}:
				for _, d := range v {
					//note: ECS Containers can have "DependsOn", but it should be ignored.
					if dStr, ok := d.(string); ok {
						refs = append(refs, dStr)
					}
				}
			default:
			}
		case "Ref":
			refs = append(refs, value.(string))
		case "Fn::GetAtt":
			switch v := value.(type) {
			case string:
				parts := strings.Split(v, ".")
				refs = append(refs, parts[0])
			case []interface{}:
				if s, ok := v[0].(string); ok {
					refs = append(refs, s)
				}
			default:
				fmt.Printf("Malformed GetAtt: %T\n", v)
			}
		case "Fn::Sub":
			switch v := value.(type) {
			case string:
				for _, groups := range subRe.FindAllStringSubmatch(v, 1) {
					refs = append(refs, groups[1])
				}
			case []interface{}:
				switch {
				case len(v) != 2:
					fmt.Printf("Malformed Sub: %T\n", v)
				default:
					switch parts := v[1].(type) {
					case map[string]interface{}:
						for _, part := range parts {
							switch p := part.(type) {
							case map[string]interface{}:
								refs = append(refs, findRefs(p, fromName)...)
							default:
								fmt.Printf("Malformed Sub: %T\n", v)
							}
						}
					default:
						fmt.Printf("Malformed Sub: %T\n", v)
					}
				}
			default:
				fmt.Printf("Malformed Sub: %T\n", v)
			}
		default:
			for _, tree := range findTrees(value) {
				refs = append(refs, findRefs(tree, fromName)...)
			}
		}
	}

	return refs
}

func findTrees(value interface{}) []map[string]interface{} {
	trees := make([]map[string]interface{}, 0)

	switch v := value.(type) {
	case map[string]interface{}:
		trees = append(trees, v)
	case []interface{}:
		for _, child := range v {
			trees = append(trees, findTrees(child)...)
		}
	}

	return trees
}

func contains(arr []string, data string) bool {
	for _, v := range arr {
		if v == data {
			return true
		}
	}
	return false
}
