package ipfix

import (
	"io"
	"testing"
)

func TestTemplateNotDefinedError_Error(t *testing.T) {
	tests := []struct {
		name string
		arg  TemplateNotDefinedError
		want string
	}{
		{
			name: "TemplateNotDefinedError 0",
			arg:  TemplateNotDefinedError{0},
			want: "ipfix template not defined 0",
		},
		{
			name: "TemplateNotDefinedError 234",
			arg:  TemplateNotDefinedError{234},
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
