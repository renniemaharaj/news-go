package store

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/renniemaharaj/news-go/internal/document"
)

// Efficient shaking method that shakes only the store instance
func (s *Instance) ShakeStore(expected map[string]struct{}) {
	s.ForEach(func(key string, r *document.Report) {
		if _, ok := expected[key]; !ok {
			delete(s.reportsByTitle, key)
			s.l.Info(fmt.Sprintf("Shook report from memory: %s", key))
		}
	})
}

// Shaking method that shakes the actual disk
func (s *Instance) ShakeDisk(expected map[string]struct{}) error {
	err := filepath.WalkDir(reportsDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			s.l.Error(fmt.Sprintf("Error accessing path %s: %s", path, err.Error()))
			return nil // Continue walking
		}

		if FileIsObject(d) {
			name := strings.TrimSuffix(d.Name(), ".json")
			if _, ok := expected[name]; !ok {
				if err := os.Remove(path); err != nil {
					s.l.Error(fmt.Sprintf("Failed to remove file: %s", path))
				} else {
					s.l.Info(fmt.Sprintf("Shook report from disk: %s", path))
				}
			}
		}

		return nil
	})

	return err
}
