package exceptions

import (
	"fmt"
	"newgit.fidibo.com/fidiborearc/krakend/plugins/saga/messages"
)

type Error struct {
	msg string
}

func (error *Error) Error() string {
	return error.msg
}

func StatusError(a string, c int) error {
	return &Error{fmt.Sprintf(messages.ClientStatusCodeError, a, c)}
}

func BackendCallError(a, b string) error {
	return &Error{fmt.Sprintf(messages.ClientServiceCallError, a, b)}
}

func CloseBodyError(a, b string) error {
	return &Error{fmt.Sprintf(messages.ClientCloseBodyError, a, b)}
}
