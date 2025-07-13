package model

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/renniemaharaj/news-go/internal/log"
	"github.com/renniemaharaj/news-go/internal/types"
)

func Prompt(msg string, l *log.Logger) (string, error) {
	// Call Python script
	l.Debug(msg)

	cmd := exec.Command("py", "internal/model/ai_transform.py", msg) // or full path
	// cmd.Stdin = bytes.NewReader([]byte(msg))

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr

	// Run the command
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("AI transform failed: %w", err)
	}

	l.Debug(out.String())

	return out.String(), nil
}

func Transform(result *types.Result, l *log.Logger) (types.Result, error) {
	var err error

	if result.TextContent == "" {
		return *result, fmt.Errorf("result had no text content: %w", err)
	}

	// 1. Title
	titlePrompt := fmt.Sprintf("Based on the following article, return a short informative title:\n\n%s", result.TextContent)
	if result.Title, err = Prompt(titlePrompt, l); err != nil {
		return *result, fmt.Errorf("title generation failed: %w", err)
	}

	// 2. Alignment
	alignPrompt := fmt.Sprintf("Rate the following text's alignment with this Christian framework: (1) God exists, (2) KJV Bible is truth, (3) Jesus Christ is God. Return an integer from 0 to 10:\n\n%s", result.TextContent)
	if alignStr, err := Prompt(alignPrompt, l); err == nil {
		fmt.Sscanf(alignStr, "%d", &result.Alignment)
	} else {
		return *result, fmt.Errorf("alignment generation failed: %w", err)
	}

	// 3. Tags
	tagsPrompt := fmt.Sprintf("Suggest 3 relevant tags (separated by commas) for the following content:\n\n%s", result.TextContent)
	if tagsStr, err := Prompt(tagsPrompt, l); err == nil {
		result.Tags = parseCommaList(tagsStr)
	} else {
		return *result, fmt.Errorf("tags generation failed: %w", err)
	}

	// 4. Political Biases
	biasPrompt := fmt.Sprintf("Identify two political biases in the following text (e.g., conservative, leftist):\n\n%s", result.TextContent)
	if biasesStr, err := Prompt(biasPrompt, l); err == nil {
		result.PoliticalBiases = parseCommaList(biasesStr)
	} else {
		return *result, fmt.Errorf("bias generation failed: %w", err)
	}

	// 5. Summary
	summaryPrompt := fmt.Sprintf("Write a thoughtful summary of the following content, highlighting any alignment or contrast with the Christian faith:\n\n%s", result.TextContent)
	if result.Summary, err = Prompt(summaryPrompt, l); err != nil {
		return *result, fmt.Errorf("summary generation failed: %w", err)
	}

	return *result, nil
}
