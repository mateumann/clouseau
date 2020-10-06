package ipfix

import (
	"bufio"
	"io/ioutil"
	"os"
	"testing"
)

func readFromFile(fileName string) ([]byte, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return bytes, err
}

func TestUnmarshal(t *testing.T) {
	// func Unmarshal(in []byte, templates *[]TemplateRecord, out *Message) error {
	simpleTemplateSet, _ := readFromFile("../../test/stream/templates_258_259.bin")
	tests := []struct {
		name string
		arg  []byte
		want Message
	}{
		{
			name: "Unmarshal a Template Set",
			arg: simpleTemplateSet,
			want: Message{
				Header: Header{
					VersionNumber:       0,
					Length:              0,
					ExportTime:          0,
					SequenceNumber:      0,
					ObservationDomainID: 0,
				},
				Sets: []Set{TemplateSet{
					SetHeader: SetHeader{
						SetID:  0,
						Length: 0,
					},
					Records:   nil,
				}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Message
			if _ := Unmarshal(tt.arg, nil, &got); got != tt.want {
				t.Errorf("FieldSpecifier.E() = %v, want %v", got, tt.want)
			}
		})
	}
}
