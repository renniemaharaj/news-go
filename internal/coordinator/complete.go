package coordinator

import (
	"time"

	"github.com/renniemaharaj/news/internal/log"
	"github.com/renniemaharaj/news/internal/types"
)

func CompleteRoutine(l *log.Logger) {

	for {
		select {
		case r := <-JobsCompleteChannel:
			completeJob(&r, l)
		case <-time.After(10 * time.Second):
			l.Debug("CompleteRoutine... waiting for JobsCompleteChannel input")
		}
	}
}

func completeJob(r *types.Report, l *log.Logger) {
	if len(r.Results) > 0 && r.Results[0].Title != "" {
		r.Title = r.Results[0].Title
		r.Date = time.Now().Format(time.RFC3339)

		Save(r, l)
	}
}
