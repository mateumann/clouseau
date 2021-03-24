package ipfix

import (
	"fmt"
	"io"
)

// An Information Element is a protocol- and encoding-independent
//      description of an attribute that may appear in an IPFIX Record.
//      Information Elements are defined in the IANA "IPFIX Information
//      Elements" registry [IANA-IPFIX].  The type associated with an
//      Information Element indicates constraints on what it may contain
//      and also determines the valid encoding mechanisms for use in
//      IPFIX.

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
