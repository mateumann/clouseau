package ipfix

import (
	"crypto/rand"
	"math/big"
	"reflect"
	"testing"
)

func TestInformationElementByID(t *testing.T) {
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

func TestInformationElementByIDAll(t *testing.T) {
	for id, want := range constInformationElements() {
		got := InformationElementByID(id)
		if want != got {
			t.Errorf("InformationElementByID(%v) = %v, want %v", id, got, want)
		}
	}
}

func TestInformationElementByIDRandom(t *testing.T) {
	var want InformationElement

	var ok bool

	ieMap := constInformationElements()

	for attempts := 100; attempts > 0; attempts-- {
		randomized, err := rand.Int(rand.Reader, big.NewInt(32768))
		if err != nil || randomized == nil {
			t.Errorf("Error occurred during getting a random number.")

			return
		}

		id := uint16(randomized.Int64())
		if want, ok = ieMap[id]; !ok {
			want = InformationElement{
				name:              "unknown",
				dataType:          reflect.Invalid,
				dataTypeSemantics: unassignedElementSemantic,
				units:             unassignedElementUnit,
			}
		}

		got := InformationElementByID(id)
		if want != got {
			t.Errorf("InformationElementByID(%v) = %v, want %v", id, got, want)
		}
	}
}

func TestInformationElement_String(t *testing.T) {
	tests := []struct {
		name string
		arg  InformationElement
		want string
	}{
		{
			name: "Information Element 0",
			arg:  InformationElementByID(0),
			want: "reserved",
		},
		{
			name: "Information Element 22",
			arg:  InformationElementByID(22),
			want: "flowStartSysUpTime",
		},
		{
			name: "Information Element 35",
			arg:  InformationElementByID(35),
			want: "samplingAlgorithm",
		},
		{
			name: "Information Element 283",
			arg:  InformationElementByID(283),
			want: "natPoolId",
		},
	}
	for _, tt := range tests {
		tt := tt // pin
		t.Run(tt.name, func(t *testing.T) {
			got := tt.arg.String()
			if got != tt.want {
				t.Errorf("%v = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
