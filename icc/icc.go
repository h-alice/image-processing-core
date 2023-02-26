package icc

import (
	"bytes"
	"pp-tools/parser/jpeg"
)

func JpegIccChunk(raw_profile []byte) (*jpeg.JpegGeneralSegment, error) {

	var err error = nil

	var buf bytes.Buffer

	_, err = buf.Write([]byte("ICC_PROFILE")) // ICC chunk header.
	if err != nil {
		return nil, err
	}

	err = buf.WriteByte('\x00') // Null byte
	if err != nil {
		return nil, err
	}

	err = buf.WriteByte('\x01') // NOTE: Profile chunk number.
	if err != nil {
		return nil, err
	}

	err = buf.WriteByte('\x01') // NOTE: Profile chunk number.
	if err != nil {
		return nil, err
	}

	_, err = buf.Write(raw_profile) // Profile content.
	if err != nil {
		return nil, err
	}

	segment, err := jpeg.NewJpegSegment("APP2", buf.Bytes())
	if err != nil {
		return nil, err
	}

	return segment, nil
}
