package log

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/google/uuid"
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

type logoutput interface {
	Printf(format string, v ...any)
}

type Logger struct {
	logFile      string
	name         string
	logger       logoutput
	file         *os.File
	parent       *Logger
	trueparent   *Logger
	childrenGrab []string
	verbose      bool
	empty        bool
	dir          string
	closed       bool
}

func NewLoggerBase(name string, dir string, verbose bool) Logger {
	logTime := time.Now().UTC().UnixNano()
	fileName := fmt.Sprintf("%s_%d.txt", name, logTime)

	_, err := os.Stat(dir)
	if err != nil {
		os.MkdirAll(dir, 0o777)
	}

	file, _ := os.OpenFile(path.Join(dir, fileName), os.O_CREATE|os.O_RDWR, 0o666)

	lg := Logger{
		logFile:      fileName,
		name:         name,
		logger:       log.New(file, "", 0),
		file:         file,
		parent:       nil,
		childrenGrab: []string{},
		empty:        false,
		verbose:      verbose,
		dir:          dir,
	}

	return lg
}

func NewLogger(name string, parent *Logger) Logger {
	if parent.empty {
		return NewLoggerEmpty()
	}

	trueparent := parent
	for trueparent.parent != nil {
		trueparent = trueparent.parent
	}

	dir := parent.dir
	verbose := parent.verbose

	logid := uuid.New().String()
	fileName := fmt.Sprintf("%s_%s.txt", name, logid)

	_, err := os.Stat(dir)
	if err != nil {
		os.MkdirAll(dir, 0o777)
	}

	file, _ := os.OpenFile(path.Join(dir, fileName), os.O_CREATE|os.O_RDWR, 0o666)

	lg := Logger{
		logFile:    fileName,
		name:       name,
		logger:     log.New(file, "", 0),
		file:       file,
		parent:     parent,
		trueparent: trueparent,
		empty:      false,
		verbose:    verbose,
		dir:        dir,
	}

	return lg
}

type emptyLog struct{}

func (el emptyLog) Printf(format string, v ...any) {}

func NewLoggerEmpty() Logger {
	lg := Logger{
		logFile: "",
		logger:  emptyLog{},
		parent:  nil,
		empty:   true,
	}

	return lg
}

func (l *Logger) EnableVerbose() {
	l.verbose = true
}

func (l *Logger) Append(str string, level LogLevel) string {
	if l.empty {
		return str
	}
	if !l.verbose && level == LEVEL_VERBOSE {
		return str
	}

	if l.parent != nil {
		l.parent.Append(str, level)
	}

	if l.closed {
		return str
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

func (l *Logger) Appendf(format string, level LogLevel, v ...any) string {
	return l.Append(fmt.Sprintf(format, v...), level)
}

func (l *Logger) Close() {
	if l.empty || l.closed {
		return
	}

	if l.trueparent != nil {
		l.trueparent.childrenGrab = append(l.trueparent.childrenGrab, l.file.Name())
	}

	for _, c := range l.childrenGrab {
		b, err := os.ReadFile(c)
		if err != nil {
			l.Appendf("failed to read child log: %s with err=%s", LEVEL_ERROR, c, err)
			continue
		}

		l.file.WriteString(fmt.Sprintf("\n\n\nCHILD LOG: %s\n", c))
		l.file.Write(b)
		os.Remove(c)
	}

	l.file.Close()
	l.closed = true

	if l.parent == nil {
		f, err := os.OpenFile(path.Join(l.dir, "@latest.txt"), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0o666)
		if err != nil {
			return
		}
		defer f.Close()

		b, err := os.ReadFile(path.Join(l.dir, l.logFile))
		if err != nil {
			return
		}

		_, err = f.Write(b)
		if err != nil {
			return
		}
	}
}
