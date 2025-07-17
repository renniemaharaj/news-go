package store

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/renniemaharaj/news-go/internal/document"
	"github.com/renniemaharaj/news-go/internal/log"
)

const reportsDir = "./reports"

var UpdateInterval = 24 * time.Hour
var AuditStoreInterval = time.Hour
var AuditDiskInterval = 24 * 30 * time.Hour

const reportTimeLayout = time.RFC3339

type Instance struct {
	reportsByTitle map[string]*document.Report
	l              *log.Logger
}

func CreateStore(l *log.Logger) *Instance {
	return &Instance{
		reportsByTitle: make(map[string]*document.Report),
		l:              l,
	}
}

// Hydrate the store map of reports from the disk
func (s *Instance) Hydrate() error {
	err := filepath.WalkDir(reportsDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			s.l.Error(fmt.Sprintf("Error accessing path %s: %s", path, err.Error()))
			return nil
		}

		if FileIsObject(d) {
			data, err := os.ReadFile(path)
			if err != nil {
				s.l.Error(fmt.Sprintf("Failed to read file: %s", path))
				return nil
			}

			if r, err := BytesToReport(data); err != nil {
				s.l.Error(fmt.Sprintf("Failed to unmarshal file: %s", path))
			} else {
				key := strings.TrimSuffix(d.Name(), ".json")
				s.reportsByTitle[key] = r
			}

		}

		return nil
	})

	if err != nil {
		s.l.Error(fmt.Sprintf("Failed to walk reports directory: %s", err.Error()))
		return err
	}

	s.l.Success(fmt.Sprintf("Loaded %d report(s)", len(s.reportsByTitle)))
	return nil
}
