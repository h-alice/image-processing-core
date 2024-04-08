package png_parser_test

import (
	"bytes"
	. "imagecore/image_parser/png"
	"log"
	"os"
	"testing"
)

func LoadTestCase(c string) []byte {
	fname := ""
	switch c {

	case "logo":
		fname = "../../resource/logo.png"

	}
	raw_bytes, err := os.ReadFile(fname)
	if err != nil {
		log.Fatalln(err)
	}
	return raw_bytes
}

func segment_pretty_print(seg PngSegment, print_data bool) {
	general_seg, _ := seg.(*PngGeneralSegment)
	log.Printf("PNG Segment %s:\n", general_seg.SegmentType)
	log.Printf("\t- Length:%d\n", general_seg.Length)
	if print_data {
		log.Printf("\t- Data:%d\n\n", *general_seg.Data)
	}
}

func TestParsePngSegment(t *testing.T) {

	sample_ihdr := []byte("\x00\x00\x00\x0D\x49\x48\x44\x52\x00\x00\x00\x64\x00\x00\x00\x64\x08\x02\x00\x00\x00\xFF\x80\x02\x03")

	seg := new(PngGeneralSegment)
	read, err := seg.ReadFrom(bytes.NewReader(sample_ihdr))
	if err != nil {
		log.Println(err)
		t.Fail()
	}
	if len(sample_ihdr) != int(read) {
		t.Fail()
	}

	segment_pretty_print(seg, true)
}

func TestPngSegmentRW(t *testing.T) {
	sample_ihdr := []byte("\x00\x00\x00\x0D\x49\x48\x44\x52\x00\x00\x00\x64\x00\x00\x00\x64\x08\x02\x00\x00\x00\xFF\x80\x02\x03")
	seg := new(PngGeneralSegment)
	read, err := seg.ReadFrom(bytes.NewReader(sample_ihdr))
	if err != nil {
		log.Println(err)
		t.Fail()
	}
	if len(sample_ihdr) != int(read) {
		t.Fail()
	}

	outbuf := bytes.NewBuffer([]byte{})
	seg.WriteTo(outbuf)
	log.Println(outbuf.Bytes())
	if !bytes.Equal(sample_ihdr, outbuf.Bytes()) {
		t.Fail()
	}

}
