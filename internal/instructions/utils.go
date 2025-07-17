package instructions

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const promptDir = "internal/instructions"

func load(path string) (string, error) {
	bytes, err := os.ReadFile(filepath.Join(promptDir, path))
	if err != nil {
		return "", fmt.Errorf("failed to load %s: %w", path, err)
	}
	return string(bytes), nil
}

func inject(template string, replacements map[string]string) string {
	for key, value := range replacements {
		template = strings.ReplaceAll(template, key, value)
	}
	return template
}
