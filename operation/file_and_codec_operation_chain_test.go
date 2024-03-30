package operation

import (
	"fmt"
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

	im, err := CreateImageFromFile(test_png_relative_path)
	if err != nil {
		t.Errorf("Error creating image from file: %v", err)
	}

	fmt.Printf("%#v\n", im)
}
