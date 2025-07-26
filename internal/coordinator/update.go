package coordinator

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/renniemaharaj/news-go/internal/config"
	"github.com/renniemaharaj/news-go/internal/document"
	"github.com/renniemaharaj/news-go/internal/loggers"
	"github.com/renniemaharaj/news-go/internal/store"
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
	loggers.LOGGER_COORDINATOR.Info("Scheduled periodic update routine")
}

// Coordinator's, store update routine
func (i *Instance) updateRoutine(storeInstance *store.Instance) {
	loggers.LOGGER_COORDINATOR.Info("Running update routine")

	// 1) Get current report titles
	upToDateTitles := storeInstance.GetUpdatedReports(loggers.LOGGER_COORDINATOR)
	upToDateMap := store.StringSliceToEmptyStructMap(upToDateTitles)

	// 2) Get all configured search queries
	searchQueries := config.Get().SearchQueries

	i.Store.HydrateTags()

	if err := godotenv.Load(); err != nil {
		loggers.LOGGER_COORDINATOR.Error("No .env file found, or error loading it: " + err.Error())
		return
	}

	if os.Getenv("ENABLE_BROWSER") == "" {
		loggers.LOGGER_COORDINATOR.Warning("Browser functionality is disabled. Set ENABLE_BROWSER to enable")
		return
	}

	// 3) For each configured search query, send only if not up-to-date
	for _, query := range searchQueries {
		if _, exists := upToDateMap[store.SanitizeFilename(query)]; !exists {
			i.r.TODO_SEARCH_CHANNEL <- document.ReportFromQuery(query)
			loggers.LOGGER_COORDINATOR.Info(fmt.Sprintf("Will updated: %s", query))
		}
	}

	// 4) Read from reporter channel and persist
	go func() {
		for report := range i.c {
			loggers.LOGGER_COORDINATOR.Success(fmt.Sprintf("Report completed: %s", report.SearchQuery))
			storeInstance.StoreReport(&report, loggers.LOGGER_COORDINATOR)
			i.Store.HydrateTags()
		}
	}()
}
