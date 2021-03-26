package ipfix

import (
	"errors"
	"io"
	"testing"
)

func TestInputError_Error(t *testing.T) {
	tests := []struct {
		name string
		arg  InputError
		want string
	}{
		{
			name: "InputError at 0 (EOF)",
			arg:  InputError{0, io.EOF},
			want: "corruption in input data at byte 0: EOF",
		},
		{
			name: "InputError at 100 (unexpected EOF)",
			arg:  InputError{100, io.ErrUnexpectedEOF},
			want: "corruption in input data at byte 100: unexpected EOF",
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

func TestInputError_Unwrap(t *testing.T) {
	tests := []struct {
		name string
		arg  InputError
		want error
	}{
		{
			name: "InputError at 0 (EOF)",
			arg:  InputError{0, io.EOF},
			want: io.EOF,
		},
		{
			name: "InputError at 100 (unexpected EOF)",
			arg:  InputError{100, io.ErrUnexpectedEOF},
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

func TestInvalidIndexError_Error(t *testing.T) {
	tests := []struct {
		name string
		arg  IndexError
		want string
	}{
		{
			name: "InvalidIndexError(2, 1)",
			arg:  IndexError{2, 1},
			want: "invalid index 2 for an object which is 1 long",
		},
		{
			name: "InvalidIndexError(-1, 12)",
			arg:  IndexError{-1, 12},
			want: "invalid index -1 for an object which is 12 long",
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

func TestInvalidSetIDError_Error(t *testing.T) {
	tests := []struct {
		name string
		arg  invalidSetIDError
		want string
	}{
		{
			name: "Set ID 33",
			arg:  invalidSetIDError{33},
			want: "invalid set id: 33",
		},
		{
			name: "Set ID 195",
			arg:  invalidSetIDError{195},
			want: "invalid set id: 195",
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

func TestTemplateNotDefinedError_Error(t *testing.T) {
	tests := []struct {
		name string
		arg  templateNotDefinedError
		want string
	}{
		{
			name: "TemplateNotDefinedError 0",
			arg:  templateNotDefinedError{0},
			want: "template not defined: 0",
		},
		{
			name: "TemplateNotDefinedError 234",
			arg:  templateNotDefinedError{234},
			want: "template not defined: 234",
		},
	}
	for _, tt := range tests {
		tt := tt // pin
		t.Run(tt.name, func(t *testing.T) {
			got := tt.arg.Error()
			if got != tt.want {
				t.Errorf("%v got %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestUnmarshalError_Error(t *testing.T) {
	tests := []struct {
		name string
		arg  unmarshalError
		want string
	}{
		{
			name: "UnmarshalError: EOF",
			arg:  unmarshalError{io.EOF},
			want: "unmarshal error: EOF",
		},
		{
			name: "UnmarshalError: unexpected EOF",
			arg:  unmarshalError{io.ErrUnexpectedEOF},
			want: "unmarshal error: unexpected EOF",
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
