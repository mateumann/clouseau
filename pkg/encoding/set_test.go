package encoding

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

func TestSet_String(t *testing.T) {
	tests := []struct {
		name string
		arg  Set
		want string
	}{
		{
			"Template Set with 1 record",
			Set{2, 16, make([]Record, 1), make([]byte, 16)},
			"template set (16 octets) 1 record",
		},
		{
			"Options Template Set with 2 record",
			Set{3, 24, make([]Record, 2), make([]byte, 16)},
			"options template set (24 octets) 2 records",
		},
		{"Data Set with 1 record", Set{258, 16, make([]Record, 1), make([]byte, 16)}, "data set (16 octets) 1 record"},
		{"Invalid Set", Set{24, 16, make([]Record, 1), make([]byte, 16)}, "invalid set #24"},
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
		want Set
		n    int
		err  error
	}{
		{
			"Set #2 with invalid optionTemplateSet",
			[]byte{0x00, 0x03, 0x00, 0x0e, 0x00, 0x01, 0x00, 0x01, 0x00, 0x01, 0x00, 0x0e},
			Set{3, 14, []Record{}, make([]byte, 12)},
			12,
			io.EOF,
		},
		{
			"Set #3 with one optionTemplateSet",
			[]byte{0x00, 0x03, 0x00, 0x0e, 0x00, 0x01, 0x00, 0x01, 0x00, 0x01, 0x00, 0x0e, 0x00, 0x01},
			Set{3, 14, []Record{
				&OptionsTemplateRecord{TemplateRecordBase{1, 1, []FieldSpecifier{{
					14, 1, 0,
				}}}, 1},
			}, make([]byte, 14)},
			14, nil,
		},
		{
			"Set #4 (invalid set type)",
			[]byte{0x00, 0x04, 0x00, 0x0e, 0x00, 0x01, 0x00, 0x01, 0x00, 0x01, 0x00, 0x0e, 0x00, 0x01},
			Set{4, 14, nil, nil},
			4, errInvalidSetType,
		},
	}

	var templates map[uint16]*TemplateRecord

	for _, tt := range tests {
		tt := tt // pin

		var s Set

		t.Run(tt.name, func(t *testing.T) {
			n, err := s.Unmarshal(bytes.NewReader(tt.arg), &templates)
			if !errors.Is(err, tt.err) {
				t.Errorf("%v got error %v, want %v", tt.name, err, tt.err)
			}
			if n != tt.n {
				t.Errorf("%v read bytes %v, want %v", tt.name, n, tt.n)
			}
			if s.SetID != tt.want.SetID ||
				s.Length != tt.want.Length ||
				len(s.Records) != len(tt.want.Records) {
				t.Errorf("%v got %v, want %v", tt.name, s, tt.want)
			}
		})
	}
}
