package unstructured

import (
	"context"
	"fmt"
	"net/http"
)

// GetDestination retrieves detailed information for a specific destination connector by its ID
func (c *Client) GetDestination(ctx context.Context, id string) (*Destination, error) {
	req, err := http.NewRequestWithContext(ctx,
		http.MethodGet,
		c.endpoint.JoinPath("destinations", id).String(),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	var destination Destination
	if err := c.do(req, &destination); err != nil {
		return nil, fmt.Errorf("failed to get destination: %w", err)
	}

	return &destination, nil
}
