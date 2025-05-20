package model

import (
	"context"
	"encoding/json"

	"github.com/renniemaharaj/news/internal/types"
	"github.com/renniemaharaj/news/internal/validation"

	"github.com/renniemaharaj/news/pkg/pool"
)

const (
	queues  = 2
	backoff = 2
)

// Prompt function interfaces with transformer package on our behalf
func Prompt(report types.Result) (types.Result, error) {
	p := pool.Instance{}
	p.InitializePool()

	// call the transformer package, queued, exponential backoff and validation
	resp, err := p.QueuedEVS(context.Background(), getInput(report), validation.Validate, queues, backoff)
	if err != nil {
		return types.Result{}, err
	}

	var reports types.Result

	// queuedEVS already handles validation into
	err = json.Unmarshal([]byte(resp), &reports)
	if err != nil {
		return types.Result{}, err
	}

	return reports, nil
}
