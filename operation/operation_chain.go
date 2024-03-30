package operation

import (
	"errors"
	"image"
)

type CurrentProcessingImage struct {
	// Image binary data, or go `image.Image` instance.
	ImageData []byte
	Image     image.Image
	IsBinary  bool
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
