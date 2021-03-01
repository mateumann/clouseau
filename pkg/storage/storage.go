package storage

import "github.com/mateumann/clouseau/pkg/encoding"

// Storage is an interface defining a IPFIX Message writer.
type Storage interface {
	Write(*encoding.Message)
}
