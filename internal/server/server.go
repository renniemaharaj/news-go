package server

import (
	"github.com/renniemaharaj/news-go/internal/http"
	"github.com/renniemaharaj/news-go/internal/log"
)

func Serve() {
	l := log.CreateLogger("Server", 100, true, false, false)
	http.ServeFrontend(l)
}
