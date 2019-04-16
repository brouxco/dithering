// Package dithering provides a customizable image ditherer
package dithering

import (
	"image"
	"image/color"
	"image/draw"
)

// Dither represent dithering algorithm implementation
type Dither struct {
	// Matrix is the error diffusion matrix
	Matrix [][]float32
	// TODO(brouxco): the shift should not be necessary, the algorithm could determine it automatically given the matrix
	Shift  int
}

// abs gives the absolute value of a signed integer
func abs(x int16) uint16 {
	if x < 0 {
		return uint16(-x)
	}
	return uint16(x)
}

// findColor determines the closest color in a palette given the pixel color and the error
//
// It returns the closest color, the updated error and the distance between the error and the color
func findColor(err color.Color, pix color.Color, pal color.Palette) (color.RGBA, PixelError, uint16) {
	var errR, errG, errB,
		pixR, pixG, pixB,
		colR, colG, colB int16
	_errR, _errG, _errB, _ := err.RGBA()
	_pixR, _pixG, _pixB, _ := pix.RGBA()

	// Low-pass filter
	errR = int16(float32(int16(_errR)) * 0.75)
	errG = int16(float32(int16(_errG)) * 0.75)
	errB = int16(float32(int16(_errB)) * 0.75)

	pixR = int16(uint8(_pixR)) + errR
	pixG = int16(uint8(_pixG)) + errG
	pixB = int16(uint8(_pixB)) + errB

	var index int
	var minDiff uint16 = 1<<16 - 1

	for i, col := range pal {
		_colR, _colG, _colB, _ := col.RGBA()

		colR = int16(uint8(_colR))
		colG = int16(uint8(_colG))
		colB = int16(uint8(_colB))
		var distance = abs(pixR-colR) + abs(pixG-colG) + abs(pixB-colB)

		if distance < minDiff {
			index = i
			minDiff = distance
		}
	}

	_colR, _colG, _colB, _ := pal[index].RGBA()

	colR = int16(uint8(_colR))
	colG = int16(uint8(_colG))
	colB = int16(uint8(_colB))

	return color.RGBA{uint8(colR), uint8(colG), uint8(colB), 255},
		PixelError{float32(pixR - colR),
			float32(pixG - colG),
			float32(pixB - colB),
			1<<16 - 1},
			minDiff
}

// Draw applies an error diffusion algorithm to the src image
func (dit Dither) Draw(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point) {
	if _, ok := dst.(*image.Paletted); !ok {
		return
	}
	p := dst.(*image.Paletted).Palette

	err := NewErrorImage(r)
	var diff uint64

	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			// using the closest color
			r, e, d := findColor(err.PixelErrorAt(x, y), src.At(x, y), p)
			dst.Set(x, y, r)
			err.SetPixelError(x, y, e)
			diff += uint64(d)

			// diffusing the error using the diffusion matrix
			for i, v1 := range dit.Matrix {
				for j, v2 := range v1 {
					err.SetPixelError(x+j+dit.Shift, y+i,
						err.PixelErrorAt(x+j+dit.Shift, y+i).Add(err.PixelErrorAt(x, y).Mul(v2)))
				}
			}
		}
	}
}
