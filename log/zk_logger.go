package log

import (
	"context"
	"fmt"
)

// ZkLogger ZooKeeper的日志记录器实现,将ZK日志输出到统一日志系统
type ZkLogger struct{}

// NewZkLogger 创建一个新的ZooKeeper日志记录器
func NewZkLogger() *ZkLogger {
	return &ZkLogger{}
}

// Printf 实现ZooKeeper日志接口,将格式化的日志输出到Info级别
func (l *ZkLogger) Printf(format string, a ...any) {
	ZInfo(context.Background(), "zookeeper output", "msg", fmt.Sprintf(format, a...))
}
