package ipfix

import (
	"errors"
	"testing"
)

func TestCorruptInputError_Error(t *testing.T) {
	tests := []struct {
		name string
		arg  CorruptInputError
		want string
	}{
		{
			name: "CorruptInputError at 0",
			arg:  CorruptInputError(0),
			want: "ipfix: illegal data at input byte 0",
		},
		{
			name: "CorruptInputError at 100",
			arg:  CorruptInputError(100),
			want: "ipfix: illegal data at input byte 100",
		},
	}
	for _, tt := range tests {
		tt := tt // pin
		t.Run(tt.name, func(t *testing.T) {
			got := tt.arg.Error()
			if got != tt.want {
				t.Errorf("%v = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestDecoder_Read(t *testing.T) {
	tests := []struct {
		name       string
		arg        []byte
		decoderOut []byte
		want       int
		err        error
	}{
		{
			name:       "Decoder Read an empty buffer with empty decoder.out",
			arg:        make([]byte, 0),
			decoderOut: make([]byte, 0),
			want:       0,
			err:        nil,
		},
		{
			name:       "Decoder Read an empty buffer with non-empty decoder.out",
			arg:        make([]byte, 0),
			decoderOut: make([]byte, 1),
			want:       0,
			err:        CorruptInputError(0),
		},
	}
	for _, tt := range tests {
		tt := tt // pin
		decoder := decoder{out: tt.decoderOut}

		t.Run(tt.name, func(t *testing.T) {
			got, err := decoder.Read(tt.arg)
			if got != tt.want {
				t.Errorf("%v = %v, want %v", tt.name, got, tt.want)
			}
			if !errors.Is(err, tt.err) {
				t.Errorf("%v error %v, want %v", tt.name, err, tt.err)
			}
		})
	}
}
