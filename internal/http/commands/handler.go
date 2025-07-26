package commands

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/renniemaharaj/grouplogs/pkg/logger"
)

func CommandHandler(con *websocket.Conn, message []byte, l *logger.Logger) {
	var c Command
	if err := json.Unmarshal(message, &c); err != nil {
		l.Error("Invalid command payload from client")
		return
	}

	l.Success(fmt.Sprintf("Command received: %s", c.Name))

	// var responses [][]byte

	switch c.Name {
	case "search":
		searchHandler(&c, con, l)
		return // async handled
	case "feed":
		FeedHandler(&c, con, l)
	case "log":
		logHandler(con, l)
	default:
		con.WriteMessage(websocket.TextMessage, []byte(buildErrorBlock("Unknown command")))
	}
}
