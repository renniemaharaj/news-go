package browser

import "github.com/renniemaharaj/news-go/internal/log"

func createLogger() *log.Logger {
	return log.CreateLogger("Browser", 100, true, false, false)
}
