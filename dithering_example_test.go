package dithering

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
)

func ExampleDither_Draw() {
	reader, err := os.Open("lenna.png")
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	src, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	dst := image.NewPaletted(image.Rect(0, 0, 640, 480), color.Palette{color.Black, color.White})

	floydSteinberg := Dither{[][]float32{
		{
			0, 0, 7.0 / 16.0,
		},
		{
			3.0 / 16.0, 5.0 / 16.0, 1.0 / 16.0,
		},
	},
	-1,
	}
	floydSteinberg.Draw(dst, dst.Bounds(), src, image.ZP)

	file, err := os.Create("result.png")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if err = png.Encode(file, dst); err != nil {
		log.Fatal(err)
	}
}
