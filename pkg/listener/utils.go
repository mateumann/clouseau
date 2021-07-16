package listener

import "net"

func UDPAddressEqual(a net.UDPAddr, b net.UDPAddr) bool {
	return a.IP.Equal(b.IP) && a.Port == b.Port && a.Zone == b.Zone
}
