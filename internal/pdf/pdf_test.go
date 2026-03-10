package pdf

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"
)

func TestFromPNG(t *testing.T) {
	img := image.NewNRGBA(image.Rect(0, 0, 2, 2))
	img.Set(0, 0, color.NRGBA{R: 255, A: 255})
	img.Set(1, 0, color.NRGBA{G: 255, A: 255})
	img.Set(0, 1, color.NRGBA{B: 255, A: 255})
	img.Set(1, 1, color.NRGBA{R: 255, G: 255, A: 128})

	var pngBuf bytes.Buffer
	if err := png.Encode(&pngBuf, img); err != nil {
		t.Fatalf("encode png: %v", err)
	}

	pdfData, err := FromPNG(pngBuf.Bytes())
	if err != nil {
		t.Fatalf("FromPNG: %v", err)
	}

	if !bytes.HasPrefix(pdfData, []byte("%PDF-1.4")) {
		t.Fatalf("expected pdf header, got %q", pdfData[:8])
	}
	if !bytes.Contains(pdfData, []byte("/Subtype /Image")) {
		t.Fatalf("expected image object in pdf")
	}
	if !bytes.Contains(pdfData, []byte("%%EOF")) {
		t.Fatalf("expected eof marker")
	}
}
