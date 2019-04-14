package dithering

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"

	_ "image/jpeg"
)

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

func findColor(err color.Color, pix color.Color, pal color.Palette) (color.RGBA, PixelError, uint16) {
	var errR, errG, errB,
		pixR, pixG, pixB,
		colR, colG, colB int16
	_errR, _errG, _errB, _ := err.RGBA()
	_pixR, _pixG, _pixB, _ := pix.RGBA()

	errR = int16(_errR)
	errG = int16(_errG)
	errB = int16(_errB)
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

func Dither(input string, output string) {
	img := loadImage(input)

	bounds := img.Bounds()

	p := color.Palette{}

	for r := 0; r < 256; r += 32 {
		for g := 0; g < 256; g += 32 {
			for b := 0; b < 256; b += 32 {
				p = append(p, color.RGBA{uint8(r), uint8(g), uint8(b), 255})
			}
		}
	}

	println("Palette size: ", len(p))

	m := [2][3]float32{
		{
			0, 0, 7.0 / 16.0,
		},
		{
			3.0 / 16.0, 5.0 / 16.0, 1.0 / 16.0,
		},
	}
	shift := -1

	result := image.NewRGBA(bounds)
	err := NewErrorImage(bounds)
	var diff uint64 = 0

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// using the closest color
			r, e, d := findColor(err.PixelErrorAt(x, y), img.At(x, y), p)
			result.SetRGBA(x, y, r)
			err.SetPixelError(x, y, e)
			diff += uint64(d)

			// diffusing the error using the diffusion matrix
			for i, v1 := range m {
				for j, v2 := range v1 {
					err.SetPixelError(x+j+shift, y+i,
						err.PixelErrorAt(x+j+shift, y+i).Add(err.PixelErrorAt(x, y).Mul(v2)))
				}
			}
		}
	}

	fmt.Printf("%f", (1.0-float64(diff)/float64(3*255*bounds.Max.X*bounds.Max.Y))*100.0)

	storeImage(output, result)
}
