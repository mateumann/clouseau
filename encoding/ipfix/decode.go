package ipfix

import (
    "bytes"
    "encoding/binary"
)

// TODO: what about padding (less bytes in data structures than declared message size)?

// Unmarshal stream of bytes into a IPFIX Message.
// A single Message consists of a Message Header and a number of Sets.
// https://www.ietf.org/rfc/rfc7011.html#section-3
func Unmarshal(in []byte, out *Message) error {
    reader := bytes.NewReader(in)

    if err := binary.Read(reader, binary.BigEndian, &out.Header); err != nil {
        return err
    }

    var setHeader SetHeader
    for bytesRead := headerLength; bytesRead < out.Header.Length; {
        // Read Set by Set starting with SetHeader
        if err := binary.Read(reader, binary.BigEndian, &setHeader); err != nil {
            return err
        }
        switch setHeader.Type() {
        case TemplateSetType:
            var templateSet TemplateSet
            err, templateSetBytesRead := readTemplateSet(reader, &templateSet)
            if err != nil {
                return err
            }
            out.Sets = append(out.Sets, templateSet)
            bytesRead += templateSetBytesRead
        case OptionsTemplateSetType:
            //var optionsTemplateRecord OptionsTemplateRecord
            //TODO
        case DataSetType:
            //var dataRecord DataRecord
        default:
        }
        //bytesRead += setHeader.Length
    }
    return nil
}

// TODO: implement as some kind of super-method that would “automagically” allow for
//       things like binary.Read(reader, binary.BinEndian, &aTemplateSet)
func readTemplateSet(reader *bytes.Reader, out *TemplateSet) (error, uint16) {
    var bytesRead uint16
    if err := binary.Read(reader, binary.BigEndian, &out.Header); err != nil {
        return err, 0
    }

    for bytesRead = setHeaderLength; bytesRead < out.Header.Length; {
        var templateRecord TemplateRecord
        if err := binary.Read(reader, binary.BigEndian, &templateRecord.Header); err != nil {
            return err, bytesRead
        }
        bytesRead += templateRecordHeaderLength
        for i := uint16(0); i < templateRecord.Header.FieldCount; i++ {
            var fieldSpecifier = FieldSpecifier{0, 0, 0}
            if err := binary.Read(reader, binary.BigEndian, &fieldSpecifier.InformationElementID); err != nil {
                return err, bytesRead
            }
            if err := binary.Read(reader, binary.BigEndian, &fieldSpecifier.FieldLength); err != nil {
                return err, bytesRead + 2
            }
            if fieldSpecifier.E() {
                if err := binary.Read(reader, binary.BigEndian, &fieldSpecifier.EnterpriseNumber); err != nil {
                    return err, bytesRead + 4
                }
            }
            templateRecord.FieldSpecifiers = append(templateRecord.FieldSpecifiers, fieldSpecifier)
            bytesRead += fieldSpecifier.Size()
        }
        out.Records = append(out.Records, templateRecord)
    }
    return nil, bytesRead
}

//func readDataRecord(reader *bytes.Reader, out *DataSetType)
