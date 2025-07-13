package health

import (
	"fmt"
	"net/http"
	"time"

	"github.com/renniemaharaj/news-go/internal/log"
)

// HealthHandler responds to healthcheck requests
func HealthHandler(version string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "OK %s", version)
	}
}

func startHealthPulse(apiURL string, l *log.Logger) {

	go func() {
		ticker := time.NewTicker(time.Minute / 2)
		defer ticker.Stop()

		client := &http.Client{}

		for range ticker.C {
			if apiURL == "" {
				l.Error("API Address not set for health check")
				continue
			}

			resp, err := client.Get(fmt.Sprintf("%s/healthcheck", apiURL))
			if err != nil {
				l.Error(fmt.Sprintf("Health check failed: %v", err))
				continue
			}
			resp.Body.Close()
			l.Success("Health check passed")
		}
	}()
}
