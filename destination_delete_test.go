package unstructured

import (
	"net/http"
	"strconv"
	"testing"
)

func TestDeleteDestination(t *testing.T) {
	t.Parallel()

	client, mux := testclient(t)

	id := "b25d4161-77a0-4e08-b65e-86f398ce15ad"

	mux.DeleteDestination = func(w http.ResponseWriter, r *http.Request) {
		if val := r.PathValue("id"); val != id {
			http.Error(w, "destination ID "+val+" not found", http.StatusNotFound)
			return
		}

		response := []byte(`{"detail": "Destination with id ` + id + ` successfully deleted."}`)

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", strconv.Itoa(len(response)))
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}

	err := client.DeleteDestination(t.Context(), "b25d4161-77a0-4e08-b65e-86f398ce15ad")
	if err != nil {
		t.Fatalf("failed to delete destination: %v", err)
	}
}
