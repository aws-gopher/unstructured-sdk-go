package unstructured

import (
	"errors"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func TestListJobs(t *testing.T) {
	t.Parallel()

	client, mux := testclient(t)

	mux.ListJobs = func(w http.ResponseWriter, _ *http.Request) {
		response := []byte(`[` +
			`  {` +
			`     "created_at": "2025-06-22T11:37:21.648Z",` +
			`     "id": "fcdc4994-eea5-425c-91fa-e03f2bd8030d",` +
			`     "status": "IN_PROGRESS",` +
			`     "runtime": null,` +
			`     "workflow_id": "16b80fee-64dc-472d-8f26-1d7729b6423d",` +
			`     "workflow_name": "test_workflow"` +
			`  }` +
			`]`)

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", strconv.Itoa(len(response)))
		w.Write(response)
	}

	jobs, err := client.ListJobs(t.Context(), &ListJobsRequest{})
	if err != nil {
		t.Fatalf("failed to list jobs: %v", err)
	}

	if len(jobs) != 1 {
		t.Fatalf("expected 1 job, got %d", len(jobs))
	}

	job := jobs[0]

	if err := errors.Join(
		eq("job.id", job.ID, "fcdc4994-eea5-425c-91fa-e03f2bd8030d"),
		eq("job.workflow_id", job.WorkflowID, "16b80fee-64dc-472d-8f26-1d7729b6423d"),
		eq("job.workflow_name", job.WorkflowName, "test_workflow"),
		eq("job.status", job.Status, JobStatusInProgress),
		equal("job.created_at", job.CreatedAt, time.Date(2025, 6, 22, 11, 37, 21, 648000000, time.UTC)),
	); err != nil {
		t.Error(err)
	}
}
