package model

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sync"

	"github.com/renniemaharaj/news-go/internal/loggers"
)

type Instance struct {
	mu sync.Mutex
}

var modelFile = "gemma.py"

func (i *Instance) Prompt_Py(msg string) (string, error) {
	i.mu.Lock()
	loggers.LOGGER_TRANSFORMER.Info("Analyzing content")
	key := <-apiKeys
	defer func() {
		apiKeys <- key
		i.mu.Unlock()
	}()

	// Call Python script
	cmd := exec.Command("py", fmt.Sprintf("internal/model/models/%s", modelFile), key, msg) // or full path
	// cmd.Stdin = bytes.NewReader([]byte(msg))

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr

	// Run the command
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("AI transform failed: %w", err)
	}

	outString := out.String()

	linted, ok := ExtractCodeBlock(outString)
	if !ok {
		return "", fmt.Errorf("AI transform failed: no code block found in output")
	}

	return linted, nil
}
