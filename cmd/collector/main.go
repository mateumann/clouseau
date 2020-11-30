package main

import (
	"fmt"
	"log"

	"github.com/mateumann/clouseau/encoding/ipfix"
	"github.com/mateumann/clouseau/net"
)

const (
	ipfixIP    = "0.0.0.0"
	ipfixPort  = 4739
	bufferSize = 8192
)

func main() {
	conn, err := net.StartCollector(ipfixIP, ipfixPort)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	fmt.Printf("Server listening %v\n", conn.LocalAddr())

	for {
		buf := make([]byte, bufferSize)
		_, remote, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Fatal(err) //nolint:gocritic
		}

		var message ipfix.Message

		err = message.UnmarshalBinary(buf)
		if err != nil {
			log.Fatal(err) //nolint:gocritic
		}

		fmt.Printf("%v:%v\t%v", remote.IP, remote.Port, message)
	}
}
