package model

import "github.com/renniemaharaj/news-go/internal/log"

func createLogger() *log.Logger {
	return log.CreateLogger("Transformer", 100, true, false, false)
}
