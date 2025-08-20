package unstructured

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Job represents a job, which is an execution of a workflow in Unstructured.io.
type Job struct {
	ID              string             `json:"id"`
	WorkflowID      string             `json:"workflow_id"`
	WorkflowName    string             `json:"workflow_name"`
	Status          JobStatus          `json:"status"`
	CreatedAt       time.Time          `json:"created_at,omitzero"`
	Runtime         *string            `json:"runtime,omitempty"`
	InputFileIDs    []string           `json:"input_file_ids,omitempty"`
	OutputNodeFiles []NodeFileMetadata `json:"output_node_files,omitempty"`
	JobType         WorkflowJobType    `json:"job_type"`
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (j *Job) UnmarshalJSON(data []byte) error {
	type mask Job

	shadowed := struct {
		*mask
		CreatedAt string `json:"created_at,omitempty"`
	}{
		mask: (*mask)(j),
	}
	if err := json.Unmarshal(data, &shadowed); err != nil {
		return fmt.Errorf("failed to unmarshal job: %w", err)
	}

	if shadowed.CreatedAt != "" {
		t, err := time.Parse("2006-01-02T15:04:05", strings.TrimSuffix(shadowed.CreatedAt, "Z"))
		if err != nil {
			return fmt.Errorf("failed to parse job creation time: %w", err)
		}

		j.CreatedAt = t
	}

	return nil
}

// JobStatus represents the status of a job (e.g., scheduled, in progress, completed, stopped, failed).
type JobStatus string

const (
	// JobStatusScheduled indicates the job is scheduled.
	JobStatusScheduled JobStatus = "SCHEDULED"
	// JobStatusInProgress indicates the job is in progress.
	JobStatusInProgress JobStatus = "IN_PROGRESS"
	// JobStatusCompleted indicates the job is completed.
	JobStatusCompleted JobStatus = "COMPLETED"
	// JobStatusStopped indicates the job is stopped.
	JobStatusStopped JobStatus = "STOPPED"
	// JobStatusFailed indicates the job has failed.
	JobStatusFailed JobStatus = "FAILED"
)

// JobProcessingStatus represents the processing status of a job (e.g., scheduled, in progress, success, etc.).
type JobProcessingStatus string

const (
	// JobProcessingStatusScheduled indicates the job is scheduled for processing.
	JobProcessingStatusScheduled JobProcessingStatus = "SCHEDULED"
	// JobProcessingStatusInProgress indicates the job is currently being processed.
	JobProcessingStatusInProgress JobProcessingStatus = "IN_PROGRESS"
	// JobProcessingStatusSuccess indicates the job was processed successfully.
	JobProcessingStatusSuccess JobProcessingStatus = "SUCCESS"
	// JobProcessingStatusCompletedWithErrors indicates the job completed with errors.
	JobProcessingStatusCompletedWithErrors JobProcessingStatus = "COMPLETED_WITH_ERRORS"
	// JobProcessingStatusStopped indicates the job was stopped.
	JobProcessingStatusStopped JobProcessingStatus = "STOPPED"
	// JobProcessingStatusFailed indicates the job failed.
	JobProcessingStatusFailed JobProcessingStatus = "FAILED"
)

// WorkflowJobType represents the type of workflow job (ephemeral, persistent, scheduled).
type WorkflowJobType string

const (
	// WorkflowJobTypeEphemeral is an ephemeral job type.
	WorkflowJobTypeEphemeral WorkflowJobType = "ephemeral"
	// WorkflowJobTypePersistent is a persistent job type.
	WorkflowJobTypePersistent WorkflowJobType = "persistent"
	// WorkflowJobTypeScheduled is a scheduled job type.
	WorkflowJobTypeScheduled WorkflowJobType = "scheduled"
)

// NodeFileMetadata represents metadata for a node file in a job.
type NodeFileMetadata struct {
	NodeID string `json:"node_id"`
	FileID string `json:"file_id"`
}

// JobDetails represents detailed information about a job, including processing status and node stats.
type JobDetails struct {
	ID               string              `json:"id"`
	ProcessingStatus JobProcessingStatus `json:"processing_status"`
	NodeStats        []JobNodeDetails    `json:"node_stats"`
	Message          *string             `json:"message,omitempty"`
}

// JobNodeDetails represents details about a job node, including status counts.
type JobNodeDetails struct {
	NodeName    *string `json:"node_name,omitempty"`
	NodeType    *string `json:"node_type,omitempty"`
	NodeSubtype *string `json:"node_subtype,omitempty"`
	Ready       int     `json:"ready"`
	InProgress  int     `json:"in_progress"`
	Success     int     `json:"success"`
	Failure     int     `json:"failure"`
}

// JobFailedFiles represents failed files for a job.
type JobFailedFiles struct {
	FailedFiles []FailedFile `json:"failed_files"`
}

// FailedFile represents a failed file in a job, including the document and error message.
type FailedFile struct {
	Document string `json:"document"`
	Error    string `json:"error"`
}
