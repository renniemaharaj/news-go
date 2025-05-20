package coordinator

import (
	"fmt"
	"time"

	"github.com/renniemaharaj/news/internal/log"
	"github.com/renniemaharaj/news/internal/types"
)

func JobsRoutine(l *log.Logger) {

	for {
		select {
		case r := <-PreResultsChannel:

			collectResults(&r, l)
			PreContentChannel <- r
		case <-time.After(10 * time.Second):
			l.Debug("JobsRoutine... waiting for PreResultsChannel input")
		}
	}
}

func collectResults(r *types.Report, l *log.Logger) {
	err := r.CollectResults(l)
	if err != nil {
		l.Error(fmt.Sprintf("Error collecting result %s", err.Error()))
	}
}
