package unstructured

import (
	"context"
	"fmt"
	"net/http"
)

// GetJob retrieves detailed information for a specific job by its ID
func (c *Client) GetJob(ctx context.Context, id string) (*Job, error) {
	req, err := http.NewRequestWithContext(ctx,
		http.MethodGet,
		c.endpoint.JoinPath("jobs", id).String(),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	var job Job
	if err := c.do(req, &job); err != nil {
		return nil, fmt.Errorf("failed to get job: %w", err)
	}

	return &job, nil
}
