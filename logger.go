package alpacamux

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	levelCritical = 50
	levelError    = 40
	levelWarning  = 30
	levelInfo     = 20
	levelDebug    = 10
	levelNotset   = 0
)

// a wrapper of fmt.Printf
type Logger struct {
	Level int
}

var log = &Logger{}
var Log = log

func (l *Logger) printf(level string, format string, a ...interface{}) {
	now := time.Now().Format("2006-01-02 15:04:05.000")
	prefix := now + " [" + level + "] "

	_, file, no, ok := runtime.Caller(2)
	if ok {
		prefix = prefix + "[" + filepath.Base(file) + ":" + strconv.Itoa(no) + "] "
	}

	fmt.Printf(prefix+format+"\n", a...)
}

func (l *Logger) SetLevel(level string) {
	if strings.EqualFold(level, "CRITICAL") {
		l.Level = levelCritical
	} else if strings.EqualFold(level, "ERROR") {
		l.Level = levelError
	} else if strings.EqualFold(level, "WARNING") {
		l.Level = levelWarning
	} else if strings.EqualFold(level, "INFO") {
		l.Level = levelInfo
	} else if strings.EqualFold(level, "DEBUG") {
		l.Level = levelDebug
	} else {
		l.Warning("Invalid level setter: %v, use INFO by default.", level)
		l.Level = levelInfo
	}
}

func (l *Logger) Critical(format string, a ...interface{}) {
	if l.Level <= levelCritical {
		l.printf("CRITICAL", format, a...)
	}
}

func (l *Logger) Error(format string, a ...interface{}) {
	if l.Level <= levelError {
		l.printf("ERROR", format, a...)
	}
}

func (l *Logger) Warning(format string, a ...interface{}) {
	if l.Level <= levelWarning {
		l.printf("WARNING", format, a...)
	}
}

func (l *Logger) Info(format string, a ...interface{}) {
	if l.Level <= levelInfo {
		l.printf("INFO", format, a...)
	}
}

func (l *Logger) Debug(format string, a ...interface{}) {
	if l.Level <= levelDebug {
		l.printf("DEBUG", format, a...)
	}
}
