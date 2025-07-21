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
	file        *os.File
	writer      *os.File
	currentLine int
	MaxLines    int
	ToStdout    bool
	JSONMode    bool
	Prefix      string
	Debugging   bool
}

var twclogsDir = "./twclogs"

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

	// Write to log file and stdout if enabled
	fmt.Fprintln(l.writer, lineString)
	if l.ToStdout {
		switch line.Level {
		case "INFO":
			fmt.Println(lineString)
		case "DEBUG":
			if l.Debugging {
				fmt.Println(lineString)
			}
		case "SUCCESS":
			color.Green(lineString)
		case "WARNING":
			color.Yellow(lineString)
		case "ERROR":
			color.Red(lineString)
		case "FATAL":
			color.Red(lineString)
			os.Exit(1)
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
