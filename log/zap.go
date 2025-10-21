package log

import (
	"context"
	"fmt"
	"github.com/Cospk/base-tools/errs"
	rotatelogs "github.com/Cospk/base-tools/log/file-rotatelogs"
	"github.com/Cospk/base-tools/mcontext"
	constant "github.com/Cospk/base-tools/utils/constants"
	"github.com/Cospk/base-tools/utils/stringutil"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	AdaptiveDefaultLevel   = LevelWarn
	AdaptiveErrorCodeLevel = map[int]int{
		errs.ErrInternalServer.Code(): LevelError,
	}
	AsyncWrite = false
)

// LogFormatter 日志格式化器接口,用于自定义日志输出格式
type LogFormatter interface {
	Format() any
}

// 日志级别常量定义,从最严重到最详细
const (
	LevelFatal = iota
	LevelPanic
	LevelError
	LevelWarn
	LevelInfo
	LevelDebug
	LevelDebugWithSQL
)

var (
	pkgLogger   Logger
	osStdout    Logger
	sp          = string(filepath.Separator)
	logLevelMap = map[int]zapcore.Level{
		LevelDebugWithSQL: zapcore.DebugLevel,
		LevelDebug:        zapcore.DebugLevel,
		LevelInfo:         zapcore.InfoLevel,
		LevelWarn:         zapcore.WarnLevel,
		LevelError:        zapcore.ErrorLevel,
		LevelPanic:        zapcore.PanicLevel,
		LevelFatal:        zapcore.FatalLevel,
	}
)

const (
	callDepth   int    = 2
	rotateCount uint   = 1
	hoursPerDay uint   = 24
	logPath     string = "./logs/"
	version     string = "undefined version"
	isSimplify         = false
)

func init() {
	InitLoggerFromConfig(
		"DefaultLogger",
		"DefaultLoggerModule",
		"", "",
		LevelDebug,
		true,
		false,
		logPath,
		rotateCount,
		hoursPerDay,
		version,
		isSimplify,
	)
}

// InitLoggerFromConfig 根据配置初始化基于Zap的日志记录器
func InitLoggerFromConfig(
	loggerPrefixName, moduleName string,
	sdkType, platformName string,
	logLevel int,
	isStdout bool,
	isJson bool,
	logLocation string,
	rotateCount uint,
	rotationTime uint,
	moduleVersion string,
	isSimplify bool,
) error {

	l, err := NewZapLogger(loggerPrefixName, moduleName, sdkType, platformName, logLevel, isStdout, isJson, logLocation, rotateCount, rotationTime, moduleVersion, isSimplify)
	if err != nil {
		return err
	}

	pkgLogger = l.WithCallDepth(callDepth)
	if isJson {
		pkgLogger = pkgLogger.WithName(moduleName)
	}
	return nil
}

// InitConsoleLogger 初始化控制台日志记录器,用于osStdout和osStderr
func InitConsoleLogger(moduleName string,
	logLevel int,
	isJson bool, moduleVersion string) error {
	l, err := NewConsoleZapLogger(moduleName, logLevel, isJson, moduleVersion, os.Stdout)
	if err != nil {
		return err
	}
	osStdout = l.WithCallDepth(callDepth)
	if isJson {
		osStdout = osStdout.WithName(moduleName)
	}

	return nil
}

// ZDebug 记录调试级别日志
func ZDebug(ctx context.Context, msg string, keysAndValues ...any) {
	pkgLogger.Debug(ctx, msg, keysAndValues...)
}

// ZInfo 记录信息级别日志
func ZInfo(ctx context.Context, msg string, keysAndValues ...any) {
	pkgLogger.Info(ctx, msg, keysAndValues...)
}

// ZWarn 记录警告级别日志
func ZWarn(ctx context.Context, msg string, err error, keysAndValues ...any) {
	pkgLogger.Warn(ctx, msg, err, keysAndValues...)
}

// ZError 记录错误级别日志
func ZError(ctx context.Context, msg string, err error, keysAndValues ...any) {
	pkgLogger.Error(ctx, msg, err, keysAndValues...)
}

