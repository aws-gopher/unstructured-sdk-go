package unstructured

// Connector type constants
const (
	ConnectorTypeAstraDB                    = "astradb"
	ConnectorTypeAzureAISearch              = "azure_ai_search"
	ConnectorTypeAzure                      = "azure"
	ConnectorTypeBox                        = "box"
	ConnectorTypeConfluence                 = "confluence"
	ConnectorTypeCouchbase                  = "couchbase"
	ConnectorTypeDatabricksVolumes          = "databricks_volumes"
	ConnectorTypeDatabricksVolumeDeltaTable = "databricks_volume_delta_tables"
	ConnectorTypeDeltaTable                 = "delta_table"
	ConnectorTypeDropbox                    = "dropbox"
	ConnectorTypeElasticsearch              = "elasticsearch"
	ConnectorTypeGCS                        = "gcs"
	ConnectorTypeGoogleDrive                = "google_drive"
	ConnectorTypeJira                       = "jira"
	ConnectorTypeKafkaCloud                 = "kafka-cloud"
	ConnectorTypeMilvus                     = "milvus"
	ConnectorTypeMongoDB                    = "mongodb"
	ConnectorTypeMotherDuck                 = "motherduck"
	ConnectorTypeNeo4j                      = "neo4j"
	ConnectorTypeOneDrive                   = "onedrive"
	ConnectorTypeOutlook                    = "outlook"
	ConnectorTypePinecone                   = "pinecone"
	ConnectorTypePostgres                   = "postgres"
	ConnectorTypeQdrantCloud                = "qdrant-cloud"
	ConnectorTypeRedis                      = "redis"
	ConnectorTypeS3                         = "s3"
	ConnectorTypeSalesforce                 = "salesforce"
	ConnectorTypeSharePoint                 = "sharepoint"
	ConnectorTypeSlack                      = "slack"
	ConnectorTypeSnowflake                  = "snowflake"
	ConnectorTypeWeaviateCloud              = "weaviate-cloud"
	ConnectorTypeZendesk                    = "zendesk"
	ConnectorTypeIBMWatsonxS3               = "ibm_watsonx_s3"
)

// DatabricksVolumesConnectorConfigInput represents the configuration for a Databricks Volumes connector.
// It contains host details, catalog information, and authentication credentials.
type DatabricksVolumesConnectorConfigInput struct {
	sourceconfig
	destinationconfig

	Host         string  `json:"host"`
	Catalog      string  `json:"catalog"`
	Schema       *string `json:"schema,omitempty"`
	Volume       string  `json:"volume"`
	VolumePath   string  `json:"volume_path"`
	ClientSecret string  `json:"client_secret"`
	ClientID     string  `json:"client_id"`
}

// Type always returns the connector type identifier for Databricks Volumes: "databricks_volumes".
func (c DatabricksVolumesConnectorConfigInput) Type() string { return ConnectorTypeDatabricksVolumes }

// ElasticsearchConnectorConfigInput represents the configuration for an Elasticsearch connector.
// It contains host details, index information, and API key authentication.
type ElasticsearchConnectorConfigInput struct {
	sourceconfig
	destinationconfig

	Hosts     []string `json:"hosts"`
	IndexName string   `json:"index_name"`
	ESAPIKey  string   `json:"es_api_key"`
}

// Type always returns the connector type identifier for Elasticsearch: "elasticsearch".
func (c ElasticsearchConnectorConfigInput) Type() string { return ConnectorTypeElasticsearch }

// MongoDBConnectorConfigInput represents the configuration for a MongoDB connector.
// It contains database connection details and collection information.
type MongoDBConnectorConfigInput struct {
	sourceconfig
	destinationconfig

	Database   string `json:"database"`
	Collection string `json:"collection"`
	URI        string `json:"uri"`
}

// Type always returns the connector type identifier for MongoDB: "mongodb".
func (c MongoDBConnectorConfigInput) Type() string { return ConnectorTypeMongoDB }

