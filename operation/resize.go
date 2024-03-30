package operation

import (
	"image"
	"log"
	"strings"

	"golang.org/x/image/draw"
)

// Get resize algorithm by name.
//
// `Catmull-Rom` is the default algorithm, which provides the best quality.
func resizeAlgotithm(algo string) (algorithm draw.Interpolator) {

	switch strings.ToLower(algo) { // Convert to lower case in case of case-insensitive.
	case "nearestneighbor": // Fastest algorithm, but worst quality.
		algorithm = draw.NearestNeighbor
	case "catmullrom": // Default algorithm, best quality.
		algorithm = draw.CatmullRom
	case "approxbiLinear": // Somehow decent quality.
		algorithm = draw.ApproxBiLinear
	default:
		log.Println("[!] Using default resize algorithm 'Catmull-Rom'.")
		algorithm = draw.CatmullRom
	}
	return
}

// Using specified resizing-factor to create resized boundary.
//
// The output image will be `factor` times smaller than the input image.
func createResizeBoundryByFactor(input image.Rectangle, factor float32) image.Rectangle {

	resized_boundary := image.Rect(
		0,                                // X0
		0,                                // Y0
		int(float32(input.Max.X)/factor), // X1
		int(float32(input.Max.Y)/factor)) // Y1

	return resized_boundary
}

// Resize image.
//
// This creates a new image with the specified boundary, and then draw the input image onto it.
// NOTE: This is an internal function, and should not be used directly.
func resizeImageInternal(in image.Image, algo string, boundary image.Rectangle) image.Image {

	canvas := image.NewRGBA(boundary)

	algorithm := resizeAlgotithm(algo)
	algorithm.Scale(canvas, canvas.Rect, in, in.Bounds(), draw.Over, nil)

	return canvas
}

// Resize image by specifying resized width.
//
// The height is automatically calculated based on the aspect ratio.
func (in CurrentProcessingImage) ResizeImageByWidth(algo string, x int) Operation {

	return func(currentImage CurrentProcessingImage) (CurrentProcessingImage, error) {

		// Input should not be binary data.
		if in.IsBinary() {
			// Change the error state.
			currentImage.errorState = ErrOperationNotSupportInBinary
			// Return error.
			return currentImage, ErrOperationNotSupportInBinary
		}

		// Do resize on `image.Image` instance.
		factor := float32(in.Image.Bounds().Max.X) / float32(x)
		boundary := createResizeBoundryByFactor(in.Image.Bounds(), factor)
		resizedImage := resizeImageInternal(in.Image, algo, boundary)
		return CurrentProcessingImage{Image: resizedImage, isBinaryData: false}, nil
	}
}

// Resize image by specifying resized height.
//
// The width is automatically calculated based on the aspect ratio.
func (in CurrentProcessingImage) ResizeImageByHeight(algo string, y int) Operation {

	return func(currentImage CurrentProcessingImage) (CurrentProcessingImage, error) {

		// Input should not be binary data.
		if in.IsBinary() {
			// Change the error state.
			currentImage.errorState = ErrOperationNotSupportInBinary
			// Return error.
			return currentImage, ErrOperationNotSupportInBinary
		}

		// Do resize on `image.Image` instance.
		factor := float32(in.Image.Bounds().Max.Y) / float32(y)
		boundary := createResizeBoundryByFactor(in.Image.Bounds(), factor)
		resizedImage := resizeImageInternal(in.Image, algo, boundary)
		return CurrentProcessingImage{Image: resizedImage, isBinaryData: false}, nil
	}
}

// Resize image by specifying resize factor.
//
// The output image will be `factor` times smaller than the input image.
func (in CurrentProcessingImage) ResizeImageByFactor(algo string, factor float32) Operation {

	return func(currentImage CurrentProcessingImage) (CurrentProcessingImage, error) {

		// Input should not be binary data.
		if in.IsBinary() {
			// Change the error state.
			currentImage.errorState = ErrOperationNotSupportInBinary
			// Return error.
			return currentImage, ErrOperationNotSupportInBinary
		}

		// Do resize on `image.Image` instance.
		boundary := createResizeBoundryByFactor(in.Image.Bounds(), factor)
		resizedImage := resizeImageInternal(in.Image, algo, boundary)
		return CurrentProcessingImage{Image: resizedImage, isBinaryData: false}, nil
	}
}

/*
func main() {
	var err error

	path := `C:\Users\h-alice\Desktop\pp-tools\resource\test2.jpg`
	path2 := `C:\Users\h-alice\Desktop\test2_.jpg`

	input, err := os.ReadFile(path)
	if err != nil {
		log.Fatalln(err)
	}

	reader := bytes.NewReader(input)

	im, err := jpeg.Decode(reader)
	if err != nil {
		log.Fatalln(err)
	}

	jpeg_opt := jpeg.Options{Quality: 100}

	var buf bytes.Buffer

	jpeg.Encode(&buf, canvas, &jpeg_opt)

	output, _ := os.Create(path2)
	defer output.Close()
	parser.jpeg
	io.Copy(output, &buf)

}
*/
