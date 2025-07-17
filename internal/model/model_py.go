package model

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sync"

	"github.com/renniemaharaj/news-go/internal/log"
)

type Instance struct {
	mu sync.Mutex
	l  *log.Logger
}

var modelFile = "local.py"

func (i *Instance) Prompt_Py(msg string) (string, error) {
	i.mu.Lock()
	i.l.Debug("Entering transformation request stage")
	defer i.mu.Unlock()

	// Call Python script
	cmd := exec.Command("py", fmt.Sprintf("internal/model/models/%s", modelFile), msg) // or full path
	// cmd.Stdin = bytes.NewReader([]byte(msg))

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr

	// Run the command
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("AI transform failed: %w", err)
	}

	outString := out.String()

	linted := LintCodeFences(&outString, "json")

	i.l.Debug(*linted)

	return *linted, nil
}
