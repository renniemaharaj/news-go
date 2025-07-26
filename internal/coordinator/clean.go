package coordinator

import (
	"time"

	"github.com/renniemaharaj/news-go/internal/config"
	"github.com/renniemaharaj/news-go/internal/loggers"
	"github.com/renniemaharaj/news-go/internal/store"
)

// Coordinator's store, scheduled audit routine
func (i *Instance) auditStoreRoutineScheduler(storeInstance *store.Instance) {
	// Now set up ticker to trigger daily updates
	auditStoreTicker := time.NewTicker(store.AuditStoreInterval)
	go func() {
		for range auditStoreTicker.C {
			i.auditStoreRoutine(storeInstance)
		}
	}()
	loggers.LOGGER_COORDINATOR.Info("Scheduled periodic audit routine")
}

func (i *Instance) auditStoreRoutine(storeInstance *store.Instance) {
	// 1) Read config and get an updated, sanitized map of it's search queries
	searchQueries := config.Get().SearchQueries
	searchQueriesMap := store.StringSliceToSanitizedEmptyStructMap(searchQueries)

	// 2) Shake the store, removing irrelevant reports
	loggers.LOGGER_COORDINATOR.Info("Shaking store")
	storeInstance.ShakeStore(searchQueriesMap)
}
