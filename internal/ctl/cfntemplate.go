// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
// Reference: https://github.com/aws-cloudformation/rain/blob/main/cft/graph/util.go

package ctl

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/aws-cloudformation/rain/cft"
	"github.com/aws-cloudformation/rain/cft/parse"
	"github.com/awslabs/diagram-as-code/internal/definition"
	"github.com/awslabs/diagram-as-code/internal/types"
	log "github.com/sirupsen/logrus"
)

func CreateDiagramFromCFnTemplate(inputfile string, outputfile *string) {

	log.Infof("input file: %s\n", inputfile)
	cfn_template, err := parse.File(inputfile)
	if err != nil {
		log.Fatal(err)
	}

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

	createTemplate(cfn_template, &template)

	var ds definition.DefinitionStructure
	var resources map[string]types.Node = make(map[string]types.Node)

	log.Info("Load DefinitionFiles section")
	loadDefinitionFiles(&template, &ds)

	log.Info("Load Resources section")
	loadResources(&template, ds, resources)

	log.Info("Associate children with parent resources")
	associateCFnChildren(&template, ds, resources)

	createDiagram(resources, outputfile)
}

func createTemplate(cfn_template cft.Template, template *TemplateStruct) {

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
			for _, parent := range findRefs(resource, logicalId) {
				findParent = true
				parent = strings.Split(parent, ".")[0]
				parents := template.Diagram.Resources[parent]
				parents.Children = append(parents.Children, logicalId)
				template.Diagram.Resources[parent] = parents
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
					refs = append(refs, d.(string))
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
