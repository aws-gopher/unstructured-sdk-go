package unstructured

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws-gopher/unstructured-sdk-go/test"
)

func testclient(t *testing.T) (*Client, *test.Mux) {
	mux := test.NewMux()

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		val := r.Header.Get("unstructured-api-key")
		if val == "" {
			http.Error(w, "Unauthorized: missing header", http.StatusUnauthorized)
			return
		}

		if val != test.FakeAPIKey {
			http.Error(w, "Unauthorized: invalid key", http.StatusUnauthorized)
			return
		}

		mux.ServeHTTP(w, r)
	}))
	t.Cleanup(server.Close)

	c, err := New(
		WithClient(server.Client()),
		WithEndpoint(server.URL),
		WithKey(test.FakeAPIKey),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	return c, mux
}
