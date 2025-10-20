package log

// 切片打印时的最大长度限制,超过此长度将被截断
const (
	slicePrintLen = 30
)

// Slice 泛型切片类型,实现LogFormatter接口用于限制日志输出长度
type Slice[T any] []T

// Format 实现LogFormatter接口,限制切片在日志中的输出长度
// 如果切片长度超过slicePrintLen,只返回前slicePrintLen个元素
func (s Slice[T]) Format() any {
	if len(s) >= slicePrintLen {
		return s[0:slicePrintLen]
	}
	return s
}
