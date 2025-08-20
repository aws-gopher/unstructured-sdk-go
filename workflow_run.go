package unstructured

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

// RunWorkflowRequest represents the request to run a workflow
type RunWorkflowRequest struct {
	ID string

	// InputFiles is a list of files to upload to the workflow.
	// The files must implement the io.Reader interface.
	InputFiles []File
}

// File represents a file to upload to the workflow.
type File interface {
	Name() string
	io.Reader
}

// FileBytes implements the File interface for an io.Reader in memory.
type FileBytes struct {
	Filename string
	Bytes    io.Reader
}

// Name returns the name of the file.
func (f *FileBytes) Name() string { return f.Filename }

// Read reads the file into the given buffer.
func (f *FileBytes) Read(p []byte) (n int, err error) {
	return f.Bytes.Read(p) //nolint:wrapcheck
}

// RunWorkflow runs a workflow by triggering a new job
func (c *Client) RunWorkflow(ctx context.Context, in *RunWorkflowRequest) (*Job, error) {
	req, err := http.NewRequestWithContext(ctx,
		http.MethodPost,
		c.endpoint.JoinPath("workflows", in.ID, "run").String(),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Determine if we need to upload files
	if len(in.InputFiles) > 0 {
		if err := addfiles(req, in.InputFiles); err != nil {
			return nil, fmt.Errorf("failed to add files to request: %w", err)
		}
	}

	var job Job
	if err := c.do(req, &job); err != nil {
		return nil, fmt.Errorf("failed to run workflow: %w", err)
	}

	return &job, nil
}

// runWorkflowWithoutFiles handles the case where no files are provided
func addfiles(req *http.Request, files []File) error {
	// Create a buffer to hold the multipart form data
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Add each file to the multipart form
	for _, f := range files {
		// Create a form file field
		part, err := writer.CreateFormFile("input_files", f.Name())
		if err != nil {
			return fmt.Errorf("failed to create form file for %s: %w", f, err)
		}

		// Copy the file content to the form part
		if _, err := io.Copy(part, f); err != nil {
			return fmt.Errorf("failed to copy file content for %s: %w", f, err)
		}
	}

	// Close the multipart writer
	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close multipart writer: %w", err)
	}

	req.Body = io.NopCloser(&buf)

	// Set the content type header for multipart form data
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return nil
}
