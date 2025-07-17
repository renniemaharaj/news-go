package storage

import (
	"github.com/renniemaharaj/news-go/internal/document"
	"github.com/renniemaharaj/news-go/internal/log"
	"github.com/renniemaharaj/news-go/internal/storage/store"
)

// Top level, get all reports command
func GetAllReports(l *log.Logger) []*document.Report {
	store := store.CreateStore(l)
	store.Hydrate()

	return store.AllReports()
}
