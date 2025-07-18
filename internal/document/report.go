package document

import (
	"fmt"

	"github.com/renniemaharaj/news-go/internal/browser"
	"github.com/renniemaharaj/news-go/internal/log"
)

// A single report structure
type Report struct {
	SearchQuery string   `json:"searchQuery"` // !System's Search query being used to get news site links
	Results     []Result `json:"results"`     // !System's search results list
	Title       string   `json:"title"`       // !Model's clear and concise report title
	Date        string   `json:"date"`        // !System's ISO 8601 format (e.g., 2025-05-03)
}

func (r *Report) CollectResults(l *log.Logger, sitesPerQuery int) error {
	results, err := browser.Get().Search(r.SearchQuery, sitesPerQuery)
	if err != nil {
		l.Error(fmt.Sprintf("Error getting google results: %s", r.SearchQuery))
		return err
	}

	for _, result := range results {
		r.Results = append(r.Results, Result{
			HREF: result,
		})
	}

	return nil
}
