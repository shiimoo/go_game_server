package slog

import (
	"fmt"
	"os"
	"sync"
)

/* Log Level */
type LogLv int

const (
	Normal = LogLv(iota)

	Debug
	Info
	Warn
	Error
	Fatal
)

var logLvStr = map[LogLv]string{
	Debug: "Debug",
	Info:  "Info",
	Warn:  "Warn",
	Error: "Error",
	Fatal: "Fatal",
}

func tranLogLv(lv LogLv) string {
	return logLvStr[lv]
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
	// *.log file path getfunc, if not, it's ""
	LogPath func() string
	// out sava lock
	outMu sync.Mutex

	Prefix func() []byte
}

func (l *Logger) Output(level LogLv, msgSplit func([]byte) []byte) {
	buf := getBuffer()
	defer putBuffer(buf)
	// add prefix
	*buf = append(*buf, l.Prefix()...)
	// add content
	*buf = msgSplit(*buf)
	*buf = append(*buf, '\n')
	l.outMu.Lock()
	defer l.outMu.Unlock()

	os.Stdout.Write(*buf) // console output
	logPath := l.LogPath()
	if logPath != "" {
		fp, _ := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND, 0644)
		if fp != nil {
			fp.Write(*buf)
		}
	}
}

// SetPrefix
func (l *Logger) SetPrefix(f func() []byte) {
	l.Prefix = f
}

// SetLogPath
func (l *Logger) SetLogPath(f func() string) {
	l.LogPath = f
}

/* log 日志记录相关方法 */

func (l *Logger) Log(v ...any) {
	l.Output(Normal, func(b []byte) []byte {
		return fmt.Append(b, v...)
	})
}

// func (l *Logger) Logf(v ...any) {
// 	l.Output(Normal, func(b []byte) []byte {
// 		return fmt.Append(b, v...)
// 	})
// }

func NewLogger() *Logger {
	l := new(Logger)
	l.LogPath = func() string {
		return ""
	}
	l.Prefix = func() []byte {
		return []byte{}
	}
	return l
}
