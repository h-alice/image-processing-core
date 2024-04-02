package operation

import (
	"image"

	"golang.org/x/image/draw"
)

// Crop image by specifying the boundary.
func cropImageInternal(input_img image.Image, x1 int, x2 int, y1 int, y2 int) (image.Image, error) {

	crop := image.Rect(x1, y1, x2, y2)

	canvas := image.NewRGBA(crop)

	draw.Draw(canvas, crop, input_img, crop.Min, draw.Over)

	return canvas, nil
}
