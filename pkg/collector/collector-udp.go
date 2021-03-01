package collector

import (
	"errors"
	"log"
	"net"
)

const (
	udpNetwork = "udp"
)

// UDPCollector is a structure with required parameters to collect IPFIX Messages over UDP.
type UDPCollector struct {
	addr net.UDPAddr
	conn *net.UDPConn
}

// NewUDPCollector creates an IPFIX UDPCollector and returns a reference to it.
func NewUDPCollector(address string, port, bufferSize int, debug bool) (*UDPCollector, error) {
	addr := net.UDPAddr{IP: net.ParseIP(address), Port: port, Zone: ""}

	conn, err := net.ListenUDP(udpNetwork, &addr)
	if err != nil {
		return &UDPCollector{}, errors.Unwrap(err)
	}

	return &UDPCollector{addr, conn}, nil
}

// Collect listens on the UDP connection for the incoming IPFIX Messages and stores them.
func (c *UDPCollector) Collect() {
	defer c.conn.Close()

	// TODO mateumann

	log.Printf("UDPCollector collecting\n")
}
