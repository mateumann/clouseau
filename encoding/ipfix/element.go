package ipfix

import (
	"fmt"
)

const (
	fieldSpecifierSize            = 4
	fieldSpecifierWEnterpriseSize = 8
)

// elementError is an error definition associated with IPFIX Elements.
type elementError struct {
	msg string
}

// Associate Error() function to the elementError.
func (e *elementError) Error() string {
	return e.msg
}

// fieldType is an abstract data types of the IPFIX information model
// https://tools.ietf.org/html/rfc7012#section-3.1
// https://www.iana.org/assignments/ipfix/ipfix.xhtml#ipfix-information-element-data-types
type fieldType uint8

const (
	// OctetArray represents a finite-length string of octets.
	OctetArray fieldType = iota
	// Uint8 represents non-negative integer value in the range of 0 to 255.
	Uint8
	// Uint16 represents non-negative integer value in the range of 0 to 65535.
	Uint16
	// Uint32 represents non-negative integer value in the range of 0 to 4294967295.
	Uint32
	// Uint64 represents non-negative integer value in the range of 0 to 18446744073709551615.
	Uint64
	// Int8 represents an integer value in the range of -128 to 127.
	Int8
	// Int16 represents an integer value in the range of -32768 to 32767.
	Int16
	// Int32 represents an integer value in the range of -2147483648 to 2147483647.
	Int32
	// Int64 represents an integer value in the range of -9223372036854775808 to 9223372036854775807.
	Int64
	// Float32 corresponds to an IEEE single-precision 32-bit floating-point type as defined in IEEE.754.2008.
	Float32
	// Float64 corresponds to an IEEE double-precision 64-bit floating-point type as defined in IEEE.754.2008.
	Float64
	// Boolean represents a binary value.  The only allowed values are "true" and "false".
	Boolean
	// MacAddress represents a MAC-48 address as defined in IEEE.802-3.2012.
	MacAddress
	// String represents a finite-length string of valid characters from the Unicode coded character set ISO.10646.
	String
	// DateTimeSeconds represents a time value expressed with second-level precision.
	DateTimeSeconds
	// DateTimeMilliseconds represents a time value expressed with millisecond-level precision.
	DateTimeMilliseconds
	// DateTimeMicroseconds represents a time value expressed with microsecond-level precision.
	DateTimeMicroseconds
	// DateTimeNanoseconds represents a time value expressed with nanosecond-level precision.
	DateTimeNanoseconds
	// IPv4Address represents an IPv4 address.
	IPv4Address
	// IPv6Address represents an IPv6 address.
	IPv6Address
	// BasicList supports structured data export as described in RFC6313.
	BasicList
	// SubTemplateList supports structured data export as described in RFC6313.
	SubTemplateList
	// SubTemplateMultiList supports structured data export as described in RFC6313.
	SubTemplateMultiList
	// UnknownFieldType represents an unknown (undefined) FiledType.
	UnknownFieldType
)

// An elementSpecifier identifies an IPFIX element.
type elementSpecifier struct {
	ElementID        uint16
	EnterpriseNumber uint32
}

// E returns an Enterprise bit of an elementSpecifier.  This is the first bit of the elementSpecifier.  If this bit
// is zero, the Information Element identifier identifies an Information Element in [IANA-IPFIX], and the four-octet
// Enterprise Number field MUST NOT be present.  If this bit is one, the Information Element identifier identifies
// an enterprise-specific Information Element, and the Enterprise Number field MUST be present.
func (elementSpec elementSpecifier) E() bool {
	return elementSpec.ElementID>>15 == 1
}

func (elementSpec elementSpecifier) String() string {
	if elementSpec.E() {
		return fmt.Sprintf("Element %d, Enterprise number %d", elementSpec.ElementID, elementSpec.EnterpriseNumber)
	}

	return fmt.Sprintf("Element %d", elementSpec.ElementID)
}

// fieldSpecifier represents an IPFIX field specification.
// https://www.ietf.org/rfc/rfc7011.html#section-3.2
// https://www.iana.org/assignments/ipfix/ipfix.xhtml
type fieldSpecifier struct {
	ElementID        uint16
	FieldLength      uint16
	EnterpriseNumber uint32
}

