package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	_ "image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/awslabs/diagram-as-code/internal/ctl"
	log "github.com/sirupsen/logrus"
)

var tmpOutputFilename = "/tmp/results/output.png"
var tmpOutputDiffFilename = "/tmp/results/diff-image.png"

func abs(x uint8) uint8 {
	if x < 0 {
		return -x
	}
	return x
}

func subColor(px1, px2 color.NRGBA) color.RGBA {
	if px1.A == px2.A {
		return color.RGBA{
			abs(px1.R - px2.R),
			abs(px1.G - px2.G),
			abs(px1.B - px2.B),
			255,
		}
	} else {
		return color.RGBA{
			abs(px1.A - px2.A),
			abs(px1.A - px2.A),
			abs(px1.A - px2.A),
			255,
		}
	}
}

func compareTwoImages(imageFilePath1, imageFilePath2 string) error {
	fmt.Printf("Compare images %s and %s\n", imageFilePath1, imageFilePath2)
	imageFile1, err := os.Open(imageFilePath1)
	if err != nil {
		return fmt.Errorf("Cannot open imageFilePath1(%s): %v", imageFilePath1, err)
	}
	defer imageFile1.Close()
	img1, _, err := image.Decode(imageFile1)
	if err != nil {
		return fmt.Errorf("Cannot decode imageFile1: %v", err)
	}

	imageFile2, err := os.Open(imageFilePath2)
	if err != nil {
		return fmt.Errorf("Cannot open imageFilePath2(%s): %v", imageFilePath2, err)
	}
	defer imageFile2.Close()
	img2, _, err := image.Decode(imageFile2)
	if err != nil {
		return fmt.Errorf("Cannot decode imageFile2: %v", err)
	}

	// Check image bounds
	if img1.Bounds() != img2.Bounds() {
		return fmt.Errorf("Image bounds mismatch: %v != %v", img1.Bounds(), img2.Bounds())
	}
	fmt.Println("Bounds OK")

	// Generate diff-image from two images
	pixels_diff_numer := 0
	img1b := img1.Bounds()
	img3 := image.NewRGBA(img1b)
	for x := 0; x < img1b.Max.X; x++ {
		for y := 0; y < img1b.Max.Y; y++ {
			px1 := img1.At(x, y)
			px2 := img2.At(x, y)
			img3.Set(x, y, subColor(px1.(color.NRGBA), px2.(color.NRGBA)))
			if px1 != px2 {
				pixels_diff_numer++
			}
		}
	}

	err = os.MkdirAll(filepath.Dir(tmpOutputDiffFilename), os.ModePerm)
	if err != nil {
		return fmt.Errorf("Cannot create directory(%s): %v", filepath.Dir(tmpOutputDiffFilename), err)
	}
	imageFile3, err := os.OpenFile(tmpOutputDiffFilename, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return fmt.Errorf("Cannot open ")
	}
	defer imageFile3.Close()
	png.Encode(imageFile3, img3)

	if pixels_diff_numer > 0 {
		return fmt.Errorf("Mismatch pixels on image %d of %d. See diff-image.png", pixels_diff_numer, img1b.Max.X*img1b.Max.Y)
	}
	fmt.Println("The generated image is an exact match")

	return nil
}

func TestFunctionality(t *testing.T) {
	log.SetLevel(log.WarnLevel)
	files, err := ioutil.ReadDir("../examples")
	if err != nil {
		t.Errorf("Cannot open examples directory, %v", err)
	}
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".yaml" {
			yamlFilename := fmt.Sprintf("../examples/%s", file.Name())
			err = os.MkdirAll(filepath.Dir(tmpOutputFilename), os.ModePerm)
			if err != nil {
				t.Errorf("Cannot create directory(%s): %v", filepath.Dir(tmpOutputFilename), err)
			}
			if strings.HasSuffix(file.Name(), "-cfn.yaml") {
				ctl.CreateDiagramFromCFnTemplate(yamlFilename, &tmpOutputFilename, true)
			} else {
				ctl.CreateDiagramFromDacFile(yamlFilename, &tmpOutputFilename)
			}
			pngFilename := strings.Replace(yamlFilename, ".yaml", ".png", 1)
			err := compareTwoImages(pngFilename, tmpOutputFilename)
			if err != nil {
				t.Errorf("Image mismatch: %v", err)
				break
			}
		}
	}
}
