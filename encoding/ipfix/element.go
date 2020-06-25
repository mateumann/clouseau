package ipfix

import (
	"fmt"
)

// FieldType is an abstract data types of the IPFIX information model
// https://tools.ietf.org/html/rfc7012#section-3.1
// https://www.iana.org/assignments/ipfix/ipfix.xhtml#ipfix-information-element-data-types
type FieldType uint8

const (
	// OctetArray represents a finite-length string of octets
	OctetArray FieldType = iota
	// Uint8 represents non-negative integer value in the range of 0 to 255
	Uint8
	// Uint16 represents non-negative integer value in the range of 0 to 65535
	Uint16
	// Uint32 represents non-negative integer value in the range of 0 to 4294967295
	Uint32
	// Uint64 represents non-negative integer value in the range of 0 to 18446744073709551615
	Uint64
	// Int8 represents an integer value in the range of -128 to 127
	Int8
	// Int16 represents an integer value in the range of -32768 to 32767
	Int16
	// Int32 represents an integer value in the range of -2147483648 to 2147483647
	Int32
	// Int64 represents an integer value in the range of -9223372036854775808 to 9223372036854775807
	Int64
	// Float32 corresponds to an IEEE single-precision 32-bit floating-point type as defined in IEEE.754.2008
	Float32
	// Float64 corresponds to an IEEE double-precision 64-bit floating-point type as defined in IEEE.754.2008
	Float64
	// Boolean represents a binary value.  The only allowed values are "true" and "false"
	Boolean
	// MacAddress represents a MAC-48 address as defined in IEEE.802-3.2012
	MacAddress
	// String represents a finite-length string of valid characters from the Unicode coded character set ISO.10646
	String
	// DateTimeSeconds represents a time value expressed with second-level precision
	DateTimeSeconds
	// DateTimeMilliseconds represents a time value expressed with millisecond-level precision
	DateTimeMilliseconds
	// DateTimeMicroseconds represents a time value expressed with microsecond-level precision
	DateTimeMicroseconds
	// DateTimeNanoseconds represents a time value expressed with nanosecond-level precision
	DateTimeNanoseconds
	// IPv4Address represents an IPv4 address
	IPv4Address
	// IPv6Address represents an IPv6 address
	IPv6Address
	// BasicList supports structured data export as described in RFC6313
	BasicList
	// SubTemplateList supports structured data export as described in RFC6313
	SubTemplateList
	// SubTemplateMultiList supports structured data export as described in RFC6313
	SubTemplateMultiList
	// UnknownFieldType represents an unknown (undefined) FiledType
	UnknownFieldType
)

// An ElementSpecifier identifies an IPFIX element
type ElementSpecifier struct {
	ElementID        uint16
	EnterpriseNumber uint32
}

// E returns an Enterprise bit of an ElementSpecifier.  This is the first bit of the ElementSpecifier.  If this bit
// is zero, the Information Element identifier identifies an Information Element in [IANA-IPFIX], and the four-octet
// Enterprise Number field MUST NOT be present.  If this bit is one, the Information Element identifier identifies
// an enterprise-specific Information Element, and the Enterprise Number field MUST be present.
func (s ElementSpecifier) E() bool {
	return s.ElementID>>15 == 1
}

func (s ElementSpecifier) String() string {
	if s.E() {
		return fmt.Sprintf("Element %d, Enterprise number %d", s.ElementID, s.EnterpriseNumber)
	}
	return fmt.Sprintf("Element %d", s.ElementID)
}

// FieldSpecifier represents an IPFIX field specification.
// https://www.ietf.org/rfc/rfc7011.html#section-3.2
// https://www.iana.org/assignments/ipfix/ipfix.xhtml
type FieldSpecifier struct {
	ElementID        uint16
	FieldLength      uint16
	EnterpriseNumber uint32
}

