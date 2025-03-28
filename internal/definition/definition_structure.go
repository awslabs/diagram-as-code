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
	for k := range b.Definitions {
		v := b.Definitions[k]
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
			sourceDefinition := b.Definitions[src]
			v.Parent = sourceDefinition
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
		v := b.Definitions[k]
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
					if b.Definitions[v.ZipFile.Source].CacheFilePath == "" {
						q = append(q, k)
						break
					}
					filePath = fmt.Sprintf("%s/%s", b.Definitions[v.ZipFile.Source].CacheFilePath, strings.TrimSuffix(v.ZipFile.Path, "/"))
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
				if b.Definitions[v.Directory.Source].CacheFilePath == "" {
					q = append(q, k)
					break
				}
			}
			v.CacheFilePath = fmt.Sprintf("%s/%s", b.Definitions[v.Directory.Source].CacheFilePath, trimmedPath)
		case "Resource", "Preset", "Group":
			if v.Icon != nil {
				if v.Icon.Path == "" {
					break
				}
				if v.Icon.Source != "" {
					if b.Definitions[v.Icon.Source].CacheFilePath == "" {
						q = append(q, k)
						break
					}
					v.CacheFilePath = fmt.Sprintf("%s/%s", b.Definitions[v.Icon.Source].CacheFilePath, v.Icon.Path)
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
