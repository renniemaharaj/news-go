package commands

import (
	"encoding/json"

	"github.com/renniemaharaj/news-go/internal/config"
	"github.com/renniemaharaj/news-go/internal/log"
)

func OptimizeHandler(c *Command, l *log.Logger) [][]byte {
	optimize := &Optimize{}
	if err := json.Unmarshal([]byte(c.Body), optimize); err != nil {
		l.Error("Malformed optimize command body")
		return [][]byte{buildReleaseBlock("releaseOptimizerPrompt")}
	}

	optimized, err := config.Get().OptimizeQueries(optimize.Preferences, optimize.Prompt)
	if err != nil {
		l.Error("Optimization failed: " + err.Error())
		return [][]byte{buildReleaseBlock("releaseOptimizerPrompt")}
	}

	optBytes, err := json.Marshal(optimized)
	if err != nil {
		l.Error("Marshal optimized failed")
		return [][]byte{buildReleaseBlock("releaseOptimizerPrompt")}
	}

	return [][]byte{
		buildDataBlock("optimizedPreferences", optBytes),
		buildReleaseBlock("releaseOptimizerPrompt"),
	}
}
