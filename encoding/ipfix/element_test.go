package ipfix

import (
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
			name: "Enterprise bit set for quite large ElementID",
			arg:  ElementSpecifier{43882, 10},
			want: true,
		},
		{
			name: "Enterprise bit set for maximum ElementID",
			arg:  ElementSpecifier{^uint16(0), 11},
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
