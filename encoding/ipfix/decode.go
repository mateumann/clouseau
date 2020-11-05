package ipfix

import (
	"strconv"
)

// decoder represents IPFIX flow already received by the collecting process.
type decoder struct {
	// err              error
	// r                io.Reader
	// buf              []byte // leftover input
	out []byte // leftover decoded output
	// messages         []message
	// data             []*dataRecord
	// templates        []*templateRecord
	// optionsTemplates []*optionsTemplateRecord
}

type CorruptInputError int64

func (e CorruptInputError) Error() string {
	return "ipfix: illegal data at input byte " + strconv.FormatInt(int64(e), 10)
}

// Read decodes a IPFIX message(s) from p buffer (implements a io.Writer interface).
func (d *decoder) Read(p []byte) (n int, err error) {
	// use leftover decoded output from last read
	if len(d.out) > 0 {
		// TODO
		return 0, CorruptInputError(0)
	}

	// TODO: read a message header, determine a message size, create buffer, read into a buffer, start parsing

	return 0, nil
}
