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

// LintCodeFences removes ```html from the start and ``` from the end of the input string.
func LintCodeFences(input *string, language string) *string {
	codeFenceStart := fmt.Sprintf("```%v", language)
	const codeFenceEnd = "```"

	// trim the starting "```html"
	*input = strings.TrimPrefix(*input, codeFenceStart)

	// trim any leading/trailing whitespace or newlines to better detect the ending code fence
	*input = strings.TrimSpace(*input)

	// trim the ending "```"
	*input = strings.TrimSuffix(*input, codeFenceEnd)

	// trim excess whitespace again
	trimmedInput := strings.TrimSpace(*input)

	return &trimmedInput
}
