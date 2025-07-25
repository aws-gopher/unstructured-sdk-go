package unstructured

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// UpdateDestinationRequest represents the request to update a destination connector.
type UpdateDestinationRequest struct {
	ID     string
	Config DestinationConfigInput
}

// UpdateDestination updates the configuration of an existing destination connector.
// It returns the updated destination connector.
func (c *Client) UpdateDestination(ctx context.Context, in UpdateDestinationRequest) (*Destination, error) {
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
		c.endpoint.JoinPath("/destinations", in.ID).String(),
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	var destination Destination
	if err := c.do(req, &destination); err != nil {
		return nil, fmt.Errorf("failed to update destination: %w", err)
	}

	return &destination, nil
}
