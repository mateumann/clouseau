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

	for attempts := 1000; attempts > 0; attempts-- {
		randomized, err := rand.Int(rand.Reader, big.NewInt(32768))
		if err != nil || randomized == nil {
			t.Errorf("Error occurred during getting a random number.")
		}

		id := uint16(randomized.Int64())
		if want, ok = ieMap[id]; !ok {
			want = InformationElement{
				name:              "UNKNOWN",
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
	var want InformationElement

	var ok bool

	ieMap := constInformationElements()

	for attempts := 1000; attempts > 0; attempts-- {
		randomized, err := rand.Int(rand.Reader, big.NewInt(32768))
		if err != nil || randomized == nil {
			t.Errorf("Error occurred during getting a random number.")
		}

		id := uint16(randomized.Int64())
		if want, ok = ieMap[id]; !ok {
			want = InformationElement{
				name:              "UNKNOWN",
				dataType:          reflect.Invalid,
				dataTypeSemantics: unassignedElementSemantic,
				units:             unassignedElementUnit,
			}
		}

		got := InformationElementByID(id)
		if want.String() != got.String() {
			t.Errorf("InformationElementByID(%v) = %v, want %v", id, got, want)
		}
	}
}
