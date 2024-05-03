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

func main() {

	var outputFile string
	var verbose bool
	var cfntemplate bool

	var rootCmd = &cobra.Command{
		Use:   "awsdac <input filename>",
		Short: "Diagram-as-code for AWS architecture.",
		Long:  "This command line interface (CLI) tool enables drawing infrastructure diagrams for Amazon Web Services through YAML code.",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			if len(args) == 0 {
				fmt.Println("Error: This tool requires an input file to run. Please provide a file path.")
				os.Exit(1)
			}

			if verbose {
				log.SetLevel(log.InfoLevel)
			} else {
				log.SetLevel(log.WarnLevel)
			}

			inputFile := args[0]

			if cfntemplate {
				ctl.CreateDiagramFromCFnTemplate(inputFile, &outputFile)
			} else {
				ctl.CreateDiagramFromYAML(inputFile, &outputFile)
			}

		},
	}

	rootCmd.PersistentFlags().StringVarP(&outputFile, "output", "o", "output.png", "Output file name")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
	rootCmd.PersistentFlags().BoolVarP(&cfntemplate, "cfn-template", "", false, "Create diagram from CloudFormation template")

	rootCmd.Execute()
}
