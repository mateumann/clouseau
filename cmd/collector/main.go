package main

import (
	"fmt"
	"os"

	"github.com/codingconcepts/env"
	"github.com/spf13/pflag"

	"github.com/mateumann/clouseau/pkg/collector"
)

type configuration struct {
	Debug          bool   `env:"DEBUG" default:"false"`
	ListenProtocol string `env:"LISTEN_PROTOCOL" default:"udp"`
	ListenAddress  string `env:"LISTEN_ADDRESS" default:"0.0.0.0"`
	ListenPort     int    `env:"LISTEN_PORT" default:"4739"`
	BufferSize     int    `env:"BUFFER_SIZE" default:"8192"`
}

func configure() *configuration {
	c := configuration{}

	if err := env.Set(&c); err != nil {
		fmt.Fprintf(os.Stderr, "error has occurred during parsing env: %v\n", err)

		return &c
	}

	pflag.BoolVar(&c.Debug, "debug", c.Debug, "Debug mode.")
	pflag.StringVar(&c.ListenProtocol, "protocol", c.ListenProtocol, "Protocol to use (only UDP is available).")
	pflag.StringVar(&c.ListenAddress, "address", c.ListenAddress, "IP address of the interface to listen to.")
	pflag.IntVar(&c.ListenPort, "port", c.ListenPort, "The UDP port to listen to.")
	pflag.IntVar(&c.BufferSize, "buffersize", c.BufferSize, "Buffer size for the incoming messages.")
	pflag.Parse()

	return &c
}

func main() {
	config := configure()

	c, err := collector.StartAndListen(
		config.ListenProtocol, config.ListenAddress, config.ListenPort, config.BufferSize, config.Debug)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	fmt.Printf("UDPCollector collecting on %s:%d\n", config.ListenAddress, config.ListenPort)

	c.Collect()
}
