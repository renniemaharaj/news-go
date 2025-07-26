package routines

import (
	"time"

	"github.com/renniemaharaj/grouplogs/pkg/logger"
	"github.com/renniemaharaj/news-go/internal/document"
)

// Coordinators Content routine
func Content(l *logger.Logger, PARENT chan document.Report, CHILD chan document.Report) {

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
func content(r *document.Report, l *logger.Logger) {
	for i := range r.Results {
		r.Results[i].Hydrate(l)
	}
}
