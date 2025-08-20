package unstructured

import (
	"errors"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func TestUpdateWorkflow(t *testing.T) {
	t.Parallel()

	client, mux := testclient(t)

	id := "16b80fee-64dc-472d-8f26-1d7729b6423d"

	mux.UpdateWorkflow = func(w http.ResponseWriter, _ *http.Request) {
		response := []byte(`{` +
			`  "id": "` + id + `",` +
			`  "name": "test_workflow",` +
			`  "status": "active",` +
			`  "workflow_type": "advanced",` +
			`  "created_at": "2025-06-22T11:37:21.648Z",` +
			`  "sources": ["f1f7b1b2-8e4b-4a2b-8f1d-3e3c7c9e5a3c"],` +
			`  "destinations": ["aeebecc7-9d8e-4625-bf1d-815c2f084869"],` +
			`  "schedule": {"crontab_entries": [{"cron_expression": "0 0 * * 0"}]},` +
			`  "workflow_nodes": []` +
			`}`)

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", strconv.Itoa(len(response)))
		w.Write(response)
	}

	updated, err := client.UpdateWorkflow(t.Context(), UpdateWorkflowRequest{
		ID:            id,
		Name:          String("test_workflow"),
		WorkflowType:  Ptr(WorkflowTypeAdvanced),
		Schedule:      String("weekly"),
		SourceID:      String("f1f7b1b2-8e4b-4a2b-8f1d-3e3c7c9e5a3c"),
		DestinationID: String("aeebecc7-9d8e-4625-bf1d-815c2f084869"),
	})
	if err != nil {
		t.Fatalf("failed to update workflow: %v", err)
	}

	if err := errors.Join(
		eq("updated_workflow.id", updated.ID, id),
		eq("updated_workflow.name", updated.Name, "test_workflow"),
		eq("updated_workflow.status", updated.Status, WorkflowStateActive),
		eq("updated_workflow.workflow_type", ToVal(updated.WorkflowType), WorkflowTypeAdvanced),
		equal("workflow.created_at", updated.CreatedAt, time.Date(2025, 6, 22, 11, 37, 21, 648000000, time.UTC)),
		eqs("updated_workflow.sources", updated.Sources, []string{"f1f7b1b2-8e4b-4a2b-8f1d-3e3c7c9e5a3c"}),
		eqs("updated_workflow.destinations", updated.Destinations, []string{"aeebecc7-9d8e-4625-bf1d-815c2f084869"}),
		eqs("updated_workflow.schedule.crontab_entries", updated.Schedule.CronTabEntries, []CronTabEntry{{CronExpression: "0 0 * * 0"}}),
	); err != nil {
		t.Error(err)
	}
}
