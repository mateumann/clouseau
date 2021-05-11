package collector

import (
	"errors"
	"testing"
)

func TestNotImplementedNetworkError(t *testing.T) {
	_, err := StartAndListen("qwe", "127.0.0.1", 20000, 0, false)
	if err.Error() != "network not implemented: qwe" {
		t.Errorf("Expected error 'network not implemented: qwe' got: '%v'", err)
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

	err = c.Close()

	if err == nil || err.Error() != "connector not initialized" {
		t.Errorf("expected 'connector not initialized' error, got %v", err)
	}
}
