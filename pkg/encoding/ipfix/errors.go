package ipfix

import (
	"fmt"
)

// IndexError describes an out of bounds error that occurred when trying to remove last `requestedIndex`
// messages from a decoder.
type IndexError struct {
	index int
	size  int
}

func (e IndexError) Error() string {
	return fmt.Sprintf("invalid index %d for an object which is %d long",
		e.index, e.size)
}

// InputError describes an error in input data.
type InputError struct {
	position int
	err      error
}

func (e InputError) Error() string {
	return fmt.Sprintf("corruption in input data at byte %v: %v", e.position, e.err)
}

// Unwrap returns the underlying error.
func (e InputError) Unwrap() error {
	return e.err
}

// invalidSetIDError is returned when there is an invalid Set.SetID in input data, which makes it
// impossible to unmarshal.
type invalidSetIDError struct {
	setID uint16
}

func (e invalidSetIDError) Error() string {
	return fmt.Sprintf("invalid set id: %d", e.setID)
}

// templateNotDefinedError is returned when no template for input data is found.
type templateNotDefinedError struct {
	templateID uint16
}

func (e templateNotDefinedError) Error() string {
	return fmt.Sprintf("template not defined: %d", e.templateID)
}

// unmarshalError is an error which is returned when there is an unmarshal error.
// For most part it wraps an error returned by the binary.Read function.
type unmarshalError struct {
	err error
}

func (e unmarshalError) Error() string {
	return fmt.Sprintf("unmarshal error: %v", e.err)
}

// Unwrap returns the underlying error.
func (e unmarshalError) Unwrap() error {
	return e.err
}
