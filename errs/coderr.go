// Package errs 提供统一的错误处理系统，支持错误码、堆栈追踪和错误关系管理。
package errs

import (
	"errors"
	"github.com/Cospk/base-tools/errs/stack"
	"strconv"
	"strings"
)

// stackSkip 堆栈追踪跳过的栈帧数量
const stackSkip = 4

// DefaultCodeRelation 全局错误码关系管理器，维护错误码的父子关系
var DefaultCodeRelation = newCodeRelation()

// CodeError 核心错误接口，扩展标准error，添加错误码、消息和详细信息
type CodeError interface {
	Code() int                          // 返回错误码
	Msg() string                        // 返回错误消息
	Detail() string                     // 返回详细信息
	WithDetail(detail string) CodeError // 添加详细信息并返回新错误
	Error                               // 嵌入Error接口
}

// NewCodeError 创建新的CodeError实例
func NewCodeError(code int, msg string) CodeError {
	return &codeError{
		code: code,
		msg:  msg,
	}
}

// codeError 是CodeError接口的实现
type codeError struct {
	code   int    // 错误码
	msg    string // 错误消息
	detail string // 详细信息
}

func (e *codeError) Code() int {
	return e.code
}

func (e *codeError) Msg() string {
	return e.msg
}

func (e *codeError) Detail() string {
	return e.detail
}

// WithDetail 添加详细信息，多次调用会累积信息
func (e *codeError) WithDetail(detail string) CodeError {
	var d string
	if e.detail == "" {
		d = detail
	} else {
		d = e.detail + ", " + detail
	}
	return &codeError{
		code:   e.code,
		msg:    e.msg,
		detail: d,
	}
}

// Wrap 包装错误并添加堆栈追踪
func (e *codeError) Wrap() error {
	return stack.New(e, stackSkip)
}

func (e *codeError) clone() *codeError {
	return &codeError{
		code:   e.code,
		msg:    e.msg,
		detail: e.detail,
	}
}

// WrapMsg 包装错误，添加消息、键值对和堆栈追踪
func (e *codeError) WrapMsg(msg string, kv ...any) error {
	retErr := e.clone()
	if msg != "" || len(kv) > 0 {
		detail := toString(msg, kv)
		if retErr.detail == "" {
			retErr.detail = detail
		} else {
			retErr.detail += ", " + detail
		}
	}
	return stack.New(retErr, stackSkip)
}

// Is 检查错误是否匹配，支持错误码相等和父子关系判断
func (e *codeError) Is(err error) bool {
	var codeErr CodeError
	ok := errors.As(Unwrap(err), &codeErr)
	if !ok {
		if err == nil && e == nil {
			return true
		}
		return false
	}
	if e == nil {
		return false
	}
	code := codeErr.Code()
	if e.code == code {
		return true
	}
	return DefaultCodeRelation.Is(e.code, code)
}

const initialCapacity = 3

// Error 返回错误的字符串表示，格式：[错误码] [消息] [详细信息]
func (e *codeError) Error() string {
	v := make([]string, 0, initialCapacity)
	v = append(v, strconv.Itoa(e.code), e.msg)

	if e.detail != "" {
		v = append(v, e.detail)
	}

	return strings.Join(v, " ")
}

// Unwrap 递归解包错误，返回最底层的原始错误
func Unwrap(err error) error {
	for err != nil {
		unwrap, ok := err.(interface {
			error
			Unwrap() error
		})
		if !ok {
			break
		}
		err = unwrap.Unwrap()
		if err == nil {
			return unwrap
		}
	}
	return err
}

// Wrap 包装标准错误并添加堆栈追踪
func Wrap(err error) error {
	if err == nil {
		return nil
	}
	return stack.New(err, stackSkip)
}

// WrapMsg 包装标准错误，添加消息、键值对和堆栈追踪
func WrapMsg(err error, msg string, kv ...any) error {
	if err == nil {
		return nil
	}
	err = NewErrorWrapper(err, toString(msg, kv))
	return stack.New(err, stackSkip)
}

// CodeRelation 错误码关系管理接口，用于建立错误码的父子关系
type CodeRelation interface {
	Add(codes ...int) error    // 添加错误码关系链，第一个是父错误码
	Is(parent, child int) bool // 判断child是否是parent的子错误码
}

func newCodeRelation() CodeRelation {
	return &codeRelation{m: make(map[int]map[int]struct{})}
}

// codeRelation 使用嵌套map存储错误码关系：map[父错误码]map[子错误码]struct{}
type codeRelation struct {
	m map[int]map[int]struct{}
}

const minimumCodesLength = 2

// Add 建立错误码的父子关系
func (r *codeRelation) Add(codes ...int) error {
	if len(codes) < minimumCodesLength {
		return New("codes length must be greater than 2", "codes", codes).Wrap()
	}
	for i := 1; i < len(codes); i++ {
		parent := codes[i-1]
		s, ok := r.m[parent]
		if !ok {
			s = make(map[int]struct{})
			r.m[parent] = s
		}
		for _, code := range codes[i:] {
			s[code] = struct{}{}
		}
	}
	return nil
}

// Is 判断child是否是parent的子错误码（或相等）
func (r *codeRelation) Is(parent, child int) bool {
	if parent == child {
		return true
	}
	s, ok := r.m[parent]
	if !ok {
		return false
	}
	_, ok = s[child]
	return ok
}
