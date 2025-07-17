package coordinator

import (
	"time"

	"github.com/renniemaharaj/news-go/internal/config"
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

	auditDiskTicker := time.NewTicker(store.AuditDiskInterval)
	go func() {
		for range auditDiskTicker.C {
			i.auditStoreRoutine(storeInstance)
		}
	}()
	i.l.Info("Scheduled periodic audit routine")
}

func (i *Instance) auditStoreRoutine(storeInstance *store.Instance) {
	// 1) Read config and get an updated, sanitized map of it's search queries
	searchQueries := config.Get().SearchQueries
	searchQueriesMap := store.StringSliceToSanitizedEmptyStructMap(searchQueries)

	// 2) Shake the store, removing irrelevant reports
	i.l.Info("Shaking store")
	storeInstance.ShakeStore(searchQueriesMap)
}
