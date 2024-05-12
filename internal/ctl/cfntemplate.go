// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
// Reference: https://github.com/aws-cloudformation/rain/blob/main/cft/graph/util.go

package ctl

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"image/color"
	"os"
	"path/filepath"
	"regexp"
	"strings"

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

func CreateDiagramFromCFnTemplate(inputfile string, outputfile *string, generateYaml bool) {

	log.Infof("input file: %s\n", inputfile)
	cfn_template, err := parse.File(inputfile)
	if err != nil {
		log.Fatal(err)
	}

	var ds definition.DefinitionStructure
	var resources map[string]types.Node = make(map[string]types.Node)

	log.Info("--- Load DefinitionFiles section ---")
	loadDefinitionFiles(&template, &ds)

	log.Info("--- Convert CloudFormation template to diagram structures ---")
	convertTemplate(cfn_template, &template, ds)

	log.Info("--- Ensuring a single parent for resources with multiple parents ---")
	ensureSingleParent(&template)

	log.Info("--- Load Resources section ---")
	loadResources(&template, ds, resources)

	log.Info("--- Associate children with parent resources ---")
	associateCFnChildren(&template, ds, resources)

	if generateYaml {
		log.Info("--- Generate yaml file from CloudFormation template ---")
		go generateYamlFromCFnTemplate(&template, *outputfile)
	}

	createDiagram(resources, outputfile)
}

func convertTemplate(cfn_template cft.Template, template *TemplateStruct, ds definition.DefinitionStructure) {

	resources_cfn_template := cfn_template.Map()["Resources"]

	if resourcesMap, ok := resources_cfn_template.(map[string]interface{}); ok {

		//Initialized with all logical IDs written in the template
		for logicalId, res := range resourcesMap {
			resource := res.(map[string]interface{})
			typeValue, _ := resource["Type"].(string)

			if _, ok := template.Diagram.Resources[logicalId]; !ok {
				template.Diagram.Resources[logicalId] = Resource{
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
				if !def.CFn.HasChildren {
					log.Infof("%s cannot have children resource.", related)
					continue
				}

				//Find parent
				findParent = true
				parent_logicalId := related
				parent_resources := template.Diagram.Resources[parent_logicalId]
				parent_resources.Children = append(parent_resources.Children, logicalId)
				template.Diagram.Resources[parent_logicalId] = parent_resources
			}

			//If there is no parent resource, consider "AWSCloud" as the parent
			if !findParent {
				parents := template.Diagram.Resources["AWSCloud"]
				parents.Children = append(parents.Children, logicalId)
				template.Diagram.Resources["AWSCloud"] = parents
			}
		}
	}
}

func ensureSingleParent(template *TemplateStruct) {
	for logicalId, resource := range template.Diagram.Resources {

		if logicalId == "Canvas" || logicalId == "AWSCloud" {
			continue
		}

		if len(resource.Children) > 1 {

			for _, childID := range resource.Children {
				child, ok := template.Diagram.Resources[childID]
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

					grandparent_resources := template.Diagram.Resources[logicalId]
					grandparent_resources.Children = newChildren
					template.Diagram.Resources[logicalId] = grandparent_resources
					resource.Children = newChildren
					log.Infof("Updated resource %s children: %v", logicalId, newChildren)
				}
			}
		}
	}
}

func associateCFnChildren(template *TemplateStruct, ds definition.DefinitionStructure, resources map[string]types.Node) {

	for logicalId, resource := range template.Resources {

		def, ok := ds.Definitions[resource.Type]

		if resource.Type == "" || !ok {
			log.Infof("%s is not defined in CloudFormation template or definition file. Skip process", logicalId)
			continue
		}

		if !def.CFn.HasChildren {
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

			resources[logicalId].AddChild(resources[child])

			if def.Border == nil {
				resources[logicalId].SetBorderColor(color.RGBA{0, 0, 0, 255})
				resources[logicalId].SetFillColor(color.RGBA{0, 0, 0, 0})
			}

			// Update yaml template for providing
			newChildren = append(newChildren, child)
			template_resource := template.Diagram.Resources[logicalId]
			template_resource.Children = newChildren
			template.Diagram.Resources[logicalId] = template_resource

		}
	}
}

func generateYamlFromCFnTemplate(template *TemplateStruct, outputfile string) {

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

	fmt.Printf("[Completed] Data written to %s\n", outputYamlFile)
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
