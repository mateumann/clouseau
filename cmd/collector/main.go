package main

import (
	"fmt"
	"log"
	"net"

	flag "github.com/spf13/pflag"

	"github.com/mateumann/clouseau/encoding/ipfix"
	"github.com/mateumann/clouseau/net/collector"
)

func main() {
	debug := flag.Bool("debug", false, "Debug mode.")
	ipfixIP := flag.String("address", "0.0.0.0", "IP address of the interface to listen to.")
	ipfixPort := flag.Int("port", 4739, "The UDP port to listen to.")
	bufferSize := flag.Int("buffersize", 8192, "Buffer size for the incoming messages.")
	flag.Parse()

	fmt.Println("debug:", *debug)

	conn, err := collector.StartCollector(*ipfixIP, *ipfixPort)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	fmt.Printf("Server listening %v\n", conn.LocalAddr())

	var udpAddr *net.UDPAddr

	for {
		buf := make([]byte, *bufferSize)
		if _, udpAddr, err = conn.ReadFromUDP(buf); err != nil {
			log.Fatal(err) //nolint:gocritic
		}

		decoder := ipfix.NewDecoder(udpAddr)

		// var templates []*templateRecord
		// var message ipfix.Message

		_, err = decoder.Read(buf)
		// err = message.UnmarshalBinary(bytes.NewReader(buf), templates)
		if err != nil {
			log.Print(fmt.Sprintf("%v\n", err))
		}

		fmt.Printf("%v\n", decoder)
	}
}
