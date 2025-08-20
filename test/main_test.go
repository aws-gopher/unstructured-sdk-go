package test

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/aws-gopher/unstructured-sdk-go"
)

func TestWorkflow(t *testing.T) {
	t.Skip()

	pretty := func(v any) string {
		data, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			t.Fatalf("failed to marshal: %v", err)
		}

		return string(data)
	}

	client, err := unstructured.New(
		// unstructured.WithClient(&http.Client{
		//	Transport: &teert{
		//		next: http.DefaultTransport,
		//		dst:  t.Output(),
		//	},
		// }),
		unstructured.WithKey(os.Getenv("UNSTRUCTURED_API_KEY")),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx := t.Context()

	workflow, err := client.CreateWorkflow(ctx, &unstructured.CreateWorkflowRequest{
		Name:         "test",
		WorkflowType: unstructured.WorkflowTypeCustom,
		WorkflowNodes: []unstructured.WorkflowNode{
			&unstructured.PartitionerAuto{
				Name: "Partitioner",
			},
		},
	})
	if err != nil {
		t.Fatalf("failed to create workflow: %v", err)
	} else {
		t.Logf("created workflow %s:\n%s", workflow.ID, pretty(workflow))
	}

	t.Cleanup(func() { _ = client.DeleteWorkflow(ctx, workflow.ID) })

	// get all the dir under ./testdata and use them in a call to run the workflow.
	dir, err := os.ReadDir("testdata")
	if err != nil {
		t.Fatalf("failed to read testdata: %v", err)
	}

	files := make([]unstructured.File, 0, len(dir))

	for _, file := range dir {
		if file.IsDir() || file.Name() == ".DS_Store" || strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		f, err := os.Open(filepath.Join("testdata", file.Name()))
		if err != nil {
			t.Fatalf("failed to open file: %v", err)
		}

		t.Cleanup(func() { _ = f.Close() })

		files = append(files, f)
	}

	t.Logf("running workflow %s with %d files", workflow.ID, len(files))

	job, err := client.RunWorkflow(ctx, &unstructured.RunWorkflowRequest{
		ID:         workflow.ID,
		InputFiles: files,
	})
	if err != nil {
		t.Errorf("failed to run workflow: %v", err)
		return
	}

	t.Logf("job %s:\n%s", job.ID, pretty(job))

	deadline, ok := t.Deadline()
	if !ok {
		deadline = time.Now().Add(5 * time.Minute)
	}

	tick := time.NewTicker(5 * time.Second)
	defer tick.Stop()

	for ctx, cancel := context.WithDeadline(ctx, deadline.Add(-1*time.Second)); ; {
		select {
		case <-ctx.Done():
			t.Error("job took too long")
			cancel()

			return

		case <-tick.C:
		}

		last := job.Status

		job, err = client.GetJob(ctx, job.ID)
		if err != nil {
			t.Errorf("failed to get job: %v", err)
			cancel()

			return
		}

		if job.Status != last {
			t.Logf("%s => %s (%s):\n%s", last, job.Status, time.Since(job.CreatedAt), pretty(job))
		}

		if job.Status == unstructured.JobStatusCompleted {
			cancel()
			break
		}
	}

	for _, node := range job.OutputNodeFiles {
		download, err := client.DownloadJob(ctx, unstructured.DownloadJobRequest{
			JobID:  job.ID,
			NodeID: node.NodeID,
			FileID: node.FileID,
		})
		if err != nil {
			t.Errorf("failed to download job: %v", err)
			return
		}

		t.Cleanup(func() { _ = download.Close() })

		f, err := os.Create(filepath.Join("testdata", node.FileID+".json"))
		if err != nil {
			t.Errorf("failed to create file: %v", err)
			return
		}

		t.Cleanup(func() { _ = f.Close() })

		if _, err := io.Copy(f, download); err != nil {
			t.Errorf("failed to read download: %v", err)
		}
	}
}
