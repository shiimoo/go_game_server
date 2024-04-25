package slog

import (
	"fmt"
	"os"
	"sync"
	"time"
)

/* Log Level */
type LogLv int

const (
	LOG = LogLv(iota)

	DEBUG
	INFO
	WARNNING
	ERORR
	FATAL
)

var logLvStr = map[LogLv]string{
	DEBUG:    "DEBUG",
	INFO:     "INFO",
	WARNNING: "WARNNING",
	ERORR:    "ERORR",
	FATAL:    "FATAL",
}

func tranLogLv(lv LogLv) string {
	if tran, ok := logLvStr[lv]; ok {
		return tran + " | "
	}
	return ""
}

/* bufferPool : Avoid repeated creation and release*/

var _bufferPool = sync.Pool{New: func() any { return new([]byte) }}

func getBuffer() *[]byte {
	bs := _bufferPool.Get().(*[]byte)
	*bs = (*bs)[:0]
	return bs
}

func putBuffer(p *[]byte) {
	// impose restrictions on base slice
	if cap(*p) > 64<<10 {
		*p = nil
	}
	_bufferPool.Put(p)
}

/* Logger */

type Logger struct {
	name string
	// out sava lock
	outMu sync.Mutex

	Prefix func() string
}

func (l *Logger) GetName() string {
	return l.name
}

func (l *Logger) Output(level LogLv, msgSplit func([]byte) []byte) {
	buf := getBuffer()
	defer putBuffer(buf)
	// add prefix
	*buf = append(*buf, l.Prefix()...)
	// add loglv
	*buf = append(*buf, tranLogLv(level)...)
	// add content
	*buf = msgSplit(*buf)
	*buf = append(*buf, '\n')
	l.outMu.Lock()
	defer l.outMu.Unlock()

	os.Stdout.Write(*buf) // console output
	logPath := defaultLogPath(l)
	if logPath != "" {
		fp, _ := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND, 0644)
		if fp != nil {
			fp.Write(*buf)
		}
	}
}

// SetPrefix
func (l *Logger) SetPrefix(f func() string) {
	l.Prefix = f
}

/* ----- log 日志记录相关方法 ----- */

func (l *Logger) Log(v ...any) {
	l.Output(LOG, func(b []byte) []byte {
		return fmt.Append(b, v...)
	})
}

func (l *Logger) Logf(format string, v ...any) {
	l.Output(LOG, func(b []byte) []byte {
		return fmt.Appendf(b, format, v...)
	})
}

func (l *Logger) Debug(v ...any) {
	l.Output(DEBUG, func(b []byte) []byte {
		return fmt.Append(b, v...)
	})
}

func (l *Logger) Debugf(format string, v ...any) {
	l.Output(DEBUG, func(b []byte) []byte {
		return fmt.Appendf(b, format, v...)
	})
}

func (l *Logger) Info(v ...any) {
	l.Output(INFO, func(b []byte) []byte {
		return fmt.Append(b, v...)
	})
}

func (l *Logger) Infof(format string, v ...any) {
	l.Output(INFO, func(b []byte) []byte {
		return fmt.Appendf(b, format, v...)
	})
}

func (l *Logger) Warn(v ...any) {
	l.Output(WARNNING, func(b []byte) []byte {
		return fmt.Append(b, v...)
	})
}

func (l *Logger) Warnf(format string, v ...any) {
	l.Output(WARNNING, func(b []byte) []byte {
		return fmt.Appendf(b, format, v...)
	})
}

func (l *Logger) Error(v ...any) {
	l.Output(ERORR, func(b []byte) []byte {
		return fmt.Append(b, v...)
	})
}

func (l *Logger) Errorf(format string, v ...any) {
	l.Output(ERORR, func(b []byte) []byte {
		return fmt.Appendf(b, format, v...)
	})
}

func (l *Logger) Fatal(v ...any) {
	l.Output(FATAL, func(b []byte) []byte {
		return fmt.Append(b, v...)
	})
	os.Exit(1)
}

func (l *Logger) Fatalf(format string, v ...any) {
	l.Output(FATAL, func(b []byte) []byte {
		return fmt.Appendf(b, format, v...)
	})
	os.Exit(1)
}

func NewLogger(key string) *Logger {
	l := new(Logger)
	l.name = key
	l.Prefix = DefaultPrefix
	loggerMgr.Store(key, l)
	return l
}

func DefaultPrefix() string {
	now := time.Now()
	return fmt.Sprintf(
		"[%s] ",
		now.Format("2006-01-02 15:04:05"),
	)
}

// *.log file path, if not, it's ""
var defaultLogPath = func(*Logger) string {
	return ""
}

func SetLogPath(f func(*Logger) string) {
	defaultLogPath = f
}
