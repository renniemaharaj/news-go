package model

import (
	"encoding/json"

	"github.com/google/generative-ai-go/genai"
	"github.com/renniemaharaj/news/internal/types"
	"github.com/renniemaharaj/news/pkg/transformer/gemi"
)

// Constructs an input for transformer communication
func getInput(result types.Result) gemi.Input {
	contentBytes, err := json.Marshal(result)
	var HISTORY = []*genai.Content{}
	if err == nil {
		return gemi.Input{
			Current: genai.Text(string(contentBytes)),
			History: HISTORY,
			Context: []map[string]string{},
		}
	}

	return gemi.Input{}
}
