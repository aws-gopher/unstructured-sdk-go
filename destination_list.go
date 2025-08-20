package unstructured

import (
	"context"
	"fmt"
	"net/http"
)

// ListDestinations retrieves a list of available destination connectors
func (c *Client) ListDestinations(ctx context.Context, typ string) ([]Destination, error) {
	req, err := http.NewRequestWithContext(ctx,
		http.MethodGet,
		c.endpoint.JoinPath("destinations").String(),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	if typ != "" {
		q := req.URL.Query()
		q.Add("destination_type", typ)
		req.URL.RawQuery = q.Encode()
	}

	var destinations []Destination
	if err := c.do(req, &destinations); err != nil {
		return nil, fmt.Errorf("failed to list destinations: %w", err)
	}

	return destinations, nil
}
