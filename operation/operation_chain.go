package operation

import (
	"errors"
	"image"
)

// The CurrentProcessingImage is a struct that holds the current image data.
//
// `isBinaryDate` ia a flag to track if the image is binary data. It should be used internally and should not be modified outside the package.
// `errorState` is used to track error in the image processing chain, this can block the execution of the chain when an error is encountered, without craching the program.
type CurrentProcessingImage struct {
	// Image binary data, or go `image.Image` instance.
	ImageData    []byte      // The binary data.
	Image        image.Image // The `image.Image` instance.
	imageFormat  string      // The image format.
	isBinaryData bool        // Flag to track if the image is binary data.
	errorState   error       // Error state, this is used to track error in the image processing chain.
}

func (c CurrentProcessingImage) IsBinary() bool {
	return c.isBinaryData
}

func (c CurrentProcessingImage) LastError() error {
	return c.errorState
}

func (c CurrentProcessingImage) ImageFormat() string {
	return c.imageFormat
}

// Define errors.
var (
	ErrOperationNotSupportInBinary = errors.New("Operation not supported in binary format, convert to `image.Image` first")
	ErrOperationNotSupportInImage  = errors.New("Operation not supported in image.Image instance, convert to binary data first")
)

// Every processing action shoult have `CurrentProcessingImage` as input.
// And return `(CurrentProcessingImage, error)` as output.
type Operation func(currentImage CurrentProcessingImage) (CurrentProcessingImage, error)

// `Then“ method is used to chain operations.
//
// Currently every operation will return a new `CurrentProcessingImage` instance.
// NOTE: In-place operation will be a future enhancement.
func (currentImage CurrentProcessingImage) Then(operations Operation) CurrentProcessingImage {

	// Check error state.
	if currentImage.LastError() != nil {
		// Refuse to execute and return the original image.
		return currentImage
	}

	// Execute operation.
	newImage, err := operations(currentImage)
	if err != nil {
		// Return the original image, with error state.
		newImage.errorState = err
		return newImage
	}

	return newImage
}

// Conditioned `Then` method.
//
// This method will execute the operation only if the condition is met.
func (currentImage CurrentProcessingImage) ThenIf(condition bool, operations Operation) CurrentProcessingImage {
	if condition {
		return currentImage.Then(operations)
	}
	return currentImage
}
