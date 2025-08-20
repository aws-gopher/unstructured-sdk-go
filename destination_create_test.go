package unstructured

import (
	"errors"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func TestCreateDestination(t *testing.T) {
	t.Parallel()

	client, mux := testclient(t)

	mux.CreateDestination = func(w http.ResponseWriter, _ *http.Request) {
		response := []byte(`{` +
			`  "config": {` +
			`    "remote_url": "s3://mock-s3-connector",` +
			`    "key": "blah",` +
			`    "secret": "blah",` +
			`    "anonymous": false` +
			`  },` +
			`  "created_at": "2023-09-15T01:06:53.146Z",` +
			`  "id": "b25d4161-77a0-4e08-b65e-86f398ce15ad",` +
			`  "name": "test_destination_name",` +
			`  "type": "s3"` +
			`}`)

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", strconv.Itoa(len(response)))
		w.Write(response)
	}

	destination, err := client.CreateDestination(t.Context(), CreateDestinationRequest{
		Name: "test_destination_name",

		Config: &S3DestinationConnectorConfigInput{
			RemoteURL: "s3://mock-s3-connector",
			Key:       String("blah"),
			Secret:    String("blah"),
		},
	})
	if err != nil {
		t.Fatalf("failed to create destination: %v", err)
	}

	if err := errors.Join(
		eq("destination.id", destination.ID, "b25d4161-77a0-4e08-b65e-86f398ce15ad"),
		eq("destination.name", destination.Name, "test_destination_name"),
		equal("destination.created_at", destination.CreatedAt, time.Date(2023, 9, 15, 1, 6, 53, 146000000, time.UTC)),
	); err != nil {
		t.Error(err)
	}

	cfg, ok := destination.Config.(*S3DestinationConnectorConfig)
	if !ok {
		t.Errorf("expected destination config to be %T, got %T", cfg, destination.Config)
	}
}
