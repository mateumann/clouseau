package ipfix

import (
	"testing"
)

func TestSetType_String(t *testing.T) {
	tests := []struct {
		name string
		arg  SetType
		want string
	}{
		{
			name: "DataSetType",
			arg:  DataSetType,
			want: "Data Set",
		},
		{
			name: "TemplateSetType",
			arg:  TemplateSetType,
			want: "Template Set",
		},
		{
			name: "OptionsTemplateSetType",
			arg:  OptionsTemplateSetType,
			want: "Options Template Set",
		},
		{
			name: "UnknownSetType",
			arg:  UnknownSetType,
			want: "Unknown Set",
		},
		{
			name: "Some other (unknown) SetType",
			arg:  71,
			want: "Unknown Set",
		},
		{
			name: "Some other (unknown, negative) SetType",
			arg:  -123,
			want: "Unknown Set",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.arg.String(); got != tt.want {
				t.Errorf("SetType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}