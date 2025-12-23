package utils

import (
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
)

type logLevel int

const (
	DEBUG logLevel = iota
	ERROR
	INFO
)

type Logger struct {
	l     *log.Logger
	level logLevel
}

func (l *Logger) Out() {
	if l.level == DEBUG {
		spew.Dump(l)
	}
}

func (l *Logger) SetLogLevel(level logLevel) {
	l.level = level
}

func (l *Logger) SetOutputFile(path string) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		l.l.Fatal(err)
	}
	l.l.SetOutput(file)
}

func (l *Logger) SetFlags(flags int) {
	l.l.SetFlags(flags)
}

func (l *Logger) SetLogger(logger *log.Logger) {
	l.l = logger
}

func InitLogger(logger *log.Logger, level logLevel, flags int, filePath string) *Logger {
	l := Logger{}
	l.SetLogger(logger)
	l.SetLogLevel(level)
	l.SetFlags(flags)
	if len(filePath) != 0 {
		l.SetOutputFile(filePath)
	}
	return &l
}

func (l *Logger) Debug(args ...interface{}) {
	if l.level <= DEBUG {
		args = append(make([]interface{}, 1, len(args)+1), args)
		args[0] = "LOG "
		l.l.Println(args...)
	}
}

func (l *Logger) Error(args ...interface{}) {
	if l.level <= ERROR {
		args = append(make([]interface{}, 1, len(args)+1), args)
		args[0] = "ERROR "
		l.l.Println(args...)
	}
}

func (l *Logger) Info(args ...interface{}) {
	if l.level <= INFO {
		args = append(make([]interface{}, 1, len(args)+1), args)
		args[0] = "INFO "
		l.l.Println(args...)
	}
}

func (l *Logger) Fatal(args ...interface{}) {
	args = append(make([]interface{}, 1, len(args)+1), args)
	args[0] = "FATAL "
	l.l.Fatal(args...)
}
