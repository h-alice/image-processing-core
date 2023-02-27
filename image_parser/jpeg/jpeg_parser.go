package jpeg

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

// Reader interface, for intrnal use only.
type Reader interface {
	io.ReadSeeker
	io.ByteReader
}

// JPEG segment interface.
type JpegSegment interface {

	// Read full segment from reader.
	Read(Reader) error

	// Return raw bytes of segment.
	RawBytes() ([]byte, error)

	// Get segment marker.
	marker() JpegMarker
}

type JpegGeneralSegment struct {
	signature [2]byte
	length    int
	data      *[]byte
}

// Read JPEG segment.
func (chunk *JpegGeneralSegment) Read(reader Reader) error {

	var err error = nil

	var reader_init_pos int64 = 0
	reader_init_pos, err = reader.Seek(0, io.SeekCurrent)

	if err != nil { // Unread and return.
		reader.Seek(reader_init_pos, io.SeekStart) // Unread.
		return err
	}

	// Read signature
	raw_sig := make([]byte, 2)
	_, err = io.ReadFull(reader, raw_sig)

	if err != nil || raw_sig[0] != '\xFF' { // Unread and return.

		switch {

		case err != nil:

		case raw_sig[0] != '\xFF':
			err = fmt.Errorf("invalid jpeg marker")
		}

		reader.Seek(reader_init_pos, io.SeekStart) // Unread.
		return err
	}

	chunk.signature = [2]byte(raw_sig)

	switch JpegMarker(raw_sig[1]) {

	case SOI, EOI, RST0, RST1, RST2, RST3, RST4, RST5, RST6, RST7:
		// Parameter-less Marker
		// NOTE: RST marker shouldn't appear here.
		chunk.length = 0
		chunk.data = nil

		return nil

	default:
		break
	}

	// Parse length
	var length int = 0
	raw_length := make([]byte, 2)
	_, err = io.ReadFull(reader, raw_length)
	length = int(binary.BigEndian.Uint16(raw_length))

	if err != nil { // Unread and return.
		reader.Seek(reader_init_pos, io.SeekStart) // Unread.
		return err
	}

	// Read data
	var data_length = length - 2
	data := make([]byte, (data_length))
	_, err = io.ReadFull(reader, data)

	if err != nil { // Unread and return.
		reader.Seek(reader_init_pos, io.SeekStart) // Unread.
		return err
	}

	chunk.length = length
	chunk.data = &data

	return nil
}

// Convert JPEG segment to raw bytes.
func (chunk *JpegGeneralSegment) RawBytes() ([]byte, error) {

	var err error = nil

	var buf bytes.Buffer

	var written int

	// Write signature
	written, err = buf.Write(chunk.signature[:])
	if err != nil || written != 2 {
		if err == nil {
			err = fmt.Errorf("can't serialize segment signature")
		}
		return nil, err
	}

	if chunk.length == 0 { // Parameter-less
		return buf.Bytes(), nil
	}

	// Write length
	_len := make([]byte, 2)
	binary.BigEndian.PutUint16(_len, uint16(chunk.length))
	written, err = buf.Write(_len)
	if err != nil || written != 2 {
		if err == nil {
			err = fmt.Errorf("can't serialize segment length")
		}
		return nil, err
	}

	// Write data
	written, err = buf.Write(*chunk.data)
	if err != nil || (written != (chunk.length - 2)) {
		if err == nil {
			err = fmt.Errorf("can't serialize segment data")
		}
		return nil, err
	}

	return buf.Bytes(), nil
}

func (chunk *JpegGeneralSegment) marker() JpegMarker {
	return JpegMarker(chunk.signature[1])
}

// Check if a segment is EOI.
func IsEndOfImage(seg JpegSegment) bool {
	return JpegMarker(seg.marker()) == EOI
}

type JpegEcsSegment struct {
	data *[]byte
}

// Read ECS segment.
func (chunk *JpegEcsSegment) Read(reader Reader) error {

	var err error

	var reader_init_pos int64 = 0
	reader_init_pos, err = reader.Seek(0, io.SeekCurrent)

	if err != nil { // Unread and return.
		return err
	}

	var _data bytes.Buffer

	for {

		var b byte
		b, err = reader.ReadByte()
		if err != nil {
			reader.Seek(reader_init_pos, io.SeekStart) // Unread
			return err
		}

		if b == '\xFF' {

			var tmp byte
			tmp, err = reader.ReadByte()
			if err != nil {
				reader.Seek(reader_init_pos, io.SeekStart) // Unread
				return err
			}
			switch JpegMarker(tmp) {

			case NUL, RST0, RST1, RST2, RST3, RST4, RST5, RST6, RST7, '\xFF':
				err = _data.WriteByte(b)
				if err != nil {
					reader.Seek(reader_init_pos, io.SeekStart) // Unread
					return err
				}

				err = _data.WriteByte(tmp)
				if err != nil {
					reader.Seek(reader_init_pos, io.SeekStart) // Unread
					return err
				}
				continue
			default:
				// Encountered a marker.
				reader.Seek(-2, io.SeekCurrent)
				_data.Bytes()

				data := _data.Bytes()
				chunk.data = &data

				return nil
			}
		} else {
			err = _data.WriteByte(b) // Default: append byte to buffer.
			if err != nil {
				reader.Seek(reader_init_pos, io.SeekStart) // Unread
				return err
			}
		}
	}
}

