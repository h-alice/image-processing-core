package operation

import (
	"os"
)

// Create image from raw binary data.
func CreateImageFromBinary(data []byte) CurrentProcessingImage {
	// Simply return the binary data.
	return CurrentProcessingImage{ImageData: data, isBinaryData: true}
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
