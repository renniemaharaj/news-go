package browser

import (
	"fmt"

	"net/http"
	"net/url"
	"strings"

	"github.com/renniemaharaj/grouplogs/pkg/logger"
	"golang.org/x/net/html"
)

// Responsible for getting link elements's href
func GetLinkAttribute(t html.Token, l *logger.Logger) (string, bool) {
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

func isLikelyThumbnail(src string) bool {
	lower := strings.ToLower(src)

	// First attempt to filter out by filenames
	return !strings.Contains(lower, "logo") &&
		!strings.Contains(lower, "icon") &&
		!strings.Contains(lower, "svg") &&
		!strings.Contains(lower, "placeholder") &&

		// Extensions Allowed
		(strings.HasSuffix(lower, ".jpg") ||
			strings.HasSuffix(lower, ".jpeg") ||
			strings.HasSuffix(lower, ".png") ||
			strings.HasSuffix(lower, ".webp") ||
			strings.HasSuffix(lower, "gif"))

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
