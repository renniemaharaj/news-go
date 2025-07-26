package coordinator

import (
	"github.com/renniemaharaj/news-go/internal/document"
	"github.com/renniemaharaj/news-go/internal/loggers"
	"github.com/renniemaharaj/news-go/internal/reporter"
	"github.com/renniemaharaj/news-go/internal/store"
)

type Instance struct {
	r *reporter.Instance
	c chan document.Report

	Store store.Instance
}

// Singleton coordinator instance
var singleton *Instance

func Initialize() {
	singleton = &Instance{}
	l := loggers.LOGGER_COORDINATOR.Info("Initializing")

	// 2) Initialize and hydrate store from disk
	singleton.Store = *store.CreateStore(l)
	singleton.Store.Hydrate()

	// 3) Shake the store in memory and schedule routine audits w/ audit disk
	go singleton.auditStoreRoutine(&singleton.Store)
	singleton.auditStoreRoutineScheduler(&singleton.Store)

	// Create the report and initialize the coordinator's reporter and receive channel
	singleton.r, singleton.c = reporter.CreateReporter(l)

	// 4) Run update once at startup and schedule update routine
	go singleton.updateRoutine(&singleton.Store)
	singleton.updateStoreRoutineScheduler(&singleton.Store)
}

// Get returns the singleton instance
func Get() *Instance {
	if singleton == nil {
		Initialize()
	}
	return singleton
}
