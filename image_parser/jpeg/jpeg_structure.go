package jpeg_parser

import (
	"bytes"
	"errors"
	"io"
	"log"
	"strings"

	"golang.org/x/exp/slices"
)

var (
	ErrInvalidAppSegmentIndex = errors.New("invalid app segment")
)

type JpegImage struct {
	Segments []JpegSegment
}

// Read JPEG and convert into segment list.
func ReadJpeg(r io.Reader) ([]JpegSegment, int64, error) {

	// Read whole file into buffer.
	buf := bytes.NewBuffer([]byte{}) // Create file content placeholder.
	file_len, err := io.Copy(buf, r) // Copy all bytes.
	if err != nil {
		return nil, 0, err
	}

	ret := make([]JpegSegment, 0)

	readseeker := bytes.NewReader(buf.Bytes()) // New reader from JPEG bytes.

	total_read := 0 // Total read byte counter.

	for { // Read loop
		tmp := JpegGeneralSegment{}                 // Empty general segment.
		read_bytes, err := tmp.ReadFrom(readseeker) // Read until segment end.
		total_read += int(read_bytes)               // Add counter.
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, int64(total_read), err
		}

		// Process segment.
		switch tmp.SegmentType {

		case jpegSOS: // Encountered start of scan segment.
			// Read ECS segment.
			tmp_ecs := new(JpegEcsSegment)                       // Empty ECS segment.
			read_bytes, err = tmp_ecs.ReadEcsSegment(readseeker) // Read
			total_read += int(read_bytes)                        // Add read counter.
			if err != nil {
				return nil, int64(total_read), err
			}
			ret = append(ret, &tmp)    // Append SOS segment.
			ret = append(ret, tmp_ecs) // Append ECS segment.
			continue

		case jpegEOI: // Encountered end of image segment.
			ret = append(ret, &tmp)                           // Append EOI segment.
			remainbyte := new(JpegRawSegment)                 // Remain bytes placeholder.
			read_byte, err := remainbyte.ReadFrom(readseeker) // Read all remain bytes.
			total_read += int(read_byte)                      // Add counter.
			if err != nil {
				if err == io.EOF { // Exit read loop if EOF.
					break
				} else {
					return nil, int64(total_read), err
				}
			}
			ret = append(ret, remainbyte) // Append segment

		default:
			ret = append(ret, &tmp) // Append segment
		}
	}

	if total_read != int(file_len) {
		// TODO: checker
		log.Println("Total read number mismatch.")

	}

	return ret, int64(total_read), err
}

// Create new general segment.
func NewGeneralSegment(segment_type byte, data []byte) *JpegGeneralSegment {

	seg := new(JpegGeneralSegment) // Create new general segment.
	seg.SegmentType = segment_type // Set segment type byte.
	if len(data) == 0 {            // Return if no data section.
		return seg
	}
	_data := make([]byte, len(data)) // Allocate data placeholder.
	copy(_data, data)                // Copy data.
	seg.Length = len(_data) + 2      // Set segment length (len(data) + 2)
	seg.Data = &_data                // Set segment data.

	return seg
}

// Create new general segment from segment name.
func NewGeneralSegmentFromSegmentName(segment_type string, data []byte) (*JpegGeneralSegment, error) {
	segment_type = strings.ToUpper(segment_type)    // Convert to upper case.
	seg_value, ok := jpeg_marker_byte[segment_type] // Get segment signature byte.
	if !ok {                                        // If segment type not exists.
		return nil, ErrInvalidSegmentType
	}

	return NewGeneralSegment(seg_value, data), nil
}

// Read JPEG from reader.
func (img *JpegImage) ReadFrom(r io.Reader) (int64, error) {
	seg_list, total_read, err := ReadJpeg(r)
	if err != nil {
		if err != io.EOF {
			return total_read, err
		}
	}

	img.Segments = seg_list
	return total_read, nil
}

// Convert parsed JPEG file to data stream.
func (img JpegImage) WriteTo(wt io.Writer) (int64, error) {
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

func (img *JpegImage) Write(p []byte) (int, error) {
	rd := bytes.NewReader(p)
	total_read, err := img.ReadFrom(rd)
	return int(total_read), err
}

func (img JpegImage) Read(p []byte) (int, error) {
	buf := bytes.NewBuffer([]byte{})
	img.WriteTo(buf)
	return io.ReadFull(buf, p)
}

// Insert APP(0-15) segment into image.
func (im *JpegImage) InsertAppSement(app_index int, data []byte) error {

	if app_index < 0 || app_index > 15 { // Check if app_index in is valid range.
		return ErrInvalidAppSegmentIndex // Return error.
	}

	target_segment_type := jpegAPP0 + byte(app_index) // Convert app_index to byte signature.

	// Scan full image and find target segment type in parsed JPEG.
	target_index := slices.IndexFunc(im.Segments, func(elem JpegSegment) bool {
		switch _t := elem.(type) { // Checking segment type.
		case *JpegGeneralSegment:
			if _t.SegmentType == target_segment_type {
				return true // Found segment, return true.
			} else {
				return false //Not found, return false.
			}
		case *JpegRawSegment:
			return false // Raw segment can't be any type of APP segment, return false.
		default:
			return false // Invalid type, return false.
		}
	})
	if target_index != -1 { // If target segment type exists, replace it.
		im.Segments[target_index] = NewGeneralSegment(target_segment_type, data) // Create new segment and replace current one.
		return nil                                                               // Return
	}

	// Find insert position.
	//  New APP segment will be inserted after SOI
	//  and before any other non-APP segments.
	target_index = slices.IndexFunc(im.Segments, func(elem JpegSegment) bool {
		switch _t := elem.(type) { // Check and cast segment type.
		case *JpegGeneralSegment: // General segment.
			switch _t.SegmentType { // Check segment signature.
			// Start of image.
			case jpegSOI: // Continue search if encounted SOI.
				return false // Return false, continue search.

			// APP segments.
			case jpegAPP0, jpegAPP1, jpegAPP2, jpegAPP3, jpegAPP4, jpegAPP5, jpegAPP6, jpegAPP7, jpegAPP8, jpegAPP9, jpegAPP10, jpegAPP11, jpegAPP12, jpegAPP13, jpegAPP14, jpegAPP15:
				// Insert in ascending order.
				if target_segment_type < _t.SegmentType {
					return true
				} else {
					return false
				}
			// Non-APP segments.
			default:
				return true // Stop search loop.
			}
		case *JpegRawSegment:
			return true // Stop if encounted raw segment.
		default:
			return true // Stop if encounted unknown segment type.
		}
	})

	// Insert segment.
	im.Segments = slices.Insert(im.Segments, target_index, JpegSegment(NewGeneralSegment(target_segment_type, data)))
	return nil
}

func (im *JpegImage) EmbedIccProfile(icc_profile []byte) error {

	buf := bytes.NewBuffer([]byte{})

	written, err := buf.Write(append([]byte("ICC_PROFILE"), []byte{'\x00', '\x01', '\x01'}...)) // ICC chunk header.
	if err != nil {
		return err
	}
	if written != 14 {
		return io.ErrShortWrite
	}

	written, err = buf.Write(icc_profile)
	if err != nil {
		return err
	}
	if written != len(icc_profile) {
		return io.ErrShortWrite
	}

	return im.InsertAppSement(2, buf.Bytes())
}
