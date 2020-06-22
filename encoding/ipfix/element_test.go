package ipfix

import (
	"fmt"
	"testing"
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

func TestFieldSpecifier_ElementSpecifier(t *testing.T) {
	var fSpec FieldSpecifier
	for i := uint16(1); i <= 1<<15; i++ {
		fSpec = FieldSpecifier{i, 11, 0}
		t.Run(fmt.Sprintf("FieldSpecifier %v", fSpec), func(t *testing.T) {
			eSpec := fSpec.ElementSpecifier()
			if eSpec.ElementID != fSpec.ElementID || eSpec.EnterpriseNumber != fSpec.EnterpriseNumber {
				t.Errorf("FieldSpecifier.ElementSpecifier() = %v for %v", eSpec, fSpec)
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
