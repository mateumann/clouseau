package ipfix

import (
	"fmt"
)

// Abstract data types of the IPFIX information model
// https://tools.ietf.org/html/rfc7012#section-3.1
// https://www.iana.org/assignments/ipfix/ipfix.xhtml#ipfix-information-element-data-types
type FieldType uint8

const (
	OctetArray FieldType = iota
	Uint8
	Uint16
	Uint32
	Uint64
	Int8
	Int16
	Int32
	Int64
	Float32
	Float64
	Boolean
	MacAddress
	String
	DateTimeSeconds
	DateTimeMilliseconds
	DateTimeMicroseconds
	DateTimeNanoseconds
	IPv4Address
	IPv6Address
	BasicList
	SubTemplateList
	SubTemplateMultiList
	Unknown
)

// An ElementSpecifier identifies an IPFIX element
type ElementSpecifier struct {
	ElementID        uint16
	EnterpriseNumber uint32
}

// Returns an Enterprise bit.  This is the first bit of the ElementSpecifier.  If this bit is zero, the Information
// Element identifier identifies an Information Element in [IANA-IPFIX], and the four-octet Enterprise Number field
// MUST NOT be present.  If this bit is one, the Information Element identifier identifies an enterprise-specific
// Information Element, and the Enterprise Number field MUST be present.
func (s ElementSpecifier) E() bool {
	return s.ElementID>>15 == 1
}

func (s ElementSpecifier) String() string {
	if s.E() {
		return fmt.Sprintf("Element %d, Enterprise number %d", s.ElementID, s.EnterpriseNumber)
	}
	return fmt.Sprintf("Element %d", s.ElementID)
}

// A FieldSpecifier represents a IPFIX field specification.
// https://www.ietf.org/rfc/rfc7011.html#section-3.2
// https://www.iana.org/assignments/ipfix/ipfix.xhtml
type FieldSpecifier struct {
	ElementID        uint16
	FieldLength      uint16
	EnterpriseNumber uint32
}

// Returns an Enterprise bit.  This is the first bit of the FieldSpecifier.  If this bit is zero, the Information
// Element identifier identifies an Information Element in [IANA-IPFIX], and the four-octet Enterprise Number field
// MUST NOT be present.  If this bit is one, the Information Element identifier identifies an enterprise-specific
// Information Element, and the Enterprise Number field MUST be present.
func (s FieldSpecifier) E() bool {
	return s.ElementID>>15 == 1
}

func (s *FieldSpecifier) ElementSpec() ElementSpecifier {
	return ElementSpecifier{s.ElementID, s.EnterpriseNumber}
}

// Returns a size of a FieldSpecifier.  This depends on Enterprise bit being set (or not).
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

type FieldDefinition struct {
	Type FieldType
	Name string
}

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

func InformationElement(s ElementSpecifier) (FieldDefinition, error) {
	if fieldDefinition, ok := ipfixInformationElements()[s]; ok {
		return fieldDefinition, nil
	}
	return FieldDefinition{Unknown, ""}, fmt.Errorf("IPFIX information element not found for %v", s)
}

func (s FieldSpecifier) InformationElement() (FieldDefinition, error) {
	return InformationElement(ElementSpecifier{s.ElementID, s.EnterpriseNumber})
}

func (f FieldDefinition) NewPlaceholder() interface{} {
	switch f.Type {
	case Uint8:
		return new(uint8)
	case Uint16:
		return new(uint16)
	case Uint32:
		return new(uint32)
	case Uint64:
		return new(uint64)
	case IPv4Address:
		return new(uint32)
	case IPv6Address:
		return new(uint64)
	case MacAddress:
		return new([6]byte)
	default:
		return nil
	}
}

/*
import "reflect"

// Abstract data types of the IPFIX information model
// https://tools.ietf.org/html/rfc7012#section-3.1
type FieldType uint16

const (
    OctetArray FieldType = iota
    Uint8
    Uint16
    Uint32
    Uint64
    Int8
    Int16
    Int32
    Int64
    Float32
    Float64
    Boolean
    MacAddress
    String
    DateTimeSeconds
    DateTimeMilliseconds
    DateTimeMicroseconds
    DateTimeNanoseconds
    IPv4Address
    IPv6Address
    BasicList
    SubTemplateList
    SubTemplateMultiList
    Unknown
)

// ElementKey represents field specifier format ids
type ElementKey struct {
    ElementID    uint16
    EnterpriseNo uint32
}

// InfoElementEntry represents standard name and
// type for a field - RFC5102
type InfoElementEntry struct {
    FieldID uint16
    Name    string
    Type    FieldType
}

// IPFIX Information Elements
// https://www.iana.org/assignments/ipfix/ipfix.xhtml
type InformationElement struct {
    ElementID uint16
    Name                 string
    Type                 reflect.Type
}

var IPFIXInformationElements = []InformationElement{
    {1, "octetDeltaCount", reflect.TypeOf(reflect.Uint64)},
    {2, "packetDeltaCount", reflect.Uint64},
    {3, "deltaFlowCount", reflect.Uint64},
    {4, "protocolIdentifier", reflect.Uint8},
    {5, "ipClassOfService", reflect.Uint8},
    {6, "tcpControlBits", reflect.Uint16},
    {7, "sourceTransportPort", reflect.Uint16},
    {8, "sourceIPv4Address", reflect.Uint32},
    {9, "sourceIPv4PrefixLength", reflect.Uint8},
    {10, "ingressInterface", reflect.Uint32},
    {11, "destinationTransportPort", reflect.Uint16},
    {12, "destinationIPv4Address", reflect.Uint32},
    {13, "destinationIPv4PrefixLength", reflect.Uint8},
    {14, "egressInterface", reflect.Uint32},
    {15, "ipNextHopIPv4Address", reflect.Uint32},
    {16, "bgpSourceAsNumber", reflect.Uint32},
    {17, "bgpDestinationAsNumber", reflect.Uint32},
    {18, "bgpNextHopIPv4Address", reflect.Uint32},
    // there are many, let's focus on the ones used by MikroTik
    {21, "flowEndSysUpTime", reflect.Uint32},
    {22, "flowStartSysUpTime", reflect.Uint32},
    {60, "ipVersion", reflect.Uint8},
    {80, "destinationMacAddress", reflect.ArrayOf(6, reflect.Uint8)},
}
*/
