package commands

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/renniemaharaj/news-go/internal/coordinator"
	"github.com/renniemaharaj/news-go/internal/document"
	"github.com/renniemaharaj/news-go/internal/log"
)

func FeedHandler(c *Command, con *websocket.Conn, l *log.Logger) {
	feed := &Feed{}
	if err := json.Unmarshal([]byte(c.Body), feed); err != nil {
		l.Error("Malformed feed command body")
	}

	coordinator.Get().Store.ForEach(func(_ string, r *document.Report) {
		if r.HasTagIntersection(feed.Preferences) {
			reportBytes, err := json.Marshal([]document.Report{*r})
			if err != nil {
				l.Error(fmt.Sprintf("Failed to marshal report: %s", err))
				return
			}

			if err := con.WriteMessage(websocket.TextMessage, reportBytes); err != nil {
				l.Error(fmt.Sprintf("Failed to send report: %s", err))
			}
		}
	})

	if feed.URLReportTitle != "" {
		reportInTitle, found := coordinator.Get().Store.ReportByTitle(feed.URLReportTitle)
		if found {
			reportBytes, err := json.Marshal([]document.Report{*reportInTitle})
			if err != nil {
				l.Error(fmt.Sprintf("Failed to marshal report: %s", err))
				return
			}

			if err := con.WriteMessage(websocket.TextMessage, reportBytes); err != nil {
				l.Error(fmt.Sprintf("Failed to send report: %s", err))
			}
		}
	}
}
