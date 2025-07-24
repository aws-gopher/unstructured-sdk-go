# Unstructured.io Go SDK

A lightweight Go client for the Unstructured.io Workflow Endpoint API with **zero external dependencies**.

[![Go Reference](https://pkg.go.dev/badge/github.com/aws-gopher/unstructured-sdk-go.svg)](https://pkg.go.dev/github.com/aws-gopher/unstructured-sdk-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/aws-gopher/unstructured-sdk-go)](https://goreportcard.com/report/github.com/aws-gopher/unstructured-sdk-go)

## Features

- ✅ **Zero external dependencies** - Uses only Go standard library
- ✅ (Nearly) Complete coverage for Unstructured.io Workflow API
- ✅ Type-safe request/response structures
- ✅ Comprehensive error handling
- ✅ Connection testing utilities
- ✅ Helper functions for optional fields
- ✅ Context support for timeouts and cancellation

## Installation

```bash
go get github.com/aws-gopher/unstructured-sdk-go
```

## Quick Start

```go
package main

import (
    "context"
    "log"
    "github.com/aws-gopher/unstructured-sdk-go"
)

func main() {
    // Create a client
    client, err := unstructured.New(
        unstructured.WithEndpoint("https://platform.unstructured.io/api/v1"),
        unstructured.WithKey("your-api-key"),
    )
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // Create a source connector (S3)
    source, err := client.CreateSource(ctx, unstructured.CreateSourceRequest{
        Name: "My S3 Source",
        Config: unstructured.S3SourceConnectorConfigInput{
            RemoteURL: "s3://my-bucket/input/",
            Key:       unstructured.String("your-access-key"),
            Secret:    unstructured.String("your-secret-key"),
        },
    })
    if err != nil {
        log.Fatal(err)
    }

    // Create a destination connector (S3)
    destination, err := client.CreateDestination(ctx, unstructured.CreateDestinationRequest{
        Name: "My S3 Destination",
        Config: unstructured.S3DestinationConnectorConfigInput{
            RemoteURL: "s3://my-bucket/output/",
            Key:       unstructured.String("your-access-key"),
            Secret:    unstructured.String("your-secret-key"),
        },
    })
    if err != nil {
        log.Fatal(err)
    }

    // Create a workflow
    workflow, err := client.CreateWorkflow(ctx, unstructured.CreateWorkflowRequest{
        Name:          "My Processing Workflow",
        SourceID:      &source.ID,
        DestinationID: &destination.ID,
        WorkflowType:  unstructured.WorkflowTypeBasic,
        WorkflowNodes: []unstructured.WorkflowNode{
            {
                Name:    "Partitioner",
                Type:    "partition",
                Subtype: "fast",
            },
            {
                Name:    "Chunker",
                Type:    "chunk",
                Subtype: "by_title",
                Settings: map[string]interface{}{
                    "chunk_size": 1000,
                    "overlap":    200,
                },
            },
        },
    })
    if err != nil {
        log.Fatal(err)
    }

    // Run the workflow
    job, err := client.RunWorkflow(ctx, workflow.ID, &unstructured.RunWorkflowRequest{
        InputFiles: []string{"document1.pdf", "document2.docx"},
    })
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Job started with ID: %s", job.ID)
}
```

## Key Concepts

### Connectors

- **Source Connectors**: Ingest files or data into Unstructured from source locations like S3, Google Drive, databases, etc.
- **Destination Connectors**: Send processed data from Unstructured to destination locations like S3, databases, vector stores, etc.

### Workflows

Workflows define how Unstructured processes your data through a series of nodes:

- **Source Node**: Represents where your files or data come from
- **Partitioner Node**: Extracts content from unstructured files and outputs structured document elements
- **Chunker Node**: Chunks partitioned data into smaller pieces for RAG applications
- **Enrichment Node**: Applies enrichments like image summaries, table summaries, NER, etc.
- **Embedder Node**: Generates vector embeddings for vector-based searches
- **Destination Node**: Represents where processed data goes

### Jobs

Jobs run workflows at specific points in time and can be monitored for status and results.

## Working with Different Connector Types

### Source Connectors

```go
// Azure Blob Storage
azureSource, err := client.CreateSource(ctx, unstructured.CreateSourceRequest{
    Name: "Azure Source",
    Config: unstructured.AzureSourceConnectorConfigInput{
        RemoteURL:        "https://myaccount.blob.core.windows.net/container/",
        ConnectionString: unstructured.String("your-connection-string"),
    },
})

// Google Drive
gdriveSource, err := client.CreateSource(ctx, unstructured.CreateSourceRequest{
    Name: "Google Drive Source",
    Config: unstructured.GoogleDriveSourceConnectorConfigInput{
        DriveID:           "your-drive-id",
        ServiceAccountKey: unstructured.String("your-service-account-key"),
        Extensions:        []string{".pdf", ".docx", ".txt"},
    },
})

// Salesforce
salesforceSource, err := client.CreateSource(ctx, unstructured.CreateSourceRequest{
    Name: "Salesforce Source",
    Config: unstructured.SalesforceSourceConnectorConfigInput{
        Username:    "your-username",
        ConsumerKey: "your-consumer-key",
        PrivateKey:  "your-private-key",
        Categories:  []string{"cases", "opportunities"},
    },
})
```

### Destination Connectors

```go
// S3 Destination
s3Dest, err := client.CreateDestination(ctx, unstructured.CreateDestinationRequest{
    Name: "S3 Destination",
    Config: unstructured.S3DestinationConnectorConfigInput{
        RemoteURL: "s3://my-bucket/processed/",
        Key:       unstructured.String("your-access-key"),
        Secret:    unstructured.String("your-secret-key"),
    },
})

// Postgres Database
postgresDest, err := client.CreateDestination(ctx, unstructured.CreateDestinationRequest{
    Name: "Postgres Destination",
    Config: unstructured.PostgresDestinationConnectorConfigInput{
        Host:      "your-postgres-host",
        Database:  "your-database",
        Port:      5432,
        Username:  "your-username",
        Password:  "your-password",
        TableName: "processed_documents",
    },
})
```

## Managing Workflows

```go
// List workflows with filtering
workflows, _ := client.ListWorkflows(ctx, &unstructured.ListWorkflowsRequest{
    Status:        &unstructured.WorkflowStateActive,
    Page:          unstructured.Int(1),
    PageSize:      unstructured.Int(10),
    SortBy:        unstructured.String("created_at"),
    SortDirection: &unstructured.SortDirectionDesc,
})

// Get workflow details
workflow, _ := client.GetWorkflow(ctx, "workflow-id")

// Update workflow
updatedWorkflow, _ := client.UpdateWorkflow(ctx, "workflow-id", unstructured.UpdateWorkflowRequest{
    Name: unstructured.String("Updated Workflow Name"),
    WorkflowNodes: []unstructured.WorkflowNode{
        {
            Name:    "Partitioner",
            Type:    "partition",
            Subtype: "fast",
        },
        {
            Name:    "Chunker",
            Type:    "chunk",
            Subtype: "by_title",
            Settings: map[string]interface{}{
                "chunk_size": 1500,
                "overlap":    300,
            },
        },
        {
            Name:    "Embedder",
            Type:    "embed",
            Subtype: "openai",
            Settings: map[string]interface{}{
                "model": "text-embedding-ada-002",
            },
        },
    },
})
```

## Monitoring Jobs

```go
// List jobs
jobs, _ := client.ListJobs(ctx, &unstructured.ListJobsRequest{
    WorkflowID: unstructured.String("workflow-id"),
    Status:     &unstructured.JobStatusCompleted,
})

// Get job details
job, _ := client.GetJob(ctx, "job-id")

// Get detailed processing information
jobDetails, _ := client.GetJobDetails(ctx, "job-id")

// Check for failed files
failedFiles, err := client.GetJobFailedFiles(ctx, "job-id")
if err == nil && len(failedFiles.FailedFiles) > 0 {
    for _, failed := range failedFiles.FailedFiles {
        log.Printf("Failed file: %s, Error: %s", failed.Document, failed.Error)
    }
}

// Download job results
reader, err := client.DownloadJob(ctx, "job-id")
if err != nil {
    log.Fatal(err)
}
defer reader.Close()

// Save to file
file, err := os.Create("job-results.json")
if err != nil {
    log.Fatal(err)
}
defer file.Close()

_, err = io.Copy(file, reader)
if err != nil {
    log.Fatal(err)
}
```

## Connection Testing

```go
// Test source connector connection
connectionCheck, err := client.CreateSourceConnectionCheck(ctx, "source-id")
if err != nil {
    log.Fatal(err)
}

// Check connection status
checkResult, err := client.GetSourceConnectionCheck(ctx, "source-id")
if err != nil {
    log.Fatal(err)
}

switch checkResult.Status {
case unstructured.ConnectionCheckStatusSuccess:
    log.Println("Connection successful")

case unstructured.ConnectionCheckStatusFailure:
    log.Printf("Connection failed: %s", *checkResult.Reason)

case unstructured.ConnectionCheckStatusScheduled:
    log.Println("Connection check in progress")
}
```

## Error Handling

The package provides comprehensive error handling with typed errors:

```go
source, err := client.CreateSource(ctx, request)
if err != nil {
    // Check for validation errors
    ve := new(HTTPValidationError)
    if errors.As(err, &ve) {
        log.Printf("Validation failed: %v", ve)
        for _, detail := range ve.Detail {
            log.Printf("  - %s at %v: %s", detail.Type, detail.Location, detail.Message)
        }
        return
    }
    
    // Handle other errors
    log.Printf("Source creation failed: %v", err)
    return
}
```

## Helper Functions

The package provides several helper functions for working with pointers to primitive types. These functions are useful when you need to pass optional values to API requests.

### Creating pointers from values

```go
str := unstructured.String("optional value")
enabled := unstructured.Bool(true)
count := unstructured.Int(42)
```

### Converting pointers back to values with safe defaults

```go
value := unstructured.ToString(str)    // returns "" if str is nil
flag := unstructured.ToBool(enabled)   // returns false if enabled is nil
number := unstructured.ToInt(count)    // returns 0 if count is nil
```

These helper functions are particularly useful when working with optional fields in request structures:

```go
workflow, err := client.CreateWorkflow(ctx, unstructured.CreateWorkflowRequest{
    Name:         "My Workflow",
    WorkflowType: unstructured.WorkflowTypeBasic,
    ReprocessAll: unstructured.Bool(true),  // Optional boolean field
    Page:         unstructured.Int(1),      // Optional integer field
})
```

## Supported File Types

The Unstructured.io platform supports a wide variety of file types including:

- **Documents**: PDF, DOCX, PPTX, XLSX, TXT, RTF
- **Images**: JPG, PNG, TIFF, BMP
- **Archives**: ZIP, TAR, RAR
- **Web**: HTML, XML, JSON
- And many more

## Authentication

The package supports API key authentication:

```go
client, err := unstructured.New(
    unstructured.WithEndpoint("https://platform.unstructured.io/api/v1"),
    unstructured.WithKey("your-api-key"),
)
```

## Rate Limiting and Best Practices

- Use `context.Context` for timeout and cancellation
- Implement proper error handling and retry logic
- Monitor job status before attempting downloads
- Use connection checks to validate connector configurations
- Consider implementing exponential backoff for retries

## Dependencies

This package has **zero external dependencies** and uses only the Go standard library:

- `context` - For context support
- `encoding/json` - For JSON marshaling/unmarshaling
- `fmt` - For string formatting
- `io` - For I/O operations
- `net/http` - For HTTP client functionality
- `strings` - For string operations
- `time` - For time-related operations

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## API Reference

For more information about the Unstructured.io API, visit the [Unstructured Workflow API docs](https://docs.unstructured.io/api-reference/workflow/overview).
