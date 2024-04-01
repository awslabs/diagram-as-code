// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package definition

type Definition struct {
	Type          string              `yaml:"Type"`
	Icon          *DefinitionIcon     `yaml:"Icon"`
	Label         *DefinitionLabel    `yaml:"Label"`
	Fill          *DefinitionFill     `yaml:"Fill"`
	Border        *DefinitionBorder   `yaml:"Border"`
	Directory     DefinitionDirectory `yaml:"Directory"`
	ZipFile       DefinitionZipFile   `yaml:"ZipFile"`
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
