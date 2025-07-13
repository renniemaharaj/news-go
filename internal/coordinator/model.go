package coordinator

import (
	"time"

	"github.com/renniemaharaj/news-go/internal/log"
	"github.com/renniemaharaj/news-go/internal/model"
	"github.com/renniemaharaj/news-go/internal/types"
)

func ModelsRoutine(l *log.Logger) {

	for {
		select {
		case r := <-PreModelChannel:
			transform(&r, l)
			JobsCompleteChannel <- r
			time.Sleep(3 * time.Second) // Hard coded delay for additional security
		case <-time.After(10 * time.Second):
			l.Debug("ModelsRoutine... waiting for PreModelChannel input")
		}
	}
}

func transform(r *types.Report, l *log.Logger) {
	var results []types.Result
	for i := range r.Results {
		result, err := model.Transform(&r.Results[i], l)

		if err != nil {
			l.Error("Error transforming result")
		}
		results = append(results, result)
	}

	r.Results = results
}
