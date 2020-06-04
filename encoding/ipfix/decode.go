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

    err := binary.Read(reader, binary.BigEndian, &out.Header)
    if err != nil {
        return err
    }

    var setHeader SetHeader
    for bytesRead := headerLength; bytesRead < out.Header.Length; {
        // Read Set by Set starting with SetHeader
        err = binary.Read(reader, binary.BigEndian, &setHeader)
        if err != nil {
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
    }
    return nil
}

// TODO: implement as some kind of super-method that would automagically allow for
//       binary.Read(reader, binary.BinEndian, &aTemplateSet)
func readTemplateSet(reader *bytes.Reader, out *TemplateSet) (error, uint16) {
    var bytesRead uint16
    err := binary.Read(reader, binary.BigEndian, &out.Header)
    if err != nil {
        return err, 0
    }

    for bytesRead = setHeaderLength; bytesRead < out.Header.Length; {
        var templateRecord TemplateRecord
        err = binary.Read(reader, binary.BigEndian, &templateRecord.TemplateID)
        if err != nil {
            return err, 0
        }
        bytesRead += 2
        err = binary.Read(reader, binary.BigEndian, &templateRecord.FieldCount)
        if err != nil {
            return err, 0
        }
        bytesRead += 2
        bytesRead += templateRecordHeaderLength
        for i := uint16(0); i < templateRecord.FieldCount; i++ {
            var fieldSpecifier FieldSpecifier
            err = binary.Read(reader, binary.BigEndian, &fieldSpecifier)
            if err != nil {
                return err, 0
            }
            templateRecord.FieldSpecifiers = append(templateRecord.FieldSpecifiers, fieldSpecifier)
            bytesRead += fieldSpecifierLength
        }
        out.Records = append(out.Records, templateRecord)
    }
    return nil, bytesRead
}

func readDataRecord(reader *bytes.Reader, out *DataSetType)
