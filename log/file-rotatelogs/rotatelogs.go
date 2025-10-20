// Package rotatelogs 是Perl的File-RotateLogs的移植版本
// (https://metacpan.org/release/File-RotateLogs),它允许你在写入输出文件时
// 根据指定的文件名模式自动轮转文件。
package rotatelogs

import (
	"fmt"
	"github.com/Cospk/base-tools/log/file-rotatelogs/internal/fileutil"
	"github.com/lestrrat-go/strftime"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"errors"
)

func (c clockFn) Now() time.Time {
	return c()
}

// New 创建一个新的RotateLogs对象。必须传入日志文件名模式。
// 可以传入可选的`Option`参数
func New(p string, options ...Option) (*RotateLogs, error) {
	globPattern := p
	for _, re := range patternConversionRegexps {
		globPattern = re.ReplaceAllString(globPattern, "*")
	}

	pattern, err := strftime.New(p)
	if err != nil {
		return nil, fmt.Errorf("invalid strftime pattern %w", err)
	}

	var clock Clock = Local
	rotationTime := 24 * time.Hour
	var rotationSize int64
	var rotationCount uint
	var linkName string
	var maxAge time.Duration
	var handler Handler
	var forceNewFile bool

	for _, o := range options {
		switch o.Name() {
		case optkeyClock:
			clock = o.Value().(Clock)
		case optkeyLinkName:
			linkName = o.Value().(string)
		case optkeyMaxAge:
			maxAge = o.Value().(time.Duration)
			if maxAge < 0 {
				maxAge = 0
			}
		case optkeyRotationTime:
			rotationTime = o.Value().(time.Duration)
			if rotationTime < 0 {
				rotationTime = 0
			}
		case optkeyRotationSize:
			rotationSize = o.Value().(int64)
			if rotationSize < 0 {
				rotationSize = 0
			}
		case optkeyRotationCount:
			rotationCount = o.Value().(uint)
		case optkeyHandler:
			handler = o.Value().(Handler)
		case optkeyForceNewFile:
			forceNewFile = true
		}
	}

	if maxAge > 0 && rotationCount > 0 {
		return nil, errors.New("options MaxAge and RotationCount cannot be both set")
	}

	if maxAge == 0 && rotationCount == 0 {
		// 如果两者都为0,给maxAge一个合理的默认值
		maxAge = 7 * 24 * time.Hour
	}

	return &RotateLogs{
		clock:         clock,
		eventHandler:  handler,
		globPattern:   globPattern,
		linkName:      linkName,
		maxAge:        maxAge,
		pattern:       pattern,
		rotationTime:  rotationTime,
		rotationSize:  rotationSize,
		rotationCount: rotationCount,
		forceNewFile:  forceNewFile,
	}, nil
}

// Write 实现io.Writer接口。它写入到当前正在使用的相应文件句柄。
// 如果达到轮转时间,目标文件将自动轮转,并在必要时清除旧文件。
func (rl *RotateLogs) Write(p []byte) (n int, err error) {
	// 防止并发写入
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	out, err := rl.getWriterNolock(false, false)
	if err != nil {
		return 0, fmt.Errorf("failed to acquite target io.Writer %w", err)
	}

	return out.Write(p)
}

// 此操作期间必须加锁
func (rl *RotateLogs) getWriterNolock(bailOnRotateFail, useGenerationalNames bool) (io.Writer, error) {
	generation := rl.generation
	previousFn := rl.curFn

	// 此文件名包含要记录到的"新"文件名,
	// 可能比rl.currentFilename更新
	baseFn := fileutil.GenerateFn(rl.pattern, rl.clock, rl.rotationTime)
	filename := baseFn
	var forceNewFile bool

	fi, err := os.Stat(rl.curFn)
	sizeRotation := false
	if err == nil && rl.rotationSize > 0 && rl.rotationSize <= fi.Size() {
		forceNewFile = true
		sizeRotation = true
	}

	if baseFn != rl.curBaseFn {
		generation = 0
		// 即使这是调用New()后的第一次写入,
		// 检查是否需要创建新文件
		if rl.forceNewFile {
			forceNewFile = true
		}
	} else {
		if !useGenerationalNames && !sizeRotation {
			// 无需操作
			return rl.outFh, nil
		}
		forceNewFile = true
		generation++
	}
	if forceNewFile {
		// 已请求一个新文件。我们不仅使用常规的strftime模式,
		// 而是使用代数名称创建新文件名,如"foo.1"、"foo.2"、"foo.3"等
		var name string
		for {
			if generation == 0 {
				name = filename
			} else {
				name = fmt.Sprintf("%s.%d", filename, generation)
			}
			if _, err := os.Stat(name); err != nil {
				filename = name

				break
			}
			generation++
		}
	}

	fh, err := fileutil.CreateFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to create a new file %s %w", filename, err)
	}

	if err := rl.rotateNolock(filename); err != nil && errors.Is(err, errors.New("the file exists")) {
		err = fmt.Errorf("failed to rotate %w", err)
		if bailOnRotateFail {
			// 轮转失败是一个问题,但仅仅因为无法重命名日志
			// 就停止应用程序并不是一个好主意。
			//
			// 我们只在明确需要时返回此错误(由bailOnRotateFail指定)
			//
			// 但是,我们*需要*在这里关闭`fh`
			if fh != nil { // 可能不会发生,但为了保险起见
				fh.Close()
			}
			return nil, err
		}
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	}

	rl.outFh.Close()
	rl.outFh = fh
	rl.curBaseFn = baseFn
	rl.curFn = filename
	rl.generation = generation

	if h := rl.eventHandler; h != nil {
		go h.Handle(&FileRotatedEvent{
			prev:    previousFn,
			current: filename,
		})
	}

	return fh, nil
}

