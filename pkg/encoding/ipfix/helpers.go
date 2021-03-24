package ipfix

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

var errInvalidSetType = errors.New("invalid set type")

// templateNotDefinedError is returned when no template for input data is found.
type templateNotDefinedError struct {
	templateID uint16
}

func (e templateNotDefinedError) Error() string {
	return fmt.Sprintf("template not defined %v", e.templateID)
}

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

func intDescribed(n int, singular, plural string) string {
	if n == 1 {
		return fmt.Sprintf("1 %v", singular)
	}

	return fmt.Sprintf("%v %v", n, plural)
}

func fieldDescribed(n uint16) string {
	return intDescribed(int(n), "field", "fields")
}

// unmarshalFields is a helper function that reads BigEndian binary stream into a bunch of fields.
// In case of an error returns UnmarshalError-wrapped error returned by the binary.Read.
func unmarshalFields(reader io.Reader, data ...interface{}) (n int, err error) {
	n = 0

	for _, elem := range data {
		if err = binary.Read(reader, binary.BigEndian, elem); err != nil {
			return n, &UnmarshalError{err}
		}

		elemSize := binary.Size(elem)

		// This will never happen because there is no such type that is "readable" but its size is not determined
		// if elemSize < 0 {
		// 	return n, &UnmarshalError{NotReadableError{reflect.TypeOf(elem)}}
		// }

		n += elemSize
	}

	return n, nil
}
