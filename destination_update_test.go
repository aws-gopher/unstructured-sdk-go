package unstructured

import (
	"errors"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func TestUpdateDestination(t *testing.T) {
	t.Parallel()

	client, mux := testclient(t)

	id := "b25d4161-77a0-4e08-b65e-86f398ce15ad"

	mux.UpdateDestination = func(w http.ResponseWriter, r *http.Request) {
		if val := r.PathValue("id"); val != id {
			http.Error(w, "destination ID "+val+" not found", http.StatusNotFound)
			return
		}

		response := []byte(`{` +
			`  "config": {` +
			`    "remote_url": "s3://mock-s3-connector",` +
			`    "key": "blah",` +
			`    "secret": "blah",` +
			`    "anonymous": false` +
			`  },` +
			`  "created_at": "2023-09-15T01:06:53.146Z",` +
			`  "id": "` + id + `",` +
			`  "name": "test_destination_name",` +
			`  "type": "s3"` +
			`}`)

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", strconv.Itoa(len(response)))
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}

	updated, err := client.UpdateDestination(testContext(t), UpdateDestinationRequest{
		ID: id,
		Config: &S3ConnectorConfig{
			RemoteURL: "s3://mock-s3-connector",
			Key:       String("blah"),
			Secret:    String("blah"),
		},
	})
	if err != nil {
		t.Fatalf("failed to update destination: %v", err)
	}

	if err := errors.Join(
		eq("updated_destination.id", updated.ID, id),
		eq("updated_destination.name", updated.Name, "test_destination_name"),
		equal("updated_destination.created_at", updated.CreatedAt, time.Date(2023, 9, 15, 1, 6, 53, 146000000, time.UTC)),
	); err != nil {
		t.Error(err)
	}

	cfg, ok := updated.Config.(*S3ConnectorConfig)
	if !ok {
		t.Errorf("expected destination config to be %T, got %T", cfg, updated.Config)
	}
}
