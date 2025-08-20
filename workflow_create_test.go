package unstructured

import (
	"errors"
	"net/http"
	"testing"
	"time"
)

func TestCreateWorkflow(t *testing.T) {
	t.Parallel()

	client, mux := testclient(t)

	mux.CreateWorkflow = func(w http.ResponseWriter, _ *http.Request) {
		response := []byte(`{` +
			`  "created_at": "2025-06-22T11:37:21.648Z",` +
			`  "destinations": ["aeebecc7-9d8e-4625-bf1d-815c2f084869"],` +
			`  "id": "16b80fee-64dc-472d-8f26-1d7729b6423d",` +
			`  "name": "test_workflow",` +
			`  "schedule": {"crontab_entries": [{"cron_expression": "0 0 * * 0"}]},` +
			`  "sources": ["f1f7b1b2-8e4b-4a2b-8f1d-3e3c7c9e5a3c"],` +
			`  "workflow_nodes": [],` +
			`  "status": "active",` +
			`  "workflow_type": "advanced"` +
			`}`)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}

	workflow, err := client.CreateWorkflow(t.Context(), &CreateWorkflowRequest{
		Name:          "test_workflow",
		WorkflowType:  WorkflowTypeAdvanced,
		Schedule:      String("weekly"),
		SourceID:      String("f1f7b1b2-8e4b-4a2b-8f1d-3e3c7c9e5a3c"),
		DestinationID: String("aeebecc7-9d8e-4625-bf1d-815c2f084869"),
	})
	if err != nil {
		t.Fatalf("failed to create workflow: %v", err)
	}

	if err := errors.Join(
		eq("new_workflow.id", workflow.ID, "16b80fee-64dc-472d-8f26-1d7729b6423d"),
		eq("new_workflow.name", workflow.Name, "test_workflow"),
		eq("new_workflow.status", workflow.Status, WorkflowStateActive),
		eq("new_workflow.workflow_type", ToVal(workflow.WorkflowType), WorkflowTypeAdvanced),
		equal("new_workflow.created_at", workflow.CreatedAt, time.Date(2025, 6, 22, 11, 37, 21, 648000000, time.UTC)),
		eqs("new_workflow.sources", workflow.Sources, []string{"f1f7b1b2-8e4b-4a2b-8f1d-3e3c7c9e5a3c"}),
		eqs("new_workflow.destinations", workflow.Destinations, []string{"aeebecc7-9d8e-4625-bf1d-815c2f084869"}),
	); err != nil {
		t.Error(err)
	}
}
