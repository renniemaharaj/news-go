package browser

import (
	"context"

	grf "github.com/renniemaharaj/go-rod-fast/pkg/browser"
	"github.com/renniemaharaj/news-go/internal/loggers"
)

// Content performs a headless spoofed browser request and returns text + thumbnails
func Content(targetURL string) (string, []string, error) {
	b := Get()
	page, err := b.NewPage(grf.ProtoTargetFromStr(targetURL), context.Background())
	checkError(err)
	defer page.Close()

	if err != nil {
		loggers.LOGGER_BROWSER.Warning("Could not get browser instance: " + err.Error())
		return "", nil, err
	}
	body, err := page.Element("body")
	if err != nil {
		loggers.LOGGER_BROWSER.Warning("Could not read site: " + err.Error())
		return "", nil, err
	}

	text, err := body.Text()
	if err != nil {
		loggers.LOGGER_BROWSER.Warning("Could not extract text: " + err.Error())
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
