package ipfix

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

// decoder represents IPFIX flow already received by the collecting process.
type decoder struct {
	// err error
	// r                io.Reader
	// buf              []byte // leftover input
	out []byte // leftover decoded output
	// messages         []message
	// data             []*dataRecord
	// templates        []*templateRecord
	// optionsTemplates []*optionsTemplateRecord
}

type CorruptInputError struct {
	position int64
	err      error
}

func (e CorruptInputError) Error() string {
	return fmt.Sprintf("ipfix: error in the input data at byte %v (%v)", e.position, e.err)
}

func (e CorruptInputError) Unwrap() error {
	return e.err
}

// Read decodes a IPFIX message(s) from p buffer (implements a io.Writer interface).
func (d *decoder) Read(r io.Reader) (n int64, err error) {
	var m message

	n = 0

	if len(d.out) > 0 {
		r = io.MultiReader(bytes.NewReader(d.out), r)
	}

	if err := binary.Read(r, binary.BigEndian, &m.version); err != nil {
		return n, CorruptInputError{0, err}
	}

	n += int64(binary.Size(m.version))

	if err := binary.Read(r, binary.BigEndian, &m.length); err != nil {
		return n, CorruptInputError{n, err}
	}

	n += int64(binary.Size(m.length))

	// TODO: read a message header, determine a message size, create buffer, read into a buffer, start parsing
	return n, nil
}
