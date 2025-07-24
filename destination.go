package unstructured

import (
	"encoding/json"
	"fmt"
	"time"
)

// destinationConfigFactories maps destination type strings to factory functions that create new instances of the appropriate concrete destination config type.
// Using a map here also provides a compile-time check that all destination type strings are unique.
var destinationConfigFactories = map[string]func() DestinationConfig{
	ConnectorTypeAstraDB:                    func() DestinationConfig { return new(AstraDBConnectorConfig) },
	ConnectorTypeAzureAISearch:              func() DestinationConfig { return new(AzureAISearchConnectorConfig) },
	ConnectorTypeCouchbase:                  func() DestinationConfig { return new(CouchbaseDestinationConnectorConfig) },
	ConnectorTypeDatabricksVolumes:          func() DestinationConfig { return new(DatabricksVolumesConnectorConfig) },
	ConnectorTypeDatabricksVolumeDeltaTable: func() DestinationConfig { return new(DatabricksVDTDestinationConnectorConfig) },
	ConnectorTypeDeltaTable:                 func() DestinationConfig { return new(DeltaTableConnectorConfig) },
	ConnectorTypeElasticsearch:              func() DestinationConfig { return new(ElasticsearchConnectorConfig) },
	ConnectorTypeGCS:                        func() DestinationConfig { return new(GCSDestinationConnectorConfig) },
	ConnectorTypeKafkaCloud:                 func() DestinationConfig { return new(KafkaCloudDestinationConnectorConfig) },
	ConnectorTypeMilvus:                     func() DestinationConfig { return new(MilvusDestinationConnectorConfig) },
	ConnectorTypeMongoDB:                    func() DestinationConfig { return new(MongoDBConnectorConfig) },
	ConnectorTypeMotherDuck:                 func() DestinationConfig { return new(MotherduckDestinationConnectorConfig) },
	ConnectorTypeNeo4j:                      func() DestinationConfig { return new(Neo4jDestinationConnectorConfig) },
	ConnectorTypeOneDrive:                   func() DestinationConfig { return new(OneDriveDestinationConnectorConfig) },
	ConnectorTypePinecone:                   func() DestinationConfig { return new(PineconeDestinationConnectorConfig) },
	ConnectorTypePostgres:                   func() DestinationConfig { return new(PostgresDestinationConnectorConfig) },
	ConnectorTypeRedis:                      func() DestinationConfig { return new(RedisDestinationConnectorConfig) },
	ConnectorTypeQdrantCloud:                func() DestinationConfig { return new(QdrantCloudDestinationConnectorConfig) },
	ConnectorTypeS3:                         func() DestinationConfig { return new(S3DestinationConnectorConfig) },
	ConnectorTypeSnowflake:                  func() DestinationConfig { return new(SnowflakeDestinationConnectorConfig) },
	ConnectorTypeWeaviateCloud:              func() DestinationConfig { return new(WeaviateDestinationConnectorConfig) },
	ConnectorTypeIBMWatsonxS3:               func() DestinationConfig { return new(IBMWatsonxS3DestinationConnectorConfig) },
}

// Destination represents a destination connector that sends processed data to various locations.
// It contains metadata about the connector and its configuration.
type Destination struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	CreatedAt time.Time         `json:"created_at,omitzero"`
	UpdatedAt time.Time         `json:"updated_at,omitzero"`
	Type      string            `json:"type"`
	Config    DestinationConfig `json:"config"`
}

// UnmarshalJSON implements custom JSON unmarshaling for Destination.
// It handles the polymorphic Config field by determining the correct type
// based on the "type" field in the JSON data.
func (d *Destination) UnmarshalJSON(data []byte) error {
	var shadow struct {
		ID        string          `json:"id"`
		Name      string          `json:"name"`
		CreatedAt time.Time       `json:"created_at"`
		UpdatedAt time.Time       `json:"updated_at"`
		Type      string          `json:"type"`
		Config    json.RawMessage `json:"config"`
	}

	if err := json.Unmarshal(data, &shadow); err != nil {
		return fmt.Errorf("failed to unmarshal destination: %w", err)
	}

	d.ID = shadow.ID
	d.Name = shadow.Name
	d.CreatedAt = shadow.CreatedAt
	d.UpdatedAt = shadow.UpdatedAt
	d.Type = shadow.Type

	// Look up the factory function for this destination type
	factory, exists := destinationConfigFactories[shadow.Type]
	if !exists {
		return fmt.Errorf("unknown destination type: %s", shadow.Type)
	}

	// Create a new instance of the appropriate config type
	config := factory()

	// Unmarshal the config data into the concrete type
	if err := json.Unmarshal(shadow.Config, config); err != nil {
		return fmt.Errorf("failed to unmarshal %s config: %w", shadow.Type, err)
	}

	d.Config = config

	return nil
}