// E returns an Enterprise bit of a FieldSpecifier.  This is the first bit of the FieldSpecifier.  If this bit is zero,
// the Information Element identifier identifies an Information Element in [IANA-IPFIX], and the four-octet Enterprise
// Number field MUST NOT be present.  If this bit is one, the Information Element identifier identifies an enterprise-
// specific Information Element, and the Enterprise Number field MUST be present.
func (s FieldSpecifier) E() bool {
	return s.ElementID>>15 == 1
}

// ElementSpecifier returns an ElementSpecifier created from FieldSpecifier ElementID and EnterpriseNumber
func (s *FieldSpecifier) ElementSpecifier() ElementSpecifier {
	return ElementSpecifier{s.ElementID, s.EnterpriseNumber}
}

// Size returns a size of a FieldSpecifier.  This depends on Enterprise bit being set (or not).
func (s FieldSpecifier) Size() uint16 {
	if s.E() {
		return 8
	}
	return 4
}

func (s FieldSpecifier) String() string {
	if s.E() {
		return fmt.Sprintf("Element %d, Enterprise number %d", s.ElementID, s.EnterpriseNumber)
	}
	return fmt.Sprintf("Element %d", s.ElementID)
}

// FieldDefinition describes a IPFIX field (single data)
type FieldDefinition struct {
	Type FieldType
	Name string
}

func (d FieldDefinition) String() string {
	if d.Type == UnknownFieldType {
		return "Unknown Field"
	}
	return fmt.Sprintf("Field '%s'", d.Name)
}

// FieldModel is a mapping of ElementSpecifier to FieldDefinition
type FieldModel map[ElementSpecifier]FieldDefinition

