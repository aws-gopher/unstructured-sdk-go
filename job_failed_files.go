package unstructured

import (
	"context"
	"fmt"
	"net/http"
)

// GetJobFailedFiles retrieves the list of failed files for a specific job by its ID.
// It returns a JobFailedFiles struct containing the failed files and error messages.
func (c *Client) GetJobFailedFiles(ctx context.Context, jobID string) (*JobFailedFiles, error) {
	req, err := http.NewRequestWithContext(ctx,
		http.MethodGet,
		c.endpoint.JoinPath("/jobs", jobID, "failed-files").String(),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	var failedFiles JobFailedFiles
	if err := c.do(req, &failedFiles); err != nil {
		return nil, fmt.Errorf("failed to get job failed files: %w", err)
	}

	return &failedFiles, nil
}
