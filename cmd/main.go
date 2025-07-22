package main

import (
	"os"

	"github.com/renniemaharaj/news-go/internal/coordinator"
	"github.com/renniemaharaj/news-go/internal/log"

	"github.com/renniemaharaj/news-go/internal/health"
	"github.com/renniemaharaj/news-go/internal/server"
)

func createLogger() *log.Logger {
	return log.CreateLogger("Coordinator", 100, true, false, false)
}

func main() {
	coordinator.Initialize()
	log.InitGlobalLogger()

	// cloudflare.Initialize()
	// router.LoginToRouter()

	go server.Serve()
	health.HealthCheckScheduler(os.Getenv("WHO_AM_I"), createLogger())
	select {}
}
