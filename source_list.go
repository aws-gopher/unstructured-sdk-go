package unstructured

import (
	"context"
	"fmt"
	"net/http"
)

// ListSources retrieves a list of available source connectors
func (c *Client) ListSources(ctx context.Context, typ string) ([]Source, error) {
	req, err := http.NewRequestWithContext(ctx,
		http.MethodGet,
		c.endpoint.JoinPath("/sources").String(),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	if typ != "" {
		q := req.URL.Query()
		q.Add("source_type", typ)
		req.URL.RawQuery = q.Encode()
	}

	var sources []Source
	if err := c.do(req, &sources); err != nil {
		return nil, fmt.Errorf("failed to list sources: %w", err)
	}

	return sources, nil
}
