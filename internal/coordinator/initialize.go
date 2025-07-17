package coordinator

import (
	"github.com/renniemaharaj/news-go/internal/document"
	"github.com/renniemaharaj/news-go/internal/log"
	"github.com/renniemaharaj/news-go/internal/reporter"
	"github.com/renniemaharaj/news-go/internal/store"
)

type Instance struct {
	l *log.Logger
	r *reporter.Instance
	c chan document.Report

	Store store.Instance
}

// Singleton coordinator instance
var singleton *Instance

func Initialize() {
	singleton = &Instance{}
	singleton.l = createLogger() // 1) Create the logger for communication
	singleton.l.Info("Initializing")

	// 2) Initialize and hydrate store from disk
	singleton.Store = *store.CreateStore(singleton.l)
	singleton.Store.Hydrate()

	// 3) Shake the store in memory and schedule routine audits w/ audit disk
	go singleton.auditStoreRoutine(&singleton.Store)
	singleton.auditStoreRoutineScheduler(&singleton.Store)

	// Create the report and initialize the coordinator's reporter and receive channel
	singleton.r, singleton.c = reporter.CreateReporter(singleton.l)

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
