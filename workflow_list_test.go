package unstructured

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestListWorkflows(t *testing.T) {
	t.Parallel()

	client, mux := testclient(t)

	mux.ListWorkflows = func(w http.ResponseWriter, _ *http.Request) {
		response := []byte(`[` +
			`  {` +
			`    "created_at": "2025-06-22T11:37:21.648Z",` +
			`    "destinations": ["aeebecc7-9d8e-4625-bf1d-815c2f084869"],` +
			`    "id": "16b80fee-64dc-472d-8f26-1d7729b6423d",` +
			`    "name": "test_workflow",` +
			`    "schedule": {"crontab_entries": [{"cron_expression": "0 0 * * 0"}]},` +
			`    "sources": ["f1f7b1b2-8e4b-4a2b-8f1d-3e3c7c9e5a3c"],` +
			`    "workflow_nodes": [],` +
			`    "status": "active",` +
			`    "workflow_type": "advanced"` +
			`  }` +
			`]`)

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", strconv.Itoa(len(response)))
		w.Write(response)
	}

	workflows, err := client.ListWorkflows(testContext(t), &ListWorkflowsRequest{
		SortBy: String("id"),
	})
	if err != nil {
		t.Fatalf("failed to list workflows: %v", err)
	}

	if len(workflows) != 1 {
		t.Fatalf("expected 1 workflow, got %d", len(workflows))
	}

	workflow := workflows[0]
	if err := errors.Join(
		eq("workflow.id", workflow.ID, "16b80fee-64dc-472d-8f26-1d7729b6423d"),
		eq("workflow.name", workflow.Name, "test_workflow"),
		eq("workflow.workflow_type", ToVal(workflow.WorkflowType), WorkflowTypeAdvanced),
		eq("workflow.status", workflow.Status, WorkflowStateActive),
		equal("workflow.created_at", workflow.CreatedAt, time.Date(2025, 6, 22, 11, 37, 21, 648000000, time.UTC)),
		eqs("workflow.schedule.crontab_entries", ToVal(workflow.Schedule).CronTabEntries, []CronTabEntry{{CronExpression: "0 0 * * 0"}}),
		eqs("workflow.sources", workflow.Sources, []string{"f1f7b1b2-8e4b-4a2b-8f1d-3e3c7c9e5a3c"}),
		eqs("workflow.destinations", workflow.Destinations, []string{"aeebecc7-9d8e-4625-bf1d-815c2f084869"}),
	); err != nil {
		t.Error(err)
	}
}

func TestListWorkflowsEmpty(t *testing.T) {
	t.Parallel()

	client, mux := testclient(t)

	mux.ListWorkflows = func(w http.ResponseWriter, _ *http.Request) {
		response := []byte(`[]`)

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", strconv.Itoa(len(response)))
		w.Write(response)
	}

	workflows, err := client.ListWorkflows(testContext(t), &ListWorkflowsRequest{
		SortBy: String("id"),
	})
	if err != nil {
		t.Fatalf("failed to list workflows: %v", err)
	}

	if len(workflows) != 0 {
		t.Fatalf("expected 0 workflows, got %d", len(workflows))
	}
}

func TestListWorkflowsErrorCode(t *testing.T) {
	t.Parallel()

	for _, code := range []int{
		http.StatusBadRequest,          // 400
		http.StatusUnauthorized,        // 401
		http.StatusForbidden,           // 403
		http.StatusNotFound,            // 404
		http.StatusInternalServerError, // 500
		http.StatusBadGateway,          // 502
		http.StatusServiceUnavailable,  // 503
		http.StatusGatewayTimeout,      // 504
	} {
		t.Run(strconv.Itoa(code), func(t *testing.T) {
			t.Parallel()

			client, mux := testclient(t)

			mux.ListWorkflows = func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(code)
			}

			_, err := client.ListWorkflows(testContext(t), &ListWorkflowsRequest{
				SortBy: String("id"),
			})
			if err == nil {
				t.Fatalf("expected error, got nil")
			}

			if !strings.Contains(err.Error(), "API error occurred") {
				t.Fatalf("expected error to contain 'API error occurred', got %v", err)
			}

			var apierr *APIError
			if !errors.As(err, &apierr) {
				t.Fatalf("expected error to be an %T, got %T", apierr, err)
			}

			if apierr.Code != code {
				t.Fatalf("expected error code to be %d, got %d", code, apierr.Code)
			}
		})
	}
}
