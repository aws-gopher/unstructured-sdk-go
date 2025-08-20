package unstructured

import (
	"context"
	"fmt"
	"net/http"
)

// DeleteDestination deletes a specific destination connector by its ID
func (c *Client) DeleteDestination(ctx context.Context, id string) error {
	req, err := http.NewRequestWithContext(ctx,
		http.MethodDelete,
		c.endpoint.JoinPath("destinations", id).String(),
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	if err := c.do(req, nil); err != nil {
		return fmt.Errorf("failed to delete destination: %w", err)
	}

	return nil
}
