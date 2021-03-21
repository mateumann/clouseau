package storage

import (
	"fmt"

	"github.com/mateumann/clouseau/pkg/encoding"
)

// DBError describes an error in storage backend.
type DBError struct {
	err error
}

func (e DBError) Error() string {
	return fmt.Sprintf("database error: %v", e.err)
}

// Unwrap returns the underlying error.
func (e DBError) Unwrap() error {
	return e.err
}

// Storage is an interface defining a IPFIX Message writer.
type Storage interface {
	Write(*encoding.Message) error
}
