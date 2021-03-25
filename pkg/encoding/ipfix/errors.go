package ipfix

import (
	"fmt"
)

// var errInvalidSetType = errors.New("invalid set type")

// UnmarshalError describes an error that occurred when data are unmarshalled.
// It wraps an error returned by the binary.Read function.
type UnmarshalError struct {
	err error
}

func (e UnmarshalError) Error() string {
	return fmt.Sprintf("ipfix unmarshal error: %v", e.err)
}

// Unwrap returns the underlying error.
func (e UnmarshalError) Unwrap() error {
	return e.err
}

// CorruptInputError describes an error in input data.
type CorruptInputError struct {
	position int
	err      error
}

func (e CorruptInputError) Error() string {
	return fmt.Sprintf("ipfix error in the input data at byte %v (%v)", e.position, e.err)
}

// Unwrap returns the underlying error.
func (e CorruptInputError) Unwrap() error {
	return e.err
}

// FreeLastMessagesError describes an out of bounds error that occurred when trying to remove last `requestedIndex`
// messages from a decoder.
type FreeLastMessagesError struct {
	requestedIndex int
	messagesLength int
}

func (e FreeLastMessagesError) Error() string {
	return fmt.Sprintf("ipfix free last message error: wanted to removed %d messages, %d exist",
		e.requestedIndex, e.messagesLength)
}

// invalidSetIDError is returned when there is an invalid Set.SetID in input data, which makes it
// impossible to unmarshal.
type invalidSetIDError struct {
	setID uint16
}

func (e invalidSetIDError) Error() string {
	return fmt.Sprintf("ipfix invalid set id: %d", e.setID)
}

// templateNotDefinedError is returned when no template for input data is found.
type templateNotDefinedError struct {
	templateID uint16
}

func (e templateNotDefinedError) Error() string {
	return fmt.Sprintf("ipfix template not defined %v", e.templateID)
}
