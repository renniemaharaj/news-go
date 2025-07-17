package browser

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/renniemaharaj/news-go/internal/log"
)

// Skip known domains that aren’t news sources
var skipDomains = map[string]struct{}{
	"www.google.com":                 {},
	"policies.google.com":            {},
	"www.google.tt":                  {},
	"accounts.google.com":            {},
	"support.google.com":             {},
	"webcache.googleusercontent.com": {},
	"photos.google.com":              {},
	"maps.google.com":                {},
}

// Browser searching method, returns search results using Rod
func Search(query string, numSitesPerQuery int, l *log.Logger) ([]string, error) {
	searchURL := "https://www.google.com/search?q=" + url.QueryEscape(query) + "&num=10&tbm=nws&tbs=qdr:d"

	path := launcher.New().Headless(true).MustLaunch()
	browser := rod.New().ControlURL(path).MustConnect()
	defer browser.MustClose()

	page := browser.MustPage(searchURL)

	// Set user agent via CDP protocol
	spoofBrowser(page, l)

	// Wait for search result blocks
	page.MustWaitLoad().MustWaitIdle()
	page.MustElement("div#search") // Ensure results are loaded

	// Extract top N results
	links := page.MustElements("a")
	var results []string

	// Link vetting
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
			// Extract from Google redirect format
			extracted = strings.SplitN(strings.TrimPrefix(raw, "/url?q="), "&", 2)[0]
		} else if strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") {
			// Use as-is
			extracted = raw
		} else {
			continue // Skip things like "/search", "#", "javascript:void(0)", etc.
		}

		parsed, err := url.Parse(extracted)
		if err != nil || parsed.Scheme == "" || parsed.Host == "" {
			l.Warning(fmt.Sprintf("Skipping invalid URL: %s", extracted))
			continue
		}

		// Optional: Skip known domains (if defined elsewhere)
		if _, skip := skipDomains[parsed.Hostname()]; skip {
			continue
		}

		// Optional: HEAD request to validate it's a real article (can be slow)
		respCheck, err := http.Head(extracted)
		if err != nil || respCheck.StatusCode == http.StatusNotFound {
			l.Warning(fmt.Sprintf("Skipping 404 or dead link: %s", extracted))
			continue
		}

		l.Info(fmt.Sprintf("Found article: %s", extracted))
		results = append(results, extracted)
	}

	return results, nil
}
