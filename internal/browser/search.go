package browser

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Browser searching method, returns search results using Rod
func (i *Instance) Search(query string, numSitesPerQuery int) ([]string, error) {
	i.m.Lock()
	defer i.m.Unlock()

	if i.rod == nil {
		Initialize()
		return Get().Search(query, numSitesPerQuery)
	}

	i.l.Info("Searching for: " + query)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	i.rod.Context(ctx)
	defer cancel()

	searchURL := "https://www.google.com/search?q=" + url.QueryEscape(query) + "&num=10&tbm=nws&tbs=qdr:d"

	page := i.rod.MustPage()
	defer page.MustClose()

	page.MustNavigate(searchURL)

	spoofBrowser(page, i.l)
	page.MustWaitLoad().MustWaitIdle()
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
			i.l.Warning(fmt.Sprintf("Skipping invalid URL: %s", extracted))
			continue
		}

		if _, skip := skipDomains[parsed.Hostname()]; skip {
			continue
		}

		respCheck, err := http.Head(extracted)
		if err != nil || respCheck.StatusCode == http.StatusNotFound {
			i.l.Warning(fmt.Sprintf("Skipping 404 or dead link: %s", extracted))
			continue
		}

		results = append(results, extracted)
	}

	i.l.Info(fmt.Sprintf("Found %d results for query: %s", len(results), query))
	return results, nil
}
