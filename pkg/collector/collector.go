package collector

import (
	"errors"
	"fmt"
	"net"
)

// InvalidAddressError is the error type usually returned by functions in the collector
// package.  It describes the invalid network address.
type InvalidAddressError struct {
	address string
}

func (e InvalidAddressError) Error() string {
	return fmt.Sprintf("invalid address: %v", e.address)
}

// NotImplementedNetworkError is an Error that occurs when collector.StartAndListen is called for a network which does
// not have a Collector implementation.
type NotImplementedNetworkError struct {
	network string
}

func (e NotImplementedNetworkError) Error() string {
	return fmt.Sprintf("network not implemented: %v", e.network)
}

var errCollectorNotInitialized = errors.New("connector not initialized")

// Collector is an interface defining a IPFIX collector which collects incoming messages.
type Collector interface {
	Close() error
	LocalAddr() net.Addr
	Collect()
}

// StartAndListen starts listening on [network]://[address]:[port] (e.g. udp://localhost:4739) for the incoming
// IPFIX messages.
func StartAndListen(network, address string, port, bufferSize int, debug bool) (Collector, error) {
	switch network {
	case "udp":
		return StartAndListenUDP(address, port, bufferSize, debug)
	default:
		return nil, NotImplementedNetworkError{network: network}
	}
}
