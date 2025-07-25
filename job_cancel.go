package unstructured

import (
	"context"
	"fmt"
	"net/http"
)

// CancelJob cancels a running job by its ID
func (c *Client) CancelJob(ctx context.Context, id string) error {
	req, err := http.NewRequestWithContext(ctx,
		http.MethodPost,
		c.endpoint.JoinPath("/jobs", id, "cancel").String(),
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	if err := c.do(req, nil); err != nil {
		return fmt.Errorf("failed to cancel job: %w", err)
	}

	return nil
}
