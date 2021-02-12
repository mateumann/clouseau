package main

import (
	"fmt"
	"log"
	"net"

	"github.com/codingconcepts/env"
	flag "github.com/spf13/pflag"

	"github.com/mateumann/clouseau/encoding/ipfix"
	"github.com/mateumann/clouseau/net/collector"
)

type configuration struct {
	Debug         bool   `env:"DEBUG" default:"false"`
	ListenAddress string `env:"LISTEN_ADDRESS" default:"0.0.0.0"`
	ListenPort    int    `env:"LISTEN_PORT" default:"4739"`
	BufferSize    int    `env:"BUFFER_SIZE" default:"8192"`
}

func configure() (*configuration, error) {
	c := configuration{}

	if err := env.Set(&c); err != nil {
		return nil, fmt.Errorf("error reading configuration from the environment: %w", err)
	}

	flag.BoolVar(&c.Debug, "debug", c.Debug, "Debug mode.")
	flag.StringVar(&c.ListenAddress, "address", c.ListenAddress, "IP address of the interface to listen to.")
	flag.IntVar(&c.ListenPort, "port", c.ListenPort, "The UDP port to listen to.")
	flag.IntVar(&c.BufferSize, "buffersize", c.BufferSize, "Buffer size for the incoming messages.")
	flag.Parse()

	return &c, nil
}

func main() {
	c, err := configure()
	if err != nil {
		log.Fatal(err)
	}

	conn, err := collector.StartCollector(c.ListenAddress, c.ListenPort)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	fmt.Printf("Server listening %v\n", conn.LocalAddr())

	var udpAddr *net.UDPAddr

	for {
		buf := make([]byte, c.BufferSize)
		if _, udpAddr, err = conn.ReadFromUDP(buf); err != nil {
			fmt.Printf("Error occurred: %v. Exiting.\n", err)

			return
		}

		decoder := ipfix.NewDecoder(udpAddr, c.Debug)

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
