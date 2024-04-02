package operation

import (
	"image"

	"golang.org/x/image/draw"
)

//type cropAlignment func(image.Rectangle, int, int) image.Rectangle

// Crop image by specifying the boundary.
func cropImageInternal(input_img image.Image, crop_boundary image.Rectangle) (image.Image, error) {

	// Reset the boundary origin to (0, 0).
	canvas_boundary := crop_boundary.Sub(crop_boundary.Min)

	// Create a new canvas with the specified boundary.
	canvas := image.NewRGBA(canvas_boundary)

	// Draw the input image onto the canvas, with the specified boundary.
	draw.Draw(canvas, canvas_boundary, input_img, crop_boundary.Min, draw.Src)

	return canvas, nil
}
