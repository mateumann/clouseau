package ipfix

import (
	"testing"
	"time"
)

func TestSetType_String(t *testing.T) {
	tests := []struct {
		name string
		arg  SetType
		want string
	}{
		{
			name: "DataSetType",
			arg:  DataSetType,
			want: "Data Set",
		},
		{
			name: "TemplateSetType",
			arg:  TemplateSetType,
			want: "Template Set",
		},
		{
			name: "OptionsTemplateSetType",
			arg:  OptionsTemplateSetType,
			want: "Options Template Set",
		},
		{
			name: "UnknownSetType",
			arg:  UnknownSetType,
			want: "Unknown Set",
		},
		{
			name: "Some other (unknown) SetType",
			arg:  71,
			want: "Unknown Set",
		},
		{
			name: "Some other (unknown, negative) SetType",
			arg:  -123,
			want: "Unknown Set",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.arg.String(); got != tt.want {
				t.Errorf("SetType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHeader_String(t *testing.T) {
	tests := []struct {
		name string
		arg  Header
		want string
	}{
		{
			name: "Perfectly valid Message Header",
			arg:  Header{10, 11, uint32(time.Date(2020, 2, 2, 2, 2, 2, 2, time.UTC).Unix()), 1, 0},
			want: "IPFIX Message Header (10) [domain 0; seq 1; exported at 02:02:02]",
		},
		{
			name: "Somewhat strange Message Header",
			arg:  Header{999, 9999, uint32(time.Date(2020, 2, 2, 2, 2, 2, 2, time.UTC).Unix()), 99, 99999},
			want: "IPFIX Message Header (999) [domain 99999; seq 99; exported at 02:02:02]",
		},
		{
			name: "Zero (weird) Message Header",
			arg:  Header{0, 0, 0, 0, 0},
			want: "IPFIX Message Header (0) [domain 0; seq 0; exported at 00:00:00]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.arg.String(); got != tt.want {
				t.Errorf("Header.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetHeader_Type(t *testing.T) {
	tests := []struct {
		name string
		arg  SetHeader
		want SetType
	}{
		{
			name: "Zero Set Header",
			arg:  SetHeader{0, 0},
			want: UnknownSetType,
		},
		{
			name: "Template Set Header",
			arg:  SetHeader{2, 0},
			want: TemplateSetType,
		},
		{
			name: "Options Template Set Header",
			arg:  SetHeader{3, 0},
			want: OptionsTemplateSetType,
		},
		{
			name: "Unknown Set Header",
			arg:  SetHeader{255, 0},
			want: UnknownSetType,
		},
		{
			name: "Data Set Header (id=256)",
			arg:  SetHeader{256, 0},
			want: DataSetType,
		},
		{
			name: "Another Data Set Header (id=258)",
			arg:  SetHeader{258, 0},
			want: DataSetType,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.arg.Type(); got != tt.want {
				t.Errorf("SetHeader.Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetHeader_String(t *testing.T) {
	tests := []struct {
		name string
		arg  SetHeader
		want string
	}{
		{
			name: "Zero Set Header",
			arg:  SetHeader{0, 0},
			want: "IPFIX Set Header (Unknown Set)",
		},
		{
			name: "Template Set Header",
			arg:  SetHeader{2, 0},
			want: "IPFIX Set Header (Template Set)",
		},
		{
			name: "Options Template Set Header",
			arg:  SetHeader{3, 0},
			want: "IPFIX Set Header (Options Template Set)",
		},
		{
			name: "Unknown Set Header",
			arg:  SetHeader{255, 0},
			want: "IPFIX Set Header (Unknown Set)",
		},
		{
			name: "Data Set Header (id=256)",
			arg:  SetHeader{256, 0},
			want: "IPFIX Set Header (Data Set)",
		},
		{
			name: "Another Data Set Header (id=258)",
			arg:  SetHeader{258, 0},
			want: "IPFIX Set Header (Data Set)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.arg.String(); got != tt.want {
				t.Errorf("SetHeader.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTemplateSet_Type(t *testing.T) {
	arg := TemplateSet{SetHeader{SetID: 2, Length: 22}, nil}
	want := TemplateSetType
	t.Run("Validate Template Set SetType()", func(t *testing.T) {
		if got := arg.Type(); got != want {
			t.Errorf("TemplateSet.Type() = %v, want %v", got, want)
		}
	})
}

func TestTemplateSet_Header(t *testing.T) {
	tests := []struct {
		name string
		arg  TemplateSet
		want SetHeader
	}{
		{
			name: "Template Set with a Zero Set Header",
			arg:  TemplateSet{SetHeader{0, 0}, nil},
			want: SetHeader{0, 0},
		},
		{
			name: "A valid Template Set (with proper SetHeader)",
			arg:  TemplateSet{SetHeader{2, 0}, nil},
			want: SetHeader{2, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.arg.Header(); got != tt.want {
				t.Errorf("TemplateSet.SetHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}