// CurrentFileName 返回RotateLogs对象当前正在写入的文件名
func (rl *RotateLogs) CurrentFileName() string {
	rl.mutex.RLock()
	defer rl.mutex.RUnlock()

	return rl.curFn
}

var patternConversionRegexps = []*regexp.Regexp{
	regexp.MustCompile(`%[%+A-Za-z]`),
	regexp.MustCompile(`\*+`),
}

type cleanupGuard struct {
	enable bool
	fn     func()
	mutex  sync.Mutex
}

func (g *cleanupGuard) Enable() {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.enable = true
}

func (g *cleanupGuard) Run() {
	g.fn()
}

// Rotate 强制轮转日志文件。如果生成的文件名因文件已存在而冲突,
// 则在日志文件末尾附加数字后缀,形式为".1"、".2"、".3"等。
//
// 此方法可以与信号处理程序结合使用,以模拟在接收到SIGHUP时生成新日志文件的服务器
func (rl *RotateLogs) Rotate() error {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()
	_, err := rl.getWriterNolock(true, true)

	return err
}

func (rl *RotateLogs) rotateNolock(filename string) error {
	lockfn := filename + `_lock`
	fh, err := os.OpenFile(lockfn, os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		// 无法加锁,直接返回
		return err
	}

	var guard cleanupGuard
	guard.fn = func() {
		fh.Close()
		os.Remove(lockfn)
	}
	defer guard.Run()

	if rl.linkName != "" {
		tmpLinkName := filename + `_symlink`

		// 根据目标位置的位置改变链接名称的生成方式。
		// 如果位置直接在主文件名的父目录下,
		// 则创建一个相对路径的符号链接
		linkDest := filename
		linkDir := filepath.Dir(rl.linkName)

		baseDir := filepath.Dir(filename)
		if strings.Contains(rl.linkName, baseDir) {
			tmp, err := filepath.Rel(linkDir, filename)
			if err != nil {
				return fmt.Errorf(`failed to evaluate relative path from %#v to %#v %w`, baseDir, rl.linkName, err)
			}

			linkDest = tmp
		}

		if err := os.Symlink(linkDest, tmpLinkName); err != nil {
			return fmt.Errorf("failed to create new symlink %w", err)
		}

		// 必须存在rl.linkName应该创建的目录
		_, err := os.Stat(linkDir)
		if err != nil { // 假设err != nil表示目录不存在
			if err := os.MkdirAll(linkDir, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s %w", linkDir, err)
			}
		}

		if err := os.Rename(tmpLinkName, rl.linkName); err != nil {
			return fmt.Errorf("failed to rename new symlink %w", err)
		}
	}

	if rl.maxAge <= 0 && rl.rotationCount <= 0 {
		return errors.New("panic: maxAge and rotationCount are both set")
	}

	matches, err := filepath.Glob(rl.globPattern)
	if err != nil {
		return err
	}

	cutoff := rl.clock.Now().Add(-1 * rl.maxAge)

	// linter告诉我预先分配这个...
	toUnlink := make([]string, 0, len(matches))
	for _, path := range matches {
		// 忽略锁文件
		if strings.HasSuffix(path, "_lock") || strings.HasSuffix(path, "_symlink") {
			continue
		}

		fi, err := os.Stat(path)
		if err != nil {
			continue
		}

		fl, err := os.Lstat(path)
		if err != nil {
			continue
		}

		if rl.maxAge > 0 && fi.ModTime().After(cutoff) {
			continue
		}

		if rl.rotationCount > 0 && fl.Mode()&os.ModeSymlink == os.ModeSymlink {
			continue
		}
		toUnlink = append(toUnlink, path)
	}

	if rl.rotationCount > 0 {
		// 仅当文件数量超过rotationCount时才删除
		if rl.rotationCount >= uint(len(toUnlink)) {
			return nil
		}

		toUnlink = toUnlink[:len(toUnlink)-int(rl.rotationCount)]
	}

	if len(toUnlink) <= 0 {
		return nil
	}

	guard.Enable()
	go func() {
		// 在单独的goroutine中删除文件
		for _, path := range toUnlink {
			os.Remove(path)
		}
	}()

	return nil
}

// Close 实现io.Closer接口。如果你对该对象执行了任何写入操作,
// 必须调用此方法。
func (rl *RotateLogs) Close() error {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	if rl.outFh == nil {
		return nil
	}

	rl.outFh.Close()
	rl.outFh = nil

	return nil
}
