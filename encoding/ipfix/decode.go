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
		messages:            make([]Message, 0),
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

	if d.buf == nil {
		d.buf = make([]byte, len(p))
		copy(d.buf, p)
	} else {
		d.buf = append(d.buf, p...)
	}

	r := bytes.NewReader(d.buf)

	var m Message

	nn, err := m.Unmarshal(r, d.templates)
	d.buf = d.buf[nn:]
	n += nn

	if err != nil && !errors.Is(errors.Unwrap(err), io.EOF) {
		d.err = err

		return n, CorruptInputError{n, errors.Unwrap(err)}
	}

	d.messages = append(d.messages, m)
	d.updateDecoderTemplates(m.sets)

	return n, nil
}

func (d *Decoder) updateDecoderTemplates(sets []set) {
	for _, s := range sets {
		switch s.setID {
		case TemplateSetID:
			d.updateTemplates(s.records)
		case OptionsTemplateSetID:
			d.updateOptionsTemplates(s.records)
		}
	}
}

func (d *Decoder) updateTemplates(records []record) {
	for _, r := range records {
		if template, ok := r.(*templateRecord); ok {
			d.templates = append(d.templates, template)
		}
	}
}

func (d *Decoder) updateOptionsTemplates(records []record) {
	for _, r := range records {
		if template, ok := r.(*optionsTemplateRecord); ok {
			d.optionsTemplates = append(d.optionsTemplates, template)
		}
	}
}
