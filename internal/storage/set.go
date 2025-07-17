package storage

import (
	"github.com/renniemaharaj/news-go/internal/document"
	"github.com/renniemaharaj/news-go/internal/log"
	"github.com/renniemaharaj/news-go/internal/storage/store"
)

// Top level, store report command
func StoreReport(r *document.Report, l *log.Logger) {
	store := store.CreateStore(l)
	store.StoreReport(r, l)
}
