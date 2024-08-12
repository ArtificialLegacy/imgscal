package log

import (
	"fmt"
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
	LEVEL_IMPORTANT
	LEVEL_VERBOSE
	LEVEL_SYSTEM
)

const LOG_DIR = "./log"

type logoutput interface {
	Printf(format string, v ...any)
}

type Logger struct {
	logFile      string
	name         string
	logger       logoutput
	file         *os.File
	Parent       *Logger
	childrenGrab []string
	verbose      bool
}

func NewLogger(name string) Logger {
	logTime := time.Now().UTC().UnixMilli()
	fileName := fmt.Sprintf("%s_%d.txt", name, logTime)

	_, err := os.Stat(LOG_DIR)
	if err != nil {
		os.MkdirAll(LOG_DIR, 0o777)
	}

	file, _ := os.OpenFile(path.Join(LOG_DIR, fileName), os.O_CREATE|os.O_RDWR, 0o666)

	lg := Logger{
		logFile:      fileName,
		name:         name,
		logger:       log.New(file, "", 0),
		file:         file,
		Parent:       nil,
		childrenGrab: []string{},
	}

	return lg
}

type emptyLog struct{}

func (el emptyLog) Printf(format string, v ...any) {}

func NewLoggerEmpty() Logger {
	lg := Logger{
		logFile: "",
		logger:  emptyLog{},
		Parent:  nil,
	}

	return lg
}

func (l Logger) EnableVerbose() Logger {
	l.verbose = true
	return l
}

func (l *Logger) Append(str string, level LogLevel) string {
	if !l.verbose && level == LEVEL_VERBOSE {
		return str
	}

	if l.Parent != nil {
		l.Parent.Append(str, level)
	}

	logTime := time.Now().Format(time.ANSIC)

	prefix := ""

	switch level {
	case LEVEL_VERBOSE:
		fallthrough
	case LEVEL_INFO:
		prefix = "# INFO"
	case LEVEL_WARN:
		prefix = "? WARN"
	case LEVEL_ERROR:
		prefix = "! ERROR"
	case LEVEL_IMPORTANT:
		prefix = "!! IMPORTANT"
	case LEVEL_SYSTEM:
		prefix = "# SYSTEM"
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

	if l.Parent == nil {
		f, err := os.OpenFile(path.Join(LOG_DIR, "@latest.txt"), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0o666)
		if err != nil {
			return
		}

		b, err := os.ReadFile(path.Join(LOG_DIR, l.logFile))
		if err != nil {
			return
		}

		_, err = f.Write(b)
		if err != nil {
			return
		}
	}
}
