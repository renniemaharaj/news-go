package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (i *Instance) Prompt(msg string) (string, error) {
	i.mu.Lock()
	i.l.Debug("Entering transformation request stage")
	defer i.mu.Unlock()

	// Define payload structure
	payload := map[string]interface{}{
		"model": "smollm3-3b", // or dynamic if needed
		// "model":       "llama-3.2-3b-instruct",
		"messages":    []map[string]string{{"role": "user", "content": msg}},
		"temperature": 0.7,
		"max_tokens":  -1,
		"stream":      false,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Create POST request
	req, err := http.NewRequest("POST", "http://localhost:1234/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("bad response (%d): %s", resp.StatusCode, body)
	}

	// Parse the response
	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no choices returned")
	}

	outString := result.Choices[0].Message.Content
	linted := LintCodeFences(&outString, "json")

	i.l.Debug(*linted)
	return *linted, nil
}
