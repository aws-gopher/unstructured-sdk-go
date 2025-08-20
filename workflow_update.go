package unstructured

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// UpdateWorkflowRequest represents the request to update a workflow.
type UpdateWorkflowRequest struct {
	ID            string         `json:"-"`
	Name          *string        `json:"name,omitempty"`
	SourceID      *string        `json:"source_id,omitempty"`
	DestinationID *string        `json:"destination_id,omitempty"`
	WorkflowType  *WorkflowType  `json:"workflow_type,omitempty"`
	WorkflowNodes []WorkflowNode `json:"workflow_nodes,omitempty"`
	Schedule      *string        `json:"schedule,omitempty"`
	ReprocessAll  *bool          `json:"reprocess_all,omitempty"`
}

// UpdateWorkflow updates the configuration of an existing workflow.
// It returns the updated workflow.
func (c *Client) UpdateWorkflow(ctx context.Context, in UpdateWorkflowRequest) (*Workflow, error) {
	body, err := json.Marshal(in)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal workflow update request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx,
		http.MethodPut,
		c.endpoint.JoinPath("workflows", in.ID).String(),
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	var workflow Workflow
	if err := c.do(req, &workflow); err != nil {
		return nil, fmt.Errorf("failed to update workflow: %w", err)
	}

	return &workflow, nil
}
