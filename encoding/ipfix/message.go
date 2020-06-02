package ipfix

// A Header represents a IPFIX header.
// see https://tools.ietf.org/html/rfc7011#section-3.1
type Header struct {
    VersionNumber       uint16
    Length              uint16
    ExportTime          uint32
    SequenceNumber      uint32
    ObservationDomainID uint32
}

// HeaderLength is a IPFIX header's length
const HeaderLength int = 2 + 2 + 4 + 4 + 4

// A FieldSpecifier
// https://tools.ietf.org/html/rfc7011#section-3.2
type FieldSpecifier struct {
    InformationElementIdentifier uint16 // the very first bit is an Enterprise bit (E)
    FieldLength                  uint16
    EnterpriseNumber             *uint32 // may or may not be empty
}

// A FieldValue represents
// https://tools.ietf.org/html/rfc7011#section-3.4.3
type FieldValue struct {
    Anything uint16 // not really a uint16
}

// A TemplateRecordHeader represents IPFIX Template Record's header.
// https://www.ietf.org/rfc/rfc7011.html#section-3.4.1
type TemplateRecordHeader struct {
    TemplateID uint16 // > 255
    FieldCount uint16
}

// A TemplateRecord represents IPFIX Template Record.
// https://www.ietf.org/rfc/rfc7011.html#section-3.4.1
type TemplateRecord struct {
    Header          TemplateRecordHeader
    FieldSpecifiers []FieldSpecifier
}

// A OptionsTemplateRecordHeader represents IPFIX Options Template Record's header.
// https://www.ietf.org/rfc/rfc7011.html#section-3.4.2
type OptionsTemplateRecordHeader struct {
    TemplateID      uint16 // > 255
    FieldCount      uint16
    ScopeFieldCount uint16
}

// A OptionsTemplateRecord represents IPFIX Options Template Record.
// https://www.ietf.org/rfc/rfc7011.html#section-3.4.2
type OptionsTemplateRecord struct {
    Header          OptionsTemplateRecordHeader
    FieldSpecifiers []FieldSpecifier
}

// A DataRecord represents IPFIX Data Record.
// https://www.ietf.org/rfc/rfc7011.html#section-3.4.3
type DataRecord struct {
    FieldValues []FieldValue
}

// A SetHeader represents a IPFIX Set's header.
// https://www.ietf.org/rfc/rfc7011.html#section-3.3.2
type SetHeader struct {
    SetID  uint16
    Length uint16
}

// A Set represents a IPFIX set.
// It can be any of these three possible types: Data Set, Template Set, or Options Template Set.
type Set interface {
    SetHeader()
}
