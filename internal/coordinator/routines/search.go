package routines

import (
	"fmt"
	"time"

	"github.com/renniemaharaj/grouplogs/pkg/logger"
	"github.com/renniemaharaj/news-go/internal/config"
	"github.com/renniemaharaj/news-go/internal/document"
)

// Coordinator's search routine
func Search(l *logger.Logger, PARENT chan document.Report, CHILD chan document.Report) {

	for {
		select {
		case r := <-PARENT:
			search(&r, l)
			CHILD <- r
		case <-time.After(10 * time.Second):
			l.Debug("#1 SEARCH: Waiting for jobs")
		}
	}
}

// Helper function
func search(r *document.Report, l *logger.Logger) {
	err := r.FindArticles(l, config.Get().NumSitesPerQuery)
	if err != nil {
		l.Error(fmt.Sprintf("Error collecting result %s", err.Error()))
	}
}
