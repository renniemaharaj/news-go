package commands

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/renniemaharaj/news-go/internal/log"
)

func CommandHandler(con *websocket.Conn, message []byte, l *log.Logger) {
	var c Command
	if err := json.Unmarshal(message, &c); err != nil {
		l.Error("Invalid command payload from client")
		return
	}

	l.Success(fmt.Sprintf("Command received: %s", c.Name))

	var responses [][]byte

	switch c.Name {
	case "search":
		searchHandler(&c, con, l)
		return // async handled
	case "optimizePreferences":
		responses = OptimizeHandler(&c, l)
	default:
		responses = [][]byte{buildErrorBlock("Unknown command")}
	}

	for _, msg := range responses {
		con.WriteMessage(websocket.TextMessage, msg)
	}
}
