package unstructured

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

// DownloadJobRequest represents a request to download a job output file.
type DownloadJobRequest struct {
	JobID  string
	NodeID string
	FileID string
}

// DownloadJob downloads the output files from a completed job
func (c *Client) DownloadJob(ctx context.Context, in DownloadJobRequest) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx,
		http.MethodGet,
		c.endpoint.JoinPath("/jobs", in.JobID, "download").String(),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	q := req.URL.Query()
	q.Add("node_id", in.NodeID)
	q.Add("file_id", in.FileID)
	req.URL.RawQuery = q.Encode()

	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute HTTP request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download job: %s", resp.Status)
	}

	return resp.Body, nil
}
