package collector_test

import (
	"testing"

	"github.com/mateumann/clouseau/pkg/collector"
)

func TestCollector(t *testing.T) {
	_, err := collector.StartAndListen("udp", "127.0.0.1", 4739, 128, true)
	if err != nil {
		t.Errorf("unable to start collector: %v", err)
	}
}
