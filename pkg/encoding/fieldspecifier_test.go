package encoding

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"math/big"
	"testing"
)

func TestFieldSpecifier_StringRand(t *testing.T) {
	var fs FieldSpecifier

	for attempts := 100; attempts > 0; attempts-- {
		randomized, err := rand.Int(rand.Reader, big.NewInt(65536))
		if err != nil || randomized == nil {
			t.Errorf("Error occurred during getting a random number.")

			return
		}

		id := uint16(randomized.Int64())
		want := ""
		fs = FieldSpecifier{
			InformationElementID: id,
			FieldLength:          id % 16,
			EnterpriseNumber:     uint32(randomized.Int64() % 2),
		}

		switch {
		case id >= 32768 && fs.EnterpriseNumber > 0:
			want = fmt.Sprintf("enterprise-specific field #%v (%v)", id, fs.EnterpriseNumber)
		case id < 32768 && fs.EnterpriseNumber == 0:
			want = fmt.Sprintf("IANA-registered field #%v", InformationElementByID(id))
		default:
			want = fmt.Sprintf("invalid field #%v", id)
		}

		got := fs.String()

		if want != got {
			t.Errorf("FieldSpecifier.String(%v) = %v, want %v", id, got, want)
		}
	}
}

func TestFieldSpecifier_e(t *testing.T) {
	var fs FieldSpecifier

	for attempts := 100; attempts > 0; attempts-- {
		randomized, err := rand.Int(rand.Reader, big.NewInt(65536))
		if err != nil || randomized == nil {
			t.Errorf("Error occurred during getting a random number.")

			return
		}

		id := uint16(randomized.Int64())
		want := id >= 32768
		fs = FieldSpecifier{
			InformationElementID: id,
			FieldLength:          0,
			EnterpriseNumber:     0,
		}
		got := fs.e()

		if want != got {
			t.Errorf("FieldSpecifier.e(%v) = %v, want %v", fs, got, want)
		}
	}
}

func TestFieldSpecifier_valid(t *testing.T) {
	var fs FieldSpecifier

	for attempts := 100; attempts > 0; attempts-- {
		randomized, err := rand.Int(rand.Reader, big.NewInt(65536))
		if err != nil || randomized == nil {
			t.Errorf("Error occurred during getting a random number.")

			return
		}

		id := uint16(randomized.Int64())
		if id >= 32768 {
			fs = FieldSpecifier{
				id, 1, uint32(1 + (randomized.Uint64() % 4294967295)),
			}
		} else {
			fs = FieldSpecifier{
				id, 1, 0,
			}
		}

		got := fs.valid()

		if got != true {
			t.Errorf("FieldSpecifier.valid(%v) = %v, want %v", fs, got, true)
		}
	}
}

func TestFieldSpecifier_invalid(t *testing.T) {
	var fs FieldSpecifier

	for attempts := 100; attempts > 0; attempts-- {
		randomized, err := rand.Int(rand.Reader, big.NewInt(65536))
		if err != nil || randomized == nil {
			t.Errorf("Error occurred during getting a random number.")

			return
		}

		id := uint16(randomized.Int64())
		if id >= 32768 {
			fs = FieldSpecifier{
				id, 1, 0,
			}
		} else {
			fs = FieldSpecifier{
				id, 1, uint32(1 + (randomized.Uint64() % 4294967295)),
			}
		}

		got := fs.valid()

		if got != false {
			t.Errorf("FieldSpecifier.valid(%v) = %v, want %v", fs, got, false)
		}
	}
}

func TestFieldSpecifier_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		arg  []byte
		want FieldSpecifier
		n    int
		err  error
	}{
		{
			name: "FieldSpecifier Simple",
			arg:  []byte{0x00, 0x42, 0x00, 0x01},
			want: FieldSpecifier{66, 1, 0},
			n:    4,
			err:  nil,
		},
		{
			name: "FieldSpecifier Enterprise Number",
			arg:  []byte{0xF0, 0x42, 0x00, 0x01, 0x00, 0x00, 0x01, 0x00},
			want: FieldSpecifier{61506, 1, 256},
			n:    8,
			err:  nil,
		},
		{
			name: "FieldSpecifier EOF",
			arg:  []byte{0xF0, 0x42, 0x00, 0x01, 0x00, 0x00, 0x01},
			want: FieldSpecifier{61506, 1, 0},
			n:    4,
			err:  io.ErrUnexpectedEOF,
		},
	}
	for _, tt := range tests {
		tt := tt // pin

		var s FieldSpecifier

		t.Run(tt.name, func(t *testing.T) {
			n, err := s.Unmarshal(bytes.NewReader(tt.arg))
			if !errors.Is(err, tt.err) {
				t.Errorf("%v got error %v, want %v", tt.name, err, tt.err)
			}
			if n != tt.n {
				t.Errorf("%v read bytes %v, want %v", tt.name, n, tt.n)
			}
			if s != tt.want {
				t.Errorf("%v got %v, want %v", tt.name, s, tt.want)
			}
		})
	}
}
