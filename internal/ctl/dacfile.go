// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package ctl

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	tmpl "text/template"

	"github.com/awslabs/diagram-as-code/internal/definition"
	"github.com/awslabs/diagram-as-code/internal/types"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
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
		defer func() {
			if closeErr := resp.Body.Close(); closeErr != nil {
				log.Warnf("Failed to close response body: %v", closeErr)
			}
		}()

		data, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
	} else {
		// Local file
		data, err = os.ReadFile(inputfile)
		if err != nil {
			return nil, err
		}
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

func CreateDiagramFromDacFile(inputfile string, outputfile *string, opts *CreateOptions) error {

	log.Infof("input file path: %s\n", inputfile)

	var template TemplateStruct

	// Get the template content
	data, err := getTemplate(inputfile)
	if err != nil {
		return fmt.Errorf("failed to get template: %w", err)
	}

	// Process the template with variables
	var processedData []byte
	if opts.IsGoTemplate {
		processedData, err = processTemplate(data)
		if processedData != nil {
			log.Infof("processed template: \n%s", string(processedData))
		}
		if err != nil {
			return fmt.Errorf("failed to process template: %w", err)
		}
	} else {
		processedData = data
	}

	// Unmarshal the processed YAML
	dec := yaml.NewDecoder(bytes.NewReader(processedData))
	dec.KnownFields(true)
	err = dec.Decode(&template)
	if err != nil {
		if !opts.IsGoTemplate && slices.Contains(processedData, '{') {
			log.Warn("Is this file a template, containing template control syntax such as {{ that according to text/template package? If so, add the -t (--tempate) option.")
		}
		return fmt.Errorf("failed to decode YAML: %w", err)
	}

	var ds definition.DefinitionStructure
	resources := make(map[string]*types.Resource)

	log.Info("Load DefinitionFiles section")
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
		// OverrideDefFile is for testing, so allow untrusted URLs
		if err := loadDefinitionFiles(&overrideDefTemplate, &ds, true); err != nil {
			return fmt.Errorf("failed to load override definition files: %w", err)
		}
		log.Infof("overrideDefTemplate: %+v", overrideDefTemplate)
	} else {
		if err := loadDefinitionFiles(&template, &ds, opts.AllowUntrustedDefinitions); err != nil {
			return fmt.Errorf("failed to load definition files: %w", err)
		}
	}

	log.Info("Load Resources section")
	if err := loadResources(&template, ds, resources); err != nil {
		return fmt.Errorf("failed to load resources: %w", err)
	}

	log.Info("Associate children with parent resources")
	if err := associateChildren(&template, resources); err != nil {
		return fmt.Errorf("failed to associate children: %w", err)
	}

	// Check for unused resources
	checkUnusedResources(&template)

	log.Info("Add Links section")
	if err := loadLinks(&template, resources); err != nil {
		return fmt.Errorf("failed to load links: %w", err)
	}

	// Reorder children based on links (UnorderedChildren feature)
	log.Info("Reorder children based on links")
	canvas, exists := resources["Canvas"]
	if exists {
		// Collect all links from all resources
		var allLinks []*types.Link
		for _, resource := range resources {
			allLinks = append(allLinks, resource.GetLinks()...)
		}
		types.ReorderChildrenByLinks(canvas, allLinks)
	}

	if err := createDiagram(resources, outputfile, opts); err != nil {
		return fmt.Errorf("failed to create diagram: %w", err)
	}
	return nil
}
