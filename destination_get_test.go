package unstructured

import (
	"errors"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func TestGetDestination(t *testing.T) {
	t.Parallel()

	client, mux := testclient(t)

	id := "0c363dec-3c70-45ee-8041-481044a6e1cc"
	mux.GetDestination = func(w http.ResponseWriter, r *http.Request) {
		if val := r.PathValue("id"); val != id {
			http.Error(w, "destination ID "+val+" not found", http.StatusNotFound)
			return
		}

		response := []byte(`{` +
			`  "config": {` +
			`    "remote_url": "s3://mock-s3-connector",` +
			`    "anonymous": false,` +
			`    "key": "**********",` +
			`    "secret": "**********",` +
			`    "token": null,` +
			`    "endpoint_url": null` +
			`  },` +
			`  "created_at": "2025-08-22T08:47:29.802Z",` +
			`  "id": "` + id + `",` +
			`  "name": "test_destination_name",` +
			`  "type": "s3"` +
			`}`)

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", strconv.Itoa(len(response)))
		w.Write(response)
	}

	destination, err := client.GetDestination(testContext(t), id)
	if err != nil {
		t.Fatalf("failed to get destination: %v", err)
	}

	if err := errors.Join(
		eq("destination.id", destination.ID, id),
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

func TestGetDestinationNotFound(t *testing.T) {
	t.Parallel()

	client, mux := testclient(t)

	id := "0c363dec-3c70-45ee-8041-481044a6e1cc"
	mux.GetDestination = func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "destination ID "+r.PathValue("id")+" not found", http.StatusNotFound)
	}

	_, err := client.GetDestination(testContext(t), id)
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
