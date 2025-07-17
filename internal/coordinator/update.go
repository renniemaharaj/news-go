package coordinator

import (
	"fmt"
	"time"

	"github.com/renniemaharaj/news-go/internal/config"
	"github.com/renniemaharaj/news-go/internal/document"
	"github.com/renniemaharaj/news-go/internal/storage/store"
)

// Update routine scheduler
func (i *Instance) updateStoreRoutineScheduler(storeInstance *store.Instance) {
	// Now set up ticker to trigger daily updates
	ticker := time.NewTicker(store.UpdateInterval)
	go func() {
		for range ticker.C {
			i.updateRoutine(storeInstance)
		}
	}()
	i.l.Info("Scheduled periodic update routine")
}

// Coordinator's, store update routine
func (i *Instance) updateRoutine(storeInstance *store.Instance) {
	i.l.Info("Running update routine")

	// 1) Get current report titles
	upToDateTitles := storeInstance.GetUpdatedReports(i.l)
	upToDateMap := store.StringSliceToEmptyStructMap(upToDateTitles)

	// 2) Get all configured search queries
	searchQueries := config.Get().SearchQueries

	// 3) For each configured search query, send only if not up-to-date
	for _, query := range searchQueries {
		if _, exists := upToDateMap[store.SanitizeFilename(query)]; !exists {
			i.r.TODO_SEARCH_CHANNEL <- document.ReportFromQuery(query)
			i.l.Info(fmt.Sprintf("Queueing search for outdated query: %s", query))
		}
	}

	// 4) Read from reporter channel and persist
	go func() {
		for report := range i.c {
			i.l.Info(fmt.Sprintf("Received updated report: %s", report.SearchQuery))
			storeInstance.StoreReport(&report, i.l)
		}
	}()
}
