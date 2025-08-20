package unstructured

import (
	"context"
	"fmt"
	"net/http"
)

// CreateSourceConnectionCheck initiates a connection check for a source connector by its ID.
// It returns a DagNodeConnectionCheck with the status of the check.
func (c *Client) CreateSourceConnectionCheck(ctx context.Context, id string) (*DagNodeConnectionCheck, error) {
	req, err := http.NewRequestWithContext(ctx,
		http.MethodPost,
		c.endpoint.JoinPath("sources", id, "connection-check").String(),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	var check DagNodeConnectionCheck
	if err := c.do(req, &check); err != nil {
		return nil, fmt.Errorf("failed to create source connection check: %w", err)
	}

	return &check, nil
}

// GetSourceConnectionCheck retrieves the status of a connection check for a source connector by its ID.
// It returns a DagNodeConnectionCheck with the current status and reason if any.
func (c *Client) GetSourceConnectionCheck(ctx context.Context, id string) (*DagNodeConnectionCheck, error) {
	req, err := http.NewRequestWithContext(ctx,
		http.MethodGet,
		c.endpoint.JoinPath("sources", id, "connection-check").String(),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	var check DagNodeConnectionCheck
	if err := c.do(req, &check); err != nil {
		return nil, fmt.Errorf("failed to get source connection check: %w", err)
	}

	return &check, nil
}
