package jpeg_parser_test

import (
	"bytes"
	"fmt"
	. "imagetools/image_parser/jpeg"
	"io"
	"log"
	"os"
	"testing"
)

func load_test_case(c string) []byte {
	fname := ""
	switch c {
	case "small":
		fname = "../../resource/test_small_jpg.jpg"
	case "logo":
		fname = "../../resource/doit-logo.jpg"

	}
	raw_bytes, err := os.ReadFile(fname)
	if err != nil {
		log.Fatalln(err)
	}
	return raw_bytes
}

func TestParseSmallJpeg(t *testing.T) {

	raw_bytes := load_test_case("small")
	img := new(JpegImage)
	_, err := img.ReadFrom(bytes.NewReader(raw_bytes))
	if err != nil {
		t.Fail()
	}

	for pos, x := range img.Segments {
		switch _t := x.(type) {
		default:
			fmt.Printf("Unexpected type: %v\n", _t)
			t.Fail()
		case *JpegGeneralSegment:
			fmt.Printf("[%d] Segment [%X]:\n\tLength: %d\n\n", pos, _t.SegmentType, _t.Length)
		case *JpegRawSegment:
			fmt.Printf("[%d] Segment [RAW]:\n\tLength: %d\n\n", pos, len(*_t.Data))
		}
	}
}

func TestParseJpeg(t *testing.T) {

	raw_bytes := load_test_case("logo")
	img := new(JpegImage)
	_, err := img.ReadFrom(bytes.NewReader(raw_bytes))
	if err != nil {
		t.Fail()
	}

	for pos, x := range img.Segments {
		switch _t := x.(type) {
		default:
			fmt.Printf("Unexpected type: %v\n", _t)
			t.Fail()
		case *JpegGeneralSegment:
			fmt.Printf("[%d] Segment [%X]:\n\tLength: %d\n\n", pos, _t.SegmentType, _t.Length)
		case *JpegRawSegment:
			fmt.Printf("[%d] Segment [RAW]:\n\tLength: %d\n\n", pos, len(*_t.Data))
		}
	}
}

func TestJpegParseRebuild(t *testing.T) {
	raw_bytes := load_test_case("logo")

	img := new(JpegImage)
	_, err := img.ReadFrom(bytes.NewReader(raw_bytes))
	if err != nil {
		t.Fail()
	}

	buf := bytes.NewBuffer([]byte{})
	_, err = img.WriteTo(buf)
	if err != nil {
		t.Fail()
	}

	// Compare
	if !bytes.Equal(buf.Bytes(), raw_bytes) {
		t.Fail()
	}
}

func TestJpegIoImpl(t *testing.T) {
	raw_bytes := load_test_case("logo")

	img := new(JpegImage)
	_, err := io.Copy(img, bytes.NewReader(raw_bytes))
	if err != nil {
		t.Fail()
	}

	buf := bytes.NewBuffer([]byte{})
	_, err = io.Copy(buf, img)
	if err != nil {
		t.Fail()
	}

	// Compare
	if !bytes.Equal(buf.Bytes(), raw_bytes) {
		t.Fail()
	}
}

func TestNewSegment(t *testing.T) {
	sample_jfif := []byte("\xFF\xE0\x00\x10\x4A\x46\x49\x46\x00\x01\x01\x01\x00\x78\x00\x78\x00\x00")
	sample_jfif_data := []byte("\x4A\x46\x49\x46\x00\x01\x01\x01\x00\x78\x00\x78\x00\x00")
	seg, err := NewGeneralSegmentFromSegmentName("APP0", sample_jfif_data)
	if err != nil {
		t.Fail()
	}
	buf := bytes.NewBuffer([]byte{})
	seg.WriteTo(buf)
	if !bytes.Equal(buf.Bytes(), sample_jfif) {
		t.Fail()
	}
}

func TestJpegInsertAppSeg1(t *testing.T) {

	pretty_print := func(l []JpegSegment) {
		for pos, x := range l {
			switch _t := x.(type) {
			default:
				fmt.Printf("Unexpected type: %v\n", _t)
				t.Fail()
			case *JpegGeneralSegment:

				var data []byte
				if _t.Data == nil {
					data = nil
				} else {
					data = *_t.Data
				}
				fmt.Printf("[%d] Segment [%X]:\n\tLength: %d\n\tData: %v\n\n", pos, _t.SegmentType, _t.Length, data)
			case *JpegRawSegment:
				fmt.Printf("[%d] Segment [RAW]:\n\tLength: %d\n\n", pos, len(*_t.Data))
			}
		}
	}

	raw_bytes := load_test_case("logo")

	img := new(JpegImage)
	_, err := io.Copy(img, bytes.NewReader(raw_bytes))
	if err != nil {
		t.Fail()
	}

	img.InsertAppSement(5, []byte("12345"))
	img.InsertAppSement(4, []byte("67890"))
	img.InsertAppSement(3, []byte("abcdefg"))
	img.InsertAppSement(6, []byte("abcdefg"))
	img.InsertAppSement(1, []byte("qwertyuiop"))
	img.InsertAppSement(1, []byte("asdfghjk"))
	img.InsertAppSement(1, []byte("zxcvbhyu"))
	pretty_print(img.Segments)

}
func TestJpegInsertICC(t *testing.T) {
	pretty_print := func(l []JpegSegment) {
		for pos, x := range l {
			switch _t := x.(type) {
			default:
				fmt.Printf("Unexpected type: %v\n", _t)
				t.Fail()
			case *JpegGeneralSegment:

				var data []byte
				if _t.Data == nil {
					data = nil
				} else {
					data = *_t.Data
				}
				fmt.Printf("[%d] Segment [%X]:\n\tLength: %d\n\tData: %v\n\n", pos, _t.SegmentType, _t.Length, data)
			case *JpegRawSegment:
				fmt.Printf("[%d] Segment [RAW]:\n\tLength: %d\n\n", pos, len(*_t.Data))
			}
		}
	}

	raw_bytes := load_test_case("logo")

	img := new(JpegImage)
	_, err := io.Copy(img, bytes.NewReader(raw_bytes))
	if err != nil {
		t.Fail()
	}

	test_icc_profile := []byte("1234567890abcdef")
	err = img.EmbedIccProfile(test_icc_profile)
	if err != nil {
		t.Fail()
	}

	pretty_print(img.Segments)
}
