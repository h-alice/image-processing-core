package operation

import (
	"image"
	"log"
	"strings"

	"golang.org/x/image/draw"
)

// Get resize algorithm.
func resizeAlgotithm(algo string) (algorithm draw.Interpolator) {

	switch strings.ToLower(algo) {
	case "nearestneighbor":
		algorithm = draw.NearestNeighbor
	case "catmullrom":
		algorithm = draw.CatmullRom
	case "approxbiLinear":
		algorithm = draw.ApproxBiLinear
	default:
		log.Println("[!] Using default resize algorithm 'Catmull-Rom'.")
		algorithm = draw.CatmullRom
	}
	return
}

// Get resize algotithm.
//
// Deprecated: Internal use only.
func createResizeBoundryByFactor(input image.Rectangle, factor float32) image.Rectangle {

	resized_boundary := image.Rect(
		0,
		0,
		int(float32(input.Max.X)/factor),
		int(float32(input.Max.Y)/factor))

	return resized_boundary
}

// Resize image.
//
// Deprecated: To keep image aspect ratio, it's not recommended to use the function directly.
func resizeImageInternal(in image.Image, algo string, boundary image.Rectangle) image.Image {

	canvas := image.NewRGBA(boundary)

	algorithm := resizeAlgotithm(algo)
	algorithm.Scale(canvas, canvas.Rect, in, in.Bounds(), draw.Over, nil)

	return canvas
}

// Resize image by specifying resized width.
func ResizeImageByWidth(in image.Image, algo string, x int) image.Image {

	factor := float32(in.Bounds().Max.X) / float32(x)
	boundary := createResizeBoundryByFactor(in.Bounds(), factor)
	return resizeImageInternal(in, algo, boundary)
}

// Resize image by specifying resized height.
func ResizeImageByHeight(in image.Image, algo string, y int) image.Image {

	factor := float32(in.Bounds().Max.Y) / float32(y)
	boundary := createResizeBoundryByFactor(in.Bounds(), factor)
	return resizeImageInternal(in, algo, boundary)
}

// Resize image by specifying resize factor.
func ResizeImageByFactor(in image.Image, algo string, factor float32) image.Image {

	boundary := createResizeBoundryByFactor(in.Bounds(), factor)
	return resizeImageInternal(in, algo, boundary)
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