// Convert ECS segment to raw bytes.
func (chunk *JpegEcsSegment) RawBytes() ([]byte, error) {

	data_length := len(*(*chunk).data)
	r := make([]byte, data_length)
	written := copy(r, *(*chunk).data)

	if written != data_length {
		err := fmt.Errorf("ecs segment buffer length mismatch")
		return nil, err
	}

	return r, nil
}

type JpegSosSegment struct {
	sos_segment *JpegGeneralSegment
	ecs         *JpegEcsSegment
}

func (chunk *JpegSosSegment) Read(reader Reader) error {
	var err error = nil

	c := new(JpegGeneralSegment)
	err = (*c).Read(reader)
	if err != nil {
		return err
	}

	ecs := new(JpegEcsSegment)
	err = ecs.Read(reader)
	if err != nil {
		return err
	}

	chunk.sos_segment = c
	chunk.ecs = ecs
	return err
}

func (chunk *JpegSosSegment) RawBytes() ([]byte, error) {

	var err error = nil

	var buf bytes.Buffer

	// Serialize SOS chunk.
	var sos_chunk []byte
	sos_chunk, err = chunk.sos_segment.RawBytes()
	if err != nil {
		return nil, err
	}

	written, err := buf.Write(sos_chunk)
	switch {
	case err != nil:
		return nil, err

	case written != len(sos_chunk):
		err = fmt.Errorf("cannot serialize sos segment")
		return nil, err

	default:
		break
	}

	// Serialize ECS chunk.
	var ecs_chunk []byte
	ecs_chunk, err = chunk.ecs.RawBytes()
	if err != nil {
		return nil, err
	}
	written, err = buf.Write(ecs_chunk)
	switch {
	case err != nil:
		return nil, err

	case written != len(ecs_chunk):
		err = fmt.Errorf("cannot serialize ecs segment")
		return nil, err

	default:
		break
	}

	return buf.Bytes(), nil
}

func (chunk *JpegSosSegment) marker() JpegMarker {
	return chunk.sos_segment.marker()
}

func JpegReadChunk(r Reader) (JpegSegment, error) {

	var p_chunk JpegSegment = nil
	var err error = nil

	sig := make([]byte, 2)
	_, err = io.ReadFull(r, sig)
	r.Seek(-2, io.SeekCurrent) // Unread signature.

	switch err {
	case nil:
		break
	case io.EOF:
		fmt.Println("[!] JPEG segment parser: End of file.")
		return nil, io.EOF
	default:
		return nil, err
	}

	switch JpegMarker(sig[1]) {

	case SOS:
		p_chunk = new(JpegSosSegment)

	default:
		p_chunk = new(JpegGeneralSegment)

	}

	err = p_chunk.Read(r)

	if err != nil {
		return nil, err
	}

	return p_chunk, nil
}

// Return a general JPEGsegment.
func NewJpegSegment(marker JpegMarker, data []byte) *JpegGeneralSegment {

	ptr_segment := new(JpegGeneralSegment)

	// Set marker.
	_marker := [2]byte{'\xFF', byte(marker)}

	// Calculate length.
	var l int = len(data)

	_data := make([]byte, l)
	copy(_data, data)

	ptr_segment = &JpegGeneralSegment{
		signature: _marker,
		length:    (l + 2),
		data:      &_data,
	}
	return ptr_segment
}

type ParsedJpeg struct {
	segments []JpegSegment
}

func (f ParsedJpeg) WriteTo(w io.Writer) (written int64, err error) {
	err = nil
	written = 0
	for _, _seg := range f.segments {
		var tmp []byte
		var _written int64
		tmp, err = _seg.RawBytes()
		if err != nil {
			return
		}
		_written, err = io.Copy(w, bytes.NewReader(tmp))
		written += _written
		if err != nil {
			return
		}
	}
	return
}

func (f *ParsedJpeg) InsertIccSegment(icc_segment *JpegGeneralSegment) (err error) {

	err = nil

	if icc_segment.marker() != APP2 {
		err = fmt.Errorf("not an icc segment")
		return
	}

	// Insert ICC profile after APP0 and APP1.
	for index, _seg := range f.segments {
		if _seg.marker() == SOI || _seg.marker() == APP0 || _seg.marker() == APP1 {
			continue
		} else if _seg.marker() == APP2 {
			return
		} else {
			//tmp := make([]JpegSegment, 0)
			//tmp := append(f.segments[0:index], icc_segment)
			//tmp = append(tmp, f.segments[index:]...)
			f.segments = append(f.segments[:index+1], f.segments[index:]...)
			f.segments[index] = icc_segment
			break
		}
	}
	return
}

func ParseJpeg(reader Reader) (*ParsedJpeg, error) {

	var err error = nil
	var parsed *ParsedJpeg
	segments := make([]JpegSegment, 0)

	for {
		var chunk JpegSegment

		chunk, err = JpegReadChunk(reader)
		if err != nil {
			return nil, err
		}

		segments = append(segments, chunk)

		// println(chunk)
		if IsEndOfImage(chunk) {
			break
		}
	}
	parsed = &ParsedJpeg{segments: segments}
	return parsed, nil
}
