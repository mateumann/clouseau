package ipfix

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"reflect"
	"time"
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

var errInvalidSetType = errors.New("invalid set type")

var errTemplateNotDefined = errors.New("template not defined")

// var errUnmarshalError = errors.New("unmarshal error")

// type errUnhandledKind struct {
// 	kind reflect.Kind
// }
//
// func (e errUnhandledKind) Error() string {
// 	return fmt.Sprintf("ipfix: unhandled kind %s", e.kind)
// }

// type errIncompatibleData struct {
// 	informationElement InformationElement
// 	data               []byte
// }
//
// func (e errIncompatibleData) Error() string {
// 	return fmt.Sprintf("ipfix: data %v incompatible with Information Element %v", e.data, e.informationElement)
// }

// UnmarshalError describes an error while unmarshalling data.
// It wraps an error returned by the binary.Read function.
type UnmarshalError struct {
	err error
}

func (e UnmarshalError) Error() string {
	return fmt.Sprintf("ipfix: unmarshalling error: %v", e.err)
}

// Unwrap returns the underlying error.
func (e UnmarshalError) Unwrap() error {
	return e.err
}

// unmarshalFields is a helper function that reads BigEndian binary stream into a bunch of fields.
// In case of an error returns UnmarshalError-wrapped error returned by the binary.Read.
func unmarshalFields(reader io.Reader, data ...interface{}) (n int, err error) {
	n = 0

	for _, elem := range data {
		if err = binary.Read(reader, binary.BigEndian, elem); err != nil {
			return n, &UnmarshalError{err}
		}

		elemSize := binary.Size(elem)

		// This will never happen because there is no such type that is "readable" but its size is not determined
		// if elemSize < 0 {
		// 	return n, &UnmarshalError{NotReadableError{reflect.TypeOf(elem)}}
		// }

		n += elemSize
	}

	return n, nil
}

// fieldSpecifier describes a single data element in a IPFIX record.
// https://tools.ietf.org/html/rfc7011#section-3.2
type fieldSpecifier struct {
	informationElementID uint16
	fieldLength          uint16
	enterpriseNumber     uint32
}

func (s *fieldSpecifier) String() string {
	if s.valid() {
		if s.e() {
			return fmt.Sprintf("enterprise-specific field #%v (%v)",
				s.informationElementID, s.enterpriseNumber)
		}

		return fmt.Sprintf("IANA-registered field #%v",
			InformationElementByID(s.informationElementID))
	}

	return fmt.Sprintf("invalid field #%v", s.informationElementID)
}

func (s *fieldSpecifier) e() bool {
	return s.informationElementID>>15 == 1
}

func (s *fieldSpecifier) valid() bool {
	if s.e() {
		return s.informationElementID >= 32768 && s.enterpriseNumber != 0
	}

	return s.informationElementID < 32768 && s.enterpriseNumber == 0
}

func (s *fieldSpecifier) Unmarshal(reader io.Reader) (n int, err error) {
	if n, err = unmarshalFields(reader, &s.informationElementID, &s.fieldLength); err != nil {
		return n, err
	}

	if s.e() {
		bytesRead, err := unmarshalFields(reader, &s.enterpriseNumber)
		n += bytesRead

		if err != nil {
			return n, err
		}
	}

	return n, err
}

// Message describes a IPFIX Message constaing of a header and zero or more sets.
// https://tools.ietf.org/html/rfc7011#section-3.1
type Message struct {
	version             uint16 // version of the IPFIX, currently 0x000a
	length              uint16 // total length of IPFIX Message, measured in octets, including the header and Set(s)
	exportTime          uint32 // time at witch the Message left the exporter, in seconds since the UNIX epoch, UTC
	sequenceNumber      uint32 // incremental sequence counter modulo 2^32
	observationDomainID uint32 // a 32-bit identifier of the observation domain, locally unique to the exporting proc.
	sets                []set
}

func (m *Message) String() string {
	return fmt.Sprintf("Message %v (%v octets, %v, %v, %v) %v",
		m.version, m.length, time.Unix(int64(m.exportTime), 0).UTC().Format(time.RFC3339),
		m.sequenceNumber, m.observationDomainID, intDescribed(len(m.sets), "set", "sets"))
}

// Unmarshal unmarshal a IPFIX Message.
func (m *Message) Unmarshal(reader io.Reader, templates []*templateRecord) (n int, err error) {
	if n, err = unmarshalFields(reader,
		&m.version, &m.length, &m.exportTime, &m.sequenceNumber, &m.observationDomainID); err != nil {
		return n, err
	}

	for n < int(m.length) {
		var s set
		bytesRead, err := s.Unmarshal(reader, templates)
		n += bytesRead

		if err != nil && errors.Is(errors.Unwrap(err), io.ErrUnexpectedEOF) && (int(s.length) > bytesRead) {
			// assume s.Unmarshal failed to read a record which was actually a padding
			buf := make([]byte, int(s.length)-bytesRead)
			if err := binary.Read(reader, binary.BigEndian, &buf); err != nil {
				return n + len(buf), UnmarshalError{err}
			}
		}

		if err != nil {
			return n, err
		}

		m.sets = append(m.sets, s)
	}

	return n, err
}

