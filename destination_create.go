package unstructured

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// CreateDestination creates a new destination connector with the specified configuration.
// It returns the created destination connector with its assigned ID and metadata.
func (c *Client) CreateDestination(ctx context.Context, in CreateDestinationRequest) (*Destination, error) {
	config, err := json.Marshal(in.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config: %w", err)
	}

	shadow := struct {
		Name   string          `json:"name"`
		Type   string          `json:"type"`
		Config json.RawMessage `json:"config"`
	}{
		Name:   in.Name,
		Type:   in.Config.Type(),
		Config: json.RawMessage(config),
	}

	body, err := json.Marshal(shadow)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal destination request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx,
		http.MethodPost,
		c.endpoint.JoinPath("destinations/").String(),
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	var destination Destination
	if err := c.do(req, &destination); err != nil {
		return nil, fmt.Errorf("failed to create destination: %w", err)
	}

	return &destination, nil
}

// CreateDestinationRequest represents a request to create a new destination connector.
// It contains the name, type, and configuration for the destination.
type CreateDestinationRequest struct {
	Name   string
	Config DestinationConfigInput
}

// DestinationConfigInput is an interface that all destination connector configurations must implement.
// It provides a way to identify the type of destination connector and marshal its configuration.
type DestinationConfigInput interface {
	isDestinationConfigInput()
	Type() string
}

type destinationconfiginput struct{}

func (s destinationconfiginput) isDestinationConfigInput() {}

// AstraDBConnectorConfigInput represents the configuration for an AstraDB destination connector.
// It contains the collection name, keyspace, batch size, API endpoint, and token.
type AstraDBConnectorConfigInput struct {
	destinationconfiginput

	CollectionName  string  `json:"collection_name"`
	Keyspace        *string `json:"keyspace,omitempty"`
	BatchSize       *int    `json:"batch_size,omitempty"`
	APIEndpoint     string  `json:"api_endpoint"`
	Token           string  `json:"token"`
	FlattenMetadata *bool   `json:"flatten_metadata,omitempty"`
}

// Type always returns the connector type identifier for AstraDB: "astradb".
func (c AstraDBConnectorConfigInput) Type() string { return ConnectorTypeAstraDB }

// AzureAISearchConnectorConfigInput represents the configuration for an Azure AI Search destination connector.
// It contains the endpoint, index name, and API key.
type AzureAISearchConnectorConfigInput struct {
	destinationconfiginput

	Endpoint string `json:"endpoint"`
	Index    string `json:"index"`
	Key      string `json:"key"`
}

// Type always returns the connector type identifier for Azure AI Search: "azure_ai_search".
func (c AzureAISearchConnectorConfigInput) Type() string { return ConnectorTypeAzureAISearch }

// CouchbaseDestinationConnectorConfigInput represents the configuration for a Couchbase destination connector.
// It contains connection details, bucket information, and authentication credentials.
type CouchbaseDestinationConnectorConfigInput struct {
	destinationconfiginput

	Bucket           string  `json:"bucket"`
	ConnectionString string  `json:"connection_string"`
	Scope            *string `json:"scope,omitempty"`
	Collection       *string `json:"collection,omitempty"`
	BatchSize        int     `json:"batch_size"`
	Username         string  `json:"username"`
	Password         string  `json:"password"`
}

// Type always returns the connector type identifier for Couchbase: "couchbase".
func (c CouchbaseDestinationConnectorConfigInput) Type() string { return ConnectorTypeCouchbase }