// DestinationConfig is an interface that all destination connector configurations implement.
// It provides a way to identify and work with different destination connector types.
type DestinationConfig interface {
	isDestinationConfig()
}

type destinationconfig struct{}

func (d destinationconfig) isDestinationConfig() {}

// AstraDBConnectorConfig represents the configuration for an AstraDB destination connector.
// It contains the collection name, keyspace, batch size, API endpoint, and token.
type AstraDBConnectorConfig struct {
	destinationconfig

	CollectionName string  `json:"collection_name"`
	Keyspace       *string `json:"keyspace,omitempty"`
	BatchSize      int     `json:"batch_size"`
	APIEndpoint    string  `json:"api_endpoint"`
	Token          string  `json:"token"`
}

// AzureAISearchConnectorConfig represents the configuration for an Azure AI Search destination connector.
// It contains the endpoint, index name, and API key.
type AzureAISearchConnectorConfig struct {
	destinationconfig

	Endpoint string `json:"endpoint"`
	Index    string `json:"index"`
	Key      string `json:"key"`
}

// CouchbaseDestinationConnectorConfig represents the configuration for a Couchbase destination connector.
// It contains connection details, bucket information, and authentication credentials.
type CouchbaseDestinationConnectorConfig struct {
	destinationconfig

	Bucket           string  `json:"bucket"`
	ConnectionString string  `json:"connection_string"`
	Scope            *string `json:"scope,omitempty"`
	Collection       *string `json:"collection,omitempty"`
	BatchSize        int     `json:"batch_size"`
	Username         string  `json:"username"`
	Password         string  `json:"password"`
}

