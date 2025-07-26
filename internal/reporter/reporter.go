package reporter

import (
	"github.com/renniemaharaj/grouplogs/pkg/logger"
	"github.com/renniemaharaj/news-go/internal/coordinator/routines"
	"github.com/renniemaharaj/news-go/internal/document"
)

// A reporter instance
type Instance struct {
	// 1. Head of the pipeline: Reports to be searched
	TODO_SEARCH_CHANNEL chan document.Report

	// 2. Reports queued for scraping
	TODO_SCRAPE_CHANNEL chan document.Report

	// 3. Reports queued for transformation
	TODO_TRANSFORM_CHANNEL chan document.Report

	// 4. Finalized reports ready for persistence or broadcast
	TODO_PERSIST_CHANNEL chan document.Report
}

// Initializes reporter channels using make
func (i *Instance) InitializeChannels(b1 int16, b2 int16, b3 int16, b4 int16) {
	i.TODO_SEARCH_CHANNEL = make(chan document.Report, b1)
	i.TODO_SCRAPE_CHANNEL = make(chan document.Report, b2)
	i.TODO_TRANSFORM_CHANNEL = make(chan document.Report, b3)
	i.TODO_PERSIST_CHANNEL = make(chan document.Report, b4)
}

// Essentially, initializes pipeline routines, expects a channel to pass finished jobs
func (i *Instance) ReadyRoutines(l *logger.Logger, CHILD chan document.Report) {

	go routines.Search(l, i.TODO_SEARCH_CHANNEL, i.TODO_SCRAPE_CHANNEL)
	go routines.Content(l, i.TODO_SCRAPE_CHANNEL, i.TODO_TRANSFORM_CHANNEL)
	go routines.Transform(l, i.TODO_TRANSFORM_CHANNEL, i.TODO_PERSIST_CHANNEL)
	go routines.Store(l, i.TODO_PERSIST_CHANNEL, CHILD)
}
