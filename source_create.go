package unstructured

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// CreateSourceRequest represents a request to create a new source connector.
// It contains the name and configuration for the source.
type CreateSourceRequest struct {
	Name   string
	Config SourceConfig
}

// CreateSource creates a new source connector with the specified configuration.
// It returns the created source connector with its assigned ID and metadata.
func (c *Client) CreateSource(ctx context.Context, in CreateSourceRequest) (*Source, error) {
	config, err := json.Marshal(in.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config: %w", err)
	}

	shadow := struct {
		Name   string          `json:"name"`
		Type   string          `json:"type"`
		Config json.RawMessage `json:"config"`
	}{
		Name:   in.Name,
		Type:   in.Config.Type(),
		Config: json.RawMessage(config),
	}

	body, err := json.Marshal(shadow)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal source request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx,
		http.MethodPost,
		c.endpoint.JoinPath("sources/").String(),
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	var source Source
	if err := c.do(req, &source); err != nil {
		return nil, fmt.Errorf("failed to create source: %w", err)
	}

	return &source, nil
}
