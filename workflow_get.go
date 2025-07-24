package unstructured

import (
	"context"
	"fmt"
	"net/http"
)

// GetWorkflow retrieves detailed information for a specific workflow by its ID
func (c *Client) GetWorkflow(ctx context.Context, workflowID string) (*Workflow, error) {
	req, err := http.NewRequestWithContext(ctx,
		http.MethodGet,
		c.endpoint.JoinPath("/workflows", workflowID).String(),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	var workflow Workflow
	if err := c.do(req, &workflow); err != nil {
		return nil, fmt.Errorf("failed to get workflow: %w", err)
	}

	return &workflow, nil
}
