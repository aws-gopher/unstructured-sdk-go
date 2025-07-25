package unstructured

import (
	"context"
	"fmt"
	"net/http"
)

// DeleteWorkflow deletes a workflow by its ID
func (c *Client) DeleteWorkflow(ctx context.Context, id string) error {
	req, err := http.NewRequestWithContext(ctx,
		http.MethodDelete,
		c.endpoint.JoinPath("/workflows", id).String(),
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	if err := c.do(req, nil); err != nil {
		return fmt.Errorf("failed to delete workflow: %w", err)
	}

	return nil
}
