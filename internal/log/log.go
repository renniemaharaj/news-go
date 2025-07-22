package log

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

type Line struct {
	Prefix string `json:"prefix"`
	Time   string `json:"time"`
	Level  string `json:"level"`
	Msg    string `json:"msg"`
}

type Logger struct {
	mu          sync.Mutex
	writer      *os.File
	currentLine int
	MaxLines    int
	ToStdout    bool
	JSONMode    bool
	Prefix      string
	Debugging   bool

	Subscribers *Subscribers
}

var GlobalLogger *Logger

var twcLogsDir = "./twcLogs"

func InitGlobalLogger() {
	GlobalLogger = &Logger{
		Prefix:    "GLOBAL",
		MaxLines:  100,
		ToStdout:  true,
		JSONMode:  false,
		Debugging: false,

		Subscribers: &Subscribers{},
	}

	GlobalLogger.rotate()
}

func (l *Logger) log(level string, msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	timestamp := time.Now().Format("2006-01-02 15:04:05.000")

	line := &Line{
		Time:   timestamp,
		Level:  strings.ToUpper(level),
		Prefix: l.Prefix,
		Msg:    msg,
	}

	lineString := fmt.Sprintf("%s: %s [%s] %s", timestamp, line.Level, l.Prefix, msg)

	if l.JSONMode {
		jsonBytes, _ := json.Marshal(line)
		lineString = string(jsonBytes)
	}

	debugFunc := func() {
		if l.Debugging {
			fmt.Println(lineString)
		}
	}

	// Broadcast to all local subscribers
	l.Subscribers.Broadcast(*line)

	// If GlobalLogger, broadcast to all global logger's subscribers
	if GlobalLogger != nil {
		GlobalLogger.Subscribers.Broadcast(*line)
	}

	// Write to log file and stdout if enabled
	fmt.Fprintln(l.writer, lineString)
	if l.ToStdout {
		switch line.Level {
		case "INFO":
			fmt.Println(lineString)
		case "SUCCESS":
			color.Green(lineString)
		case "WARNING":
			color.Yellow(lineString)
		case "ERROR":
			color.Red(lineString)
		case "FATAL":
			color.Red(lineString)
			os.Exit(1)
		case "DEBUG":
			debugFunc()
		default:
			fmt.Println(lineString)
		}
	}
	l.currentLine++

	// Rotate if line limit is reached
	if l.currentLine >= l.MaxLines {
		l.rotate()
	}
}
