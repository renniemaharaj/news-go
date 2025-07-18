package browser

import (
	"context"
	"fmt"
	"time"
)

// import "fmt"

// Content performs a headless spoofed browser request and returns text + thumbnails
func (i *Instance) Content(targetURL string) (string, []string, error) {
	i.m.Lock()
	defer i.m.Unlock()
	// Panic-safe wrapper
	defer func() {
		if r := recover(); r != nil {
			i.l.Error(fmt.Sprintf("Recovered from panic: %v", r))

			// Safely reinitialize browser
			i.m.Lock()
			defer i.m.Unlock()

			i.l.Info("Reinitializing browser after panic")
			if err := Initialize(); err != nil {
				i.l.Error("Failed to reinitialize browser: " + err.Error())
				return
			}
		}
	}()

	if i.rod == nil {
		Initialize()
		return Get().Content(targetURL)
	}

	i.l.Info("Fetching content from: " + targetURL)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	i.rod.Context(ctx)
	defer cancel()

	page := i.rod.MustPage().MustWaitLoad()
	defer page.MustClose()

	page.MustNavigate(targetURL)

	spoofBrowser(page, i.l)
	page.WaitLoad()

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
