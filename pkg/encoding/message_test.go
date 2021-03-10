package encoding

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math/big"
	mbits "math/bits"
	"reflect"
	"testing"
	"time"
)

func TestTemplateNotDefinedError_Error(t *testing.T) {
	tests := []struct {
		name string
		arg  templateNotDefinedError
		want string
	}{
		{
			name: "TemplateNotDefinedError 0",
			arg:  templateNotDefinedError{0},
			want: "template not defined 0",
		},
		{
			name: "TemplateNotDefinedError 234",
			arg:  templateNotDefinedError{234},
			want: "template not defined 234",
		},
	}
	for _, tt := range tests {
		tt := tt // pin
		t.Run(tt.name, func(t *testing.T) {
			got := tt.arg.Error()
			if got != tt.want {
				t.Errorf("%v got %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

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
				[]set{{2, 280, make([]record, 2), make([]byte, 10)}},
			},
			want: "Message 10 (296 octets, 2020-06-15T06:31:13Z, 0, 0) 1 set",
		},
		{
			name: "A Message with one data set",
			arg: Message{
				0x000a, 412, 0x5ee7163d, 195, 0,
				[]set{{259, 396, make([]record, 1), make([]byte, 10)}},
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
			arg: []byte{
				0x00, 0x0a, 0x00, 0x10, 0x5e, 0xe7, 0x15, 0xb1, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
			want: Message{
				0x000a, 16, 0x5ee715b1, 0, 0, nil,
			},
			err: nil,
		},
		{
			name: "An invalid Message (unexpected EOF inside a template set)",
			arg: []byte{
				0x00, 0x0a, 0x01, 0x28, 0x5e, 0xe7, 0x15, 0xb1, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x02, 0x00, 0x08, 0x02, 0x02, 0x00, 0x05, 0x0f,
			},
			want: Message{
				10, 296, 0x5ee715b1, 0, 0, nil,
			},
			err: io.ErrUnexpectedEOF,
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
			if _, err := m.Unmarshal(bytes.NewReader(tt.arg), &templates); !errors.Is(errors.Unwrap(err), tt.err) {
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
				[]set{{2, 4, nil, nil}},
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
				[]set{{2, 0, nil, nil}},
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
			if !errors.Is(errors.Unwrap(err), tt.err) {
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

//nolint: funlen
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
				2, 40, []record{
					&templateRecord{templateRecordBase{258, 5, []fieldSpecifier{
						{60, 1, 0},
						{22, 4, 0},
						{21, 4, 0},
						{2, 4, 0},
						{1, 4, 0},
					}}},
					&templateRecord{templateRecordBase{259, 2, []fieldSpecifier{
						{6, 1, 0},
						{80, 6, 0},
					}}},
				}, []byte{
					0x3c, 0x00, 0x01, 0x00, 0x16, 0x00, 0x04,
					0x00, 0x15, 0x00, 0x04, 0x00, 0x02, 0x00, 0x04, 0x00, 0x01, 0x00, 0x04, 0x01, 0x03, 0x00, 0x02,
					0x00, 0x06, 0x00, 0x01, 0x00, 0x50, 0x00, 0x06,
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
			"Template Set with 1 record",
			set{2, 16, make([]record, 1), make([]byte, 16)},
			"template set (16 octets) 1 record",
		},
		{
			"Options Template Set with 2 record",
			set{3, 24, make([]record, 2), make([]byte, 16)},
			"options template set (24 octets) 2 records",
		},
		{"Data Set with 1 record", set{258, 16, make([]record, 1), make([]byte, 16)}, "data set (16 octets) 1 record"},
		{"Invalid Set", set{24, 16, make([]record, 1), make([]byte, 16)}, "invalid set #24"},
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
			"Set #2 with invalid optionTemplateSet",
			[]byte{0x00, 0x03, 0x00, 0x0e, 0x00, 0x01, 0x00, 0x01, 0x00, 0x01, 0x00, 0x0e},
			set{3, 14, []record{}, make([]byte, 12)},
			12,
			io.EOF,
		},
		{
			"Set #3 with one optionTemplateSet",
			[]byte{0x00, 0x03, 0x00, 0x0e, 0x00, 0x01, 0x00, 0x01, 0x00, 0x01, 0x00, 0x0e, 0x00, 0x01},
			set{3, 14, []record{
				&optionsTemplateRecord{templateRecordBase{1, 1, []fieldSpecifier{{
					14, 1, 0,
				}}}, 1},
			}, make([]byte, 14)},
			14, nil,
		},
		{
			"Set #4 (invalid set type)",
			[]byte{0x00, 0x04, 0x00, 0x0e, 0x00, 0x01, 0x00, 0x01, 0x00, 0x01, 0x00, 0x0e, 0x00, 0x01},
			set{4, 14, nil, nil},
			4, errInvalidSetType,
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
			"Template Record #258 with 1 field specifier",
			templateRecord{templateRecordBase{258, 1, make([]fieldSpecifier, 1)}},
			"template record #258 (1 field)",
		},
		{
			"Template Record #259 with 19 field specifiers",
			templateRecord{templateRecordBase{259, 19, make([]fieldSpecifier, 19)}},
			"template record #259 (19 fields)",
		},
		{
			"Template Record #260 with 0 field specifiers",
			templateRecord{templateRecordBase{260, 0, make([]fieldSpecifier, 0)}},
			"template record #260 (0 fields)",
		},
		{
			"Template Record #24 with 1 field specifier",
			templateRecord{templateRecordBase{24, 1, make([]fieldSpecifier, 1)}},
			"invalid template record #24",
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

func TestTemplateRecord_Size(t *testing.T) {
	tests := []struct {
		name string
		arg  templateRecord
		want int
	}{
		{
			"Template Record #258 with 1 field specifier",
			templateRecord{templateRecordBase{258, 1, make([]fieldSpecifier, 1)}},
			4 + 8,
		},
		{
			"Template Record #259 with 19 field specifiers",
			templateRecord{templateRecordBase{259, 19, make([]fieldSpecifier, 19)}},
			4 + 19*8,
		},
		{
			"Template Record #260 with 0 field specifiers",
			templateRecord{templateRecordBase{260, 0, make([]fieldSpecifier, 0)}},
			4 + 0,
		},
	}
	for _, tt := range tests {
		tt := tt // pin
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.arg.Size(); got != tt.want {
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
			"Template Record #258 with 1 field specifiers",
			[]byte{0x01, 0x02, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00},
			templateRecord{templateRecordBase{258, 1, []fieldSpecifier{{0, 0, 0}}}},
			8, nil,
		},
		{
			"An invalid Template Record (EOF)",
			[]byte{0x01, 0x04},
			templateRecord{templateRecordBase{260, 0, nil}},
			2, io.EOF,
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
			"Options Template Record #258 with 3 field specifier",
			optionsTemplateRecord{templateRecordBase{258, 3, make([]fieldSpecifier, 3)}, 2},
			"options template record #258 (scope 2 fields of 3 fields)",
		},
		{
			"Options Template Record #24 with 1 field specifier",
			optionsTemplateRecord{templateRecordBase{24, 1, make([]fieldSpecifier, 1)}, 1},
			"invalid options template record #24",
		},
		{
			"Options Template Record #384 with 0 scope fields",
			optionsTemplateRecord{templateRecordBase{384, 1, make([]fieldSpecifier, 1)}, 0},
			"invalid options template record #384",
		},
		{
			"Options Template Record #385 with 2 field specifier",
			optionsTemplateRecord{templateRecordBase{385, 2, make([]fieldSpecifier, 2)}, 3},
			"invalid options template record #385",
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

func TestOptionsTemplateRecord_Size(t *testing.T) {
	tests := []struct {
		name string
		arg  optionsTemplateRecord
		want int
	}{
		{
			"Options Template Record #258 with 3 field specifier",
			optionsTemplateRecord{templateRecordBase{258, 3, make([]fieldSpecifier, 3)}, 2},
			6 + 3*8,
		},
		{
			"Options Template Record #24 with 1 field specifier",
			optionsTemplateRecord{templateRecordBase{24, 1, make([]fieldSpecifier, 1)}, 1},
			6 + 8,
		},
		{
			"Options Template Record #384 with 0 scope fields",
			optionsTemplateRecord{templateRecordBase{384, 1, make([]fieldSpecifier, 1)}, 0},
			6 + 8,
		},
		{
			"Options Template Record #385 with 2 field specifier",
			optionsTemplateRecord{templateRecordBase{385, 2, make([]fieldSpecifier, 2)}, 3},
			6 + 2*8,
		},
	}
	for _, tt := range tests {
		tt := tt // pin
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.arg.Size(); got != tt.want {
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
			want: optionsTemplateRecord{templateRecordBase{261, 3, []fieldSpecifier{
				{61681, 1, 305419896},
				{61682, 15, 305419897},
				{61683, 64, 305419896},
			}}, 2},
			n:   30,
			err: nil,
		},
		{
			name: "Template Record #262 with 2 field specifiers and 3 scope fields",
			arg:  []byte{0x01, 0x06, 0x00, 0x02, 0x00, 0x03},
			want: optionsTemplateRecord{templateRecordBase{262, 2, nil}, 3},
			n:    6,
			err:  io.EOF,
		},
		{
			name: "An invalid Template Record (unexpected EOF)",
			arg:  []byte{0x01, 0x07, 0x00, 0x00, 0xff},
			want: optionsTemplateRecord{templateRecordBase{263, 0, nil}, 0},
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

func TestDataRecord_String(t *testing.T) {
	tests := []struct {
		name string
		arg  dataRecord
		want string
	}{
		{
			name: "data record with 3 field values",
			arg: dataRecord{
				fieldValues: []interface{}{
					uint8(6),
					uint32(0x18fa5d76),
					uint32(0x18fa5d76),
				},
			},
			want: "data record (3 fields)",
		},
		{
			name: "data record with 0 field values",
			arg: dataRecord{
				fieldValues: []interface{}{},
			},
			want: "data record (0 fields)",
		},
		{
			name: "data record with 1 field values",
			arg: dataRecord{
				fieldValues: []interface{}{
					uint32(0x18fa5d76),
				},
			},
			want: "data record (1 field)",
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

func TestDataRecord_Size(t *testing.T) {
	tests := []struct {
		name string
		arg  dataRecord
		want int
	}{
		{
			"Data Record with 3 integers",
			dataRecord{[]interface{}{1, 2, 3}},
			3 * mbits.UintSize / 8,
		},
	}
	for _, tt := range tests {
		tt := tt // pin
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.arg.Size(); got != tt.want {
				t.Errorf("%v got %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestGetTemplate(t *testing.T) {
	template1 := templateRecord{templateRecordBase{
		templateID:      258,
		fieldCount:      1,
		fieldSpecifiers: []fieldSpecifier{{0, 0, 0}},
	}}
	template2 := templateRecord{templateRecordBase{
		templateID:      259,
		fieldCount:      1,
		fieldSpecifiers: []fieldSpecifier{{0, 0, 0}},
	}}
	templates := []*templateRecord{&template1, &template2}

	tests := []struct {
		name  string
		arg   uint16
		found bool
		want  *templateRecord
	}{
		{
			name:  "get template #258",
			arg:   258,
			want:  &template1,
			found: true,
		},
		{
			name:  "get template #259",
			arg:   259,
			want:  &template2,
			found: true,
		},
		{
			name:  "get inexistent template",
			arg:   999,
			want:  nil,
			found: false,
		},
	}
	for _, tt := range tests {
		tt := tt // pin
		t.Run(tt.name, func(t *testing.T) {
			got, err := getTemplate(tt.arg, templates)
			if got != tt.want {
				t.Errorf("%v got %v, want %v", tt.name, got, tt.want)
			}
			if err != tt.found {
				t.Errorf("%v error %v, want %v", tt.name, err, tt.found)
			}
		})
	}
}

func TestDataRecord_Unmarshal(t *testing.T) {
	templates := []*templateRecord{
		{templateRecordBase{
			258, 1,
			[]fieldSpecifier{{0x3c, 1, 0}},
		}},
		{templateRecordBase{
			999, 1,
			[]fieldSpecifier{{0x00, 1, 0}},
		}},
	}

	tests := []struct {
		name     string
		arg      []byte
		template *templateRecord
		want     dataRecord
		n        int
		err      error
	}{
		{
			"single byte (ipVersion)",
			[]byte{0x06},
			templates[0],
			dataRecord{fieldValues: []interface{}{uint8(6)}},
			1, nil,
		},
		{
			"too much data for a template (ipVersion)",
			[]byte{0x06, 0x00},
			templates[0],
			dataRecord{fieldValues: []interface{}{uint8(6)}},
			1, nil,
		},
		{"nil template", []byte{0x06, 0x00}, nil, dataRecord{}, 0, templateNotDefinedError{}},
		// {"invalid template", []byte{0x01}, templates[1], dataRecord{}, 1, errUnmarshalError},
	}
	for _, tt := range tests {
		tt := tt // pin

		var r dataRecord

		t.Run(tt.name, func(t *testing.T) {
			n, err := r.Unmarshal(bytes.NewReader(tt.arg), tt.template)
			if !errors.Is(err, tt.err) {
				t.Errorf("%v got error %v, want %v", tt.name, err, tt.err)
			}
			if n != tt.n {
				t.Errorf("%v read bytes %v, want %v", tt.name, n, tt.n)
			}
		})
	}
}

func TestDataRecord_Unmarshal_Mild(t *testing.T) {
	templates := []*templateRecord{
		{templateRecordBase{
			258, 1,
			[]fieldSpecifier{{0x3c, 1, 0}},
		}},
		{templateRecordBase{
			259, 7,
			[]fieldSpecifier{
				{0x003c, 1, 0},
				{0x0016, 4, 0},
				{0x0015, 4, 0},
				{0x0002, 4, 0},
				{0x0001, 4, 0},
				{0x0007, 2, 0},
				{0x000b, 2, 0},
			},
		}},
	}

	tests := []struct {
		name     string
		arg      []byte
		template *templateRecord
		want     dataRecord
		n        int
		err      error
	}{
		{
			name: "seven fields data set #259",
			arg: []byte{
				0x06, 0x18, 0xfa, 0x5d, 0x76, 0x18, 0xfa, 0x5d, 0x76, 0x00, 0x00,
				0x00, 0x01, 0x00, 0x00, 0x00, 0x48, 0x00, 0x00, 0x00, 0x00,
			},
			template: templates[0],
			want: dataRecord{fieldValues: []interface{}{
				uint8(6), uint32(0x18fa5d76),
				uint32(0x18fa5d76), uint32(1), uint32(72), uint16(0), uint16(0),
			}},
			n:   1,
			err: nil,
		},
	}
	for _, tt := range tests {
		tt := tt // pin

		var r dataRecord

		t.Run(tt.name, func(t *testing.T) {
			n, err := r.Unmarshal(bytes.NewReader(tt.arg), tt.template)
			if !errors.Is(err, tt.err) {
				t.Errorf("%v got error %v, want %v", tt.name, err, tt.err)
			}
			if n != tt.n {
				t.Errorf("%v read bytes %v, want %v", tt.name, n, tt.n)
			}
		})
	}
}

func TestCreateValue(t *testing.T) {
	tests := []struct {
		name string
		arg  InformationElement
		data []byte
		want interface{}
	}{
		{"hashDigestOutput (bool)", InformationElementByID(333), []byte{0x01}, true},
		{"dot1qDEI (bool)", InformationElementByID(388), []byte{0x00}, false},
		{"mibObjectValueInteger (int32)", InformationElementByID(434), []byte{0xff, 0xfe, 0xfd, 0xfc}, int32(-66052)},
		{"mibObjectValueInteger (int32)", InformationElementByID(434), []byte{0xff, 0xfe, 0xfd, 0xfc}, int32(-66052)},
		{"protocolIdentifier (uint8)", InformationElementByID(60), []byte{0x06}, uint8(6)},
		{"udpSourcePort (uint16)", InformationElementByID(180), []byte{0x01, 0x02}, uint16(258)},
		{"tcpSequenceNumber (uint32)", InformationElementByID(184), []byte{0x01, 0x02, 0x03, 0x04}, uint32(16909060)},
		{
			"exportedOctetTotalCount (uint64)", InformationElementByID(40),
			[]byte{0xff, 0xfe, 0xfd, 0x01, 0x02, 0x03},
			uint64(281470631477763),
		},
		{
			"absoluteError (float64)", InformationElementByID(320),
			[]byte{0xaa, 0x55, 0xaa, 0x55, 0xaa},
			float64(3.614493686586e-312),
		},
		{"interfaceName (string)", InformationElementByID(82), []byte{0x65, 0x74, 0x68, 0x30}, "eth0"},
		{
			"destinationMacAddress (array)", InformationElementByID(80),
			[]byte{0x5a, 0xea, 0xf0, 0xa8, 0xd3, 0x7b},
			[]byte{0x5a, 0xea, 0xf0, 0xa8, 0xd3, 0x7b},
		},
		{"layer2packetSectionData (slice)", InformationElementByID(104), []byte{0xca, 0xfe}, []byte{0xca, 0xfe}},
		// {"reserved InformationElement", InformationElementByID(0), []byte{0xca, 0xfe}, nil},
	}
	for _, tt := range tests {
		tt := tt // pin

		t.Run(tt.name, func(t *testing.T) {
			got := createValue(tt.arg, tt.data)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%v created value %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
