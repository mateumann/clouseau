package encoding

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

// templateNotDefinedError is returned when no template for input data is found.
type templateNotDefinedError struct {
	templateID uint16
}

func (e templateNotDefinedError) Error() string {
	return fmt.Sprintf("template not defined %v", e.templateID)
}

// UnmarshalError describes an error that occurred when data are unmarshalled.
// It wraps an error returned by the binary.Read function.
type UnmarshalError struct {
	err error
}

func (e UnmarshalError) Error() string {
	return fmt.Sprintf("ipfix: unmarshal error: %v", e.err)
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

// FieldSpecifier describes a single data element in a IPFIX record.
// https://tools.ietf.org/html/rfc7011#section-3.2
type FieldSpecifier struct {
	InformationElementID uint16 `rethinkdb:"element_id"`
	FieldLength          uint16 `rethinkdb:"length"`
	EnterpriseNumber     uint32 `rethinkdb:"enterprise_number"`
}

func (s *FieldSpecifier) String() string {
	if s.valid() {
		if s.e() {
			return fmt.Sprintf("enterprise-specific field #%v (%v)",
				s.InformationElementID, s.EnterpriseNumber)
		}

		return fmt.Sprintf("IANA-registered field #%v",
			InformationElementByID(s.InformationElementID))
	}

	return fmt.Sprintf("invalid field #%v", s.InformationElementID)
}

func (s *FieldSpecifier) e() bool {
	return s.InformationElementID>>15 == 1
}

func (s *FieldSpecifier) valid() bool {
	if s.e() {
		return s.InformationElementID >= 32768 && s.EnterpriseNumber != 0
	}

	return s.InformationElementID < 32768 && s.EnterpriseNumber == 0
}

// Unmarshal reads byte stream from the reader and parses it into a FieldSpecifier.
func (s *FieldSpecifier) Unmarshal(reader io.Reader) (n int, err error) {
	if n, err = unmarshalFields(reader, &s.InformationElementID, &s.FieldLength); err != nil {
		return n, err
	}

	if s.e() {
		bytesRead, err := unmarshalFields(reader, &s.EnterpriseNumber)
		n += bytesRead

		if err != nil {
			return n, err
		}
	}

	return n, err
}

// Message describes a IPFIX Message containing of a header and zero or more sets.
// https://tools.ietf.org/html/rfc7011#section-3.1
type Message struct {
	SourceAddr          string    `rethinkdb:"source"`  // source address of the Message
	Version             uint16    `rethinkdb:"version"` // version of the IPFIX, currently 0x000a
	Length              uint16    `rethinkdb:"length"`  // total length of Message, including the header and Set(s)
	exportTime          uint32    // time at witch Message left the exporter, Unix Epoch, UTC
	ExportedAt          time.Time `rethinkdb:"export_at"` // time at witch Message left the exporter, UTC
	SequenceNumber      uint32    `rethinkdb:"seq"`       // incremental sequence counter modulo 2^32
	ObservationDomainID uint32    `rethinkdb:"domain_id"` // a 32-bit identifier of the observation domain
	Sets                []Set     `rethinkdb:"sets"`
}

func (m *Message) String() string {
	return fmt.Sprintf("Message %v (%v octets, %v, %v, %v) %v",
		m.Version, m.Length, m.ExportedAt.Format(time.RFC3339),
		m.SequenceNumber, m.ObservationDomainID, intDescribed(len(m.Sets), "set", "sets"))
}

func (m Message) valid() bool {
	return m.Version == 0x000a && m.Length > 0
}

// Unmarshal reads byte stream from the reader and parses it into a IPFIX Message. `templates` are used to store IPFIX
// (Options) Template Records used to decode raw data into a IPFIX Data Records.
func (m *Message) Unmarshal(reader io.Reader, source string, templates *map[uint16]*TemplateRecord) (n int, err error) {
	if n, err = unmarshalFields(reader,
		&m.Version, &m.Length, &m.exportTime, &m.SequenceNumber, &m.ObservationDomainID); err != nil {
		m.ExportedAt = time.Unix(int64(m.exportTime), 0).UTC()

		return n, err
	}

	m.ExportedAt = time.Unix(int64(m.exportTime), 0).UTC()
	m.SourceAddr = source

	for n < int(m.Length) {
		bytesRead, err := m.unmarshalSet(reader, templates)
		n += bytesRead

		if err != nil {
			return n, err
		}
	}

	return n, err
}

func (m *Message) unmarshalSet(reader io.Reader, templates *map[uint16]*TemplateRecord) (n int, err error) {
	var s Set
	n, err = s.Unmarshal(reader, templates)

	if err != nil {
		nFromHandler, errFromHandler := s.handleUnmarshalError(reader, n, err)
		if errFromHandler != nil {
			return n + nFromHandler, UnmarshalError{errFromHandler}
		}

		n += nFromHandler
	}

	m.Sets = append(m.Sets, s)

	return n, nil
}

// Record describes one of three ipfix records.
type Record interface {
	Size() int
	Unmarshal(io.Reader, *TemplateRecord) (int, error)
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

// TemplateRecordBase contains fields definitions common for `templateRecord` and `optionsTemplateRecord`.
type TemplateRecordBase struct {
	TemplateID      uint16           `rethinkdb:"template_id"`
	FieldCount      uint16           `rethinkdb:"field_count"` // number of all fields in this record
	FieldSpecifiers []FieldSpecifier `rethinkdb:"field_specifiers"`
}

func (r *TemplateRecordBase) unmarshalFieldSpecifiers(reader io.Reader) (n int, err error) {
	for i := uint16(0); i < r.FieldCount; i++ {
		var fs FieldSpecifier
		bytesRead, err := fs.Unmarshal(reader)
		n += bytesRead

		if err != nil {
			return n, err
		}

		r.FieldSpecifiers = append(r.FieldSpecifiers, fs)
	}

	return n, err
}

// TemplateRecord contains any combination of InformationElement-s.
// https://tools.ietf.org/html/rfc7011#section-3.4.1
type TemplateRecord struct {
	TemplateRecordBase
}

func (r *TemplateRecord) String() string {
	if !r.valid() {
		return fmt.Sprintf("invalid template record #%v", r.TemplateID)
	}

	return fmt.Sprintf("template record #%v (%v)", r.TemplateID, fieldDescribed(r.FieldCount))
}

func (r *TemplateRecord) valid() bool {
	return r.TemplateID >= 256 && r.FieldCount == uint16(len(r.FieldSpecifiers))
}

// Size returns a binary size of the templateRecord.
func (r TemplateRecord) Size() int {
	n := binary.Size(r.TemplateID) + binary.Size(r.FieldCount)

	for _, s := range r.FieldSpecifiers {
		n += binary.Size(s)
	}

	return n
}

// Unmarshal unmarshal a IPFIX templateRecord.
func (r *TemplateRecord) Unmarshal(reader io.Reader, _ *TemplateRecord) (n int, err error) {
	if n, err = unmarshalFields(reader, &r.TemplateID, &r.FieldCount); err != nil {
		return n, err
	}

	bytesRead, err := r.unmarshalFieldSpecifiers(reader)
	n += bytesRead

	return n, err
}

// OptionsTemplateRecord provides additional information about the context of the reported InformationElement-s.
// https://tools.ietf.org/html/rfc7011#section-3.4.2
type OptionsTemplateRecord struct {
	TemplateRecordBase
	ScopeFieldCount uint16 `rethinkdb:"scope_field_count"` // number of scope fields
}

func (r *OptionsTemplateRecord) String() string {
	if !r.valid() {
		return fmt.Sprintf("invalid options template record #%v", r.TemplateID)
	}

	return fmt.Sprintf("options template record #%v (scope %v of %v)", r.TemplateID,
		fieldDescribed(r.ScopeFieldCount), fieldDescribed(r.FieldCount))
}

func (r *OptionsTemplateRecord) valid() bool {
	return r.TemplateID >= 256 &&
		r.FieldCount == uint16(len(r.FieldSpecifiers)) &&
		r.ScopeFieldCount > 0 &&
		r.ScopeFieldCount <= r.FieldCount
}

// Size returns a binary size of the optionsTemplateRecord.
func (r OptionsTemplateRecord) Size() int {
	n := binary.Size(r.TemplateID) + binary.Size(r.FieldCount) + binary.Size(r.ScopeFieldCount)

	for _, s := range r.FieldSpecifiers {
		n += binary.Size(s)
	}

	return n
}

// Unmarshal unmarshal a IPFIX templateRecord.
func (r *OptionsTemplateRecord) Unmarshal(reader io.Reader, _ *TemplateRecord) (n int, err error) {
	if n, err = unmarshalFields(reader, &r.TemplateID, &r.FieldCount, &r.ScopeFieldCount); err != nil {
		return n, err
	}

	bytesRead, err := r.unmarshalFieldSpecifiers(reader)
	n += bytesRead

	return n, err
}

// DataRecord consist only of one or more fieldValues.
// https://tools.ietf.org/html/rfc7011#section-3.4.3
type DataRecord struct {
	FieldValues []interface{} `rethinkdb:"field_values"`
}

func (r *DataRecord) String() string {
	return fmt.Sprintf("data record (%v)", fieldDescribed(uint16(len(r.FieldValues))))
}

// Size returns a binary size of the dataRecord.
func (r DataRecord) Size() int {
	n := 0
	for _, v := range r.FieldValues {
		n += int(reflect.ValueOf(v).Type().Size())
	}

	return n
}

// Unmarshal unmarshal a IPFIX dataRecord.
func (r *DataRecord) Unmarshal(reader io.Reader, template *TemplateRecord) (n int, err error) {
	if template == nil {
		return 0, templateNotDefinedError{0}
	}

	n = 0

	for _, f := range template.FieldSpecifiers {
		e := InformationElementByID(f.InformationElementID)
		buf := make([]byte, f.FieldLength)

		if err = binary.Read(reader, binary.BigEndian, &buf); err != nil {
			return n, &UnmarshalError{err}
		}

		n += int(f.FieldLength)

		elem := createValue(e, buf)

		r.FieldValues = append(r.FieldValues, elem)
	}

	return n, nil
}

// nolint:exhaustive,gomnd
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
	default:
		return data
	}
}

func fillUp(n int, data []byte) []byte {
	if n <= len(data) {
		return data
	}

	buf := make([]byte, n-len(data))

	return append(buf, data...)
}
