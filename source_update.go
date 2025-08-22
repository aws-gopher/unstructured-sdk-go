package unstructured

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// UpdateSourceRequest represents the request to update a source connector.
type UpdateSourceRequest struct {
	ID     string
	Config SourceConfig
}

// UpdateSource updates the configuration of an existing source connector.
// It returns the updated source connector.
func (c *Client) UpdateSource(ctx context.Context, in UpdateSourceRequest) (*Source, error) {
	config, err := json.Marshal(in.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config: %w", err)
	}

	wrapper := struct {
		Config json.RawMessage `json:"config"`
	}{
		Config: json.RawMessage(config),
	}

	body, err := json.Marshal(wrapper)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal update request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx,
		http.MethodPut,
		c.endpoint.JoinPath("sources", in.ID).String(),
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	var source Source
	if err := c.do(req, &source); err != nil {
		return nil, fmt.Errorf("failed to update source: %w", err)
	}

	return &source, nil
}
