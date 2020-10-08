package ipfix

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestElementSpecifier_E(t *testing.T) {
	tests := []struct {
		name string
		arg  ElementSpecifier
		want bool
	}{
		{
			name: "Enterprise bit not set for small ElementID",
			arg:  ElementSpecifier{0, 0},
			want: false,
		},
		{
			name: "Enterprise bit not set for larger ElementID",
			arg:  ElementSpecifier{19349, 0},
			want: false,
		},
		{
			name: "Enterprise bit not set for ElementID = 2^15 - 1",
			arg:  ElementSpecifier{1<<15 - 1, 0},
			want: false,
		},
		{
			name: "Enterprise bit set for ElementID = 2^15",
			arg:  ElementSpecifier{1 << 15, 10},
			want: true,
		},
		{
			name: "Enterprise bit set for quite large ElementID",
			arg:  ElementSpecifier{43882, 11},
			want: true,
		},
		{
			name: "Enterprise bit set for maximum ElementID",
			arg:  ElementSpecifier{^uint16(0), 12},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.arg.E(); got != tt.want {
				t.Errorf("ElementSpecifier.E() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestElementSpecifier_String(t *testing.T) {
	tests := []struct {
		name string
		arg  ElementSpecifier
		want string
	}{
		{
			name: "ElementSpecifier with no Enterprise bit set (ElementID = 0)",
			arg:  ElementSpecifier{0, 0},
			want: "Element 0",
		},
		{
			name: "ElementSpecifier with no Enterprise bit set (ElementID = 19349)",
			arg:  ElementSpecifier{19349, 0},
			want: "Element 19349",
		},
		{
			name: "ElementSpecifier with Enterprise bit set for (ElementID = 43882)",
			arg:  ElementSpecifier{43882, 10},
			want: "Element 43882, Enterprise number 10",
		},
		{
			name: "ElementSpecifier with Enterprise bit set for maximum ElementID",
			arg:  ElementSpecifier{^uint16(0), 11},
			want: "Element 65535, Enterprise number 11",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.arg.String(); got != tt.want {
				t.Errorf("ElementSpecifier.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFieldSpecifier_SemiRandom(t *testing.T) {
	var fSpec FieldSpecifier
	rand.Seed(time.Now().UnixNano())
	elementIds := []uint16{0, 1 << 15 - 1, 1 << 15, 1 << 16 - 1}
	for i := 4 ; i < 24 ; i++ {
		elementIds = append(elementIds, uint16(rand.Intn(1<<16 +1)))
	}
	for _, i := range elementIds{
		fSpec = FieldSpecifier{i, 11, 0}
		t.Run(fmt.Sprintf("FieldSpecifier %v", fSpec), func(t *testing.T) {
			eSpec := fSpec.ElementSpecifier()
			if eSpec.ElementID != fSpec.ElementID || eSpec.EnterpriseNumber != fSpec.EnterpriseNumber {
				t.Errorf("FieldSpecifier.ElementSpecifier() = %v for %v", eSpec, fSpec)
			}
			if ((i < 1<<15) && fSpec.E()) || ((i >= 1<<15) && ! fSpec.E()) {
				t.Errorf("FieldSpecifier.E() = %v, which is unexpected for %v", fSpec.E(), fSpec)
			}
			if (i >= 1<<15) && fSpec.EnterpriseNumber != 0 {
				t.Errorf("FieldSpecifier.EnterpriseNumber() = %v, should be 0 for %v", fSpec.EnterpriseNumber, fSpec)
			}
		})
	}
}

func TestFieldSpecifier_E(t *testing.T) {
	tests := []struct {
		name string
		arg  FieldSpecifier
		want bool
	}{
		{
			name: "Enterprise bit not set for small ElementID",
			arg:  FieldSpecifier{3, 8, 0},
			want: false,
		},
		{
			name: "Enterprise bit not set for larger ElementID",
			arg:  FieldSpecifier{19349, 1, 0},
			want: false,
		},
		{
			name: "Enterprise bit not set for ElementID = 2^15 - 1",
			arg:  FieldSpecifier{1<<15 - 1, 1, 0},
			want: false,
		},
		{
			name: "Enterprise bit set for ElementID = 2^15",
			arg:  FieldSpecifier{1 << 15, 8, 10},
			want: true,
		},
		{
			name: "Enterprise bit set for quite large ElementID",
			arg:  FieldSpecifier{43882, 16, 11},
			want: true,
		},
		{
			name: "Enterprise bit set for maximum ElementID",
			arg:  FieldSpecifier{^uint16(0), 32, 12},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.arg.E(); got != tt.want {
				t.Errorf("FieldSpecifier.E() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFieldSpecifier_Size(t *testing.T) {
	tests := []struct {
		name string
		arg  FieldSpecifier
		want uint16
	}{
		{
			name: "Size of 4 octets for small ElementID",
			arg:  FieldSpecifier{3, 8, 0},
			want: uint16(4),
		},
		{
			name: "Size of 4 octets for larger ElementID",
			arg:  FieldSpecifier{19349, 1, 0},
			want: uint16(4),
		},
		{
			name: "Size of 4 octets for ElementID = 2^15 - 1",
			arg:  FieldSpecifier{1<<15 - 1, 1, 0},
			want: uint16(4),
		},
		{
			name: "Size of 8 octets for ElementID = 2^15",
			arg:  FieldSpecifier{1 << 15, 8, 10},
			want: uint16(8),
		},
		{
			name: "Size of 8 octets for quite large ElementID",
			arg:  FieldSpecifier{43882, 16, 11},
			want: uint16(8),
		},
		{
			name: "Size of 8 octets for maximum ElementID",
			arg:  FieldSpecifier{^uint16(0), 32, 12},
			want: uint16(8),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.arg.Size(); got != tt.want {
				t.Errorf("FieldSpecifier.Size() == %v, want % v", got, tt.want)
			}
		})
	}
}

func TestElementSpecifier_InformationElement(t *testing.T) {
	tests := []struct {
		name       string
		arg        ElementSpecifier
		want       FieldDefinition
		errMessage string
	}{
		{
			name: "Unknown Field",
			arg:  ElementSpecifier{0, 0},
			want: FieldDefinition{UnknownFieldType, ""},
		},
		{
			name: "octetDeltaCount Field",
			arg:  ElementSpecifier{1, 0},
			want: FieldDefinition{Uint64, "octetDeltaCount"},
		},
		{
			name: "packetDeltaCount Field",
			arg:  ElementSpecifier{2, 0},
			want: FieldDefinition{Uint64, "packetDeltaCount"},
		},
		{
			name: "deltaFlowCount Field",
			arg:  ElementSpecifier{3, 0},
			want: FieldDefinition{Uint64, "deltaFlowCount"},
		},
		{
			name: "protocolIdentifier Field",
			arg:  ElementSpecifier{4, 0},
			want: FieldDefinition{Uint8, "protocolIdentifier"},
		},
		{
			name: "ipClassOfService Field",
			arg:  ElementSpecifier{5, 0},
			want: FieldDefinition{Uint8, "ipClassOfService"},
		},
		{
			name: "tcpControlBits Field",
			arg:  ElementSpecifier{6, 0},
			want: FieldDefinition{Uint16, "tcpControlBits"},
		},
		{
			name: "sourceTransportPort Field",
			arg:  ElementSpecifier{7, 0},
			want: FieldDefinition{Uint16, "sourceTransportPort"},
		},
		{
			name: "sourceIPv4Address Filed",
			arg:  ElementSpecifier{8, 0},
			want: FieldDefinition{IPv4Address, "sourceIPv4Address"},
		},
		{
			name: "sourceIPv4PrefixLength Filed",
			arg:  ElementSpecifier{9, 0},
			want: FieldDefinition{Uint8, "sourceIPv4PrefixLength"},
		},
		{
			name: "ingressInterface Filed",
			arg:  ElementSpecifier{10, 0},
			want: FieldDefinition{Uint32, "ingressInterface"},
		},
		{
			name: "destinationTransportPort Filed",
			arg:  ElementSpecifier{11, 0},
			want: FieldDefinition{Uint16, "destinationTransportPort"},
		},
		{
			name: "destinationIPv4Address Filed",
			arg:  ElementSpecifier{12, 0},
			want: FieldDefinition{IPv4Address, "destinationIPv4Address"},
		},
		{
			name: "destinationIPv4PrefixLength Filed",
			arg:  ElementSpecifier{13, 0},
			want: FieldDefinition{Uint8, "destinationIPv4PrefixLength"},
		},
		{
			name: "egressInterface Filed",
			arg:  ElementSpecifier{14, 0},
			want: FieldDefinition{Uint32, "egressInterface"},
		},
		{
			name: "ipNextHopIPv4Address Filed",
			arg:  ElementSpecifier{15, 0},
			want: FieldDefinition{IPv4Address, "ipNextHopIPv4Address"},
		},
		{
			name: "flowEndSysUpTime Filed",
			arg:  ElementSpecifier{21, 0},
			want: FieldDefinition{Uint32, "flowEndSysUpTime"},
		},
		{
			name: "flowStartSysUpTime Filed",
			arg:  ElementSpecifier{22, 0},
			want: FieldDefinition{Uint32, "flowStartSysUpTime"},
		},
		{
			name: "sourceMacAddress Filed",
			arg:  ElementSpecifier{56, 0},
			want: FieldDefinition{MacAddress, "sourceMacAddress"},
		},
		{
			name: "ipVersion Filed",
			arg:  ElementSpecifier{60, 0},
			want: FieldDefinition{Uint8, "ipVersion"},
		},
		{
			name: "ipNextHopIPv6Address Filed",
			arg:  ElementSpecifier{62, 0},
			want: FieldDefinition{IPv6Address, "ipNextHopIPv6Address"},
		},
		{
			name: "destinationMacAddress Filed",
			arg:  ElementSpecifier{80, 0},
			want: FieldDefinition{MacAddress, "destinationMacAddress"},
		},
		{
			name: "postSourceMacAddress Filed",
			arg:  ElementSpecifier{81, 0},
			want: FieldDefinition{MacAddress, "postSourceMacAddress"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := InformationElement(tt.arg); got != tt.want {
				t.Errorf("InformationElement(%v) == %v, want %v", tt.arg, got, tt.want)
			}
		})
	}
}

func TestFieldSpecifier_String(t *testing.T) {
	tests := []struct {
		name string
		arg  FieldSpecifier
		want string
	}{
		{
			name: "FieldSpecifier with no Enterprise bit set (ElementID = 0)",
			arg:  FieldSpecifier{0, 0, 0},
			want: "Element 0",
		},
		{
			name: "FieldSpecifier with no Enterprise bit set (ElementID = 19349)",
			arg:  FieldSpecifier{19349, 11, 0},
			want: "Element 19349",
		},
		{
			name: "FieldSpecifier with Enterprise bit set for (ElementID = 43882)",
			arg:  FieldSpecifier{43882, 12, 10},
			want: "Element 43882, Enterprise number 10",
		},
		{
			name: "FieldSpecifier with Enterprise bit set for maximum ElementID",
			arg:  FieldSpecifier{^uint16(0), ^uint16(0), 11},
			want: "Element 65535, Enterprise number 11",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.arg.String(); got != tt.want {
				t.Errorf("FieldSpecifier.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFieldDefinition_String(t *testing.T) {
	tests := []struct {
		name string
		arg  FieldDefinition
		want string
	}{
		{
			name: "Unknown FieldDefinition",
			arg:  FieldDefinition{UnknownFieldType, ""},
			want: "Unknown Field",
		},
		{
			name: "FieldDefinition for octetDeltaCount",
			arg:  FieldDefinition{Uint64, "octetDeltaCount"},
			want: "Field 'octetDeltaCount'",
		},
		{
			name: "FieldDefinition for sourceIPv6Address",
			arg:  FieldDefinition{IPv6Address, "sourceIPv6Address"},
			want: "Field 'sourceIPv6Address'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.arg.String(); got != tt.want {
				t.Errorf("FieldDefinition.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFieldSpecifier_InformationElement(t *testing.T) {
	tests := []struct {
		name       string
		arg        FieldSpecifier
		want       FieldDefinition
		errMessage string
	}{
		{
			name: "Unknown Field",
			arg:  FieldSpecifier{0, 0, 0},
			want: FieldDefinition{UnknownFieldType, ""},
		},
		{
			name: "octetDeltaCount Field",
			arg:  FieldSpecifier{1, 8, 0},
			want: FieldDefinition{Uint64, "octetDeltaCount"},
		},
		{
			name: "packetDeltaCount Field",
			arg:  FieldSpecifier{2, 8, 0},
			want: FieldDefinition{Uint64, "packetDeltaCount"},
		},
		{
			name: "deltaFlowCount Field",
			arg:  FieldSpecifier{3, 8, 0},
			want: FieldDefinition{Uint64, "deltaFlowCount"},
		},
		{
			name: "protocolIdentifier Field",
			arg:  FieldSpecifier{4, 1, 0},
			want: FieldDefinition{Uint8, "protocolIdentifier"},
		},
		{
			name: "ipClassOfService Field",
			arg:  FieldSpecifier{5, 1, 0},
			want: FieldDefinition{Uint8, "ipClassOfService"},
		},
		{
			name: "tcpControlBits Field",
			arg:  FieldSpecifier{6, 2, 0},
			want: FieldDefinition{Uint16, "tcpControlBits"},
		},
		{
			name: "sourceTransportPort Field",
			arg:  FieldSpecifier{7, 2, 0},
			want: FieldDefinition{Uint16, "sourceTransportPort"},
		},
		{
			name: "sourceIPv4Address Filed",
			arg:  FieldSpecifier{8, 4, 0},
			want: FieldDefinition{IPv4Address, "sourceIPv4Address"},
		},
		{
			name: "sourceIPv4PrefixLength Filed",
			arg:  FieldSpecifier{9, 1, 0},
			want: FieldDefinition{Uint8, "sourceIPv4PrefixLength"},
		},
		{
			name: "ingressInterface Filed",
			arg:  FieldSpecifier{10, 4, 0},
			want: FieldDefinition{Uint32, "ingressInterface"},
		},
		{
			name: "destinationTransportPort Filed",
			arg:  FieldSpecifier{11, 2, 0},
			want: FieldDefinition{Uint16, "destinationTransportPort"},
		},
		{
			name: "destinationIPv4Address Filed",
			arg:  FieldSpecifier{12, 4, 0},
			want: FieldDefinition{IPv4Address, "destinationIPv4Address"},
		},
		{
			name: "destinationIPv4PrefixLength Filed",
			arg:  FieldSpecifier{13, 1, 0},
			want: FieldDefinition{Uint8, "destinationIPv4PrefixLength"},
		},
		{
			name: "egressInterface Filed",
			arg:  FieldSpecifier{14, 4, 0},
			want: FieldDefinition{Uint32, "egressInterface"},
		},
		{
			name: "ipNextHopIPv4Address Filed",
			arg:  FieldSpecifier{15, 4, 0},
			want: FieldDefinition{IPv4Address, "ipNextHopIPv4Address"},
		},
		{
			name: "flowEndSysUpTime Filed",
			arg:  FieldSpecifier{21, 4, 0},
			want: FieldDefinition{Uint32, "flowEndSysUpTime"},
		},
		{
			name: "flowStartSysUpTime Filed",
			arg:  FieldSpecifier{22, 4, 0},
			want: FieldDefinition{Uint32, "flowStartSysUpTime"},
		},
		{
			name: "sourceMacAddress Filed",
			arg:  FieldSpecifier{56, 6, 0},
			want: FieldDefinition{MacAddress, "sourceMacAddress"},
		},
		{
			name: "ipVersion Filed",
			arg:  FieldSpecifier{60, 1, 0},
			want: FieldDefinition{Uint8, "ipVersion"},
		},
		{
			name: "ipNextHopIPv6Address Filed",
			arg:  FieldSpecifier{62, 8, 0},
			want: FieldDefinition{IPv6Address, "ipNextHopIPv6Address"},
		},
		{
			name: "destinationMacAddress Filed",
			arg:  FieldSpecifier{80, 6, 0},
			want: FieldDefinition{MacAddress, "destinationMacAddress"},
		},
		{
			name: "postSourceMacAddress Filed",
			arg:  FieldSpecifier{81, 6, 0},
			want: FieldDefinition{MacAddress, "postSourceMacAddress"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := tt.arg.InformationElement(); got != tt.want {
				t.Errorf("FieldSpecifier.InformationElement() == %v, want %v", got, tt.want)
			}
		})
	}
}