// record describes one of three ipfix records.
type record interface {
	Size() int
	Unmarshal(io.Reader, *templateRecord) (int, error)
}

// set describes a generic term for a collection of records that have a similar structure.  There are three different
// types of set: template sets, options template sets, and data sets.  Each of these sets consists of a set header and
// one or more records.  The exporting process may insert some padding octets, so that the subsequent set starts at an
// aligned boundary.
// https://tools.ietf.org/html/rfc7011#section-3.3.1, https://tools.ietf.org/html/rfc7011#section-3.3.2
type set struct {
	setID   uint16
	length  uint16 // total length of the set, including set header, all records, and the optional padding; in octets
	records []record
}

// Set ID constants defining valid IPFIX set identifiers (or a minimum valid set id in case of MinimumDataSetID).
const (
	TemplateSetID        = 2   // TemplateSetID identifies a IPFIX Template Set
	OptionsTemplateSetID = 3   // OptionsTemplateSetID identifies a IPFIX Options Template Set
	MinimumDataSetID     = 256 // set ids >= MinimumDataSetID identify IPFIX Data Sets
)

func (s *set) String() string {
	if !s.valid() {
		return fmt.Sprintf("invalid set #%v", s.setID)
	}

	var setGenre string

	switch {
	case s.setID == TemplateSetID:
		setGenre = "template set"
	case s.setID == OptionsTemplateSetID:
		setGenre = "options template set"
	case s.setID >= MinimumDataSetID:
		setGenre = "data set"
	}

	return fmt.Sprintf("%v (%v octets) %v", setGenre, s.length, intDescribed(len(s.records), "record", "records"))
}

func (s *set) valid() bool {
	return (s.setID == 2 || s.setID == 3 || s.setID >= 256) && s.length > 3
}

