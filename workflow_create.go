package unstructured

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// CreateWorkflowRequest represents the request to create a workflow
type CreateWorkflowRequest struct {
	Name          string         `json:"name"`
	SourceID      *string        `json:"source_id,omitempty"`
	DestinationID *string        `json:"destination_id,omitempty"`
	WorkflowType  WorkflowType   `json:"workflow_type"`
	WorkflowNodes []WorkflowNode `json:"workflow_nodes,omitempty"`
	Schedule      *string        `json:"schedule,omitempty"`
	ReprocessAll  *bool          `json:"reprocess_all,omitempty"`
}

// CreateWorkflow creates a new workflow
func (c *Client) CreateWorkflow(ctx context.Context, in CreateWorkflowRequest) (*Workflow, error) {
	body, err := json.Marshal(in)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal workflow request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx,
		http.MethodPost,
		c.endpoint.JoinPath("/workflows").String(),
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	var workflow Workflow
	if err := c.do(req, &workflow); err != nil {
		return nil, fmt.Errorf("failed to create workflow: %w", err)
	}

	return &workflow, nil
}
