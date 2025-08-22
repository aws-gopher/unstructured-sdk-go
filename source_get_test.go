package unstructured

import (
	"errors"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func TestGetSource(t *testing.T) {
	t.Parallel()

	client, mux := testclient(t)

	id := "a15d4161-77a0-4e08-b65e-86f398ce15ad"
	mux.GetSource = func(w http.ResponseWriter, r *http.Request) {
		if val := r.PathValue("id"); val != id {
			http.Error(w, "source ID "+val+" not found", http.StatusNotFound)
			return
		}

		response := []byte(`{` +
			`  "config": {` +
			`    "client_id": "foo",` +
			`    "tenant": "foo",` +
			`    "authority_url": "foo",` +
			`    "user_pname": "foo",` +
			`    "client_cred": "foo",` +
			`    "recursive": false,` +
			`    "path": "foo"` +
			`  },` +
			`  "created_at": "2023-09-15T01:06:53.146Z",` +
			`  "id": "` + id + `",` +
			`  "name": "test_source_name",` +
			`  "type": "onedrive"` +
			`}`)

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", strconv.Itoa(len(response)))
		w.Write(response)
	}

	source, err := client.GetSource(t.Context(), id)
	if err != nil {
		t.Fatalf("failed to get source: %v", err)
	}

	if err := errors.Join(
		eq("source.id", source.ID, id),
		eq("source.name", source.Name, "test_source_name"),
		equal("source.created_at", source.CreatedAt, time.Date(2023, 9, 15, 1, 6, 53, 146000000, time.UTC)),
	); err != nil {
		t.Error(err)
	}

	cfg, ok := source.Config.(*OneDriveConnectorConfig)
	if !ok {
		t.Errorf("expected source config to be %T, got %T", cfg, source.Config)
	}
}

func TestGetSourceNotFound(t *testing.T) {
	t.Parallel()

	client, mux := testclient(t)

	id := "a15d4161-77a0-4e08-b65e-86f398ce15ad"
	mux.GetSource = func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "source ID "+r.PathValue("id")+" not found", http.StatusNotFound)
	}

	_, err := client.GetSource(t.Context(), id)
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
