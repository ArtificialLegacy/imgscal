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

const LOG_DIR = "./log"

type logoutput interface {
	Printf(format string, v ...any)
	SetFlags(flag int)
	SetOutput(w io.Writer)
}

type Logger struct {
	logFile      string
	name         string
	logger       logoutput
	file         *os.File
	Parent       *Logger
	childrenGrab []string
}

func NewLogger(name string) Logger {
	logTime := time.Now().UTC().UnixMilli()
	fileName := fmt.Sprintf("%s_%d.txt", name, logTime)

	file, _ := os.OpenFile(path.Join(LOG_DIR, fileName), os.O_CREATE, 0o666)

	lg := Logger{
		logFile:      fileName,
		name:         name,
		logger:       log.New(file, "", 0),
		file:         file,
		Parent:       nil,
		childrenGrab: []string{},
	}

	_, err := os.Stat(LOG_DIR)
	if err != nil {
		os.MkdirAll(LOG_DIR, 0o666)
	}

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
		Parent:  nil,
	}

	return lg
}

func (l *Logger) Append(str string, level LogLevel) string {
	if l.Parent != nil {
		l.Parent.Append(str, level)
	}

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

func (l *Logger) Close() {
	if l.Parent != nil {
		l.Parent.childrenGrab = append(l.Parent.childrenGrab, l.file.Name())
	}

	for _, c := range l.childrenGrab {
		b, err := os.ReadFile(c)
		if err != nil {
			l.Append(fmt.Sprintf("failed to read child log: %s with err=%s", c, err), LEVEL_ERROR)
			continue
		}

		l.file.WriteString(fmt.Sprintf("\n\n\nCHILD LOG: %s\n", c))
		l.file.Write(b)
		os.Remove(c)
	}

	l.file.Close()
}
