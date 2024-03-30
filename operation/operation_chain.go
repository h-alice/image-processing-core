package operation

import (
	"errors"
	"image"
)

// The CurrentProcessingImage is a struct that holds the current image data.
//
// `isBinaryDate` ia a flag to track if the image is binary data. It should be used internally and should not be modified outside the package.
type CurrentProcessingImage struct {
	// Image binary data, or go `image.Image` instance.
	ImageData    []byte      // The binary data.
	Image        image.Image // The `image.Image` instance.
	isBinaryData bool        // Flag to track if the image is binary data.
}

func (c CurrentProcessingImage) IsBinary() bool {
	return c.isBinaryData
}

// Define errors.
var (
	ErrOperationNotSupportInBinary = errors.New("Operation not supported in binary format, convert to `image.Image` first")
	ErrOperationNotSupportInImage  = errors.New("Operation not supported in image.Image instance, convert to binary data first")
)

// Every processing action shoult have `CurrentProcessingImage` as input.
// And return `(CurrentProcessingImage, error)` as output.
type Operation func(currentImage CurrentProcessingImage) (CurrentProcessingImage, error)

//func (currentImage *CurrentProcessingImage) Then()
