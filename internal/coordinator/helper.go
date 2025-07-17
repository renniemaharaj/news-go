package coordinator

import "github.com/renniemaharaj/news-go/internal/log"

func createLogger() *log.Logger {
	return log.CreateLogger("Coordinator", 100, true, false, false)
}
