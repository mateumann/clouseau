package main

import (
	"fmt"
	"log"
	"net"

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

	var udpAddr *net.UDPAddr

	for {
		buf := make([]byte, bufferSize)
		if _, udpAddr, err = conn.ReadFromUDP(buf); err != nil {
			log.Fatal(err) //nolint:gocritic
		}

		decoder := ipfix.NewDecoder(udpAddr)

		// var templates []*templateRecord
		// var message ipfix.Message

		_, err = decoder.Read(buf)
		// err = message.UnmarshalBinary(bytes.NewReader(buf), templates)
		if err != nil {
			log.Print(fmt.Sprintf("%v\n", err)) //nolint:gocritic
		}

		fmt.Printf("%v\n", decoder)
	}
}
