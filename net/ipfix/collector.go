package ipfix

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/mateumann/clouseau/encoding/ipfix"
)

const (
	collectorNetwork = "udp"
)

// Collector describes a UDP listener service which listens to the incoming data flow and calls handler function.
type Collector struct {
	address    net.UDPAddr
	conn       *net.UDPConn
	bufferSize int
	debug      bool
}

// NewUDPCollector starts a UDP server listening on a given IP:port address.
func NewUDPCollector(ip string, port int, bufferSize int, debug bool) (*Collector, error) {
	addr := net.UDPAddr{IP: net.ParseIP(ip), Port: port, Zone: ""}

	conn, err := net.ListenUDP(collectorNetwork, &addr)
	if err != nil {
		return nil, errors.Unwrap(err)
	}

	if debug {
		// TODO mateumann: use logger
		fmt.Printf("Server listening %v\n", conn.LocalAddr())
	}

	return &Collector{addr, conn, bufferSize, debug}, nil
}

func (c *Collector) Close() error {
	return c.conn.Close()
}

// HandleFunc registers a io.Writer interface which will store received ipfix messages.
// It also acts as kind of server function (for time being probably).
func (c *Collector) HandleFunc(f func(io.Writer)) {
	for {
		buf := make([]byte, c.bufferSize)
		if _, _, err := c.conn.ReadFromUDP(buf); err != nil {
			// TODO mateumann: use logger
			fmt.Printf("Error occurred: %v\n", err)
		}

		decoder := ipfix.NewDecoder(&c.address, c.debug) // TODO mateumann: is `&` needed?

		// var templates []*templateRecord
		// var message ipfix.Message

		_, err := decoder.Read(buf)
		// err = message.UnmarshalBinary(bytes.NewReader(buf), templates)
		if err != nil {
			log.Print(fmt.Sprintf("%v\n", err))
		}

		// TODO mateumann: call w.Write somehow on successful read

		fmt.Printf("%v\n", decoder)
	}
}

func (c *Collector) HandleFuncDeprecated(writer io.Writer, reader io.Reader, bufferSize int, debug bool) error {
	var udpAddr *net.UDPAddr

	for {
		buf := make([]byte, bufferSize)
		if _, _, err := c.conn.ReadFromUDP(buf); err != nil {
			fmt.Printf("Error occurred: %v. Exiting.\n", err)

			return errors.Unwrap(err)
		}

		decoder := ipfix.NewDecoder(udpAddr, debug)

		// var templates []*templateRecord
		// var message ipfix.Message

		_, err := decoder.Read(buf)
		// err = message.UnmarshalBinary(bytes.NewReader(buf), templates)
		if err != nil {
			log.Print(fmt.Sprintf("%v\n", err))
		}

		fmt.Printf("%v\n", decoder)
	}
}
