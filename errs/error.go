package errs

import (
	"bytes"
	"errors"
	"fmt"
)

// Error 基础错误接口，扩展标准error，提供类型判断、包装和添加上下文的能力
type Error interface {
	Is(err error) bool          // 判断错误类型
	Wrap() error                // 添加堆栈追踪
	WrapMsg(msg string, kv ...any) error // 添加消息和堆栈追踪
	error
}

// New 创建新的Error实例，支持可选的键值对
func New(s string, kv ...any) Error {
	return &errorString{
		s: toString(s, kv),
	}
}

type errorString struct {
	s string
}

func (e *errorString) Is(err error) bool {
	if err == nil {
		return false
	}
	var t *errorString
	ok := errors.As(err, &t)
	return ok && e.s == t.s
}

func (e *errorString) Error() string {
	return e.s
}

func (e *errorString) Wrap() error {
	return Wrap(e)
}

func (e *errorString) WrapMsg(msg string, kv ...any) error {
	return WrapMsg(e, msg, kv...)
}

// toString 将消息和键值对转换为字符串，格式："msg, key1=value1, key2=value2"
func toString(s string, kv []any) string {
	if len(kv) == 0 {
		return s
	} else {
		var buf bytes.Buffer
		buf.WriteString(s)

		for i := 0; i < len(kv); i += 2 {
			if buf.Len() > 0 {
				buf.WriteString(", ")
			}

			key := fmt.Sprintf("%v", kv[i])
			buf.WriteString(key)
			buf.WriteString("=")

			if i+1 < len(kv) {
				value := fmt.Sprintf("%v", kv[i+1])
				buf.WriteString(value)
			} else {
				buf.WriteString("MISSING")
			}
		}
		return buf.String()
	}
}
