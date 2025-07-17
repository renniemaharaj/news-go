package store

import (
	"strings"

	"github.com/renniemaharaj/news-go/internal/document"
)

func (s *Instance) AllReports() []*document.Report {
	reports := make([]*document.Report, 0, len(s.reportsByTitle))
	for _, r := range s.reportsByTitle {
		reports = append(reports, r)
	}
	return reports
}

func (s *Instance) ReportByTitle(title string) (*document.Report, bool) {
	r, ok := s.reportsByTitle[strings.ToLower(title)]
	return r, ok
}

func (s *Instance) ResultByTitles(reportTitle, resultTitle string) (*document.Result, bool) {
	report, ok := s.ReportByTitle(reportTitle)
	if !ok {
		return nil, false
	}
	for _, res := range report.Results {
		if strings.EqualFold(res.Title, resultTitle) {
			return &res, true
		}
	}
	return nil, false
}
