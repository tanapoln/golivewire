package golivewire

import "fmt"

var (
	ErrBadRequest          = newHTTPError(400)
	ErrNotFound            = newHTTPError(404)
	ErrInternalServerError = newHTTPError(500)
)

type HTTPError interface {
	error
	HTTPStatusCode() int
}

type httpErr struct {
	statusCode int
	msg        string
}

func (h httpErr) HTTPStatusCode() int {
	return h.statusCode
}

func (h httpErr) Error() string {
	msg := fmt.Sprintf("Error %d", h.statusCode)
	if h.msg != "" {
		msg += ": " + h.msg
	}
	return msg
}

func (h httpErr) Message(msg string) httpErr {
	h.msg = msg
	return h
}

func (h httpErr) Err(err error) httpErr {
	h.msg = err.Error()
	return h
}

func newHTTPError(code int) httpErr {
	return httpErr{
		statusCode: code,
	}
}
