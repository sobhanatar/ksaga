package messages

import "fmt"

type Error struct {
	msg string
}

func (error *Error) Error() string {
	return error.msg
}

func StatusError(a string, c int) error {
	return &Error{fmt.Sprintf(ClientStatusCodeError, a, c)}
}

func BackendCallError(a, b string) error {
	return &Error{fmt.Sprintf(ClientServiceCallError, a, b)}
}

func CloseBodyError(a, b string) error {
	return &Error{fmt.Sprintf(ClientCloseBodyError, a, b)}
}
