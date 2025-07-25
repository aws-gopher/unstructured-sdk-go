package unstructured

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

// DownloadJob downloads the output files from a completed job
func (c *Client) DownloadJob(ctx context.Context, id string) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx,
		http.MethodGet,
		c.endpoint.JoinPath("/jobs", id, "download").String(),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute HTTP request: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download job: %s", resp.Status)
	}

	return resp.Body, nil
}
