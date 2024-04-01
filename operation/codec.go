package operation

import (
	"bytes"
	"errors"
	"image"
	"strings"

	"image/jpeg"
	"image/png"
)

// Image binary header.
var (
	JPEG_HEADER = []byte("\xff\xd8")
	PNG_HEADER  = []byte("\x89\x50\x4E\x47\x0D\x0A\x1A\x0A")
)

// Register common image format.
func init() {
	image.RegisterFormat("jpeg", string(JPEG_HEADER), jpeg.Decode, jpeg.DecodeConfig)
	image.RegisterFormat("png", string(PNG_HEADER), png.Decode, png.DecodeConfig)
}

// Define some errors.
var (
	ErrEncodingFormatNotSupported = errors.New("encoding format not supported")
)

type EncoderOption struct {
	// For JPEG encoder.
	Quality int
}

// Decode image from given `CurrentProcessingImage` instance.
func Decode() Operation {
	return func(currentImage CurrentProcessingImage) (CurrentProcessingImage, error) {
		// Input image should in binary format.
		if !currentImage.IsBinary() {
			// Change the error state.
			currentImage.errorState = ErrOperationNotSupportInImage
			// Return error.
			return currentImage, ErrOperationNotSupportInImage
		}

		// Create reader from binary data.
		r := bytes.NewReader(currentImage.ImageData)

		image, format, err := image.Decode(r)
		if err != nil {
			// Change the error state.
			currentImage.errorState = err
			// Return error.
			return currentImage, err
		}

		return CurrentProcessingImage{Image: image, isBinaryData: false, imageFormat: format}, nil
	}
}

// Encode image to given format.
func Encode(format string, opt *EncoderOption) Operation {

	return func(currentImage CurrentProcessingImage) (CurrentProcessingImage, error) {
		// Input should not be binary data.
		if currentImage.IsBinary() {
			// Change the error state.
			currentImage.errorState = ErrOperationNotSupportInBinary
			// Return error.
			return currentImage, ErrOperationNotSupportInBinary
		}

		if opt == nil {
			opt = new(EncoderOption)
		}

		// Create binary buffer for output.
		buf := new(bytes.Buffer)

		// Encodes the image to desired format.
		switch strings.ToLower(format) { // Convert to lower case.
		case "jpg", "jpeg":
			quality := 100
			if opt.Quality == 0 {
				quality = 100
			} else {
				quality = opt.Quality
			}

			err := jpeg.Encode(buf, currentImage.Image, &jpeg.Options{Quality: quality})
			if err != nil {
				// Change the error state.
				currentImage.errorState = err
				// Return error.
				return currentImage, err
			}
		case "png":
			// Quality option is ignored.
			err := png.Encode(buf, currentImage.Image)
			if err != nil {
				// Change the error state.
				currentImage.errorState = err
				// Return error.
				return currentImage, err
			}
		default:
			err := ErrEncodingFormatNotSupported
			if err != nil {
				// Change the error state.
				currentImage.errorState = err
				// Return error.
				return currentImage, err
			}
		}

		// Get the bytes from buffer.
		binary_content := make([]byte, buf.Len())
		copy(binary_content, buf.Bytes())

		// Return the new image.
		return CurrentProcessingImage{ImageData: binary_content, isBinaryData: true, imageFormat: format}, nil
	}

}
