//go:build integration

package test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/aws-gopher/unstructured-sdk-go"
)

func TestDestinationPermutations(t *testing.T) {
	t.Parallel()

	if os.Getenv("UNSTRUCTURED_API_KEY") == "" {
		t.Skip("skipping because UNSTRUCTURED_API_KEY is not set")
	}

	client, err := unstructured.New()
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	for name, src := range map[string]unstructured.DestinationConfig{
		"astra-db": &unstructured.AstraDBConnectorConfig{
			CollectionName: "foo",
			APIEndpoint:    "https://foo.apps.astra.datastax.com",
			Token:          "foo",
		},

		"azure-ai-search": &unstructured.AzureAISearchConnectorConfig{
			Endpoint: "https://foo.search.windows.net",
			Index:    "foo",
			Key:      "foo",
		},

		"couchbase": &unstructured.CouchbaseConnectorConfig{
			Bucket:           "foo",
			ConnectionString: "couchbase://foo",
			Username:         "foo",
			Password:         "foo",
			BatchSize:        100,
		},

		// server responds 500
		// "databricks-volume-delta-table": unstructured.DatabricksVDTDestinationConnectorConfig{
		// 	ServerHostname: "foo.cloud.databricks.com",
		// 	HTTPPath:       "/sql/1.0/warehouses/foo",
		// 	Token:          S("foo"),
		// 	Catalog:        "foo",
		// 	Volume:         "foo",
		// },

		"delta-table": &unstructured.DeltaTableConnectorConfig{
			AwsAccessKeyID:     "foo",
			AwsSecretAccessKey: "foo",
			AwsRegion:          "us-east-1",
			TableURI:           "s3://foo/table",
		},

		"elasticsearch": &unstructured.ElasticsearchConnectorConfig{
			Hosts:     []string{"https://foo.elastic-cloud.com"},
			IndexName: "foo",
			ESAPIKey:  "foo",
		},

		"gcs": &unstructured.GCSConnectorConfig{
			RemoteURL:         "gs://foo",
			ServiceAccountKey: "foo",
		},

		// server responds 412 asking for `bootstrap_server` instead of `bootstrap_servers`
		// "kafka-cloud": unstructured.KafkaCloudDestinationConnectorConfig{
		// 	BootstrapServers: "foo.cloud.confluent.io",
		// 	Topic:            "foo",
		// 	KafkaAPIKey:      "foo",
		// 	Secret:           "foo",
		// },

		"milvus-token": &unstructured.MilvusDestinationConnectorConfig{
			URI:            "https://foo.zilliz.com",
			CollectionName: "foo",
			RecordIDKey:    "foo",
			Token:          S("foo"),
		},
		"milvus-password": &unstructured.MilvusDestinationConnectorConfig{
			URI:            "https://foo.zilliz.com",
			CollectionName: "foo",
			RecordIDKey:    "foo",
			User:           S("foo"),
			Password:       S("foo"),
		},

		"mongo-db": &unstructured.MongoDBConnectorConfig{
			Database:   "foo",
			Collection: "foo",
			URI:        "mongodb://foo:27017/foo",
		},

		// server responds 422: Destination Connector type motherduck not supported
		// "mother-duck": unstructured.MotherduckDestinationConnectorConfig{
		// 	Account:  "foo",
		// 	Role:     "foo",
		// 	User:     "foo",
		// 	Password: "foo",
		// 	Host:     "foo.duckdb.io",
		// 	Database: "foo",
		// },

		"neo4j": &unstructured.Neo4jDestinationConnectorConfig{
			URI:      "bolt://foo:7687",
			Database: "foo",
			Username: "foo",
			Password: "foo",
		},

		"one-drive": &unstructured.OneDriveConnectorConfig{
			ClientID:     "foo",
			UserPName:    "foo",
			Tenant:       "foo",
			AuthorityURL: "https://login.microsoftonline.com/foo",
			ClientCred:   "foo",
			RemoteURL:    S("onedrive://foo"),
		},

		"pinecone": &unstructured.PineconeDestinationConnectorConfig{
			IndexName: "foo",
			APIKey:    "foo",
			Namespace: "foo",
		},

		"postgres": &unstructured.PostgresConnectorConfig{
			Host:      "foo.com",
			Database:  "foo",
			Port:      5432,
			Username:  "foo",
			Password:  "foo",
			TableName: "foo",
			BatchSize: 100,
		},

		"redis": &unstructured.RedisDestinationConnectorConfig{
			Host:     "foo.com",
			Username: S("foo"),
			Password: S("foo"),
		},

		"qdrant-cloud": &unstructured.QdrantCloudDestinationConnectorConfig{
			URL:            "https://foo.qdrant.io",
			APIKey:         "foo",
			CollectionName: "foo",
		},

		"s3": &unstructured.S3ConnectorConfig{
			RemoteURL: "s3://foo",
			Key:       S("foo"),
			Secret:    S("foo"),
		},

		// server responds 500
		// "snowflake": unstructured.SnowflakeDestinationConnectorConfig{
		// 	Account:  "foo",
		// 	Role:     "foo",
		// 	User:     "foo",
		// 	Password: "foo",
		// 	Host:     "foo.snowflakecomputing.com",
		// 	Database: "foo",
		// },

		"weaviate-cloud": &unstructured.WeaviateDestinationConnectorConfig{
			ClusterURL: "https://foo.weaviate.network",
			APIKey:     "foo",
		},

		"ibm-watsonx-s3": &unstructured.IBMWatsonxS3DestinationConnectorConfig{
			IAMApiKey:             "foo",
			AccessKeyID:           "foo",
			SecretAccessKey:       "foo",
			IcebergEndpoint:       "https://foo.iceberg.cloud.ibm.com",
			ObjectStorageEndpoint: "https://foo.s3.cloud.ibm.com",
			ObjectStorageRegion:   "us-east",
			Catalog:               "foo",
			Namespace:             "foo",
			Table:                 "foo",
		},

		// server responds 500
		// "databricks-volumes": unstructured.DatabricksVolumesConnectorConfig{
		// 	Host:         "foo.cloud.databricks.com",
		// 	Catalog:      "foo",
		// 	Volume:       "foo",
		// 	VolumePath:   "/foo",
		// 	ClientSecret: "foo",
		// 	ClientID:     "foo",
		// },
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			destination, err := client.CreateDestination(testContext(t), unstructured.CreateDestinationRequest{
				Name:   fmt.Sprintf("test-%s-%s", name, randText()),
				Config: src,
			})
			if err != nil {
				t.Fatalf("failed to create destination: %v", err)
			}

			t.Cleanup(func() { _ = client.DeleteDestination(context.Background(), destination.ID) })
		})
	}
}
