package server

import (
	"github.com/renniemaharaj/news-go/internal/http"
	"github.com/renniemaharaj/news-go/internal/loggers"
)

func Serve() {
	http.ServeFrontend(loggers.LOGGER_SERVER)
}
