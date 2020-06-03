package ipfix

import (
    "bytes"
    "encoding/binary"
)

//Unmarshal stream of bytes into an empty interface, based somehow on JSON's Unmarshal
func Unmarshal(in []byte, out *Message) error {
    reader := bytes.NewReader(in)

    err := binary.Read(reader, binary.BigEndian, &out.Header)
    if err != nil {
        return err
    }

    return nil
}

//
//func Unmarshal(data []byte, v interface{}) error {
//	// Check for well-formedness.
//	// Avoids filling out half a data structure
//	// before discovering a JSON syntax error.
//	var d decodeState
//	err := checkValid(data, &d.scan)
//	if err != nil {
//		return err
//	}
//
//	d.init(data)
//	return d.unmarshal(v)
//}
//
//
//func unmarshalIPFIXHeader(data []byte, len int) (*adding.IPFIXHeader, error) {
//    if len < adding.IPFIXHeaderLength {
//        return nil, errors.New("not enough data to unmarshal IPFIX Header")
//    }
//
//    var header adding.IPFIXHeader
//
//    reader := bytes.NewReader(data)
//
//    err := binary.Read(reader, binary.BigEndian, &header.Version)
//    if err != nil {
//        return nil, err
//    }
//
//    err = binary.Read(reader, binary.BigEndian, &header.Length)
//    if err != nil {
//        return nil, err
//    }
//
//    err = binary.Read(reader, binary.BigEndian, &header.ExportTime)
//    if err != nil {
//        return nil, err
//    }
//
//    err = binary.Read(reader, binary.BigEndian, &header.SequenceNumber)
//    if err != nil {
//        return nil, err
//    }
//
//    err = binary.Read(reader, binary.BigEndian, &header.ObservationDomainID)
//    if err != nil {
//        return nil, err
//    }
//
//    err = binary.Read(reader, binary.BigEndian, &header.SetID)
//    if err != nil {
//        return nil, err
//    }
//
//    err = binary.Read(reader, binary.BigEndian, &header.SetLength)
//    if err != nil {
//        return nil, err
//    }
//
//    bytesToRead := len - adding.IPFIXHeaderLength
//
//    payload := make([]byte, bytesToRead)
//
//    err = binary.Read(reader, binary.BigEndian, &payload)
//    if err != nil {
//        return nil, err
//    }
//
//    header.Payload = payload
//
//    return &header, nil
//}
//
