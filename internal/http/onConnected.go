package http

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/renniemaharaj/news-go/internal/coordinator"
	"github.com/renniemaharaj/news-go/internal/log"
)

// Handler called on connection
func connectedHandler(con *websocket.Conn, l *log.Logger) {
	reports := coordinator.Get().Store.AllReports()

	reportBytes, err := json.Marshal(reports)
	if err != nil {
		l.Error("Error occurred while marshalling reports")
	}

	l.Debug("Writing reports to connected client on connected")
	con.WriteMessage(websocket.TextMessage, reportBytes)
}
