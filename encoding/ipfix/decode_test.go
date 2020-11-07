package ipfix

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

func TestCorruptInputError_Error(t *testing.T) {
	tests := []struct {
		name string
		arg  CorruptInputError
		want string
	}{
		{
			name: "CorruptInputError at 0 (EOF)",
			arg:  CorruptInputError{0, io.EOF},
			want: "ipfix: error in the input data at byte 0 (EOF)",
		},
		{
			name: "CorruptInputError at 100 (unexpected EOF)",
			arg:  CorruptInputError{100, io.ErrUnexpectedEOF},
			want: "ipfix: error in the input data at byte 100 (unexpected EOF)",
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

func TestCorruptInputError_Unwrap(t *testing.T) {
	tests := []struct {
		name string
		arg  CorruptInputError
		want error
	}{
		{
			name: "CorruptInputError at 0 (EOF)",
			arg:  CorruptInputError{0, io.EOF},
			want: io.EOF,
		},
		{
			name: "CorruptInputError at 100 (unexpected EOF)",
			arg:  CorruptInputError{100, io.ErrUnexpectedEOF},
			want: io.ErrUnexpectedEOF,
		},
	}
	for _, tt := range tests {
		tt := tt // pin
		t.Run(tt.name, func(t *testing.T) {
			got := tt.arg.Unwrap()
			if !errors.Is(tt.want, got) {
				t.Errorf("%v = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestDecoder_Read_FewBytes(t *testing.T) {
	tests := []struct {
		name string
		arg  io.Reader
		want int64
		err  error
	}{
		{
			name: "Decoder Read an empty buffer",
			arg:  bytes.NewReader([]byte{}),
			want: 0,
			err:  CorruptInputError{0, io.EOF},
		},
		{
			name: "Decoder Read 3-byte buffer",
			arg:  bytes.NewReader([]byte{1, 2, 3}),
			want: 2,
			err:  CorruptInputError{2, io.ErrUnexpectedEOF},
		},
	}
	for _, tt := range tests {
		tt := tt // pin
		decoder := decoder{}

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

func TestDecoder_Read_OutBuffer(t *testing.T) {
	tests := []struct {
		name       string
		arg        io.Reader
		decoderOut []byte
		want       int64
		err        error
	}{
		{
			name:       "Decoder Read an empty buffer, 2-byte decoder.out",
			arg:        bytes.NewReader([]byte{}),
			decoderOut: []byte{99, 100},
			want:       2,
			err:        CorruptInputError{2, io.EOF},
		},
		{
			name:       "Decoder Read 3-byte buffer, 1-byte decoder.out",
			arg:        bytes.NewReader([]byte{1, 2, 3}),
			decoderOut: []byte{99},
			want:       4,
			err:        nil,
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
