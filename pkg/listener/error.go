package listener

import "fmt"

type Error struct {
	op  string
	err error
}

func (e *Error) Error() string {
	if e.err != nil {
		return fmt.Sprintf("listener error %s: %v", e.op, e.err)
	}

	return fmt.Sprintf("listener error %s", e.op)
}

func (e *Error) Unwrap() error {
	return e.err
}
