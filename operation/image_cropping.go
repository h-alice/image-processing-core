package operation

import (
	"image"

	"golang.org/x/image/draw"
)

//type cropAlignment func(image.Rectangle, int, int) image.Rectangle

// Crop image by specifying the boundary.
func cropImageInternal(input_img image.Image, crop_boundary image.Rectangle) (image.Image, error) {

	canvas := image.NewRGBA(crop_boundary.Sub(crop_boundary.Min))

	draw.Draw(canvas, crop_boundary, input_img, crop_boundary.Min, draw.Src)

	return canvas, nil
}
