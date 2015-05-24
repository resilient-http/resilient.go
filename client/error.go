package client

type Error struct {
	timeout bool
	err     error
}

func (e *Error) Timeout() bool {
	return e.timeout
}

func (e *Error) Error() string {
	return e.err.Error()
}
