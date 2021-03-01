package storage

import "github.com/mateumann/clouseau/pkg/encoding"

type Storage interface {
	Write(*encoding.Message)
}
