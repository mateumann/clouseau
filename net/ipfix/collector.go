package ipfix

import (
	"log"
	"net"
)

const (
	ipfixNetwork = "udp"
)

func StartCollector(ip string, port int) (*net.UDPConn, error) {
	address := net.UDPAddr{IP: net.ParseIP(ip), Port: port, Zone: ""}

	return net.ListenUDP(ipfixNetwork, &address)
}

func CloseCollector(conn *net.UDPConn) {
	if err := conn.Close(); err != nil {
		log.Fatal(err)
	}
}
