package collector

import (
	"errors"
	"fmt"

	"github.com/mateumann/clouseau/pkg/storage"
)

// invalidAddressError is the error type usually returned by functions in the collector
// package.  It describes the invalid network address.
type invalidAddressError struct {
	address string
}

func (e invalidAddressError) Error() string {
	return fmt.Sprintf("invalid address: %v", e.address)
}

// notImplementedNetworkError is an Error that occurs when collector.StartAndListen is called for a network which does
// not have a Collector implementation.
type notImplementedNetworkError struct {
	network string
}

func (e notImplementedNetworkError) Error() string {
	return fmt.Sprintf("network not implemented: %v", e.network)
}

// Error represents general collector error.  This is a wrapper for errors from other packages.
type Error struct {
	err error
}

// Unwrap returns the original error.
func (e Error) Unwrap() error {
	return e.err
}

func (e Error) Error() string {
	return fmt.Sprintf("collector error: %v", e.err)
}

var errCollectorNotInitialized = errors.New("connector not initialized")

// Collector is an interface defining a IPFIX collector which collects incoming messages.
type Collector interface {
	Close() error
	Collect(s storage.Storage, debug bool)
}

// StartAndListen starts listening on [network]://[address]:[port] (e.g. udp://localhost:4739) for the incoming
// IPFIX messages.
func StartAndListen(network, address string, port, bufferSize int) (Collector, error) {
	switch network {
	case "udp":
		return StartAndListenUDP(address, port, bufferSize)
	default:
		return nil, notImplementedNetworkError{network: network}
	}
}
