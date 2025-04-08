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
				error_message := "awsdac: This tool requires an input file to run. Please provide a file path.\n"
				fmt.Println(error_message)
				cmd.Help()

				os.Exit(1)
			}

			inputFile := args[0]
			if !ctl.IsURL(inputFile) {

				if _, err := os.Stat(inputFile); os.IsNotExist(err) {
					fmt.Printf("awsdac: Input file '%s' does not exist.\n", inputFile)
					os.Exit(1)
				}
			}

			return nil

		},
		Run: func(cmd *cobra.Command, args []string) {

			if verbose {
				log.SetLevel(log.InfoLevel)
			} else {
				log.SetLevel(log.WarnLevel)
			}

			inputFile := args[0]

			if cfnTemplate {
				ctl.CreateDiagramFromCFnTemplate(inputFile, &outputFile, generateDacFile)
			} else {
				ctl.CreateDiagramFromDacFile(inputFile, &outputFile, isGoTemplate, overrideDefFile)
			}

		},
	}

	rootCmd.PersistentFlags().StringVarP(&outputFile, "output", "o", "output.png", "Output file name")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
	rootCmd.PersistentFlags().BoolVarP(&cfnTemplate, "cfn-template", "c", false, "[beta] Create diagram from CloudFormation template")
	rootCmd.PersistentFlags().BoolVarP(&generateDacFile, "dac-file", "d", false, "[beta] Generate YAML file in dac (diagram-as-code) format from CloudFormation template")
	rootCmd.PersistentFlags().StringVarP(&overrideDefFile, "override-def-file", "", "", "For testing purpose, override DefinitionFiles to another url/local file")
	rootCmd.PersistentFlags().BoolVarP(&isGoTemplate, "template", "t", false, "Processes the input file as a template according to text/template.")


	rootCmd.Execute()
}
