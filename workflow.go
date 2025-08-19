package unstructured

import "time"

// Workflow represents a workflow, which defines a series of processing steps for data in Unstructured.io.
// A workflow connects sources, destinations, and processing nodes.
type Workflow struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	Sources       []string          `json:"sources"`
	Destinations  []string          `json:"destinations"`
	WorkflowType  *WorkflowType     `json:"workflow_type,omitempty"`
	WorkflowNodes WorkflowNodes     `json:"workflow_nodes"`
	Schedule      *WorkflowSchedule `json:"schedule,omitempty"`
	Status        WorkflowState     `json:"status"`
	CreatedAt     time.Time         `json:"created_at,omitzero"`
	UpdatedAt     time.Time         `json:"updated_at,omitzero"`
	ReprocessAll  *bool             `json:"reprocess_all,omitempty"`
}

// WorkflowType represents the type of workflow (e.g., basic, advanced, platinum, custom).
type WorkflowType string

const (
	// WorkflowTypeBasic is a basic workflow type.
	WorkflowTypeBasic WorkflowType = "basic"
	// WorkflowTypeAdvanced is an advanced workflow type.
	WorkflowTypeAdvanced WorkflowType = "advanced"
	// WorkflowTypePlatinum is a platinum workflow type.
	WorkflowTypePlatinum WorkflowType = "platinum"
	// WorkflowTypeCustom is a custom workflow type.
	WorkflowTypeCustom WorkflowType = "custom"
)

// WorkflowState represents the state of a workflow (active or inactive).
type WorkflowState string

const (
	// WorkflowStateActive indicates the workflow is active.
	WorkflowStateActive WorkflowState = "active"
	// WorkflowStateInactive indicates the workflow is inactive.
	WorkflowStateInactive WorkflowState = "inactive"
)

// WorkflowSchedule represents a workflow schedule, which can include cron tab entries.
type WorkflowSchedule struct {
	CronTabEntries []CronTabEntry `json:"crontab_entries"`
}

// CronTabEntry represents a cron tab entry for scheduling workflows.
type CronTabEntry struct {
	CronExpression string `json:"cron_expression"`
}

// SortDirection represents the sort direction for listing workflows.
type SortDirection string

const (
	// SortDirectionAsc sorts results in ascending order.
	SortDirectionAsc SortDirection = "asc"
	// SortDirectionDesc sorts results in descending order.
	SortDirectionDesc SortDirection = "desc"
)
