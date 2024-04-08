package operation

import (
	"bytes"
	icc "imagecore/icc"
	image_parser "imagecore/image_parser"
)

func EmbedProfile(profile_name string) Operation {

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

		// Parse binary image to segments.
		parsed_image, err := image_parser.Parse(r)
		if err != nil {
			// Change the error state.
			currentImage.errorState = err
			// Return error.
			return currentImage, err
		}

		err = icc.EmbedIccProfile(profile_name, parsed_image)
		if err != nil {
			// Change the error state.
			currentImage.errorState = err
			// Return error.
			return currentImage, err
		}

		// Create a buffer to hold the image data.
		buf := new(bytes.Buffer)
		_, err = parsed_image.WriteTo(buf)
		if err != nil {
			// Change the error state.
			currentImage.errorState = err
			// Return error.
			return currentImage, err
		}

		return CurrentProcessingImage{ImageData: buf.Bytes(), isBinaryData: true}, nil
	}

}
