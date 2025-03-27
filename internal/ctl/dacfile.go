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

func CreateDiagramFromDacFile(inputfile string, outputfile *string, overrideDefFile string) {

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
	var resources map[string]*types.Resource = make(map[string]*types.Resource)

	log.Info("Load DefinitionFiles section")
	if overrideDefFile != "" {
		var overrideDefTemplate TemplateStruct
		if IsURL(overrideDefFile) {
			log.Infof("As given overrideDefFile, use %s as URL instead of %v", overrideDefFile, &template.DefinitionFiles)
			var defFile = DefinitionFile{
				Type: "URL",
				Url: overrideDefFile,
			}
			overrideDefTemplate.Diagram.DefinitionFiles = append(overrideDefTemplate.Diagram.DefinitionFiles, defFile)
		} else {
			log.Infof("As given overrideDefFile, use %s as LocalFile instead of %v", overrideDefFile, &template.DefinitionFiles)
			var defFile = DefinitionFile{
				Type: "LocalFile",
				LocalFile: overrideDefFile,
			}
			overrideDefTemplate.Diagram.DefinitionFiles = append(overrideDefTemplate.Diagram.DefinitionFiles, defFile)
		}
		loadDefinitionFiles(&overrideDefTemplate, &ds)
		log.Infof("overrideDefTemplate: %+v", overrideDefTemplate)
	} else {
		loadDefinitionFiles(&template, &ds)
	}

	log.Info("Load Resources section")
	loadResources(&template, ds, resources)

	log.Info("Associate children with parent resources")
	associateChildren(&template, resources)

	log.Info("Add Links section")
	loadLinks(&template, resources)

	createDiagram(resources, outputfile)
}
