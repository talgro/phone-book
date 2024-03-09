package stdout

import (
	"context"
	"fmt"
	"log"
	"os"
)

const (
	LevelDebug   = "debug"
	LevelInfo    = "info"
	LevelError   = "error"
	LevelWarning = "warning"
)

type logger struct {
	logger *log.Logger
}

func NewLogger() *logger {
	return &logger{
		logger: log.New(os.Stdout, "", log.LstdFlags),
	}
}

func (l *logger) Info(_ context.Context, msg string, keyvals ...interface{}) {
	l.log(LevelInfo, msg, nil, keyvals...)
}

func (l *logger) Error(_ context.Context, err error, keyvals ...interface{}) {
	l.log(LevelError, "", err, keyvals...)
}

func (l *logger) Warning(_ context.Context, err error, keyvals ...interface{}) {
	l.log(LevelWarning, "", err, keyvals...)
}

func (l *logger) Debug(_ context.Context, msg string, keyvals ...interface{}) {
	l.log(LevelDebug, msg, nil, keyvals...)
}

func (l *logger) log(level, msg string, err error, keyvals ...interface{}) {
	var logMsg string
	if err != nil {
		logMsg = fmt.Sprintf("%s: %v", msg, err)
	} else {
		logMsg = msg
	}
	l.logger.Printf("[%s] %s %v\n", level, logMsg, keyvals)
}
