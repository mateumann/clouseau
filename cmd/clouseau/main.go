package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"
	"time"
)

const IPFIXHeaderLength int = 20

// IPFIXHeader defines a structure which stores IPFIX's Header
type IPFIXHeader struct {
	Version             uint16
	MessageLength       uint16
	ExportTimestamp     uint32
	SequenceNumber      uint32
	ObservationDomainID uint32
	SetID               uint16
	SetLength           uint16
	Payload             []byte
}

func (h IPFIXHeader) String() string {

	timestamp := time.Unix(int64(h.ExportTimestamp), 0)
	return fmt.Sprintf("v%v; message length %v; exported at %v; seq %v; domain %v; setid %v; set length %v", h.Version, h.MessageLength, timestamp, h.SequenceNumber, h.ObservationDomainID, h.SetID, h.SetLength)
}

func unmarshalIPFIXHeader(data []byte, len int) (*IPFIXHeader, error) {
	if len < IPFIXHeaderLength {
		return nil, errors.New("Not enough data to unmarshal IPFIX Header")
	}

	var header IPFIXHeader

	reader := bytes.NewReader(data)

	err := binary.Read(reader, binary.BigEndian, &header.Version)
	if err != nil {
		return nil, err
	}

	err = binary.Read(reader, binary.BigEndian, &header.MessageLength)
	if err != nil {
		return nil, err
	}

	err = binary.Read(reader, binary.BigEndian, &header.ExportTimestamp)
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

	bytesToRead := len - IPFIXHeaderLength

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
		Port: 2055,
		IP:   net.ParseIP("0.0.0.0"),
	})

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	fmt.Printf("Server listening %s\n", conn.LocalAddr().String())

	for {
		buf := make([]byte, IPFIXHeaderLength+8192)
		buflen, remote, err := conn.ReadFromUDP(buf[:])
		header, err := unmarshalIPFIXHeader(buf, buflen)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("received %d bytes from %s; header %s; len(payload) %d\n", buflen, remote, header, len(header.Payload))
	}

}