// E returns an Enterprise bit of a fieldSpecifier.  This is the first bit of the fieldSpecifier.  If this bit is zero,
// the Information Element identifier identifies an Information Element in [IANA-IPFIX], and the four-octet Enterprise
// Number field MUST NOT be present.  If this bit is one, the Information Element identifier identifies an enterprise-
// specific Information Element, and the Enterprise Number field MUST be present.
func (fieldSpec fieldSpecifier) E() bool {
	return fieldSpec.ElementID>>15 == 1
}

// elementSpecifier returns an elementSpecifier created from fieldSpecifier ElementID and EnterpriseNumber.
func (fieldSpec *fieldSpecifier) ElementSpecifier() elementSpecifier {
	return elementSpecifier{fieldSpec.ElementID, fieldSpec.EnterpriseNumber}
}

// Size returns a size of a fieldSpecifier.  This depends on Enterprise bit being set (or not).
func (fieldSpec fieldSpecifier) Size() uint16 {
	if fieldSpec.E() {
		return fieldSpecifierWEnterpriseSize
	}

	return fieldSpecifierSize
}

func (fieldSpec fieldSpecifier) String() string {
	return fieldSpec.ElementSpecifier().String()
}

// fieldDefinition describes a IPFIX field (single data).
type fieldDefinition struct {
	Type fieldType
	Name string
}

func (fieldDef fieldDefinition) String() string {
	if fieldDef.Type == UnknownFieldType {
		return "Unknown Field"
	}

	return fmt.Sprintf("Field '%s'", fieldDef.Name)
}

// fieldModel is a mapping of elementSpecifier to fieldDefinition.
type fieldModel map[elementSpecifier]fieldDefinition

