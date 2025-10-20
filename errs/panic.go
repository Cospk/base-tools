package errs

import (
	"fmt"
	"github.com/Cospk/base-tools/errs/stack"
)

// ErrPanic 将panic恢复的值转换为CodeError
func ErrPanic(r any) error {
	return ErrPanicMsg(r, ServerInternalError, "panic error", 9)
}

// ErrPanicMsg 将panic恢复的值转换为自定义CodeError
// skip参数用于调整堆栈追踪的起点，默认值9适用于大多数场景
func ErrPanicMsg(r any, code int, msg string, skip int) error {
	if r == nil {
		return nil
	}
	err := &codeError{
		code:   code,
		msg:    msg,
		detail: fmt.Sprint(r),
	}
	return stack.New(err, skip)
}
