package ipfix

import (
    "fmt"
    "time"
)

// A SetType represents IPFIX Set type.  There are three different types: Data Set, Template Set and
// Options Template Set which are determined by SetID field of SetHeader:
//   * A value of 2 is reserved for Template Sets.
//   * A value of 3 is reserved for Options Template Sets.
//   * Values from 4 to 255 are reserved for future use.
//   * Values 256 and above are used for Data Sets
//   * The Set ID values of 0 and 1 are not used, for historical reasons.
// https://www.ietf.org/rfc/rfc7011.html#section-3.3.2
type SetType int8

const (
    DataSetType SetType = iota
    TemplateSetType
    OptionsTemplateSetType
    UnknownSetType SetType = -128
)

func (t SetType) String() string {
    switch t {
    case DataSetType:
        return "Data Set"
    case TemplateSetType:
        return "Template Set"
    case OptionsTemplateSetType:
        return "Options Template Set"
    default:
        return "Unknown Set"
    }
}

// A Message represents a IPFIX Message.  It consists of a single Message Header and multiple Sets.
// https://www.ietf.org/rfc/rfc7011.html#section-3
type Message struct {
    Header Header
    Sets   []Set
}

// A Header represents a IPFIX header.
// see https://tools.ietf.org/html/rfc7011#section-3.1
type Header struct {
    VersionNumber       uint16
    Length              uint16
    ExportTime          uint32
    SequenceNumber      uint32
    ObservationDomainID uint32
}

// A headerLength defines the length of a IPFIX Message Header.
const headerLength uint16 = 16

// Stringer implementation for IPFIX Header
func (h Header) String() string {
    timestamp := time.Unix(int64(h.ExportTime), 0)
    return fmt.Sprintf("IPFIX Message Header (%v) [domain %v; seq %v; exported at %v]",
        h.VersionNumber, h.ObservationDomainID, h.SequenceNumber, timestamp.Format("15:04:05"))
}

// A SetHeader represents a IPFIX Set Header
// https://www.ietf.org/rfc/rfc7011.html#section-3.3.2
type SetHeader struct {
    SetID  uint16
    Length uint16
}

// A setHeaderLength defines the length of a IPFIX Set Header.
const setHeaderLength uint16 = 4

// Stringer implementation for IPFIX SetHeader
func (h SetHeader) String() string {
    return fmt.Sprintf("IPFIX Set Header (%s)", h.Type())
}

func (h SetHeader) Type() SetType {
    switch {
    case h.SetID == 2:
        return TemplateSetType
    case h.SetID == 3:
        return OptionsTemplateSetType
    case h.SetID >= 256:
        return DataSetType
    default:
        return UnknownSetType
    }
}

// A Set represents any of three IPFIX sets: Data Set, Template Set or Options Template Set.
type Set interface {
    Type() SetType
}

// A FieldSpecifier represents a IPFIX field specification.
// https://www.ietf.org/rfc/rfc7011.html#section-3.2
// https://www.iana.org/assignments/ipfix/ipfix.xhtml
type FieldSpecifier struct {
    InformationElementID uint16
    FieldLength          uint16
    EnterpriseNumber     uint32
}

// Returns an Enterprise bit.  This is the first bit of the FieldSpecifier.  If this bit is zero, the Information
// Element identifier identifies an Information Element in [IANA-IPFIX], and the four-octet Enterprise Number field
// MUST NOT be present.  If this bit is one, the Information Element identifier identifies an enterprise-specific
// Information Element, and the Enterprise Number field MUST be present.
func (s FieldSpecifier) E() bool {
    return s.InformationElementID>>15 == 1
}

// A fieldSpecifierLength defines the length of a IPFIX Field Specifier.
const fieldSpecifierLength uint16 = 8

// A TemplateSetType represents a IPFIX Template Set.  It consists of a SetHeader and number of TemplateRecords.
// https://www.ietf.org/rfc/rfc7011.html#section-3.4.1
type TemplateSet struct {
    Header SetHeader
    Records []TemplateRecord
}

type TemplateRecordHeader struct {
    TemplateID uint16
    FieldCount uint16
}

// A templateRecordHeaderLength defines the length of a IPFIX Template Record Header.
const templateRecordHeaderLength uint16 = 4

// A TemplateRecord represents a IPFIX Template Record.
// https://www.ietf.org/rfc/rfc7011.html#section-3.4.1
type TemplateRecord struct {
    TemplateID      uint16
    FieldCount      uint16
    FieldSpecifiers []FieldSpecifier
}

func (r TemplateSet) Type() SetType {
    return TemplateSetType
}

// A DataRecord represents a single portion of data transmitted in IPFIX Data Sets.
// https://www.ietf.org/rfc/rfc7011.html#section-3.4.3
//   * TemplateID defines a TemplateRecord used to parse data
type DataRecord struct {
    TemplateID  uint16
    FieldValues []FieldValue
}

type FieldValue interface {
}
