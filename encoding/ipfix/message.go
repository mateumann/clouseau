package ipfix

import "fmt"

// A Template is an ordered sequence of <type, length> pairs used to
//      completely specify the structure and semantics of a particular set
//      of information that needs to be communicated from an IPFIX Device
//      to a Collector.  Each Template is uniquely identifiable by means
//      of a Template ID.

// An IPFIX Message is a message that originates at the Exporting
//      Process and carries the IPFIX records of this Exporting Process,
//      and whose destination is a Collecting Process.  An IPFIX Message
//      is encapsulated at the transport layer.

// The Message Header is the first part of an IPFIX Message; the
//      Message Header provides basic information about the message, such
//      as the IPFIX version, length of the message, message sequence
//      number, etc.

// A Template Record defines the structure and interpretation of
//      fields in a Data Record.

// A Data Record is a record that contains values of the parameters
//      corresponding to a Template Record.

// An Options Template Record is a Template Record that defines the
//      structure and interpretation of fields in a Data Record, including
//      defining how to scope the applicability of the Data Record.

// A Set is a collection of records that have a similar structure,
//      prefixed by a header.  In an IPFIX Message, zero or more Sets
//      follow the Message Header.  There are three different types of
//      Sets: Template Sets, Options Template Sets, and Data Sets.

// A Template Set is a collection of one or more Template Records
//      that have been grouped together in an IPFIX Message.

// A Data Set is one or more Data Records, of the same type, that are
//      grouped together in an IPFIX Message.  Each Data Record is
//      previously defined by a Template Record or an Options Template
//      Record.

// An Information Element is a protocol- and encoding-independent
//      description of an attribute that may appear in an IPFIX Record.
//      Information Elements are defined in the IANA "IPFIX Information
//      Elements" registry [IANA-IPFIX].  The type associated with an
//      Information Element indicates constraints on what it may contain
//      and also determines the valid encoding mechanisms for use in
//      IPFIX.

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
//      Template Sets, as shown in Figure E.  Such a message can be used
//      to define or redefine Templates and Options Templates in bulk.
//
//      +--------+-------------------------------------------------+
//      |        | +----------+     +----------+      +----------+ |
//      |Message | | Template |     | Template |      | Options  | |
//      | Header | | Set      | ... | Set      | ...  | Template | |
//      |        | |          |     |          |      | Set      | |
//      |        | +----------+     +----------+      +----------+ |
//      +--------+-------------------------------------------------+

// fieldSpecifier describes a single data element in a IPFIX record
// https://tools.ietf.org/html/rfc7011#section-3.2
type fieldSpecifier struct {
	informationElementID uint16
	fieldLength          uint16
	enterpriseNumber     uint32
}

func (s *fieldSpecifier) e() bool {
	return s.informationElementID>>15 == 1
}

func (s *fieldSpecifier) String() string {
	if !s.e() {
		return fmt.Sprintf("IANA-registered field %v (%v octets)",
			InformationElementByID(s.informationElementID), s.fieldLength)
	}

	return fmt.Sprintf("Enterprise-specific field #%v (%v, %v octets)",
		s.informationElementID, s.enterpriseNumber, s.fieldLength)
}

// messageHeader describes a IPFIX message header
// https://tools.ietf.org/html/rfc7011#section-3.1
type messageHeader struct {
	version             uint16 // version of the IPFIX, currently 0x000a
	length              uint16 // total length of IPFIX Message, measured in octets, including the header and Set(s)
	exportTime          uint32 // time at witch the message left the exporter, in seconds since the UNIX epoch, UTC
	sequenceNumber      uint32 // incremental sequence counter modulo 2^32
	observationDomainID uint32 // a 32-bit identifier of the observation domain, locally unique to the exporting proc.
}

func (h *messageHeader) String() string {
	return fmt.Sprintf("header %v (%v octets, %v, %v, %v)",
		h.version, h.length, h.exportTime, h.sequenceNumber, h.observationDomainID)
}
