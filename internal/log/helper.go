package log

import (
	"os"
	"path/filepath"
	"time"
)

func CreateLogger(prefix string, maxLines int, toStdout bool, jsonMode bool, debugging bool) *Logger {
	if err := os.MkdirAll(twclogsDir, 0755); err != nil {
		panic("Failed to create log directory: " + err.Error())
	}

	l := &Logger{
		Prefix:    prefix,
		MaxLines:  maxLines,
		ToStdout:  toStdout,
		JSONMode:  jsonMode,
		Debugging: debugging,
	}
	l.rotate()

	return l
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

func (l *Logger) Info(msg string) {
	l.log("info", msg)
}

func (l *Logger) Debug(msg string) {
	if l.Debugging {
		l.log("debug", msg)
	}

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
