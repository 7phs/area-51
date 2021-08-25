package server

import "fmt"

var (
	_ error = Unexpected{}
)

type Unexpected struct {
	err error
	msg string
}

func ErrUnexpected(msg string, err error) Unexpected {
	return Unexpected{
		err: err,
		msg: msg,
	}
}

func (e Unexpected) Error() string {
	return fmt.Sprintf("unextected error: '%s' [%v]", e.msg, e.err)
}

func (e Unexpected) Unwrap() error {
	return e.err
}
