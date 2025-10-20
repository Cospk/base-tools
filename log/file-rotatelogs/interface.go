package rotatelogs

import (
	"os"
	"sync"
	"time"

	"github.com/lestrrat-go/strftime"
)

// Handler 事件处理器接口
type Handler interface {
	Handle(Event)
}

// HandlerFunc 事件处理函数类型
type HandlerFunc func(Event)

// Event 事件接口
type Event interface {
	Type() EventType
}

// EventType 事件类型枚举
type EventType int

const (
	InvalidEventType     EventType = iota // 无效事件类型
	FileRotatedEventType                  // 文件轮转事件类型
)

// FileRotatedEvent 文件轮转事件,包含轮转前后的文件名
type FileRotatedEvent struct {
	prev    string // 前一个文件名
	current string // 当前新文件名
}

// RotateLogs 表示一个自动轮转的日志文件,当向其写入时会自动进行轮转
type RotateLogs struct {
	clock         Clock              // 时钟接口,用于获取当前时间
	curFn         string             // 当前文件名
	curBaseFn     string             // 当前基础文件名
	globPattern   string             // 用于匹配日志文件的glob模式
	generation    int                // 代数编号
	linkName      string             // 符号链接名称
	maxAge        time.Duration      // 日志文件最大保留时间
	mutex         sync.RWMutex       // 读写互斥锁
	eventHandler  Handler            // 事件处理器
	outFh         *os.File           // 输出文件句柄
	pattern       *strftime.Strftime // strftime格式模式
	rotationTime  time.Duration      // 轮转时间间隔
	rotationSize  int64              // 轮转文件大小阈值
	rotationCount uint               // 保留的日志文件数量
	forceNewFile  bool               // 是否强制创建新文件
}

// Clock 时钟接口,用于RotateLogs对象确定当前时间
type Clock interface {
	Now() time.Time
}

type clockFn func() time.Time

// UTC 满足Clock接口的对象,返回UTC时区的当前时间
var UTC = clockFn(func() time.Time { return time.Now().UTC() })

// Local 满足Clock接口的对象,返回本地时区的当前时间
var Local = clockFn(time.Now)

// Option 用于向RotateLogs构造函数传递可选参数的接口
type Option interface {
	Name() string
	Value() interface{}
}
