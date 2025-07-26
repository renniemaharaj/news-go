package commands

import (
	"fmt"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/renniemaharaj/grouplogs/pkg/logger"
	"github.com/renniemaharaj/news-go/internal/coordinator"
	"github.com/renniemaharaj/news-go/internal/document"
	"github.com/renniemaharaj/news-go/internal/reporter"
)

func searchHandler(c *Command, con *websocket.Conn, l *logger.Logger) {
	r, EXPOSE := reporter.CreateReporter(l)

	go func() {
		report := <-EXPOSE

		if len(report.Results) > 0 {
			coordinator.Get().Store.StoreReport(&report, l)

			href := fmt.Sprintf("/read/%s/%s",
				url.PathEscape(report.SearchQuery),
				url.PathEscape(report.Results[0].Title),
			)

			con.WriteMessage(websocket.TextMessage, buildDataBlockString("href", href))
			return
		}

		con.WriteMessage(websocket.TextMessage, buildErrorBlock("Your request was not successful"))
	}()

	r.TODO_SEARCH_CHANNEL <- document.ReportFromQuery(c.Body)
}
