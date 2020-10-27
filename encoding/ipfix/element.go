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

// func (fieldSpec fieldSpecifier) String() string {
// 	return fieldSpec.ElementSpecifier().String()
// }

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

// Kind of immutable map (as there is no such thing as constant map in Go).
func ipfixInformationElements(elementID uint16, enterpriseNumber uint32) fieldDefinition {
	if enterpriseNumber == 0 {
		switch elementID {
		case 1:
			return fieldDefinition{Uint64, "octetDeltaCount"}
		case 2:
			return fieldDefinition{Uint64, "packetDeltaCount"}
		case 3:
			return fieldDefinition{Uint64, "deltaFlowCount"}
		case 4:
			return fieldDefinition{Uint8, "protocolIdentifier"}
		case 5:
			return fieldDefinition{Uint8, "ipClassOfService"}
		case 6:
			return fieldDefinition{Uint16, "tcpControlBits"}
		case 7:
			return fieldDefinition{Uint16, "sourceTransportPort"}
		case 8:
			return fieldDefinition{IPv4Address, "sourceIPv4Address"}
		case 9:
			return fieldDefinition{Uint8, "sourceIPv4PrefixLength"}
		case 10:
			return fieldDefinition{Uint32, "ingressInterface"}
		case 11:
			return fieldDefinition{Uint16, "destinationTransportPort"}
		case 12:
			return fieldDefinition{IPv4Address, "destinationIPv4Address"}
		case 13:
			return fieldDefinition{Uint8, "destinationIPv4PrefixLength"}
		case 14:
			return fieldDefinition{Uint32, "egressInterface"}
		case 15:
			return fieldDefinition{IPv4Address, "ipNextHopIPv4Address"}
		// many, many more, but let's focus on those actually used by my MikroTik
		case 21:
			return fieldDefinition{Uint32, "flowEndSysUpTime"}
		case 22:
			return fieldDefinition{Uint32, "flowStartSysUpTime"}
		// many, many more, but let's focus on those actually used by my MikroTik
		case 27:
			return fieldDefinition{IPv6Address, "sourceIPv6Address"}
		case 28:
			return fieldDefinition{IPv6Address, "destinationIPv6Address"}
		case 29:
			return fieldDefinition{Uint8, "sourceIPv6PrefixLength"}
		case 30:
			return fieldDefinition{Uint8, "destinationIPv6PrefixLength"}
		case 31:
			return fieldDefinition{Uint32, "flowLabelIPv6"}
		case 32:
			return fieldDefinition{Uint16, "icmpTypeCodeIPv4"}
		case 33:
			return fieldDefinition{Uint8, "igmpType"}
		// many, many more, but let's focus on those actually used by my MikroTik
		case 56:
			return fieldDefinition{MacAddress, "sourceMacAddress"}
		// many, many more, but let's focus on those actually used by my MikroTik
		case 60:
			return fieldDefinition{Uint8, "ipVersion"}
		// many, many more, but let's focus on those actually used by my MikroTik
		case 62:
			return fieldDefinition{IPv6Address, "ipNextHopIPv6Address"}
		// many, many more, but let's focus on those actually used by my MikroTik
		case 80:
			return fieldDefinition{MacAddress, "destinationMacAddress"}
		case 81:
			return fieldDefinition{MacAddress, "postSourceMacAddress"}
		// many, many more, but let's focus on those actually used by my MikroTik
		case 176:
			return fieldDefinition{Uint8, "icmpTypeIPv4"}
		case 177:
			return fieldDefinition{Uint8, "icmpCodeIPv4"}
		case 178:
			return fieldDefinition{Uint8, "icmpTypeIPv6"}
		case 179:
			return fieldDefinition{Uint8, "icmpCodeIPv6"}
		// many, many more, but let's focus on those actually used by my MikroTik
		case 184:
			return fieldDefinition{Uint32, "tcpSequenceNumber"}
		case 185:
			return fieldDefinition{Uint32, "tcpAcknowledgementNumber"}
		case 186:
			return fieldDefinition{Uint16, "tcpWindowSize"}
		// many, many more, but let's focus on those actually used by my MikroTik
		case 189:
			return fieldDefinition{Uint8, "ipHeaderLength"}
		// many, many more, but let's focus on those actually used by my MikroTik
		case 192:
			return fieldDefinition{Uint8, "ipTTL"}
		// many, many more, but let's focus on those actually used by my MikroTik
		case 205:
			return fieldDefinition{Uint16, "udpMessageLength"}
		case 206:
			return fieldDefinition{Uint8, "isMulticast"}
		// many, many more, but let's focus on those actually used by my MikroTik
		case 224:
			return fieldDefinition{Uint64, "ipTotalLength"}
		case 225:
			return fieldDefinition{IPv4Address, "postNATSourceIPv4Address"}
		case 226:
			return fieldDefinition{IPv4Address, "postNATDestinationIPv4Address"}
		case 227:
			return fieldDefinition{Uint16, "postNAPTSourceTransportPort"}
		case 228:
			return fieldDefinition{Uint16, "postNAPTDestinationTransportPort"}
		}
	}
	return fieldDefinition{UnknownFieldType, "unknown field"}
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
