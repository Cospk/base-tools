// Package mcontext 提供context.Context的扩展功能，用于在请求链路中传递操作ID、用户ID、平台信息等
package mcontext

import (
    "context"
    "github.com/Cospk/base-tools/errs"
    constant "github.com/Cospk/base-tools/utils/constants"
)

type ctxKey string

var (
    keyOperationID    ctxKey = ctxKey(constant.OperationID)
    keyOpUserID       ctxKey = ctxKey(constant.OpUserID)
    keyOpUserPlatform ctxKey = ctxKey(constant.OpUserPlatform)
    keyConnID         ctxKey = ctxKey(constant.ConnID)
    keyTriggerID      ctxKey = ctxKey(constant.TriggerID)
    keyRemoteAddr     ctxKey = ctxKey(constant.RemoteAddr)
)

// mapper 定义必需的上下文字段
var mapper = []string{constant.OperationID, constant.OpUserID, constant.OpUserPlatform, constant.ConnID}

// WithOpUserIDContext 设置操作用户ID到context
func WithOpUserIDContext(ctx context.Context, opUserID string) context.Context {
    return context.WithValue(ctx, keyOpUserID, opUserID)
}

// WithOpUserPlatformContext 设置用户平台到context
func WithOpUserPlatformContext(ctx context.Context, platform string) context.Context {
    return context.WithValue(ctx, keyOpUserPlatform, platform)
}

// WithTriggerIDContext 设置触发器ID到context
func WithTriggerIDContext(ctx context.Context, triggerID string) context.Context {
    return context.WithValue(ctx, keyTriggerID, triggerID)
}

// NewCtx 创建新的context并设置operationID，用于链路追踪
func NewCtx(operationID string) context.Context {
    c := context.Background()
    ctx := context.WithValue(c, keyOperationID, operationID)
    return SetOperationID(ctx, operationID)
}

// SetOperationID 设置操作ID
func SetOperationID(ctx context.Context, operationID string) context.Context {
    return context.WithValue(ctx, keyOperationID, operationID)
}

// SetOpUserID 设置操作用户ID
func SetOpUserID(ctx context.Context, opUserID string) context.Context {
    return context.WithValue(ctx, keyOpUserID, opUserID)
}

// SetConnID 设置连接ID
func SetConnID(ctx context.Context, connID string) context.Context {
    return context.WithValue(ctx, keyConnID, connID)
}

// GetOperationID 从context获取操作ID
func GetOperationID(ctx context.Context) string {
    if v, ok := ctx.Value(keyOperationID).(string); ok {
        return v
    }
    s, _ := ctx.Value(constant.OperationID).(string)
    return s
}

// GetOpUserID 从context获取用户ID
func GetOpUserID(ctx context.Context) string {
    if v, ok := ctx.Value(keyOpUserID).(string); ok {
        return v
    }
    s, _ := ctx.Value(constant.OpUserID).(string)
    return s
}

// GetConnID 从context获取连接ID
func GetConnID(ctx context.Context) string {
    if v, ok := ctx.Value(keyConnID).(string); ok {
        return v
    }
    s, _ := ctx.Value(constant.ConnID).(string)
    return s
}

// GetTriggerID 从context获取触发器ID
func GetTriggerID(ctx context.Context) string {
    if v, ok := ctx.Value(keyTriggerID).(string); ok {
        return v
    }
    s, _ := ctx.Value(constant.TriggerID).(string)
    return s
}

// GetOpUserPlatform 从context获取用户平台
func GetOpUserPlatform(ctx context.Context) string {
    if v, ok := ctx.Value(keyOpUserPlatform).(string); ok {
        return v
    }
    s, _ := ctx.Value(constant.OpUserPlatform).(string)
    return s
}

// GetRemoteAddr 从context获取远程地址
func GetRemoteAddr(ctx context.Context) string {
    if v, ok := ctx.Value(keyRemoteAddr).(string); ok {
        return v
    }
    s, _ := ctx.Value(constant.RemoteAddr).(string)
    return s
}

// GetMustCtxInfo 获取必需的上下文信息，如果缺少任何字段则返回错误
func GetMustCtxInfo(ctx context.Context) (operationID, opUserID, platform, connID string, err error) {
    operationID, ok := ctx.Value(keyOperationID).(string)
    if !ok {
        operationID, ok = ctx.Value(constant.OperationID).(string)
    }
    if !ok {
        err = errs.ErrArgs.WrapMsg("ctx missing operationID")
        return
    }
    opUserID, ok1 := ctx.Value(keyOpUserID).(string)
    if !ok1 {
        opUserID, ok1 = ctx.Value(constant.OpUserID).(string)
    }
    if !ok1 {
        err = errs.ErrArgs.WrapMsg("ctx missing opUserID")
        return
    }
    platform, ok2 := ctx.Value(keyOpUserPlatform).(string)
    if !ok2 {
        platform, ok2 = ctx.Value(constant.OpUserPlatform).(string)
    }
    if !ok2 {
        err = errs.ErrArgs.WrapMsg("ctx missing platform")
        return
    }
    if v, ok := ctx.Value(keyConnID).(string); ok {
        connID = v
    } else {
        connID, _ = ctx.Value(constant.ConnID).(string)
    }
    return
}

// GetCtxInfos 获取上下文信息，只有operationID是必需的
func GetCtxInfos(ctx context.Context) (operationID, opUserID, platform, connID string, err error) {
    operationID, ok := ctx.Value(keyOperationID).(string)
    if !ok {
        operationID, ok = ctx.Value(constant.OperationID).(string)
    }
    if !ok {
        err = errs.ErrArgs.WrapMsg("ctx missing operationID")
        return
    }
    if v, ok := ctx.Value(keyOpUserID).(string); ok { opUserID = v } else { opUserID, _ = ctx.Value(constant.OpUserID).(string) }
    if v, ok := ctx.Value(keyOpUserPlatform).(string); ok { platform = v } else { platform, _ = ctx.Value(constant.OpUserPlatform).(string) }
    if v, ok := ctx.Value(keyConnID).(string); ok { connID = v } else { connID, _ = ctx.Value(constant.ConnID).(string) }
    return
}

// WithMustInfoCtx 从字符串切片创建包含必需信息的context
func WithMustInfoCtx(values []string) context.Context {
    ctx := context.Background()
    // 将字符串键映射为类型安全的键
    keyMap := map[string]ctxKey{
        constant.OperationID:    keyOperationID,
        constant.OpUserID:       keyOpUserID,
        constant.OpUserPlatform: keyOpUserPlatform,
        constant.ConnID:         keyConnID,
    }
    for i, v := range values {
        if k, ok := keyMap[mapper[i]]; ok {
            ctx = context.WithValue(ctx, k, v)
        } else {
            ctx = context.WithValue(ctx, mapper[i], v)
        }
    }
    return ctx
}
