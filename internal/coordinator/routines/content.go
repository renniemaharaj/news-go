package routines

import (
	"time"

	"github.com/renniemaharaj/news-go/internal/document"
	"github.com/renniemaharaj/news-go/internal/log"
)

// Coordinators Content routine
func Content(l *log.Logger, PARENT chan document.Report, CHILD chan document.Report) {

	for {
		select {
		case r := <-PARENT:
			content(&r, l)
			CHILD <- r
		case <-time.After(10 * time.Second):
			l.Debug("#2 CONTENT: Waiting for jobs")
		}
	}
}

// Helper function
func content(r *document.Report, l *log.Logger) {
	for i := range r.Results {
		r.Results[i].RequestContent(l)
	}
}
