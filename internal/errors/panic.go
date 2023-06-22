package errors

import (
	"errors"
	"fmt"
	"runtime/debug"
)

func NewPanicError(x any) error {
	stack := fmt.Sprintf("panic: %+v;\nstack: %s", x, string(debug.Stack()))
	return errors.New(stack)
}
