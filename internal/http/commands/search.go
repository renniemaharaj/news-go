package commands

import (
	"fmt"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/renniemaharaj/news-go/internal/config"
	"github.com/renniemaharaj/news-go/internal/coordinator"
	"github.com/renniemaharaj/news-go/internal/document"
	"github.com/renniemaharaj/news-go/internal/log"
	"github.com/renniemaharaj/news-go/internal/reporter"
)

func searchHandler(c *Command, con *websocket.Conn, l *log.Logger) {
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
			adaptPrompt := "Please extend master list to include the attached search query"
			config.Get().OptimizeQueries([]string{report.SearchQuery}, adaptPrompt)

			return
		}

		con.WriteMessage(websocket.TextMessage, buildErrorBlock("Your request was not successful"))
	}()

	r.TODO_SEARCH_CHANNEL <- document.ReportFromQuery(c.Body)
}
