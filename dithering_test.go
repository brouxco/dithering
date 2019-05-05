package dithering

import (
	"image"
	"image/color"
	"testing"
)

// TestDither_DrawCheckerboard tests that the result of Floyd-Steinberg
// error diffusion of a uniform 50% gray source image with a black-and-white
// palette is a checkerboard pattern.
func TestDither_DrawCheckerboard(t *testing.T) {
	b := image.Rect(0, 0, 640, 480)
	src := &image.Uniform{color.Gray{0x7f}}
	dst := image.NewPaletted(b, color.Palette{color.Black, color.White})

	floydSteinberg := NewDither(FloydSteinberg)
	floydSteinberg.Draw(dst, b, src, image.ZP)

	nErr := 0
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			got := dst.Pix[dst.PixOffset(x, y)]
			want := uint8(x+y) % 2
			if got != want {
				t.Errorf("at (%d, %d): got %d, want %d", x, y, got, want)
				if nErr++; nErr == 10 {
					t.Fatal("there may be more errors")
				}
			}
		}
	}
}
