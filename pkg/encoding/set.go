package encoding

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

// A Template is an ordered sequence of <type, length> pairs used to
//      completely specify the structure and semantics of a particular set
//      of information that needs to be communicated from an IPFIX Device
//      to a Collector.  Each Template is uniquely identifiable by means
//      of a Template ID.

// An IPFIX Message is a Message that originates at the Exporting
//      Process and carries the IPFIX records of this Exporting Process,
//      and whose destination is a Collecting Process.  An IPFIX Message
//      is encapsulated at the transport layer.

// The Message Header is the first part of an IPFIX Message; the
//      Message Header provides basic information about the Message, such
//      as the IPFIX version, length of the Message, Message sequence
//      number, etc.

// A Template Record defines the structure and interpretation of
//      fields in a Data Record.

// A Data Record is a record that contains values of the parameters
//      corresponding to a Template Record.

// An Options Template Record is a Template Record that defines the
//      structure and interpretation of fields in a Data Record, including
//      defining how to scope the applicability of the Data Record.

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

// An Information Element is a protocol- and encoding-independent
//      description of an attribute that may appear in an IPFIX Record.
//      Information Elements are defined in the IANA "IPFIX Information
//      Elements" registry [IANA-IPFIX].  The type associated with an
//      Information Element indicates constraints on what it may contain
//      and also determines the valid encoding mechanisms for use in
//      IPFIX.

// Figure A shows a summary of IPFIX terminology.
//
//    +------------------+---------------------------------------------+
//    |                  |                 Contents                    |
//    |                  +--------------------+------------------------+
//    |       Set        |      Template      |         Record         |
//    +------------------+--------------------+------------------------+
//    |     Data Set     |          /         |     Data Record(s)     |
//    +------------------+--------------------+------------------------+
//    |   Template Set   | Template Record(s) |           /            |
//    +------------------+--------------------+------------------------+
//    | Options Template |  Options Template  |           /            |
//    |       Set        |      Record(s)     |                        |
//    +------------------+--------------------+------------------------+

// The format of the IPFIX Message is shown in Figure B.
//
//         +----------------------------------------------------+
//         | Message Header                                     |
//         +----------------------------------------------------+
//         | Set                                                |
//         +----------------------------------------------------+
//         | Set                                                |
//         +----------------------------------------------------+
//           ...
//         +----------------------------------------------------+
//         | Set                                                |
//         +----------------------------------------------------+

// An IPFIX Message consisting of interleaved Template, Data, and
//      Options Template Sets, as shown in Figure C.
//
//     +--------+--------------------------------------------------------+
//     |        | +----------+ +---------+     +-----------+ +---------+ |
//     |Message | | Template | | Data    |     | Options   | | Data    | |
//     | Header | | Set      | | Set     | ... | Template  | | Set     | |
//     |        | |          | |         |     | Set       | |         | |
//     |        | +----------+ +---------+     +-----------+ +---------+ |
//     +--------+--------------------------------------------------------+

// An IPFIX Message consisting entirely of Data Sets, sent after the
//      appropriate Template Records have been defined and transmitted to
//      the Collecting Process, as shown in Figure D.
//
//       +--------+----------------------------------------------+
//       |        | +---------+     +---------+      +---------+ |
//       |Message | | Data    |     | Data    |      | Data    | |
//       | Header | | Set     | ... | Set     | ...  | Set     | |
//       |        | +---------+     +---------+      +---------+ |
//       +--------+----------------------------------------------+

// An IPFIX Message consisting entirely of Template and Options
//      Template Sets, as shown in Figure E.  Such a Message can be used
//      to define or redefine Templates and Options Templates in bulk.
//
//      +--------+-------------------------------------------------+
//      |        | +----------+     +----------+      +----------+ |
//      |Message | | Template |     | Template |      | Options  | |
//      | Header | | Set      | ... | Set      | ...  | Template | |
//      |        | |          |     |          |      | Set      | |
//      |        | +----------+     +----------+      +----------+ |
//      +--------+-------------------------------------------------+

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
		return fmt.Sprintf("invalid set #%v", s.SetID)
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

	return fmt.Sprintf("%v (%v octets) %v", setGenre, s.Length, intDescribed(len(s.Records), "record", "records"))
}

func (s *Set) valid() bool {
	return (s.SetID == 2 || s.SetID == 3 || s.SetID >= 256) && s.Length > 3
}

// Unmarshal unmarshal a IPFIX Message.
func (s *Set) Unmarshal(reader io.Reader, templates *map[uint16]*TemplateRecord) (n int, err error) {
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

	return 0, templateNotDefinedError{templateID: s.SetID}
}

func (s *Set) handleUnmarshalError(reader io.Reader, n int, err error) (int, error) {
	switch {
	case errors.Is(err, templateNotDefinedError{templateID: s.SetID}):
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
