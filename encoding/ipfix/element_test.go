package ipfix

import (
	"reflect"
	"testing"
)

func TestFuncInformationElement(t *testing.T) {
	tests := []struct {
		name string
		arg  uint16
		want InformationElement
	}{
		{
			name: "Information Element 8 (sourceIPv4Address)",
			arg:  8,
			want: InformationElement{
				name:              "sourceIPv4Address",
				dataType:          reflect.Uint32,
				dataTypeSemantics: defaultElementSemantic,
				units:             noneElementUnit,
			},
		},
	}
	for _, tt := range tests {
		tt := tt // pin
		t.Run(tt.name, func(t *testing.T) {
			if got := InformationElementByID(tt.arg); got != tt.want {
				t.Errorf("InformationElementByID(%v) = %v, want %v", tt.arg, got, tt.want)
			}
		})
	}
}
