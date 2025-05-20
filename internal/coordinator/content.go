package coordinator

import (
	"time"

	"github.com/renniemaharaj/news/internal/log"
	"github.com/renniemaharaj/news/internal/types"
)

func ContentRoutine(l *log.Logger) {

	for {
		select {
		case r := <-PreContentChannel:
			scrapeContent(&r, l)
			PreModelChannel <- r
		case <-time.After(10 * time.Second):
			l.Debug("ScrapeRoutine... waiting for PreContentChannel input")
		}
	}
}

func scrapeContent(r *types.Report, l *log.Logger) {
	for _, result := range r.Results {
		result.RequestContent(l)
	}
}
