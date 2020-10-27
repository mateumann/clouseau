package ipfix

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestElementError_Error(t *testing.T) {
	tests := []struct {
		name string
		arg  elementError
		want string
	}{
		{
			name: "Error one",
			arg: elementError{
				msg: "an error message one",
			},
			want: "an error message one",
		},
		{
			name: "Error two",
			arg: elementError{
				msg: "an error message two",
			},
			want: "an error message two",
		},
	}

	for _, tt := range tests {
		tt := tt // pin!
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.arg.Error(); got != tt.want {
				t.Errorf("elementError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestElementSpecifier_E(t *testing.T) {
	tests := []struct {
		name string
		arg  elementSpecifier
		want bool
	}{
		{
			name: "Enterprise bit not set for small elementID",
			arg:  elementSpecifier{0, 0},
			want: false,
		},
		{
			name: "Enterprise bit not set for larger elementID",
			arg:  elementSpecifier{19349, 0},
			want: false,
		},
		{
			name: "Enterprise bit not set for elementID = 2^15 - 1",
			arg:  elementSpecifier{1<<15 - 1, 0},
			want: false,
		},
		{
			name: "Enterprise bit set for elementID = 2^15",
			arg:  elementSpecifier{1 << 15, 10},
			want: true,
		},
		{
			name: "Enterprise bit set for quite large elementID",
			arg:  elementSpecifier{43882, 11},
			want: true,
		},
		{
			name: "Enterprise bit set for maximum elementID",
			arg:  elementSpecifier{^uint16(0), 12},
			want: true,
		},
	}
	for _, tt := range tests {
		tt := tt // pin!
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.arg.E(); got != tt.want {
				t.Errorf("elementSpecifier.E() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestElementSpecifier_String(t *testing.T) {
	tests := []struct {
		name string
		arg  elementSpecifier
		want string
	}{
		{
			name: "elementSpecifier with no Enterprise bit set (elementID = 0)",
			arg:  elementSpecifier{0, 0},
			want: "Element 0",
		},
		{
			name: "elementSpecifier with no Enterprise bit set (elementID = 19349)",
			arg:  elementSpecifier{19349, 0},
			want: "Element 19349",
		},
		{
			name: "elementSpecifier with Enterprise bit set for (elementID = 43882)",
			arg:  elementSpecifier{43882, 10},
			want: "Element 43882, Enterprise number 10",
		},
		{
			name: "elementSpecifier with Enterprise bit set for maximum elementID",
			arg:  elementSpecifier{^uint16(0), 11},
			want: "Element 65535, Enterprise number 11",
		},
	}
	for _, tt := range tests {
		tt := tt // pin!
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.arg.String(); got != tt.want {
				t.Errorf("elementSpecifier.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFieldSpecifier_SemiRandom(t *testing.T) {
	var fSpec fieldSpecifier

	rand.Seed(time.Now().UnixNano())

	elementIds := []uint16{0, 1<<15 - 1, 1 << 15, 1<<16 - 1}

	for i := 4; i < 24; i++ {
		elementIds = append(elementIds, uint16(rand.Intn(1<<16+1))) //nolint:gosec
	}

	for _, i := range elementIds {
		i := i // pin!
		fSpec = fieldSpecifier{i, 11, 0}
		t.Run(fmt.Sprintf("fieldSpecifier %v", fSpec), func(t *testing.T) {
			eSpec := fSpec.ElementSpecifier()
			if eSpec.elementID != fSpec.elementID || eSpec.enterpriseNumber != fSpec.enterpriseNumber {
				t.Errorf("fieldSpecifier.elementSpecifier() = %v for %v", eSpec, fSpec)
			}
			if ((i < 1<<15) && fSpec.E()) || ((i >= 1<<15) && !fSpec.E()) {
				t.Errorf("fieldSpecifier.E() = %v, which is unexpected for %v", fSpec.E(), fSpec)
			}
			if (i >= 1<<15) && fSpec.enterpriseNumber != 0 {
				t.Errorf("fieldSpecifier.enterpriseNumber() = %v, should be 0 for %v", fSpec.enterpriseNumber, fSpec)
			}
		})
	}
}

func TestFieldSpecifier_E(t *testing.T) {
	tests := []struct {
		name string
		arg  fieldSpecifier
		want bool
	}{
		{
			name: "Enterprise bit not set for small elementID",
			arg:  fieldSpecifier{3, 8, 0},
			want: false,
		},
		{
			name: "Enterprise bit not set for larger elementID",
			arg:  fieldSpecifier{19349, 1, 0},
			want: false,
		},
		{
			name: "Enterprise bit not set for elementID = 2^15 - 1",
			arg:  fieldSpecifier{1<<15 - 1, 1, 0},
			want: false,
		},
		{
			name: "Enterprise bit set for elementID = 2^15",
			arg:  fieldSpecifier{1 << 15, 8, 10},
			want: true,
		},
		{
			name: "Enterprise bit set for quite large elementID",
			arg:  fieldSpecifier{43882, 16, 11},
			want: true,
		},
		{
			name: "Enterprise bit set for maximum elementID",
			arg:  fieldSpecifier{^uint16(0), 32, 12},
			want: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.arg.E(); got != tt.want {
				t.Errorf("fieldSpecifier.E() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFieldSpecifier_Size(t *testing.T) {
	tests := []struct {
		name string
		arg  fieldSpecifier
		want uint16
	}{
		{
			name: "Size of 4 octets for small elementID",
			arg:  fieldSpecifier{3, 8, 0},
			want: uint16(4),
		},
		{
			name: "Size of 4 octets for larger elementID",
			arg:  fieldSpecifier{19349, 1, 0},
			want: uint16(4),
		},
		{
			name: "Size of 4 octets for elementID = 2^15 - 1",
			arg:  fieldSpecifier{1<<15 - 1, 1, 0},
			want: uint16(4),
		},
		{
			name: "Size of 8 octets for elementID = 2^15",
			arg:  fieldSpecifier{1 << 15, 8, 10},
			want: uint16(8),
		},
		{
			name: "Size of 8 octets for quite large elementID",
			arg:  fieldSpecifier{43882, 16, 11},
			want: uint16(8),
		},
		{
			name: "Size of 8 octets for maximum elementID",
			arg:  fieldSpecifier{^uint16(0), 32, 12},
			want: uint16(8),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.arg.Size(); got != tt.want {
				t.Errorf("fieldSpecifier.Size() == %v, want %v", got, tt.want)
			}
		})
	}
}

//nolint:funlen
func TestElementSpecifier_InformationElement(t *testing.T) {
	tests := []struct {
		name       string
		arg        elementSpecifier
		want       fieldDefinition
		errMessage string
	}{
		{
			name: "Unknown Field",
			arg:  elementSpecifier{0, 0},
			want: fieldDefinition{UnknownFieldType, ""},
		},
		{
			name: "octetDeltaCount Field",
			arg:  elementSpecifier{1, 0},
			want: fieldDefinition{Uint64, "octetDeltaCount"},
		},
		{
			name: "packetDeltaCount Field",
			arg:  elementSpecifier{2, 0},
			want: fieldDefinition{Uint64, "packetDeltaCount"},
		},
		{
			name: "deltaFlowCount Field",
			arg:  elementSpecifier{3, 0},
			want: fieldDefinition{Uint64, "deltaFlowCount"},
		},
		{
			name: "protocolIdentifier Field",
			arg:  elementSpecifier{4, 0},
			want: fieldDefinition{Uint8, "protocolIdentifier"},
		},
		{
			name: "ipClassOfService Field",
			arg:  elementSpecifier{5, 0},
			want: fieldDefinition{Uint8, "ipClassOfService"},
		},
		{
			name: "tcpControlBits Field",
			arg:  elementSpecifier{6, 0},
			want: fieldDefinition{Uint16, "tcpControlBits"},
		},
		{
			name: "sourceTransportPort Field",
			arg:  elementSpecifier{7, 0},
			want: fieldDefinition{Uint16, "sourceTransportPort"},
		},
		{
			name: "sourceIPv4Address Filed",
			arg:  elementSpecifier{8, 0},
			want: fieldDefinition{IPv4Address, "sourceIPv4Address"},
		},
		{
			name: "sourceIPv4PrefixLength Filed",
			arg:  elementSpecifier{9, 0},
			want: fieldDefinition{Uint8, "sourceIPv4PrefixLength"},
		},
		{
			name: "ingressInterface Filed",
			arg:  elementSpecifier{10, 0},
			want: fieldDefinition{Uint32, "ingressInterface"},
		},
		{
			name: "destinationTransportPort Filed",
			arg:  elementSpecifier{11, 0},
			want: fieldDefinition{Uint16, "destinationTransportPort"},
		},
		{
			name: "destinationIPv4Address Filed",
			arg:  elementSpecifier{12, 0},
			want: fieldDefinition{IPv4Address, "destinationIPv4Address"},
		},
		{
			name: "destinationIPv4PrefixLength Filed",
			arg:  elementSpecifier{13, 0},
			want: fieldDefinition{Uint8, "destinationIPv4PrefixLength"},
		},
		{
			name: "egressInterface Filed",
			arg:  elementSpecifier{14, 0},
			want: fieldDefinition{Uint32, "egressInterface"},
		},
		{
			name: "ipNextHopIPv4Address Filed",
			arg:  elementSpecifier{15, 0},
			want: fieldDefinition{IPv4Address, "ipNextHopIPv4Address"},
		},
		{
			name: "flowEndSysUpTime Filed",
			arg:  elementSpecifier{21, 0},
			want: fieldDefinition{Uint32, "flowEndSysUpTime"},
		},
		{
			name: "flowStartSysUpTime Filed",
			arg:  elementSpecifier{22, 0},
			want: fieldDefinition{Uint32, "flowStartSysUpTime"},
		},
		{
			name: "sourceMacAddress Filed",
			arg:  elementSpecifier{56, 0},
			want: fieldDefinition{MacAddress, "sourceMacAddress"},
		},
		{
			name: "ipVersion Filed",
			arg:  elementSpecifier{60, 0},
			want: fieldDefinition{Uint8, "ipVersion"},
		},
		{
			name: "ipNextHopIPv6Address Filed",
			arg:  elementSpecifier{62, 0},
			want: fieldDefinition{IPv6Address, "ipNextHopIPv6Address"},
		},
		{
			name: "destinationMacAddress Filed",
			arg:  elementSpecifier{80, 0},
			want: fieldDefinition{MacAddress, "destinationMacAddress"},
		},
		{
			name: "postSourceMacAddress Filed",
			arg:  elementSpecifier{81, 0},
			want: fieldDefinition{MacAddress, "postSourceMacAddress"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := informationElement(tt.arg); got != tt.want {
				t.Errorf("informationElement(%v) == %v, want %v", tt.arg, got, tt.want)
			}
		})
	}
}

func TestFieldSpecifier_String(t *testing.T) {
	tests := []struct {
		name string
		arg  fieldSpecifier
		want string
	}{
		{
			name: "fieldSpecifier with no Enterprise bit set (elementID = 0)",
			arg:  fieldSpecifier{0, 0, 0},
			want: "Element 0",
		},
		{
			name: "fieldSpecifier with no Enterprise bit set (elementID = 19349)",
			arg:  fieldSpecifier{19349, 11, 0},
			want: "Element 19349",
		},
		{
			name: "fieldSpecifier with Enterprise bit set for (elementID = 43882)",
			arg:  fieldSpecifier{43882, 12, 10},
			want: "Element 43882, Enterprise number 10",
		},
		{
			name: "fieldSpecifier with Enterprise bit set for maximum elementID",
			arg:  fieldSpecifier{^uint16(0), ^uint16(0), 11},
			want: "Element 65535, Enterprise number 11",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.arg.String(); got != tt.want {
				t.Errorf("fieldSpecifier.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFieldDefinition_String(t *testing.T) {
	tests := []struct {
		name string
		arg  fieldDefinition
		want string
	}{
		{
			name: "Unknown fieldDefinition",
			arg:  fieldDefinition{UnknownFieldType, ""},
			want: "Unknown Field",
		},
		{
			name: "fieldDefinition for octetDeltaCount",
			arg:  fieldDefinition{Uint64, "octetDeltaCount"},
			want: "Field 'octetDeltaCount'",
		},
		{
			name: "fieldDefinition for sourceIPv6Address",
			arg:  fieldDefinition{IPv6Address, "sourceIPv6Address"},
			want: "Field 'sourceIPv6Address'",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.arg.String(); got != tt.want {
				t.Errorf("fieldDefinition.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

//nolint:funlen
func TestFieldSpecifier_InformationElement(t *testing.T) {
	tests := []struct {
		name       string
		arg        fieldSpecifier
		want       fieldDefinition
		errMessage string
	}{
		{
			name: "Unknown Field",
			arg:  fieldSpecifier{0, 0, 0},
			want: fieldDefinition{UnknownFieldType, ""},
		},
		{
			name: "octetDeltaCount Field",
			arg:  fieldSpecifier{1, 8, 0},
			want: fieldDefinition{Uint64, "octetDeltaCount"},
		},
		{
			name: "packetDeltaCount Field",
			arg:  fieldSpecifier{2, 8, 0},
			want: fieldDefinition{Uint64, "packetDeltaCount"},
		},
		{
			name: "deltaFlowCount Field",
			arg:  fieldSpecifier{3, 8, 0},
			want: fieldDefinition{Uint64, "deltaFlowCount"},
		},
		{
			name: "protocolIdentifier Field",
			arg:  fieldSpecifier{4, 1, 0},
			want: fieldDefinition{Uint8, "protocolIdentifier"},
		},
		{
			name: "ipClassOfService Field",
			arg:  fieldSpecifier{5, 1, 0},
			want: fieldDefinition{Uint8, "ipClassOfService"},
		},
		{
			name: "tcpControlBits Field",
			arg:  fieldSpecifier{6, 2, 0},
			want: fieldDefinition{Uint16, "tcpControlBits"},
		},
		{
			name: "sourceTransportPort Field",
			arg:  fieldSpecifier{7, 2, 0},
			want: fieldDefinition{Uint16, "sourceTransportPort"},
		},
		{
			name: "sourceIPv4Address Filed",
			arg:  fieldSpecifier{8, 4, 0},
			want: fieldDefinition{IPv4Address, "sourceIPv4Address"},
		},
		{
			name: "sourceIPv4PrefixLength Filed",
			arg:  fieldSpecifier{9, 1, 0},
			want: fieldDefinition{Uint8, "sourceIPv4PrefixLength"},
		},
		{
			name: "ingressInterface Filed",
			arg:  fieldSpecifier{10, 4, 0},
			want: fieldDefinition{Uint32, "ingressInterface"},
		},
		{
			name: "destinationTransportPort Filed",
			arg:  fieldSpecifier{11, 2, 0},
			want: fieldDefinition{Uint16, "destinationTransportPort"},
		},
		{
			name: "destinationIPv4Address Filed",
			arg:  fieldSpecifier{12, 4, 0},
			want: fieldDefinition{IPv4Address, "destinationIPv4Address"},
		},
		{
			name: "destinationIPv4PrefixLength Filed",
			arg:  fieldSpecifier{13, 1, 0},
			want: fieldDefinition{Uint8, "destinationIPv4PrefixLength"},
		},
		{
			name: "egressInterface Filed",
			arg:  fieldSpecifier{14, 4, 0},
			want: fieldDefinition{Uint32, "egressInterface"},
		},
		{
			name: "ipNextHopIPv4Address Filed",
			arg:  fieldSpecifier{15, 4, 0},
			want: fieldDefinition{IPv4Address, "ipNextHopIPv4Address"},
		},
		{
			name: "flowEndSysUpTime Filed",
			arg:  fieldSpecifier{21, 4, 0},
			want: fieldDefinition{Uint32, "flowEndSysUpTime"},
		},
		{
			name: "flowStartSysUpTime Filed",
			arg:  fieldSpecifier{22, 4, 0},
			want: fieldDefinition{Uint32, "flowStartSysUpTime"},
		},
		{
			name: "sourceMacAddress Filed",
			arg:  fieldSpecifier{56, 6, 0},
			want: fieldDefinition{MacAddress, "sourceMacAddress"},
		},
		{
			name: "ipVersion Filed",
			arg:  fieldSpecifier{60, 1, 0},
			want: fieldDefinition{Uint8, "ipVersion"},
		},
		{
			name: "ipNextHopIPv6Address Filed",
			arg:  fieldSpecifier{62, 8, 0},
			want: fieldDefinition{IPv6Address, "ipNextHopIPv6Address"},
		},
		{
			name: "destinationMacAddress Filed",
			arg:  fieldSpecifier{80, 6, 0},
			want: fieldDefinition{MacAddress, "destinationMacAddress"},
		},
		{
			name: "postSourceMacAddress Filed",
			arg:  fieldSpecifier{81, 6, 0},
			want: fieldDefinition{MacAddress, "postSourceMacAddress"},
		},
	}
	for _, tt := range tests {
		tt := tt // pin!
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := tt.arg.InformationElement(); got != tt.want {
				t.Errorf("fieldSpecifier.informationElement() == %v, want %v", got, tt.want)
			}
		})
	}
}
