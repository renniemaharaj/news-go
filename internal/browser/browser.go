package browser

import (
	"sync"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"

	"github.com/renniemaharaj/news-go/internal/log"
)

type Instance struct {
	rod     *rod.Browser
	once    sync.Once
	initErr error
	l       *log.Logger
	m       sync.Mutex
}

var singleton *Instance

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

// Initialize launches the singleton browser instance
func Initialize() error {
	singleton = &Instance{}
	singleton.l = createLogger()
	singleton.once.Do(func() {
		path := launcher.New().Headless(true).
			Leakless(true).
			Set("disable-blink-features", "AutomationControlled").
			MustLaunch()

		singleton.rod = rod.New().ControlURL(path)
		singleton.initErr = singleton.rod.Connect()
		if singleton.initErr != nil {
			singleton.l.Error("Browser failed to connect: " + singleton.initErr.Error())
		}
	})
	return singleton.initErr
}

// Get returns the singleton instance, initializing if necessary
func Get() *Instance {
	if singleton == nil {
		if err := Initialize(); err != nil {
			return nil
		}
	}
	return singleton
}
