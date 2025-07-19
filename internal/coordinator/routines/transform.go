package routines

import (
	"time"

	"github.com/renniemaharaj/news-go/internal/document"
	"github.com/renniemaharaj/news-go/internal/log"
	"github.com/renniemaharaj/news-go/internal/model"
)

// Coordinators transform routine
func Transform(l *log.Logger, PARENT chan document.Report, CHILD chan document.Report) {

	for {
		select {
		case r := <-PARENT:
			transform(&r, l)
			CHILD <- r
			time.Sleep(3 * time.Second)
		case <-time.After(10 * time.Second):
			l.Debug("#3 TRANSFORM: waiting for jobs")
		}
	}
}

// Helper function
func transform(r *document.Report, l *log.Logger) {
	var results []document.Result
	for i := range r.Results {
		result, err := model.Transform(&r.Results[i], l)

		if err != nil {
			l.Error(err.Error())
			continue
		}
		results = append(results, result)
	}

	r.Results = results
}
