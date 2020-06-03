package ipfix

import (
    "fmt"
    "time"
)

// A Message represents a IPFIX Message.  It consists of a single Message Header and multiple Sets
// https://www.ietf.org/rfc/rfc7011.html#section-3
type Message struct {
    Header Header
    Sets   []Set
}

// A Header represents a IPFIX header.
// see https://tools.ietf.org/html/rfc7011#section-3.1
type Header struct {
    VersionNumber       uint16
    Length              uint16
    ExportTime          uint32
    SequenceNumber      uint32
    ObservationDomainID uint32
}

// A Set represents any of three IPFIX sets: Data Set, Template Set or Options Template Set
type Set interface{}

//
func (h Header) String() string {
    timestamp := time.Unix(int64(h.ExportTime), 0)
    return fmt.Sprintf("IPFIX Message Header (%v) [domain %v; seq %v; exported at %v]", h.VersionNumber, h.ObservationDomainID, h.SequenceNumber, timestamp.Format("15:04:05"))
}
