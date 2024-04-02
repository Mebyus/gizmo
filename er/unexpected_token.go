package er

import "io"

type UnexpectedTokenError struct {
}

var _ Error = &UnexpectedTokenError{}

func (e *UnexpectedTokenError) Error() string {
	return ""
}

func (e *UnexpectedTokenError) Kind() Kind {
	return UnexpectedToken
}

func (e *UnexpectedTokenError) Format(w io.Writer, format string) error {
	return nil
}
