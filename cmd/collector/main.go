package main

import (
    "bytes"
    "encoding/binary"
    "errors"
    "fmt"
    "github.com/mateumann/clouseau/pkg/adding"
    "log"
    "net"
)

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
