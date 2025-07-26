package commands

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/renniemaharaj/grouplogs/pkg/logger"
	"github.com/renniemaharaj/news-go/internal/coordinator"
	"github.com/renniemaharaj/news-go/internal/document"
)

func FeedHandler(c *Command, con *websocket.Conn, l *logger.Logger) {
	feed := &Feed{}
	if err := json.Unmarshal([]byte(c.Body), feed); err != nil {
		l.Error("Malformed feed command body")
	}

	//It's safe to check without verifying urlReportTitle, empty string return not found
	reportInTitle, found := coordinator.Get().Store.ReportByTitle(feed.URLReportTitle)

	// If the report is found, send the report
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

	coordinator.Get().Store.ForEach(func(_ string, r *document.Report) {
		// If we have reached the maximum number of items to send at once, stop processing
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
}