// Unmarshal unmarshal a IPFIX Message.
func (s *set) Unmarshal(reader io.Reader, templates []*templateRecord) (n int, err error) {
	if n, err = unmarshalFields(reader, &s.setID, &s.length); err != nil {
		return n, err
	}

	var bytesRead int

	for n < int(s.length) {
		switch {
		case s.setID == TemplateSetID:
			bytesRead, err = s.unmarshalTemplateRecord(reader)

			if err == nil {
				templates = append(templates, s.records[len(s.records)-1].(*templateRecord))
			}

		case s.setID == OptionsTemplateSetID:
			bytesRead, err = s.unmarshalOptionsTemplateRecord(reader)

		case s.setID >= MinimumDataSetID:
			if template, ok := getTemplate(s.setID, templates); ok {
				bytesRead, err = s.unmarshalDataRecord(reader, template)
			} else {
				return n, UnmarshalError{errTemplateNotDefined}
			}

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

func (s *set) unmarshalTemplateRecord(reader io.Reader) (n int, err error) {
	r := &templateRecord{}
	n, err = r.Unmarshal(reader, nil)

	if err != nil {
		return n, err
	}

	s.records = append(s.records, r)

	return n, nil
}

func (s *set) unmarshalOptionsTemplateRecord(reader io.Reader) (n int, err error) {
	r := &optionsTemplateRecord{}
	n, err = r.Unmarshal(reader, nil)

	if err != nil {
		return n, err
	}

	s.records = append(s.records, r)

	return n, nil
}

func getTemplate(templateID uint16, templates []*templateRecord) (*templateRecord, bool) {
	for _, template := range templates {
		if template.templateID == templateID {
			return template, true
		}
	}

	return nil, false
}

func (s *set) unmarshalDataRecord(reader io.Reader, template *templateRecord) (n int, err error) {
	r := &dataRecord{}
	n, err = r.Unmarshal(reader, template)

	if err != nil {
		return n, err
	}

	s.records = append(s.records, r)

	return n, nil
}

// templateRecord contains any combination of InformationElement-s.
// https://tools.ietf.org/html/rfc7011#section-3.4.1
type templateRecord struct {
	templateID      uint16
	fieldCount      uint16
	fieldSpecifiers []fieldSpecifier
}

func (r *templateRecord) String() string {
	if !r.valid() {
		return fmt.Sprintf("invalid template record #%v", r.templateID)
	}

	return fmt.Sprintf("template record #%v (%v)", r.templateID, fieldDescribed(r.fieldCount))
}

func (r *templateRecord) valid() bool {
	return r.templateID >= 256 && r.fieldCount == uint16(len(r.fieldSpecifiers))
}

// Size returns a binary size of the templateRecord.
func (r templateRecord) Size() int {
	n := binary.Size(r.templateID) + binary.Size(r.fieldCount)

	for _, s := range r.fieldSpecifiers {
		n += binary.Size(s)
	}

	return n
}

// Unmarshal unmarshal a IPFIX templateRecord.
func (r *templateRecord) Unmarshal(reader io.Reader, _ *templateRecord) (n int, err error) {
	if n, err = unmarshalFields(reader, &r.templateID, &r.fieldCount); err != nil {
		return n, err
	}

	for i := uint16(0); i < r.fieldCount; i++ {
		var fs fieldSpecifier
		bytesRead, err := fs.Unmarshal(reader)
		n += bytesRead

		if err != nil {
			return n, err
		}

		r.fieldSpecifiers = append(r.fieldSpecifiers, fs)
	}

	return n, err
}

// optionsTemplateRecord provides additional information about the context of the reported InformationElement-s.
// https://tools.ietf.org/html/rfc7011#section-3.4.2
type optionsTemplateRecord struct {
	templateID      uint16
	fieldCount      uint16 // number of all fields in this record, including the scope fields
	scopeFieldCount uint16 // number of scope fields (those are the first N fields of the record)
	fieldSpecifiers []fieldSpecifier
}

func (r *optionsTemplateRecord) String() string {
	if !r.valid() {
		return fmt.Sprintf("invalid options template record #%v", r.templateID)
	}

	return fmt.Sprintf("options template record #%v (scope %v of %v)", r.templateID,
		fieldDescribed(r.scopeFieldCount), fieldDescribed(r.fieldCount))
}

func (r *optionsTemplateRecord) valid() bool {
	return r.templateID >= 256 &&
		r.fieldCount == uint16(len(r.fieldSpecifiers)) &&
		r.scopeFieldCount > 0 &&
		r.scopeFieldCount <= r.fieldCount
}

// Size returns a binary size of the optionsTemplateRecord.
func (r optionsTemplateRecord) Size() int {
	n := binary.Size(r.templateID) + binary.Size(r.fieldCount) + binary.Size(r.scopeFieldCount)

	for _, s := range r.fieldSpecifiers {
		n += binary.Size(s)
	}

	return n
}

// Unmarshal unmarshal a IPFIX templateRecord.
func (r *optionsTemplateRecord) Unmarshal(reader io.Reader, _ *templateRecord) (n int, err error) {
	if n, err = unmarshalFields(reader, &r.templateID, &r.fieldCount, &r.scopeFieldCount); err != nil {
		return n, err
	}

	for i := uint16(0); i < r.fieldCount; i++ {
		var fs fieldSpecifier
		bytesRead, err := fs.Unmarshal(reader)
		n += bytesRead

		if err != nil {
			return n, err
		}

		r.fieldSpecifiers = append(r.fieldSpecifiers, fs)
	}

	return n, err
}

// dataRecord consist only of one or more fieldValues.
// https://tools.ietf.org/html/rfc7011#section-3.4.3
type dataRecord struct {
	fieldValues []interface{}
}

func (r *dataRecord) String() string {
	return fmt.Sprintf("data record (%v)", fieldDescribed(uint16(len(r.fieldValues))))
}

// Size returns a binary size of the dataRecord.
func (r dataRecord) Size() int {
	n := 0
	for _, v := range r.fieldValues {
		n += int(reflect.ValueOf(v).Type().Size())
	}

	return n
}

// Unmarshal unmarshal a IPFIX dataRecord.
func (r *dataRecord) Unmarshal(reader io.Reader, template *templateRecord) (n int, err error) {
	if template == nil {
		return 0, errTemplateNotDefined
	}

	n = 0

	for _, f := range template.fieldSpecifiers {
		e := InformationElementByID(f.informationElementID)
		buf := make([]byte, f.fieldLength)

		if err = binary.Read(reader, binary.BigEndian, &buf); err != nil {
			return n, &UnmarshalError{err}
		}

		n += int(f.fieldLength)

		elem := createValue(e, buf)

		r.fieldValues = append(r.fieldValues, elem)
	}

	return n, nil
}

// nolint: exhaustive, gomnd
func createValue(e InformationElement, data []byte) interface{} {
	switch e.dataType {
	case reflect.Bool:
		return data[0] != 0
	case reflect.Int32:
		return int32(binary.BigEndian.Uint32(fillUp(4, data)))
	case reflect.Uint8:
		return data[0]
	case reflect.Uint16:
		return binary.BigEndian.Uint16(fillUp(2, data))
	case reflect.Uint32:
		return binary.BigEndian.Uint32(fillUp(4, data))
	case reflect.Uint64:
		return binary.BigEndian.Uint64(fillUp(8, data))
	case reflect.Float64:
		return math.Float64frombits(binary.BigEndian.Uint64(fillUp(8, data)))
	case reflect.String:
		return string(data)
	}

	return data
}

func fillUp(n int, data []byte) []byte {
	if n <= len(data) {
		return data
	}

	buf := make([]byte, n-len(data))

	return append(buf, data...)
}
