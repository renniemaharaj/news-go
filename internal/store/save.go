package store

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/renniemaharaj/grouplogs/pkg/logger"
	"github.com/renniemaharaj/news-go/internal/document"
)

func (s *Instance) StoreReport(r *document.Report, l *logger.Logger) {
	// Create reports directory if it doesn't exist
	if err := os.MkdirAll(reportsDir, 0755); err != nil {
		l.Error(fmt.Sprintf(("Failed to create reports directory: %s"), err.Error()))
		return
	}

	// Marshal report to JSON bytes
	jsonBytes, err := json.Marshal(r)
	if err != nil {
		l.Error(fmt.Sprintf(("Failed to marshal report: %s"), err.Error()))
		return
	}

	sanitizedTitle := SanitizeFilename(r.SearchQuery)

	// Create file name based on searchQuery
	fileName := fmt.Sprintf("%s/%s.json", reportsDir, sanitizedTitle)

	// Write to file
	if err := os.WriteFile(fileName, jsonBytes, 0644); err != nil {
		l.Error(fmt.Sprintf(("Failed to write report: %s"), err.Error()))
		l.Error(fileName)
		return
	}

	s.reportsByTitle[sanitizedTitle] = r
	l.Success(fmt.Sprintf(("Saved report: %s"), fileName))
}
