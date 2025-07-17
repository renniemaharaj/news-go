package http

import (
	"fmt"
	"html"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/renniemaharaj/news-go/internal/config"
	"github.com/renniemaharaj/news-go/internal/log"
)

// The upgrader to be used later for upgrading with cors configuration
var wsUpgrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Accepting all origin
	},
}

// Upgrade handler function
func upgradeHandler(w http.ResponseWriter, r *http.Request) {
	connection, _ := wsUpgrade.Upgrade(w, r, nil)

	l := log.CreateLogger("Socket", 100, true, false, false)
	l.Info(fmt.Sprintf("%s client connected", html.EscapeString(r.URL.Path)))

	connectedHandler(connection, l)

	handler := CreateMessageHandler(2 * time.Second) // min 2s between messages

	for {
		mt, message, err := connection.ReadMessage()
		if err != nil || mt == websocket.CloseMessage {
			break // Exit the blocking loop if client tries to close
		}

		handler.HandleMessage(connection, message, l)
	}
	connection.Close() // Close connection
}

var distDir = "./dist"

// Primary open function for starting server
func ServeFrontend(l *log.Logger) {
	http.HandleFunc("/ws", upgradeHandler) // WebSocket route

	// Serve static assets in /dist (e.g., /assets/*.js, .css, etc.)
	fileServer := http.FileServer(http.Dir(distDir))
	http.Handle("/assets/", fileServer) // handles JS, CSS, images, etc.
	http.Handle("/favicon.ico", fileServer)

	// Catch-all for React/Vite SPA routes — serve index.html
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, fmt.Sprintf("%s/index.html", distDir))
	})

	port := config.Get().Port
	l.Info(fmt.Sprintf("Starting: view the frontend on localhost%s\n", port))
	if err := http.ListenAndServe(port, nil); err != nil {
		l.Error(fmt.Sprintf("Error starting server: %s\n", err))
		return
	}

}
