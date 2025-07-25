package unstructured

import "time"

// DagNodeConnectionCheck represents a connection check result for a DAG node (source or destination connector).
type DagNodeConnectionCheck struct {
	ID         string                `json:"id"`
	Status     ConnectionCheckStatus `json:"status"`
	Reason     *string               `json:"reason,omitempty"`
	CreatedAt  time.Time             `json:"created_at,omitempty"`
	ReportedAt *string               `json:"reported_at,omitempty"`
}

// ConnectionCheckStatus represents the status of a connection check (scheduled, success, or failure).
type ConnectionCheckStatus string

const (
	// ConnectionCheckStatusScheduled indicates the connection check is scheduled.
	ConnectionCheckStatusScheduled ConnectionCheckStatus = "SCHEDULED"
	// ConnectionCheckStatusSuccess indicates the connection check succeeded.
	ConnectionCheckStatusSuccess ConnectionCheckStatus = "SUCCESS"
	// ConnectionCheckStatusFailure indicates the connection check failed.
	ConnectionCheckStatusFailure ConnectionCheckStatus = "FAILURE"
)
