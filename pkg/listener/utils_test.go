package listener

import (
	"net"
	"testing"
)

// nolint: funlen
func TestUDPAddressEqual(t *testing.T) {
	testCases := []struct {
		a    net.UDPAddr
		b    net.UDPAddr
		want bool
	}{
		{
			a: net.UDPAddr{
				IP:   net.IPv4(127, 0, 0, 1),
				Port: 12345,
				Zone: "",
			},
			b: net.UDPAddr{
				IP:   net.ParseIP("127.0.0.1"),
				Port: 12345,
				Zone: "",
			},
			want: true,
		},
		{
			a: net.UDPAddr{
				IP:   net.IPv6zero,
				Port: 12345,
				Zone: "eth0",
			},
			b: net.UDPAddr{
				IP:   net.ParseIP("::0"),
				Port: 12345,
				Zone: "eth0",
			},
			want: true,
		},
		{
			a: net.UDPAddr{
				IP:   net.IPv4zero,
				Port: 12345,
				Zone: "eth0",
			},
			b: net.UDPAddr{
				IP:   net.IPv6zero,
				Port: 12345,
				Zone: "eth0",
			},
			want: false,
		},
		{
			a: net.UDPAddr{
				IP:   net.IPv4allsys,
				Port: 12345,
				Zone: "",
			},
			b: net.UDPAddr{
				IP:   net.IPv4allsys,
				Port: 54321,
				Zone: "",
			},
			want: false,
		},
		{
			a: net.UDPAddr{
				IP:   net.IPv6loopback,
				Port: 12345,
				Zone: "",
			},
			b: net.UDPAddr{
				IP:   net.IPv6loopback,
				Port: 12345,
				Zone: "eth0",
			},
			want: false,
		},
	}

	for _, tc := range testCases {
		got := UDPAddressEqual(tc.a, tc.b)

		if got != tc.want {
			t.Errorf("%v == %v should be %v, got %v", tc.a, tc.b, tc.want, got)
		}
	}
}