// Kind of immutable map (as there is no such thing as constant map in Go).
func ipfixInformationElements() fieldModel {
	return fieldModel{
		elementSpecifier{1, 0}:  fieldDefinition{Uint64, "octetDeltaCount"},
		elementSpecifier{2, 0}:  fieldDefinition{Uint64, "packetDeltaCount"},
		elementSpecifier{3, 0}:  fieldDefinition{Uint64, "deltaFlowCount"},
		elementSpecifier{4, 0}:  fieldDefinition{Uint8, "protocolIdentifier"},
		elementSpecifier{5, 0}:  fieldDefinition{Uint8, "ipClassOfService"},
		elementSpecifier{6, 0}:  fieldDefinition{Uint16, "tcpControlBits"},
		elementSpecifier{7, 0}:  fieldDefinition{Uint16, "sourceTransportPort"},
		elementSpecifier{8, 0}:  fieldDefinition{IPv4Address, "sourceIPv4Address"},
		elementSpecifier{9, 0}:  fieldDefinition{Uint8, "sourceIPv4PrefixLength"},
		elementSpecifier{10, 0}: fieldDefinition{Uint32, "ingressInterface"},
		elementSpecifier{11, 0}: fieldDefinition{Uint16, "destinationTransportPort"},
		elementSpecifier{12, 0}: fieldDefinition{IPv4Address, "destinationIPv4Address"},
		elementSpecifier{13, 0}: fieldDefinition{Uint8, "destinationIPv4PrefixLength"},
		elementSpecifier{14, 0}: fieldDefinition{Uint32, "egressInterface"},
		elementSpecifier{15, 0}: fieldDefinition{IPv4Address, "ipNextHopIPv4Address"},
		// many, many more, but let's focus on those actually used by my MikroTik
		elementSpecifier{21, 0}: fieldDefinition{Uint32, "flowEndSysUpTime"},
		elementSpecifier{22, 0}: fieldDefinition{Uint32, "flowStartSysUpTime"},
		// many, many more, but let's focus on those actually used by my MikroTik
		elementSpecifier{27, 0}: fieldDefinition{IPv6Address, "sourceIPv6Address"},
		elementSpecifier{28, 0}: fieldDefinition{IPv6Address, "destinationIPv6Address"},
		elementSpecifier{29, 0}: fieldDefinition{Uint8, "sourceIPv6PrefixLength"},
		elementSpecifier{30, 0}: fieldDefinition{Uint8, "destinationIPv6PrefixLength"},
		elementSpecifier{31, 0}: fieldDefinition{Uint32, "flowLabelIPv6"},
		elementSpecifier{32, 0}: fieldDefinition{Uint16, "icmpTypeCodeIPv4"},
		elementSpecifier{33, 0}: fieldDefinition{Uint8, "igmpType"},
		// many, many more, but let's focus on those actually used by my MikroTik
		elementSpecifier{56, 0}: fieldDefinition{MacAddress, "sourceMacAddress"},
		// many, many more, but let's focus on those actually used by my MikroTik
		elementSpecifier{60, 0}: fieldDefinition{Uint8, "ipVersion"},
		// many, many more, but let's focus on those actually used by my MikroTik
		elementSpecifier{62, 0}: fieldDefinition{IPv6Address, "ipNextHopIPv6Address"},
		// many, many more, but let's focus on those actually used by my MikroTik
		elementSpecifier{80, 0}: fieldDefinition{MacAddress, "destinationMacAddress"},
		elementSpecifier{81, 0}: fieldDefinition{MacAddress, "postSourceMacAddress"},
		// many, many more, but let's focus on those actually used by my MikroTik
		elementSpecifier{176, 0}: fieldDefinition{Uint8, "icmpTypeIPv4"},
		elementSpecifier{177, 0}: fieldDefinition{Uint8, "icmpCodeIPv4"},
		elementSpecifier{178, 0}: fieldDefinition{Uint8, "icmpTypeIPv6"},
		elementSpecifier{179, 0}: fieldDefinition{Uint8, "icmpCodeIPv6"},
		// many, many more, but let's focus on those actually used by my MikroTik
		elementSpecifier{184, 0}: fieldDefinition{Uint32, "tcpSequenceNumber"},
		elementSpecifier{185, 0}: fieldDefinition{Uint32, "tcpAcknowledgementNumber"},
		elementSpecifier{186, 0}: fieldDefinition{Uint16, "tcpWindowSize"},
		// many, many more, but let's focus on those actually used by my MikroTik
		elementSpecifier{189, 0}: fieldDefinition{Uint8, "ipHeaderLength"},
		// many, many more, but let's focus on those actually used by my MikroTik
		elementSpecifier{192, 0}: fieldDefinition{Uint8, "ipTTL"},
		// many, many more, but let's focus on those actually used by my MikroTik
		elementSpecifier{205, 0}: fieldDefinition{Uint16, "udpMessageLength"},
		elementSpecifier{206, 0}: fieldDefinition{Uint8, "isMulticast"},
		// many, many more, but let's focus on those actually used by my MikroTik
		elementSpecifier{224, 0}: fieldDefinition{Uint64, "ipTotalLength"},
		elementSpecifier{225, 0}: fieldDefinition{IPv4Address, "postNATSourceIPv4Address"},
		elementSpecifier{226, 0}: fieldDefinition{IPv4Address, "postNATDestinationIPv4Address"},
		elementSpecifier{227, 0}: fieldDefinition{Uint16, "postNAPTSourceTransportPort"},
		elementSpecifier{228, 0}: fieldDefinition{Uint16, "postNAPTDestinationTransportPort"},
	}
}

// informationElement retrieves a fieldDefinition from IPFIX fieldModel map.
func informationElement(elementSpec elementSpecifier) (fieldDefinition, error) {
	if fieldDefinition, ok := ipfixInformationElements()[elementSpec]; ok {
		return fieldDefinition, nil
	}

	return fieldDefinition{UnknownFieldType, ""},
		&elementError{fmt.Sprintf("ipfix: field definition not found for %v", elementSpec)}
}

// informationElement returns a fieldDefinition basing on provided fieldSpecifier.
func (fieldSpec fieldSpecifier) InformationElement() (fieldDefinition, error) {
	return informationElement(elementSpecifier{fieldSpec.ElementID, fieldSpec.EnterpriseNumber})
}
