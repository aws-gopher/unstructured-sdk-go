package unstructured

import (
	"context"
	"fmt"
	"net/http"
)

// GetJobDetails retrieves detailed processing information for a specific job by its ID.
// It returns a JobDetails struct with node stats and processing status.
func (c *Client) GetJobDetails(ctx context.Context, id string) (*JobDetails, error) {
	req, err := http.NewRequestWithContext(ctx,
		http.MethodGet,
		c.endpoint.JoinPath("jobs", id, "details").String(),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	var details JobDetails
	if err := c.do(req, &details); err != nil {
		return nil, fmt.Errorf("failed to get job details: %w", err)
	}

	return &details, nil
}
