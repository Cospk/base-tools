package log

import (
	"fmt"

	"go.uber.org/zap/zapcore"
)

// ANSI前景色常量定义，用于终端彩色日志输出
const (
	Black   Color = iota + 30
	Red           // 错误日志
	Green
	Yellow  // 警告日志
	Blue    // 信息日志
	Magenta
	Cyan
	White   // 调试日志
)

var (
	// 日志级别到颜色的映射
	_levelToColor = map[zapcore.Level]Color{
		zapcore.DebugLevel:  White,
		zapcore.InfoLevel:   Blue,
		zapcore.WarnLevel:   Yellow,
		zapcore.ErrorLevel:  Red,
		zapcore.DPanicLevel: Red,
		zapcore.PanicLevel:  Red,
		zapcore.FatalLevel:  Red,
	}

	_unknownLevelColor           = make(map[zapcore.Level]string, len(_levelToColor))
	_levelToLowercaseColorString = make(map[zapcore.Level]string, len(_levelToColor))
	_levelToCapitalColorString   = make(map[zapcore.Level]string, len(_levelToColor))
)

// init 初始化日志级别到彩色字符串的映射
func init() {
	for level, color := range _levelToColor {
		_levelToLowercaseColorString[level] = color.Add(level.String())
		_levelToCapitalColorString[level] = color.Add(level.CapitalString())
	}
}

// Color 表示终端文本颜色，使用ANSI转义序列
type Color uint8

// Add 为字符串添加ANSI颜色码，格式：\x1b[颜色码m文本\x1b[0m
func (c Color) Add(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", uint8(c), s)
}
