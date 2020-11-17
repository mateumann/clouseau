package main

import (
	"fmt"
	"log"
	"net"
)

const (
	ipfixPort = 4739
)

func main() {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		Port: ipfixPort,
		IP:   net.ParseIP("0.0.0.0"),
	})
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	fmt.Printf("Server listening %s\n", conn.LocalAddr().String())
}
