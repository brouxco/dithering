package dithering

import (
	"image"
	"image/color"
)

// PixelError represents the error for each canal in the image
// when dithering an image
// Errors are floats because they are the result of a division
type PixelError struct {
	R, G, B, A float32
}

func (c PixelError) RGBA() (r, g, b, a uint32) {
	return uint32(c.R), uint32(c.G), uint32(c.B), uint32(c.A)
}

func pixelErrorModel(c color.Color) color.Color {
	if _, ok := c.(PixelError); ok {
		return c
	}
	r, g, b, a := c.RGBA()
	return PixelError{float32(r), float32(g), float32(b), float32(a)}
}

// ErrorImage is an in-memory image whose At method returns dithering.PixelError values
type ErrorImage struct {
	// Pix holds the image's pixels, in R, G, B, A order. The pixel at
	// (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*4].
	Pix []float32
	// Stride is the Pix stride (in bytes) between vertically adjacent pixels.
	Stride int
	// Rect is the image's bounds.
	Rect image.Rectangle
}

func (p *ErrorImage) ColorModel() color.Model {
	return color.ModelFunc(pixelErrorModel)
}

func (p *ErrorImage) Bounds() image.Rectangle { return p.Rect }

func (p *ErrorImage) At(x, y int) color.Color {
	return p.PixelErrorAt(x, y)
}

func (p *ErrorImage) PixelErrorAt(x, y int) PixelError {
	if !(image.Point{x, y}.In(p.Rect)) {
		return PixelError{}
	}
	i := p.PixOffset(x, y)
	return PixelError{p.Pix[i+0], p.Pix[i+1], p.Pix[i+2], p.Pix[i+3]}
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (p *ErrorImage) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*4
}

func (p *ErrorImage) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	c1 := color.ModelFunc(pixelErrorModel).Convert(c).(PixelError)
	p.Pix[i+0] = c1.R
	p.Pix[i+1] = c1.G
	p.Pix[i+2] = c1.B
	p.Pix[i+3] = c1.A
}

func (p *ErrorImage) SetPixelError(x, y int, c PixelError) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	p.Pix[i+0] = c.R
	p.Pix[i+1] = c.G
	p.Pix[i+2] = c.B
	p.Pix[i+3] = c.A
}

func NewErrorImage(r image.Rectangle) *ErrorImage {
	w, h := r.Dx(), r.Dy()
	buf := make([]float32, 4*w*h)
	return &ErrorImage{buf, 4 * w, r}
}
