package collector

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
	address net.UDPAddr
	conn    *net.UDPConn
}

// StartCollector starts a UDP server listening on a given IP:port address.
func StartCollector(ip string, port int) (*Collector, error) {
	addr := net.UDPAddr{IP: net.ParseIP(ip), Port: port, Zone: ""}

	conn, err := net.ListenUDP(collectorNetwork, &addr)
	if err != nil {
		return nil, errors.Unwrap(err)
	}

	// http.HandleFunc()
	return &Collector{address: addr, conn: conn}, nil
}

func (c *Collector) Close() error {
	return c.conn.Close()
}

func (c *Collector) HandleFunc(writer io.Writer, reader io.Reader, bufferSize int, debug bool) error {
	fmt.Printf("Server listening %v\n", c.conn.LocalAddr())

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
