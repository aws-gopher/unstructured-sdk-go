/*
Package unstructured provides a Go client for the Unstructured.io Workflow Endpoint API.

The Unstructured.io Workflow Endpoint enables you to work with connectors, workflows, and jobs programmatically.
This package provides a complete Go implementation of the API, allowing you to:

  - Create and manage source connectors that ingest files or data from various locations
  - Create and manage destination connectors that send processed data to different destinations
  - Define and manage workflows that specify how Unstructured processes your data
  - Run jobs that execute workflows at specific points in time

# Key Concepts

# Connectors

  - Source Connectors: Ingest files or data into Unstructured from source locations like S3, Google Drive, databases, etc.
  - Destination Connectors: Send processed data from Unstructured to destination locations like S3, databases, vector stores, etc.

# Workflows

Workflows define how Unstructured processes your data through a series of nodes:

  - Source Node: Represents where your files or data come from
  - Partitioner Node: Extracts content from unstructured files and outputs structured document elements
  - Chunker Node: Chunks partitioned data into smaller pieces for RAG applications
  - Enrichment Node: Applies enrichments like image summaries, table summaries, NER, etc.
  - Embedder Node: Generates vector embeddings for vector-based searches
  - Destination Node: Represents where processed data goes

# Jobs

Jobs run workflows at specific points in time and can be monitored for status and results.

Quick Start

	package main

	import (
		"context"
		"log"
		"github.com/your-org/unstructured-sdk-go"
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

# Working with Different Connector Types

Source Connectors

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

Destination Connectors

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
			Host:     "your-postgres-host",
			Database: "your-database",
			Port:     5432,
			Username: "your-username",
			Password: "your-password",
			TableName: "processed_documents",
		},
	})

Managing Workflows

	// List workflows with filtering
	workflows, err := client.ListWorkflows(ctx, &unstructured.ListWorkflowsRequest{
		Status:        &unstructured.WorkflowStateActive,
		Page:          unstructured.Int(1),
		PageSize:      unstructured.Int(10),
		SortBy:        unstructured.String("created_at"),
		SortDirection: &unstructured.SortDirectionDesc,
	})

	// Get workflow details
	workflow, err := client.GetWorkflow(ctx, "workflow-id")

	// Update workflow
	updatedWorkflow, err := client.UpdateWorkflow(ctx, "workflow-id", unstructured.UpdateWorkflowRequest{
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

Monitoring Jobs

	// List jobs
	jobs, err := client.ListJobs(ctx, &unstructured.ListJobsRequest{
		WorkflowID: unstructured.String("workflow-id"),
		Status:     &unstructured.JobStatusCompleted,
	})

	// Get job details
	job, err := client.GetJob(ctx, "job-id")

	// Get detailed processing information
	jobDetails, err := client.GetJobDetails(ctx, "job-id")

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

Connection Testing

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

# Error Handling

The package provides comprehensive error handling:

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

# Supported File Types

The Unstructured.io platform supports a wide variety of file types including:

  - Documents: PDF, DOCX, PPTX, XLSX, TXT, RTF
  - Images: JPG, PNG, TIFF, BMP
  - Archives: ZIP, TAR, RAR
  - Web: HTML, XML, JSON
  - And many more

# Rate Limiting and Best Practices

  - Use [context.Context] for timeout and cancellation
  - Implement proper error handling and retry logic
  - Monitor job status before attempting downloads
  - Use connection checks to validate connector configurations
  - Consider implementing exponential backoff for retries

# Authentication

The package supports API key authentication:

	client, err := unstructured.New(
		unstructured.WithEndpoint("https://platform.unstructured.io/api/v1"),
		unstructured.WithKey("your-api-key"),
	)

# Helper Functions

The package provides several helper functions for working with pointers to primitive types.
These functions are useful when you need to pass optional values to API requests.

Creating pointers from values:

	str := unstructured.String("optional value")
	enabled := unstructured.Bool(true)
	count := unstructured.Int(42)

Converting pointers back to values with safe defaults:

	value := unstructured.ToString(str) // returns "" if str is nil
	flag := unstructured.ToBool(enabled) // returns false if enabled is nil
	number := unstructured.ToInt(count)  // returns 0 if count is nil

These helper functions are particularly useful when working with optional fields in request structures:

	workflow, err := client.CreateWorkflow(ctx, unstructured.CreateWorkflowRequest{
		Name:          "My Workflow",
		WorkflowType:  unstructured.WorkflowTypeBasic,
		ReprocessAll:  unstructured.Bool(true),  // Optional boolean field
		Page:          unstructured.Int(1),      // Optional integer field
	})

For more information about the Unstructured.io API, visit:
https://docs.unstructured.io/api-reference/workflow/overview
*/
package unstructured
