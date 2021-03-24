package ipfix

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math/big"
	"testing"
	"time"
)

func TestMessage_String(t *testing.T) {
	tests := []struct {
		name string
		arg  Message
		want string
	}{
		{
			name: "A Message with one template set",
			arg: Message{
				"udp:127.0.0.1:4739", 0x000a, 296, 0x5ee715b1, time.Unix(0x5ee715b1, 0).UTC(), 0, 0,
				[]Set{{2, 280, make([]Record, 2), make([]byte, 10)}},
			},
			want: "ipfix message 10 (296 octets, 2020-06-15T06:31:13Z, 0, 0) 1 set",
		},
		{
			name: "A Message with one data set",
			arg: Message{
				"udp:127.0.0.1:4739", 0x000a, 412, 0x5ee7163d, time.Unix(0x5ee7163d, 0).UTC(), 195, 0,
				[]Set{{259, 396, make([]Record, 1), make([]byte, 10)}},
			},
			want: "ipfix message 10 (412 octets, 2020-06-15T06:33:33Z, 195, 0) 1 set",
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
			"udp:127.0.0.1:4739", 0x000a, uint16(id % 16), id, time.Unix(int64(id), 0).UTC(), 1, 1, nil,
		}

		want := fmt.Sprintf("ipfix message %v (%s, %v, %v, %v) 0 sets",
			h.Version, intDescribed(int(h.Length), "octet", "octets"),
			time.Unix(int64(h.exportTime), 0).UTC().Format(time.RFC3339), h.SequenceNumber, h.ObservationDomainID)
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
				"udp:127.0.0.1:4739", 0x000a, 16, 0x5ee715b1, time.Unix(0x5ee715b1, 0).UTC(), 0, 0, nil,
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
				"udp:127.0.0.1:4739", 10, 296, 0x5ee715b1, time.Unix(0x5ee715b1, 0).UTC(), 0, 0, nil,
			},
			err: io.ErrUnexpectedEOF,
		},
		{
			name: "An invalid Message (unexpected EOF in message header)",
			arg:  []byte{0x00, 0x0a, 0x01, 0x28, 0x5e, 0xe7, 0x15, 0xb1, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			want: Message{"udp:127.0.0.1:4739", 10, 296, 0x5ee715b1, time.Unix(0x5ee715b1, 0).UTC(), 0, 0, nil},
			err:  io.ErrUnexpectedEOF,
		},
	}

	var templates map[uint16]*TemplateRecord

	for _, tt := range tests {
		tt := tt // pin

		var m Message

		t.Run(tt.name, func(t *testing.T) {
			if _, err := m.Unmarshal(bytes.NewReader(tt.arg), "udp:127.0.0.1:4739",
				&templates); !errors.Is(err, tt.err) {
				t.Errorf("%v got error %v (%T), want %v", tt.name, err, err, tt.err)
			}
			if m.Version != tt.want.Version ||
				m.Length != tt.want.Length ||
				m.ExportedAt != tt.want.ExportedAt ||
				m.SequenceNumber != tt.want.SequenceNumber ||
				m.ObservationDomainID != tt.want.ObservationDomainID ||
				len(m.Sets) != len(tt.want.Sets) {
				t.Errorf("%v got %v, want %v", tt.name, m, tt.want)
			}
		})
	}
}

