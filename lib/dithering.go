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

func abs(x int16) int16 {
	if x < 0 {
		return -x
	}
	return x
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

func findColor(err color.Color, pix color.Color) (color.RGBA, PixelError) {
	var errR, errG, errB, pixR, pixG, pixB int16
	_errR, _errG, _errB, _ := err.RGBA()
	_pixR, _pixG, _pixB, _ := pix.RGBA()

	errR = int16(_errR)
	errG = int16(_errG)
	errB = int16(_errB)
	pixR = int16(uint8(_pixR))
	pixG = int16(uint8(_pixG))
	pixB = int16(uint8(_pixB))

	diff_black := abs(errR-pixR) + abs(errG-pixG) + abs(errB-pixB)
	diff_white := abs((255+errR)-pixR) + abs((255+errG)-pixG) + abs((255+errB)-pixB)

	if diff_black < diff_white {
		return color.RGBA{0, 0, 0, 255},
			PixelError{float32(bound((-pixR) + errR)),
				float32(bound((-pixG) + errG)),
				float32(bound((-pixB) + errB)),
				1 << 16 - 1}
	} else {
		return color.RGBA{255, 255, 255, 255},
			PixelError{float32(bound((255 - pixR) + errR)),
				float32(bound((255 - pixG) + errG)),
				float32(bound((255 - pixB) + errB)),
				1 << 16 - 1}
	}
}

func Dither(input string, output string) {
	img := loadImage(input)

	bounds := img.Bounds()

	result := image.NewRGBA(bounds)
	err := NewErrorImage(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// using the closest color
			// TODO: palettes
			r, e := findColor(err.PixelErrorAt(x, y), img.At(x, y))
			result.SetRGBA(x, y, r)
			err.SetPixelError(x, y, e)

			// diffusing the error using the diffusion matrix
			// TODO: matrices
			if x+1 < bounds.Max.X {
				err.SetPixelError(x+1, y, err.PixelErrorAt(x, y))
			}
		}
	}

	storeImage(output, err)
}
