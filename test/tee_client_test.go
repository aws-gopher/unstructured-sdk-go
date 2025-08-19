package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
)

type teert struct {
	dst  io.Writer
	next http.RoundTripper
}

func (c *teert) RoundTrip(req *http.Request) (*http.Response, error) {
	data, err := httputil.DumpRequestOut(req, false)
	if err != nil {
		return nil, fmt.Errorf("failed to dump request: %w", err)
	}

	c.dst.Write(data)

	if req.Body != nil {
		data, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read request body: %w", err)
		}

		req.Body = io.NopCloser(bytes.NewBuffer(data))

		switch req.Header.Get("Content-Type") {
		case "application/json":
			enc := json.NewEncoder(c.dst)
			enc.SetIndent("", "  ")

			if err := enc.Encode(json.RawMessage(data)); err != nil {
				return nil, fmt.Errorf("failed to encode request body: %w", err)
			}

		default:
			c.dst.Write([]byte("[content omitted]"))
		}
	}

	var res *http.Response

	res, err = c.next.RoundTrip(req)
	if err != nil {
		return nil, fmt.Errorf("failed to round trip request: %w", err)
	}

	data, err = httputil.DumpResponse(res, false)
	if err != nil {
		return nil, fmt.Errorf("failed to dump response: %w", err)
	}

	c.dst.Write([]byte("\n\n"))
	c.dst.Write(data)

	if res.Body != nil {
		data, err = io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		if len(data) == 0 {
			return res, nil
		}

		res.Body = io.NopCloser(bytes.NewBuffer(data))

		switch res.Header.Get("Content-Type") {
		case "application/json":
			enc := json.NewEncoder(c.dst)
			enc.SetIndent("", "  ")

			if err := enc.Encode(json.RawMessage(data)); err != nil {
				return nil, fmt.Errorf("failed to encode response body: %w", err)
			}

		default:
			c.dst.Write([]byte("[content omitted]"))
		}
	}

	c.dst.Write([]byte("\n\n"))

	return res, nil
}
