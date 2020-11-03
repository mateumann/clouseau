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

func (s *fieldSpecifier) String() string {
	if s.valid() {
		if s.e() {
			return fmt.Sprintf("enterprise-specific field #%v (%v)",
				s.informationElementID, s.enterpriseNumber)
		}

		return fmt.Sprintf("IANA-registered field #%v",
			InformationElementByID(s.informationElementID))
	}

	return fmt.Sprintf("invalid field #%v", s.informationElementID)
}

func (s *fieldSpecifier) e() bool {
	return s.informationElementID>>15 == 1
}

func (s *fieldSpecifier) valid() bool {
	if s.e() {
		return s.informationElementID >= 32768 && s.enterpriseNumber != 0
	}

	return s.informationElementID < 32768 && s.enterpriseNumber == 0
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

// set describes a generic term for a collection of records that have a similar structure.  There are three different
// types of set: template sets, options template sets, and data sets.  Each of these sets consists of a set header and
// one or more records.  The exporting process may insert some padding octets, so that the subsequent set starts at an
// aligned boundary.
// https://tools.ietf.org/html/rfc7011#section-3.3.1
type set struct {
	header  setHeader
	records []interface{}
}

func (s *set) String() string {
	recordsDescription := "records"
	if len(s.records) == 1 {
		recordsDescription = "record"
	}

	return fmt.Sprintf("%v %v %v", s.header.String(), len(s.records), recordsDescription)
}

// setHeader describes a set by identifying it by a setID and declaring its length.
// https://tools.ietf.org/html/rfc7011#section-3.3.2
type setHeader struct {
	setID  uint16
	length uint16 // total length of the set, including setHeader, all records, and the optional padding; in octets
}

// Set ID constants defining valid IPFIX set identifiers (or a minimum valid set id in case of MinimumDataSetID).
const (
	TemplateSetID        = 2   // TemplateSetID identifies a IPFIX Template Set
	OptionsTemplateSetID = 3   // OptionsTemplateSetID identifies a IPFIX Options Template Set
	MinimumDataSetID     = 256 // set ids >= MinimumDataSetID identify IPFIX Data Sets
)

func (h *setHeader) String() string {
	if h.valid() {
		var setGenre string

		switch {
		case h.setID == TemplateSetID:
			setGenre = "template set"
		case h.setID == OptionsTemplateSetID:
			setGenre = "options template set"
		case h.setID >= MinimumDataSetID:
			setGenre = "data set"
		}

		return fmt.Sprintf("%v (%v octets)", setGenre, h.length)
	}

	return fmt.Sprintf("invalid set #%v", h.setID)
}

func (h *setHeader) valid() bool {
	return (h.setID == 2 || h.setID == 3 || h.setID >= 256) && h.length > 3
}
