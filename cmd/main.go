package main

import (
	"os"

	"github.com/renniemaharaj/news-go/internal/coordinator"
	"github.com/renniemaharaj/news-go/internal/health"
	"github.com/renniemaharaj/news-go/internal/server"
)

func main() {
	coordinator.Initialize()

	go server.Serve()
	health.HealthCheckScheduler(os.Getenv("WHO_AM_I_API_URL"), createLogger())
	select {}
}
