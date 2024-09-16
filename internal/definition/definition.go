// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package definition

import "fmt"

type Definition struct {
	Type          string              `yaml:"Type"`
	Icon          *DefinitionIcon     `yaml:"Icon"`
	Label         *DefinitionLabel    `yaml:"Label"`
	Fill          *DefinitionFill     `yaml:"Fill"`
	Border        *DefinitionBorder   `yaml:"Border"`
	Directory     DefinitionDirectory `yaml:"Directory"`
	ZipFile       DefinitionZipFile   `yaml:"ZipFile"`
	CFn           DefinitionCFn       `yaml:"CFn"`
	Parent        *Definition
	CacheFilePath string
}

type DefinitionLabel struct {
	Title string `yaml:"Title"`
	Color string `yaml:"Color"`
	Font  string `yaml:"Font"`
}

type DefinitionFill struct {
	Color string `yaml:"Color"`
}

type DefinitionBorder struct {
	Color string `yaml:"Color"`
	Type  string `yaml:"Type"`
}

// [TODO] make interface
type DefinitionIcon struct {
	Source string `yaml:"Source"`
	Path   string `yaml:"Path"`
}

type DefinitionDirectory struct {
	Source string `yaml:"Source"`
	Path   string `yaml:"Path"`
}

type DefinitionZipFile struct {
	SourceType string `yaml:"SourceType"`
	Source     string `yaml:"Source"`
	Path       string `yaml:"Path"`
	Url        string `yaml:"Url"`
}

type DefinitionCFn struct {
	HasChildren bool `yaml:"HasChildren"`
}

func (d *Definition) String() string {
	res := "Definition{\n"
	if d.Type != "" {
		res += d.Type
	}
	if d.Icon != nil {
		res += fmt.Sprintf("  Icon: %v\n", d.Icon)
	}
	if d.Label != nil {
		res += fmt.Sprintf("  Label: %v\n", d.Label)
	}
	if d.Fill != nil {
		res += fmt.Sprintf("  Fill: %v\n", d.Fill)
	}
	if d.Border != nil {
		res += fmt.Sprintf("  Border: %v\n", d.Border)
	}
	res += fmt.Sprintf("  Directory: %v\n", d.Directory)
	res += fmt.Sprintf("  ZipFile: %v\n", d.ZipFile)
	if d.Parent != nil {
		res += fmt.Sprintf("  Parent: %v\n", d.Parent)
	}
	if d.CacheFilePath != "" {
		res += d.CacheFilePath
	}
	res += "}\n"
	return res
}
