package ctl

import (
	"fmt"
	"os"

	dacpdf "github.com/awslabs/diagram-as-code/internal/pdf"
)

func CreatePDFFromDacFile(inputfile string, outputfile *string, opts *CreateOptions) error {
	return createPDF(func(tmpOutput string) error {
		return CreateDiagramFromDacFile(inputfile, &tmpOutput, opts)
	}, outputfile, opts)
}

func CreatePDFFromCFnTemplate(inputfile string, outputfile *string, generateDacFile bool, opts *CreateOptions) error {
	return createPDF(func(tmpOutput string) error {
		return CreateDiagramFromCFnTemplate(inputfile, &tmpOutput, generateDacFile, opts)
	}, outputfile, opts)
}

func createPDF(generatePNG func(tmpOutput string) error, outputfile *string, opts *CreateOptions) error {
	if outputfile == nil {
		return fmt.Errorf("output file is required")
	}
	if generatePNG == nil {
		return fmt.Errorf("png generator is required")
	}
	if opts == nil {
		opts = &CreateOptions{}
	}

	tmpPNG, err := os.CreateTemp("", "dac-pdf-*.png")
	if err != nil {
		return fmt.Errorf("failed to create temp png file: %w", err)
	}
	tmpPNG.Close()
	defer os.Remove(tmpPNG.Name())

	if err := generatePNG(tmpPNG.Name()); err != nil {
		return err
	}

	pngData, err := os.ReadFile(tmpPNG.Name())
	if err != nil {
		return fmt.Errorf("failed to read generated png: %w", err)
	}

	pdfData, err := dacpdf.FromPNG(pngData)
	if err != nil {
		return fmt.Errorf("failed to convert png to pdf: %w", err)
	}

	if err := CheckOutputFileOverwrite(*outputfile, opts.OverwriteMode); err != nil {
		return err
	}

	if err := os.WriteFile(*outputfile, pdfData, 0600); err != nil {
		return fmt.Errorf("failed to write pdf output: %w", err)
	}

	return nil
}
