package operation

import (
	"fmt"
	"image"
	"io"

	"image/jpeg"
	"image/png"
)

func init() {
	image.RegisterFormat("jpeg", "\xff\xd8", jpeg.Decode, jpeg.DecodeConfig)
	image.RegisterFormat("png", "\x89\x50\x4E\x47\x0D\x0A\x1A\x0A", png.Decode, png.DecodeConfig)
}

func Decode(r io.Reader) (image.Image, string, error) {
	image, format, err := image.Decode(r)
	return image, format, err
}

type EncoderOption struct {
	// For JPEG encoder.
	Quality *int
}

func Encode(w io.Writer, input_image *image.Image, format string, opt *EncoderOption) (err error) {
	if input_image == nil {
		err = fmt.Errorf("np input image")
		return
	}
	if opt == nil {
		opt = new(EncoderOption)
	}

	switch format {
	case "jpg", "jpeg":
		quality := 100
		if opt.Quality == nil {
			quality = 100
		} else {
			quality = *opt.Quality
		}

		err = jpeg.Encode(w, *input_image, &jpeg.Options{Quality: quality})
		return
	case "png":
		// Quality option is ignored.
		err = png.Encode(w, *input_image)
		return
	default:
		err = fmt.Errorf("unsupported format")
		return
	}

}
