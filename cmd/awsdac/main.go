// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"

	"github.com/awslabs/diagram-as-code/internal/ctl"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var version = "dev"

func main() {

	var outputFile string
	var verbose bool
	var cfnTemplate bool
	var generateDacFile bool
	var overrideDefFile string
	var isGoTemplate bool

	var rootCmd = &cobra.Command{
		Use:     "awsdac <input filename>",
		Version: version,
		Short:   "Diagram-as-code for AWS architecture.",
		Long:    "This command line interface (CLI) tool enables drawing infrastructure diagrams for Amazon Web Services through YAML code.",
		Args:    cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {

			if len(args) == 0 {
				return fmt.Errorf("awsdac: This tool requires an input file to run. Please provide a file path")
			}

			inputFile := args[0]
			if !ctl.IsURL(inputFile) {
				if _, err := os.Stat(inputFile); os.IsNotExist(err) {
					return fmt.Errorf("awsdac: Input file '%s' does not exist", inputFile)
				}
			}

			return nil

		},
		RunE: func(cmd *cobra.Command, args []string) error {

			if verbose {
				log.SetLevel(log.InfoLevel)
			} else {
				log.SetLevel(log.WarnLevel)
			}

			inputFile := args[0]

			if cfnTemplate {
				opts := ctl.CreateOptions{
					OverrideDefFile: overrideDefFile,
				}
				if err := ctl.CreateDiagramFromCFnTemplate(inputFile, &outputFile, generateDacFile, &opts); err != nil {
					return fmt.Errorf("failed to create diagram from CloudFormation template: %w", err)
				}
			} else {
				opts := ctl.CreateOptions{
					IsGoTemplate:    isGoTemplate,
					OverrideDefFile: overrideDefFile,
				}
				if err := ctl.CreateDiagramFromDacFile(inputFile, &outputFile, &opts); err != nil {
					return fmt.Errorf("failed to create diagram: %w", err)
				}
			}

			return nil
		},
	}

	rootCmd.PersistentFlags().StringVarP(&outputFile, "output", "o", "output.png", "Output file name")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
	rootCmd.PersistentFlags().BoolVarP(&cfnTemplate, "cfn-template", "c", false, "[beta] Create diagram from CloudFormation template")
	rootCmd.PersistentFlags().BoolVarP(&generateDacFile, "dac-file", "d", false, "[beta] Generate YAML file in dac (diagram-as-code) format from CloudFormation template")
	rootCmd.PersistentFlags().StringVarP(&overrideDefFile, "override-def-file", "", "", "For testing purpose, override DefinitionFiles to another url/local file")
	rootCmd.PersistentFlags().BoolVarP(&isGoTemplate, "template", "t", false, "Processes the input file as a template according to text/template.")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
