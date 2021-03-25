package ipfix

import (
	"fmt"
	"io"
	"time"
)

// An IPFIX Message is a Message that originates at the Exporting
//      Process and carries the IPFIX records of this Exporting Process,
//      and whose destination is a Collecting Process.  An IPFIX Message
//      is encapsulated at the transport layer.

// The Message Header is the first part of an IPFIX Message; the
//      Message Header provides basic information about the Message, such
//      as the IPFIX version, length of the Message, Message sequence
//      number, etc.

// Figure A shows a summary of IPFIX terminology.
//
//    +------------------+---------------------------------------------+
//    |                  |                 Contents                    |
//    |                  +--------------------+------------------------+
//    |       Set        |      Template      |         Record         |
//    +------------------+--------------------+------------------------+
//    |     Data Set     |          /         |     Data Record(s)     |
//    +------------------+--------------------+------------------------+
//    |   Template Set   | Template Record(s) |           /            |
//    +------------------+--------------------+------------------------+
//    | Options Template |  Options Template  |           /            |
//    |       Set        |      Record(s)     |                        |
//    +------------------+--------------------+------------------------+

// The format of the IPFIX Message is shown in Figure B.
//
//         +----------------------------------------------------+
//         | Message Header                                     |
//         +----------------------------------------------------+
//         | Set                                                |
//         +----------------------------------------------------+
//         | Set                                                |
//         +----------------------------------------------------+
//           ...
//         +----------------------------------------------------+
//         | Set                                                |
//         +----------------------------------------------------+

// An IPFIX Message consisting of interleaved Template, Data, and
//      Options Template Sets, as shown in Figure C.
//
//     +--------+--------------------------------------------------------+
//     |        | +----------+ +---------+     +-----------+ +---------+ |
//     |Message | | Template | | Data    |     | Options   | | Data    | |
//     | Header | | Set      | | Set     | ... | Template  | | Set     | |
//     |        | |          | |         |     | Set       | |         | |
//     |        | +----------+ +---------+     +-----------+ +---------+ |
//     +--------+--------------------------------------------------------+

// An IPFIX Message consisting entirely of Data Sets, sent after the
//      appropriate Template Records have been defined and transmitted to
//      the Collecting Process, as shown in Figure D.
//
//       +--------+----------------------------------------------+
//       |        | +---------+     +---------+      +---------+ |
//       |Message | | Data    |     | Data    |      | Data    | |
//       | Header | | Set     | ... | Set     | ...  | Set     | |
//       |        | +---------+     +---------+      +---------+ |
//       +--------+----------------------------------------------+

// An IPFIX Message consisting entirely of Template and Options
//      Template Sets, as shown in Figure E.  Such a Message can be used
//      to define or redefine Templates and Options Templates in bulk.
//
//      +--------+-------------------------------------------------+
//      |        | +----------+     +----------+      +----------+ |
//      |Message | | Template |     | Template |      | Options  | |
//      | Header | | Set      | ... | Set      | ...  | Template | |
//      |        | |          |     |          |      | Set      | |
//      |        | +----------+     +----------+      +----------+ |
//      +--------+-------------------------------------------------+

// Message describes a IPFIX Message containing of a header and zero or more sets.
// https://tools.ietf.org/html/rfc7011#section-3.1
type Message struct {
	SourceAddr          string    `rethinkdb:"source"`  // source address of the Message
	Version             uint16    `rethinkdb:"version"` // version of the IPFIX, currently 0x000a
	Length              uint16    `rethinkdb:"length"`  // total length of Message, including the header and Set(s)
	exportTime          uint32    // time at witch Message left the exporter, Unix Epoch, UTC
	ExportedAt          time.Time `rethinkdb:"export_at"` // time at witch Message left the exporter, UTC
	SequenceNumber      uint32    `rethinkdb:"seq"`       // incremental sequence counter modulo 2^32
	ObservationDomainID uint32    `rethinkdb:"domain_id"` // a 32-bit identifier of the observation domain
	Sets                []Set     `rethinkdb:"sets"`
}

func (m *Message) String() string {
	return fmt.Sprintf("ipfix message %d (%s, %v, %v, %v) %v",
		m.Version, intDescribed(int(m.Length), "octet", "octets"), m.ExportedAt.Format(time.RFC3339),
		m.SequenceNumber, m.ObservationDomainID, intDescribed(len(m.Sets), "set", "sets"))
}

func (m Message) valid() bool {
	return m.Version == 0x000a && m.Length > 0
}

// unmarshal reads byte stream from the reader and parses it into a IPFIX Message. `templates` are used to store IPFIX
// (Options) Template Records used to decode raw data into a IPFIX Data Records.
func (m *Message) unmarshal(reader io.Reader, source string, templates *map[uint16]*TemplateRecord) (n int, err error) {
	n, err = unmarshalFields(reader,
		&m.Version, &m.Length, &m.exportTime, &m.SequenceNumber, &m.ObservationDomainID)
	m.ExportedAt = time.Unix(int64(m.exportTime), 0).UTC()

	if err != nil {
		return n, err
	}

	m.SourceAddr = source

	// if err != nil {
	// 	return n, err
	// }

	for n < int(m.Length) {
		bytesRead, err := m.unmarshalSet(reader, templates)
		n += bytesRead

		if err != nil {
			return n, err
		}
	}

	return n, err
}

func (m *Message) unmarshalSet(reader io.Reader, templates *map[uint16]*TemplateRecord) (n int, err error) {
	var s Set
	n, err = s.unmarshal(reader, templates)

	if err != nil {
		return n, err
	}

	m.Sets = append(m.Sets, s)

	return n, nil
}
