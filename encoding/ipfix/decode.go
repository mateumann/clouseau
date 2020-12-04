package ipfix

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
)

const (
	// NoSpecificObservationDomainID represents no specific observation domain id relevant for the entire IPFIX Message
	// e.g. when exporting the Exporting Process Statistics, or in the case of a hierarchy of Collectors when
	// aggregated Data Records are exported.
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
	source              *net.UDPAddr
	observationDomainID uint32
	err                 error
	buf                 []byte // leftover decoded output
	messages            []Message
	data                []*dataRecord
	templates           []*templateRecord
	optionsTemplates    []*optionsTemplateRecord
	// r                io.Reader
}

func (d *Decoder) String() string {
	return fmt.Sprintf("Decoder (%s, %d): %s", d.source, d.observationDomainID,
		intDescribed(len(d.messages), "message", "messages"))
}

// NewDecoder initializes a Decoder with provided source (net.UDPAddr) and default values of other parameters.
func NewDecoder(s *net.UDPAddr) *Decoder {
	return &Decoder{
		source:              s,
		observationDomainID: NoSpecificObservationDomainID,
		err:                 nil,
		buf:                 make([]byte, initialDecoderBufSize),
		messages:            make([]Message, 1),
		data:                make([]*dataRecord, 0),
		templates:           make([]*templateRecord, 0),
		optionsTemplates:    make([]*optionsTemplateRecord, 0),
	}
}

// Read decodes a IPFIX Message(s) from p buffer.
func (d *Decoder) Read(p []byte) (n int, err error) {
	if d.err != nil {
		return 0, d.err
	}

	var m Message

	n = 0

	if d.buf == nil {
		d.buf = make([]byte, len(p))
		copy(d.buf, p)
	} else {
		d.buf = append(d.buf, p...)
	}

	r := bytes.NewReader(d.buf)

	nn, err := m.Unmarshal(r, d.templates)
	d.buf = d.buf[nn:]

	if err != nil {
		if errors.Is(errors.Unwrap(err), io.EOF) {
			return n + nn, nil
		}

		d.err = err

		return n + nn, CorruptInputError{n + nn, errors.Unwrap(err)}
	}

	n += nn

	return n, nil
}