// ZPanic 记录panic级别日志
func ZPanic(ctx context.Context, msg string, err error, keysAndValues ...any) {
	pkgLogger.Error(ctx, msg, err, keysAndValues...)
}

// ZAdaptive 根据错误码自适应日志级别
func ZAdaptive(ctx context.Context, msg string, err error, keysAndValues ...any) {
	var level int
	if cErr, ok := errs.Unwrap(err).(errs.CodeError); ok {
		level, ok = AdaptiveErrorCodeLevel[cErr.Code()]
		if !ok {
			level = AdaptiveDefaultLevel
		}
	} else {
		level = AdaptiveDefaultLevel
	}
	switch level {
	case LevelDebug:

		pkgLogger.Debug(ctx, msg, appendError(keysAndValues, err)...)
	case LevelInfo:
		pkgLogger.Info(ctx, msg, appendError(keysAndValues, err)...)
	case LevelWarn:
		pkgLogger.Warn(ctx, msg, err, keysAndValues...)
	case LevelError:
		pkgLogger.Error(ctx, msg, err, keysAndValues...)
	case LevelPanic:
		pkgLogger.Error(ctx, msg, err, keysAndValues...)
	default:
	}
}

// CInfo 记录控制台信息日志,输出到osStdout
func CInfo(ctx context.Context, msg string, keysAndValues ...any) {
	if osStdout == nil {
		return
	}
	osStdout.Info(ctx, msg, keysAndValues...)
}

// Flush 刷新所有缓冲的日志到输出
func Flush() {
	if pkgLogger == nil {
		return
	}
	pkgLogger.Flush()
}

// ZapLogger 基于Zap的日志记录器实现,支持日志轮转和上下文信息
type ZapLogger struct {
	zap              *zap.SugaredLogger // Zap的糖化日志器
	level            zapcore.Level      // 日志级别
	moduleName       string             // 模块名称
	moduleVersion    string             // 模块版本
	loggerPrefixName string             // 日志文件前缀名
	rotationTime     time.Duration      // 日志轮转时间间隔
	sdkType          string             // SDK类型
	platformName     string             // 平台名称
	isSimplify       bool               // 是否简化日志输出
}

// NewZapLogger 创建一个新的Zap日志记录器实例
func NewZapLogger(
	loggerPrefixName, moduleName string, sdkType, platformName string,
	logLevel int,
	isStdout bool,
	isJson bool,
	logLocation string,
	rotateCount uint,
	rotationTime uint,
	moduleVersion string,
	isSimplify bool,
) (*ZapLogger, error) {
	zapConfig := zap.Config{
		Level:             zap.NewAtomicLevelAt(logLevelMap[logLevel]),
		DisableStacktrace: true,
	}
	if isJson {
		zapConfig.Encoding = "json"
	} else {
		zapConfig.Encoding = "console"
	}
	zl := &ZapLogger{level: logLevelMap[logLevel],
		moduleName:       moduleName,
		loggerPrefixName: loggerPrefixName,
		rotationTime:     time.Duration(rotationTime) * time.Hour,
		moduleVersion:    moduleVersion,
		sdkType:          sdkType,
		platformName:     platformName,
		isSimplify:       isSimplify,
	}
	opts, err := zl.cores(isStdout, isJson, logLocation, rotateCount)
	if err != nil {
		return nil, err
	}
	l, err := zapConfig.Build(opts)
	if err != nil {
		return nil, err
	}
	zl.zap = l.Sugar()
	return zl, nil
}

// NewConsoleZapLogger 创建一个输出到控制台的Zap日志记录器
func NewConsoleZapLogger(
	moduleName string,
	logLevel int,
	isJson bool,
	moduleVersion string,
	outPut *os.File) (*ZapLogger, error) {
	zapConfig := zap.Config{
		Level:             zap.NewAtomicLevelAt(logLevelMap[logLevel]),
		DisableStacktrace: true,
	}
	if isJson {
		zapConfig.Encoding = "json"
	} else {
		zapConfig.Encoding = "console"
	}
	zl := &ZapLogger{level: logLevelMap[logLevel], moduleName: moduleName, moduleVersion: moduleVersion}
	opts, err := zl.consoleCores(outPut, isJson)
	if err != nil {
		return nil, err
	}
	l, err := zapConfig.Build(opts)
	if err != nil {
		return nil, err
	}
	zl.zap = l.Sugar()
	return zl, nil
}

