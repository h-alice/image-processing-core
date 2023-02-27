package icc

import (
	"bytes"
)

func JpegIccSegment(raw_profile []byte) ([]byte, error) {

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

	return buf.Bytes(), nil
}
