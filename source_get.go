package unstructured

import (
	"context"
	"fmt"
	"net/http"
)

// GetSource retrieves detailed information for a specific source connector by its ID
func (c *Client) GetSource(ctx context.Context, id string) (*Source, error) {
	req, err := http.NewRequestWithContext(ctx,
		http.MethodGet,
		c.endpoint.JoinPath("/sources", id).String(),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	var source Source
	if err := c.do(req, &source); err != nil {
		return nil, fmt.Errorf("failed to get source: %w", err)
	}

	return &source, nil
}
