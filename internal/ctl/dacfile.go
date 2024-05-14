// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package ctl

import (
	"io"
	"net/http"
	"os"

	"github.com/awslabs/diagram-as-code/internal/definition"
	"github.com/awslabs/diagram-as-code/internal/types"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func CreateDiagramFromDacFile(inputfile string, outputfile *string) {

	log.Infof("input file path: %s\n", inputfile)

	var template TemplateStruct

	if IsURL(inputfile) {
		// URL from remote
		resp, err := http.Get(inputfile)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		err = yaml.Unmarshal(data, &template)
		if err != nil {
			log.Fatal(err)
		}

	} else {
		// Local file
		data, err := os.ReadFile(inputfile)
		if err != nil {
			log.Fatal(err)
		}
		err = yaml.Unmarshal([]byte(data), &template)
		if err != nil {
			log.Fatal(err)
		}
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
