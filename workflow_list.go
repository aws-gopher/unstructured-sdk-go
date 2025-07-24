package unstructured

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// ListWorkflowsRequest represents the request to list workflows with optional filters.
type ListWorkflowsRequest struct {
	DagNodeConfigurationID   *string
	SourceID                 *string
	DestinationID            *string
	Status                   *WorkflowState
	Page                     *int
	PageSize                 *int
	CreatedSince             *time.Time
	CreatedBefore            *time.Time
	Name                     *string
	SortBy                   *string
	SortDirection            *SortDirection
	ShowOnlySoftDeleted      *bool
	ShowRecommenderWorkflows *bool
}

// ListWorkflows retrieves a list of workflows with optional filtering and pagination.
func (c *Client) ListWorkflows(ctx context.Context, in *ListWorkflowsRequest) ([]Workflow, error) {
	req, err := http.NewRequestWithContext(ctx,
		http.MethodGet,
		c.endpoint.JoinPath("/workflows").String(),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	if in != nil {
		req.URL.RawQuery = buildWorkflowListQuery(in).Encode()
	}

	var workflows []Workflow
	if err := c.do(req, &workflows); err != nil {
		return nil, fmt.Errorf("failed to list workflows: %w", err)
	}

	return workflows, nil
}

// buildWorkflowListQuery builds the query parameters for the workflow list request.
func buildWorkflowListQuery(in *ListWorkflowsRequest) url.Values {
	q := make(url.Values)

	if in.DagNodeConfigurationID != nil {
		q.Add("dag_node_configuration_id", *in.DagNodeConfigurationID)
	}

	if in.SourceID != nil {
		q.Add("source_id", *in.SourceID)
	}

	if in.DestinationID != nil {
		q.Add("destination_id", *in.DestinationID)
	}

	if in.Status != nil {
		q.Add("status", string(*in.Status))
	}

	if in.Page != nil {
		q.Add("page", strconv.Itoa(*in.Page))
	}

	if in.PageSize != nil {
		q.Add("page_size", strconv.Itoa(*in.PageSize))
	}

	if in.CreatedSince != nil {
		q.Add("created_since", in.CreatedSince.Format(time.RFC3339))
	}

	if in.CreatedBefore != nil {
		q.Add("created_before", in.CreatedBefore.Format(time.RFC3339))
	}

	if in.Name != nil {
		q.Add("name", *in.Name)
	}

	if in.SortBy != nil {
		q.Add("sort_by", *in.SortBy)
	}

	if in.SortDirection != nil {
		q.Add("sort_direction", string(*in.SortDirection))
	}

	if in.ShowOnlySoftDeleted != nil {
		q.Add("show_only_soft_deleted", strconv.FormatBool(*in.ShowOnlySoftDeleted))
	}

	if in.ShowRecommenderWorkflows != nil {
		q.Add("show_recommender_workflows", strconv.FormatBool(*in.ShowRecommenderWorkflows))
	}

	return q
}
