package cloudflare

import "github.com/renniemaharaj/news-go/internal/log"

func createLogger() *log.Logger {
	return log.CreateLogger("Cloudflare", 100, true, false, false)
}
