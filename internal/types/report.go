package types

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/renniemaharaj/news-go/internal/browser"
	"github.com/renniemaharaj/news-go/internal/log"
)

// The model should not populate !System's
// The model should not populate !Editor's

// The model should populate !Model's

// A wrapper for each search result and it's information
type Result struct {
	Title           string   `json:"title"`           // !Model's title based on text content
	Commentaries    []string `json:"commentaries"`    // !Editor's commentaries
	TextContent     string   `json:"textContent"`     // !System's Accumulated text content
	Alignment       int      `json:"alignment"`       // !Model's framework alignment score (0–10)
	Summary         string   `json:"summary"`         // !Model's thoughtful summary with framework reflection
	Images          []string `json:"images"`          // !System's relevant image URLs (e.g., thumbnails)
	HREF            string   `json:"href"`            // !System's assigned link
	Tags            []string `json:"tags"`            // !Model's categories/topics, including framework themes
	PoliticalBiases []string `json:"politicalBiases"` // !Model's examination and tagging for political biases eg: leftist, right, progressive, conservative
}

// Function requests and returns http response to reduce required requests
func (r *Result) RequestContent(l *log.Logger) error {
	resp, err := browser.Request(r.HREF, l)
	if err != nil {
		l.Error(fmt.Sprintf("Error scraping site body from: %s", r.HREF))
		return err
	}

	// Create two copies: one for text, one for images
	copies, err := browser.DeepCopyBody(resp.Body, 2)
	if err != nil {
		l.Error(fmt.Sprintf("Error duplicating response body for: %s", r.HREF))
		return err
	}

	r.PopulateTextContent(copies[0], l)
	r.PopulateImages(copies[1], resp.Request.URL.String(), l)

	return nil
}

// Function to population report result text content
func (r *Result) PopulateTextContent(body io.Reader, l *log.Logger) {
	if closer, ok := body.(io.Closer); ok {
		defer closer.Close()
	}

	textContent, err := browser.ExtractTextContent(body)
	if err != nil {
		l.Error(fmt.Sprintf("Error scraping for text content from: %s", r.HREF))
		return
	}

	r.TextContent = textContent

}

// Function to populate report result images
func (r *Result) PopulateImages(body io.Reader, base string, l *log.Logger) {
	if closer, ok := body.(io.Closer); ok {
		defer closer.Close()
	}

	images, err := browser.ExtractThumbnails(body, base)
	if err != nil {
		l.Error(fmt.Sprintf("Error scraping for images from: %s", r.HREF))
		return
	}

	r.Images = images
}

// A single report structure
type Report struct {
	SearchQuery string   `json:"searchQuery"` // !System's Search query being used to get news site links
	Results     []Result `json:"results"`     // !System's search results list
	Title       string   `json:"title"`       // !Model's clear and concise report title
	Date        string   `json:"date"`        // !System's ISO 8601 format (e.g., 2025-05-03)
}

func (r *Report) FromSearchQuery(s string) Report {
	return Report{
		SearchQuery: s,
	}
}

func (r *Report) CollectResults(l *log.Logger) error {
	results, err := browser.Search(r.SearchQuery, 2, l)
	if err != nil {
		l.Error(fmt.Sprintf("Error getting google results: %s", r.SearchQuery))
		return err
	}

	resultBytes, err := json.Marshal(results)
	if err != nil {
		l.Error("Error marshalling results")
	}
	l.Info(fmt.Sprintf("Google results: %s", string(resultBytes)))

	for _, result := range results {
		l.Info(fmt.Sprintf("Got %s from collection", result))
		r.Results = append(r.Results, Result{
			HREF: result,
		})
	}

	return nil
}
