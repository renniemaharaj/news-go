package browser

import (
	"context"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/renniemaharaj/news-go/internal/log"
)

type Instance struct {
	rod     *rod.Browser
	once    sync.Once
	initErr error
	l       *log.Logger
}

var singleton *Instance

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

// Read performs a headless spoofed browser request and returns text + thumbnails
func (i *Instance) Read(targetURL string) (string, []string, error) {

	page := i.rod.MustPage("")
	defer page.MustClose()

	spoofBrowser(page, i.l)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err := page.Context(ctx).Navigate(targetURL)
	if err != nil {
		i.l.Error("Navigation failed: " + err.Error())
		return "", nil, err
	}

	page.MustWaitLoad()

	body, err := page.Element("body")
	if err != nil {
		i.l.Warning("Could not read site: " + err.Error())
		return "", nil, err
	}

	text, err := body.Text()
	if err != nil {
		i.l.Warning("Could not extract text: " + err.Error())
		return "", nil, err
	}

	imgEls, _ := page.Elements("img")
	var thumbs []string
	for _, img := range imgEls {
		src, _ := img.Attribute("src")
		if src != nil && isLikelyThumbnail(*src) {
			thumbs = append(thumbs, resolveURL(*src, targetURL))
		}
	}

	return text, thumbs, nil
}
