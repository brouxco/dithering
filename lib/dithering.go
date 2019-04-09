package dithering

import (
	"image"
	"image/png"
	"log"
	"os"

	_ "image/png"
)

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
	defer file.Close();

	if err = png.Encode(file, img); err != nil {
		log.Fatal(err)
	}
}

func Dither(input string, output string){
	img := loadImage(input)

	storeImage(output, img)
}