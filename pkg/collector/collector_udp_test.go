package collector

import (
	"errors"
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

func TestStartAndListen_NotImplemented(t *testing.T) {
	c, err := StartAndListen("foo", "127.0.0.1", 20000, 0, false)
	if err == nil || !errors.Is(err, NotImplementedNetworkError{"foo"}) {
		t.Errorf("expected: 'network not implemented: foo' error, got: %v", err)
	}

	if c != nil {
		t.Errorf("expected: nil collector, got: %v", c)
	}
}

func TestStartAndListen_InvalidAddress(t *testing.T) {
	c, err := StartAndListen("udp", "bar", 20000, 0, false)
	if err == nil || !errors.Is(err, InvalidAddressError{"bar"}) {
		t.Errorf("expected: 'invalid address: bar' error, got: %v", err)
	}

	if addr := c.LocalAddr(); addr != nil {
		t.Errorf("expected local address: 'nil', got: %v", addr)
	}

	err = c.Close()

	if err == nil || err.Error() != "connector not initialized" {
		t.Errorf("expected 'connector not initialized' error, got %v", err)
	}
}

func TestStartAndListen_CannotListen(t *testing.T) {
	c, err := StartAndListen("udp", "127.0.0.1", -1, 0, false)
	if err == nil || err.Error() != "bind: invalid argument" {
		t.Errorf("expected: 'bind: invalid argument' error, got: %v", err)
	}

	if addr := c.LocalAddr(); addr != nil {
		t.Errorf("expected local address: 'nil', got: %v", addr)
	}

	err = c.Close()

	if err == nil || err.Error() != "connector not initialized" {
		t.Errorf("expected 'connector not initialized' error, got %v", err)
	}
}