// DatabricksVDTDestinationConnectorConfigInput represents the configuration for a Databricks Volume Delta Tables destination connector.
// It contains server details, authentication, and table configuration.
type DatabricksVDTDestinationConnectorConfigInput struct {
	destinationconfiginput

	ServerHostname string  `json:"server_hostname"`
	HTTPPath       string  `json:"http_path"`
	Token          *string `json:"token,omitempty"`
	ClientID       *string `json:"client_id,omitempty"`
	ClientSecret   *string `json:"client_secret,omitempty"`
	Catalog        string  `json:"catalog"`
	Database       *string `json:"database,omitempty"`
	TableName      *string `json:"table_name,omitempty"`
	Schema         *string `json:"schema,omitempty"`
	Volume         string  `json:"volume"`
	VolumePath     *string `json:"volume_path,omitempty"`
}

// Type always returns the connector type identifier for Databricks Volume Delta Tables: "databricks_volume_delta_tables".
func (c DatabricksVDTDestinationConnectorConfigInput) Type() string {
	return ConnectorTypeDatabricksVolumeDeltaTable
}

// DeltaTableConnectorConfigInput represents the configuration for a Delta Table destination connector.
// It contains AWS credentials and table URI for Delta Lake storage.
type DeltaTableConnectorConfigInput struct {
	destinationconfiginput

	AwsAccessKeyID     string `json:"aws_access_key_id"`
	AwsSecretAccessKey string `json:"aws_secret_access_key"`
	AwsRegion          string `json:"aws_region"`
	TableURI           string `json:"table_uri"`
}

// Type always returns the connector type identifier for Delta Table: "delta_table".
func (c DeltaTableConnectorConfigInput) Type() string { return ConnectorTypeDeltaTable }

// GCSDestinationConnectorConfigInput represents the configuration for a Google Cloud Storage destination connector.
// It contains the remote URL and service account key for authentication.
type GCSDestinationConnectorConfigInput struct {
	destinationconfiginput

	RemoteURL         string `json:"remote_url"`
	ServiceAccountKey string `json:"service_account_key"`
}

// Type always returns the connector type identifier for Google Cloud Storage: "gcs".
func (c GCSDestinationConnectorConfigInput) Type() string { return ConnectorTypeGCS }

// KafkaCloudDestinationConnectorConfigInput represents the configuration for a Kafka Cloud destination connector.
// It contains broker details, topic information, and authentication credentials.
type KafkaCloudDestinationConnectorConfigInput struct {
	destinationconfiginput

	BootstrapServers string  `json:"bootstrap_servers"`
	Port             *int    `json:"port,omitempty"`
	GroupID          *string `json:"group_id,omitempty"`
	Topic            string  `json:"topic"`
	KafkaAPIKey      string  `json:"kafka_api_key"`
	Secret           string  `json:"secret"`
	BatchSize        *int    `json:"batch_size,omitempty"`
}

// Type always returns the connector type identifier for Kafka Cloud: "kafka-cloud".
func (c KafkaCloudDestinationConnectorConfigInput) Type() string { return ConnectorTypeKafkaCloud }

// MilvusDestinationConnectorConfigInput represents the configuration for a Milvus destination connector.
// It contains connection details, collection information, and authentication.
type MilvusDestinationConnectorConfigInput struct {
	destinationconfiginput

	URI            string  `json:"uri"`
	User           *string `json:"user,omitempty"`
	Token          *string `json:"token,omitempty"`
	Password       *string `json:"password,omitempty"`
	DBName         *string `json:"db_name,omitempty"`
	CollectionName string  `json:"collection_name"`
	RecordIDKey    string  `json:"record_id_key"`
}

// Type always returns the connector type identifier for Milvus: "milvus".
func (c MilvusDestinationConnectorConfigInput) Type() string { return ConnectorTypeMilvus }

// MotherduckDestinationConnectorConfigInput represents the configuration for a MotherDuck destination connector.
// It contains database connection details and authentication credentials.
type MotherduckDestinationConnectorConfigInput struct {
	destinationconfiginput

	Account     string  `json:"account"`
	Role        string  `json:"role"`
	User        string  `json:"user"`
	Password    string  `json:"password"`
	Host        string  `json:"host"`
	Port        *int    `json:"port,omitempty"`
	Database    string  `json:"database"`
	Schema      *string `json:"schema,omitempty"`
	TableName   *string `json:"table_name,omitempty"`
	BatchSize   *int    `json:"batch_size,omitempty"`
	RecordIDKey *string `json:"record_id_key,omitempty"`
}

