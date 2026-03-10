package ctl

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/png"
	"os"
	"os/exec"
	"time"

	dacpdf "github.com/awslabs/diagram-as-code/internal/pdf"
)

var (
	drawioLookPath = exec.LookPath
	drawioRun = func(ctx context.Context, bin string, args ...string) error {
		cmd := exec.CommandContext(ctx, bin, args...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			if len(output) == 0 {
				return err
			}
			return fmt.Errorf("%w: %s", err, string(output))
		}
		return nil
	}
	drawioGenerator = CreateDrawioFromDacFile
)

func findDrawioBinary() (string, error) {
	if bin, err := drawioLookPath("drawio"); err == nil {
		return bin, nil
	}
	const macAppBinary = "/Applications/draw.io.app/Contents/MacOS/draw.io"
	if _, err := os.Stat(macAppBinary); err == nil {
		return macAppBinary, nil
	}
	return "", fmt.Errorf("draw.io binary not found")
}

func canUseDrawioExport() bool {
	_, err := findDrawioBinary()
	return err == nil
}

func exportDrawioFile(inputfile, outputfile, format string) error {
	bin, err := findDrawioBinary()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	args := []string{"--export", "--format", format, "--output", outputfile, inputfile}
	if err := drawioRun(ctx, bin, args...); err != nil {
		return fmt.Errorf("draw.io export failed: %w", err)
	}
	return nil
}

func resizePNGFile(path string, width, height int) error {
	if width == 0 && height == 0 {
		return nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read exported png: %w", err)
	}
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to decode exported png: %w", err)
	}

	bounds := img.Bounds()
	rgba := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rgba.Set(x-bounds.Min.X, y-bounds.Min.Y, img.At(x, y))
		}
	}

	resized := resizeImage(rgba, width, height)

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to open resized png output: %w", err)
	}
	defer f.Close()

	if err := png.Encode(f, resized); err != nil {
		return fmt.Errorf("failed to encode resized png: %w", err)
	}
	return nil
}

func createDiagramViaDrawio(inputfile string, outputfile *string, opts *CreateOptions) error {
	if outputfile == nil {
		return fmt.Errorf("output file is required")
	}
	if opts == nil {
		opts = &CreateOptions{}
	}
	if err := CheckOutputFileOverwrite(*outputfile, opts.OverwriteMode); err != nil {
		return err
	}

	tmpDrawio, err := os.CreateTemp("", "dac-export-*.drawio")
	if err != nil {
		return fmt.Errorf("failed to create temp drawio file: %w", err)
	}
	tmpDrawio.Close()
	defer os.Remove(tmpDrawio.Name())

	tmpDrawioPath := tmpDrawio.Name()
	if err := drawioGenerator(inputfile, &tmpDrawioPath, opts); err != nil {
		return fmt.Errorf("failed to create drawio source: %w", err)
	}
	if err := exportDrawioFile(tmpDrawioPath, *outputfile, "png"); err != nil {
		return err
	}
	return resizePNGFile(*outputfile, opts.Width, opts.Height)
}

func createPDFViaDrawio(inputfile string, outputfile *string, opts *CreateOptions) error {
	if outputfile == nil {
		return fmt.Errorf("output file is required")
	}
	if opts == nil {
		opts = &CreateOptions{}
	}
	if err := CheckOutputFileOverwrite(*outputfile, opts.OverwriteMode); err != nil {
		return err
	}

	tmpDrawio, err := os.CreateTemp("", "dac-export-*.drawio")
	if err != nil {
		return fmt.Errorf("failed to create temp drawio file: %w", err)
	}
	tmpDrawio.Close()
	defer os.Remove(tmpDrawio.Name())

	tmpDrawioPath := tmpDrawio.Name()
	if err := drawioGenerator(inputfile, &tmpDrawioPath, opts); err != nil {
		return fmt.Errorf("failed to create drawio source: %w", err)
	}

	if opts.Width == 0 && opts.Height == 0 {
		return exportDrawioFile(tmpDrawioPath, *outputfile, "pdf")
	}

	tmpPNG, err := os.CreateTemp("", "dac-export-*.png")
	if err != nil {
		return fmt.Errorf("failed to create temp png file: %w", err)
	}
	tmpPNG.Close()
	defer os.Remove(tmpPNG.Name())

	if err := exportDrawioFile(tmpDrawioPath, tmpPNG.Name(), "png"); err != nil {
		return err
	}
	if err := resizePNGFile(tmpPNG.Name(), opts.Width, opts.Height); err != nil {
		return err
	}
	pngData, err := os.ReadFile(tmpPNG.Name())
	if err != nil {
		return fmt.Errorf("failed to read drawio png export: %w", err)
	}
	pdfData, err := dacpdf.FromPNG(pngData)
	if err != nil {
		return fmt.Errorf("failed to convert drawio png to pdf: %w", err)
	}
	if err := os.WriteFile(*outputfile, pdfData, 0600); err != nil {
		return fmt.Errorf("failed to write pdf output: %w", err)
	}
	return nil
}
