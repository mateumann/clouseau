package main

import (
	"fmt"
	"log"
	"net"

	flag "github.com/spf13/pflag"

	"github.com/mateumann/clouseau/encoding/ipfix"
	"github.com/mateumann/clouseau/net/collector"
)

type configuration struct {
	debug         bool
	listenAddress string
	listenPort    int
	bufferSize    int
}

func configure() *configuration {
	c := configuration{}
	flag.BoolVar(&c.debug, "debug", false, "Debug mode.")
	flag.StringVar(&c.listenAddress, "address", "0.0.0.0", "IP address of the interface to listen to.")
	flag.IntVar(&c.listenPort, "port", 4739, "The UDP port to listen to.")
	flag.IntVar(&c.bufferSize, "buffersize", 8192, "Buffer size for the incoming messages.")
	flag.Parse()

	return &c
}

func main() {
	c := configure()

	conn, err := collector.StartCollector(c.listenAddress, c.listenPort)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	fmt.Printf("Server listening %v\n", conn.LocalAddr())

	var udpAddr *net.UDPAddr

	for {
		buf := make([]byte, c.bufferSize)
		if _, udpAddr, err = conn.ReadFromUDP(buf); err != nil {
			fmt.Printf("Error occurred: %v. Exiting.\n", err)

			return
		}

		decoder := ipfix.NewDecoder(udpAddr, c.debug)

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
