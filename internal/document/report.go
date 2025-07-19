package document

import (
	"fmt"

	"github.com/renniemaharaj/news-go/internal/browser"
	"github.com/renniemaharaj/news-go/internal/log"
	"github.com/renniemaharaj/news-go/internal/utils"
)

// A single report structure
type Report struct {
	SearchQuery string   `json:"searchQuery"` // !System's Search query being used to get news site links
	Results     []Result `json:"results"`     // !System's search results list
	Title       string   `json:"title"`       // !Model's clear and concise report title
	Date        string   `json:"date"`        // !System's ISO 8601 format (e.g., 2025-05-03)
}

// Reports forEach method
func (r *Report) ForEach(callback func(int, *Result)) {
	for k, r := range r.Results {
		callback(k, &r)
	}
}

// Reports filter method, returns filtered slice
func (r *Report) Filter(predicate func(int, *Result) bool) []*Result {
	var filtered []*Result
	for k, r := range r.Results {
		if predicate(k, &r) {
			filtered = append(filtered, &r)
		}
	}
	return filtered
}

func (r *Report) HasTagIntersection(userTags []string) bool {
	// Convert user tags to set
	userTagSet := utils.StringSliceToMap(userTags)

	// Walk report results and check tags
	for _, result := range r.Results {
		for _, tag := range result.Tags {
			if _, ok := userTagSet[tag]; ok {
				return true
			}
		}
	}
	return false
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
