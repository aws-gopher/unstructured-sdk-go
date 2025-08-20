package unstructured

import (
	"net/http"
	"strconv"
	"testing"
)

func TestCancelJob(t *testing.T) {
	t.Parallel()

	client, mux := testclient(t)

	id := "fcdc4994-eea5-425c-91fa-e03f2bd8030d"

	mux.CancelJob = func(w http.ResponseWriter, r *http.Request) {
		if val := r.PathValue("id"); val != id {
			http.Error(w, "job ID "+val+" not found", http.StatusNotFound)
			return
		}

		response := []byte(`{` +
			`  "id": "` + id + `",` +
			`  "status": "cancelled",` +
			`  "message": "Job successfully cancelled."` +
			`}`)

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", strconv.Itoa(len(response)))
		w.Write(response)
	}

	err := client.CancelJob(t.Context(), id)
	if err != nil {
		t.Fatalf("failed to cancel job: %v", err)
	}
}