// Type always returns the connector type identifier for Motherduck: "motherduck".
func (c MotherduckDestinationConnectorConfigInput) Type() string { return ConnectorTypeMotherDuck }

// Neo4jDestinationConnectorConfigInput represents the configuration for a Neo4j destination connector.
// It contains database connection details and authentication credentials.
type Neo4jDestinationConnectorConfigInput struct {
	destinationconfiginput

	URI       string `json:"uri"`
	Database  string `json:"database"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	BatchSize *int   `json:"batch_size,omitempty"`
}

// Type always returns the connector type identifier for Neo4j: "neo4j".
func (c Neo4jDestinationConnectorConfigInput) Type() string { return ConnectorTypeNeo4j }

// OneDriveDestinationConnectorConfigInput represents the configuration for a OneDrive destination connector.
// It contains Microsoft Graph API authentication and file storage details.
type OneDriveDestinationConnectorConfigInput struct {
	destinationconfiginput

	ClientID     string `json:"client_id"`
	UserPName    string `json:"user_pname"`
	Tenant       string `json:"tenant"`
	AuthorityURL string `json:"authority_url"`
	ClientCred   string `json:"client_cred"`
	RemoteURL    string `json:"remote_url"`
}

// Type always returns the connector type identifier for OneDrive: "onedrive".
func (c OneDriveDestinationConnectorConfigInput) Type() string { return ConnectorTypeOneDrive }

// PineconeDestinationConnectorConfigInput represents the configuration for a Pinecone destination connector.
// It contains index details, API key, and namespace information.
type PineconeDestinationConnectorConfigInput struct {
	destinationconfiginput

	IndexName string `json:"index_name"`
	APIKey    string `json:"api_key"`
	Namespace string `json:"namespace"`
	BatchSize *int   `json:"batch_size,omitempty"`
}

// Type always returns the connector type identifier for Pinecone: "pinecone".
func (c PineconeDestinationConnectorConfigInput) Type() string { return ConnectorTypePinecone }

// PostgresDestinationConnectorConfigInput represents the configuration for a PostgreSQL destination connector.
// It contains database connection details and table configuration.
type PostgresDestinationConnectorConfigInput struct {
	destinationconfiginput

	Host      string `json:"host"`
	Database  string `json:"database"`
	Port      int    `json:"port"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	TableName string `json:"table_name"`
	BatchSize int    `json:"batch_size"`
}

// Type always returns the connector type identifier for PostgreSQL: "postgres".
func (c PostgresDestinationConnectorConfigInput) Type() string { return ConnectorTypePostgres }

// RedisDestinationConnectorConfigInput represents the configuration for a Redis destination connector.
// It contains connection details, database selection, and authentication.
type RedisDestinationConnectorConfigInput struct {
	destinationconfiginput

	Host      string  `json:"host"`
	Port      *int    `json:"port,omitempty"`
	Username  *string `json:"username,omitempty"`
	Password  *string `json:"password,omitempty"`
	URI       *string `json:"uri,omitempty"`
	Database  *int    `json:"database,omitempty"`
	SSL       *bool   `json:"ssl,omitempty"`
	BatchSize *int    `json:"batch_size,omitempty"`
}

// Type always returns the connector type identifier for Redis: "redis".
func (c RedisDestinationConnectorConfigInput) Type() string { return ConnectorTypeRedis }

// QdrantCloudDestinationConnectorConfigInput represents the configuration for a Qdrant Cloud destination connector.
// It contains API endpoint, collection details, and authentication.
type QdrantCloudDestinationConnectorConfigInput struct {
	destinationconfiginput

	URL            string `json:"url"`
	APIKey         string `json:"api_key"`
	CollectionName string `json:"collection_name"`
	BatchSize      *int   `json:"batch_size,omitempty"`
}

