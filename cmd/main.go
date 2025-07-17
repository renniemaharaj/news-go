package main

import (
	"github.com/renniemaharaj/news-go/internal/coordinator"
	"github.com/renniemaharaj/news-go/internal/server"
)

func main() {
	coordinator.Initialize()

	go server.Serve()
	select {}
}
