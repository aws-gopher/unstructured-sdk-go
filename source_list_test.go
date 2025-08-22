package unstructured

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestListSources(t *testing.T) {
	t.Parallel()

	client, mux := testclient(t)

	mux.ListSources = func(w http.ResponseWriter, _ *http.Request) {
		response := []byte(`[` +
			`  {` +
			`    "config": {` +
			`      "client_id": "foo",` +
			`      "tenant": "foo",` +
			`      "authority_url": "foo",` +
			`      "user_pname": "foo",` +
			`      "client_cred": "foo",` +
			`      "recursive": false,` +
			`      "path": "foo"` +
			`    },` +
			`    "created_at": "2023-09-15T01:06:53.146Z",` +
			`    "id": "a15d4161-77a0-4e08-b65e-86f398ce15ad",` +
			`    "name": "test_source_name",` +
			`    "type": "onedrive"` +
			`  }` +
			`]`)

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", strconv.Itoa(len(response)))
		w.Write(response)
	}

	sources, err := client.ListSources(t.Context(), "")
	if err != nil {
		t.Fatalf("failed to list sources: %v", err)
	}

	if len(sources) != 1 {
		t.Fatalf("expected 1 source, got %d", len(sources))
	}

	source := sources[0]
	if err := errors.Join(
		eq("source.id", source.ID, "a15d4161-77a0-4e08-b65e-86f398ce15ad"),
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

func TestListSourcesEmpty(t *testing.T) {
	t.Parallel()

	client, mux := testclient(t)

	mux.ListSources = func(w http.ResponseWriter, _ *http.Request) {
		response := []byte(`[]`)

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", strconv.Itoa(len(response)))
		w.Write(response)
	}

	sources, err := client.ListSources(t.Context(), "")
	if err != nil {
		t.Fatalf("failed to list sources: %v", err)
	}

	if len(sources) != 0 {
		t.Fatalf("expected 0 sources, got %d", len(sources))
	}
}

func TestListSourcesErrorCode(t *testing.T) {
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

			mux.ListSources = func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(code)
			}

			_, err := client.ListSources(t.Context(), "")
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
