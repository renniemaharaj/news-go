package validation

import (
	"encoding/json"
	"fmt"

	"github.com/renniemaharaj/news-go/internal/document"
)

func Validate(resp string) error {
	var result document.Result
	err := json.Unmarshal([]byte(resp), &result)
	if err != nil {
		return err
	}

	// Validate required Model fields
	if result.Alignment < 0 || result.Alignment > 10 {
		return fmt.Errorf("alignment must be between 0 and 10")
	}

	if result.Summary == "" {
		return fmt.Errorf("summary is required")
	}

	if len(result.Tags) == 0 {
		return fmt.Errorf("tags are required")
	}

	if len(result.PoliticalBiases) == 0 {
		return fmt.Errorf("political biases are required")
	}

	// Note: Commentaries, TextContent, Images, and HREF can be omitted as they are
	// not Model-generated fields

	return nil
}
