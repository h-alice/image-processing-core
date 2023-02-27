package imageparser

import (
	"io"
)

// Reader interface, for intrnal use only.
type reader interface {
	io.ReadSeeker
	io.ByteReader
}

type ParserdImage interface {
}
