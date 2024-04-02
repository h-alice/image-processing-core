package operation

import (
	"errors"
	"image"

	"golang.org/x/image/draw"
)

// Define errors.
var (
	ErrInvalidCropBoundary    = errors.New("invalid crop boundary")
	ErrCroppingAreaOutOfBound = errors.New("crop boundary is outside of the original image")
)

// Defines a prototype for crop alignment function. (e.g. Center, TopLeft, BottomRight)
// The first parameter is the boundary of the original image.
// The second and third parameters are the width and height of the cropping area, respectively.
type CropAlignment func(image.Rectangle, int, int) image.Rectangle

func (CropAlignment) CenterAlignment(original_image_boundary image.Rectangle, cropped_width, cropped_height int) image.Rectangle {

	// Calculate the center point of the original image.
	center := image.Point{
		X: original_image_boundary.Dx() / 2,
		Y: original_image_boundary.Dy() / 2,
	}

	// Calculate the starting point of the cropping area.
	crop_start := image.Point{
		X: center.X - (cropped_width / 2),
		Y: center.Y - (cropped_height / 2),
	}

	// Calculate the ending point of the cropping area.
	crop_end := image.Point{
		X: crop_start.X + cropped_width,
		Y: crop_start.Y + cropped_height,
	}

	return image.Rect(crop_start.X, crop_start.Y, crop_end.X, crop_end.Y)
}

func (CropAlignment) TopLeftAlignment(original_image_boundary image.Rectangle, cropped_width, cropped_height int) image.Rectangle {

	// This is the most common-used alignment.

	// Calculate the starting point of the cropping area.
	crop_start := image.Point{
		X: 0,
		Y: 0,
	}

	// Calculate the ending point of the cropping area.
	crop_end := image.Point{
		X: crop_start.X + cropped_width,
		Y: crop_start.Y + cropped_height,
	}

	return image.Rect(crop_start.X, crop_start.Y, crop_end.X, crop_end.Y)
}

// Crop image by specifying the boundary.
func cropImageInternal(input_img image.Image, crop_boundary image.Rectangle) (image.Image, error) {

	// Check if the cropping area is inside the original image.
	if !crop_boundary.In(input_img.Bounds()) {
		return nil, ErrCroppingAreaOutOfBound
	}

	// Reset the boundary origin to (0, 0).
	canvas_boundary := crop_boundary.Sub(crop_boundary.Min)

	// Create a new canvas with the specified boundary.
	canvas := image.NewRGBA(canvas_boundary)

	// Draw the input image onto the canvas, with the specified boundary.
	draw.Draw(canvas, canvas_boundary, input_img, crop_boundary.Min, draw.Src)

	return canvas, nil
}

func Crop(crop_width int, crop_height int, alignment CropAlignment) Operation {

	return func(currentImage CurrentProcessingImage) (CurrentProcessingImage, error) {

		// Check if the input image is instance of `image.Image`.
		if currentImage.IsBinary() {
			// Change the error state.
			currentImage.errorState = ErrOperationNotSupportInBinary
			// Return error.
			return currentImage, ErrOperationNotSupportInBinary
		}

		// Get the boundary of the original image.
		original_image_boundary := currentImage.Image.Bounds()

		// Calculate the cropping area boundary.
		crop_boundary := alignment(original_image_boundary, crop_width, crop_height)

		cropped_image, err := cropImageInternal(currentImage.Image, crop_boundary)
		if err != nil {
			// Change the error state.
			currentImage.errorState = err
			// Return error.
			return currentImage, err
		}

		return CurrentProcessingImage{Image: cropped_image, isBinaryData: false}, nil
	}
}
