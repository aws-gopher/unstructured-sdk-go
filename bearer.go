package unstructured

import (
	"net/http"
)

type bearer struct {
	key string
	rt  http.RoundTripper
}

// RoundTrip implements the http.RoundTripper interface.
func (b *bearer) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Unstructured-API-Key", b.key)

	// This is implementing the http.RoundTripper interface, errors should be passed through as-is
	return b.rt.RoundTrip(req) //nolint:wrapcheck
}
