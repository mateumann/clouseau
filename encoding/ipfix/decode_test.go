package ipfix

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestUnmarshal_ValueIsNil(t *testing.T) {
	var decode Decode

	buffer := make([]byte, 60)

	err := decode.Unmarshal(buffer, nil)

	if err == nil {
		t.Errorf("Error should have not been nil: %v", err)
	}
}

func TestInvalidUnmarshalError_Error(t *testing.T) {
	tests := []struct {
		name string
		arg  interface{}
		want string
	}{
		{
			name: "Unmarshal nil",
			arg:  nil,
			want: "ipfix: Unmarshal(nil)",
		},
		{
			name: "Unmarshal integer",
			arg:  1,
			want: "ipfix: Unmarshal(non-pointer int)",
		},
		{
			name: "Unmarshal string",
			arg:  "A message",
			want: "ipfix: Unmarshal(non-pointer string)",
		},
		{
			name: "ipfix: Unmarshal ipfix.Message",
			arg: Message{
				Header: MessageHeader{},
				Sets:   nil,
			},
			want: "ipfix: Unmarshal(non-pointer ipfix.Message)",
		},
	}

	for _, tt := range tests {
		tt := tt // pin!

		t.Run(tt.name, func(t *testing.T) {
			err := &InvalidUnmarshalError{reflect.TypeOf(tt.arg)}
			if err.Error() != tt.want {
				t.Errorf("%v: got %v, want %v", tt.name, err, tt.want)
			}
		})
	}
}

func TestUnmarshal_ReadError(t *testing.T) {
	var decode Decode

	var message Message

	buffer := make([]byte, 60)

	err := decode.Unmarshal(buffer, &message)
	if err != nil {
		t.Errorf("Error during Unmarshallings: %v", err)
	}
}

