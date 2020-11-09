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

func TestDecoder_Write_FewBytes(t *testing.T) {
	tests := []struct {
		name string
		arg  []byte
		want int
		err  error
	}{
		{
			name: "Decoder Write an empty buffer",
			arg:  []byte{},
			want: 0,
			err:  nil,
		},
		{
			name: "Decoder Write 1-byte buffer",
			arg:  []byte{0x00},
			want: 0,
			err:  CorruptInputError{0, io.ErrUnexpectedEOF},
		},
		{
			name: "Decoder Write 3-byte buffer",
			arg:  []byte{0x00, 0x0a, 0x01},
			want: 2,
			err:  CorruptInputError{2, io.ErrUnexpectedEOF},
		},
		{
			name: "Decoder Write 8-byte buffer",
			arg:  []byte{0x00, 0x0a, 0x01, 0x28, 0x5e, 0xe7, 0x16, 0x26},
			want: 8,
			err:  nil,
		},
	}
	for _, tt := range tests {
		tt := tt // pin
		decoder := decoder{}

		t.Run(tt.name, func(t *testing.T) {
			got, err := decoder.Write(tt.arg)
			if got != tt.want {
				t.Errorf("%v = %v, want %v", tt.name, got, tt.want)
			}
			if !errors.Is(err, tt.err) {
				t.Errorf("%v error %v, want %v", tt.name, err, tt.err)
			}
		})
	}
}

func TestDecoder_Write_Buffer(t *testing.T) {
	tests := []struct {
		name             string
		arg              []byte
		decoderBufBefore []byte
		decoderBufAfter  []byte
		want             int
		err              error
	}{
		{
			name:             "Decoder Write an empty buffer, 2-byte decoder.buf",
			arg:              []byte{},
			decoderBufBefore: []byte{0x00, 0x0a},
			decoderBufAfter:  []byte{},
			want:             2,
			err:              nil,
		},
		{
			name:             "Decoder Write 3-byte buffer, 1-byte decoder.buf",
			arg:              []byte{0x0a, 0x01, 0x28},
			decoderBufBefore: []byte{0x00},
			decoderBufAfter:  []byte{},
			want:             4,
			err:              nil,
		},
		{
			name:             "Decoder Write 3-byte buffer, 3-byte decoder.buf",
			arg:              []byte{0x28, 0x5e, 0xe7},
			decoderBufBefore: []byte{0x00, 0x0a, 0x01},
			decoderBufAfter:  []byte{0x5e, 0xe7},
			want:             4,
			err:              CorruptInputError{4, io.ErrUnexpectedEOF},
		},
	}
	for _, tt := range tests {
		tt := tt // pin
		decoder := decoder{buf: tt.decoderBufBefore}

		t.Run(tt.name, func(t *testing.T) {
			got, err := decoder.Write(tt.arg)
			if got != tt.want {
				t.Errorf("%v = %v, want %v", tt.name, got, tt.want)
			}
			if !bytes.Equal(decoder.buf, tt.decoderBufAfter) {
				t.Errorf("%v decoder's buf afger %v, want %v", tt.name, decoder.buf, tt.decoderBufAfter)
			}
			if !errors.Is(err, tt.err) {
				t.Errorf("%v error %v, want %v", tt.name, err, tt.err)
			}
		})
	}
}
