package commands

import (
	"encoding/json"
	"fmt"
)

type Command struct {
	Name string `json:"name"`
	Body string `json:"body"`
}

type Optimize struct {
	Preferences []string `json:"preferences"`
	Prompt      string   `json:"prompt"`
}

func buildDataBlock(key string, value json.RawMessage) []byte {
	return fmt.Appendf(make([]byte, 0, 64), `[{"%s":%s}]`, key, value)
}

func buildDataBlockString(key string, value string) []byte {
	return fmt.Appendf(make([]byte, 0, 64), `[{"%s":"%s"}]`, key, value)
}

func buildReleaseBlock(name string) []byte {
	return fmt.Appendf(make([]byte, 0, 64), `[{"%s":""}]`, name)
}

func buildErrorBlock(msg string) []byte {
	return fmt.Appendf(make([]byte, 0, 64), `[{"error":"%s"}]`, msg)
}
