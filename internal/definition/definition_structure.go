// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package definition

import (
	"fmt"
	"os"
	"strings"

	"github.com/awslabs/diagram-as-code/internal/cache"
	"golang.org/x/exp/maps"
	"gopkg.in/yaml.v3"
)

type DefinitionStructure struct {
	Definitions map[string]*Definition `yaml:"Definitions"`
}

func (ds *DefinitionStructure) LoadDefinitions(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("cannot open Definition File(%s): %v", filePath, err)
	}

	var b DefinitionStructure

	err = yaml.Unmarshal([]byte(data), &b)
	if err != nil {
		return fmt.Errorf("cannot yaml.Unmarshal Definition File(%s): %v", filePath, err)
	}

	// Linking definitions
	for k, v := range b.Definitions {
		src := func() string {
			switch v.Type {
			case "Resource", "Preset":
				if v.Icon != nil {
					return v.Icon.Source
				}
			case "Directory":
				return v.Directory.Source
			case "Zip":
				return v.ZipFile.Source
			}
			return ""
		}()
		if src != "" {
			sourceDefinition, exists := b.Definitions[src]
			if exists {
				v.Parent = sourceDefinition
			}
		}
		b.Definitions[k] = v
	}

	// Downdload files and extract ZIP
	q := []string{}
	for k := range b.Definitions {
		q = append(q, k)
	}
	for len(q) > 0 {
		k := q[0]
		v, ok := b.Definitions[k]
		if !ok {
			return fmt.Errorf("definition key %s not found in definitions map", k)
		}
		switch v.Type {
		case "Zip":
			switch v.ZipFile.SourceType {
			case "url":
				if v.ZipFile.Url == "" {
					return fmt.Errorf("Zip(url) needs ZipFile.URL")
				}
				filePath, err := cache.FetchFile(v.ZipFile.Url)
				if err != nil {
					return fmt.Errorf("cannot FetchFile(%s): %v", v.ZipFile.Url, err)
				}
				v.CacheFilePath, err = cache.ExtractZipFile(filePath)
				if err != nil {
					return fmt.Errorf("cannot ExtractZipFile(%s): %v", filePath, err)
				}

			case "file":
				filePath := ""
				if v.ZipFile.Source == "" {
					//return fmt.Errorf("Zip(file) needs ZipFile.Source")
					filePath = strings.TrimSuffix(v.ZipFile.Path, "/")
				} else {
					sourceDef, ok := b.Definitions[v.ZipFile.Source]
					if !ok {
						return fmt.Errorf("ZipFile source %s not found in definitions", v.ZipFile.Source)
					}
					if sourceDef.CacheFilePath == "" {
						q = append(q, k)
						break
					}
					filePath = fmt.Sprintf("%s/%s", sourceDef.CacheFilePath, strings.TrimSuffix(v.ZipFile.Path, "/"))
				}
				v.CacheFilePath, err = cache.ExtractZipFile(filePath)
				if err != nil {
					return fmt.Errorf("cannot ExtractZipFile(%s): %v", filePath, err)
				}
			}
		case "Directory":
			trimmedPath := strings.TrimSuffix(v.Directory.Path, "/")
			if trimmedPath == "" {
				//lint:ignore ST1005 Directory is a proper noun, so will ignore capitalization rule.
				return fmt.Errorf("Directory %s has only slash or empty path", q)
			}
			if v.Directory.Source != "" {
				sourceDef, ok := b.Definitions[v.Directory.Source]
				if !ok {
					return fmt.Errorf("Directory source %s not found in definitions", v.Directory.Source)
				}
				if sourceDef.CacheFilePath == "" {
					q = append(q, k)
					break
				}
				v.CacheFilePath = fmt.Sprintf("%s/%s", sourceDef.CacheFilePath, trimmedPath)
			} else {
				v.CacheFilePath = trimmedPath
			}
		case "Resource", "Preset", "Group":
			if v.Icon != nil {
				if v.Icon.Path == "" {
					break
				}
				if v.Icon.Source != "" {
					sourceDef, ok := b.Definitions[v.Icon.Source]
					if !ok {
						return fmt.Errorf("Icon source %s not found in definitions", v.Icon.Source)
					}
					if sourceDef.CacheFilePath == "" {
						q = append(q, k)
						break
					}
					v.CacheFilePath = fmt.Sprintf("%s/%s", sourceDef.CacheFilePath, v.Icon.Path)
				}
			}
		}
		q = q[1:]
	}
	if ds.Definitions == nil {
		ds.Definitions = map[string]*Definition{}
	}
	maps.Copy(ds.Definitions, b.Definitions)
	return nil
}
