package ipfix

import (
	"net"
)

const (
	ipfixNetwork = "udp"
)

// StartCollector starts a UDP server listening on a given IP:port address.
func StartCollector(ip string, port int) (*net.UDPConn, error) {
	address := net.UDPAddr{IP: net.ParseIP(ip), Port: port, Zone: ""}

	return net.ListenUDP(ipfixNetwork, &address)
}
