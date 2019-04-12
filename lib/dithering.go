package dithering

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"os"

	_ "image/jpeg"
)

func bound(x int16) int16 {
	if x > 255 {
		return 255
	} else if x < -255 {
		return -255
	} else {
		return x
	}
}

func abs(x int16) uint16 {
	if x < 0 {
		return uint16(-x)
	}
	return uint16(x)
}

func loadImage(filename string) image.Image {
	reader, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	img, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	return img
}

func storeImage(filename string, img image.Image) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if err = png.Encode(file, img); err != nil {
		log.Fatal(err)
	}
}

func findColor(err color.Color, pix color.Color, pal color.Palette) (color.RGBA, PixelError) {
	var errR, errG, errB,
		pixR, pixG, pixB,
		colR, colG, colB int16
	_errR, _errG, _errB, _ := err.RGBA()
	_pixR, _pixG, _pixB, _ := pix.RGBA()

	errR = int16(_errR)
	errG = int16(_errG)
	errB = int16(_errB)
	pixR = int16(uint8(_pixR))
	pixG = int16(uint8(_pixG))
	pixB = int16(uint8(_pixB))

	var index int
	var minDiff uint16 = 1<<16 - 1

	for i, col := range pal {
		_colR, _colG, _colB, _ := col.RGBA()

		colR = int16(uint8(_colR))
		colG = int16(uint8(_colG))
		colB = int16(uint8(_colB))
		var distance = abs(colR+errR-pixR) + abs(colG+errG-pixG) + abs(colB+errB-pixB)

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
		PixelError{float32(bound((colR - pixR) + errR)),
			float32(bound((colG - pixG) + errG)),
			float32(bound((colB - pixB) + errB)),
			1<<16 - 1}
}

func Dither(input string, output string) {
	img := loadImage(input)

	bounds := img.Bounds()

	p := color.Palette{
		color.RGBA{0, 0, 0, 255},
		color.RGBA{255, 255, 255, 255},
		color.RGBA{255, 0, 0, 255},
		color.RGBA{0, 255, 0, 255},
		color.RGBA{0, 0, 255, 255},
	}

	result := image.NewRGBA(bounds)
	err := NewErrorImage(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// using the closest color
			// TODO: palettes
			r, e := findColor(err.PixelErrorAt(x, y), img.At(x, y), p)
			result.SetRGBA(x, y, r)
			err.SetPixelError(x, y, e)

			// diffusing the error using the diffusion matrix
			// TODO: matrices
			if x+1 < bounds.Max.X {
				err.SetPixelError(x+1, y, err.PixelErrorAt(x, y))
			}
		}
	}

	storeImage(output, result)
}
