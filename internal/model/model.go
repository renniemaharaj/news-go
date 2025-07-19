package model

import (
	"encoding/json"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var singleton *Instance
var apiKeys = make(chan string, 100) // buffered channel for efficiency

func Initialize() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, or error loading it:", err)
	}

	// Read and parse JSON array
	var keys []string
	if raw := os.Getenv("API_KEYS"); raw != "" {
		err := json.Unmarshal([]byte(raw), &keys)
		if err != nil {
			log.Println("Failed to parse API_KEYS:", err)
		} else {
			for _, key := range keys {
				apiKeys <- key
			}
		}
	} else {
		log.Println("API_KEYS not set in .env")
	}

	singleton = &Instance{}
	singleton.l = createLogger()
}

// Get returns the singleton instance
func Get() *Instance {
	if singleton == nil {
		Initialize()
	}
	return singleton
}
