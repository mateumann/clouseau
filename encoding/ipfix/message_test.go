package ipfix

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
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

func TestFieldSpecifier_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		arg  []byte
		want fieldSpecifier
		n    int
		err  error
	}{
		{
			name: "FieldSpecifier Simple",
			arg:  []byte{0x00, 0x42, 0x00, 0x01},
			want: fieldSpecifier{66, 1, 0},
			n:    4,
			err:  nil,
		},
		{
			name: "FieldSpecifier Enterprise Number",
			arg:  []byte{0xF0, 0x42, 0x00, 0x01, 0x00, 0x00, 0x01, 0x00},
			want: fieldSpecifier{61506, 1, 256},
			n:    8,
			err:  nil,
		},
		{
			name: "FieldSpecifier EOF",
			arg:  []byte{0xF0, 0x42, 0x00, 0x01, 0x00, 0x00, 0x01},
			want: fieldSpecifier{61506, 1, 0},
			n:    4,
			err:  io.ErrUnexpectedEOF,
		},
	}
	for _, tt := range tests {
		tt := tt // pin

		var s fieldSpecifier

		t.Run(tt.name, func(t *testing.T) {
			n, err := s.Unmarshal(bytes.NewReader(tt.arg))
			if !errors.Is(err, tt.err) {
				t.Errorf("%v got error %v, want %v", tt.name, err, tt.err)
			}
			if n != tt.n {
				t.Errorf("%v read bytes %v, want %v", tt.name, n, tt.n)
			}
			if s != tt.want {
				t.Errorf("%v got %v, want %v", tt.name, s, tt.want)
			}
		})
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

func TestMessage_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		arg  []byte
		want Message
		err  error
	}{
		{
			name: "A valid Message with no sets",
			arg:  []byte{0x00, 0x0a, 0x00, 0x10, 0x5e, 0xe7, 0x15, 0xb1, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			want: Message{0x000a, 16, 0x5ee715b1, 0, 0, nil},
			err:  nil,
		},
		{
			name: "An invalid Message (unexpected EOF inside a template set)",
			arg: []byte{
				0x00, 0x0a, 0x01, 0x28, 0x5e, 0xe7, 0x15, 0xb1, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x02, 0x00, 0x08, 0x02, 0x02, 0x00, 0x05, 0x0f,
			},
			want: Message{10, 296, 0x5ee715b1, 0, 0, nil},
			err:  io.ErrUnexpectedEOF,
		},
		{
			name: "An invalid Message (unexpected EOF in message header)",
			arg:  []byte{0x00, 0x0a, 0x01, 0x28, 0x5e, 0xe7, 0x15, 0xb1, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			want: Message{10, 296, 0x5ee715b1, 0, 0, nil},
			err:  io.ErrUnexpectedEOF,
		},
	}

	var templates []*templateRecord

	for _, tt := range tests {
		tt := tt // pin

		var m Message

		t.Run(tt.name, func(t *testing.T) {
			if _, err := m.Unmarshal(bytes.NewReader(tt.arg), &templates); !errors.Is(err, tt.err) {
				t.Errorf("%v got error %v (%T), want %v", tt.name, err, err, tt.err)
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

func TestMessage_Unmarshal_OneSet(t *testing.T) {
	tests := []struct {
		name string
		arg  []byte
		want Message
		n    int
		err  error
	}{
		{
			name: "A valid Message header and a set's header",
			arg: []byte{
				0x00, 0x0a, 0x00, 0x14, 0x5e, 0xe7, 0x15, 0xb1, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x02, 0x00, 0x04,
			},
			want: Message{
				0x000a, 20, 0x5ee715b1, 0, 0,
				[]set{{2, 4, nil}},
			},
			n:   20,
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
			n:   18,
			err: io.ErrUnexpectedEOF,
		},
	}

	var templates []*templateRecord

	for _, tt := range tests {
		tt := tt // pin

		var m Message

		t.Run(tt.name, func(t *testing.T) {
			n, err := m.Unmarshal(bytes.NewReader(tt.arg), &templates)
			if !errors.Is(err, tt.err) {
				t.Errorf("%v got error %v, want %v", tt.name, err, tt.err)
			}
			if n != tt.n {
				t.Errorf("%v read bytes %v, want %v", tt.name, n, tt.n)
			}
			for i, s := range m.sets {
				if s.setID != tt.want.sets[i].setID ||
					s.length != tt.want.sets[i].length ||
					len(s.records) != len(tt.want.sets[i].records) {
					t.Errorf("Unmarshal data template set[%v] %v, want %v", i, s, tt.want.sets[i])
				}
			}
		})
	}
}

func TestMessage_UnmarshalDataTemplate(t *testing.T) {
	arg := []byte{
		0x00, 0x0a, 0x00, 0x38, 0x5e, 0xe7, 0x16, 0x26, 0x00, 0x00, 0x00, 0xa8, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x02, 0x00, 0x28, 0x01, 0x02, 0x00, 0x05, 0x00, 0x3c, 0x00, 0x01, 0x00, 0x16, 0x00, 0x04,
		0x00, 0x15, 0x00, 0x04, 0x00, 0x02, 0x00, 0x04, 0x00, 0x01, 0x00, 0x04, 0x01, 0x03, 0x00, 0x02,
		0x00, 0x06, 0x00, 0x01, 0x00, 0x50, 0x00, 0x06,
	}
	want := Message{
		0x000a, 56, 0x5ee71626, 168, 0,
		[]set{
			{
				2, 40, []interface{}{
					templateRecord{258, 5, []fieldSpecifier{
						{60, 1, 0},
						{22, 4, 0},
						{21, 4, 0},
						{2, 4, 0},
						{1, 4, 0},
					}},
					templateRecord{259, 2, []fieldSpecifier{
						{6, 1, 0},
						{80, 6, 0},
					}},
				},
			},
		},
	}

	var m Message

	var templates []*templateRecord

	n, err := m.Unmarshal(bytes.NewReader(arg), &templates)

	if !errors.Is(err, nil) {
		t.Errorf("Unmarshal data template got error %v, want %v", err, nil)
	}

	if n != binary.Size(arg) {
		t.Errorf("Unmarshal data template read bytes %v, want %v", n, binary.Size(arg))
	}

	if m.version != want.version ||
		m.length != want.length ||
		m.exportTime != want.exportTime ||
		m.sequenceNumber != want.sequenceNumber ||
		m.observationDomainID != want.observationDomainID ||
		len(m.sets) != len(want.sets) {
		t.Errorf("Unmarshal data template message %v, want %v", m, want)
	}

	for i, s := range m.sets {
		if s.setID != want.sets[i].setID ||
			s.length != want.sets[i].length ||
			len(s.records) != len(want.sets[i].records) {
			t.Errorf("Unmarshal data template set[%v] %v, want %v", i, s, want.sets[i])
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

func TestSet_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		arg  []byte
		want set
		n    int
		err  error
	}{
		{
			name: "Set #3 with one optionTemplateSet",
			arg:  []byte{0x00, 0x03, 0x00, 0x0e, 0x00, 0x01, 0x00, 0x01, 0x00, 0x01, 0x00, 0x0e, 0x00, 0x01},
			want: set{3, 14, []interface{}{
				optionsTemplateRecord{1, 1, 1, []fieldSpecifier{{
					14, 1, 0,
				}}},
			}},
			n:   14,
			err: nil,
		},
		{
			name: "Set #4 (invalid set type)",
			arg:  []byte{0x00, 0x04, 0x00, 0x0e, 0x00, 0x01, 0x00, 0x01, 0x00, 0x01, 0x00, 0x0e, 0x00, 0x01},
			want: set{4, 14, nil},
			n:    4,
			err:  errInvalidSetType,
		},
	}

	var templates []*templateRecord

	for _, tt := range tests {
		tt := tt // pin

		var s set

		t.Run(tt.name, func(t *testing.T) {
			n, err := s.Unmarshal(bytes.NewReader(tt.arg), &templates)
			if !errors.Is(err, tt.err) {
				t.Errorf("%v got error %v, want %v", tt.name, err, tt.err)
			}
			if n != tt.n {
				t.Errorf("%v read bytes %v, want %v", tt.name, n, tt.n)
			}
			if s.setID != tt.want.setID ||
				s.length != tt.want.length ||
				len(s.records) != len(tt.want.records) {
				t.Errorf("%v got %v, want %v", tt.name, s, tt.want)
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

func TestTemplateRecord_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		arg  []byte
		want templateRecord
		n    int
		err  error
	}{
		{
			name: "Template Record #258 with 1 field specifiers",
			arg:  []byte{0x01, 0x02, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00},
			want: templateRecord{
				258, 1, []fieldSpecifier{{0, 0, 0}},
			},
			n:   8,
			err: nil,
		},
		{
			name: "An invalid Template Record (EOF)",
			arg:  []byte{0x01, 0x04},
			want: templateRecord{
				260, 0, nil,
			},
			n:   2,
			err: io.EOF,
		},
	}
	for _, tt := range tests {
		tt := tt // pin

		var r templateRecord

		t.Run(tt.name, func(t *testing.T) {
			n, err := r.Unmarshal(bytes.NewReader(tt.arg), nil)
			if !errors.Is(err, tt.err) {
				t.Errorf("%v got error %v, want %v", tt.name, err, tt.err)
			}
			if n != tt.n {
				t.Errorf("%v read bytes %v, want %v", tt.name, n, tt.n)
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

func TestOptionsTemplateRecord_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		arg  []byte
		want optionsTemplateRecord
		n    int
		err  error
	}{
		{
			name: "Options Template Record #261 with 3 field specifiers and 2 scope fields",
			arg: []byte{
				0x01, 0x05, 0x00, 0x03, 0x00, 0x02, 0xF0, 0xF1, 0x00, 0x01, 0x12, 0x34, 0x56, 0x78, 0xF0, 0xF2,
				0x00, 0x0F, 0x12, 0x34, 0x56, 0x79, 0xF0, 0xF3, 0x00, 0x40, 0x12, 0x34, 0x56, 0x78,
			},
			want: optionsTemplateRecord{261, 3, 2, []fieldSpecifier{
				{61681, 1, 305419896},
				{61682, 15, 305419897},
				{61683, 64, 305419896},
			}},
			n:   30,
			err: nil,
		},
		{
			name: "Template Record #262 with 2 field specifiers and 3 scope fields",
			arg:  []byte{0x01, 0x06, 0x00, 0x02, 0x00, 0x03},
			want: optionsTemplateRecord{262, 2, 3, nil},
			n:    6,
			err:  io.EOF,
		},
		{
			name: "An invalid Template Record (unexpected EOF)",
			arg:  []byte{0x01, 0x07, 0x00, 0x00, 0xff},
			want: optionsTemplateRecord{263, 0, 0, nil},
			n:    4,
			err:  io.ErrUnexpectedEOF,
		},
	}
	for _, tt := range tests {
		tt := tt // pin

		var r optionsTemplateRecord

		t.Run(tt.name, func(t *testing.T) {
			n, err := r.Unmarshal(bytes.NewReader(tt.arg), nil)
			if !errors.Is(err, tt.err) {
				t.Errorf("%v got error %v, want %v", tt.name, err, tt.err)
			}
			if n != tt.n {
				t.Errorf("%v read bytes %v, want %v", tt.name, n, tt.n)
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
