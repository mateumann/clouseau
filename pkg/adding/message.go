package adding

import (
    "fmt"
    "time"
)

// IPFIXHeader defines a structure which stores a IPFIX Header
type _IPFIXHeader struct {
	Version             uint16
	Length              uint16
	ExportTime          uint32
	SequenceNumber      uint32
	ObservationDomainID uint32
	SetID               uint16
	SetLength           uint16
	Payload             []byte
}

// IPFIXHeaderLength is a IPFIX Header's length
const _IPFIXHeaderLength int = 20

func (h IPFIXHeader) String() string {
    timestamp := time.Unix(int64(h.ExportTime), 0)
    return fmt.Sprintf("v%v; length %v [octets, header+sets]; exported at %v; seq %v; domain %v; setid %v; set length %v", h.Version, h.Length, timestamp, h.SequenceNumber, h.ObservationDomainID, h.SetID, h.SetLength)
}
