package ipfix

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

func TestTemplateRecord_String(t *testing.T) {
	tests := []struct {
		name string
		arg  TemplateRecord
		want string
	}{
		{
			"Template Record #258 with 1 field specifier",
			TemplateRecord{TemplateRecordBase{258, 1, make([]FieldSpecifier, 1)}},
			"template record #258 (1 field)",
		},
		{
			"Template Record #259 with 19 field specifiers",
			TemplateRecord{TemplateRecordBase{259, 19, make([]FieldSpecifier, 19)}},
			"template record #259 (19 fields)",
		},
		{
			"Template Record #260 with 0 field specifiers",
			TemplateRecord{TemplateRecordBase{260, 0, make([]FieldSpecifier, 0)}},
			"template record #260 (0 fields)",
		},
		{
			"Template Record #24 with 1 field specifier",
			TemplateRecord{TemplateRecordBase{24, 1, make([]FieldSpecifier, 1)}},
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
		arg  TemplateRecord
		want int
	}{
		{
			"Template Record #258 with 1 field specifier",
			TemplateRecord{TemplateRecordBase{258, 1, make([]FieldSpecifier, 1)}},
			4 + 8,
		},
		{
			"Template Record #259 with 19 field specifiers",
			TemplateRecord{TemplateRecordBase{259, 19, make([]FieldSpecifier, 19)}},
			4 + 19*8,
		},
		{
			"Template Record #260 with 0 field specifiers",
			TemplateRecord{TemplateRecordBase{260, 0, make([]FieldSpecifier, 0)}},
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
		want TemplateRecord
		n    int
		err  error
	}{
		{
			"Template Record #258 with 1 field specifiers",
			[]byte{0x01, 0x02, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00},
			TemplateRecord{TemplateRecordBase{258, 1, []FieldSpecifier{{0, 0, 0}}}},
			8, nil,
		},
		{
			"An invalid Template Record (EOF)",
			[]byte{0x01, 0x04},
			TemplateRecord{TemplateRecordBase{260, 0, nil}},
			2, io.EOF,
		},
	}
	for _, tt := range tests {
		tt := tt // pin

		var r TemplateRecord

		t.Run(tt.name, func(t *testing.T) {
			n, err := r.Unmarshal(bytes.NewReader(tt.arg), nil)
			if !errors.Is(err, tt.err) {
				t.Errorf("%v got error %v, want %v", tt.name, err, tt.err)
			}
			if n != tt.n {
				t.Errorf("%v read bytes %v, want %v", tt.name, n, tt.n)
			}
			if r.TemplateID != tt.want.TemplateID ||
				r.FieldCount != tt.want.FieldCount ||
				len(r.FieldSpecifiers) != len(tt.want.FieldSpecifiers) {
				t.Errorf("%v got %v, want %v", tt.name, r, tt.want)
			}
		})
	}
}

func TestOptionsTemplateRecord_String(t *testing.T) {
	tests := []struct {
		name string
		arg  OptionsTemplateRecord
		want string
	}{
		{
			"Options Template Record #258 with 3 field specifier",
			OptionsTemplateRecord{TemplateRecordBase{258, 3, make([]FieldSpecifier, 3)}, 2},
			"options template record #258 (scope 2 fields of 3 fields)",
		},
		{
			"Options Template Record #24 with 1 field specifier",
			OptionsTemplateRecord{TemplateRecordBase{24, 1, make([]FieldSpecifier, 1)}, 1},
			"invalid options template record #24",
		},
		{
			"Options Template Record #384 with 0 scope fields",
			OptionsTemplateRecord{TemplateRecordBase{384, 1, make([]FieldSpecifier, 1)}, 0},
			"invalid options template record #384",
		},
		{
			"Options Template Record #385 with 2 field specifier",
			OptionsTemplateRecord{TemplateRecordBase{385, 2, make([]FieldSpecifier, 2)}, 3},
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
		arg  OptionsTemplateRecord
		want int
	}{
		{
			"Options Template Record #258 with 3 field specifier",
			OptionsTemplateRecord{TemplateRecordBase{258, 3, make([]FieldSpecifier, 3)}, 2},
			6 + 3*8,
		},
		{
			"Options Template Record #24 with 1 field specifier",
			OptionsTemplateRecord{TemplateRecordBase{24, 1, make([]FieldSpecifier, 1)}, 1},
			6 + 8,
		},
		{
			"Options Template Record #384 with 0 scope fields",
			OptionsTemplateRecord{TemplateRecordBase{384, 1, make([]FieldSpecifier, 1)}, 0},
			6 + 8,
		},
		{
			"Options Template Record #385 with 2 field specifier",
			OptionsTemplateRecord{TemplateRecordBase{385, 2, make([]FieldSpecifier, 2)}, 3},
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
		want OptionsTemplateRecord
		n    int
		err  error
	}{
		{
			name: "Options Template Record #261 with 3 field specifiers and 2 scope fields",
			arg: []byte{
				0x01, 0x05, 0x00, 0x03, 0x00, 0x02, 0xF0, 0xF1, 0x00, 0x01, 0x12, 0x34, 0x56, 0x78, 0xF0, 0xF2,
				0x00, 0x0F, 0x12, 0x34, 0x56, 0x79, 0xF0, 0xF3, 0x00, 0x40, 0x12, 0x34, 0x56, 0x78,
			},
			want: OptionsTemplateRecord{TemplateRecordBase{261, 3, []FieldSpecifier{
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
			want: OptionsTemplateRecord{TemplateRecordBase{262, 2, nil}, 3},
			n:    6,
			err:  io.EOF,
		},
		{
			name: "An invalid Template Record (unexpected EOF)",
			arg:  []byte{0x01, 0x07, 0x00, 0x00, 0xff},
			want: OptionsTemplateRecord{TemplateRecordBase{263, 0, nil}, 0},
			n:    4,
			err:  io.ErrUnexpectedEOF,
		},
	}
	for _, tt := range tests {
		tt := tt // pin

		var r OptionsTemplateRecord

		t.Run(tt.name, func(t *testing.T) {
			n, err := r.Unmarshal(bytes.NewReader(tt.arg), nil)
			if !errors.Is(err, tt.err) {
				t.Errorf("%v got error %v, want %v", tt.name, err, tt.err)
			}
			if n != tt.n {
				t.Errorf("%v read bytes %v, want %v", tt.name, n, tt.n)
			}
			if r.TemplateID != tt.want.TemplateID ||
				r.FieldCount != tt.want.FieldCount ||
				r.ScopeFieldCount != tt.want.ScopeFieldCount ||
				len(r.FieldSpecifiers) != len(tt.want.FieldSpecifiers) {
				t.Errorf("%v got %v, want %v", tt.name, r, tt.want)
			}
		})
	}
}
