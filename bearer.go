package unstructured

import (
	"net/http"
)

type bearer struct {
	key string
	rt  http.RoundTripper
}

// HeaderKey is "Unstructured-API-Key", which is the header where Unstructured expects to find the API key.
const HeaderKey = "Unstructured-API-Key"

// RoundTrip implements the http.RoundTripper interface.
func (b *bearer) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set(HeaderKey, b.key)

	// This is implementing the http.RoundTripper interface, errors should be passed through as-is
	return b.rt.RoundTrip(req) //nolint:wrapcheck
}
