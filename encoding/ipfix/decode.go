package ipfix

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

// TODO: what about padding (less bytes in data structures than declared message size)?

// Unmarshal stream of bytes into a IPFIX Message.
// A single Message consists of a Message Header and a number of Sets.
// https://www.ietf.org/rfc/rfc7011.html#section-3
func Unmarshal(in []byte, templates *[]TemplateRecord, out *Message) error {
	reader := bytes.NewReader(in)

	if err := binary.Read(reader, binary.BigEndian, &out.Header); err != nil {
		return err
	}

	var setHeader SetHeader
	// Read Set by Set starting with SetHeader
	if err := binary.Read(reader, binary.BigEndian, &setHeader); err != nil {
		return err
	}

	for bytesRead := headerLength + setHeaderLength; bytesRead < out.Header.Length; {
		switch setHeader.Type() {
		case TemplateSetType:
			var templateSet TemplateSet
			setBytesRead, err := readTemplateSet(reader, &templateSet)

			if err != nil {
				return err
			}

			out.Sets = append(out.Sets, templateSet)

			if templates != nil {
				*templates = append(*templates, templateSet.Records...)
			} else {
				return errors.New("invalid syntax: provide a non-empty list of TemplateRecords") // nolint:goerr113
			}

			bytesRead += setBytesRead
		case OptionsTemplateSetType:
		case DataSetType:
			var dataSet DataSet
			dataSet.SetHeader = setHeader
			// if templates are empty, there is nothing but reading dataSet.SetHeader.Length octets to be done
			if templates == nil {
				stream := make([]uint8, dataSet.SetHeader.Length-setHeaderLength)
				if err := binary.Read(reader, binary.BigEndian, &stream); err != nil {
					return err
				}
				// TODO
			}

		case UnknownSetType:
		default:
			// TODO there is (could be) some issue with padding, maybe not always bytesRead would be equeal to
			//      setBytesRead but rather setHeader.Length ???
		}
	}

	return nil
}

func readTemplateSet(reader io.Reader, out *TemplateSet) (uint16, error) {
	// TODO: implement as some kind of super-method that would “automagically” allow for
	//       things like binary.Read(reader, binary.BinEndian, &aTemplateSet)
	var bytesRead uint16

	if err := binary.Read(reader, binary.BigEndian, &out.SetHeader); err != nil {
		return 0, err
	}

	for bytesRead = setHeaderLength; bytesRead < out.SetHeader.Length; {
		var templateRecord TemplateRecord
		// if err := binary.Read(reader, binary.BigEndian, &templateRecord.Header); err != nil {
		// 	return bytesRead, err
		// }
		// bytesRead += templateRecordHeaderLength
		for i := uint16(0); i < out.Header().Length; i++ {
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

		out.Records = append(out.Records, templateRecord)
	}

	return bytesRead, nil
}
