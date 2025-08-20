//go:build integration

package test

import (
	"context"
	"crypto/rand"
	"fmt"
	"os"
	"testing"

	"github.com/aws-gopher/unstructured-sdk-go"
)

func TestSourcePermutations(t *testing.T) {
	t.Parallel()

	if os.Getenv("UNSTRUCTURED_API_KEY") == "" {
		t.Skip("skipping because UNSTRUCTURED_API_KEY is not set")
	}

	client, err := unstructured.New()
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	for name, src := range map[string]unstructured.SourceConfigInput{
		"azure-account-key": unstructured.AzureSourceConnectorConfigInput{
			RemoteURL:   "az://foo",
			AccountName: S("foo"),
			AccountKey:  S("foo"),
		},
		"azure-connection-string": unstructured.AzureSourceConnectorConfigInput{
			RemoteURL:        "az://foo",
			ConnectionString: S("foo"),
		},
		"azure-sas-token": unstructured.AzureSourceConnectorConfigInput{
			RemoteURL:   "az://foo",
			AccountName: S("foo"),
			SASToken:    S("foo"),
		},

		"box": unstructured.BoxSourceConnectorConfigInput{
			BoxAppConfig: "foo",
			RemoteURL:    "box://foo",
		},

		// server responds 500
		// "confluence": unstructured.ConfluenceSourceConnectorConfigInput{
		// 	URL:      "https://foo.atlassian.net",
		// 	Username: "foo",
		// 	Password: S("foo"),
		// },

		"couchbase": unstructured.CouchbaseSourceConnectorConfigInput{
			Bucket:           "foo",
			ConnectionString: "couchbase://foo",
			Username:         "foo",
			Password:         "foo",
			CollectionID:     "foo",
			BatchSize:        100,
		},

		// server responds 500
		// "databricks-volumes": unstructured.DatabricksVolumesConnectorConfigInput{
		// 	Host:         "foo.cloud.databricks.com",
		// 	Catalog:      "foo",
		// 	Volume:       "foo",
		// 	VolumePath:   "/foo",
		// 	ClientSecret: "foo",
		// 	ClientID:     "foo",
		// },

		"dropbox": unstructured.DropboxSourceConnectorConfigInput{
			Token:     "foo",
			RemoteURL: "dropbox://foo",
		},

		"elasticsearch": unstructured.ElasticsearchConnectorConfigInput{
			Hosts:     []string{"https://foo.elastic-cloud.com"},
			IndexName: "foo",
			ESAPIKey:  "foo",
		},

		"gcs": unstructured.GCSSourceConnectorConfigInput{
			RemoteURL:         "gs://foo",
			ServiceAccountKey: "foo",
		},

		"google-drive": unstructured.GoogleDriveSourceConnectorConfigInput{
			DriveID:           "foo",
			ServiceAccountKey: S("foo"),
		},

		"jira": unstructured.JiraSourceConnectorConfigInput{
			URL:      "https://foo.atlassian.net",
			Username: "foo",
			Password: S("foo"),
		},

		// server responds 412 asking for `bootstrap_server` instead of `bootstrap_servers`
		// "kafka-cloud": unstructured.KafkaCloudSourceConnectorConfigInput{
		// 	BootstrapServers: "foo.cloud.confluent.io",
		// 	Topic:            "foo",
		// 	KafkaAPIKey:      "foo",
		// 	Secret:           "foo",
		// },

		"mongodb": unstructured.MongoDBConnectorConfigInput{
			Database:   "foo",
			Collection: "foo",
			URI:        "mongodb://foo",
		},

		"onedrive": unstructured.OneDriveSourceConnectorConfigInput{
			ClientID:     "foo",
			UserPName:    "foo",
			Tenant:       "foo",
			AuthorityURL: "https://login.microsoftonline.com/foo",
			ClientCred:   "foo",
			Path:         "/foo",
		},

		"outlook": unstructured.OutlookSourceConnectorConfigInput{
			ClientID:       "foo",
			ClientCred:     "foo",
			UserEmail:      "foo@example.com",
			OutlookFolders: []string{"Inbox"},
		},

		"postgres": unstructured.PostgresSourceConnectorConfigInput{
			Host:      "foo.com",
			Database:  "foo",
			Port:      5432,
			Username:  "foo",
			Password:  "foo",
			TableName: "foo",
			BatchSize: 100,
		},

		"s3": unstructured.S3SourceConnectorConfigInput{
			RemoteURL: "s3://foo",
			Key:       S("foo"),
			Secret:    S("foo"),
		},

		"salesforce": unstructured.SalesforceSourceConnectorConfigInput{
			Username:    "foo",
			ConsumerKey: "foo",
			PrivateKey:  "foo",
			Categories:  []string{"foo"},
		},

		"sharepoint": unstructured.SharePointSourceConnectorConfigInput{
			Site:       "https://foo.sharepoint.com/sites/foo",
			Tenant:     "foo",
			UserPName:  "foo",
			ClientID:   "foo",
			ClientCred: "foo",
		},

		// server responds 500
		// "snowflake": unstructured.SnowflakeSourceConnectorConfigInput{
		// 	Account:   "foo",
		// 	Role:      "foo",
		// 	User:      "foo",
		// 	Password:  "foo",
		// 	Host:      "foo.snowflakecomputing.com",
		// 	Database:  "foo",
		// 	TableName: S("foo"),
		// 	IDColumn:  S("foo"),
		// },

		"zendesk": unstructured.ZendeskSourceConnectorConfigInput{
			Subdomain: "foo",
			Email:     "foo@example.com",
			APIToken:  "foo",
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			source, err := client.CreateSource(t.Context(), unstructured.CreateSourceRequest{
				Name:   fmt.Sprintf("test-%s-%s", name, rand.Text()),
				Config: src,
			})
			if err != nil {
				t.Fatalf("failed to create source: %v", err)
			}

			t.Cleanup(func() { _ = client.DeleteSource(context.Background(), source.ID) })
		})
	}
}
