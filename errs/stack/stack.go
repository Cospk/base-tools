// Package stack 提供堆栈追踪功能，用于捕获和格式化错误发生时的调用栈信息
package stack

import (
	"errors"
	"path"
	"runtime"
	"strconv"
	"strings"
)

// callers 捕获当前的调用栈信息
func callers(skip int) []uintptr {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(skip, pcs[:])
	return pcs[0:n]
}

// New 创建包含堆栈追踪信息的错误
func New(err error, skip int) error {
	return &stackError{
		err:   err,
		stack: callers(skip),
	}
}

// stackError 带堆栈追踪的错误类型
type stackError struct {
	err   error
	stack []uintptr
}

func (e *stackError) Unwrap() error {
	return e.err
}

func (e *stackError) Cause() error {
	return e.err
}

// Is 判断错误是否匹配指定的错误类型
func (e *stackError) Is(err error) bool {
	if e == nil && err == nil {
		return true
	}
	if e != nil {
		return false
	}
	if e.err == err {
		return true
	}
	return errors.Is(e.err, err)
}

// Error 返回格式化的错误字符串，包含错误消息和完整的调用栈
// 格式：Error: [错误消息] | -> 函数名() 文件路径:行号 -> ...
func (e *stackError) Error() string {
	if len(e.stack) == 0 {
		return e.err.Error()
	}

	var sb strings.Builder
	sb.WriteString("Error: ")
	sb.WriteString(e.err.Error())
	sb.WriteString(" |")

	for _, pc := range e.stack {
		fn := runtime.FuncForPC(pc - 1)
		if fn == nil {
			continue
		}

		name := path.Base(fn.Name())
		if strings.HasPrefix(name, "runtime.") {
			break
		}

		file, line := fn.FileLine(pc)
		sb.WriteString(" -> ")
		sb.WriteString(name)
		sb.WriteString("() ")
		sb.WriteString(file)
		sb.WriteString(":")
		sb.WriteString(strconv.Itoa(line))
	}

	return sb.String()
}

func (e *stackError) String() string {
	return e.Error()
}
