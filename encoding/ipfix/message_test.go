package ipfix

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"testing"
)

func TestFieldSpecifier_String(t *testing.T) {
	var fs fieldSpecifier

	for attempts := 100; attempts > 0; attempts-- {
		randomized, err := rand.Int(rand.Reader, big.NewInt(65536))
		if err != nil || randomized == nil {
			t.Errorf("Error occurred during getting a random number.")

			return
		}

		id := uint16(randomized.Int64())
		want := ""
		fs = fieldSpecifier{
			informationElementID: id,
			fieldLength:          id % 16,
			enterpriseNumber:     uint32(randomized.Int64() % 2),
		}

		switch {
		case id >= 32768 && fs.enterpriseNumber > 0:
			want = fmt.Sprintf("enterprise-specific field #%v (%v)", id, fs.enterpriseNumber)
		case id < 32768 && fs.enterpriseNumber == 0:
			want = fmt.Sprintf("IANA-registered field #%v", InformationElementByID(id))
		default:
			want = fmt.Sprintf("invalid field #%v", id)
		}

		got := fs.String()

		if want != got {
			t.Errorf("FieldSpecifier.String(%v) = %v, want %v", id, got, want)
		}
	}
}

func TestFieldSpecifier_e(t *testing.T) {
	var fs fieldSpecifier

	for attempts := 100; attempts > 0; attempts-- {
		randomized, err := rand.Int(rand.Reader, big.NewInt(65536))
		if err != nil || randomized == nil {
			t.Errorf("Error occurred during getting a random number.")

			return
		}

		id := uint16(randomized.Int64())
		want := id >= 32768
		fs = fieldSpecifier{
			informationElementID: id,
			fieldLength:          0,
			enterpriseNumber:     0,
		}
		got := fs.e()

		if want != got {
			t.Errorf("FieldSpecifier.e(%v) = %v, want %v", fs, got, want)
		}
	}
}

func TestFieldSpecifier_valid(t *testing.T) {
	var fs fieldSpecifier

	for attempts := 100; attempts > 0; attempts-- {
		randomized, err := rand.Int(rand.Reader, big.NewInt(65536))
		if err != nil || randomized == nil {
			t.Errorf("Error occurred during getting a random number.")

			return
		}

		id := uint16(randomized.Int64())
		if id >= 32768 {
			fs = fieldSpecifier{
				informationElementID: id,
				fieldLength:          1,
				enterpriseNumber:     uint32(1 + (randomized.Uint64() % 4294967295)),
			}
		} else {
			fs = fieldSpecifier{
				informationElementID: id,
				fieldLength:          1,
				enterpriseNumber:     0,
			}
		}

		got := fs.valid()

		if got != true {
			t.Errorf("FieldSpecifier.valid(%v) = %v, want %v", fs, got, true)
		}
	}
}

func TestFieldSpecifier_invalid(t *testing.T) {
	var fs fieldSpecifier

	for attempts := 100; attempts > 0; attempts-- {
		randomized, err := rand.Int(rand.Reader, big.NewInt(65536))
		if err != nil || randomized == nil {
			t.Errorf("Error occurred during getting a random number.")

			return
		}

		id := uint16(randomized.Int64())
		if id >= 32768 {
			fs = fieldSpecifier{
				informationElementID: id,
				fieldLength:          1,
				enterpriseNumber:     0,
			}
		} else {
			fs = fieldSpecifier{
				informationElementID: id,
				fieldLength:          1,
				enterpriseNumber:     uint32(1 + (randomized.Uint64() % 4294967295)),
			}
		}

		got := fs.valid()

		if got != false {
			t.Errorf("FieldSpecifier.valid(%v) = %v, want %v", fs, got, false)
		}
	}
}

func TestMessageHeader_String(t *testing.T) {
	var h messageHeader

	for attempts := 100; attempts > 0; attempts-- {
		randomized, err := rand.Int(rand.Reader, big.NewInt(4294967296))
		if err != nil || randomized == nil {
			t.Errorf("Error occurred during getting a random number.")

			return
		}

		id := uint32(randomized.Int64())
		h = messageHeader{
			version:             0x000a,
			length:              uint16(id % 16),
			exportTime:          id,
			sequenceNumber:      1,
			observationDomainID: 1,
		}

		want := fmt.Sprintf("header %v (%v octets, %v, %v, %v)",
			h.version, h.length, h.exportTime, h.sequenceNumber, h.observationDomainID)
		got := h.String()

		if want != got {
			t.Errorf("MessageHeader.String(%v) = %v, want %v", id, got, want)
		}
	}
}

func TestSet_String(t *testing.T) {
	tests := []struct {
		name string
		arg  set
		want string
	}{
		{
			name: "Template Set with 1 record",
			arg: set{
				header: setHeader{
					setID:  2,
					length: 16,
				},
				records: make([]interface{}, 1),
			},
			want: "template set (16 octets) 1 record",
		},
		{
			name: "Options Template Set with 2 record",
			arg: set{
				header: setHeader{
					setID:  3,
					length: 24,
				},
				records: make([]interface{}, 2),
			},
			want: "options template set (24 octets) 2 records",
		},
		{
			name: "Data Set with 1 record",
			arg: set{
				header: setHeader{
					setID:  258,
					length: 16,
				},
				records: make([]interface{}, 1),
			},
			want: "data set (16 octets) 1 record",
		},
		{
			name: "Invalid Set",
			arg: set{
				header: setHeader{
					setID:  24,
					length: 16,
				},
				records: make([]interface{}, 1),
			},
			want: "invalid set #24 1 record",
		},
	}
	for _, tt := range tests {
		tt := tt // pin
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.arg.String(); got != tt.want {
				t.Errorf("%v got %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
