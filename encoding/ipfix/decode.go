package ipfix

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

// decoder represents IPFIX flow already received by the collecting process.
type decoder struct {
	// err error
	// r                io.Reader
	// buf              []byte // leftover input
	buf []byte // leftover decoded output
	// messages         []message
	// data             []*dataRecord
	// templates        []*templateRecord
	// optionsTemplates []*optionsTemplateRecord
}

type CorruptInputError struct {
	position int
	err      error
}

func (e CorruptInputError) Error() string {
	return fmt.Sprintf("ipfix: error in the input data at byte %v (%v)", e.position, e.err)
}

func (e CorruptInputError) Unwrap() error {
	return e.err
}

// Write decodes a IPFIX message(s) from p buffer (implements a io.Writer interface).
func (d *decoder) Write(p []byte) (n int, err error) {
	var m message

	var nn int

	n = 0

	if d.buf == nil {
		d.buf = make([]byte, len(p))
		copy(d.buf, p)
	} else {
		d.buf = append(d.buf, p...)
	}

	r := bytes.NewReader(d.buf)

	if nn, err = d.fetchElement(r, &m.version); err != nil {
		if errors.Is(errors.Unwrap(err), io.EOF) {
			return n, nil
		}

		return n, CorruptInputError{n, errors.Unwrap(err)}
	}

	n += nn

	if nn, err = d.fetchElement(r, &m.length); err != nil {
		if errors.Is(errors.Unwrap(err), io.EOF) {
			return n, nil
		}

		return n, CorruptInputError{n, errors.Unwrap(err)}
	}

	n += nn

	if nn, err = d.fetchElement(r, &m.exportTime); err != nil {
		if errors.Is(errors.Unwrap(err), io.EOF) {
			return n, nil
		}

		return n, CorruptInputError{n, errors.Unwrap(err)}
	}

	n += nn

	return n, nil
}

func (d *decoder) fetchElement(r io.Reader, data interface{}) (n int, err error) {
	if err := binary.Read(r, binary.BigEndian, data); err != nil {
		return 0, CorruptInputError{1, err}
	}

	n = binary.Size(data)
	d.buf = d.buf[n:]

	return n, nil
}
