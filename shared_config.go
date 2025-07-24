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

// Shared connector config types that work for both source and destination

type sharedconfiginput struct{}

func (s sharedconfiginput) isSourceConfigInput()      {}
func (s sharedconfiginput) isDestinationConfigInput() {}

// DatabricksVolumesConnectorConfigInput represents the configuration for a Databricks Volumes connector.
// It contains host details, catalog information, and authentication credentials.
type DatabricksVolumesConnectorConfigInput struct {
	sharedconfiginput

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
	sharedconfiginput

	Hosts     []string `json:"hosts"`
	IndexName string   `json:"index_name"`
	ESAPIKey  string   `json:"es_api_key"`
}

// Type always returns the connector type identifier for Elasticsearch: "elasticsearch".
func (c ElasticsearchConnectorConfigInput) Type() string { return ConnectorTypeElasticsearch }

// MongoDBConnectorConfigInput represents the configuration for a MongoDB connector.
// It contains database connection details and collection information.
type MongoDBConnectorConfigInput struct {
	sharedconfiginput

	Database   string `json:"database"`
	Collection string `json:"collection"`
	URI        string `json:"uri"`
}

// Type always returns the connector type identifier for MongoDB: "mongodb".
func (c MongoDBConnectorConfigInput) Type() string { return ConnectorTypeMongoDB }

// Shared connector config types that work for both source and destination
type sharedconfig struct{}

func (s sharedconfig) isSourceConfig()      {}
func (s sharedconfig) isDestinationConfig() {}

// DatabricksVolumesConnectorConfig represents the configuration for a Databricks Volumes connector.
// It contains host details, catalog information, and authentication credentials.
type DatabricksVolumesConnectorConfig struct {
	sharedconfig

	Host         string  `json:"host"`
	Catalog      string  `json:"catalog"`
	Schema       *string `json:"schema,omitempty"`
	Volume       string  `json:"volume"`
	VolumePath   string  `json:"volume_path"`
	ClientSecret string  `json:"client_secret"`
	ClientID     string  `json:"client_id"`
}

// ElasticsearchConnectorConfig represents the configuration for an Elasticsearch connector.
// It contains host details, index information, and API key authentication.
type ElasticsearchConnectorConfig struct {
	sharedconfig

	Hosts     []string `json:"hosts"`
	IndexName string   `json:"index_name"`
	ESAPIKey  string   `json:"es_api_key"`
}

// MongoDBConnectorConfig represents the configuration for a MongoDB connector.
// It contains database connection details and collection information.
type MongoDBConnectorConfig struct {
	sharedconfig

	Database   string `json:"database"`
	Collection string `json:"collection"`
	URI        string `json:"uri"`
}
