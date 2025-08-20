package unstructured

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestGetJob(t *testing.T) {
	t.Parallel()

	client, mux := testclient(t)

	id := "fcdc4994-eea5-425c-91fa-e03f2bd8030d"
	mux.GetJob = func(w http.ResponseWriter, r *http.Request) {
		if val := r.PathValue("id"); val != id {
			http.Error(w, "job ID "+val+" not found", http.StatusNotFound)
			return
		}

		response := []byte(`{` +
			`  "created_at": "2025-06-22T11:37:21.648Z",` +
			`  "id": "` + id + `",` +
			`  "status": "SCHEDULED",` +
			`  "runtime": null,` +
			`  "workflow_id": "16b80fee-64dc-472d-8f26-1d7729b6423d",` +
			`  "workflow_name": "test_workflow"` +
			`}`)

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", strconv.Itoa(len(response)))
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}

	job, err := client.GetJob(t.Context(), id)
	if err != nil {
		t.Fatalf("failed to get job: %v", err)
	}

	if err := errors.Join(
		eq("job.id", job.ID, id),
		eq("job.workflow_id", job.WorkflowID, "16b80fee-64dc-472d-8f26-1d7729b6423d"),
		eq("job.workflow_name", job.WorkflowName, "test_workflow"),
		eq("job.status", job.Status, JobStatusScheduled),
		equal("job.created_at", job.CreatedAt, time.Date(2025, 6, 22, 11, 37, 21, 648000000, time.UTC)),
	); err != nil {
		t.Error(err)
	}
}

func TestGetJobNotFound(t *testing.T) {
	t.Parallel()

	client, mux := testclient(t)

	id := "fcdc4994-eea5-425c-91fa-e03f2bd8030d"
	mux.GetJob = func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "job ID "+r.PathValue("id")+" not found", http.StatusNotFound)
	}

	_, err := client.GetJob(t.Context(), id)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	var apierr *APIError
	if !errors.As(err, &apierr) {
		t.Fatalf("expected error to be an %T, got %T", apierr, err)
	}

	if apierr.Code != http.StatusNotFound {
		t.Fatalf("expected error code to be %d, got %d", http.StatusNotFound, apierr.Code)
	}
}

func TestGetJobError(t *testing.T) {
	t.Parallel()

	id := "fcdc4994-eea5-425c-91fa-e03f2bd8030d"

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

			mux.GetJob = func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(code)
			}

			_, err := client.GetJob(t.Context(), id)
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
