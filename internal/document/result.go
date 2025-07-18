package document

import (
	"fmt"

	"github.com/renniemaharaj/news-go/internal/browser"
	"github.com/renniemaharaj/news-go/internal/log"
)

// A wrapper for each search result and it's information
type Result struct {
	Title           string   `json:"title"`           // !Model's title based on text content
	Commentaries    []string `json:"commentaries"`    // !Editor's commentaries
	TextContent     string   `json:"-"`               // !System's Accumulated text content
	Alignment       int      `json:"alignment"`       // !Model's framework alignment score (0–10)
	Summary         string   `json:"summary"`         // !Model's thoughtful summary with framework reflection
	Images          []string `json:"images"`          // !System's relevant image URLs (e.g., thumbnails)
	HREF            string   `json:"href"`            // !System's assigned link
	Tags            []string `json:"tags"`            // !Model's categories/topics, including framework themes
	PoliticalBiases []string `json:"politicalBiases"` // !Model's examination and tagging for political biases eg: leftist, right, progressive, conservative
}

// Function requests and returns http response to reduce required requests
func (r *Result) RequestContent(l *log.Logger) error {
	textContent, images, err := browser.Get().Content(r.HREF)
	if err != nil {
		l.Error(fmt.Sprintf("Failed to read site: %s", r.HREF))
		return err
	}

	r.TextContent = textContent
	r.Images = images
	return nil
}
