package ipfix

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"reflect"
)

// A Data Record is a record that contains values of the parameters
//      corresponding to a Template Record.

// DataRecord consist only of one or more fieldValues.
// https://tools.ietf.org/html/rfc7011#section-3.4.3
type DataRecord struct {
	FieldValues []interface{} `rethinkdb:"field_values"`
}

func (r *DataRecord) String() string {
	return fmt.Sprintf("ipfix data record (%v)", fieldDescribed(uint16(len(r.FieldValues))))
}

// Size returns a binary size of the dataRecord.
func (r DataRecord) Size() int {
	n := 0
	for _, v := range r.FieldValues {
		n += int(reflect.ValueOf(v).Type().Size())
	}

	return n
}

// Unmarshal unmarshal a IPFIX dataRecord.
func (r *DataRecord) Unmarshal(reader io.Reader, template *TemplateRecord) (n int, err error) {
	if template == nil {
		return 0, templateNotDefinedError{0}
	}

	n = 0

	for _, f := range template.FieldSpecifiers {
		e := InformationElementByID(f.InformationElementID)
		buf := make([]byte, f.FieldLength)

		if err = binary.Read(reader, binary.BigEndian, &buf); err != nil {
			return n, &UnmarshalError{err}
		}

		n += int(f.FieldLength)

		elem := createValue(e, buf)

		r.FieldValues = append(r.FieldValues, elem)
	}

	return n, nil
}

// nolint:exhaustive,gomnd
func createValue(e InformationElement, data []byte) interface{} {
	switch e.dataType {
	case reflect.Bool:
		return data[0] != 0
	case reflect.Int32:
		return int32(binary.BigEndian.Uint32(fillUp(4, data)))
	case reflect.Uint8:
		return data[0]
	case reflect.Uint16:
		return binary.BigEndian.Uint16(fillUp(2, data))
	case reflect.Uint32:
		return binary.BigEndian.Uint32(fillUp(4, data))
	case reflect.Uint64:
		return binary.BigEndian.Uint64(fillUp(8, data))
	case reflect.Float64:
		return math.Float64frombits(binary.BigEndian.Uint64(fillUp(8, data)))
	case reflect.String:
		return string(data)
	default:
		return data
	}
}

func fillUp(n int, data []byte) []byte {
	if n <= len(data) {
		return data
	}

	buf := make([]byte, n-len(data))

	return append(buf, data...)
}
