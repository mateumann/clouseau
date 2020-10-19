package ipfix

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

// TODO: what about padding (less bytes in data structures than declared message size)?

// Unmarshal stream of bytes into a IPFIX Message.
// A single Message consists of a MessageHeader and an array of Set-s.
// https://www.ietf.org/rfc/rfc7011.html#section-3
func Unmarshal(in []byte, templates *[]TemplateRecord, out *Message) error {
	if out == nil {
		return errors.New("invalid syntax: output message pointer is nil") // nolint:goerr113
	}

	if templates == nil {
		return errors.New("invalid syntax: templates list is nil") // nolint:goerr113
	}

	reader := bytes.NewReader(in)

	if err := binary.Read(reader, binary.BigEndian, &out.Header); err != nil {
		return err
	}

	// There can be more than one sets in each Message, so…

	for bytesRead := messageHeaderLength; bytesRead < out.Header.Length; {
		var setHeader SetHeader

		// Read Set by Set starting with header
		if err := binary.Read(reader, binary.BigEndian, &setHeader); err != nil {
			return err
		}

		bytesRead += setHeaderLength

		switch setHeader.Type() {
		case TemplateSetType:
			templateSet := TemplateSet{
				header:  setHeader,
				records: nil,
			}

			setBytesRead, err := readTemplateSet(reader, &templateSet)
			if err != nil {
				return err
			}

			out.Sets = append(out.Sets, templateSet)

			*templates = append(*templates, templateSet.records...)

			bytesRead += setBytesRead

		case OptionsTemplateSetType:
		case DataSetType:
			var dataSet DataSet
			dataSet.header = setHeader
			// if templates are empty, there is nothing but reading dataSet.header.Length octets to be done
			if len(*templates) == 0 {
				stream := make([]uint8, dataSet.header.Length-setHeaderLength)
				if err := binary.Read(reader, binary.BigEndian, &stream); err != nil {
					return err
				}

				bytesRead += dataSet.header.Length - setHeaderLength

				out.Sets = append(out.Sets, dataSet)
			}

		case UnknownSetType:
		default:
			bytesRead++ // FIXME this is a temporary solution while not everything is yet implemented
		}
	}

	return nil
}

// TODO: implement as some kind of super-method that would “automagically” allow for
//       things like binary.Read(reader, binary.BinEndian, &aTemplateSet)
func readTemplateSet(reader io.Reader, out *TemplateSet) (uint16, error) {
	var bytesRead uint16

	for bytesRead = 0; bytesRead < out.header.Length-setHeaderLength; {
		var templateRecord TemplateRecord

		if err := binary.Read(reader, binary.BigEndian, &templateRecord.TemplateHeader); err != nil {
			return 0, err
		}

		bytesRead += templateRecordHeaderLength

		for i := uint16(0); i < templateRecord.TemplateHeader.FieldCount; i++ {
			fieldSpecifier := FieldSpecifier{0, 0, 0}

			if err := binary.Read(reader, binary.BigEndian, &fieldSpecifier.ElementID); err != nil {
				return bytesRead, err
			}

			if err := binary.Read(reader, binary.BigEndian, &fieldSpecifier.FieldLength); err != nil {
				return bytesRead + 2, err
			}

			if fieldSpecifier.E() {
				if err := binary.Read(reader, binary.BigEndian, &fieldSpecifier.EnterpriseNumber); err != nil {
					return bytesRead + 4, err
				}
			}

			templateRecord.FieldSpecifiers = append(templateRecord.FieldSpecifiers, fieldSpecifier)
			bytesRead += fieldSpecifier.Size()
		}

		out.records = append(out.records, templateRecord)
	}

	return bytesRead, nil
}
