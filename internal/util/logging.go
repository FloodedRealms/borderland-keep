package util

import (
	"log"
)

type Logger struct {
	PrintDebug bool
}

func NewLogger(debug bool) *Logger {
	return &Logger{debug}
}

func (m *Logger) Debug(args ...interface{}) {
	if m.PrintDebug {
		m.Print(args...)
	}
}

func (m *Logger) Print(args ...interface{}) {
	log.Print(args...)
}
