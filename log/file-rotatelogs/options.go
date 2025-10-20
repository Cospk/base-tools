package rotatelogs

import (
	"github.com/Cospk/base-tools/log/file-rotatelogs/internal/option"
	"time"
)

// 选项键常量定义
const (
	optkeyClock         = "clock"
	optkeyHandler       = "handler"
	optkeyLinkName      = "link-name"
	optkeyMaxAge        = "max-age"
	optkeyRotationTime  = "rotation-time"
	optkeyRotationSize  = "rotation-size"
	optkeyRotationCount = "rotation-count"
	optkeyForceNewFile  = "force-new-file"
)

// WithClock 创建一个新的Option,设置RotateLogs对象用于确定当前时间的时钟接口
//
// 默认使用rotatelogs.Local,它返回本地时区的当前时间。
// 如果想使用UTC时间,可以将rotatelogs.UTC作为参数传递给此选项,并传递给构造函数。
func WithClock(c Clock) Option {
	return option.New(optkeyClock, c)
}

// WithLocation 创建一个新的Option,设置RotateLogs对象用于确定当前时间的"Clock"接口
//
// 此选项通过始终返回给定位置的时间来工作。
func WithLocation(loc *time.Location) Option {
	return option.New(optkeyClock, clockFn(func() time.Time {
		return time.Now().In(loc)
	}))
}

// WithLinkName 创建一个新的Option,设置链接到当前正在使用的文件名的符号链接名称
func WithLinkName(s string) Option {
	return option.New(optkeyLinkName, s)
}

// WithMaxAge 创建一个新的Option,设置日志文件在从文件系统清除之前的最大保留时间
func WithMaxAge(d time.Duration) Option {
	return option.New(optkeyMaxAge, d)
}

// WithRotationTime 创建一个新的Option,设置日志轮转的时间间隔
func WithRotationTime(d time.Duration) Option {
	return option.New(optkeyRotationTime, d)
}

// WithRotationSize 创建一个新的Option,设置触发日志轮转的文件大小阈值
func WithRotationSize(s int64) Option {
	return option.New(optkeyRotationSize, s)
}

// WithRotationCount 创建一个新的Option,设置在从文件系统清除之前应保留的文件数量
func WithRotationCount(n uint) Option {
	return option.New(optkeyRotationCount, n)
}

// WithHandler 创建一个新的Option,指定在事件发生时被调用的Handler对象
// 目前支持`FileRotated`事件
func WithHandler(h Handler) Option {
	return option.New(optkeyHandler, h)
}

// ForceNewFile 确保每次调用New()时都创建一个新文件。
// 如果基础文件名已存在,则执行隐式轮转
func ForceNewFile() Option {
	return option.New(optkeyForceNewFile, true)
}
