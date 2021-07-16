package listener_test

import (
	"testing"

	"github.com/mateumann/clouseau/pkg/listener"
)

func TestNewUDPListener(t *testing.T) {
	l, err := listener.NewUDPListener("127.0.0.1:12345")
	if err != nil {
		t.Error(err)
	}

	err = l.Listen(0)
	if err != nil {
		t.Error(err)
	}
}
