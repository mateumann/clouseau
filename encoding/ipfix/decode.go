package ipfix

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
)

const (
	NoSpecificObservationDomainID = 0
	initialDecoderBufSize         = 2048
)

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

// Decoder represents IPFIX flow already received by the collecting process.
type Decoder struct {
	source              net.UDPAddr
	observationDomainID uint32
	err                 error
	buf                 []byte // leftover decoded output
	// r                io.Reader
	// messages         []Message
	// data             []*dataRecord
	// templates        []*templateRecord
	// optionsTemplates []*optionsTemplateRecord
}

func NewDecoder(s net.UDPAddr) *Decoder {
	d := Decoder{
		source:              s,
		observationDomainID: NoSpecificObservationDomainID,
		buf:                 make([]byte, initialDecoderBufSize),
	}

	return &d
}

// Write decodes a IPFIX Message(s) from p buffer (implements a io.Writer interface).
func (d *Decoder) Write(p []byte) (n int, err error) {
	if d.err != nil {
		return 0, d.err
	}

	var m Message

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

func (d *Decoder) fetchElements(r io.Reader, data ...interface{}) (n int, err error) {
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
