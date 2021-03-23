package collector

import (
	"errors"
	"testing"

	"github.com/mateumann/clouseau/pkg/encoding"
)

type MockStorageOK struct{}

func (db MockStorageOK) Write(m *encoding.Message) (err error) {
	return nil
}

type MockStorageErr struct{}

// nolint:goerr113
func (db MockStorageErr) Write(m *encoding.Message) (err error) {
	return errors.New("no good")
}

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

func TestUDPCollector_Collect_exit_on_nil_connection(t *testing.T) {
	c := UDPCollector{}
	c.Collect(nil)
	t.Log("success")
}

func TestDecode(t *testing.T) {
	d := encoding.NewDecoder(false)
	n := decode(d, []byte{1, 2, 3, 4, 5})

	if n != 0 {
		t.Errorf("wrong number of decoded messages %v, want %v", n, 1)
	}
}

/*
func TestStore(t *testing.T) {
	// TODO test func store(s storage.Storage, decoder *encoding.Decoder, n int) error
}

func TestCleanUp(t *testing.T) {
	// TODO test func cleanUp(decoder *encoding.Decoder, n int)
}
*/
