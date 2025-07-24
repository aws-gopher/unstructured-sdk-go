package unstructured

import (
	"context"
	"fmt"
	"net/http"
)

// CreateDestinationConnectionCheck initiates a connection check for a destination connector by its ID.
// It returns a DagNodeConnectionCheck with the status of the check.
func (c *Client) CreateDestinationConnectionCheck(ctx context.Context, id string) (*DagNodeConnectionCheck, error) {
	req, err := http.NewRequestWithContext(ctx,
		http.MethodPost,
		c.endpoint.JoinPath("/destinations", id, "connection-check").String(),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	var check DagNodeConnectionCheck
	if err := c.do(req, &check); err != nil {
		return nil, fmt.Errorf("failed to create destination connection check: %w", err)
	}

	return &check, nil
}

// GetDestinationConnectionCheck retrieves the status of a connection check for a destination connector by its ID.
// It returns a DagNodeConnectionCheck with the current status and reason if any.
func (c *Client) GetDestinationConnectionCheck(ctx context.Context, id string) (*DagNodeConnectionCheck, error) {
	req, err := http.NewRequestWithContext(ctx,
		http.MethodGet,
		c.endpoint.JoinPath("/destinations", id, "connection-check").String(),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	var check DagNodeConnectionCheck
	if err := c.do(req, &check); err != nil {
		return nil, fmt.Errorf("failed to get destination connection check: %w", err)
	}

	return &check, nil
}