//nolint: funlen,gocognit
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
				"", 0x000a, 20, 0x5ee715b1, time.Unix(0x5ee715b1, 0).UTC(), 0, 0,
				[]Set{{2, 4, nil, nil}},
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
				"", 0x000a, 296, 0x5ee715b1, time.Unix(0x5ee715b1, 0).UTC(), 0, 0,
				[]Set{{2, 0, nil, nil}},
			},
			n:   18,
			err: io.ErrUnexpectedEOF,
		},
		{
			name: "A Message with data record with unknown template id",
			arg: []byte{
				0x00, 0x0a, 0x00, 0x17, 0x5e, 0xe7, 0x15, 0xb1, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x01, 0x02, 0x00, 0x07, 0x06, 0x0f, 0x11,
			},
			want: Message{
				"", 0x000a, 23, 0x5ee715b1, time.Unix(0x5ee715b1, 0).UTC(), 0, 0,
				[]Set{{258, 7, nil, []byte{6, 15, 17}}},
			},
			n:   23,
			err: nil,
		},
	}

	var templates map[uint16]*TemplateRecord

	for _, tt := range tests {
		tt := tt // pin

		var m Message

		t.Run(tt.name, func(t *testing.T) {
			n, err := m.Unmarshal(bytes.NewReader(tt.arg), "", &templates)
			if !errors.Is(err, tt.err) {
				t.Errorf("%v got error %v, want %v", tt.name, err, tt.err)
			}
			if n != tt.n {
				t.Errorf("%v read bytes %v, want %v", tt.name, n, tt.n)
			}
			for i, s := range m.Sets {
				if s.SetID != tt.want.Sets[i].SetID ||
					s.Length != tt.want.Sets[i].Length ||
					len(s.Records) != len(tt.want.Sets[i].Records) {
					t.Errorf("%v got set[%v] %v, want %v", tt.name, i, s, tt.want.Sets[i])
				}
				if len(s.RawData) != len(tt.want.Sets[i].RawData) {
					t.Errorf("%v got %d items of raw data, want %v", tt.name,
						len(s.RawData), len(tt.want.Sets[i].RawData))
				}
				for j, r := range s.RawData {
					if r != tt.arg[len(tt.arg)-len(s.RawData)+j] {
						t.Errorf("%v got raw data %v at %d, want %v",
							tt.name, r, j, tt.arg[len(tt.arg)-len(s.RawData)+j])
					}
				}
			}
		})
	}
}

func TestMessage_UnmarshalDataTemplate(t *testing.T) {
	arg := []byte{
		0x00, 0x0a, 0x00, 0x38, 0x5e, 0xe7, 0x16, 0x26, 0x00, 0x00, 0x00, 0xa8, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x02, 0x00, 0x28, 0x01, 0x02, 0x00, 0x05, 0x00, 0x3c, 0x00, 0x01, 0x00, 0x16, 0x00, 0x04,
		0x00, 0x15, 0x00, 0x04, 0x00, 0x02, 0x00, 0x04, 0x00, 0x01, 0x00, 0x04, 0x01, 0x03, 0x00, 0x02,
		0x00, 0x06, 0x00, 0x01, 0x00, 0x50, 0x00, 0x06,
	}
	want := Message{
		"udp:127.0.0.1:4739", 0x000a, 56, 0x5ee71626, time.Unix(0x5ee71626, 0).UTC(), 168, 0,
		[]Set{
			{
				2, 40, []Record{
					&TemplateRecord{TemplateRecordBase{258, 5, []FieldSpecifier{
						{60, 1, 0},
						{22, 4, 0},
						{21, 4, 0},
						{2, 4, 0},
						{1, 4, 0},
					}}},
					&TemplateRecord{TemplateRecordBase{259, 2, []FieldSpecifier{
						{6, 1, 0},
						{80, 6, 0},
					}}},
				}, []byte{
					0x3c, 0x00, 0x01, 0x00, 0x16, 0x00, 0x04, 0x00, 0x15, 0x00, 0x04, 0x00, 0x02, 0x00, 0x04, 0x00,
					0x01, 0x00, 0x04, 0x01, 0x03, 0x00, 0x02, 0x00, 0x06, 0x00, 0x01, 0x00, 0x50, 0x00, 0x06,
				},
			},
		},
	}

	var m Message

	var templates map[uint16]*TemplateRecord

	n, err := m.Unmarshal(bytes.NewReader(arg), "udp:127.0.0.1:4739", &templates)

	if !errors.Is(err, nil) {
		t.Errorf("Unmarshal data template got error %v, want %v", err, nil)
	}

	if n != binary.Size(arg) {
		t.Errorf("Unmarshal data template read bytes %v, want %v", n, binary.Size(arg))
	}

	if m.Version != want.Version ||
		m.Length != want.Length ||
		m.ExportedAt != want.ExportedAt ||
		m.SequenceNumber != want.SequenceNumber ||
		m.ObservationDomainID != want.ObservationDomainID ||
		len(m.Sets) != len(want.Sets) {
		t.Errorf("Unmarshal data template message %v, want %v", m, want)
	}

	for i, s := range m.Sets {
		if s.SetID != want.Sets[i].SetID ||
			s.Length != want.Sets[i].Length ||
			len(s.Records) != len(want.Sets[i].Records) {
			t.Errorf("Unmarshal data template set[%v] %v, want %v", i, s, want.Sets[i])
		}
	}
}
