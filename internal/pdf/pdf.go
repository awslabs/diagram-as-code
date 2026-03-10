package pdf

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"strings"
)

// FromPNG wraps a PNG image into a single-page PDF.
func FromPNG(pngData []byte) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(pngData))
	if err != nil {
		return nil, fmt.Errorf("decode png: %w", err)
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("invalid image size")
	}

	var raw bytes.Buffer
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rgba := color.RGBAModel.Convert(img.At(x, y)).(color.RGBA)
			alpha := float64(rgba.A) / 255.0
			raw.WriteByte(flatten(rgba.R, alpha))
			raw.WriteByte(flatten(rgba.G, alpha))
			raw.WriteByte(flatten(rgba.B, alpha))
		}
	}

	var compressed bytes.Buffer
	zw := zlib.NewWriter(&compressed)
	if _, err := zw.Write(raw.Bytes()); err != nil {
		return nil, fmt.Errorf("compress image: %w", err)
	}
	if err := zw.Close(); err != nil {
		return nil, fmt.Errorf("close compressor: %w", err)
	}

	objects := [][]byte{
		[]byte("<< /Type /Catalog /Pages 2 0 R >>"),
		[]byte("<< /Type /Pages /Kids [3 0 R] /Count 1 >>"),
		[]byte(fmt.Sprintf("<< /Type /Page /Parent 2 0 R /MediaBox [0 0 %d %d] /Resources << /XObject << /Im0 5 0 R >> >> /Contents 4 0 R >>", width, height)),
		streamObject([]byte(fmt.Sprintf("q\n%d 0 0 %d 0 0 cm\n/Im0 Do\nQ\n", width, height))),
		streamObjectWithDict(
			fmt.Sprintf("<< /Type /XObject /Subtype /Image /Width %d /Height %d /ColorSpace /DeviceRGB /BitsPerComponent 8 /Filter /FlateDecode /Length %d >>", width, height, compressed.Len()),
			compressed.Bytes(),
		),
	}

	var buf bytes.Buffer
	buf.WriteString("%PDF-1.4\n")
	buf.WriteString("%\xFF\xFF\xFF\xFF\n")

	offsets := make([]int, 0, len(objects)+1)
	offsets = append(offsets, 0)
	for index, obj := range objects {
		offsets = append(offsets, buf.Len())
		fmt.Fprintf(&buf, "%d 0 obj\n", index+1)
		buf.Write(obj)
		buf.WriteString("\nendobj\n")
	}

	xrefOffset := buf.Len()
	fmt.Fprintf(&buf, "xref\n0 %d\n", len(objects)+1)
	buf.WriteString("0000000000 65535 f \n")
	for _, offset := range offsets[1:] {
		fmt.Fprintf(&buf, "%010d 00000 n \n", offset)
	}
	fmt.Fprintf(&buf, "trailer\n<< /Size %d /Root 1 0 R >>\nstartxref\n%d\n%%%%EOF", len(objects)+1, xrefOffset)

	return buf.Bytes(), nil
}

func flatten(channel uint8, alpha float64) uint8 {
	return uint8(alpha*float64(channel) + (1.0-alpha)*255.0)
}

func streamObject(data []byte) []byte {
	return streamObjectWithDict(fmt.Sprintf("<< /Length %d >>", len(data)), data)
}

func streamObjectWithDict(dict string, data []byte) []byte {
	var buf bytes.Buffer
	buf.WriteString(dict)
	buf.WriteString("\nstream\n")
	buf.Write(data)
	if !strings.HasSuffix(buf.String(), "\n") {
		buf.WriteByte('\n')
	}
	buf.WriteString("endstream")
	return buf.Bytes()
}
