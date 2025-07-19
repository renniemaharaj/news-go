package model

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func parseCommaList(s string) []string {
	parts := strings.Split(s, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

// The regex pattern to extract the inner content of a code block.
// Place somewhere in your .env file or otherwise
// (?s)\s*```(?:json)?\s*(\{.*?\})\s*```
func ExtractCodeBlock(input string) (string, bool) {
	pattern := os.Getenv("CODE_FENCE_INNER_REGEX")
	if pattern == "" {
		log.Fatal("CODE_FENCE_INNER_REGEX not set in environment")
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		log.Fatalf("Invalid regex pattern: %v", err)
	}
	matches := re.FindStringSubmatch(input)

	if len(matches) < 2 {
		return "", false
	}

	return matches[1], true
}

// LintCodeFences removes ```<lang> and ``` fences from an input string.
func LintCodeFences(input *string, language string) *string {
	// Normalize input
	trimmed := strings.TrimSpace(*input)

	// Remove starting fence (with or without language, with optional newline)
	prefixWithLang := fmt.Sprintf("```%s\n", language)
	prefixPlain := "```\n"
	if strings.HasPrefix(trimmed, prefixWithLang) {
		trimmed = strings.TrimPrefix(trimmed, prefixWithLang)
	} else if strings.HasPrefix(trimmed, prefixPlain) {
		trimmed = strings.TrimPrefix(trimmed, prefixPlain)
	} else if strings.HasPrefix(trimmed, "```"+language) {
		trimmed = strings.TrimPrefix(trimmed, "```"+language)
	} else if strings.HasPrefix(trimmed, "```") {
		trimmed = strings.TrimPrefix(trimmed, "```")
	}

	// Remove trailing fence (with or without preceding newline)
	trimmed = strings.TrimSpace(trimmed)
	trimmed = strings.TrimSuffix(trimmed, "```")

	trimmed = strings.TrimSpace(trimmed)
	return &trimmed
}