func (l *ZapLogger) cores(isStdout bool, isJson bool, logLocation string, rotateCount uint) (zap.Option, error) {
	c := zap.NewProductionEncoderConfig()
	c.EncodeTime = l.timeEncoder
	c.EncodeDuration = zapcore.StringDurationEncoder
	c.MessageKey = "msg"
	c.LevelKey = "level"
	c.TimeKey = "time"
	c.CallerKey = "caller"
	c.NameKey = "logger"
	var fileEncoder zapcore.Encoder

	if isJson {
		c.EncodeLevel = zapcore.CapitalLevelEncoder
		fileEncoder = zapcore.NewJSONEncoder(c)
		fileEncoder.AddInt("PID", os.Getpid())
		fileEncoder.AddString("version", l.moduleVersion)
	} else {
		c.EncodeLevel = l.capitalColorLevelEncoder
		c.EncodeCaller = l.customCallerEncoder
		fileEncoder = zapcore.NewConsoleEncoder(c)
	}

	fileEncoder = &alignEncoder{Encoder: fileEncoder}
	writer, err := l.getWriter(logLocation, rotateCount)
	if err != nil {
		return nil, err
	}

	if !isStdout && AsyncWrite {
		writer = &zapcore.BufferedWriteSyncer{
			WS:            writer,
			FlushInterval: time.Second * 2,
			Size:          1024 * 512,
		}
	}

	var cores []zapcore.Core
	if logLocation != "" {
		cores = []zapcore.Core{
			zapcore.NewCore(fileEncoder, writer, zap.NewAtomicLevelAt(l.level)),
		}
	}

	if isStdout {
		cores = append(cores, zapcore.NewCore(fileEncoder, zapcore.Lock(os.Stdout), zap.NewAtomicLevelAt(l.level)))
	}

	return zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return zapcore.NewTee(cores...)
	}), nil
}

func (l *ZapLogger) consoleCores(outPut *os.File, isJson bool) (zap.Option, error) {
	c := zap.NewProductionEncoderConfig()
	c.EncodeTime = l.timeEncoder
	c.EncodeDuration = zapcore.StringDurationEncoder
	c.MessageKey = "msg"
	c.LevelKey = "level"
	c.TimeKey = "time"
	c.CallerKey = "caller"
	c.NameKey = "logger"
	var fileEncoder zapcore.Encoder
	if isJson {
		c.EncodeLevel = zapcore.CapitalLevelEncoder
		fileEncoder = zapcore.NewJSONEncoder(c)
		fileEncoder.AddInt("PID", os.Getpid())
		fileEncoder.AddString("version", l.moduleVersion)
	} else {
		c.EncodeLevel = l.capitalColorLevelEncoder
		c.EncodeCaller = l.customCallerEncoder
		fileEncoder = zapcore.NewConsoleEncoder(c)
	}
	var cores []zapcore.Core
	cores = append(cores, zapcore.NewCore(fileEncoder, zapcore.Lock(outPut), zap.NewAtomicLevelAt(l.level)))

	return zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return zapcore.NewTee(cores...)
	}), nil
}

func (l *ZapLogger) customCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	if l.sdkType != "" && l.platformName != "" {
		fixedLength := 50
		sdkPlatform := fmt.Sprintf("[%s/%s]", l.sdkType, l.platformName)
		sdkPlatformFormatted := stringutil.FormatString(sdkPlatform, fixedLength, true)
		enc.AppendString(sdkPlatformFormatted)

		trimmedPath := caller.TrimmedPath()
		trimmedPath = "[" + trimmedPath + "]"
		s := stringutil.FormatString(trimmedPath, fixedLength, true)
		enc.AppendString(s)
	} else {
		fixedLength := 50
		trimmedPath := caller.TrimmedPath()
		trimmedPath = "[" + trimmedPath + "]"
		s := stringutil.FormatString(trimmedPath, fixedLength, true)
		enc.AppendString(s)
	}
}

