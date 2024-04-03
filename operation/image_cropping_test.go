package operation

import (
	"image"
	"testing"
)

func TestImageCroppingInternal(t *testing.T) {

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

func TestImageCroppingTopLeft(t *testing.T) {

	// Load image.
	img, err := CreateImageFromFile("test_resources/test_ayaya.png")
	if err != nil {
		t.Fatalf("Failed to load image: %v", err)
	}

	img = img.Then(Decode())

	cropped_img := img.Then(Crop(12, 13, "topleft"))

	if cropped_img.LastError() != nil {
		t.Fatalf("Failed to crop image: %v", cropped_img.LastError())
	}

	// Calculate the cropping area.
	crop_area := TopLeftAlignment(img.Image.Bounds(), 12, 13)

	// Check if cropped image bounds are correct.
	if crop_area.Min.X != 0 || crop_area.Min.Y != 0 || crop_area.Max.X != 12 || crop_area.Max.Y != 13 {
		t.Fatalf("Cropping area is not correct: %v", crop_area)
	} else {
		t.Logf("Cropping area: %v", crop_area)
	}

	// Compare images, pixel-by-pixel.
	for x := 0; x < cropped_img.Image.Bounds().Dx(); x++ {
		for y := 0; y < cropped_img.Image.Bounds().Dy(); y++ {

			// Calculate the starting point of cropping area.
			crop_start := image.Point{X: crop_area.Min.X, Y: crop_area.Min.Y}

			// Calculate the current evaluating point on the original image by shifting the starting point of cropping area.
			current_point_original_image := crop_start.Add(image.Point{X: x, Y: y})

			// NOTE: This is a tricky part. If using At(x, y) to compare the pixel, it may not work as expected.
			//       So, we use RGBA64At(x, y) instead.
			if img.Image.Bounds().RGBA64At(current_point_original_image.X, current_point_original_image.Y) != cropped_img.Image.Bounds().RGBA64At(x, y) {
				t.Logf("Image 1: %v", img.Image.At(x, y))
				t.Logf("Image 2: %v", cropped_img.Image.At(x, y))
				t.Fatalf("Images are not equal at (%d, %d)", x, y)
			}
		}
	}

}

func TestImageCroppingCenter(t *testing.T) {

	// Load image.
	img, err := CreateImageFromFile("test_resources/test_ayaya.png")
	if err != nil {
		t.Fatalf("Failed to load image: %v", err)
	}

	img = img.Then(Decode())

	cropped_img := img.Then(Crop(11, 15, "center"))

	if cropped_img.LastError() != nil {
		t.Fatalf("Failed to crop image: %v", cropped_img.LastError())
	}

	// Calculate the cropping area.
	crop_area := CenterAlignment(img.Image.Bounds(), 11, 15)

	// Check if cropped image bounds are correct.
	// 14 - int(11 / 2) = 9
	// 14 - int(15 / 2) = 7
	if crop_area.Min.X != 9 || crop_area.Min.Y != 7 || crop_area.Max.X != 20 || crop_area.Max.Y != 22 {
		t.Fatalf("Cropping area is not correct: %v", crop_area)
	} else {
		t.Logf("Cropping area: %v", crop_area)
	}

	// Compare images, pixel-by-pixel.
	for x := 0; x < cropped_img.Image.Bounds().Dx(); x++ {
		for y := 0; y < cropped_img.Image.Bounds().Dy(); y++ {

			// Calculate the starting point of cropping area.
			crop_start := image.Point{X: crop_area.Min.X, Y: crop_area.Min.Y}

			// Calculate the current evaluating point on the original image by shifting the starting point of cropping area.
			current_point_original_image := crop_start.Add(image.Point{X: x, Y: y})

			// NOTE: This is a tricky part. If using At(x, y) to compare the pixel, it may not work as expected.
			//       So, we use RGBA64At(x, y) instead.
			if img.Image.Bounds().RGBA64At(current_point_original_image.X, current_point_original_image.Y) != cropped_img.Image.Bounds().RGBA64At(x, y) {
				t.Logf("Image 1: %v", img.Image.At(x, y))
				t.Logf("Image 2: %v", cropped_img.Image.At(x, y))
				t.Fatalf("Images are not equal at (%d, %d)", x, y)
			}
		}
	}

}

func TestInvalidBoundary(t *testing.T) {

	// Load image.
	img, err := CreateImageFromFile("test_resources/test_ayaya.png")
	if err != nil {
		t.Fatalf("Failed to load image: %v", err)
	}

	img = img.Then(Decode())

	// Test invalid boundary.
	crop_area := image.Rect(0, 0, 100, 100)
	_, err = cropImageInternal(img.Image, crop_area)
	if err == nil {
		t.Fatalf("Expected error, got nil")
	} else if err == ErrCroppingAreaOutOfBound {
		t.Logf("Got error as expected: %v", err)
	} else {
		t.Fatalf("Got unexpected error: %v", err)
	}

}
