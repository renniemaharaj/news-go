package log

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Logger struct {
	mu          sync.Mutex
	file        *os.File
	writer      *os.File
	currentLine int
	MaxLines    int
	ToStdout    bool
	JSONMode    bool
}

var prefixes = map[string]string{
	"debug":   "🐞",
	"info":    "ℹ️",
	"warning": "⚠️",
	"error":   "❌",
	"success": "✅",
}

var twclogsDir = "twclogs"

func NewLogger(maxLines int, toStdout bool, jsonMode bool) *Logger {
	if err := os.MkdirAll(twclogsDir, 0755); err != nil {
		panic("Failed to create log directory: " + err.Error())
	}

	l := &Logger{
		MaxLines: maxLines,
		ToStdout: toStdout,
		JSONMode: jsonMode,
	}
	l.rotate()

	return l
}

func (l *Logger) GetLogFileName() string {
	return l.file.Name()
}

func (l *Logger) rotate() {
	if l.file != nil {
		l.file.Close()
	}

	filename := filepath.Join(twclogsDir, "log-"+time.Now().Format("2006-01-02-15-04-05")+".log")

	f, err := os.Create(filename)
	if err != nil {
		panic("Failed to create log file: " + err.Error())
	}

	l.file = f
	l.writer = f
	l.currentLine = 0
}

func (l *Logger) log(level string, msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	prefix := prefixes[level]

	line := fmt.Sprintf("%s [%s] %s %s", timestamp, strings.ToUpper(level), prefix, msg)

	if l.JSONMode {
		logObj := map[string]interface{}{
			"time":  timestamp,
			"level": level,
			"msg":   msg,
		}
		jsonBytes, _ := json.Marshal(logObj)
		line = string(jsonBytes)
	}

	// Write to log file
	fmt.Fprintln(l.writer, line)
	l.currentLine++

	// Also print to stdout
	if l.ToStdout {
		fmt.Println(line)
	}

	// Rotate if line limit is reached
	if l.currentLine >= l.MaxLines {
		l.rotate()
	}
}

func (l *Logger) Info(msg string) {
	l.log("info", msg)
}

func (l *Logger) Debug(msg string) {
	l.log("debug", msg)
}

func (l *Logger) Success(msg string) {
	l.log("success", msg)
}

func (l *Logger) Warning(msg string) {
	l.log("warning", msg)
}

func (l *Logger) Error(msg string) {
	l.log("error", msg)
}

func (l *Logger) Fatal(e error) {
	l.log("error", e.Error())
}
