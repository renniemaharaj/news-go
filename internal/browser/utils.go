package browser

import (
	"fmt"

	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"

	"github.com/renniemaharaj/news-go/internal/log"
)

func GetLinkAttribute(t html.Token, l *log.Logger) (string, bool) {
	if t.Data != "a" {
		return "", false
	}

	for _, attr := range t.Attr {
		if attr.Key == "href" && strings.HasPrefix(attr.Val, "/url?q=") {
			extracted := strings.Split(attr.Val[7:], "&")[0]
			if strings.HasPrefix(extracted, "http") {
				u, err := url.Parse(extracted)
				if err != nil {
					return "", false
				}

				if _, skip := skipDomains[u.Hostname()]; skip {
					return "", false
				}

				// Check for 404
				respCheck, err := http.Head(extracted)
				if err != nil || respCheck.StatusCode == http.StatusNotFound {
					l.Warning(fmt.Sprintf("Skipping 404 result: %s", u))
					return "", false
				}

				return strings.TrimSpace(extracted), true
			}
		}
	}

	return "", false
}

// Builds a google-based search url using query for five news report
func searchNewsURL(query string) string {
	queryWithContext := fmt.Sprintf("%s - %s", query, time.Now().Format("2006-01-02"))
	joinedQuery := strings.ReplaceAll(queryWithContext, " ", "+")
	return "https://www.google.com/search?q=" + joinedQuery + "&num=5&tbm=nws"
}

func isLikelyThumbnail(src string) bool {
	lower := strings.ToLower(src)

	return !strings.Contains(lower, "logo") &&
		!strings.Contains(lower, "icon") &&
		!strings.Contains(lower, "svg") &&
		!strings.Contains(lower, "placeholder") &&

		//Extensions

		(strings.HasSuffix(lower, ".jpg") ||
			strings.HasSuffix(lower, ".jpeg") ||
			strings.HasSuffix(lower, ".png") ||
			strings.HasSuffix(lower, ".webp"))
}

func resolveURL(link string, base string) string {
	uri, err := url.Parse(link)
	if err != nil {
		return link
	}
	baseURL, err := url.Parse(base)
	if err != nil {
		return link
	}

	return baseURL.ResolveReference(uri).String()
}
