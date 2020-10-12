package ipfix

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestUnmarshal_OutIsNil(t *testing.T) {
	buffer := make([]byte, 60)
	templates := make([]TemplateRecord, 2)

	err := Unmarshal(buffer, &templates, nil)

	if err == nil {
		t.Errorf("Error should have not been nil: %v", err)
	}
}

func TestUnmarshal_ReadError(t *testing.T) {
	buffer := make([]byte, 60)
	templates := make([]TemplateRecord, 2)

	var message Message

	err := Unmarshal(buffer, &templates, &message)
	if err != nil {
		t.Errorf("Error during Unmarshallings: %v", err)
	}
}

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
				Header: Header{
					VersionNumber:       10,
					Length:              296,
					ExportTime:          uint32(time.Date(2020, 0o6, 15, 0o6, 31, 13, 0, time.UTC).Unix()),
					SequenceNumber:      0,
					ObservationDomainID: 0,
				},
				Sets: []Set{
					TemplateSet{
						SetHeader: SetHeader{
							SetID:  258,
							Length: 35,
						},
						Records: nil,
					},
					TemplateSet{
						SetHeader: SetHeader{
							SetID:  259,
							Length: 32,
						},
						Records: nil,
					},
				},
			},
		},
	}

	var got Message

	templates := make([]TemplateRecord, 0)

	for _, tt := range tests {
		tt := tt // pin!
		t.Run(tt.name, func(t *testing.T) {
			if tt.arg == nil {
				t.Errorf("%v: does not have valid input arguments", tt.name)

				return
			}
			_ = Unmarshal(tt.arg, &templates, &got)
			if got.Header != tt.want.Header {
				t.Errorf("%v: got %v, want %v", tt.name, got.Header, tt.want.Header)
			}
			for index, set := range got.Sets {
				if set.Header() != tt.want.Sets[index].Header() {
					t.Errorf("%v: got %v, want %v", tt.name, set.Header(), tt.want.Sets[index].Header())
				}
			}
			// TODO compare number of Tempates (Records?) in a TemplateSets to their declared length
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
