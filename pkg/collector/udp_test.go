package collector

import (
	"net"
	"testing"

	"github.com/mateumann/clouseau/pkg/encoding"
)

func TestUDPCollector_StartAndListenUDP(t *testing.T) {
	c, err := StartAndListen("udp", "127.0.0.1", 20000, 0, false)
	if err != nil {
		t.Errorf("unable to start UDP collector: %v", err)
	}

	if c == nil {
		t.Error("empty (nil) collector")

		return
	}

	if c.LocalAddr().String() != "127.0.0.1:20000" {
		t.Error("invalid connection address")
	}

	err = c.Close()
	if err != nil {
		t.Errorf("unable to close UDP connection: %v", err)
	}

	err = c.Close()

	if err == nil || err.Error() != "close udp 127.0.0.1:20000: use of closed network connection" {
		t.Errorf("expected use of closed network connection error, got %v", err)
	}
}

/*
func TestUDPCollector_Collect_nil_connection(t *testing.T) {
	c := UDPCollector{}
	c.Collect()
}
*/

func TestUDPCollector_decode(t *testing.T) {
	c := UDPCollector{}
	d := encoding.NewDecoder(&net.UDPAddr{IP: net.IPv4zero, Port: 20000, Zone: ""}, false)
	c.decode(d, []byte{1, 2, 3, 4, 5})
}
