package ipfix

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

// A Set is a collection of records that have a similar structure,
//      prefixed by a header.  In an IPFIX Message, zero or more Sets
//      follow the Message Header.  There are three different types of
//      Sets: Template Sets, Options Template Sets, and Data Sets.

// A Template Set is a collection of one or more Template Records
//      that have been grouped together in an IPFIX Message.

// A Data Set is one or more Data Records, of the same type, that are
//      grouped together in an IPFIX Message.  Each Data Record is
//      previously defined by a Template Record or an Options Template
//      Record.

// A Data Set is one or more Data Records, of the same type, that are
//      grouped together in an IPFIX Message.  Each Data Record is
//      previously defined by a Template Record or an Options Template
//      Record.

var errInvalidSetType = errors.New("invalid set type")

// TemplateNotDefinedError is returned when no template for input data is found.
type TemplateNotDefinedError struct {
	templateID uint16
}

func (e TemplateNotDefinedError) Error() string {
	return fmt.Sprintf("ipfix template not defined %v", e.templateID)
}

// UnmarshalError describes an error that occurred when data are unmarshalled.
// It wraps an error returned by the binary.Read function.
type UnmarshalError struct {
	err error
}

func (e UnmarshalError) Error() string {
	return fmt.Sprintf("ipfix unmarshal error: %v", e.err)
}

// Unwrap returns the underlying error.
func (e UnmarshalError) Unwrap() error {
	return e.err
}

// Set describes a generic term for a collection of records that have a similar structure.  There are three different
// types of Set: template sets, options template sets, and data sets.  Each of these sets consists of a Set header and
// one or more records.  The exporting process may insert some padding octets, so that the subsequent Set starts at an
// aligned boundary.
// https://tools.ietf.org/html/rfc7011#section-3.3.1, https://tools.ietf.org/html/rfc7011#section-3.3.2
type Set struct {
	SetID   uint16   `rethinkdb:"set_id"`
	Length  uint16   `rethinkdb:"length"` // the length, including set header, all records, and the optional padding
	Records []Record `rethinkdb:"records"`
	RawData []byte   `rethinkdb:"data,omitempty"`
}

// Set ID constants defining valid IPFIX set identifiers (or a minimum valid set id in case of MinimumDataSetID).
const (
	TemplateSetID        = 2   // TemplateSetID identifies a IPFIX Template Set
	OptionsTemplateSetID = 3   // OptionsTemplateSetID identifies a IPFIX Options Template Set
	MinimumDataSetID     = 256 // set ids >= MinimumDataSetID identify IPFIX Data Sets
)

func (s *Set) String() string {
	if !s.valid() {
		return fmt.Sprintf("ipfix invalid set #%v", s.SetID)
	}

	var setGenre string

	switch {
	case s.SetID == TemplateSetID:
		setGenre = "template set"
	case s.SetID == OptionsTemplateSetID:
		setGenre = "options template set"
	case s.SetID >= MinimumDataSetID:
		setGenre = "data set"
	}

	return fmt.Sprintf("ipfix %v (%s) %s", setGenre, intDescribed(int(s.Length), "octet", "octets"),
		intDescribed(len(s.Records), "record", "records"))
}

func (s *Set) valid() bool {
	return (s.SetID == 2 || s.SetID == 3 || s.SetID >= 256) && s.Length > 3
}

// Unmarshal unmarshal a IPFIX Message.
func (s *Set) Unmarshal(reader io.Reader, templates *map[uint16]*TemplateRecord) (n int, err error) {
	n, err = s.unmarshal(reader, templates)

	if err != nil {
		nFromHandler, errFromHandler := s.handleUnmarshalError(reader, n, err)
		if errFromHandler != nil {
			return n + nFromHandler, UnmarshalError{errFromHandler}
		}

		n += nFromHandler
	}

	return n, nil
}

func (s *Set) unmarshal(reader io.Reader, templates *map[uint16]*TemplateRecord) (n int, err error) {
	if n, err = unmarshalFields(reader, &s.SetID, &s.Length); err != nil {
		return n, err
	}

	var bytesRead int

	for n < int(s.Length) {
		switch {
		case s.SetID == TemplateSetID:
			bytesRead, err = s.unmarshalTemplateRecord(reader, templates)

		case s.SetID == OptionsTemplateSetID:
			bytesRead, err = s.unmarshalOptionsTemplateRecord(reader)

		case s.SetID >= MinimumDataSetID:
			bytesRead, err = s.unmarshalDataRecord(reader, *templates)

		default:
			return n, UnmarshalError{errInvalidSetType}
		}

		n += bytesRead

		if err != nil {
			return n, UnmarshalError{err}
		}
	}

	return n, nil
}

func (s *Set) unmarshalTemplateRecord(reader io.Reader, templates *map[uint16]*TemplateRecord) (n int, err error) {
	r := &TemplateRecord{}
	n, err = r.Unmarshal(reader, nil)

	if err != nil {
		return n, err
	}

	s.Records = append(s.Records, r)

	if *templates == nil {
		*templates = make(map[uint16]*TemplateRecord)
	}

	(*templates)[r.TemplateID] = r

	return n, nil
}

func (s *Set) unmarshalOptionsTemplateRecord(reader io.Reader) (n int, err error) {
	r := &OptionsTemplateRecord{}
	n, err = r.Unmarshal(reader, nil)

	if err != nil {
		return n, err
	}

	s.Records = append(s.Records, r)

	return n, nil
}

func (s *Set) unmarshalDataRecord(reader io.Reader, templates map[uint16]*TemplateRecord) (n int, err error) {
	if template, ok := templates[s.SetID]; ok {
		// if template, ok := getTemplate(s.SetID, templates); ok {
		r := &DataRecord{}
		n, err = r.Unmarshal(reader, template)

		if err != nil {
			return n, err
		}

		s.Records = append(s.Records, r)

		return n, nil
	}

	return 0, TemplateNotDefinedError{templateID: s.SetID}
}

func (s *Set) handleUnmarshalError(reader io.Reader, n int, err error) (int, error) {
	switch {
	case errors.Is(err, TemplateNotDefinedError{templateID: s.SetID}):
		// template for the data set not known, store raw data instead
		buf := make([]byte, int(s.Length)-n)
		if e := binary.Read(reader, binary.BigEndian, &buf); e != nil {
			return len(buf), UnmarshalError{e}
		}

		s.RawData = buf

		return len(buf), nil
	case errors.Is(err, io.ErrUnexpectedEOF) && (int(s.Length) > n):
		// assume s.Unmarshal failed to read a record which was actually a padding
		buf := make([]byte, int(s.Length)-n)

		if e := binary.Read(reader, binary.BigEndian, &buf); e != nil {
			return len(buf), UnmarshalError{e}
		}

		return len(buf), nil
	default:
		return 0, err
	}
}

// Record describes one of three ipfix records.
type Record interface {
	Size() int
	Unmarshal(io.Reader, *TemplateRecord) (int, error)
}
