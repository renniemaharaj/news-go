package routines

import (
	"time"

	"github.com/renniemaharaj/news-go/internal/document"
	"github.com/renniemaharaj/news-go/internal/log"
)

// Coordinators Scrape routine
func Scrape(l *log.Logger, PARENT chan document.Report, CHILD chan document.Report) {

	for {
		select {
		case r := <-PARENT:
			l.Info("SCRAPE_WORKER: Using browser")
			scrape(&r, l)
			CHILD <- r
			l.Info("SCRAPE_WORKER: Completed a job")
		case <-time.After(10 * time.Second):
			l.Debug("SCRAPE_WORKER: Waiting for jobs")
		}
	}
}

// Helper function
func scrape(r *document.Report, l *log.Logger) {
	for i := range r.Results {
		r.Results[i].RequestContent(l)
	}
}
