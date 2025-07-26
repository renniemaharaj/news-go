package model

import (
	"encoding/json"
	"os"

	"github.com/joho/godotenv"
	"github.com/renniemaharaj/news-go/internal/loggers"
)

var singleton *Instance
var apiKeys = make(chan string, 100) // buffered channel for efficiency

func Initialize() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		loggers.LOGGER_TRANSFORMER.Fatal(err)
	}

	// Read and parse JSON array
	var keys []string
	if raw := os.Getenv("API_KEYS"); raw != "" {
		err := json.Unmarshal([]byte(raw), &keys)
		if err != nil {
			loggers.LOGGER_TRANSFORMER.Fatal(err)
		} else {
			for _, key := range keys {
				apiKeys <- key
			}
		}
	} else {
		loggers.LOGGER_TRANSFORMER.Error("API_KEYS not set in .env")
	}

	singleton = &Instance{}
}

// Get returns the singleton instance
func Get() *Instance {
	if singleton == nil {
		Initialize()
	}
	return singleton
}
