package ipfix

import (
	"bytes"
	"errors"
	"io"
	"net"
	"testing"
	"time"
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

func TestNewDecoder(t *testing.T) {
	tests := []struct {
		name string
		arg  net.UDPAddr
		want Decoder
	}{
		{
			name: "Decoder on localhost:4739",
			arg:  net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 4739, Zone: ""},
			want: Decoder{
				source:              net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 4739, Zone: ""},
				observationDomainID: 0,
				err:                 nil,
				buf:                 make([]byte, 2048),
			},
		},
	}
	for _, tt := range tests {
		tt := tt // pin

		t.Run(tt.name, func(t *testing.T) {
			got := NewDecoder(tt.arg)
			if !got.source.IP.Equal(tt.want.source.IP) ||
				got.source.Port != tt.want.source.Port ||
				got.source.Zone != tt.want.source.Zone ||
				got.observationDomainID != tt.want.observationDomainID {
				t.Errorf("%v = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestDecoder_Read_FewBytes(t *testing.T) {
	tests := []struct {
		name string
		arg  []byte
		want int
		err  error
	}{
		{
			name: "Decoder Read an empty buffer",
			arg:  []byte{},
			want: 0,
			err:  nil,
		},
		{
			name: "Decoder Read 1-byte buffer",
			arg:  []byte{0x00},
			want: 0,
			err:  CorruptInputError{0, io.ErrUnexpectedEOF},
		},
		{
			name: "Decoder Read 3-byte buffer",
			arg:  []byte{0x00, 0x0a, 0x01},
			want: 2,
			err:  CorruptInputError{2, io.ErrUnexpectedEOF},
		},
		{
			name: "Decoder Read 8-byte buffer",
			arg:  []byte{0x00, 0x0a, 0x01, 0x28, 0x5e, 0xe7, 0x16, 0x26},
			want: 8,
			err:  nil,
		},
	}
	for _, tt := range tests {
		tt := tt // pin
		decoder := Decoder{}

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

func TestDecoder_Read_Buffer(t *testing.T) {
	tests := []struct {
		name             string
		arg              []byte
		decoderBufBefore []byte
		decoderBufAfter  []byte
		want             int
		err              error
	}{
		{
			name:             "Decoder Read an empty buffer, 2-byte Decoder.buf",
			arg:              []byte{},
			decoderBufBefore: []byte{0x00, 0x0a},
			decoderBufAfter:  []byte{},
			want:             2,
			err:              nil,
		},
		{
			name:             "Decoder Read 3-byte buffer, 1-byte Decoder.buf",
			arg:              []byte{0x0a, 0x01, 0x28},
			decoderBufBefore: []byte{0x00},
			decoderBufAfter:  []byte{},
			want:             4,
			err:              nil,
		},
		{
			name:             "Decoder Read 3-byte buffer, 3-byte Decoder.buf",
			arg:              []byte{0x28, 0x5e, 0xe7},
			decoderBufBefore: []byte{0x00, 0x0a, 0x01},
			decoderBufAfter:  []byte{0x5e, 0xe7},
			want:             4,
			err:              CorruptInputError{4, io.ErrUnexpectedEOF},
		},
		{
			name:             "Decoder Read Message header from 8-byte buffer and 8-byte Decoder.buf",
			arg:              []byte{0x00, 0x0a, 0x00, 0x10, 0x5e, 0xe7, 0x15, 0xb1},
			decoderBufBefore: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			decoderBufAfter:  []byte{},
			want:             16,
			err:              nil,
		},
	}
	for _, tt := range tests {
		tt := tt // pin
		decoder := Decoder{buf: tt.decoderBufBefore}

		t.Run(tt.name, func(t *testing.T) {
			got, err := decoder.Read(tt.arg)
			if got != tt.want {
				t.Errorf("%v = %v, want %v", tt.name, got, tt.want)
			}
			if !bytes.Equal(decoder.buf, tt.decoderBufAfter) {
				t.Errorf("%v Decoder's buf after %v, want %v", tt.name, decoder.buf, tt.decoderBufAfter)
			}
			if !errors.Is(err, tt.err) {
				t.Errorf("%v error %v, want %v", tt.name, err, tt.err)
			}
		})
	}
}

func TestDecoder_Read_WithError(t *testing.T) {
	tests := []struct {
		name       string
		arg        []byte
		decoderErr error
		want       int
		err        error
	}{
		{
			name:       "Decoder Read with nil Decoder.err",
			arg:        []byte{0x00, 0x0a},
			decoderErr: nil,
			want:       2,
			err:        nil,
		},
		{
			name:       "Decoder Read with non-nil Decoder.err",
			arg:        []byte{0x00, 0x0a, 0x01, 0x28},
			decoderErr: CorruptInputError{0, nil},
			want:       0,
			err:        CorruptInputError{0, nil},
		},
	}
	for _, tt := range tests {
		tt := tt // pin
		decoder := Decoder{err: tt.decoderErr}

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

func TestDecoder_Read_WholeMessageHeaders(t *testing.T) {
	tests := []struct {
		name string
		arg  []byte
		want Message
		err  error
	}{
		{
			name: "Decoder Read read a Message header",
			arg: []byte{
				0x00, 0x0a, 0x01, 0x28, 0x5e, 0xe7, 0x15, 0xb1,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
			want: Message{
				version: 10,
				length:  296,
				exportTime: uint32(time.Date(2020, 6, 15,
					6, 31, 13, 0, time.UTC).Unix()),
				sequenceNumber:      0,
				observationDomainID: 0,
				sets:                nil,
			},
			err: nil,
		},
	}
	for _, tt := range tests {
		tt := tt // pin
		decoder := Decoder{}

		t.Run(tt.name, func(t *testing.T) {
			got, err := decoder.Read(tt.arg)
			if got != 16 {
				t.Errorf("%v read not exactly Message header length bytes", tt.name)
			}
			if !errors.Is(err, tt.err) {
				t.Errorf("%v error %v, want %v", tt.name, err, tt.err)
			}
			if decoder.err != nil {
				t.Errorf("%v Decoder error: %v", tt.name, decoder.err)
			}
		})
	}
}
