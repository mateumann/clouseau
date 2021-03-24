package ipfix

import (
	"encoding/binary"
	"fmt"
	"io"
)

// A Template is an ordered sequence of <type, length> pairs used to
//      completely specify the structure and semantics of a particular set
//      of information that needs to be communicated from an IPFIX Device
//      to a Collector.  Each Template is uniquely identifiable by means
//      of a Template ID.

// A Template Record defines the structure and interpretation of
//      fields in a Data Record.

// An Options Template Record is a Template Record that defines the
//      structure and interpretation of fields in a Data Record, including
//      defining how to scope the applicability of the Data Record.

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
		return fmt.Sprintf("ipfix invalid template record #%v", r.TemplateID)
	}

	return fmt.Sprintf("ipfix template record #%v (%v)", r.TemplateID, fieldDescribed(r.FieldCount))
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
		return fmt.Sprintf("ipfix invalid options template record #%v", r.TemplateID)
	}

	return fmt.Sprintf("ipfix options template record #%v (scope %v of %v)", r.TemplateID,
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
