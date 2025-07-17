package routines

import (
	"fmt"
	"time"

	"github.com/renniemaharaj/news-go/internal/config"
	"github.com/renniemaharaj/news-go/internal/document"
	"github.com/renniemaharaj/news-go/internal/log"
)

// Coordinator's search routine
func Search(l *log.Logger, PARENT chan document.Report, CHILD chan document.Report) {

	for {
		select {
		case r := <-PARENT:
			l.Info("#1 SEARCH: Doing web searches")
			search(&r, l)
			CHILD <- r
			l.Info("#1 SEARCH: Completed a job")
		case <-time.After(10 * time.Second):
			l.Debug("#1 SEARCH: Waiting for jobs")
		}
	}
}

// Helper function
func search(r *document.Report, l *log.Logger) {
	err := r.CollectResults(l, config.Get().NumSitesPerQuery)
	if err != nil {
		l.Error(fmt.Sprintf("Error collecting result %s", err.Error()))
	}
}