// DatabricksVolumesConnectorConfig represents the configuration for a Databricks Volumes connector.
// It contains host details, catalog information, and authentication credentials.
type DatabricksVolumesConnectorConfig struct {
	sourceconfig
	destinationconfig

	Host         string  `json:"host"`
	Catalog      string  `json:"catalog"`
	Schema       *string `json:"schema,omitempty"`
	Volume       string  `json:"volume"`
	VolumePath   string  `json:"volume_path"`
	ClientSecret string  `json:"client_secret"`
	ClientID     string  `json:"client_id"`
}

var _ SourceConfig = (*DatabricksVolumesConnectorConfig)(nil)

// Type always returns the connector type identifier for Databricks Volumes: "databricks_volumes".
func (c DatabricksVolumesConnectorConfig) Type() string { return ConnectorTypeDatabricksVolumes }

// ElasticsearchConnectorConfig represents the configuration for an Elasticsearch connector.
// It contains host details, index information, and API key authentication.
type ElasticsearchConnectorConfig struct {
	sourceconfig
	destinationconfig

	Hosts     []string `json:"hosts"`
	IndexName string   `json:"index_name"`
	ESAPIKey  string   `json:"es_api_key"`
}

var _ SourceConfig = (*ElasticsearchConnectorConfig)(nil)

// Type always returns the connector type identifier for Elasticsearch: "elasticsearch".
func (c ElasticsearchConnectorConfig) Type() string { return ConnectorTypeElasticsearch }

// MongoDBConnectorConfig represents the configuration for a MongoDB connector.
// It contains database connection details and collection information.
type MongoDBConnectorConfig struct {
	sourceconfig
	destinationconfig

	Database   string `json:"database"`
	Collection string `json:"collection"`
	URI        string `json:"uri"`
}

var _ SourceConfig = (*MongoDBConnectorConfig)(nil)

// Type always returns the connector type identifier for MongoDB: "mongodb".
func (c MongoDBConnectorConfig) Type() string { return ConnectorTypeMongoDB }

// CouchbaseConnectorConfig represents the configuration for a Couchbase connector.
// It contains connection details, bucket information, and authentication credentials.
type CouchbaseConnectorConfig struct {
	sourceconfig
	destinationconfig

	Bucket           string  `json:"bucket"`
	ConnectionString string  `json:"connection_string"`
	Scope            *string `json:"scope,omitempty"`
	Collection       *string `json:"collection,omitempty"`
	BatchSize        int     `json:"batch_size"`
	Username         string  `json:"username"`
	Password         string  `json:"password"`
	CollectionID     *string `json:"collection_id,omitempty"`
}

var _ SourceConfig = (*CouchbaseConnectorConfig)(nil)
var _ DestinationConfig = (*CouchbaseConnectorConfig)(nil)

// Type always returns the connector type identifier for Couchbase: "couchbase".
func (c CouchbaseConnectorConfig) Type() string { return ConnectorTypeCouchbase }

// S3ConnectorConfig represents the configuration for an S3 connector.
// It supports both AWS S3 and S3-compatible storage services.
type S3ConnectorConfig struct {
	sourceconfig
	destinationconfig

	RemoteURL   string  `json:"remote_url"`
	Anonymous   *bool   `json:"anonymous,omitempty"`
	Key         *string `json:"key,omitempty"`
	Secret      *string `json:"secret,omitempty"`
	Token       *string `json:"token,omitempty"`
	EndpointURL *string `json:"endpoint_url,omitempty"`
	Recursive   *bool   `json:"recursive,omitempty"`
}

var _ SourceConfig = (*S3ConnectorConfig)(nil)
var _ DestinationConfig = (*S3ConnectorConfig)(nil)

// Type always returns the connector type identifier for S3: "s3".
func (c S3ConnectorConfig) Type() string { return ConnectorTypeS3 }

// GCSConnectorConfig represents the configuration for a Google Cloud Storage connector.
// It contains the remote URL and service account key for authentication.
type GCSConnectorConfig struct {
	sourceconfig
	destinationconfig

	RemoteURL         string `json:"remote_url"`
	ServiceAccountKey string `json:"service_account_key"`
	Recursive         *bool  `json:"recursive,omitempty"`
}

var _ SourceConfig = (*GCSConnectorConfig)(nil)
var _ DestinationConfig = (*GCSConnectorConfig)(nil)

