package coordinator

import (
	"time"

	"github.com/renniemaharaj/news-go/internal/log"
	"github.com/renniemaharaj/news-go/internal/types"
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
	for i := range r.Results {
		r.Results[i].RequestContent(l)
	}
}
