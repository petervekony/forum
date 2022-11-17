package server

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"

	"golang.org/x/image/draw"
)

func ProcessImage(file string) error {
	img, err := os.Open(file)
	if err != nil {
		return err
	}
	defer img.Close()
	decoded, err := jpeg.Decode(img)
	if err != nil {
		return err
	}
	fmt.Println(decoded.Bounds().Max)
	resized := image.NewRGBA(image.Rect(0, 0, decoded.Bounds().Max.X/2, decoded.Bounds().Max.Y/2))
	draw.NearestNeighbor.Scale(resized, resized.Rect, decoded, decoded.Bounds(), draw.Over, nil)
	output, err := os.Create("image_resized.jpg")
	if err != nil {
		return err
	}
	defer output.Close()
	jpeg.Encode(output, resized, nil)
	return nil
}

// Given image is 5120x2880
// Want to have it 2560x1440
