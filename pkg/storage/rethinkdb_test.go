package storage

import (
	"errors"
	"testing"
)

func TestNotImplementedNetworkError(t *testing.T) {
	_, err := NewRethinkDB("127.0.0.1:28015", "clouseau_test")

	if err == nil {
		t.Errorf("Expected RethinkDBError got: '%v'", err)
	}

	if err.Error() != "rethinkdb error: rethinkdb: dial tcp 127.0.0.1:28015: connect: connection refused" {
		t.Errorf("Expected RethinkDBError got: '%v'", err)
	}

	if errors.Unwrap(err).Error() != "rethinkdb: dial tcp 127.0.0.1:28015: connect: connection refused" {
		t.Errorf("Expected RethinkDBError got: '%v'", err)
	}
}
