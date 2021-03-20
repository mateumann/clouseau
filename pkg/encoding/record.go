package encoding

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"reflect"
)

// A Template is an ordered sequence of <type, length> pairs used to
//      completely specify the structure and semantics of a particular set
//      of information that needs to be communicated from an IPFIX Device
//      to a Collector.  Each Template is uniquely identifiable by means
//      of a Template ID.

// An IPFIX Message is a Message that originates at the Exporting
//      Process and carries the IPFIX records of this Exporting Process,
//      and whose destination is a Collecting Process.  An IPFIX Message
//      is encapsulated at the transport layer.

// The Message Header is the first part of an IPFIX Message; the
//      Message Header provides basic information about the Message, such
//      as the IPFIX version, length of the Message, Message sequence
//      number, etc.

// A Template Record defines the structure and interpretation of
//      fields in a Data Record.

// A Data Record is a record that contains values of the parameters
//      corresponding to a Template Record.

// An Options Template Record is a Template Record that defines the
//      structure and interpretation of fields in a Data Record, including
//      defining how to scope the applicability of the Data Record.

// A Set is a collection of records that have a similar structure,
//      prefixed by a header.  In an IPFIX Message, zero or more Sets
//      follow the Message Header.  There are three different types of
//      Sets: Template Sets, Options Template Sets, and Data Sets.

// A Template Set is a collection of one or more Template Records
//      that have been grouped together in an IPFIX Message.

// A Data Set is one or more Data Records, of the same type, that are
//      grouped together in an IPFIX Message.  Each Data Record is
//      previously defined by a Template Record or an Options Template
//      Record.

// An Information Element is a protocol- and encoding-independent
//      description of an attribute that may appear in an IPFIX Record.
//      Information Elements are defined in the IANA "IPFIX Information
//      Elements" registry [IANA-IPFIX].  The type associated with an
//      Information Element indicates constraints on what it may contain
//      and also determines the valid encoding mechanisms for use in
//      IPFIX.

// Figure A shows a summary of IPFIX terminology.
//
//    +------------------+---------------------------------------------+
//    |                  |                 Contents                    |
//    |                  +--------------------+------------------------+
//    |       Set        |      Template      |         Record         |
//    +------------------+--------------------+------------------------+
//    |     Data Set     |          /         |     Data Record(s)     |
//    +------------------+--------------------+------------------------+
//    |   Template Set   | Template Record(s) |           /            |
//    +------------------+--------------------+------------------------+
//    | Options Template |  Options Template  |           /            |
//    |       Set        |      Record(s)     |                        |
//    +------------------+--------------------+------------------------+

// The format of the IPFIX Message is shown in Figure B.
//
//         +----------------------------------------------------+
//         | Message Header                                     |
//         +----------------------------------------------------+
//         | Set                                                |
//         +----------------------------------------------------+
//         | Set                                                |
//         +----------------------------------------------------+
//           ...
//         +----------------------------------------------------+
//         | Set                                                |
//         +----------------------------------------------------+

// An IPFIX Message consisting of interleaved Template, Data, and
//      Options Template Sets, as shown in Figure C.
//
//     +--------+--------------------------------------------------------+
//     |        | +----------+ +---------+     +-----------+ +---------+ |
//     |Message | | Template | | Data    |     | Options   | | Data    | |
//     | Header | | Set      | | Set     | ... | Template  | | Set     | |
//     |        | |          | |         |     | Set       | |         | |
//     |        | +----------+ +---------+     +-----------+ +---------+ |
//     +--------+--------------------------------------------------------+

// An IPFIX Message consisting entirely of Data Sets, sent after the
//      appropriate Template Records have been defined and transmitted to
//      the Collecting Process, as shown in Figure D.
//
//       +--------+----------------------------------------------+
//       |        | +---------+     +---------+      +---------+ |
//       |Message | | Data    |     | Data    |      | Data    | |
//       | Header | | Set     | ... | Set     | ...  | Set     | |
//       |        | +---------+     +---------+      +---------+ |
//       +--------+----------------------------------------------+

// An IPFIX Message consisting entirely of Template and Options
//      Template Sets, as shown in Figure E.  Such a Message can be used
//      to define or redefine Templates and Options Templates in bulk.
//
//      +--------+-------------------------------------------------+
//      |        | +----------+     +----------+      +----------+ |
//      |Message | | Template |     | Template |      | Options  | |
//      | Header | | Set      | ... | Set      | ...  | Template | |
//      |        | |          |     |          |      | Set      | |
//      |        | +----------+     +----------+      +----------+ |
//      +--------+-------------------------------------------------+

// DataRecord consist only of one or more fieldValues.
// https://tools.ietf.org/html/rfc7011#section-3.4.3
type DataRecord struct {
	FieldValues []interface{} `rethinkdb:"field_values"`
}

func (r *DataRecord) String() string {
	return fmt.Sprintf("data record (%v)", fieldDescribed(uint16(len(r.FieldValues))))
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
