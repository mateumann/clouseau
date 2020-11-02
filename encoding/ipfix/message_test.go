package ipfix

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"testing"
)

func TestFieldSpecifier_e(t *testing.T) {
	var fs fieldSpecifier

	for attempts := 100; attempts > 0; attempts-- {
		randomized, err := rand.Int(rand.Reader, big.NewInt(65536))
		if err != nil || randomized == nil {
			t.Errorf("Error occurred during getting a random number.")
		}

		id := uint16(randomized.Int64())
		want := id >= 32768
		fs = fieldSpecifier{
			informationElementID: id,
			fieldLength:          0,
			enterpriseNumber:     0,
		}
		got := fs.e()

		if want != got {
			t.Errorf("FieldSpecifier.e(%v) = %v, want %v", id, got, want)
		}
	}
}

func TestFieldSpecifier_String(t *testing.T) {
	var fs fieldSpecifier

	for attempts := 100; attempts > 0; attempts-- {
		randomized, err := rand.Int(rand.Reader, big.NewInt(65536))
		if err != nil || randomized == nil {
			t.Errorf("Error occurred during getting a random number.")
		}

		id := uint16(randomized.Int64())
		description := ""
		fs = fieldSpecifier{
			informationElementID: id,
			fieldLength:          id % 16,
			enterpriseNumber:     999,
		}

		if id >= 32768 {
			description = fmt.Sprintf("Enterprise-specific field #%v (%v, ", id, 999)
		} else {
			description = fmt.Sprintf("IANA-registered field %v (", InformationElementByID(id))
		}

		want := fmt.Sprintf("%v%v octets)", description, fs.fieldLength)
		got := fs.String()

		if want != got {
			t.Errorf("FieldSpecifier.String(%v) = %v, want %v", id, got, want)
		}
	}
}

func TestMessageHeader_String(t *testing.T) {
	var h messageHeader

	for attempts := 100; attempts > 0; attempts-- {
		randomized, err := rand.Int(rand.Reader, big.NewInt(4294967296))
		if err != nil || randomized == nil {
			t.Errorf("Error occurred during getting a random number.")
		}

		id := uint32(randomized.Int64())
		h = messageHeader{
			version:             0x000a,
			length:              uint16(id % 16),
			exportTime:          id,
			sequenceNumber:      1,
			observationDomainID: 1,
		}

		want := fmt.Sprintf("header %v (%v octets, %v, %v, %v)",
			h.version, h.length, h.exportTime, h.sequenceNumber, h.observationDomainID)
		got := h.String()

		if want != got {
			t.Errorf("MessageHeader.String(%v) = %v, want %v", id, got, want)
		}
	}
}
