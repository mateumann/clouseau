package collector

import (
	"errors"
	"io"
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

// StartAndListenUDP creates an IPFIX UDPCollector and returns a reference to it.
func StartAndListenUDP(address string, port, bufferSize int, debug bool) (*UDPCollector, error) {
	ip := net.ParseIP(address)
	if ip == nil {
		return &UDPCollector{}, InvalidAddressError{address: address}
	}

	addr := &net.UDPAddr{IP: ip, Port: port, Zone: ""}

	conn, err := net.ListenUDP(udpNetwork, addr)
	if err != nil {
		return &UDPCollector{}, errors.Unwrap(err)
	}

	return &UDPCollector{addr, conn, bufferSize, debug}, nil
}

func (c *UDPCollector) Close() error {
	if c == nil || c.conn == nil {
		return errCollectorNotInitialized
	}

	return c.conn.Close()
}

func (c *UDPCollector) LocalAddr() net.Addr {
	if c == nil || c.conn == nil {
		return nil
	}

	return c.conn.LocalAddr()
}

// Collect listens on the UDP connection for the incoming IPFIX Messages and stores them.
func (c *UDPCollector) Collect() {
	if c.conn == nil {
		log.Println("uninitialized UDP connection")

		return
	}
	defer c.conn.Close()

	decoder := encoding.NewDecoder(c.addr, c.debug)

	for {
		buf := make([]byte, c.bufferSize)
		if bytesCount, fromAddr, err := c.conn.ReadFromUDP(buf); err != nil {
			log.Printf("error while reading from UDP occurred: %v\n", err)
		} else {
			log.Printf("%d bytes read from %v\n", bytesCount, fromAddr)
			c.decode(decoder, buf[:bytesCount])
		}
	}
}

func (c *UDPCollector) decode(decoder io.Reader, p []byte) {
	messagesCount, err := decoder.Read(p)
	if err != nil {
		log.Printf("error while decoding IPFIX message: %v\n", err)
	} else if c.debug {
		log.Printf("%d messages decoded\n", messagesCount)
	}
}
