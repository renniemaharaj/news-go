package browser

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/renniemaharaj/news-go/internal/log"
)

// spoofBrowser applies anti-bot spoofing to the browser page.
// It spoofs user-agent, language, platform, and viewport to mimic a real browser.
func spoofBrowser(page *rod.Page, l *log.Logger) {
	// Spoof user agent
	err := proto.NetworkSetUserAgentOverride{
		UserAgent:      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36",
		AcceptLanguage: "en-US,en;q=0.9",
		Platform:       "Win32",
	}.Call(page)
	if err != nil {
		l.Warning("Failed to spoof user-agent: " + err.Error())
	}

	// Spoof screen resolution and viewport
	err = proto.EmulationSetDeviceMetricsOverride{
		Width:             1920,
		Height:            1080,
		DeviceScaleFactor: 1.0,
		Mobile:            false,
	}.Call(page)
	if err != nil {
		l.Warning("Failed to spoof device metrics: " + err.Error())
	}

	// Override navigator properties using JavaScript
	// This runs before page loads to trick JS-based bot detection
	page.MustEvalOnNewDocument(`
		Object.defineProperty(navigator, 'platform', {
			get: () => 'Win32'
		});
		Object.defineProperty(navigator, 'userAgent', {
			get: () => 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36'
		});
		Object.defineProperty(navigator, 'language', {
			get: () => 'en-US'
		});
		Object.defineProperty(navigator, 'languages', {
			get: () => ['en-US', 'en']
		});
	`)
}
