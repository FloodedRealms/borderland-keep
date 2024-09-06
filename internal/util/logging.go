package util

import (
	"fmt"
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

func (m *Logger) PrintWithSessionId(sessionId string, args ...interface{}) {
	s := fmt.Sprintf("Session %s", sessionId)
	log.Print(s)
	log.Print(args...)
}
