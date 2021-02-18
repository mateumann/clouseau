package main

import (
	"fmt"
	"io"
	"log"

	"github.com/codingconcepts/env"
	flag "github.com/spf13/pflag"

	"github.com/mateumann/clouseau/net/ipfix"
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
	config, err := configure()
	if err != nil {
		log.Fatal(err)
	}

	c, err := ipfix.NewUDPCollector(config.ListenAddress, config.ListenPort, config.BufferSize, config.Debug)
	if err != nil {
		log.Fatal(err)
	}

	c.HandleFunc(ipfixHandler)
	// err = c.HandleFunc_OLD(nil, nil,
	// 	config.BufferSize,
	// 	config.Debug,
	// )
	// if err != nil {
	// 	c.Close()
	// 	log.Fatal(err)
	// }

	c.Close()
}

func ipfixHandler(w io.Writer) {
	fmt.Fprintf(w, "ipfixHandler here\n")
}
