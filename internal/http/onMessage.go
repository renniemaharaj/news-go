package http

import (
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/renniemaharaj/news-go/internal/http/commands"
	"github.com/renniemaharaj/news-go/internal/log"
)

type Instance struct {
	mu          sync.Mutex
	locked      bool
	lastHandled time.Time
	minInterval time.Duration
}

// Creates message handler
func CreateMessageHandler(minInterval time.Duration) *Instance {
	return &Instance{
		minInterval: minInterval,
	}
}

// Primary message handling function that forces limits and protects command handler
func (i *Instance) HandleMessage(con *websocket.Conn, message []byte, l *log.Logger) {
	now := time.Now()

	i.mu.Lock()
	if i.locked || now.Sub(i.lastHandled) < i.minInterval {
		i.mu.Unlock()
		l.Error("Skipping message: handler is busy or message came too soon")
		return
	}

	i.locked = true
	i.lastHandled = now
	i.mu.Unlock()

	defer func() {
		i.mu.Lock()
		i.locked = false
		i.mu.Unlock()
	}()

	l.Debug(fmt.Sprintf("Got message from client: %s", message))

	commands.CommandHandler(con, message, l)
}
