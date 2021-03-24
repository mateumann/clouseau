package ipfix

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
)

type contextKey uint8

const (
	sourceAddress contextKey = iota

	unknownSource = "unknown"
)

// Decoder represents IPFIX flow already received by the collecting process.
type Decoder struct {
	ctx              context.Context
	debug            bool
	err              error
	buf              []byte // leftover decoded output
	messages         []*Message
	data             []*DataRecord
	templates        map[uint16]*TemplateRecord
	optionsTemplates map[uint16]*OptionsTemplateRecord
}

func (d *Decoder) String() string {
	return fmt.Sprintf("ipfix decoder: %s, %s, %s",
		intDescribed(len(d.templates), "template", "templates"),
		intDescribed(len(d.optionsTemplates), "options template", "options templates"),
		intDescribed(len(d.messages), "message", "messages"))
}

// NewContext creates and returns a new Decoder's Context which carries a value of the request source address.
func (d *Decoder) NewContext(addr string) context.Context {
	d.ctx = context.WithValue(d.ctx, sourceAddress, addr)

	return d.ctx
}

// sourceAddress returns the source address value stored in ctx, if any.
func (d Decoder) sourceAddress() string {
	if d.ctx == nil {
		return unknownSource
	}

	a, ok := d.ctx.Value(sourceAddress).(string)

	if !ok {
		return unknownSource
	}

	return a
}

// Read reads (and decodes) a IPFIX Message(s) from p buffer.
func (d *Decoder) Read(p []byte) (n int, err error) {
	if d.err != nil {
		return 0, d.err
	}

	if d.buf == nil {
		d.buf = make([]byte, len(p))
		copy(d.buf, p)
	} else {
		d.buf = append(d.buf, p...)
	}

	r := bytes.NewReader(d.buf)

	var m Message

	nn, err := m.unmarshal(r, d.sourceAddress(), &d.templates)
	d.buf = d.buf[nn:]
	n += nn

	if err != nil && !errors.Is(err, io.EOF) {
		d.err = err

		return 0, CorruptInputError{n, errors.Unwrap(err)}
	}

	if m.valid() {
		d.appendMessage(m)

		return 1, nil
	}

	return 0, nil
}

func (d *Decoder) appendMessage(m Message) {
	d.messages = append(d.messages, &m)
}

// LastMessages returns a slice of the last `n` decded Messages.
func (d *Decoder) LastMessages(n int) []*Message {
	return d.messages[len(d.messages)-n:]
}

// FreeLastMessages removes last `n` Messages.
func (d *Decoder) FreeLastMessages(n int) error {
	if len(d.messages) < n {
		return FreeLastMessagesError{requestedIndex: n, messagesLength: len(d.messages)}
	}

	d.messages = d.messages[:len(d.messages)-n]

	return nil
}

// NewDecoder initializes a Decoder with provided source (net.UDPAddr) and default values of other parameters.
func NewDecoder(debug bool) *Decoder {
	return &Decoder{
		ctx:              context.Background(),
		err:              nil,
		buf:              make([]byte, 0),
		messages:         make([]*Message, 0),
		data:             make([]*DataRecord, 0),
		templates:        make(map[uint16]*TemplateRecord),
		optionsTemplates: make(map[uint16]*OptionsTemplateRecord),
		debug:            debug,
	}
}
