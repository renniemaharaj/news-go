package router

import "github.com/renniemaharaj/news-go/internal/log"

func createLogger() *log.Logger {
	return log.CreateLogger("Port Forwarding", 100, true, false, false)
}
