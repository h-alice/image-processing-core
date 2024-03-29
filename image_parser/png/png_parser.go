package png_parser

import (
	"encoding/binary"
	"errors"
	"hash/crc32"
	"io"
	"log"
)

var (
	ErrCrcCheckFailed = errors.New("segment checksum failed")
)

type PngSegment interface {
	ReadFrom(io.Reader) (int64, error)
	WriteTo(io.Writer) (int64, error)
}

type PngGeneralSegment struct {
	SegmentType string
	Length      int
	Data        *[]byte
}

func (seg PngGeneralSegment) Crc32Checksum() uint32 {

	return crc32.Checksum(append([]byte(seg.SegmentType), *(seg.Data)...), crc32.MakeTable(crc32.IEEE))
}

func (seg *PngGeneralSegment) ReadFrom(reader io.Reader) (int64, error) {

	total_read := int64(0)
	seg_len := make([]byte, 4) // Segment length placeholder.
	read, err := io.ReadFull(reader, seg_len)
	total_read += int64(read)

	if err != nil {
		return total_read, err
	}

	length := int(binary.BigEndian.Uint32(seg_len)) // Convert bytes to int.
	seg.Length = length                             // Set length.
	log.Println(length)

	segment_type := make([]byte, 4) // Segment type placeholder.
	read, err = io.ReadFull(reader, segment_type)
	total_read += int64(read)
	if err != nil {
		return total_read, err
	}

	seg.SegmentType = string(segment_type)

	data := make([]byte, seg.Length) // Segment data placeholder.
	read, err = io.ReadFull(reader, data)
	total_read += int64(read)
	if err != nil {
		return total_read, err
	}

	seg.Data = &data // Set segment data.

	crc_checksum := make([]byte, 4)
	read, err = io.ReadFull(reader, crc_checksum) // Read CRC32 checksum.
	total_read += int64(read)
	if err != nil {
		return total_read, err
	}

	// Check CRC.
	if binary.BigEndian.Uint32(crc_checksum) != seg.Crc32Checksum() {
		return total_read, ErrCrcCheckFailed
	}

	return total_read, nil
}

func (seg PngGeneralSegment) WriteTo(writer io.Writer) (int64, error) {

	total_written := int64(0)

	// Write length.
	_len := make([]byte, 4)
	binary.BigEndian.PutUint32(_len, uint32(len(*seg.Data)))
	n, err := writer.Write(_len)
	total_written += int64(n)
	if err != nil {
		return total_written, err
	}
	if n != 4 {
		return total_written, io.ErrShortWrite
	}

	// Write signature.
	n, err = writer.Write([]byte(seg.SegmentType))
	total_written += int64(n)
	if err != nil {
		return total_written, err
	}
	if n != 4 {
		return total_written, io.ErrShortWrite
	}

	//Write data.
	n, err = writer.Write(*seg.Data)
	total_written += int64(n)
	if err != nil {
		return total_written, err
	}
	if n != seg.Length {
		return total_written, io.ErrShortWrite
	}

	// Write CRC32.
	_crc := make([]byte, 4)
	binary.BigEndian.PutUint32(_crc, seg.Crc32Checksum())

	n, err = writer.Write(_crc)
	total_written += int64(n)
	if err != nil {
		return total_written, err
	}
	if n != 4 {
		return total_written, io.ErrShortWrite
	}

	return total_written, nil
}
