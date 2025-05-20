package main

import "github.com/renniemaharaj/news/pkg/news"

func main() {
	n := news.Instance{}

	n.CreateLogger()
	go n.GoRoutines()
	n.HydrateJobs()

	select {} // blocks forever
}
