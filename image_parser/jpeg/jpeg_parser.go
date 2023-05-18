package jpeg_parser

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

var (
	ErrInvalidJpegHeader  = errors.New("invalid jpeg header")
	ErrInvalidSegmentType = errors.New("invalid jpeg segment type")
)

type JpegSegment interface {
	ReadFrom(io.Reader) (int64, error)
	WriteTo(io.Writer) (int64, error)
}

type JpegRawSegment struct {
	Data *[]byte
}

type JpegEcsSegment = JpegRawSegment

func (seg *JpegEcsSegment) ReadEcsSegment(readseeker io.ReadSeeker) (int64, error) {

	tmp_hi := make([]byte, 1)           // High byte placeholder.
	tmp_lo := make([]byte, 1)           // Low byte placeholder.
	ecsbuf := bytes.NewBuffer([]byte{}) // Raw data buffer.
	total_read := int64(0)              // Total byte read.

	for {

		// Read high byte
		read_byte, err := io.ReadFull(readseeker, tmp_hi)
		total_read += int64(read_byte)
		if err != nil {
			return total_read, err
		}

		if tmp_hi[0] == '\xFF' { //
			read_byte, err := io.ReadFull(readseeker, tmp_lo)
			total_read += int64(read_byte)
			if err != nil {
				return total_read, err
			}

			switch tmp_lo[0] { // Check low byte.

			case jpegNUL, jpegRST0, jpegRST1, jpegRST2, jpegRST3, jpegRST4, jpegRST5, jpegRST6, jpegRST7, '\xFF':
				// Not a header signature -> write 2 bytes to buffer.

				err := ecsbuf.WriteByte(tmp_hi[0]) // Write high byte.
				if err != nil {
					return total_read, err
				}

				err = ecsbuf.WriteByte(tmp_lo[0]) // Write low byte
				if err != nil {
					return total_read, err
				}
			default:
				// Encountered a marker.
				readseeker.Seek(-2, io.SeekCurrent) // Unread 2 bytes.
				raw_data := ecsbuf.Bytes()
				seg.Data = &raw_data
				total_read -= 2
				return total_read, err
			}
		} else { // High byte is not '\xFF' -> write to buffer.
			err := ecsbuf.WriteByte(tmp_hi[0])
			if err != nil {
				return total_read, err
			}
		}
	}
}

func (seg *JpegRawSegment) ReadFrom(rd io.Reader) (int64, error) {

	var buf bytes.Buffer
	n, err := io.Copy(&buf, rd)
	if err != nil {
		return n, err
	}

	if n == 0 {
		return n, io.EOF
	} else {
		raw_data := buf.Bytes()
		seg.Data = &raw_data
		return n, err
	}
}

func (seg JpegRawSegment) WriteTo(wt io.Writer) (int64, error) {

	_len := len(*seg.Data)
	n, err := wt.Write(*seg.Data)
	if err != nil {
		return int64(n), err
	}
	if n != _len {
		return int64(n), io.ErrShortWrite
	}

	return int64(n), err
}

type JpegGeneralSegment struct {
	SegmentType byte
	Length      int
	Data        *[]byte
}

// Read one segment from reader.
func (seg *JpegGeneralSegment) ReadFrom(rd io.Reader) (int64, error) {

	var total_read_bytes int64 = 0

	// Read signature high byte.
	sig := make([]byte, 1)         // One byte placeholder.
	n, err := io.ReadFull(rd, sig) // Read one byte.
	total_read_bytes += int64(n)   // Add read counter.
	if err != nil {
		return total_read_bytes, err
	}

	if sig[0] != '\xFF' { // Inavlid segment.
		return total_read_bytes, ErrInvalidJpegHeader
	}

	// Read signature low byte.
	n, err = io.ReadFull(rd, sig) // Read one byte.
	total_read_bytes += int64(n)  // Add read counter.
	if err != nil {
		return total_read_bytes, err
	}
	seg.SegmentType = sig[0] // Set segment type byte.

	// Check segment type
	switch JpegMarker(seg.SegmentType) {

	case jpegSOI, jpegEOI, jpegRST0, jpegRST1, jpegRST2, jpegRST3, jpegRST4, jpegRST5, jpegRST6, jpegRST7:
		// Parameter-less Marker
		// NOTE: RST marker shouldn't appear here.
		seg.Length = 0
		seg.Data = nil

		return total_read_bytes, err

	default:
		break
	}

	// Read length
	var length int = 0                   // Parsed length placeholder.
	raw_length := make([]byte, 2)        // Raw length bytes placeholder.
	n, err = io.ReadFull(rd, raw_length) // Read 2-bytes raw length.
	total_read_bytes += int64(n)         // Add read counter.
	if err != nil {
		return total_read_bytes, err
	}

	length = int(binary.BigEndian.Uint16(raw_length)) // Convert bytes to int.
	seg.Length = length                               // Set length.

	// Read data
	var data_length = seg.Length - 2    // Data length = (segment length) - 2
	data := make([]byte, (data_length)) // Data placeholder.
	n, err = io.ReadFull(rd, data)      // Read data.
	total_read_bytes += int64(n)        // Add read counter.
	if err != nil {
		return total_read_bytes, err
	}

	seg.Data = &data // Set data.

	return total_read_bytes, err
}

func (seg JpegGeneralSegment) WriteTo(wt io.Writer) (int64, error) {

	var total_written_bytes int64 = 0

	// Write signature high byte.
	n, err := wt.Write([]byte{'\xFF'})
	total_written_bytes += int64(n)
	if err != nil {
		return total_written_bytes, err
	}
	if n != 1 {
		return total_written_bytes, io.ErrShortWrite
	}

	// Write signature low byte.
	n, err = wt.Write([]byte{seg.SegmentType})
	total_written_bytes += int64(n)
	if err != nil {
		return total_written_bytes, err
	}
	if n != 1 {
		return total_written_bytes, io.ErrShortWrite
	}

	// If no data to write, return.
	if seg.Length == 0 {
		return total_written_bytes, err
	}

	// Write length.
	_len := make([]byte, 2)
	binary.BigEndian.PutUint16(_len, uint16(seg.Length))
	n, err = wt.Write(_len)
	total_written_bytes += int64(n)
	if err != nil {
		return total_written_bytes, err
	}
	if n != 2 {
		return total_written_bytes, io.ErrShortWrite
	}

	// Write data
	n, err = wt.Write(*seg.Data)
	total_written_bytes += int64(n)
	if err != nil {
		return total_written_bytes, err
	}
	if n != (seg.Length - 2) {
		return total_written_bytes, io.ErrShortWrite
	}

	return total_written_bytes, err
}
