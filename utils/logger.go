package utils

import (
	"fmt"
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
	if len(filePath) != 0 {
		l.SetOutputFile(filePath)
	}
	l.l.Println()
	l.SetFlags(flags)
	return &l
}

func (l *Logger) Debug(format string, args ...interface{}) {
	if l.level <= DEBUG {
		l.l.Println("DEBUG: ", fmt.Sprintf(format, args...))
	}
}

func (l *Logger) Error(format string, args ...interface{}) {
	if l.level <= ERROR {
		l.l.Println("ERROR: ", fmt.Sprintf(format, args...))
	}
}

func (l *Logger) Info(format string, args ...interface{}) {
	if l.level <= INFO {
		l.l.Println("INFO: ", fmt.Sprintf(format, args...))
	}
}

func (l *Logger) Fatal(format string, args ...interface{}) {
	l.l.Fatal("FATAL: ", fmt.Sprintf(format, args...))
}
