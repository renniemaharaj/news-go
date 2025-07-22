package commands

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/renniemaharaj/news-go/internal/log"
)

func logHandler(con *websocket.Conn, l *log.Logger) {
	subscription := log.GlobalLogger.Subscribers.Subscribe()

	l.Info("A client subscribed to the global logger")
	for {
		logArr := &[]log.Line{
			<-subscription.C,
		}
		logBytes, _ := json.Marshal(logArr)
		if err := con.WriteMessage(websocket.TextMessage, []byte(logBytes)); err != nil {
			log.GlobalLogger.Subscribers.Unsubscribe(subscription)
			l.Info("A client was unsubscribed from the global logger")
			break
		}
	}
}
