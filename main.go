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

const netFlowHeaderLength int = 20

// NetFlowHeader defines a structure which stores NetFlow's Header
type NetFlowHeader struct {
	Version   uint16
	Count     uint16
	SysUpTime uint32
	UnixSec   uint32
	Seq       uint32
	SourceID  uint32
	Payload   []byte
}

func (h NetFlowHeader) String() string {

	uptime := time.Duration(int64(h.SysUpTime) * 1000 * 1000)
	seconds := time.Unix(int64(h.UnixSec), 0)
	return fmt.Sprintf("v%v (%v records to follow); uptime %v; seconds %v; seq %v; source %v ", h.Version, h.Count, uptime, seconds, h.Seq, h.SourceID)
}

func unmarshalNetFlowHeader(data []byte, len int) (*NetFlowHeader, error) {
	if len < netFlowHeaderLength {
		return nil, errors.New("Not enough data to unmarshal NetFlow Header")
	}

	var header NetFlowHeader

	reader := bytes.NewReader(data)

	err := binary.Read(reader, binary.BigEndian, &header.Version)
	if err != nil {
		return nil, err
	}

	err = binary.Read(reader, binary.BigEndian, &header.Count)
	if err != nil {
		return nil, err
	}

	err = binary.Read(reader, binary.BigEndian, &header.SysUpTime)
	if err != nil {
		return nil, err
	}

	err = binary.Read(reader, binary.BigEndian, &header.UnixSec)
	if err != nil {
		return nil, err
	}

	err = binary.Read(reader, binary.BigEndian, &header.Seq)
	if err != nil {
		return nil, err
	}

	err = binary.Read(reader, binary.BigEndian, &header.SourceID)
	if err != nil {
		return nil, err
	}

	bytesToRead := len - netFlowHeaderLength

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
		buf := make([]byte, 20+8*1024)
		buflen, remote, err := conn.ReadFromUDP(buf[:])
		header, err := unmarshalNetFlowHeader(buf, buflen)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("received %d bytes from %s; header %s; len(payload) %d\n", buflen, remote, header, len(header.Payload))
	}

}
