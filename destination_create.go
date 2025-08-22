package unstructured

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// CreateDestination creates a new destination connector with the specified configuration.
// It returns the created destination connector with its assigned ID and metadata.
func (c *Client) CreateDestination(ctx context.Context, in CreateDestinationRequest) (*Destination, error) {
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
		return nil, fmt.Errorf("failed to marshal destination request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx,
		http.MethodPost,
		c.endpoint.JoinPath("destinations/").String(),
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	var destination Destination
	if err := c.do(req, &destination); err != nil {
		return nil, fmt.Errorf("failed to create destination: %w", err)
	}

	return &destination, nil
}

// CreateDestinationRequest represents a request to create a new destination connector.
// It contains the name, type, and configuration for the destination.
type CreateDestinationRequest struct {
	Name   string
	Config DestinationConfig
}
