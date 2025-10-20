// Package mcontext 提供context.Context的扩展功能，用于在请求链路中传递操作ID、用户ID、平台信息等
package mcontext

import (
	"context"
	"github.com/Cospk/base-tools/errs"
	constant "github.com/Cospk/base-tools/utils/constants"
)

// mapper 定义必需的上下文字段
var mapper = []string{constant.OperationID, constant.OpUserID, constant.OpUserPlatform, constant.ConnID}

// WithOpUserIDContext 设置操作用户ID到context
func WithOpUserIDContext(ctx context.Context, opUserID string) context.Context {
	return context.WithValue(ctx, constant.OpUserID, opUserID)
}

// WithOpUserPlatformContext 设置用户平台到context
func WithOpUserPlatformContext(ctx context.Context, platform string) context.Context {
	return context.WithValue(ctx, constant.OpUserPlatform, platform)
}

// WithTriggerIDContext 设置触发器ID到context
func WithTriggerIDContext(ctx context.Context, triggerID string) context.Context {
	return context.WithValue(ctx, constant.TriggerID, triggerID)
}

// NewCtx 创建新的context并设置operationID，用于链路追踪
func NewCtx(operationID string) context.Context {
	c := context.Background()
	ctx := context.WithValue(c, constant.OperationID, operationID)
	return SetOperationID(ctx, operationID)
}

// SetOperationID 设置操作ID
func SetOperationID(ctx context.Context, operationID string) context.Context {
	return context.WithValue(ctx, constant.OperationID, operationID)
}

// SetOpUserID 设置操作用户ID
func SetOpUserID(ctx context.Context, opUserID string) context.Context {
	return context.WithValue(ctx, constant.OpUserID, opUserID)
}

// SetConnID 设置连接ID
func SetConnID(ctx context.Context, connID string) context.Context {
	return context.WithValue(ctx, constant.ConnID, connID)
}

// GetOperationID 从context获取操作ID
func GetOperationID(ctx context.Context) string {
	s, _ := ctx.Value(constant.OperationID).(string)
	return s
}

// GetOpUserID 从context获取用户ID
func GetOpUserID(ctx context.Context) string {
	s, _ := ctx.Value(constant.OpUserID).(string)
	return s
}

// GetConnID 从context获取连接ID
func GetConnID(ctx context.Context) string {
	s, _ := ctx.Value(constant.ConnID).(string)
	return s
}

// GetTriggerID 从context获取触发器ID
func GetTriggerID(ctx context.Context) string {
	s, _ := ctx.Value(constant.TriggerID).(string)
	return s
}

// GetOpUserPlatform 从context获取用户平台
func GetOpUserPlatform(ctx context.Context) string {
	s, _ := ctx.Value(constant.OpUserPlatform).(string)
	return s
}

// GetRemoteAddr 从context获取远程地址
func GetRemoteAddr(ctx context.Context) string {
	s, _ := ctx.Value(constant.RemoteAddr).(string)
	return s
}

// GetMustCtxInfo 获取必需的上下文信息，如果缺少任何字段则返回错误
func GetMustCtxInfo(ctx context.Context) (operationID, opUserID, platform, connID string, err error) {
	operationID, ok := ctx.Value(constant.OperationID).(string)
	if !ok {
		err = errs.ErrArgs.WrapMsg("ctx missing operationID")
		return
	}
	opUserID, ok1 := ctx.Value(constant.OpUserID).(string)
	if !ok1 {
		err = errs.ErrArgs.WrapMsg("ctx missing opUserID")
		return
	}
	platform, ok2 := ctx.Value(constant.OpUserPlatform).(string)
	if !ok2 {
		err = errs.ErrArgs.WrapMsg("ctx missing platform")
		return
	}
	connID, _ = ctx.Value(constant.ConnID).(string)
	return
}

// GetCtxInfos 获取上下文信息，只有operationID是必需的
func GetCtxInfos(ctx context.Context) (operationID, opUserID, platform, connID string, err error) {
	operationID, ok := ctx.Value(constant.OperationID).(string)
	if !ok {
		err = errs.ErrArgs.WrapMsg("ctx missing operationID")
		return
	}
	opUserID, _ = ctx.Value(constant.OpUserID).(string)
	platform, _ = ctx.Value(constant.OpUserPlatform).(string)
	connID, _ = ctx.Value(constant.ConnID).(string)
	return
}

// WithMustInfoCtx 从字符串切片创建包含必需信息的context
func WithMustInfoCtx(values []string) context.Context {
	ctx := context.Background()
	for i, v := range values {
		ctx = context.WithValue(ctx, mapper[i], v)
	}
	return ctx
}
