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
	err error
	// r                io.Reader
	// buf              []byte // leftover input
	buf []byte // leftover decoded output
	// messages         []message
	// data             []*dataRecord
	// templates        []*templateRecord
	// optionsTemplates []*optionsTemplateRecord
}

// CorruptInputError describes an error in input data.
type CorruptInputError struct {
	position int
	err      error
}

func (e CorruptInputError) Error() string {
	return fmt.Sprintf("ipfix: error in the input data at byte %v (%v)", e.position, e.err)
}

// Unwrap returns the underlying error.
func (e CorruptInputError) Unwrap() error {
	return e.err
}

// Write decodes a IPFIX message(s) from p buffer (implements a io.Writer interface).
func (d *decoder) Write(p []byte) (n int, err error) {
	if d.err != nil {
		return 0, d.err
	}

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

	if nn, err = d.fetchElements(r,
		&m.version, &m.length, &m.exportTime, &m.sequenceNumber, &m.observationDomainID); err != nil {
		if errors.Is(errors.Unwrap(err), io.EOF) {
			return n + nn, nil
		}

		d.err = err

		return n + nn, CorruptInputError{n + nn, errors.Unwrap(err)}
	}

	n += nn

	return n, nil
}

func (d *decoder) fetchElements(r io.Reader, data ...interface{}) (n int, err error) {
	n = 0

	for _, elem := range data {
		if err := binary.Read(r, binary.BigEndian, elem); err != nil {
			return n, CorruptInputError{n + 1, err}
		}

		nn := binary.Size(elem)
		d.buf = d.buf[nn:]
		n += nn
	}

	return n, nil
}
