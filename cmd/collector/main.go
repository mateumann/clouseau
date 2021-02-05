package main

import (
	"fmt"
	"log"
	"net"

	flag "github.com/spf13/pflag"

	"github.com/mateumann/clouseau/encoding/ipfix"
	"github.com/mateumann/clouseau/net/collector"
)

type collectorArgs struct {
	debug         bool
	listenAddress string
	listenPort    int
	bufferSize    int
}

func parseFlags() (a collectorArgs) {
	flag.BoolVar(&a.debug, "debug", false, "Debug mode.")
	flag.StringVar(&a.listenAddress, "address", "0.0.0.0", "IP address of the interface to listen to.")
	flag.IntVar(&a.listenPort, "port", 4739, "The UDP port to listen to.")
	flag.IntVar(&a.bufferSize, "buffersize", 8192, "Buffer size for the incoming messages.")
	flag.Parse()

	return a
}

func main() {
	args := parseFlags()

	conn, err := collector.StartCollector(args.listenAddress, args.listenPort)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	fmt.Printf("Server listening %v\n", conn.LocalAddr())

	var udpAddr *net.UDPAddr

	for {
		buf := make([]byte, args.bufferSize)
		if _, udpAddr, err = conn.ReadFromUDP(buf); err != nil {
			log.Fatal(err) //nolint:gocritic
		}

		decoder := ipfix.NewDecoder(udpAddr, args.debug)

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
