package ipfix

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// TODO: what about padding (less bytes in data structures than declared message size)?

// An InvalidUnmarshalError describes an invalid argument passed to Unmarshal.
// (The argument to Unmarshal must be a non-nil pointer to the Message struct.)
type InvalidUnmarshalError struct{}

func (e *InvalidUnmarshalError) Error() string {
	return "ipfix: Unmarshal(nil)"
}

// Decode represents the data used to decode IPFIX messages.
type Decode struct {
	templates []templateRecord
}

// Unmarshal stream of bytes into a IPFIX Message.
// A single Message consists of a messageHeader and an array of set-s.
// https://www.ietf.org/rfc/rfc7011.html#section-3
func (d *Decode) Unmarshal(in []byte, message *Message) error {
	if message == nil {
		return &InvalidUnmarshalError{}
	}

	reader := bytes.NewReader(in)

	if err := binary.Read(reader, binary.BigEndian, &message.Header); err != nil {
		return err
	}

	for bytesRead := messageHeaderLength; bytesRead < int64(message.Header.Length); {
		setBytesRead, err := d.readSet(reader, message)
		if err != nil {
			return err
		}

		bytesRead += setBytesRead
	}

	return nil
}

// readSet reads a single set from the reader and stores it in the Message or
// Decode structures depending on its type.
func (d *Decode) readSet(reader *bytes.Reader, message *Message) (int64, error) {
	var setHeader setHeader
	if err := binary.Read(reader, binary.BigEndian, &setHeader); err != nil {
		return 0, err
	}

	switch setHeader.Type() {
	case TemplateSetType:
		templateSet := templateSet{
			header:  setHeader,
			records: nil,
		}

		_, err := readTemplateSet(reader, &templateSet)
		if err != nil {
			return reader.Size() - int64(reader.Len()), err
		}

		message.Sets = append(message.Sets, templateSet)

		d.templates = append(d.templates, templateSet.records...)

	case OptionsTemplateSetType:
	case DataSetType:
		var dataSet dataSet
		dataSet.header = setHeader
		// if templates are empty, there is nothing but reading dataSet.header.Length octets to be done
		if len(d.templates) == 0 {
			stream := make([]uint8, dataSet.header.Length-setHeaderLength)
			if err := binary.Read(reader, binary.BigEndian, &stream); err != nil {
				return reader.Size() - int64(reader.Len()), err
			}

			message.Sets = append(message.Sets, dataSet)
		}

	case UnknownSetType:
	default:
		// FIXME this is a temporary solution while not everything is yet implemented
		fmt.Printf("ipfix Decode.readSet at %v", reader.Len())
	}

	return reader.Size() - int64(reader.Len()), nil
}

// TODO: implement as some kind of super-method that would “automagically” allow for
//       things like binary.Read(reader, binary.BinEndian, &aTemplateSet)
func readTemplateSet(reader *bytes.Reader, out *templateSet) (int64, error) {
	var bytesRead uint16

	for bytesRead = 0; bytesRead < out.header.Length-setHeaderLength; {
		var templateRecord templateRecord

		if err := binary.Read(reader, binary.BigEndian, &templateRecord.TemplateHeader); err != nil {
			return 0, err
		}

		bytesRead += templateRecordHeaderLength

		for i := uint16(0); i < templateRecord.TemplateHeader.FieldCount; i++ {
			fieldSpecifier := fieldSpecifier{0, 0, 0}

			if err := binary.Read(reader, binary.BigEndian, &fieldSpecifier.ElementID); err != nil {
				return reader.Size() - int64(reader.Len()), err
			}

			if err := binary.Read(reader, binary.BigEndian, &fieldSpecifier.FieldLength); err != nil {
				return reader.Size() - int64(reader.Len()), err
			}

			if fieldSpecifier.E() {
				if err := binary.Read(reader, binary.BigEndian, &fieldSpecifier.EnterpriseNumber); err != nil {
					return reader.Size() - int64(reader.Len()), err
				}
			}

			templateRecord.FieldSpecifiers = append(templateRecord.FieldSpecifiers, fieldSpecifier)
			bytesRead += fieldSpecifier.Size()
		}

		out.records = append(out.records, templateRecord)
	}

	return reader.Size() - int64(reader.Len()), nil
}
