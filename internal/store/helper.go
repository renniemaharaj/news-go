package store

import (
	"time"

	"github.com/renniemaharaj/news-go/internal/document"
	"github.com/renniemaharaj/news-go/internal/log"
)

// Reports forEach method
func (s *Instance) ForEach(callback func(string, *document.Report)) {
	for k, r := range s.reportsByTitle {
		callback(k, r)
	}
}

// Reports filter method, returns filtered slice
func (s *Instance) Filter(predicate func(string, *document.Report) bool) []*document.Report {
	var filtered []*document.Report
	for k, r := range s.reportsByTitle {
		if predicate(k, r) {
			filtered = append(filtered, r)
		}
	}
	return filtered
}

// Returns a []string of updated reports
func (i *Instance) GetUpdatedReports(l *log.Logger) []string {
	titles := make([]string, 0)

	i.ForEach(func(k string, r *document.Report) {
		reportTime, err := time.Parse(reportTimeLayout, r.Date)
		if err != nil {
			l.Error("Forced to assume report as outdated: time parsing error")
		}

		if time.Since(reportTime) < UpdateInterval {
			titles = append(titles, k)
		}
	})

	return titles
}
