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
			want: "ipfix template not defined 0",
		},
		{
			name: "TemplateNotDefinedError 234",
			arg:  templateNotDefinedError{234},
			want: "ipfix template not defined 234",
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
		arg  UnmarshalError
		want string
	}{
		{
			name: "UnmarshalError: EOF",
			arg:  UnmarshalError{io.EOF},
			want: "ipfix unmarshal error: EOF",
		},
		{
			name: "UnmarshalError: unexpected EOF",
			arg:  UnmarshalError{io.ErrUnexpectedEOF},
			want: "ipfix unmarshal error: unexpected EOF",
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
		arg  CorruptInputError
		want string
	}{
		{
			name: "CorruptInputError at 0 (EOF)",
			arg:  CorruptInputError{0, io.EOF},
			want: "ipfix error in the input data at byte 0 (EOF)",
		},
		{
			name: "CorruptInputError at 100 (unexpected EOF)",
			arg:  CorruptInputError{100, io.ErrUnexpectedEOF},
			want: "ipfix error in the input data at byte 100 (unexpected EOF)",
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

func TestFreeLastMessagesError_Error(t *testing.T) {
	tests := []struct {
		name string
		arg  FreeLastMessagesError
		want string
	}{
		{
			name: "FreeLastMessagesError(2, 1)",
			arg:  FreeLastMessagesError{2, 1},
			want: "ipfix free last message error: wanted to removed 2 messages, 1 exist",
		},
		{
			name: "FreeLastMessagesError(0, 12)",
			arg:  FreeLastMessagesError{0, 12},
			want: "ipfix free last message error: wanted to removed 0 messages, 12 exist",
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
