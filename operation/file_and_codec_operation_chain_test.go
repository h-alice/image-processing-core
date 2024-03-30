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

	if im.LastError() != nil {
		t.Errorf("Expected no error, got: %v", im.LastError())
	}

	im, err = CreateImageFromFile(test_jpg_relative_path)
	if err != nil {
		t.Errorf("Error creating image from file: %v", err)
	}

	// Check basic properties.
	if !im.IsBinary() {
		t.Errorf("Expected image to be binary data")
	}

	if im.LastError() != nil {
		t.Errorf("Expected no error, got: %v", im.LastError())
	}

}

func TestFileDecoding(t *testing.T) {
	test_png_relative_path := "./test_resources/test_ayaya.png"
	test_jpg_relative_path := "./test_resources/test_ayaya.jpg"

	im, _ := CreateImageFromFile(test_png_relative_path) // This is tested in previous case.

	im_decoded := im.Then(Decode())

	// Check image properties.
	if im_decoded.LastError() != nil {
		t.Errorf("Expected no error, got: %v", im_decoded.LastError())
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
	if im_decoded.LastError() != nil {
		t.Errorf("Expected no error, got: %v", im_decoded.LastError())
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
	if im_decoded.LastError() == nil {
		t.Errorf("Expected error, got none")
	} else {
		t.Logf("Expected error, got: %v", im_decoded.LastError())
	}

	if !im_decoded.IsBinary() {
		t.Errorf("Expected image to be binary data")
	}
}

func TestEncoding(t *testing.T) {
	test_png_relative_path := "./test_resources/test_ayaya.png"
	test_jpg_relative_path := "./test_resources/test_ayaya.jpg"

	// Test 1: Read PNG and encode to JPG.
	im, _ := CreateImageFromFile(test_png_relative_path) // This is tested in previous case.

	im_decoded := im.Then(Decode())

	// Re-encode image to jpg.
	im_encoded := im_decoded.Then(Encode("jpg", nil))

	// Check image properties.
	if im_encoded.LastError() != nil {
		t.Errorf("Expected no error, got: %v", im_encoded.LastError())
	}

	if !im_encoded.IsBinary() {
		t.Errorf("Expected image to be binary data")
	}

	// Re-encoded image should be jpg.
	if im_encoded.Then(Decode()).ImageFormat() != "jpg" && im_encoded.Then(Decode()).ImageFormat() != "jpeg" {
		t.Errorf("Expected image format to be jpg, got: %v", im_encoded.ImageFormat())
	}

	// Test 2: Read JPG and encode to PNG.
	im, _ = CreateImageFromFile(test_jpg_relative_path) // This is tested in previous case.

	im_decoded = im.Then(Decode())

	// Re-encode image to png.
	im_encoded = im_decoded.Then(Encode("png", nil))

	// Check image properties.
	if im_encoded.LastError() != nil {
		t.Errorf("Expected no error, got: %v", im_encoded.LastError())
	}

	if !im_encoded.IsBinary() {
		t.Errorf("Expected image to be binary data")
	}

	// Re-encoded image should be png.
	if im_encoded.Then(Decode()).ImageFormat() != "png" {
		t.Errorf("Expected image format to be png, got: %v", im_encoded.ImageFormat())
	}

}

func TestInvalidEncoding(t *testing.T) {
	test_png_relative_path := "./test_resources/test_ayaya.png"

	// Test : Read PNG and encode to TXT, which is not supported.
	im, _ := CreateImageFromFile(test_png_relative_path) // This is tested in previous case.

	im_decoded := im.Then(Decode())

	// Re-encode image to txt, which is not supported.
	im_encoded := im_decoded.Then(Encode("txt", nil))

	// Check image properties.
	if im_encoded.LastError() == ErrEncodingFormatNotSupported {
		t.Logf("Expected error, got: %v", im_encoded.LastError())
	} else {
		t.Errorf("Expected error, got: %v", im_encoded.LastError())
	}

}
