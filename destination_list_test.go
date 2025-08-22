package unstructured

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestListDestinations(t *testing.T) {
	t.Parallel()

	client, mux := testclient(t)

	mux.ListDestinations = func(w http.ResponseWriter, _ *http.Request) {
		response := []byte(`[` +
			`  {` +
			`    "config": {` +
			`      "remote_url": "s3://mock-s3-connector",` +
			`      "anonymous": false,` +
			`      "key": "**********",` +
			`      "secret": "**********",` +
			`      "token": null,` +
			`      "endpoint_url": null` +
			`    },` +
			`    "created_at": "2025-08-22T08:47:29.802Z",` +
			`    "id": "0c363dec-3c70-45ee-8041-481044a6e1cc",` +
			`    "name": "test_destination_name",` +
			`    "type": "s3"` +
			`  }` +
			`]`)

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", strconv.Itoa(len(response)))
		w.Write(response)
	}

	destinations, err := client.ListDestinations(t.Context(), "")
	if err != nil {
		t.Fatalf("failed to list destinations: %v", err)
	}

	if len(destinations) != 1 {
		t.Fatalf("expected 1 destination, got %d", len(destinations))
	}

	destination := destinations[0]
	if err := errors.Join(
		eq("destination.id", destination.ID, "0c363dec-3c70-45ee-8041-481044a6e1cc"),
		eq("destination.name", destination.Name, "test_destination_name"),
		eq("destination.type", destination.Type, "s3"),
		equal("destination.created_at", destination.CreatedAt, time.Date(2025, 8, 22, 8, 47, 29, 802000000, time.UTC)),
	); err != nil {
		t.Error(err)
	}

	cfg, ok := destination.Config.(*S3ConnectorConfig)
	if !ok {
		t.Errorf("expected destination config to be %T, got %T", cfg, destination.Config)
	}
}

func TestListDestinationsEmpty(t *testing.T) {
	t.Parallel()

	client, mux := testclient(t)

	mux.ListDestinations = func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", "2")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
	}

	destinations, err := client.ListDestinations(t.Context(), "")
	if err != nil {
		t.Fatalf("failed to list destinations: %v", err)
	}

	if len(destinations) != 0 {
		t.Fatalf("expected 0 destinations, got %d", len(destinations))
	}
}

func TestListDestinationsErrorCode(t *testing.T) {
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

			mux.ListDestinations = func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(code)
			}

			_, err := client.ListDestinations(t.Context(), "")
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
