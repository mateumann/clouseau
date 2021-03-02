package collector

import (
	"errors"
	"log"
	"net"

	"github.com/mateumann/clouseau/pkg/encoding"
)

const (
	udpNetwork = "udp"
)

// UDPCollector is a structure with required parameters to collect IPFIX Messages over UDP.
type UDPCollector struct {
	addr       *net.UDPAddr
	conn       *net.UDPConn
	bufferSize int
	debug      bool
}

// NewUDPCollector creates an IPFIX UDPCollector and returns a reference to it.
func NewUDPCollector(address string, port, bufferSize int, debug bool) (*UDPCollector, error) {
	addr := &net.UDPAddr{IP: net.ParseIP(address), Port: port, Zone: ""}

	conn, err := net.ListenUDP(udpNetwork, addr)
	if err != nil {
		return &UDPCollector{}, errors.Unwrap(err)
	}

	return &UDPCollector{addr, conn, bufferSize, debug}, nil
}

// Collect listens on the UDP connection for the incoming IPFIX Messages and stores them.
func (c *UDPCollector) Collect() {
	defer c.conn.Close()

	decoder := encoding.NewDecoder(c.addr, c.debug)

	for {
		buf := make([]byte, c.bufferSize)
		if bytesCount, fromAddr, err := c.conn.ReadFromUDP(buf); err != nil {
			log.Printf("error while reading from UDP occurred: %v\n", err)
		} else {
			log.Printf("%d bytes read from %v\n", bytesCount, fromAddr)
		}

		messagesCount, err := decoder.Read(buf)
		if err != nil {
			log.Printf("error while decoding IPFIX message: %v\n", err)
		} else {
			log.Printf("%d messages decoded\n", messagesCount)
		}
	}
}
