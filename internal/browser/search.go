package browser

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/renniemaharaj/news-go/internal/log"
	"golang.org/x/net/html"
)

// Browser searching method, return search results
func Search(query string, numSitesPerQuery int, l *log.Logger) ([]string, error) {
	searchURL := "https://www.google.com/search?q=" + strings.ReplaceAll(query, " ", "+") + "&num=5&tbm=nws"
	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	z := html.NewTokenizer(resp.Body)
	var results []string

	for {
		if len(results) >= numSitesPerQuery {
			break
		}

		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			return results, nil
		case html.StartTagToken:
			t := z.Token()
			if t.Data == "a" {
				for _, attr := range t.Attr {
					if attr.Key == "href" && strings.HasPrefix(attr.Val, "/url?q=") {
						extracted := strings.Split(attr.Val[7:], "&")[0]
						if strings.HasPrefix(extracted, "http") {
							u, err := url.Parse(extracted)
							if err != nil {
								continue
							}
							domain := u.Hostname()
							if _, skip := skipDomains[domain]; skip {
								continue
							}

							// Check for 404 before including
							respCheck, err := http.Head(extracted)
							if err != nil || respCheck.StatusCode == http.StatusNotFound {
								l.Warning(fmt.Sprintf("⚠️ Skipping 404 result: %s", u))
								continue
							}

							results = append(results, strings.TrimSpace(extracted))
							break
						}
					}
				}
			}
		}
	}

	return results, nil
}
