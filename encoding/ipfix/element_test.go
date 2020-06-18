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
			arg:  ElementSpecifier{43882, 0},
			want: true,
		},
		{
			name: "Enterprise bit set for maximum ElementID",
			arg:  ElementSpecifier{^uint16(0), 0},
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
