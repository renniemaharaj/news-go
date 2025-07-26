package main

import (
	"os"

	"github.com/renniemaharaj/news-go/internal/coordinator"
	"github.com/renniemaharaj/news-go/internal/loggers"

	"github.com/renniemaharaj/news-go/internal/health"
	"github.com/renniemaharaj/news-go/internal/server"
)

func main() {
	loggers.Initialize()
	coordinator.Initialize()

	go server.Serve()
	health.HealthCheckScheduler(os.Getenv("WHO_AM_I"), loggers.LOGGER_COORDINATOR)
	select {}
}
