package routines

import (
	"time"

	"github.com/renniemaharaj/news-go/internal/document"
	"github.com/renniemaharaj/news-go/internal/log"
)

// Coordinators persist routine, and file routine in pipeline
func Store(l *log.Logger, PARENT chan document.Report, CHILD chan document.Report) {

	for {
		select {
		case r := <-PARENT:
			store(&r, l)
			CHILD <- r
		case <-time.After(10 * time.Second):
			l.Debug("#4 STORE: waiting for jobs")
		}
	}
}

// Helper function
func store(r *document.Report, l *log.Logger) {
	if len(r.Results) > 0 && r.Results[0].Title != "" {
		r.Title = r.Results[0].Title
		r.Date = time.Now().Format(time.RFC3339)
	}
}
