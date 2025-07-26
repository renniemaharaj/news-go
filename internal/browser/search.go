package browser

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	grf "github.com/renniemaharaj/go-rod-fast/pkg/browser"
	"github.com/renniemaharaj/news-go/internal/loggers"
)

// Browser searching method, returns search results using Rod
func Search(query string, numSitesPerQuery int) ([]string, error) {
	searchURL := "https://www.google.com/search?q=" + url.QueryEscape(query) + "&num=10&tbm=nws&tbs=qdr:d"

	b := Get()
	page, err := b.NewPage(grf.ProtoTargetFromStr(searchURL), context.Background())
	checkError(err)
	defer page.Close()

	page.MustElement("div#search")

	links := page.MustElements("a")
	var results []string

	for _, link := range links {
		if len(results) >= numSitesPerQuery {
			break
		}
		href, err := link.Attribute("href")
		if err != nil || href == nil {
			continue
		}

		raw := strings.TrimSpace(*href)
		var extracted string
		if strings.HasPrefix(raw, "/url?q=") {
			extracted = strings.SplitN(strings.TrimPrefix(raw, "/url?q="), "&", 2)[0]
		} else if strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") {
			extracted = raw
		} else {
			continue
		}

		parsed, err := url.Parse(extracted)
		if err != nil || parsed.Scheme == "" || parsed.Host == "" {
			loggers.LOGGER_BROWSER.Warning(fmt.Sprintf("Skipping invalid URL: %s", extracted))
			continue
		}

		if _, skip := skipDomains[parsed.Hostname()]; skip {
			continue
		}

		respCheck, err := http.Head(extracted)
		if err != nil || respCheck.StatusCode == http.StatusNotFound {
			loggers.LOGGER_BROWSER.Warning(fmt.Sprintf("Skipping 404 or dead link: %s", extracted))
			continue
		}

		results = append(results, extracted)
	}

	loggers.LOGGER_BROWSER.Info(fmt.Sprintf("Found %d results for query: %s", len(results), query))
	return results, nil
}
