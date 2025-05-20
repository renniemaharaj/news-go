package news

import (
	"fmt"

	"github.com/renniemaharaj/news/internal/config"
	"github.com/renniemaharaj/news/internal/coordinator"
	"github.com/renniemaharaj/news/internal/log"
	"github.com/renniemaharaj/news/internal/types"
)

type Instance struct {
	l *log.Logger
}

func (n *Instance) HydrateJobs() {
	cfg, err := config.Load("./config.json")

	if err != nil {
		n.l.Error("Error loading config")
	}

	for _, query := range cfg.SearchQueries {
		report := types.Report{}
		coordinator.PreResultsChannel <- report.FromSearchQuery(query)
		n.l.Success(fmt.Sprintf("Hydrated channels with: %s", query))
	}
}

func (n *Instance) CreateLogger() {
	n.l = log.NewLogger(100, true, false)
}

func (n *Instance) GoRoutines() {

	go coordinator.CompleteRoutine(n.l)
	go coordinator.ModelsRoutine(n.l)
	go coordinator.ContentRoutine(n.l)
	go coordinator.JobsRoutine(n.l)
}
