package ipfix

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"math/big"
	"testing"
	"time"
)

func TestUnmarshalError_Error(t *testing.T) {
	tests := []struct {
		name string
		arg  UnmarshalError
		want string
	}{
		{
			name: "UnmarshalError: EOF",
			arg:  UnmarshalError{io.EOF},
			want: "ipfix: unmarshalling error: EOF",
		},
		{
			name: "UnmarshalError: unexpected EOF",
			arg:  UnmarshalError{io.ErrUnexpectedEOF},
			want: "ipfix: unmarshalling error: unexpected EOF",
		},
	}
	for _, tt := range tests {
		tt := tt // pin
		t.Run(tt.name, func(t *testing.T) {
			got := tt.arg.Error()
			if got != tt.want {
				t.Errorf("%v = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

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
				id, 1, uint32(1 + (randomized.Uint64() % 4294967295)),
			}
		} else {
			fs = fieldSpecifier{
				id, 1, 0,
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
				id, 1, 0,
			}
		} else {
			fs = fieldSpecifier{
				id, 1, uint32(1 + (randomized.Uint64() % 4294967295)),
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
		arg  Message
		want string
	}{
		{
			name: "A Message with one template set",
			arg: Message{
				0x000a, 296, 0x5ee715b1, 0, 0,
				[]set{{2, 280, make([]interface{}, 2)}},
			},
			want: "Message 10 (296 octets, 2020-06-15T06:31:13Z, 0, 0) 1 set",
		},
		{
			name: "A Message with one data set",
			arg: Message{
				0x000a, 412, 0x5ee7163d, 195, 0,
				[]set{{259, 396, make([]interface{}, 1)}},
			},
			want: "Message 10 (412 octets, 2020-06-15T06:33:33Z, 195, 0) 1 set",
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
	var h Message

	for attempts := 100; attempts > 0; attempts-- {
		randomized, err := rand.Int(rand.Reader, big.NewInt(4294967296))
		if err != nil || randomized == nil {
			t.Errorf("Error occurred during getting a random number.")

			return
		}

		id := uint32(randomized.Int64())
		h = Message{
			0x000a, uint16(id % 16), id, 1, 1, nil,
		}

		want := fmt.Sprintf("Message %v (%v octets, %v, %v, %v) 0 sets",
			h.version, h.length, time.Unix(int64(h.exportTime), 0).UTC().Format(time.RFC3339),
			h.sequenceNumber, h.observationDomainID)
		got := h.String()

		if want != got {
			t.Errorf("Message.String(%v) = %v, want %v", id, got, want)
		}
	}
}

func TestMessage_UnmarshalBinary(t *testing.T) {
	tests := []struct {
		name string
		arg  []byte
		want Message
		err  error
	}{
		{
			name: "A valid Message header",
			arg: []byte{
				0x00, 0x0a, 0x01, 0x28, 0x5e, 0xe7, 0x15, 0xb1, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
			want: Message{
				0x000a, 296, 0x5ee715b1, 0, 0, nil,
			},
			err: nil,
		},
		{
			name: "An invalid Message header (unexpected EOF)",
			arg: []byte{
				0x00, 0x0a, 0x01, 0x28, 0x5e, 0xe7,
			},
			want: Message{
				0x000a, 296, 0, 0, 0, nil,
			},
			err: io.ErrUnexpectedEOF,
		},
	}
	for _, tt := range tests {
		tt := tt // pin

		var m Message

		t.Run(tt.name, func(t *testing.T) {
			if err := m.UnmarshalBinary(tt.arg); !errors.Is(err, tt.err) {
				t.Errorf("%v got error %v, want %v", tt.name, err, tt.err)
			}
			if m.version != tt.want.version ||
				m.length != tt.want.length ||
				m.exportTime != tt.want.exportTime ||
				m.sequenceNumber != tt.want.sequenceNumber ||
				m.observationDomainID != tt.want.observationDomainID ||
				len(m.sets) != len(tt.want.sets) {
				t.Errorf("%v got %v, want %v", tt.name, m, tt.want)
			}
		})
	}
}

func TestMessage_UnmarshalBinary_OneSet(t *testing.T) {
	tests := []struct {
		name string
		arg  []byte
		want Message
		err  error
	}{
		{
			name: "A valid Message header and a set's header",
			arg: []byte{
				0x00, 0x0a, 0x01, 0x28, 0x5e, 0xe7, 0x15, 0xb1, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x02, 0x01, 0x18,
			},
			want: Message{
				0x000a, 296, 0x5ee715b1, 0, 0,
				[]set{{2, 280, nil}},
			},
			err: nil,
		},
		{
			name: "A valid Message header and an invalid set's header",
			arg: []byte{
				0x00, 0x0a, 0x01, 0x28, 0x5e, 0xe7, 0x15, 0xb1, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x02, 0x01,
			},
			want: Message{
				0x000a, 296, 0x5ee715b1, 0, 0,
				[]set{{2, 0, nil}},
			},
			err: io.ErrUnexpectedEOF,
		},
	}
	for _, tt := range tests {
		tt := tt // pin

		var m Message

		t.Run(tt.name, func(t *testing.T) {
			if err := m.UnmarshalBinary(tt.arg); !errors.Is(err, tt.err) {
				t.Errorf("%v got error %v, want %v", tt.name, err, tt.err)
			}
			for i, s := range m.sets {
				if s.setID != tt.want.sets[i].setID ||
					s.length != tt.want.sets[i].length ||
					len(s.records) != len(tt.want.sets[i].records) {
					t.Errorf("UnmarshalBinary data template set[%v] %v, want %v", i, s, tt.want.sets[i])
				}
			}
		})
	}
}

func TestMessage_UnmarshalBinaryDataTemplate(t *testing.T) {
	arg := []byte{
		0x00, 0x0a, 0x01, 0x28, 0x5e, 0xe7, 0x15, 0xb1, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x02, 0x01, 0x18,
	}
	want := Message{
		version:             0x000a,
		length:              296,
		exportTime:          0x5ee715b1,
		sequenceNumber:      0,
		observationDomainID: 0,
		sets: []set{
			{
				setID:   2,
				length:  280,
				records: nil,
			},
		},
	}

	var m Message

	if err := m.UnmarshalBinary(arg); !errors.Is(err, nil) {
		t.Errorf("UnmarshalBinary data template got error %v, want %v", err, nil)
	}

	if m.version != want.version ||
		m.length != want.length ||
		m.exportTime != want.exportTime ||
		m.sequenceNumber != want.sequenceNumber ||
		m.observationDomainID != want.observationDomainID ||
		len(m.sets) != len(want.sets) {
		t.Errorf("UnmarshalBinary data template message %v, want %v", m, want)
	}

	for i, s := range m.sets {
		if s.setID != want.sets[i].setID ||
			s.length != want.sets[i].length ||
			len(s.records) != len(want.sets[i].records) {
			t.Errorf("UnmarshalBinary data template set[%v] %v, want %v", i, m, want)
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
				2, 16, make([]interface{}, 1),
			},
			want: "template set (16 octets) 1 record",
		},
		{
			name: "Options Template Set with 2 record",
			arg: set{
				3, 24, make([]interface{}, 2),
			},
			want: "options template set (24 octets) 2 records",
		},
		{
			name: "Data Set with 1 record",
			arg: set{
				258, 16, make([]interface{}, 1),
			},
			want: "data set (16 octets) 1 record",
		},
		{
			name: "Invalid Set",
			arg: set{
				24, 16, make([]interface{}, 1),
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
				258, 1, make([]fieldSpecifier, 1),
			},
			want: "template record #258 (1 field)",
		},
		{
			name: "Template Record #259 with 19 field specifiers",
			arg: templateRecord{
				259, 19, make([]fieldSpecifier, 19),
			},
			want: "template record #259 (19 fields)",
		},
		{
			name: "Template Record #260 with 0 field specifiers",
			arg: templateRecord{
				260, 0, make([]fieldSpecifier, 0),
			},
			want: "template record #260 (0 fields)",
		},
		{
			name: "Template Record #24 with 1 field specifier",
			arg: templateRecord{
				24, 1, make([]fieldSpecifier, 1),
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

func TestTemplateRecord_UnmarshalBinary(t *testing.T) {
	tests := []struct {
		name string
		arg  []byte
		want templateRecord
		err  error
	}{
		{
			name: "Template Record #258 with 35 field specifiers",
			arg:  []byte{0x01, 0x02, 0x00, 0x23},
			want: templateRecord{
				258, 35, nil,
			},
			err: nil,
		},
		{
			name: "Template Record #259 with 32 field specifiers",
			arg:  []byte{0x01, 0x03, 0x00, 0x20},
			want: templateRecord{
				259, 32, nil,
			},
			err: nil,
		},
		{
			name: "An invalid Template Record (EOF)",
			arg:  []byte{0x01, 0x04},
			want: templateRecord{
				260, 0, nil,
			},
			err: io.EOF,
		},
	}
	for _, tt := range tests {
		tt := tt // pin

		var r templateRecord

		t.Run(tt.name, func(t *testing.T) {
			if err := r.UnmarshalBinary(tt.arg); !errors.Is(err, tt.err) {
				t.Errorf("%v got error %v, want %v", tt.name, err, tt.err)
			}
			if r.templateID != tt.want.templateID ||
				r.fieldCount != tt.want.fieldCount ||
				len(r.fieldSpecifiers) != len(tt.want.fieldSpecifiers) {
				t.Errorf("%v got %v, want %v", tt.name, r, tt.want)
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
				258, 3, 2, make([]fieldSpecifier, 3),
			},
			want: "options template record #258 (scope 2 fields of 3 fields)",
		},
		{
			name: "Options Template Record #24 with 1 field specifier",
			arg: optionsTemplateRecord{
				24, 1, 1, make([]fieldSpecifier, 1),
			},
			want: "invalid options template record #24",
		},
		{
			name: "Options Template Record #384 with 0 scope fields",
			arg: optionsTemplateRecord{
				384, 1, 0, make([]fieldSpecifier, 1),
			},
			want: "invalid options template record #384",
		},
		{
			name: "Options Template Record #385 with 2 field specifier",
			arg: optionsTemplateRecord{
				385, 2, 3, make([]fieldSpecifier, 2),
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

func TestOptionsTemplateRecord_UnmarshalBinary(t *testing.T) {
	tests := []struct {
		name string
		arg  []byte
		want optionsTemplateRecord
		err  error
	}{
		{
			name: "Options Template Record #261 with 5 field specifiers and 2 scope fields",
			arg:  []byte{0x01, 0x05, 0x00, 0x05, 0x00, 0x02},
			want: optionsTemplateRecord{
				261, 5, 2, nil,
			},
			err: nil,
		},
		{
			name: "Template Record #262 with 2 field specifiers and 3 scope fields",
			arg:  []byte{0x01, 0x06, 0x00, 0x02, 0x00, 0x03},
			want: optionsTemplateRecord{
				262, 2, 3, nil,
			},
			err: nil,
		},
		{
			name: "An invalid Template Record (unexpected EOF)",
			arg:  []byte{0x01, 0x07, 0x00, 0x00, 0xff},
			want: optionsTemplateRecord{
				263, 0, 0, nil,
			},
			err: io.ErrUnexpectedEOF,
		},
	}
	for _, tt := range tests {
		tt := tt // pin

		var r optionsTemplateRecord

		t.Run(tt.name, func(t *testing.T) {
			if err := r.UnmarshalBinary(tt.arg); !errors.Is(err, tt.err) {
				t.Errorf("%v got error %v, want %v", tt.name, err, tt.err)
			}
			if r.templateID != tt.want.templateID ||
				r.fieldCount != tt.want.fieldCount ||
				r.scopeFieldCount != tt.want.scopeFieldCount ||
				len(r.fieldSpecifiers) != len(tt.want.fieldSpecifiers) {
				t.Errorf("%v got %v, want %v", tt.name, r, tt.want)
			}
		})
	}
}
