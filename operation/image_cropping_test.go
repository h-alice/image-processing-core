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

	// Compare images.
	if cropped_img.Bounds().Dx() != 14 || cropped_img.Bounds().Dy() != 14 {
		t.Fatalf("Cropped image bounds are not correct: %v", cropped_img.Bounds())
	} else {

		t.Logf("Image bounds: %v", cropped_img.Bounds())
	}
	for x := 0; x < cropped_img.Bounds().Dx(); x++ {
		for y := 0; y < cropped_img.Bounds().Dy(); y++ {
			// NOTE: This is a tricky part. If using At(x, y) to compare the pixel, it may not work as expected.
			//       So, we use RGBA64At(x, y) instead.
			if img.Image.Bounds().RGBA64At(x, y) != cropped_img.Bounds().RGBA64At(x, y) {
				t.Logf("Image 1: %v", img.Image.At(x, y))
				t.Logf("Image 2: %v", cropped_img.At(x, y))
				t.Fatalf("Images are not equal at (%d, %d)", x, y)
			}
		}
	}
}
