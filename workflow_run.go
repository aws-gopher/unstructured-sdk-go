package unstructured

import (
	"context"
	"fmt"
	"net/http"
)

// RunWorkflowRequest represents the request to run a workflow
type RunWorkflowRequest struct {
	InputFiles []string
}

// RunWorkflow runs a workflow by triggering a new job
func (c *Client) RunWorkflow(ctx context.Context, workflowID string, _ *RunWorkflowRequest) (*Job, error) {
	// For now, we'll implement a simple version without file uploads
	// The actual implementation would need multipart form data handling
	req, err := http.NewRequestWithContext(ctx,
		http.MethodPost,
		c.endpoint.JoinPath("/workflows", workflowID, "run").String(),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	var job Job
	if err := c.do(req, &job); err != nil {
		return nil, fmt.Errorf("failed to run workflow: %w", err)
	}

	return &job, nil
}
