package model

import (
	"encoding/json"
	"fmt"

	"github.com/renniemaharaj/news-go/internal/instructions"
	"github.com/renniemaharaj/news-go/internal/log"
)

func Optimize(masterList, userList []string, userPrompt string, l *log.Logger) (*Optimized, error) {
	var err error

	if userPrompt == "" {
		return nil, fmt.Errorf("result had no text content: %w", err)
	}

	transformPrompt, err := instructions.BuildOptimizationPrompt(masterList, userList, userPrompt)
	if err != nil {
		return nil, fmt.Errorf("result had no text content: %w", err)
	}

	optimizeResponse, err := Get().Prompt(transformPrompt)

	if err != nil {
		return nil, fmt.Errorf("transformation failed: %w", err)
	}

	optimized := Optimized{}
	err = json.Unmarshal([]byte(optimizeResponse), &optimized)
	if err != nil {
		return nil, fmt.Errorf("transformation failed, got malformed object")
	}

	return &optimized, nil
}
