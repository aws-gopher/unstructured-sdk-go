package unstructured

import (
	"net/http"
	"strconv"
	"testing"
)

func TestDeleteSource(t *testing.T) {
	t.Parallel()

	client, mux := testclient(t)

	id := "a15d4161-77a0-4e08-b65e-86f398ce15ad"

	mux.DeleteSource = func(w http.ResponseWriter, r *http.Request) {
		if val := r.PathValue("id"); val != id {
			http.Error(w, "source ID "+val+" not found", http.StatusNotFound)
			return
		}

		response := []byte(`{"detail": "Source with id ` + id + ` successfully deleted."}`)

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", strconv.Itoa(len(response)))
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}

	err := client.DeleteSource(t.Context(), id)
	if err != nil {
		t.Fatalf("failed to delete source: %v", err)
	}
}
