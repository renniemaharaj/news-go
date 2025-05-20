package browser

import (
	"fmt"
	"net/http"

	"github.com/renniemaharaj/news/internal/log"
)

// Domains to skip during search results
var skipDomains = map[string]struct{}{
	"maps.google.com":   {},
	"photos.google.com": {},
}

// Creates an HTTP request with a user-agent
func httpRequester(target string) *http.Request {
	req, _ := http.NewRequest("GET", target, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	return req
}

// Main function to perform a request
func Request(url string, l *log.Logger) (*http.Response, error) {
	l.Info(fmt.Sprintf("Visiting site for body scraping: %s", url))

	req := httpRequester(url)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		l.Error(fmt.Sprintf("Error scraping site body from: %s", url))
		return nil, err
	}

	return resp, nil
}
