package unstructured

import (
	"context"
	"fmt"
	"net/http"
)

// ListJobsRequest represents the request to list jobs with optional filters.
type ListJobsRequest struct {
	WorkflowID *string
	Status     *JobStatus
}

// ListJobs retrieves a list of jobs with optional filtering.
func (c *Client) ListJobs(ctx context.Context, in *ListJobsRequest) ([]Job, error) {
	req, err := http.NewRequestWithContext(ctx,
		http.MethodGet,
		c.endpoint.JoinPath("jobs", "").String(),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	if in != nil {
		q := req.URL.Query()

		if in.WorkflowID != nil {
			q.Add("workflow_id", *in.WorkflowID)
		}

		if in.Status != nil {
			q.Add("status", string(*in.Status))
		}

		req.URL.RawQuery = q.Encode()
	}

	var jobs []Job
	if err := c.do(req, &jobs); err != nil {
		return nil, fmt.Errorf("failed to list jobs: %w", err)
	}

	return jobs, nil
}
