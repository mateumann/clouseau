package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/codingconcepts/env"
	flag "github.com/spf13/pflag"

	"github.com/mateumann/clouseau/ipfix"
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

	f, err := os.OpenFile("/tmp/collector.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o600)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	addr := net.UDPAddr{IP: net.ParseIP(config.ListenAddress), Port: config.ListenPort, Zone: ""}

	err = ipfix.ListenAndCollect(addr, f, config.BufferSize, config.Debug)
	if err != nil {
		log.Printf("Error occurred: %v\n", err)
	}
}

// 	c, err := ipfix.NewUDPCollector(config.ListenAddress, config.ListenPort, config.BufferSize, config.Debug)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	f, err := os.OpenFile("/tmp/collector.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o600)
// 	if err != nil {
// 		c.Close()
// 		log.Fatal(err)
// 	}

// 	defer c.Close()
// 	defer f.Close()

// 	// NOTE mateumann, 2021-02-19
// 	// Now the idea is to have Decoder.Read injected as a reader function and YetToBeWrittenDatabaseAdapter.Write
// 	// as a writer function.  Each time collector receives a message (net.UDPConn.ReadFromUDP) it passes it to the
// 	// reader function and then calls the writer function.  Must also handle partial messages - should not always
// 	// call the reader but only if a whole message (or a few messages) are read/decoded.  Probably need to rewrite
// 	// Decoder.Read to return number of read messages instead of number of read bytes.
// 	// Might be worth calling the writer function (and the reader as well?) in a separate goroutine.

// 	// err = c.HandleFunc_OLD(nil, nil,
// 	// 	config.BufferSize,
// 	// 	config.Debug,
// 	// )
// 	// if err != nil {
// 	// 	c.Close()
// 	// 	log.Fatal(err)
// 	// }

// 	c.Serve(ipfixHandler(f))
// }

// func ipfixHandler(w io.Writer) io.Writer {
// 	fmt.Fprintf(w, "ipfixHandler here\n")

// 	return w
// }