// DatabricksVDTDestinationConnectorConfig represents the configuration for a Databricks Volume Delta Tables destination connector.
// It contains server details, authentication, and table configuration.
type DatabricksVDTDestinationConnectorConfig struct {
	destinationconfig

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

// DeltaTableConnectorConfig represents the configuration for a Delta Table destination connector.
// It contains AWS credentials and table URI for Delta Lake storage.
type DeltaTableConnectorConfig struct {
	destinationconfig

	AwsAccessKeyID     string `json:"aws_access_key_id"`
	AwsSecretAccessKey string `json:"aws_secret_access_key"`
	AwsRegion          string `json:"aws_region"`
	TableURI           string `json:"table_uri"`
}

// GCSDestinationConnectorConfig represents the configuration for a Google Cloud Storage destination connector.
// It contains the remote URL and service account key for authentication.
type GCSDestinationConnectorConfig struct {
	destinationconfig

	RemoteURL         string `json:"remote_url"`
	ServiceAccountKey string `json:"service_account_key"`
}

// KafkaCloudDestinationConnectorConfig represents the configuration for a Kafka Cloud destination connector.
// It contains broker details, topic information, and authentication credentials.
type KafkaCloudDestinationConnectorConfig struct {
	destinationconfig

	BootstrapServers string  `json:"bootstrap_servers"`
	Port             *int    `json:"port,omitempty"`
	GroupID          *string `json:"group_id,omitempty"`
	Topic            string  `json:"topic"`
	KafkaAPIKey      string  `json:"kafka_api_key"`
	Secret           string  `json:"secret"`
	BatchSize        *int    `json:"batch_size,omitempty"`
}

// MilvusDestinationConnectorConfig represents the configuration for a Milvus destination connector.
// It contains connection details, collection information, and authentication.
type MilvusDestinationConnectorConfig struct {
	destinationconfig

	URI            string  `json:"uri"`
	User           *string `json:"user,omitempty"`
	Token          *string `json:"token,omitempty"`
	Password       *string `json:"password,omitempty"`
	DBName         *string `json:"db_name,omitempty"`
	CollectionName string  `json:"collection_name"`
	RecordIDKey    string  `json:"record_id_key"`
}

// Neo4jDestinationConnectorConfig represents the configuration for a Neo4j destination connector.
// It contains database connection details and authentication credentials.
type Neo4jDestinationConnectorConfig struct {
	destinationconfig

	URI       string `json:"uri"`
	Database  string `json:"database"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	BatchSize *int   `json:"batch_size,omitempty"`
}

// MotherduckDestinationConnectorConfig represents the configuration for a MotherDuck destination connector.
// It contains database connection details and authentication credentials.
type MotherduckDestinationConnectorConfig struct {
	destinationconfig

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

// OneDriveDestinationConnectorConfig represents the configuration for a OneDrive destination connector.
// It contains Microsoft Graph API authentication and file storage details.
type OneDriveDestinationConnectorConfig struct {
	destinationconfig

	ClientID     string `json:"client_id"`
	UserPName    string `json:"user_pname"`
	Tenant       string `json:"tenant"`
	AuthorityURL string `json:"authority_url"`
	ClientCred   string `json:"client_cred"`
	RemoteURL    string `json:"remote_url"`
}

// PineconeDestinationConnectorConfig represents the configuration for a Pinecone destination connector.
// It contains index details, API key, and namespace information.
type PineconeDestinationConnectorConfig struct {
	destinationconfig

	IndexName string `json:"index_name"`
	APIKey    string `json:"api_key"`
	Namespace string `json:"namespace"`
	BatchSize *int   `json:"batch_size,omitempty"`
}

// PostgresDestinationConnectorConfig represents the configuration for a PostgreSQL destination connector.
// It contains database connection details and table configuration.
type PostgresDestinationConnectorConfig struct {
	destinationconfig

	Host      string `json:"host"`
	Database  string `json:"database"`
	Port      int    `json:"port"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	TableName string `json:"table_name"`
	BatchSize int    `json:"batch_size"`
}

// RedisDestinationConnectorConfig represents the configuration for a Redis destination connector.
// It contains connection details, database selection, and authentication.
type RedisDestinationConnectorConfig struct {
	destinationconfig

	Host      string  `json:"host"`
	Port      *int    `json:"port,omitempty"`
	Username  *string `json:"username,omitempty"`
	Password  *string `json:"password,omitempty"`
	URI       *string `json:"uri,omitempty"`
	Database  *int    `json:"database,omitempty"`
	SSL       *bool   `json:"ssl,omitempty"`
	BatchSize *int    `json:"batch_size,omitempty"`
}

// QdrantCloudDestinationConnectorConfig represents the configuration for a Qdrant Cloud destination connector.
// It contains API endpoint, collection details, and authentication.
type QdrantCloudDestinationConnectorConfig struct {
	destinationconfig

	URL            string `json:"url"`
	APIKey         string `json:"api_key"`
	CollectionName string `json:"collection_name"`
	BatchSize      *int   `json:"batch_size,omitempty"`
}

// S3DestinationConnectorConfig represents the configuration for an Amazon S3 destination connector.
// It supports both AWS S3 and S3-compatible storage services for storing processed data.
type S3DestinationConnectorConfig struct {
	destinationconfig

	RemoteURL   string  `json:"remote_url"`
	Anonymous   bool    `json:"anonymous"`
	Key         *string `json:"key,omitempty"`
	Secret      *string `json:"secret,omitempty"`
	Token       *string `json:"token,omitempty"`
	EndpointURL *string `json:"endpoint_url,omitempty"`
}

// SnowflakeDestinationConnectorConfig represents the configuration for a Snowflake destination connector.
// It contains account details, authentication, and table configuration.
type SnowflakeDestinationConnectorConfig struct {
	destinationconfig

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

// WeaviateDestinationConnectorConfig represents the configuration for a Weaviate destination connector.
// It contains cluster URL, API key, and collection information.
type WeaviateDestinationConnectorConfig struct {
	destinationconfig

	ClusterURL string  `json:"cluster_url"`
	APIKey     string  `json:"api_key"`
	Collection *string `json:"collection,omitempty"`
}

// IBMWatsonxS3DestinationConnectorConfig represents the configuration for an IBM Watsonx S3 destination connector.
// It contains IBM Cloud authentication, storage endpoints, and table configuration.
type IBMWatsonxS3DestinationConnectorConfig struct {
	destinationconfig

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
