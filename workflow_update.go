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
	Name          *string
	SourceID      *string
	DestinationID *string
	WorkflowType  *WorkflowType
	WorkflowNodes []WorkflowNode
	Schedule      *string
	ReprocessAll  *bool
}

// UpdateWorkflow updates the configuration of an existing workflow.
// It returns the updated workflow.
func (c *Client) UpdateWorkflow(ctx context.Context, workflowID string, in UpdateWorkflowRequest) (*Workflow, error) {
	body, err := json.Marshal(in)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal workflow update request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx,
		http.MethodPut,
		c.endpoint.JoinPath("/workflows", workflowID).String(),
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
