package gateway

import (
	"fmt"
	"golang.org/x/exp/slog"
)

type Errors struct {
	msg string
}

func NewErrors(msg string, err error) *Errors {
	return &Errors{msg: fmt.Sprintf("%s: %v", msg, err)}
}

func (e *Errors) Error() {

	slog.Error(e.msg)
}
