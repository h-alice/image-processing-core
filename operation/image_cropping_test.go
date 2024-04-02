package operation

import (
	"testing"
)

func TestImageCropping(t *testing.T) {

	// Load image.
	img, err := CreateImageFromFile("test_resources/test_ayaya.jpg")
	if err != nil {
		t.Fatalf("Failed to load image: %v", err)
	}

	img = img.Then(Decode())

	// Crop image.
	cropped_img, err := cropImageInternal(img.Image, 0, 14, 0, 14)
	if err != nil {
		t.Fatalf("Failed to crop image: %v", err)
	}

	// Set cropped image to current image.
	img.Image = cropped_img

	// Encode image.
	img = img.Then(Encode("jpeg", nil))

	// Write image to file.
	err = img.Then(WriteImageToFile("test_resources/test1_cropped.jpg")).LastError()
	if err != nil {
		t.Fatalf("Failed to write image: %v", err)
	}

	// Remove file.
	//err = os.Remove("test_resources/test1_cropped.jpg")
	//if err != nil {
	//	t.Fatalf("Failed to remove file: %v", err)
	//}
}
