package ctl

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

func TestCreateDiagramFromDacFilePrefersDrawioExport(t *testing.T) {
	originalLookPath := drawioLookPath
	originalRun := drawioRun
	originalGenerator := drawioGenerator
	t.Cleanup(func() {
		drawioLookPath = originalLookPath
		drawioRun = originalRun
		drawioGenerator = originalGenerator
	})

	generatorCalled := false
	drawioLookPath = func(file string) (string, error) {
		return "/usr/bin/drawio", nil
	}
	drawioGenerator = func(inputfile string, outputfile *string, opts *CreateOptions) error {
		generatorCalled = true
		return os.WriteFile(*outputfile, []byte("<mxfile/>"), 0600)
	}
	drawioRun = func(ctx context.Context, bin string, args ...string) error {
		if bin != "/usr/bin/drawio" {
			t.Fatalf("unexpected draw.io binary: %s", bin)
		}
		if len(args) != 6 || args[2] != "png" {
			t.Fatalf("unexpected draw.io args: %v", args)
		}
		return writeTestPNG(args[4])
	}

	inputPath := filepath.Join(t.TempDir(), "input.yaml")
	if err := os.WriteFile(inputPath, []byte("this is not valid yaml"), 0600); err != nil {
		t.Fatalf("failed to write input file: %v", err)
	}

	outputPath := filepath.Join(t.TempDir(), "output.png")
	if err := CreateDiagramFromDacFile(inputPath, &outputPath, &CreateOptions{
		PreferDrawioExport: true,
		OverwriteMode:      Force,
	}); err != nil {
		t.Fatalf("CreateDiagramFromDacFile returned error: %v", err)
	}

	if !generatorCalled {
		t.Fatal("expected draw.io generator to be used")
	}
	if _, err := os.Stat(outputPath); err != nil {
		t.Fatalf("expected output png to exist: %v", err)
	}
}

func TestCreatePDFFromDacFilePrefersDrawioExport(t *testing.T) {
	originalLookPath := drawioLookPath
	originalRun := drawioRun
	originalGenerator := drawioGenerator
	t.Cleanup(func() {
		drawioLookPath = originalLookPath
		drawioRun = originalRun
		drawioGenerator = originalGenerator
	})

	generatorCalled := false
	drawioLookPath = func(file string) (string, error) {
		return "/usr/bin/drawio", nil
	}
	drawioGenerator = func(inputfile string, outputfile *string, opts *CreateOptions) error {
		generatorCalled = true
		return os.WriteFile(*outputfile, []byte("<mxfile/>"), 0600)
	}
	drawioRun = func(ctx context.Context, bin string, args ...string) error {
		if bin != "/usr/bin/drawio" {
			t.Fatalf("unexpected draw.io binary: %s", bin)
		}
		if len(args) != 6 || args[2] != "pdf" {
			t.Fatalf("unexpected draw.io args: %v", args)
		}
		return os.WriteFile(args[4], []byte("%PDF-1.4\n%%EOF\n"), 0600)
	}

	inputPath := filepath.Join(t.TempDir(), "input.yaml")
	if err := os.WriteFile(inputPath, []byte("this is not valid yaml"), 0600); err != nil {
		t.Fatalf("failed to write input file: %v", err)
	}

	outputPath := filepath.Join(t.TempDir(), "output.pdf")
	if err := CreatePDFFromDacFile(inputPath, &outputPath, &CreateOptions{
		PreferDrawioExport: true,
		OverwriteMode:      Force,
	}); err != nil {
		t.Fatalf("CreatePDFFromDacFile returned error: %v", err)
	}

	if !generatorCalled {
		t.Fatal("expected draw.io generator to be used")
	}
	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read output pdf: %v", err)
	}
	if !bytes.HasPrefix(data, []byte("%PDF-1.4")) {
		t.Fatalf("expected pdf output, got: %q", data)
	}
}

func writeTestPNG(path string) error {
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	img.Set(0, 0, color.RGBA{R: 255, A: 255})
	img.Set(1, 0, color.RGBA{G: 255, A: 255})
	img.Set(0, 1, color.RGBA{B: 255, A: 255})
	img.Set(1, 1, color.RGBA{R: 255, G: 255, A: 255})

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}