// Kind of immutable map (as there is no such thing as constant map in Go)
func ipfixInformationElements() FieldModel {
	return FieldModel{
		ElementSpecifier{1, 0}:  FieldDefinition{Uint64, "octetDeltaCount"},
		ElementSpecifier{2, 0}:  FieldDefinition{Uint64, "packetDeltaCount"},
		ElementSpecifier{3, 0}:  FieldDefinition{Uint64, "deltaFlowCount"},
		ElementSpecifier{4, 0}:  FieldDefinition{Uint8, "protocolIdentifier"},
		ElementSpecifier{5, 0}:  FieldDefinition{Uint8, "ipClassOfService"},
		ElementSpecifier{6, 0}:  FieldDefinition{Uint16, "tcpControlBits"},
		ElementSpecifier{7, 0}:  FieldDefinition{Uint16, "sourceTransportPort"},
		ElementSpecifier{8, 0}:  FieldDefinition{IPv4Address, "sourceIPv4Address"},
		ElementSpecifier{9, 0}:  FieldDefinition{Uint8, "sourceIPv4PrefixLength"},
		ElementSpecifier{10, 0}: FieldDefinition{Uint32, "ingressInterface"},
		ElementSpecifier{11, 0}: FieldDefinition{Uint16, "destinationTransportPort"},
		ElementSpecifier{12, 0}: FieldDefinition{IPv4Address, "destinationIPv4Address"},
		ElementSpecifier{13, 0}: FieldDefinition{Uint8, "destinationIPv4PrefixLength"},
		ElementSpecifier{14, 0}: FieldDefinition{Uint32, "egressInterface"},
		ElementSpecifier{15, 0}: FieldDefinition{IPv4Address, "ipNextHopIPv4Address"},
		// many, many more, but let's focus on those actually used by my MikroTik
		ElementSpecifier{21, 0}: FieldDefinition{Uint32, "flowEndSysUpTime"},
		ElementSpecifier{22, 0}: FieldDefinition{Uint32, "flowStartSysUpTime"},
		// many, many more, but let's focus on those actually used by my MikroTik
		ElementSpecifier{27, 0}: FieldDefinition{IPv6Address, "sourceIPv6Address"},
		ElementSpecifier{28, 0}: FieldDefinition{IPv6Address, "destinationIPv6Address"},
		ElementSpecifier{29, 0}: FieldDefinition{Uint8, "sourceIPv6PrefixLength"},
		ElementSpecifier{30, 0}: FieldDefinition{Uint8, "destinationIPv6PrefixLength"},
		ElementSpecifier{31, 0}: FieldDefinition{Uint32, "flowLabelIPv6"},
		ElementSpecifier{32, 0}: FieldDefinition{Uint16, "icmpTypeCodeIPv4"},
		ElementSpecifier{33, 0}: FieldDefinition{Uint8, "igmpType"},
		// many, many more, but let's focus on those actually used by my MikroTik
		ElementSpecifier{56, 0}: FieldDefinition{MacAddress, "sourceMacAddress"},
		// many, many more, but let's focus on those actually used by my MikroTik
		ElementSpecifier{60, 0}: FieldDefinition{Uint8, "ipVersion"},
		// many, many more, but let's focus on those actually used by my MikroTik
		ElementSpecifier{62, 0}: FieldDefinition{IPv6Address, "ipNextHopIPv6Address"},
		// many, many more, but let's focus on those actually used by my MikroTik
		ElementSpecifier{80, 0}: FieldDefinition{MacAddress, "destinationMacAddress"},
		ElementSpecifier{81, 0}: FieldDefinition{MacAddress, "postSourceMacAddress"},
		// many, many more, but let's focus on those actually used by my MikroTik
		ElementSpecifier{176, 0}: FieldDefinition{Uint8, "icmpTypeIPv4"},
		ElementSpecifier{177, 0}: FieldDefinition{Uint8, "icmpCodeIPv4"},
		ElementSpecifier{178, 0}: FieldDefinition{Uint8, "icmpTypeIPv6"},
		ElementSpecifier{179, 0}: FieldDefinition{Uint8, "icmpCodeIPv6"},
		// many, many more, but let's focus on those actually used by my MikroTik
		ElementSpecifier{184, 0}: FieldDefinition{Uint32, "tcpSequenceNumber"},
		ElementSpecifier{185, 0}: FieldDefinition{Uint32, "tcpAcknowledgementNumber"},
		ElementSpecifier{186, 0}: FieldDefinition{Uint16, "tcpWindowSize"},
		// many, many more, but let's focus on those actually used by my MikroTik
		ElementSpecifier{189, 0}: FieldDefinition{Uint8, "ipHeaderLength"},
		// many, many more, but let's focus on those actually used by my MikroTik
		ElementSpecifier{192, 0}: FieldDefinition{Uint8, "ipTTL"},
		// many, many more, but let's focus on those actually used by my MikroTik
		ElementSpecifier{205, 0}: FieldDefinition{Uint16, "udpMessageLength"},
		ElementSpecifier{206, 0}: FieldDefinition{Uint8, "isMulticast"},
		// many, many more, but let's focus on those actually used by my MikroTik
		ElementSpecifier{224, 0}: FieldDefinition{Uint64, "ipTotalLength"},
		ElementSpecifier{225, 0}: FieldDefinition{IPv4Address, "postNATSourceIPv4Address"},
		ElementSpecifier{226, 0}: FieldDefinition{IPv4Address, "postNATDestinationIPv4Address"},
		ElementSpecifier{227, 0}: FieldDefinition{Uint16, "postNAPTSourceTransportPort"},
		ElementSpecifier{228, 0}: FieldDefinition{Uint16, "postNAPTDestinationTransportPort"},
	}
}

// InformationElement retrieves a FieldDefinition from IPFIX FieldModel map
func InformationElement(s ElementSpecifier) (FieldDefinition, error) {
	if fieldDefinition, ok := ipfixInformationElements()[s]; ok {
		return fieldDefinition, nil
	}
	return FieldDefinition{UnknownFieldType, ""}, fmt.Errorf("IPFIX information element not found for %v", s)
}

// InformationElement returns a FieldDefinition basing on provided FieldSpecifier
func (s FieldSpecifier) InformationElement() (FieldDefinition, error) {
	return InformationElement(ElementSpecifier{s.ElementID, s.EnterpriseNumber})
}
