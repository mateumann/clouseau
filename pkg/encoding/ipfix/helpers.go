package ipfix

import (
	"encoding/binary"
	"fmt"
	"io"
)

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
