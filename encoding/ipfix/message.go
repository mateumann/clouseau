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
	// DataSetType represents Data Set type.
	DataSetType SetType = iota
	// TemplateSetType represents Template Set type.
	TemplateSetType
	// OptionsTemplateSetType represents Options Template Set type.
	OptionsTemplateSetType
	// UnknownSetType represents an unknown (or invalid) Set type.
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
	case UnknownSetType:
		return "Unknown Set"
	}

	return "Unknown Set"
}

// A Message represents a IPFIX Message.  It consists of a single Message MessageHeader and multiple Sets.
// https://www.ietf.org/rfc/rfc7011.html#section-3
type Message struct {
	Header MessageHeader
	Sets   []Set
}

// A MessageHeader represents a IPFIX header.
// see https://tools.ietf.org/html/rfc7011#section-3.1
type MessageHeader struct {
	VersionNumber       uint16
	Length              uint16
	ExportTime          uint32
	SequenceNumber      uint32
	ObservationDomainID uint32
}

// A messageHeaderLength defines the length of a IPFIX Message MessageHeader.
const messageHeaderLength int64 = 16

// Stringer implementation for IPFIX MessageHeader.
func (h MessageHeader) String() string {
	timestamp := time.Unix(int64(h.ExportTime), 0).UTC()

	return fmt.Sprintf("IPFIX Message Header (%v) [domain %v; seq %v; exported at %v]",
		h.VersionNumber, h.ObservationDomainID, h.SequenceNumber, timestamp.Format("15:04:05"))
}

// A SetHeader represents a IPFIX Set MessageHeader.
// https://www.ietf.org/rfc/rfc7011.html#section-3.3.2
type SetHeader struct {
	SetID  uint16
	Length uint16
}

// A setHeaderLength defines the length of a IPFIX Set MessageHeader.
const setHeaderLength uint16 = 4

// Type returns SetType based on SetId provided in IPFIX Set MessageHeader.
//nolint:gomnd
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

// Stringer implementation for IPFIX SetHeader.
func (h SetHeader) String() string {
	return fmt.Sprintf("IPFIX %s Header (id=%v, length=%v)", h.Type(), h.SetID, h.Length)
}

// A Set represents any of three IPFIX sets: Data Set, Template Set or Options Template Set.
type Set interface {
	Type() SetType
	Header() SetHeader
	TemplateRecords() []TemplateRecord
}

// A TemplateSet type represents a IPFIX Template Set.  It consists of a SetHeader and number of TemplateRecords.
// https://www.ietf.org/rfc/rfc7011.html#section-3.4.1
type TemplateSet struct {
	header  SetHeader
	records []TemplateRecord
}

// Type function returns TemplateSetType for any TemplateSet.
func (s TemplateSet) Type() SetType {
	return TemplateSetType
}

// Header function returns Template Set's IPFIX Header.
func (s TemplateSet) Header() SetHeader {
	return s.header
}

// TemplateRecords function returns Template Set's records.
func (s TemplateSet) TemplateRecords() []TemplateRecord {
	return s.records
}

// TemplateRecordHeader represents a IPFIX Template Record MessageHeader.
type TemplateRecordHeader struct {
	TemplateID uint16
	FieldCount uint16
}

// A templateRecordHeaderLength defines the length of a IPFIX Template Record MessageHeader.
const templateRecordHeaderLength uint16 = 4

// A TemplateRecord represents a IPFIX Template Record.
// https://www.ietf.org/rfc/rfc7011.html#section-3.4.1
type TemplateRecord struct {
	TemplateHeader  TemplateRecordHeader
	FieldSpecifiers []FieldSpecifier
}

// A DataSet type represents a IPFIX Data Set.  It consists of a SetHeader and number of DataRecords.
// https://tools.ietf.org/html/rfc7011#section-3.4.3
type DataSet struct {
	header  SetHeader
	records []DataRecord
}

// Type function returns DataSetType for any DataSet.
func (s DataSet) Type() SetType {
	return DataSetType
}

// Header function returns Data Set's IPFIX SetHeader.
func (s DataSet) Header() SetHeader {
	return s.header
}

// TemplateRecords function returns nil.
func (DataSet) TemplateRecords() []TemplateRecord {
	return nil
}

// A DataRecord represents a single portion of data transmitted in IPFIX Data Sets.
// https://www.ietf.org/rfc/rfc7011.html#section-3.4.3
//   * TemplateID defines a TemplateRecord used to parse data
type DataRecord interface{}
