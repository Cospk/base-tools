package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestSliceFormat 测试 Slice 类型的 Format 方法
func TestSliceFormat(t *testing.T) {
	t.Run("Small slice should not be truncated", func(t *testing.T) {
		// 小于30个元素的切片应该保持原样
		slice := Slice[int]{1, 2, 3, 4, 5}
		result := slice.Format()

		resultSlice, ok := result.(Slice[int])
		assert.True(t, ok, "Format should return Slice[int]")
		assert.Equal(t, 5, len(resultSlice), "Length should be 5")
		assert.Equal(t, slice, resultSlice, "Small slice should not be modified")
	})

	t.Run("Exactly 30 elements should not be truncated", func(t *testing.T) {
		// 正好30个元素
		slice := make(Slice[int], 30)
		for i := 0; i < 30; i++ {
			slice[i] = i
		}

		result := slice.Format()
		resultSlice := result.(Slice[int])
		assert.Equal(t, 30, len(resultSlice), "Length should be 30")
	})

	t.Run("Large slice should be truncated to 30", func(t *testing.T) {
		// 大于30个元素的切片应该被截断
		slice := make(Slice[int], 100)
		for i := 0; i < 100; i++ {
			slice[i] = i
		}

		result := slice.Format()
		resultSlice, ok := result.(Slice[int])
		assert.True(t, ok, "Format should return Slice[int]")
		assert.Equal(t, slicePrintLen, len(resultSlice), "Length should be truncated to 30")

		// 验证前30个元素内容正确
		for i := 0; i < slicePrintLen; i++ {
			assert.Equal(t, i, resultSlice[i], "Element %d should match", i)
		}
	})

	t.Run("String slice should work", func(t *testing.T) {
		// 测试字符串切片
		slice := make(Slice[string], 50)
		for i := 0; i < 50; i++ {
			slice[i] = "item"
		}

		result := slice.Format()
		resultSlice := result.(Slice[string])
		assert.Equal(t, slicePrintLen, len(resultSlice), "Length should be 30")
	})

	t.Run("Struct slice should work", func(t *testing.T) {
		// 测试结构体切片
		type User struct {
			ID   int
			Name string
		}

		slice := make(Slice[User], 40)
		for i := 0; i < 40; i++ {
			slice[i] = User{ID: i, Name: "user"}
		}

		result := slice.Format()
		resultSlice := result.(Slice[User])
		assert.Equal(t, slicePrintLen, len(resultSlice), "Length should be 30")
		assert.Equal(t, 0, resultSlice[0].ID, "First element should be correct")
	})

	t.Run("Empty slice should work", func(t *testing.T) {
		// 空切片
		slice := Slice[int]{}
		result := slice.Format()
		resultSlice := result.(Slice[int])
		assert.Equal(t, 0, len(resultSlice), "Empty slice should remain empty")
	})

	t.Run("Nil slice should work", func(t *testing.T) {
		// nil 切片
		var slice Slice[int]
		result := slice.Format()
		resultSlice := result.(Slice[int])
		assert.Equal(t, 0, len(resultSlice), "Nil slice should have length 0")
	})
}

// TestSliceAsLogFormatter 测试 Slice 实现了 LogFormatter 接口
func TestSliceAsLogFormatter(t *testing.T) {
	var _ LogFormatter = Slice[int]{}

	// 确保可以作为 LogFormatter 使用
	var formatter LogFormatter = Slice[int]{1, 2, 3}
	result := formatter.Format()

	_, ok := result.(Slice[int])
	assert.True(t, ok, "Format result should be Slice[int]")
}

// TestSlicePrintLen 测试常量定义
func TestSlicePrintLen(t *testing.T) {
	assert.Equal(t, 30, slicePrintLen, "slicePrintLen should be 30")
}
