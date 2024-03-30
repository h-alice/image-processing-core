package operation

import (
	"os"
	"testing"
)

func TestPreTestTesting(t *testing.T) {
	test_png_relative_path := "./test_resources/test_ayaya.png"
	_, err := os.ReadFile(test_png_relative_path)
	if err != nil {
		t.Errorf("(Pre-testing) Error reading file: %v", err)
	}

	test_jpg_relative_path := "./test_resources/test_ayaya.jpg"
	_, err = os.ReadFile(test_jpg_relative_path)
	if err != nil {
		t.Errorf("(Pre-testing) Error reading file: %v", err)
	}
}

func TestCreateImageFromFile(t *testing.T) {
	test_png_relative_path := "./test_resources/test_ayaya.png"
	test_jpg_relative_path := "./test_resources/test_ayaya.jpg"

	im, err := CreateImageFromFile(test_png_relative_path)
	if err != nil {
		t.Errorf("Error creating image from file: %v", err)
	}

	// Check basic properties.
	if !im.IsBinary() {
		t.Errorf("Expected image to be binary data")
	}

	if im.lastError() != nil {
		t.Errorf("Expected no error, got: %v", im.lastError())
	}

	im, err = CreateImageFromFile(test_jpg_relative_path)
	if err != nil {
		t.Errorf("Error creating image from file: %v", err)
	}

	// Check basic properties.
	if !im.IsBinary() {
		t.Errorf("Expected image to be binary data")
	}

	if im.lastError() != nil {
		t.Errorf("Expected no error, got: %v", im.lastError())
	}

}

func TestFileDecoding(t *testing.T) {
	test_png_relative_path := "./test_resources/test_ayaya.png"
	test_jpg_relative_path := "./test_resources/test_ayaya.jpg"

	im, _ := CreateImageFromFile(test_png_relative_path) // This is tested in previous case.

	im_decoded := im.Then(Decode())

	// Check image properties.
	if im_decoded.lastError() != nil {
		t.Errorf("Expected no error, got: %v", im_decoded.lastError())
	}

	if im_decoded.IsBinary() {
		t.Errorf("Expected image to be image.Image instance")
	}

	if im_decoded.ImageFormat() != "png" {
		t.Errorf("Expected image format to be png, got: %v", im_decoded.ImageFormat())
	}

	// Jpg test.

	im, _ = CreateImageFromFile(test_jpg_relative_path) // This is tested in previous case.

	im_decoded = im.Then(Decode())

	// Check image properties.
	if im_decoded.lastError() != nil {
		t.Errorf("Expected no error, got: %v", im_decoded.lastError())
	}

	if im_decoded.IsBinary() {
		t.Errorf("Expected image to be image.Image instance")
	}

	if im_decoded.ImageFormat() != "jpeg" && im_decoded.ImageFormat() != "jpg" {
		t.Errorf("Expected image format to be jpeg, got: %v", im_decoded.ImageFormat())
	}
}

func TestInvaildDecoding(t *testing.T) {
	test_txt_relative_path := "./test_resources/.gitignore"

	im, err := CreateImageFromFile(test_txt_relative_path)
	if err != nil {
		t.Errorf("Error creating image from file: %v", err)
	}

	im_decoded := im.Then(Decode())

	// Check image properties.
	if im_decoded.lastError() == nil {
		t.Errorf("Expected error, got none")
	} else {
		t.Logf("Expected error, got: %v", im_decoded.lastError())
	}

	if !im_decoded.IsBinary() {
		t.Errorf("Expected image to be binary data")
	}
}
