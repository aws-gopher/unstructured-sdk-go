package unstructured

import (
	"errors"
	"net/http"
	"testing"
)

func TestRunWorkflow(t *testing.T) {
	t.Parallel()

	client, mux := testclient(t)

	id := "16b80fee-64dc-472d-8f26-1d7729b6423d"

	mux.RunWorkflow = func(w http.ResponseWriter, r *http.Request) {
		if val := r.PathValue("id"); val != id {
			http.Error(w, "workflow ID "+val+" not found", http.StatusNotFound)
			return
		}

		response := []byte(`{` +
			`  "created_at": "2025-06-22T11:37:21.648Z",` +
			`  "id": "fcdc4994-eea5-425c-91fa-e03f2bd8030d",` +
			`  "status": "IN_PROGRESS",` +
			`  "runtime": null,` +
			`  "workflow_id": "` + id + `",` +
			`  "workflow_name": "test_workflow"` +
			`}`)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		w.Write(response)
	}

	job, err := client.RunWorkflow(testContext(t), &RunWorkflowRequest{ID: id})
	if err != nil {
		t.Fatalf("failed to run workflow: %v", err)
	}

	if err := errors.Join(
		eq("new_job.id", job.ID, "fcdc4994-eea5-425c-91fa-e03f2bd8030d"),
		eq("new_job.workflow_id", job.WorkflowID, id),
		eq("new_job.workflow_name", job.WorkflowName, "test_workflow"),
		eq("new_job.status", job.Status, JobStatusInProgress),
	); err != nil {
		t.Error(err)
	}
}
