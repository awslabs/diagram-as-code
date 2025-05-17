// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package ctl

import (
	"bytes"
	"io"
	"net/http"
	"os"
	tmpl "text/template"

	"golang.org/x/exp/slices"
	"github.com/awslabs/diagram-as-code/internal/definition"
	"github.com/awslabs/diagram-as-code/internal/types"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func getTemplate(inputfile string) ([]byte, error) {
	var data []byte
	var err error

	if IsURL(inputfile) {
		// URL from remote
		resp, err := http.Get(inputfile)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		data, err = io.ReadAll(resp.Body)
	} else {
		// Local file
		data, err = os.ReadFile(inputfile)
	}
	if err != nil {
		return nil, err
	}
	return data, nil
}

func processTemplate(templateData []byte) ([]byte, error) {
	// Create a new template
	tmpl, err := tmpl.New("dacfile").Funcs(funcMap).Parse(string(templateData))
	if err != nil {
		log.Infof("%v", tmpl)
		return nil, err
	}

	// Create a buffer to store the processed template
	var processed bytes.Buffer

	// Execute the template with the provided variables
	err = tmpl.Execute(&processed, nil)
	if err != nil {
		return nil, err
	}

	return processed.Bytes(), nil
}

func CreateDiagramFromDacFile(inputfile string, outputfile *string, isGoTemplate bool, overrideDefFile string) {

	log.Infof("input file path: %s\n", inputfile)

	var template TemplateStruct

	// Get the template content
	data, err := getTemplate(inputfile)
	if err != nil {
		log.Fatal(err)
	}

	// Process the template with variables
	var processedData []byte
	if isGoTemplate {
		processedData, err = processTemplate(data)
		if processedData != nil {
			log.Infof("processed template: \n%s", string(processedData))
		}
		if err != nil {
			log.Fatal(err)
		}
	} else {
		processedData = data
	}

	// Unmarshal the processed YAML
	dec := yaml.NewDecoder(bytes.NewReader(processedData))
	dec.KnownFields(true)
	err = dec.Decode(&template)
	if err != nil {
		if ! isGoTemplate && slices.Contains(processedData, '{') {
			log.Warn("Is this file a template, containing template control syntax such as {{ that according to text/template package? If so, add the -t (--tempate) option.")
		}
		log.Fatal(err)
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
