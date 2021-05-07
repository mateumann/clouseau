package storage

import (
	"errors"
	"fmt"

	"github.com/mateumann/clouseau/pkg/encoding/ipfix"
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

// NewDBError returns a DBError which wraps an err.
func NewDBError(err error) DBError {
	return DBError{err}
}

// ErrStorageNotInitialized is an error returned when storage object is not initialized.
var ErrStorageNotInitialized = errors.New("storage not initialized")

// Storage is an interface defining a IPFIX Message writer.
type Storage interface {
	Write(*ipfix.Message) error
}
