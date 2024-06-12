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
)

type Logger struct {
	logFile string
	logger  *log.Logger
}

func NewLogger(dir string) Logger {
	logTime := time.Now().UTC().UnixMilli()

	lg := Logger{
		logFile: fmt.Sprintf("%d.txt", logTime),
		logger:  log.Default(),
	}

	_, err := os.Stat(dir)
	if err != nil {
		os.MkdirAll(dir, 0o666)
	}

	file, _ := os.OpenFile(path.Join(dir, lg.logFile), os.O_CREATE, 0o666)

	lg.logger.SetOutput(file)

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
