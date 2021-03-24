package ipfix

import (
	"bytes"
	"errors"
	mbits "math/bits"
	"reflect"
	"testing"
)

func TestDataRecord_String(t *testing.T) {
	tests := []struct {
		name string
		arg  DataRecord
		want string
	}{
		{
			name: "data record with 3 field values",
			arg: DataRecord{
				FieldValues: []interface{}{
					uint8(6),
					uint32(0x18fa5d76),
					uint32(0x18fa5d76),
				},
			},
			want: "ipfix data record (3 fields)",
		},
		{
			name: "data record with 0 field values",
			arg: DataRecord{
				FieldValues: []interface{}{},
			},
			want: "ipfix data record (0 fields)",
		},
		{
			name: "data record with 1 field values",
			arg: DataRecord{
				FieldValues: []interface{}{
					uint32(0x18fa5d76),
				},
			},
			want: "ipfix data record (1 field)",
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
		arg  DataRecord
		want int
	}{
		{
			"Data Record with 3 integers",
			DataRecord{[]interface{}{1, 2, 3}},
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

func TestDataRecord_Unmarshal(t *testing.T) {
	templates := []*TemplateRecord{
		{TemplateRecordBase{
			258, 1,
			[]FieldSpecifier{{0x3c, 1, 0}},
		}},
		{TemplateRecordBase{
			999, 1,
			[]FieldSpecifier{{0x00, 1, 0}},
		}},
	}

	tests := []struct {
		name     string
		arg      []byte
		template *TemplateRecord
		want     DataRecord
		n        int
		err      error
	}{
		{
			"single byte (ipVersion)",
			[]byte{0x06},
			templates[0],
			DataRecord{FieldValues: []interface{}{uint8(6)}},
			1, nil,
		},
		{
			"too much data for a template (ipVersion)",
			[]byte{0x06, 0x00},
			templates[0],
			DataRecord{FieldValues: []interface{}{uint8(6)}},
			1, nil,
		},
		{"nil template", []byte{0x06, 0x00}, nil, DataRecord{}, 0, templateNotDefinedError{}},
	}
	for _, tt := range tests {
		tt := tt // pin

		var r DataRecord

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
	templates := []*TemplateRecord{
		{TemplateRecordBase{
			258, 1,
			[]FieldSpecifier{{0x3c, 1, 0}},
		}},
		{TemplateRecordBase{
			259, 7,
			[]FieldSpecifier{
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
		template *TemplateRecord
		want     DataRecord
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
			want: DataRecord{FieldValues: []interface{}{
				uint8(6), uint32(0x18fa5d76),
				uint32(0x18fa5d76), uint32(1), uint32(72), uint16(0), uint16(0),
			}},
			n:   1,
			err: nil,
		},
	}
	for _, tt := range tests {
		tt := tt // pin

		var r DataRecord

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
