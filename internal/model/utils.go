package model

import (
	"fmt"
	"strings"
)

func parseCommaList(s string) []string {
	parts := strings.Split(s, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
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
