package main

import (
    "bytes"
    "encoding/binary"
    "errors"
    "fmt"
    "log"
    "net"

    "encoding/ipfix"
)

func unmarshalIPFIXHeader(data []byte, len int) (*adding.IPFIXHeader, error) {
	if len < adding.IPFIXHeaderLength {
		return nil, errors.New("not enough data to unmarshal IPFIX Header")
	}

	var header adding.IPFIXHeader

	reader := bytes.NewReader(data)

	err := binary.Read(reader, binary.BigEndian, &header.Version)
	if err != nil {
		return nil, err
	}

	err = binary.Read(reader, binary.BigEndian, &header.Length)
	if err != nil {
		return nil, err
	}

	err = binary.Read(reader, binary.BigEndian, &header.ExportTime)
	if err != nil {
		return nil, err
	}

	err = binary.Read(reader, binary.BigEndian, &header.SequenceNumber)
	if err != nil {
		return nil, err
	}

	err = binary.Read(reader, binary.BigEndian, &header.ObservationDomainID)
	if err != nil {
		return nil, err
	}

	err = binary.Read(reader, binary.BigEndian, &header.SetID)
	if err != nil {
		return nil, err
	}

	err = binary.Read(reader, binary.BigEndian, &header.SetLength)
	if err != nil {
		return nil, err
	}

	bytesToRead := len - adding.IPFIXHeaderLength

	payload := make([]byte, bytesToRead)

	err = binary.Read(reader, binary.BigEndian, &payload)
	if err != nil {
		return nil, err
	}

	header.Payload = payload

	return &header, nil
}

func main() {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		Port: 4739,
		IP:   net.ParseIP("0.0.0.0"),
	})

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	fmt.Printf("Server listening %s\n", conn.LocalAddr().String())

	for {
		buf := make([]byte, adding.IPFIXHeaderLength+8192)
		buflen, remote, err := conn.ReadFromUDP(buf[:])
		header, err := unmarshalIPFIXHeader(buf, buflen)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("received %d bytes from %s; header %s; len(payload) %d\n", buflen, remote, header, len(header.Payload))
	}

}