// SDKLog SDK专用日志方法,支持原生调用信息记录
func SDKLog(ctx context.Context, logLevel int, file string, line int, msg string, err error, keysAndValues []any) {
	nativeCallerKey := "native_caller"
	nativeCaller := fmt.Sprintf("[%s:%d]", file, line)

	kv := []any{nativeCallerKey, nativeCaller}
	kv = append(kv, keysAndValues...)

	switch logLevel {
	case LevelDebugWithSQL:
		ZDebug(ctx, msg, kv...)
	case LevelInfo:
		ZInfo(ctx, msg, kv...)
	case LevelWarn:
		ZWarn(ctx, msg, err, kv...)
	case LevelError:
		ZError(ctx, msg, err, kv...)
	}
}

func (l *ZapLogger) timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	layout := "2006-01-02 15:04:05.000"
	type appendTimeEncoder interface {
		AppendTimeLayout(time.Time, string)
	}
	if enc, ok := enc.(appendTimeEncoder); ok {
		enc.AppendTimeLayout(t, layout)
		return
	}
	enc.AppendString(t.Format(layout))
}

func (l *ZapLogger) getWriter(logLocation string, rorateCount uint) (zapcore.WriteSyncer, error) {
	var path string
	if l.rotationTime%(time.Hour*time.Duration(hoursPerDay)) == 0 {
		path = logLocation + sp + l.loggerPrefixName + ".%Y-%m-%d"
	} else if l.rotationTime%time.Hour == 0 {
		path = logLocation + sp + l.loggerPrefixName + ".%Y-%m-%d_%H"
	} else {
		path = logLocation + sp + l.loggerPrefixName + ".%Y-%m-%d_%H_%M_%S"
	}
	logf, err := rotatelogs.New(path,
		rotatelogs.WithRotationCount(rorateCount),
		rotatelogs.WithRotationTime(l.rotationTime),
	)
	if err != nil {
		return nil, err
	}
	return zapcore.AddSync(logf), nil
}

func (l *ZapLogger) capitalColorLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	s, ok := _levelToCapitalColorString[level]
	if !ok {
		s = _unknownLevelColor[zapcore.ErrorLevel]
	}
	pid := stringutil.FormatString(fmt.Sprintf("[PID:%d]", os.Getpid()), 15, true)
	color := _levelToColor[level]
	enc.AppendString(s)
	enc.AppendString(color.Add(pid))
	if l.moduleName != "" {
		moduleName := stringutil.FormatString(l.moduleName, 25, true)
		enc.AppendString(color.Add(moduleName))
	}
	if l.moduleVersion != "" {
		moduleVersion := stringutil.FormatString(fmt.Sprintf("[%s]", l.moduleVersion), 30, true)
		enc.AppendString(moduleVersion)
	}
}

// ToZap 返回底层的Zap SugaredLogger实例
func (l *ZapLogger) ToZap() *zap.SugaredLogger {
	return l.zap
}

// Debug 记录调试级别日志
func (l *ZapLogger) Debug(ctx context.Context, msg string, keysAndValues ...any) {
	if l.level > zapcore.DebugLevel {
		return
	}
	keysAndValues = l.kvAppend(ctx, keysAndValues)
	l.zap.Debugw(msg, keysAndValues...)
}

// Info 记录信息级别日志
func (l *ZapLogger) Info(ctx context.Context, msg string, keysAndValues ...any) {
	if l.level > zapcore.InfoLevel {
		return
	}
	keysAndValues = l.kvAppend(ctx, keysAndValues)
	l.zap.Infow(msg, keysAndValues...)
}

// Warn 记录警告级别日志
func (l *ZapLogger) Warn(ctx context.Context, msg string, err error, keysAndValues ...any) {
	if l.level > zapcore.WarnLevel {
		return
	}
	keysAndValues = l.kvAppend(ctx, appendError(keysAndValues, err))
	l.zap.Warnw(msg, keysAndValues...)
}

