package commands

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/renniemaharaj/grouplogs/pkg/logger"
	"github.com/renniemaharaj/news-go/internal/loggers"
)

func logHandler(con *websocket.Conn, l *logger.Logger) {
	subscription := loggers.GROUP_PUBLIC.Delegate.Subscribe()
	l.Info("A client subscribed to the public logger")

	for {
		logArr := &[]logger.Line{
			<-subscription.C,
		}
		logBytes, _ := json.Marshal(logArr)
		if err := con.WriteMessage(websocket.TextMessage, []byte(logBytes)); err != nil {
			loggers.GROUP_PUBLIC.Delegate.Unsubscribe(subscription)
			l.Info("A client was unsubscribed from the public logger")
			break
		}
	}
}
