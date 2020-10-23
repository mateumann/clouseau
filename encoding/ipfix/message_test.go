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
			want: "Data set",
		},
		{
			name: "TemplateSetType",
			arg:  TemplateSetType,
			want: "Template set",
		},
		{
			name: "OptionsTemplateSetType",
			arg:  OptionsTemplateSetType,
			want: "Options Template set",
		},
		{
			name: "UnknownSetType",
			arg:  UnknownSetType,
			want: "Unknown set",
		},
		{
			name: "Some other (unknown) SetType",
			arg:  71,
			want: "Unknown set",
		},
		{
			name: "Some other (unknown, negative) SetType",
			arg:  -123,
			want: "Unknown set",
		},
	}
	for _, tt := range tests {
		tt := tt // pin!
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
		arg  messageHeader
		want string
	}{
		{
			name: "Perfectly valid Message Header",
			arg:  messageHeader{10, 11, uint32(time.Date(2020, 2, 2, 2, 2, 2, 2, time.UTC).Unix()), 1, 0},
			want: "IPFIX Message Header (10) [domain 0; seq 1; exported at 02:02:02]",
		},
		{
			name: "Somewhat strange Message Header",
			arg:  messageHeader{999, 9999, uint32(time.Date(2020, 2, 2, 2, 2, 2, 2, time.UTC).Unix()), 99, 99999},
			want: "IPFIX Message Header (999) [domain 99999; seq 99; exported at 02:02:02]",
		},
		{
			name: "Zero (weird) Message Header",
			arg:  messageHeader{0, 0, 0, 0, 0},
			want: "IPFIX Message Header (0) [domain 0; seq 0; exported at 00:00:00]",
		},
	}
	for _, tt := range tests {
		tt := tt // pin!
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.arg.String(); got != tt.want {
				t.Errorf("messageHeader.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetHeader_Type(t *testing.T) {
	tests := []struct {
		name string
		arg  setHeader
		want SetType
	}{
		{
			name: "Zero set messageHeader",
			arg:  setHeader{0, 0},
			want: UnknownSetType,
		},
		{
			name: "Template set messageHeader",
			arg:  setHeader{2, 0},
			want: TemplateSetType,
		},
		{
			name: "Options Template set messageHeader",
			arg:  setHeader{3, 0},
			want: OptionsTemplateSetType,
		},
		{
			name: "Unknown set messageHeader",
			arg:  setHeader{255, 0},
			want: UnknownSetType,
		},
		{
			name: "Data set messageHeader (id=256)",
			arg:  setHeader{256, 0},
			want: DataSetType,
		},
		{
			name: "Another Data set messageHeader (id=258)",
			arg:  setHeader{258, 0},
			want: DataSetType,
		},
	}
	for _, tt := range tests {
		tt := tt // pin!
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.arg.Type(); got != tt.want {
				t.Errorf("header.Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetHeader_String(t *testing.T) {
	tests := []struct {
		name string
		arg  setHeader
		want string
	}{
		{
			name: "Zero set messageHeader",
			arg:  setHeader{0, 0},
			want: "IPFIX Unknown set Header (id=0, length=0)",
		},
		{
			name: "Template set messageHeader",
			arg:  setHeader{2, 0},
			want: "IPFIX Template set Header (id=2, length=0)",
		},
		{
			name: "Options Template set messageHeader",
			arg:  setHeader{3, 0},
			want: "IPFIX Options Template set Header (id=3, length=0)",
		},
		{
			name: "Unknown set messageHeader",
			arg:  setHeader{255, 0},
			want: "IPFIX Unknown set Header (id=255, length=0)",
		},
		{
			name: "Data set messageHeader (id=256)",
			arg:  setHeader{256, 0},
			want: "IPFIX Data set Header (id=256, length=0)",
		},
		{
			name: "Another Data set messageHeader (id=258)",
			arg:  setHeader{258, 0},
			want: "IPFIX Data set Header (id=258, length=0)",
		},
	}
	for _, tt := range tests {
		tt := tt // pin!
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.arg.String(); got != tt.want {
				t.Errorf("header.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTemplateSet_Type(t *testing.T) {
	arg := templateSet{setHeader{SetID: 2, Length: 22}, nil}
	want := TemplateSetType

	t.Run("Validate Template set SetType()", func(t *testing.T) {
		if got := arg.Type(); got != want {
			t.Errorf("templateSet.Type() = %v, want %v", got, want)
		}
	})
}

func TestTemplateSet_Header(t *testing.T) {
	tests := []struct {
		name string
		arg  templateSet
		want setHeader
	}{
		{
			name: "Template set with a Zero set messageHeader",
			arg:  templateSet{setHeader{0, 0}, nil},
			want: setHeader{0, 0},
		},
		{
			name: "A valid Template set (with proper header)",
			arg:  templateSet{setHeader{2, 0}, nil},
			want: setHeader{2, 0},
		},
	}
	for _, tt := range tests {
		tt := tt // pin!
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.arg.Header(); got != tt.want {
				t.Errorf("templateSet.header() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDataSet_Type(t *testing.T) {
	arg := dataSet{setHeader{SetID: 2, Length: 22}, nil}
	want := DataSetType

	t.Run("Validate Data set SetType()", func(t *testing.T) {
		if got := arg.Type(); got != want {
			t.Errorf("dataSet.Type() = %v, want %v", got, want)
		}
	})
}

func TestDataSet_Header(t *testing.T) {
	tests := []struct {
		name string
		arg  dataSet
		want setHeader
	}{
		{
			name: "Data set with a Zero set messageHeader",
			arg:  dataSet{setHeader{0, 0}, nil},
			want: setHeader{0, 0},
		},
		{
			name: "A valid Data set (with proper header)",
			arg:  dataSet{setHeader{259, 0}, nil},
			want: setHeader{259, 0},
		},
	}
	for _, tt := range tests {
		tt := tt // pin!
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.arg.Header(); got != tt.want {
				t.Errorf("dataSet.header() = %v, want %v", got, tt.want)
			}
		})
	}
}
