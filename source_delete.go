package unstructured

import (
	"context"
	"fmt"
	"net/http"
)

// DeleteSource deletes a specific source connector identified by its ID
func (c *Client) DeleteSource(ctx context.Context, id string) error {
	req, err := http.NewRequestWithContext(ctx,
		http.MethodDelete,
		c.endpoint.JoinPath("/sources", id).String(),
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	if err := c.do(req, nil); err != nil {
		return fmt.Errorf("failed to delete source: %w", err)
	}

	return nil
}
