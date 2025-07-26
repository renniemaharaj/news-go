package browser

import (
	grf "github.com/renniemaharaj/go-rod-fast/pkg/browser"
	"github.com/renniemaharaj/news-go/internal/loggers"
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

var singleton *grf.Browser

func checkError(err error) {
	if err != nil {
		loggers.LOGGER_BROWSER.Fatal(err)
	}
}

func Get() *grf.Browser {
	if singleton == nil {
		singleton = grf.Create(false)
	}

	return singleton
}
