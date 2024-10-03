package utils

import (
	"encoding/json"
	"io"
	"runtime/debug"
	"sync"
	"time"
)

type Level int8

const (
	LevelInfo Level = iota
	LevelError
	LevelFatal
)

func (l Level) String() string {
	switch l {
	case LevelInfo:
		return "INFO"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return ""
	}
}

type LogData struct {
	Level      string            `json:"level"`
	Time       string            `json:"time"`
	Message    string            `json:"message"`
	Properties map[string]string `json:"properties,omitempty"`
	Trace      string            `json:"trace,omitempty"`
}

type Logger struct {
	out   io.Writer
	mutex sync.Mutex
}

func NewLogger(out io.Writer) *Logger {
	return &Logger{
		out: out,
	}
}

func (l *Logger) Log(level Level, message string, properties map[string]string) (int, error) {
	logData := LogData{
		Level:      level.String(),
		Time:       time.Now().Format(time.RFC3339),
		Message:    message,
		Properties: properties,
	}

	if level >= LevelError {
		logData.Trace = string(debug.Stack())
	}

	var logEntry []byte
	logEntry, err := json.Marshal(&logData)

	if err != nil {
		logEntry = []byte(LevelError.String() + ": unable to log message: " + err.Error())
	}

	l.mutex.Lock()
	defer l.mutex.Unlock()
	return l.out.Write(append(logEntry, '\n'))
}
