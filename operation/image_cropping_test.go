package operation

import (
	"image"
	"testing"
)

func TestImageCropping(t *testing.T) {

	// Load image.
	img, err := CreateImageFromFile("test_resources/test_ayaya.png")
	if err != nil {
		t.Fatalf("Failed to load image: %v", err)
	}

	img = img.Then(Decode())

	// Crop image.
	crop_area := image.Rect(7, 13, 20, 24)
	t.Logf("Cropping area: %v", crop_area)

	cropped_img, err := cropImageInternal(img.Image, crop_area)
	if err != nil {
		t.Fatalf("Failed to crop image: %v", err)
	}

	// Check if cropped image bounds are correct.
	if cropped_img.Bounds().Dx() != crop_area.Dx() || cropped_img.Bounds().Dy() != crop_area.Dy() {
		t.Fatalf("Cropped image bounds are not correct: %v", cropped_img.Bounds())
	} else {
		t.Logf("Image bounds: %v", cropped_img.Bounds())
	}

	// Compare images, pixel-by-pixel.
	for x := 0; x < cropped_img.Bounds().Dx(); x++ {
		for y := 0; y < cropped_img.Bounds().Dy(); y++ {

			// Calculate the starting point of cropping area.
			crop_start := image.Point{X: crop_area.Min.X, Y: crop_area.Min.Y}

			// Calculate the current evaluating point on the original image by shifting the starting point of cropping area.
			current_point_original_image := crop_start.Add(image.Point{X: x, Y: y})

			// NOTE: This is a tricky part. If using At(x, y) to compare the pixel, it may not work as expected.
			//       So, we use RGBA64At(x, y) instead.
			if img.Image.Bounds().RGBA64At(current_point_original_image.X, current_point_original_image.Y) != cropped_img.Bounds().RGBA64At(x, y) {
				t.Logf("Image 1: %v", img.Image.At(x, y))
				t.Logf("Image 2: %v", cropped_img.At(x, y))
				t.Fatalf("Images are not equal at (%d, %d)", x, y)
			}
		}
	}
}