// Type always returns the connector type identifier for Qdrant Cloud: "qdrant-cloud".
func (c QdrantCloudDestinationConnectorConfigInput) Type() string { return ConnectorTypeQdrantCloud }

// S3DestinationConnectorConfigInput represents the configuration for an Amazon S3 destination connector.
// It supports both AWS S3 and S3-compatible storage services for storing processed data.
type S3DestinationConnectorConfigInput struct {
	destinationconfiginput

	RemoteURL   string  `json:"remote_url"`
	Anonymous   *bool   `json:"anonymous,omitempty"`
	Key         *string `json:"key,omitempty"`
	Secret      *string `json:"secret,omitempty"`
	Token       *string `json:"token,omitempty"`
	EndpointURL *string `json:"endpoint_url,omitempty"`
}

// Type always returns the connector type identifier for S3: "s3".
func (c S3DestinationConnectorConfigInput) Type() string { return ConnectorTypeS3 }

// SnowflakeDestinationConnectorConfigInput represents the configuration for a Snowflake destination connector.
// It contains account details, authentication, and table configuration.
type SnowflakeDestinationConnectorConfigInput struct {
	destinationconfiginput

	Account     string  `json:"account"`
	Role        string  `json:"role"`
	User        string  `json:"user"`
	Password    string  `json:"password"`
	Host        string  `json:"host"`
	Port        *int    `json:"port,omitempty"`
	Database    string  `json:"database"`
	Schema      *string `json:"schema,omitempty"`
	TableName   *string `json:"table_name,omitempty"`
	BatchSize   *int    `json:"batch_size,omitempty"`
	RecordIDKey *string `json:"record_id_key,omitempty"`
}

// Type always returns the connector type identifier for Snowflake: "snowflake".
func (c SnowflakeDestinationConnectorConfigInput) Type() string { return ConnectorTypeSnowflake }

// WeaviateDestinationConnectorConfigInput represents the configuration for a Weaviate destination connector.
// It contains cluster URL, API key, and collection information.
type WeaviateDestinationConnectorConfigInput struct {
	destinationconfiginput

	ClusterURL string  `json:"cluster_url"`
	APIKey     string  `json:"api_key"`
	Collection *string `json:"collection,omitempty"`
}

// Type always returns the connector type identifier for Weaviate Cloud: "weaviate-cloud".
func (c WeaviateDestinationConnectorConfigInput) Type() string { return ConnectorTypeWeaviateCloud }

// IBMWatsonxS3DestinationConnectorConfigInput represents the configuration for an IBM Watsonx S3 destination connector.
// It contains IBM Cloud authentication, storage endpoints, and table configuration.
type IBMWatsonxS3DestinationConnectorConfigInput struct {
	destinationconfiginput

	IAMApiKey             string  `json:"iam_api_key"`
	AccessKeyID           string  `json:"access_key_id"`
	SecretAccessKey       string  `json:"secret_access_key"`
	IcebergEndpoint       string  `json:"iceberg_endpoint"`
	ObjectStorageEndpoint string  `json:"object_storage_endpoint"`
	ObjectStorageRegion   string  `json:"object_storage_region"`
	Catalog               string  `json:"catalog"`
	MaxRetriesConnection  *int    `json:"max_retries_connection,omitempty"`
	Namespace             string  `json:"namespace"`
	Table                 string  `json:"table"`
	MaxRetries            *int    `json:"max_retries,omitempty"`
	RecordIDKey           *string `json:"record_id_key,omitempty"`
}

// Type always returns the connector type identifier for IBM Watsonx S3: "ibm_watsonx_s3".
func (c IBMWatsonxS3DestinationConnectorConfigInput) Type() string { return ConnectorTypeIBMWatsonxS3 }
