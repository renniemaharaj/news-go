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
	i.l.Info("Analyzing content")
	key := <-apiKeys
	defer func() {
		apiKeys <- key
		i.mu.Unlock()
	}()

	// Gemini API endpoint with model and API key
	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent"

	// Construct Gemini request body
	payload := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"text": msg},
				},
			},
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-goog-api-key", key)

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

	// Parse the Gemini response
	var result struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no content returned")
	}

	outString := result.Candidates[0].Content.Parts[0].Text
	linted, ok := ExtractCodeBlock(outString)
	if !ok {
		i.l.Error(outString)
		return "", fmt.Errorf("AI transform failed: no code block found in output")
	}

	return linted, nil
}
