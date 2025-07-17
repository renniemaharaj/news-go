package routines

import (
	"time"

	"github.com/renniemaharaj/news-go/internal/log"
	// "github.com/renniemaharaj/news-go/internal/storage"
	"github.com/renniemaharaj/news-go/internal/document"
)

// Coordinators persist routine, and file routine in pipeline
func Persist(l *log.Logger, PARENT chan document.Report, CHILD chan document.Report) {

	for {
		select {
		case r := <-PARENT:
			l.Info("PERSIST_WORKER: Received completed job")
			persist(&r, l)
			CHILD <- r
			l.Info("PERSIST_WORKER: Persisted completed job")
		case <-time.After(10 * time.Second):
			l.Debug("PERSIST_WORKER: waiting for jobs")
		}
	}
}

// Helper function
func persist(r *document.Report, l *log.Logger) {
	if len(r.Results) > 0 && r.Results[0].Title != "" {
		r.Title = r.Results[0].Title
		r.Date = time.Now().Format(time.RFC3339)
	}
}
