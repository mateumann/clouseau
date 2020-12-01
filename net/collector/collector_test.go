package collector

import (
	"testing"
)

func TestCollector_Start_Close(t *testing.T) {
	conn, err := StartCollector("127.0.0.1", 20000)
	if err != nil {
		t.Errorf("unable to start collector: %v", err)
	}

	if conn == nil {
		t.Error("empty (nil) collector connection")

		return
	}

	if conn.LocalAddr().String() != "127.0.0.1:20000" {
		t.Error("invalid connection address")
	}

	err = conn.Close()
	if err != nil {
		t.Errorf("unable to close UDP connection: %v", err)
	}

	err = conn.Close()

	if err == nil || err.Error() != "close udp 127.0.0.1:20000: use of closed network connection" {
		t.Errorf("expected use of closed network connection error, got %v", err)
	}
}
