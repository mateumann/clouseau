package ipfix

import (
	"fmt"
	"time"
)

// A SetType represents IPFIX set type.  There are three different types: Data set, Template set and
// Options Template set which are determined by SetID field of setHeader:
//   * A value of 2 is reserved for Template Sets.
//   * A value of 3 is reserved for Options Template Sets.
//   * Values from 4 to 255 are reserved for future use.
//   * Values 256 and above are used for Data Sets
//   * The set ID values of 0 and 1 are not used, for historical reasons.
// https://www.ietf.org/rfc/rfc7011.html#section-3.3.2
type SetType int8

const (
	// DataSetType represents Data set type.
	DataSetType SetType = iota
	// TemplateSetType represents Template set type.
	TemplateSetType
	// OptionsTemplateSetType represents Options Template set type.
	OptionsTemplateSetType
	// UnknownSetType represents an unknown (or invalid) set type.
	UnknownSetType SetType = -128
)

func (t SetType) String() string {
	switch t {
	case DataSetType:
		return "Data set"
	case TemplateSetType:
		return "Template set"
	case OptionsTemplateSetType:
		return "Options Template set"
	case UnknownSetType:
		return "Unknown set"
	}

	return "Unknown set"
}

// A Message represents a IPFIX Message.  It consists of a single Message messageHeader and multiple Sets.
// https://www.ietf.org/rfc/rfc7011.html#section-3
type Message struct {
	Header messageHeader
	Sets   []set
}

// A messageHeader represents a IPFIX header.
// see https://tools.ietf.org/html/rfc7011#section-3.1
type messageHeader struct {
	VersionNumber       uint16
	Length              uint16
	ExportTime          uint32
	SequenceNumber      uint32
	ObservationDomainID uint32
}

// A messageHeaderLength defines the length of a IPFIX Message messageHeader.
const messageHeaderLength int64 = 16

// Stringer implementation for IPFIX messageHeader.
func (h messageHeader) String() string {
	timestamp := time.Unix(int64(h.ExportTime), 0).UTC()

	return fmt.Sprintf("IPFIX Message Header (%v) [domain %v; seq %v; exported at %v]",
		h.VersionNumber, h.ObservationDomainID, h.SequenceNumber, timestamp.Format("15:04:05"))
}

// A setHeader represents a IPFIX set messageHeader.
// https://www.ietf.org/rfc/rfc7011.html#section-3.3.2
type setHeader struct {
	SetID  uint16
	Length uint16
}

// A setHeaderLength defines the length of a IPFIX set messageHeader.
const setHeaderLength uint16 = 4

// Type returns SetType based on SetId provided in IPFIX set messageHeader.
//nolint:gomnd
func (h setHeader) Type() SetType {
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

// Stringer implementation for IPFIX setHeader.
func (h setHeader) String() string {
	return fmt.Sprintf("IPFIX %s Header (id=%v, length=%v)", h.Type(), h.SetID, h.Length)
}

// A set represents any of three IPFIX sets: Data set, Template set or Options Template set.
type set interface {
	Type() SetType
	Header() setHeader
	TemplateRecords() []templateRecord
}

// A templateSet type represents a IPFIX Template set.  It consists of a setHeader and number of TemplateRecords.
// https://www.ietf.org/rfc/rfc7011.html#section-3.4.1
type templateSet struct {
	header  setHeader
	records []templateRecord
}

// Type function returns TemplateSetType for any templateSet.
func (s templateSet) Type() SetType {
	return TemplateSetType
}

// Header function returns Template set's IPFIX Header.
func (s templateSet) Header() setHeader {
	return s.header
}

// TemplateRecords function returns Template set's records.
func (s templateSet) TemplateRecords() []templateRecord {
	return s.records
}

// templateRecordHeader represents a IPFIX Template Record messageHeader.
type templateRecordHeader struct {
	TemplateID uint16
	FieldCount uint16
}

// A templateRecordHeaderLength defines the length of a IPFIX Template Record messageHeader.
const templateRecordHeaderLength uint16 = 4

// A templateRecord represents a IPFIX Template Record.
// https://www.ietf.org/rfc/rfc7011.html#section-3.4.1
type templateRecord struct {
	TemplateHeader  templateRecordHeader
	FieldSpecifiers []fieldSpecifier
}

// A dataSet type represents a IPFIX Data set.  It consists of a setHeader and number of DataRecords.
// https://tools.ietf.org/html/rfc7011#section-3.4.3
type dataSet struct {
	header  setHeader
	records []dataRecord
}

// Type function returns DataSetType for any dataSet.
func (s dataSet) Type() SetType {
	return DataSetType
}

// Header function returns Data set's IPFIX setHeader.
func (s dataSet) Header() setHeader {
	return s.header
}

// TemplateRecords function returns nil.
func (dataSet) TemplateRecords() []templateRecord {
	return nil
}

// A dataRecord represents a single portion of data transmitted in IPFIX Data Sets.
// https://www.ietf.org/rfc/rfc7011.html#section-3.4.3
//   * TemplateID defines a templateRecord used to parse data
type dataRecord interface{}
