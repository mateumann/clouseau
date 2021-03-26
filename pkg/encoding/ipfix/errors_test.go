package ipfix

import (
	"errors"
	"io"
	"testing"
)

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

func TestCorruptInputError_Error(t *testing.T) {
	tests := []struct {
		name string
		arg  InputError
		want string
	}{
		{
			name: "CorruptInputError at 0 (EOF)",
			arg:  InputError{0, io.EOF},
			want: "corruption in input data at byte 0: EOF",
		},
		{
			name: "CorruptInputError at 100 (unexpected EOF)",
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

func TestCorruptInputError_Unwrap(t *testing.T) {
	tests := []struct {
		name string
		arg  InputError
		want error
	}{
		{
			name: "CorruptInputError at 0 (EOF)",
			arg:  InputError{0, io.EOF},
			want: io.EOF,
		},
		{
			name: "CorruptInputError at 100 (unexpected EOF)",
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
			name: "FreeLastMessagesError(2, 1)",
			arg:  IndexError{2, 1},
			want: "invalid index 2 for an object which is 1 long",
		},
		{
			name: "FreeLastMessagesError(0, 12)",
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
