package model

import (
	"encoding/json"
	"fmt"

	"github.com/renniemaharaj/news-go/internal/document"
	"github.com/renniemaharaj/news-go/internal/instructions"
	"github.com/renniemaharaj/news-go/internal/log"
)

func Transform(result *document.Result, l *log.Logger) (document.Result, error) {
	var err error

	if result.TextContent == "" {
		return *result, fmt.Errorf("result had no text content")
	}

	transformPrompt, err := instructions.BuildTransformPrompt(result.TextContent, result.Images)
	if err != nil {
		return *result, fmt.Errorf("result had no text content: %s", err.Error())
	}

	transformResponse, err := Get().Prompt_Py(transformPrompt)

	if err != nil {
		return *result, fmt.Errorf("transformation failed: %s", err.Error())
	}

	transformed := Transformed{}
	err = json.Unmarshal([]byte(transformResponse), &transformed)
	if err != nil {
		l.Error(transformResponse)
		return *result, fmt.Errorf("transformation failed, got malformed object")
	}

	if transformed.InSufficientContent {
		return *result, fmt.Errorf("report flagged as insufficient content or not analyzable")
	}

	result.Title = transformed.Title
	result.Alignment = transformed.Alignment
	result.Tags = transformed.Tags
	result.PoliticalBiases = transformed.PoliticalBiases
	result.Summary = transformed.Summary
	result.Images = transformed.Images

	//Cleanup result raw text
	result.TextContent = ""

	l.Success("A result was successfully transformed!")
	return *result, nil
}
