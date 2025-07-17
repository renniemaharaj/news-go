package config

import "github.com/renniemaharaj/news-go/internal/log"

func createLogger() *log.Logger {
	return log.CreateLogger("Config", 100, true, false, false)
}
