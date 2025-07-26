package reporter

import (
	"github.com/renniemaharaj/grouplogs/pkg/logger"
	"github.com/renniemaharaj/news-go/internal/document"
)

// Channel creator for the reporter
func CreateChannel() chan document.Report {
	return make(chan document.Report, 100)
}

// Report creator helper, creates, initializes and returns an exposed channel
func CreateReporter(l *logger.Logger) (*Instance, chan document.Report) {
	r := Instance{}
	c := CreateChannel()

	r.InitializeChannels(100, 100, 100, 100)

	r.ReadyRoutines(l, c)

	return &r, c
}
