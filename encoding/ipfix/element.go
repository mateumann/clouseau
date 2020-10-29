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

// // An elementSpecifier identifies an IPFIX element.
// type elementSpecifier struct {
// 	elementID        uint16
// 	enterpriseNumber uint32
// }

// // E returns an Enterprise bit of an elementSpecifier.  This is the first bit of the elementSpecifier.  If this bit
// // is zero, the Information Element identifier identifies an Information Element in [IANA-IPFIX], and the four-octet
// // Enterprise Number field MUST NOT be present.  If this bit is one, the Information Element identifier identifies
// // an enterprise-specific Information Element, and the Enterprise Number field MUST be present.
// func (elementSpec elementSpecifier) E() bool {
// 	return elementSpec.elementID>>15 == 1
// }
//
// func (elementSpec elementSpecifier) String() string {
// 	if elementSpec.E() {
// 		return fmt.Sprintf("Element %d, Enterprise number %d", elementSpec.elementID, elementSpec.enterpriseNumber)
// 	}
//
// 	return fmt.Sprintf("Element %d", elementSpec.elementID)
// }

// fieldSpecifier represents an IPFIX field specification.
// https://www.ietf.org/rfc/rfc7011.html#section-3.2
// https://www.iana.org/assignments/ipfix/ipfix.xhtml
type fieldSpecifier struct {
	elementID        uint16
	length           uint16
	enterpriseNumber uint32
}

// E returns an Enterprise bit of a fieldSpecifier.  This is the first bit of the fieldSpecifier.  If this bit is zero,
// the Information Element identifier identifies an Information Element in [IANA-IPFIX], and the four-octet Enterprise
// Number field MUST NOT be present.  If this bit is one, the Information Element identifier identifies an enterprise-
// specific Information Element, and the Enterprise Number field MUST be present.
func (fieldSpec fieldSpecifier) E() bool {
	return fieldSpec.elementID>>15 == 1
}

// // elementSpecifier returns an elementSpecifier created from fieldSpecifier elementID and enterpriseNumber.
// func (fieldSpec *fieldSpecifier) ElementSpecifier() elementSpecifier {
// 	return elementSpecifier{fieldSpec.elementID, fieldSpec.enterpriseNumber}
// }

// Size returns a size of a fieldSpecifier.  This depends on Enterprise bit being set (or not).
func (fieldSpec fieldSpecifier) Size() uint16 {
	if fieldSpec.E() {
		return fieldSpecifierWEnterpriseSize
	}

	return fieldSpecifierSize
}

func (fieldSpec fieldSpecifier) String() string {
	if fieldSpec.E() {
		return fmt.Sprintf("Element %d, Enterprise number %d", fieldSpec.elementID, fieldSpec.enterpriseNumber)
	}

	return fmt.Sprintf("Element %d", fieldSpec.elementID)
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

func informationElement(elementID uint16, enterpriseNumber uint32) fieldDefinition {
	if enterpriseNumber == 0 {
		if val, ok := ipfixInformationElement()[elementID]; ok {
			return val
		}
	}
	return fieldDefinition{UnknownFieldType, ""}
}

// Kind of immutable map (as there is no such thing as constant map in Go).
func ipfixInformationElement() map[uint16]fieldDefinition {
	return map[uint16]fieldDefinition{
		1:  {Uint64, "octetDeltaCount"},
		2:  {Uint64, "packetDeltaCount"},
		3:  {Uint64, "deltaFlowCount"},
		4:  {Uint8, "protocolIdentifier"},
		5:  {Uint8, "ipClassOfService"},
		6:  {Uint16, "tcpControlBits"},
		7:  {Uint16, "sourceTransportPort"},
		8:  {IPv4Address, "sourceIPv4Address"},
		9:  {Uint8, "sourceIPv4PrefixLength"},
		10: {Uint32, "ingressInterface"},
		11: {Uint16, "destinationTransportPort"},
		12: {IPv4Address, "destinationIPv4Address"},
		13: {Uint8, "destinationIPv4PrefixLength"},
		14: {Uint32, "egressInterface"},
		15: {IPv4Address, "ipNextHopIPv4Address"},
		// many, many more, but let's focus on those actually used by my MikroTik
		21: {Uint32, "flowEndSysUpTime"},
		22: {Uint32, "flowStartSysUpTime"},
		// many, many more, but let's focus on those actually used by my MikroTik
		27: {IPv6Address, "sourceIPv6Address"},
		28: {IPv6Address, "destinationIPv6Address"},
		29: {Uint8, "sourceIPv6PrefixLength"},
		30: {Uint8, "destinationIPv6PrefixLength"},
		31: {Uint32, "flowLabelIPv6"},
		32: {Uint16, "icmpTypeCodeIPv4"},
		33: {Uint8, "igmpType"},
		// many, many more, but let's focus on those actually used by my MikroTik
		56: {MacAddress, "sourceMacAddress"},
		// many, many more, but let's focus on those actually used by my MikroTik
		60: {Uint8, "ipVersion"},
		// many, many more, but let's focus on those actually used by my MikroTik
		62: {IPv6Address, "ipNextHopIPv6Address"},
		// many, many more, but let's focus on those actually used by my MikroTik
		80: {MacAddress, "destinationMacAddress"},
		81: {MacAddress, "postSourceMacAddress"},
		// many, many more, but let's focus on those actually used by my MikroTik
		176: {Uint8, "icmpTypeIPv4"},
		177: {Uint8, "icmpCodeIPv4"},
		178: {Uint8, "icmpTypeIPv6"},
		179: {Uint8, "icmpCodeIPv6"},
		// many, many more, but let's focus on those actually used by my MikroTik
		184: {Uint32, "tcpSequenceNumber"},
		185: {Uint32, "tcpAcknowledgementNumber"},
		186: {Uint16, "tcpWindowSize"},
		// many, many more, but let's focus on those actually used by my MikroTik
		189: {Uint8, "ipHeaderLength"},
		// many, many more, but let's focus on those actually used by my MikroTik
		192: {Uint8, "ipTTL"},
		// many, many more, but let's focus on those actually used by my MikroTik
		205: {Uint16, "udpMessageLength"},
		206: {Uint8, "isMulticast"},
		// many, many more, but let's focus on those actually used by my MikroTik
		224: {Uint64, "ipTotalLength"},
		225: {IPv4Address, "postNATSourceIPv4Address"},
		226: {IPv4Address, "postNATDestinationIPv4Address"},
		227: {Uint16, "postNAPTSourceTransportPort"},
		228: {Uint16, "postNAPTDestinationTransportPort"},
	}
}

// // informationElement retrieves a fieldDefinition from IPFIX fieldModel map.
// func informationElement(elementSpec elementSpecifier) (fieldDefinition, error) {
// 	if fieldDefinition, ok := ipfixInformationElements()[elementSpec]; ok {
// 		return fieldDefinition, nil
// 	}
//
// 	return fieldDefinition{UnknownFieldType, ""},
// 		&elementError{fmt.Sprintf("ipfix: field definition not found for %v", elementSpec)}
// }
//
// // informationElement returns a fieldDefinition basing on provided fieldSpecifier.
// func (fieldSpec fieldSpecifier) InformationElement() (fieldDefinition, error) {
// 	return informationElement(elementSpecifier{fieldSpec.elementID, fieldSpec.enterpriseNumber})
// }
