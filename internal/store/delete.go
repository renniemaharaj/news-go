package store

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/renniemaharaj/news-go/internal/document"
	"github.com/renniemaharaj/news-go/internal/loggers"
)

func (s *Instance) DeleteByTitle(title string) error {
	key := SanitizeFilename(title)
	_, exists := s.reportsByTitle[key]
	if !exists {
		return fmt.Errorf("report %q not found", title)
	}

	filePath := filepath.Join(reportsDir, key+".json")
	if err := os.Remove(filePath); err != nil {
		loggers.LOGGER_STORE.Error(fmt.Sprintf("Failed to remove file: %s", filePath))
		return err
	}

	delete(s.reportsByTitle, key)
	loggers.LOGGER_STORE.Info(fmt.Sprintf("Removed report: %s", title))
	return nil
}

func (s *Instance) DeleteResultByTitles(reportTitle, resultTitle string) error {
	reportKey := SanitizeFilename(reportTitle)
	report, exists := s.reportsByTitle[reportKey]
	if !exists {
		return fmt.Errorf("report %q not found", reportTitle)
	}

	// Filter out the result
	newResults := make([]document.Result, 0, len(report.Results))
	found := false
	for _, res := range report.Results {
		if !strings.EqualFold(res.Title, resultTitle) {
			newResults = append(newResults, res)
		} else {
			found = true
		}
	}

	if found {
		report.Results = newResults
		s.StoreReport(report, loggers.LOGGER_STORE)

		loggers.LOGGER_STORE.Info(fmt.Sprintf("Removed result %q from report %q", resultTitle, reportTitle))
		return nil

	}

	return fmt.Errorf("result %q not found in report %q", resultTitle, reportTitle)
}
