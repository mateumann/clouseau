package ipfix

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"testing"
	"time"
)

func TestFieldSpecifier_StringRand(t *testing.T) {
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

func TestMessage_String(t *testing.T) {
	tests := []struct {
		name string
		arg  message
		want string
	}{
		{
			name: "A message with one template set",
			arg: message{
				version:             0x000a,
				length:              296,
				exportTime:          0x5ee715b1,
				sequenceNumber:      0,
				observationDomainID: 0,
				sets: []set{
					{
						setID:   2,
						length:  280,
						records: make([]interface{}, 2),
					},
				},
			},
			want: "message 10 (296 octets, 2020-06-15T06:31:13Z, 0, 0) 1 set",
		},
		{
			name: "A message with one data set",
			arg: message{
				version:             0x000a,
				length:              412,
				exportTime:          0x5ee7163d,
				sequenceNumber:      195,
				observationDomainID: 0,
				sets: []set{
					{
						setID:   259,
						length:  396,
						records: make([]interface{}, 1),
					},
				},
			},
			want: "message 10 (412 octets, 2020-06-15T06:33:33Z, 195, 0) 1 set",
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

func TestMessage_StringRand(t *testing.T) {
	var h message

	for attempts := 100; attempts > 0; attempts-- {
		randomized, err := rand.Int(rand.Reader, big.NewInt(4294967296))
		if err != nil || randomized == nil {
			t.Errorf("Error occurred during getting a random number.")

			return
		}

		id := uint32(randomized.Int64())
		h = message{
			version:             0x000a,
			length:              uint16(id % 16),
			exportTime:          id,
			sequenceNumber:      1,
			observationDomainID: 1,
		}

		want := fmt.Sprintf("message %v (%v octets, %v, %v, %v) 0 sets",
			h.version, h.length, time.Unix(int64(h.exportTime), 0).UTC().Format(time.RFC3339),
			h.sequenceNumber, h.observationDomainID)
		got := h.String()

		if want != got {
			t.Errorf("Message.String(%v) = %v, want %v", id, got, want)
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
				setID:   2,
				length:  16,
				records: make([]interface{}, 1),
			},
			want: "template set (16 octets) 1 record",
		},
		{
			name: "Options Template Set with 2 record",
			arg: set{
				setID:   3,
				length:  24,
				records: make([]interface{}, 2),
			},
			want: "options template set (24 octets) 2 records",
		},
		{
			name: "Data Set with 1 record",
			arg: set{
				setID:   258,
				length:  16,
				records: make([]interface{}, 1),
			},
			want: "data set (16 octets) 1 record",
		},
		{
			name: "Invalid Set",
			arg: set{
				setID:   24,
				length:  16,
				records: make([]interface{}, 1),
			},
			want: "invalid set #24",
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

func TestTemplateRecord_String(t *testing.T) {
	tests := []struct {
		name string
		arg  templateRecord
		want string
	}{
		{
			name: "Template Record #258 with 1 field specifier",
			arg: templateRecord{
				templateID:      258,
				fieldCount:      1,
				fieldSpecifiers: make([]fieldSpecifier, 1),
			},
			want: "template record #258 (1 field)",
		},
		{
			name: "Template Record #259 with 19 field specifier",
			arg: templateRecord{
				templateID:      259,
				fieldCount:      19,
				fieldSpecifiers: make([]fieldSpecifier, 19),
			},
			want: "template record #259 (19 fields)",
		},
		{
			name: "Template Record #260 with 0 field specifier",
			arg: templateRecord{
				templateID:      260,
				fieldCount:      0,
				fieldSpecifiers: make([]fieldSpecifier, 0),
			},
			want: "template record #260 (0 fields)",
		},
		{
			name: "Template Record #24 with 1 field specifier",
			arg: templateRecord{

				templateID:      24,
				fieldCount:      1,
				fieldSpecifiers: make([]fieldSpecifier, 1),
			},
			want: "invalid template record #24",
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

func TestOptionsTemplateRecord_String(t *testing.T) {
	tests := []struct {
		name string
		arg  optionsTemplateRecord
		want string
	}{
		{
			name: "Options Template Record #258 with 3 field specifier",
			arg: optionsTemplateRecord{
				templateID:      258,
				fieldCount:      3,
				scopeFieldCount: 2,
				fieldSpecifiers: make([]fieldSpecifier, 3),
			},
			want: "options template record #258 (scope 2 fields of 3 fields)",
		},
		{
			name: "Options Template Record #24 with 1 field specifier",
			arg: optionsTemplateRecord{

				templateID:      24,
				fieldCount:      1,
				scopeFieldCount: 1,
				fieldSpecifiers: make([]fieldSpecifier, 1),
			},
			want: "invalid options template record #24",
		},
		{
			name: "Options Template Record #384 with 0 scope fields",
			arg: optionsTemplateRecord{

				templateID:      384,
				fieldCount:      1,
				scopeFieldCount: 0,
				fieldSpecifiers: make([]fieldSpecifier, 1),
			},
			want: "invalid options template record #384",
		},
		{
			name: "Options Template Record #385 with 2 field specifier",
			arg: optionsTemplateRecord{
				templateID:      385,
				fieldCount:      2,
				scopeFieldCount: 3,
				fieldSpecifiers: make([]fieldSpecifier, 2),
			},
			want: "invalid options template record #385",
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
