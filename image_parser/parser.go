package image_parser

import (
	"bytes"
	"errors"
	jpeg_parser "imagecore/image_parser/jpeg"
	"io"
)

var (
	ErrUnsupportedFileType = errors.New("file type not supported")
)

type ParserdImage interface {
	EmbedIccProfile(icc_profile []byte) error
	ReadFrom(io.Reader) (int64, error)
	WriteTo(io.Writer) (int64, error)
}

func Parse(rd io.Reader) (ParserdImage, error) {

	buf := bytes.NewBuffer([]byte{})
	_, err := io.Copy(buf, rd)
	if err != nil {
		return nil, err
	}

	switch {
	// JPEG
	case bytes.Equal(buf.Bytes()[0:2], []byte{'\xFF', '\xD8'}):
		parsed_jpeg := new(jpeg_parser.JpegImage)
		_, err := parsed_jpeg.ReadFrom(buf)
		if err != nil {
			return nil, err
		}
		return parsed_jpeg, nil
	// PNG
	case bytes.Equal(buf.Bytes()[0:8], []byte("\x89\x50\x4E\x47\x0D\x0A\x1A\x0A")):
		return nil, ErrUnsupportedFileType

	default:
		return nil, ErrUnsupportedFileType
	}
}