//nolint:funlen,gocognit
func TestUnmarshal(t *testing.T) {
	tests := []struct {
		name string
		arg  []byte
		want Message
	}{
		{
			name: "Unmarshal a Template Set",
			arg:  readFromFile("../../testdata/stream/templates_258_259.bin"),
			want: Message{
				Header: MessageHeader{
					VersionNumber:       10,
					Length:              296,
					ExportTime:          uint32(time.Date(2020, 6, 15, 6, 31, 13, 0, time.UTC).Unix()),
					SequenceNumber:      0,
					ObservationDomainID: 0,
				},
				// missing Flow Set header: Id (uint16) = 2, Length (uint16) = 280
				Sets: []Set{
					TemplateSet{
						header: SetHeader{
							SetID:  2,
							Length: 280,
						},
						records: []TemplateRecord{
							{
								TemplateHeader: TemplateRecordHeader{
									TemplateID: 258,
									FieldCount: 35,
								},
								FieldSpecifiers: []FieldSpecifier{
									{
										ElementID:        60,
										FieldLength:      1,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        22,
										FieldLength:      4,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        21,
										FieldLength:      4,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        2,
										FieldLength:      4,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        1,
										FieldLength:      4,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        7,
										FieldLength:      2,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        11,
										FieldLength:      2,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        10,
										FieldLength:      4,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        14,
										FieldLength:      4,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        4,
										FieldLength:      1,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        5,
										FieldLength:      1,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        6,
										FieldLength:      1,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        80,
										FieldLength:      6,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        81,
										FieldLength:      6,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        56,
										FieldLength:      6,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        8,
										FieldLength:      4,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        12,
										FieldLength:      4,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        15,
										FieldLength:      4,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        9,
										FieldLength:      1,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        13,
										FieldLength:      1,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        192,
										FieldLength:      1,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        206,
										FieldLength:      1,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        189,
										FieldLength:      1,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        224,
										FieldLength:      8,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        205,
										FieldLength:      2,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        184,
										FieldLength:      4,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        185,
										FieldLength:      4,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        186,
										FieldLength:      2,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        33,
										FieldLength:      1,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        176,
										FieldLength:      1,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        177,
										FieldLength:      1,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        225,
										FieldLength:      4,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        226,
										FieldLength:      4,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        227,
										FieldLength:      2,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        228,
										FieldLength:      2,
										EnterpriseNumber: 0,
									},
								},
							},
							{
								TemplateHeader: TemplateRecordHeader{
									TemplateID: 259,
									FieldCount: 32,
								},
								FieldSpecifiers: []FieldSpecifier{
									{
										ElementID:        60,
										FieldLength:      1,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        22,
										FieldLength:      4,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        21,
										FieldLength:      4,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        2,
										FieldLength:      4,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        1,
										FieldLength:      4,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        7,
										FieldLength:      2,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        11,
										FieldLength:      2,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        10,
										FieldLength:      4,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        14,
										FieldLength:      4,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        4,
										FieldLength:      1,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        5,
										FieldLength:      1,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        6,
										FieldLength:      1,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        80,
										FieldLength:      6,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        81,
										FieldLength:      6,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        56,
										FieldLength:      6,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        27,
										FieldLength:      16,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        28,
										FieldLength:      16,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        62,
										FieldLength:      16,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        29,
										FieldLength:      1,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        30,
										FieldLength:      1,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        192,
										FieldLength:      1,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        206,
										FieldLength:      1,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        189,
										FieldLength:      1,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        224,
										FieldLength:      8,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        205,
										FieldLength:      2,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        184,
										FieldLength:      4,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        185,
										FieldLength:      4,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        186,
										FieldLength:      2,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        33,
										FieldLength:      1,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        178,
										FieldLength:      1,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        179,
										FieldLength:      1,
										EnterpriseNumber: 0,
									},
									{
										ElementID:        31,
										FieldLength:      4,
										EnterpriseNumber: 0,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Unmarshal a single Data Set with padding",
			arg:  readFromFile("../../testdata/stream/data_258_1set_padding.bin"),
			want: Message{
				Header: MessageHeader{
					VersionNumber:       10,
					Length:              124,
					ExportTime:          uint32(time.Date(2020, 6, 15, 6, 33, 23, 0, time.UTC).Unix()),
					SequenceNumber:      188,
					ObservationDomainID: 0,
				},
				Sets: []Set{
					DataSet{
						header: SetHeader{
							SetID:  258,
							Length: 108,
						},
						records: []DataRecord{},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		var decode Decode

		var got Message

		tt := tt // pin!

		t.Run(tt.name, func(t *testing.T) {
			if tt.arg == nil {
				t.Errorf("%v: does not have valid input arguments", tt.name)

				return
			}
			_ = decode.Unmarshal(tt.arg, &got)
			if got.Header != tt.want.Header {
				t.Errorf("%v: got %v, want %v", tt.name, got.Header, tt.want.Header)
			}
			if len(got.Sets) != len(tt.want.Sets) {
				t.Errorf("%v: got %v Sets, want %v Sets", tt.name, len(got.Sets), len(tt.want.Sets))
			}
			for setIndex, set := range got.Sets {
				if set.Type() != tt.want.Sets[setIndex].Type() {
					t.Errorf("%v: got %v, want %v", tt.name, set.Type(), tt.want.Sets[setIndex].Type())
				}
				if set.Header() != tt.want.Sets[setIndex].Header() {
					t.Errorf("%v: got %v, want %v", tt.name, set.Header(), tt.want.Sets[setIndex].Header())
				}
				for recordIndex, record := range set.TemplateRecords() {
					if record.TemplateHeader != tt.want.Sets[setIndex].TemplateRecords()[recordIndex].TemplateHeader {
						t.Errorf("%v: got %v, want %v", tt.name, record.TemplateHeader,
							tt.want.Sets[setIndex].TemplateRecords()[recordIndex].TemplateHeader)
					}
					for fsIndex, fieldSpecifier := range record.FieldSpecifiers {
						if fieldSpecifier != tt.want.Sets[setIndex].TemplateRecords()[recordIndex].FieldSpecifiers[fsIndex] {
							t.Errorf("%v: got %v, want %v", tt.name, fieldSpecifier,
								tt.want.Sets[setIndex].TemplateRecords()[recordIndex].FieldSpecifiers[fsIndex])
						}
					}
				}
			}
			// TODO compare number of Tempates (records?) in a TemplateSets to their declared length
		})
	}
}

func readFromFile(fileName string) []byte {
	file, err := os.Open(fileName)
	if err != nil {
		return nil
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil
	}

	return bytes
}
