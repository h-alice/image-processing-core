package operation

import (
	"errors"
	"io"
	"os"
)

// Define errors.
var (
	ErrWrittenLengthMismatch = errors.New("written length mismatch")
)

// Create image from raw binary data.
func CreateImageFromBinary(data []byte) CurrentProcessingImage {
	// Simply return the binary data.
	return CurrentProcessingImage{ImageData: data, isBinaryData: true}
}

// Create image from binary reader.
func CreateImageFromReader(reader io.Reader) (CurrentProcessingImage, error) {
	// Read data.
	data, err := io.ReadAll(reader)
	if err != nil {
		return CurrentProcessingImage{}, err
	}

	// Return image.
	return CreateImageFromBinary(data), nil
}

// Create image from file.
func CreateImageFromFile(path string) (CurrentProcessingImage, error) {
	// Read file.
	input, err := os.ReadFile(path)
	if err != nil {
		return CurrentProcessingImage{}, err
	}

	// Return image.
	return CreateImageFromBinary(input), nil
}

// Write image data to file.
func WriteImageToFile(path string) Operation {
	return func(currentImage CurrentProcessingImage) (CurrentProcessingImage, error) {
		// Check image format, it should be binary data.
		if !currentImage.IsBinary() {
			// Change the error state.
			currentImage.errorState = ErrOperationNotSupportInImage
			// Return error.
			return currentImage, ErrOperationNotSupportInImage
		}

		// Write to file.
		ofp, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
		defer func() {
			ofp.Close()
		}()

		if err != nil {
			// Change the error state.
			currentImage.errorState = err
			// Return error.
			return currentImage, err
		}

		// Write data.
		written, err := ofp.Write(currentImage.ImageData)
		if err != nil {
			// Change the error state.
			currentImage.errorState = err
			// Return error.
			return currentImage, err
		}

		// Check written length.
		if written != len(currentImage.ImageData) {
			// Change the error state.
			currentImage.errorState = ErrWrittenLengthMismatch
			// Return error.
			return currentImage, ErrWrittenLengthMismatch
		}

		// Return the image.
		return currentImage, nil
	}
}

// Write image data to writer.
func WriteImageToWriter(writer io.Writer) Operation {
	return func(currentImage CurrentProcessingImage) (CurrentProcessingImage, error) {
		// Check image format, it should be binary data.
		if !currentImage.IsBinary() {
			// Change the error state.
			currentImage.errorState = ErrOperationNotSupportInImage
			// Return error.
			return currentImage, ErrOperationNotSupportInImage
		}

		// Write to writer.
		written, err := writer.Write(currentImage.ImageData)
		if err != nil {
			// Change the error state.
			currentImage.errorState = err
			// Return error.
			return currentImage, err
		}

		// Check written length.
		if written != len(currentImage.ImageData) {
			// Change the error state.
			currentImage.errorState = ErrWrittenLengthMismatch
			// Return error.
			return currentImage, ErrWrittenLengthMismatch
		}

		// Return the image.
		return currentImage, nil
	}
}
