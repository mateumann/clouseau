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
    EnterpriseNumber             uint32
}

// A SetHeader represents a IPFIX Set's header.
// https://www.ietf.org/rfc/rfc7011.html#section-3.3.2
type SetHeader struct {
    SetID  uint16
    Length uint16
}

// A TemplateRecordHeader represents IPFIX Template Record's header.
// https://www.ietf.org/rfc/rfc7011.html#section-3.4.1
type TemplateRecordHeader struct {
    TemplateID uint16
    FieldCount uint16
}

// A


// A Set represents a IPFIX set.
// It can be any of these three possible types: Data Set, Template Set, or Options Template Set.
type Set struct {

}
