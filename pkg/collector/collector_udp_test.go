package collector

import (
	"testing"
)

func TestStartAndListen_UDP(t *testing.T) {
	c, err := StartAndListen("udp", "127.0.0.1", 20000, 0, false)
	if err != nil {
		t.Errorf("unable to start collector: %v", err)
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
