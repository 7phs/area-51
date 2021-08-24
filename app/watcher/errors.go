package watcher

import "fmt"

var (
	_ error = Unexpected{}
)

type Unexpected struct {
	err  error
	path string
}

func ErrUnexpected(err error) Unexpected {
	return Unexpected{
		err: err,
	}
}

func (e Unexpected) Error() string {
	return fmt.Sprintf("unextected error: '%s' [%v]", e.path, e.err)
}

func (e Unexpected) Unwrap() error {
	return e.err
}

type AlreadyShutdown string

func ErrAlreadyShutdown() AlreadyShutdown {
	return AlreadyShutdown("")
}

func (e AlreadyShutdown) Error() string {
	return "already shutdown"
}
