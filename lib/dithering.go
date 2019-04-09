package dithering

import (
	"image"
	"image/png"
	"log"
	"os"

	_ "image/png"
)

func Dither(input string, output string){
	reader, err := os.Open(input)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	img, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create(output)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close();

	if err = png.Encode(file, img); err != nil {
		log.Fatal(err)
	}
}