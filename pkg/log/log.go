package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"time"
)

type LogLevel int

const (
	LEVEL_INFO LogLevel = iota
	LEVEL_WARN
	LEVEL_ERROR
)

type logoutput interface {
	Printf(format string, v ...any)
	SetFlags(flag int)
	SetOutput(w io.Writer)
}

type Logger struct {
	logFile string
	logger  logoutput
}

func NewLogger(dir string) Logger {
	logTime := time.Now().UTC().UnixMilli()

	lg := Logger{
		logFile: fmt.Sprintf("%d.txt", logTime),
		logger:  log.Default(),
	}

	lg.logger.SetFlags(0)

	_, err := os.Stat(dir)
	if err != nil {
		os.MkdirAll(dir, 0o666)
	}

	file, _ := os.OpenFile(path.Join(dir, lg.logFile), os.O_CREATE, 0o666)

	lg.logger.SetOutput(file)

	return lg
}

type emptyLog struct {
}

func (el emptyLog) Printf(format string, v ...any) {}
func (el emptyLog) SetFlags(flag int)              {}
func (el emptyLog) SetOutput(w io.Writer)          {}

func NewLoggerEmpty() Logger {
	lg := Logger{
		logFile: "",
		logger:  emptyLog{},
	}

	return lg
}

func (l *Logger) Append(str string, level LogLevel) string {
	logTime := time.Now().Format(time.ANSIC)

	prefix := ""

	switch level {
	case LEVEL_INFO:
		prefix = "# INFO"
	case LEVEL_WARN:
		prefix = "? WARN"
	case LEVEL_ERROR:
		prefix = "! ERROR"
	}

	l.logger.Printf("%s: [%s] > '%s'\n", prefix, logTime, str)
	return str
}
