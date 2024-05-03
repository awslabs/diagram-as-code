// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package ctl

import (
	"os"

	"github.com/awslabs/diagram-as-code/internal/definition"
	"github.com/awslabs/diagram-as-code/internal/types"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func CreateDiagramFromYAML(inputfile string, outputfile *string) {

	log.Infof("input file: %s\n", inputfile)
	data, err := os.ReadFile(inputfile)
	if err != nil {
		log.Fatal(err)
	}

	var template TemplateStruct

	err = yaml.Unmarshal([]byte(data), &template)
	if err != nil {
		log.Fatal(err)
	}

	var ds definition.DefinitionStructure
	var resources map[string]types.Node = make(map[string]types.Node)

	log.Info("Load DefinitionFiles section")
	loadDefinitionFiles(&template, &ds)

	log.Info("Load Resources section")
	loadResources(&template, ds, resources)

	log.Info("Associate children with parent resources")
	associateChildren(&template, resources)

	log.Info("Add Links section")
	loadLinks(&template, resources)

	createDiagram(resources, outputfile)
}
