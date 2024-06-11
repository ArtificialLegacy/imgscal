package log

import (
	"bytes"
	"fmt"
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

type Logger struct {
	logFile string
	buff    bytes.Buffer
}

func NewLogger() Logger {
	logTime := time.Now().UTC().UnixMilli()

	return Logger{
		logFile: fmt.Sprintf("%d.txt", logTime),
	}
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

	l.buff.WriteString(fmt.Sprintf("%s: [%s] > '%s'\n", prefix, logTime, str))
	return str
}

func (l *Logger) Dump(dir string) error {
	_, err := os.Stat(dir)
	if err != nil {
		os.MkdirAll(dir, 0o666)
	}

	file, err := os.OpenFile(path.Join(dir, l.logFile), os.O_CREATE, 0o666)
	if err != nil {
		return err
	}

	file.Write(l.buff.Bytes())

	return nil
}
