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
			name: "Information Element 3 (deltaFlowCount)",
			arg:  3,
			want: InformationElement{
				name:              "deltaFlowCount",
				dataType:          reflect.Uint64,
				dataTypeSemantics: deltaCounter,
				units:             flows,
			},
		},
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
		{
			name: "Information Element 189 (ipHeaderLength)",
			arg:  189,
			want: InformationElement{
				name:              "ipHeaderLength",
				dataType:          reflect.Uint8,
				dataTypeSemantics: defaultElementSemantic,
				units:             octets,
			},
		},
		{
			name: "Information Element 206 (isMulticast)",
			arg:  206,
			want: InformationElement{
				name:              "isMulticast",
				dataType:          reflect.Uint8,
				dataTypeSemantics: flags,
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
