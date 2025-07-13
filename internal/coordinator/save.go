package coordinator

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/renniemaharaj/news-go/internal/log"
	"github.com/renniemaharaj/news-go/internal/types"
)

func Save(r *types.Report, l *log.Logger) {
	// Create reports directory if it doesn't exist
	if err := os.MkdirAll("./reports", 0755); err != nil {
		l.Error(fmt.Sprintf(("Failed to create reports directory: %s"), err.Error()))
		return
	}

	// Marshal report to JSON bytes
	jsonBytes, err := json.Marshal(r)
	if err != nil {
		l.Error(fmt.Sprintf(("Failed to marshal report: %s"), err.Error()))
		return
	}

	// Create file name based on title and date
	fileName := fmt.Sprintf("./reports/%s_%s.json",
		strings.ReplaceAll(r.Title, " ", "_"),
		time.Now().Format("20060102_150405"))

	// Write to file
	if err := os.WriteFile(fileName, jsonBytes, 0644); err != nil {
		l.Error(fmt.Sprintf(("Failed to write report: %s"), err.Error()))
		return
	}

	l.Info(fmt.Sprintf(("Saved report: %s"), fileName))
}
