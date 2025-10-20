package splitter

// SplitResult 保存分割操作的结果，包含一个字符串切片。
type SplitResult struct {
	Item []string
}

// Splitter 负责将字符串切片分割成多个部分。
type Splitter struct {
	splitCount int      // 要分割成的部分数量
	data       []string // 要分割的原始数据
}

// NewSplitter 使用指定的分割数量和数据创建新的 Splitter 实例。
func NewSplitter(splitCount int, data []string) *Splitter {
	return &Splitter{splitCount: splitCount, data: data}
}

// GetSplitResult 执行分割操作并返回 SplitResult 切片作为结果。
// 每个 SplitResult 包含一个字符串切片，代表原始数据的一部分。
func (s *Splitter) GetSplitResult() (result []*SplitResult) {
	remain := len(s.data) % s.splitCount
	integer := len(s.data) / s.splitCount

	for i := 0; i < integer; i++ {
		r := new(SplitResult)
		r.Item = s.data[i*s.splitCount : (i+1)*s.splitCount]
		result = append(result, r)
	}
	if remain > 0 {
		r := new(SplitResult)
		r.Item = s.data[integer*s.splitCount:]
		result = append(result, r)
	}
	return result
}
