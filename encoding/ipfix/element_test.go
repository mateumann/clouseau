package ipfix

import (
	"testing"
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
			if got := informationElement(tt.arg.elementID, tt.arg.enterpriseNumber); got != tt.want {
				t.Errorf("informationElement(%v) == %v, want %v", tt.arg, got, tt.want)
			}
		})
	}
}
