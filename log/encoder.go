package log

import (
	"fmt"

	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

// alignEncoder 自定义日志编码器包装器，用于对齐日志消息
type alignEncoder struct {
	zapcore.Encoder
}

// EncodeEntry 重写编码方法，将日志消息左对齐到50个字符宽度
func (ae *alignEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	entry.Message = fmt.Sprintf("%-50s", entry.Message)
	return ae.Encoder.EncodeEntry(entry, fields)
}