// Type always returns the connector type identifier for GCS: "gcs".
func (c GCSConnectorConfig) Type() string { return ConnectorTypeGCS }

// KafkaCloudConnectorConfig represents the configuration for a Kafka Cloud connector.
// It contains broker details, topic information, and authentication credentials.
type KafkaCloudConnectorConfig struct {
	sourceconfig
	destinationconfig

	BootstrapServers     string  `json:"bootstrap_servers"`
	Port                 *int    `json:"port,omitempty"`
	GroupID              *string `json:"group_id,omitempty"`
	Topic                string  `json:"topic"`
	KafkaAPIKey          string  `json:"kafka_api_key"`
	Secret               string  `json:"secret"`
	NumMessagesToConsume *int    `json:"num_messages_to_consume,omitempty"`
	BatchSize            *int    `json:"batch_size,omitempty"`
}

var _ SourceConfig = (*KafkaCloudConnectorConfig)(nil)
var _ DestinationConfig = (*KafkaCloudConnectorConfig)(nil)

// Type always returns the connector type identifier for Kafka Cloud: "kafka-cloud".
func (c KafkaCloudConnectorConfig) Type() string { return ConnectorTypeKafkaCloud }

// PostgresConnectorConfig represents the configuration for a PostgreSQL connector.
// It contains database connection details and table configuration.
type PostgresConnectorConfig struct {
	sourceconfig
	destinationconfig

	Host      string   `json:"host"`
	Database  string   `json:"database"`
	Port      int      `json:"port"`
	Username  string   `json:"username"`
	Password  string   `json:"password"`
	TableName string   `json:"table_name"`
	BatchSize int      `json:"batch_size"`
	IDColumn  *string  `json:"id_column,omitempty"`
	Fields    []string `json:"fields,omitempty"`
}

var _ SourceConfig = (*PostgresConnectorConfig)(nil)
var _ DestinationConfig = (*PostgresConnectorConfig)(nil)

// Type always returns the connector type identifier for PostgreSQL: "postgres".
func (c PostgresConnectorConfig) Type() string { return ConnectorTypePostgres }

// SnowflakeConnectorConfig represents the configuration for a Snowflake connector.
// It contains account details, authentication, and table configuration.
type SnowflakeConnectorConfig struct {
	sourceconfig
	destinationconfig

	Account     string   `json:"account"`
	Role        string   `json:"role"`
	User        string   `json:"user"`
	Password    string   `json:"password"`
	Host        string   `json:"host"`
	Port        *int     `json:"port,omitempty"`
	Database    string   `json:"database"`
	Schema      *string  `json:"schema,omitempty"`
	TableName   *string  `json:"table_name,omitempty"`
	BatchSize   *int     `json:"batch_size,omitempty"`
	IDColumn    *string  `json:"id_column,omitempty"`
	Fields      []string `json:"fields,omitempty"`
	RecordIDKey *string  `json:"record_id_key,omitempty"`
}

var _ SourceConfig = (*SnowflakeConnectorConfig)(nil)
var _ DestinationConfig = (*SnowflakeConnectorConfig)(nil)

// Type always returns the connector type identifier for Snowflake: "snowflake".
func (c SnowflakeConnectorConfig) Type() string { return ConnectorTypeSnowflake }

// OneDriveConnectorConfig represents the configuration for a OneDrive connector.
// It contains Microsoft Graph API authentication and file access settings.
type OneDriveConnectorConfig struct {
	sourceconfig
	destinationconfig

	ClientID     string  `json:"client_id"`
	UserPName    string  `json:"user_pname"`
	Tenant       string  `json:"tenant"`
	AuthorityURL string  `json:"authority_url"`
	ClientCred   string  `json:"client_cred"`
	Recursive    *bool   `json:"recursive,omitempty"`
	Path         *string `json:"path,omitempty"`
	RemoteURL    *string `json:"remote_url,omitempty"`
}

var _ SourceConfig = (*OneDriveConnectorConfig)(nil)
var _ DestinationConfig = (*OneDriveConnectorConfig)(nil)

// Type always returns the connector type identifier for OneDrive: "onedrive".
func (c OneDriveConnectorConfig) Type() string { return ConnectorTypeOneDrive }
