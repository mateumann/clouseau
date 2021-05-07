package collector

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/mateumann/clouseau/pkg/encoding/ipfix"
	"github.com/mateumann/clouseau/pkg/storage"
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

// Close closes a UDP connection.
func (c *UDPCollector) Close() error {
	if c == nil || c.conn == nil {
		return errCollectorNotInitialized
	}

	if err := c.conn.Close(); err != nil {
		return Error{err}
	}

	return nil
}

// LocalAddr returns an address of a UDP connection.
func (c *UDPCollector) LocalAddr() net.Addr {
	if c == nil || c.conn == nil {
		return nil
	}

	return c.conn.LocalAddr()
}

// Collect listens on the UDP connection for the incoming IPFIX Messages and stores them.
func (c *UDPCollector) Collect(s storage.Storage) {
	if c.conn == nil {
		log.Println("uninitialized UDP connection")

		return
	}
	defer c.conn.Close()

	decoder := ipfix.NewDecoder(c.debug)

	var bytesCount int

	var sourceAddr *net.UDPAddr

	var err error

	buf := make([]byte, c.bufferSize)

	log.Printf("UDPCollector collecting on '%s', buffer size %d bytes, debug %v", c.conn.LocalAddr(), len(buf), c.debug)

	for {
		if bytesCount, sourceAddr, err = c.conn.ReadFromUDP(buf); err != nil {
			log.Printf("error while reading from UDP occurred: %v\n", err)

			continue
		}

		c.process(decoder, s, buf, bytesCount, sourceAddr)
	}
}

func (c UDPCollector) process(decoder *ipfix.Decoder, s storage.Storage,
	buf []byte, bytesCount int, sourceAddr *net.UDPAddr) {
	if c.debug {
		log.Printf("UDPCollector got %d bytes from %v", bytesCount, sourceAddr)
	}

	decoder.NewContext(fmt.Sprintf("%s:%s", udpNetwork, sourceAddr))

	n := decode(decoder, buf[:bytesCount])

	if err := store(s, decoder, n); err != nil {
		log.Printf("error saving: %v\n", err)
	} else if err = cleanUp(decoder, n); err != nil {
		log.Printf("error cleaning up: %v\n", err)
	}
}

func decode(decoder io.Reader, p []byte) int {
	n, err := decoder.Read(p)
	if err != nil {
		log.Printf("error while decoding IPFIX message: %v\n", err)
	}

	return n
}

func store(s storage.Storage, decoder *ipfix.Decoder, n int) error {
	messagesDecoded := decoder.LastMessages(n)
	for _, m := range messagesDecoded {
		if err := s.Write(m); err != nil {
			return storage.DBError{}
		}
	}

	return nil
}

func cleanUp(decoder *ipfix.Decoder, n int) error {
	if err := decoder.FreeLastMessages(n); err != nil {
		return Error{err}
	}

	return nil
}
