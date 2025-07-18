package config

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/renniemaharaj/news-go/internal/log"
)

type Instance struct {
	Port             string   `json:"port"`
	Debugging        bool     `json:"debugging"`
	SearchQueries    []string `json:"searchQueries"`
	NumSitesPerQuery int      `json:"num_sites_per_query"`

	l  *log.Logger
	mu sync.Mutex
}

const configPath = "./config.json"

var singleton *Instance

// Initialize must be called once at program start
func Initialize() {
	file, err := os.ReadFile(configPath)
	if err != nil {
		panic(fmt.Sprintf("Expected configuration file at: %s", configPath))
	}

	var cfg Instance
	if err := json.Unmarshal(file, &cfg); err != nil {
		panic("Configuration is corrupted")
	}

	cfg.l = createLogger()
	singleton = &cfg
}

// Get returns the singleton config instance
func Get() *Instance {
	if singleton == nil {
		Initialize()
	}
	return singleton
}

// Write persists the current config to disk (mutex-protected)
func (i *Instance) Write() error {
	i.mu.Lock()
	defer i.mu.Unlock()

	data, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	i.l.Success("Configuration file written")
	return nil
}
