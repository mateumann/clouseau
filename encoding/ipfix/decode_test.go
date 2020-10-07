package ipfix

import (
	"io/ioutil"
	"os"
	"testing"
    "time"
)

func readFromFile(t* testing.T, fileName string) []byte {
	file, err := os.Open(fileName)
	if err != nil {
        t.Errorf("Unable to read test stream file %v", fileName)
        return nil
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
    if err != nil {
        t.Errorf("Unable to read test stream file %v", fileName)
        return nil
    }

	return bytes
}

func TestUnmarshal(t *testing.T) {
	tests := []struct {
		name string
		arg  []byte
		want Message
	}{
		{
			name: "Unmarshal a Template Set",
			arg: readFromFile(t, "../../test/stream/templates_258_259.bin"),
			want: Message{
				Header: Header{
                    VersionNumber: 10,
                    Length:        296,
                    ExportTime:    uint32(time.Date(2020, 06, 15, 06, 31, 13, 0, time.UTC).Unix()),
                    SequenceNumber:      0,
                    ObservationDomainID: 0,
                },
				Sets: []Set{TemplateSet{
					SetHeader: SetHeader{
						SetID:  2,
						Length: 280,
					},
					Records:   nil,
				}},
			},
		},
	}
    var got Message
    var templates = make([]TemplateRecord, 10)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		    if tt.arg == nil {
                t.Errorf("%v: does not have valid input arguments", tt.name)
                return
            }
            _ = Unmarshal(tt.arg, &templates, &got)
            if got.Header != tt.want.Header {
                t.Errorf("%v: got %v, want %v", tt.name, got.Header, tt.want.Header)
            }
            //for index, set := range got.Sets {
            //    if set.Header() != tt.want.Sets[index].Header() {
            //        t.Errorf("%v: got %v, want %v", tt.name, set.Header(), tt.want.Sets[index].Header())
            //    }
            //}
        })
	}
}
