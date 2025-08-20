package unstructured

import (
	"errors"
	"net/http"
	"testing"
	"time"
)

func TestCreateSource(t *testing.T) {
	t.Parallel()

	client, mux := testclient(t)

	mux.CreateSource = func(w http.ResponseWriter, _ *http.Request) {
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
			`  "id": "a15d4161-77a0-4e08-b65e-86f398ce15ad",` +
			`  "name": "test_source_name",` +
			`  "type": "onedrive"` +
			`}`)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}

	source, err := client.CreateSource(t.Context(), CreateSourceRequest{
		Name: "test_source_name",
		Config: &OneDriveSourceConnectorConfigInput{
			ClientID:     "foo",
			Tenant:       "foo",
			AuthorityURL: "foo",
			UserPName:    "foo",
			ClientCred:   "foo",
			Path:         "foo",
		},
	})
	if err != nil {
		t.Fatalf("failed to create source: %v", err)
	}

	if err := errors.Join(
		eq("source.id", source.ID, "a15d4161-77a0-4e08-b65e-86f398ce15ad"),
		eq("source.name", source.Name, "test_source_name"),
		equal("source.created_at", source.CreatedAt, time.Date(2023, 9, 15, 1, 6, 53, 146000000, time.UTC)),
	); err != nil {
		t.Error(err)
	}

	cfg, ok := source.Config.(*OneDriveSourceConnectorConfig)
	if !ok {
		t.Errorf("expected source config to be %T, got %T", cfg, source.Config)
	}
}
