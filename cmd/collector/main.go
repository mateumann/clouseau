package main

import (
	"fmt"
	"log"

	"github.com/mateumann/clouseau/encoding/ipfix"
	"github.com/mateumann/clouseau/net/collector"
)

const (
	ipfixIP    = "0.0.0.0"
	ipfixPort  = 4739
	bufferSize = 8192
)

func main() {
	conn, err := collector.StartCollector(ipfixIP, ipfixPort)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	fmt.Printf("Server listening %v\n", conn.LocalAddr())

	for {
		buf := make([]byte, bufferSize)
		if _, _, err := conn.ReadFromUDP(buf); err != nil {
			log.Fatal(err) //nolint:gocritic
		}

		var message ipfix.Message

		err = message.UnmarshalBinary(buf)
		if err != nil {
			log.Print(fmt.Sprintf("%v\n", err)) //nolint:gocritic
		}

		fmt.Printf("%s\n", message.String())
	}
}
