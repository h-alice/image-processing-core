package png_parser

import (
	"bytes"
	"errors"
	"io"
)

var (
	ErrSignatureMismatch = errors.New("png signature mismatch")
)

var PNG_HEADER = []byte{'\x89', '\x50', '\x4E', '\x47', '\x0D', '\x0A', '\x1A', '\x0A'}

type PngImage struct {
	Segments []PngSegment
}

func ReadPng(r io.Reader) ([]PngSegment, int64, error) {

	total_read := int64(0)

	signature := make([]byte, 8)
	read, err := io.ReadFull(r, signature)
	total_read += int64(read)
	if err != nil {
		return nil, total_read, err
	}
	if !bytes.Equal(signature, PNG_HEADER) {
		return nil, total_read, ErrSignatureMismatch
	}

	seg_list := make([]PngSegment, 0)

	for {
		seg := new(PngGeneralSegment)
		read, err := seg.ReadFrom(r)
		total_read += read
		if err != nil {
			return nil, total_read, err
		}

		seg_list = append(seg_list, seg)

		if seg.SegmentType == "IEND" {
			break
		}
	}
	return seg_list, total_read, nil
}

func (img *PngImage) ReadFrom(r io.Reader) (int64, error) {
	seg_list, total_read, err := ReadPng(r)
	if err != nil {
		if err != io.EOF {
			return total_read, err
		}
	}

	img.Segments = seg_list
	return total_read, nil
}

func (img PngImage) WriteTo(wt io.Writer) (int64, error) {
	total_written := int64(0)
	for _, seg := range img.Segments {
		written, err := seg.WriteTo(wt)
		total_written += written
		if err != nil {
			return total_written, err
		}
	}
	return total_written, nil
}

func NewGeneralSegment(segment_type string, data []byte) *PngGeneralSegment {
	seg := PngGeneralSegment{
		Length:      len(data),
		SegmentType: segment_type,
		Data:        &data,
	}
	return &seg
}

func (im *PngImage) EmbedIccProfile(icc_profile []byte) {

}