// Error 记录错误级别日志
func (l *ZapLogger) Error(ctx context.Context, msg string, err error, keysAndValues ...any) {
	if l.level > zapcore.ErrorLevel {
		return
	}
	keysAndValues = l.kvAppend(ctx, appendError(keysAndValues, err))
	l.zap.Errorw(msg, keysAndValues...)
}

// Panic 记录panic级别日志
func (l *ZapLogger) Panic(ctx context.Context, msg string, err error, keysAndValues ...any) {
	if l.level > zapcore.PanicLevel {
		return
	}
	keysAndValues = l.kvAppend(ctx, appendError(keysAndValues, err))
	l.zap.Panicw(msg, keysAndValues...)
}

// kvAppend 向键值对切片追加上下文信息(operationID, userID等)
func (l *ZapLogger) kvAppend(ctx context.Context, keysAndValues []any) []any {
    if ctx == nil {
        return keysAndValues
    }
    operationID := mcontext.GetOperationID(ctx)
    opUserID := mcontext.GetOpUserID(ctx)
    connID := mcontext.GetConnID(ctx)
    triggerID := mcontext.GetTriggerID(ctx)
    opUserPlatform := mcontext.GetOpUserPlatform(ctx)
    remoteAddr := mcontext.GetRemoteAddr(ctx)

    // 兼容测试或外部代码使用原始字符串键注入上下文的场景
    if operationID == "" {
        if v, ok := ctx.Value("OperationID").(string); ok {
            operationID = v
        }
    }
    if opUserID == "" {
        if v, ok := ctx.Value("OpUserID").(string); ok {
            opUserID = v
        }
    }
    if connID == "" {
        if v, ok := ctx.Value("ConnID").(string); ok {
            connID = v
        }
    }

	if l.isSimplify {
		if len(keysAndValues)%2 == 0 {
			for i := 1; i < len(keysAndValues); i += 2 {
				if val, ok := keysAndValues[i].(LogFormatter); ok && val != nil {
					keysAndValues[i] = val.Format()
				}
			}
		} else {
			ZError(ctx, "keysAndValues length is not even", errs.ErrInternalServer.Wrap())
		}
	}

    if opUserID != "" {
        keysAndValues = append([]any{constant.OpUserID, opUserID}, keysAndValues...)
    }
    if operationID != "" {
        keysAndValues = append([]any{constant.OperationID, operationID}, keysAndValues...)
	}
	if connID != "" {
		keysAndValues = append([]any{constant.ConnID, connID}, keysAndValues...)
	}
	if triggerID != "" {
		keysAndValues = append([]any{constant.TriggerID, triggerID}, keysAndValues...)
	}
	if opUserPlatform != "" {
		keysAndValues = append([]any{constant.OpUserPlatform, opUserPlatform}, keysAndValues...)
	}
	if remoteAddr != "" {
		keysAndValues = append([]any{constant.RemoteAddr, remoteAddr}, keysAndValues...)
	}
	return keysAndValues
}

// WithValues 返回一个附加了键值对的新Logger实例
func (l *ZapLogger) WithValues(keysAndValues ...any) Logger {
	dup := *l
	dup.zap = l.zap.With(keysAndValues...)
	return &dup
}

// WithName 返回一个带有指定名称的新Logger实例
func (l *ZapLogger) WithName(name string) Logger {
	dup := *l
	dup.zap = l.zap.Named(name)
	return &dup
}

// WithCallDepth 返回一个调整了调用深度的新Logger实例,用于正确显示调用位置
func (l *ZapLogger) WithCallDepth(depth int) Logger {
	dup := *l
	dup.zap = l.zap.WithOptions(zap.AddCallerSkip(depth))
	return &dup
}

// Flush 刷新日志缓冲区,将所有待写入的日志写入到目标
func (l *ZapLogger) Flush() {
	if err := l.zap.Sync(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "failed to flush zap logger", err)
	}
}

// appendError 将错误信息追加到键值对切片中
func appendError(keysAndValues []any, err error) []any {
	if err != nil {
		keysAndValues = append(keysAndValues, "error", err.Error())
	}
	return keysAndValues
}
