package unstructured

import (
	"encoding/json"
	"fmt"
	"time"
)

// sourceConfigFactories maps source type strings to factory functions
// that create new instances of the appropriate concrete source config type.
var sourceConfigFactories = map[string]func() SourceConfig{
	ConnectorTypeAzure:             func() SourceConfig { return &AzureSourceConnectorConfig{} },
	ConnectorTypeBox:               func() SourceConfig { return &BoxSourceConnectorConfig{} },
	ConnectorTypeConfluence:        func() SourceConfig { return &ConfluenceSourceConnectorConfig{} },
	ConnectorTypeCouchbase:         func() SourceConfig { return &CouchbaseSourceConnectorConfig{} },
	ConnectorTypeDatabricksVolumes: func() SourceConfig { return &DatabricksVolumesConnectorConfig{} },
	ConnectorTypeDropbox:           func() SourceConfig { return &DropboxSourceConnectorConfig{} },
	ConnectorTypeElasticsearch:     func() SourceConfig { return &ElasticsearchConnectorConfig{} },
	ConnectorTypeGCS:               func() SourceConfig { return &GCSSourceConnectorConfig{} },
	ConnectorTypeGoogleDrive:       func() SourceConfig { return &GoogleDriveSourceConnectorConfig{} },
	ConnectorTypeJira:              func() SourceConfig { return &JiraSourceConnectorConfig{} },
	ConnectorTypeKafkaCloud:        func() SourceConfig { return &KafkaCloudSourceConnectorConfig{} },
	ConnectorTypeMongoDB:           func() SourceConfig { return &MongoDBConnectorConfig{} },
	ConnectorTypeOneDrive:          func() SourceConfig { return &OneDriveSourceConnectorConfig{} },
	ConnectorTypeOutlook:           func() SourceConfig { return &OutlookSourceConnectorConfig{} },
	ConnectorTypePostgres:          func() SourceConfig { return &PostgresSourceConnectorConfig{} },
	ConnectorTypeS3:                func() SourceConfig { return &S3SourceConnectorConfig{} },
	ConnectorTypeSalesforce:        func() SourceConfig { return &SalesforceSourceConnectorConfig{} },
	ConnectorTypeSharePoint:        func() SourceConfig { return &SharePointSourceConnectorConfig{} },
	ConnectorTypeSlack:             func() SourceConfig { return &SlackSourceConnectorConfig{} },
	ConnectorTypeSnowflake:         func() SourceConfig { return &SnowflakeSourceConnectorConfig{} },
	ConnectorTypeZendesk:           func() SourceConfig { return &ZendeskSourceConnectorConfig{} },
}

// Source represents a source connector that ingests files or data from various locations.
// It contains metadata about the connector and its configuration.
type Source struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	CreatedAt time.Time    `json:"created_at,omitzero"`
	UpdatedAt time.Time    `json:"updated_at,omitzero"`
	Type      string       `json:"type"`
	Config    SourceConfig `json:"config"`
}

// UnmarshalJSON implements custom JSON unmarshaling for Source.
// It handles the polymorphic Config field by determining the correct type
// based on the "type" field in the JSON data.
func (s *Source) UnmarshalJSON(data []byte) error {
	var shadow struct {
		ID        string          `json:"id"`
		Name      string          `json:"name"`
		CreatedAt time.Time       `json:"created_at,omitzero"`
		UpdatedAt time.Time       `json:"updated_at,omitzero"`
		Type      string          `json:"type"`
		Config    json.RawMessage `json:"config"`
	}

	if err := json.Unmarshal(data, &shadow); err != nil {
		return fmt.Errorf("failed to unmarshal source: %w", err)
	}

	s.ID = shadow.ID
	s.Name = shadow.Name
	s.CreatedAt = shadow.CreatedAt
	s.UpdatedAt = shadow.UpdatedAt
	s.Type = shadow.Type

	// Look up the factory function for this source type
	factory, exists := sourceConfigFactories[shadow.Type]
	if !exists {
		return fmt.Errorf("unknown source type: %s", shadow.Type)
	}

	// Create a new instance of the appropriate config type
	config := factory()

	// Unmarshal the config data into the concrete type
	if err := json.Unmarshal(shadow.Config, config); err != nil {
		return fmt.Errorf("failed to unmarshal %s config: %w", shadow.Type, err)
	}

	s.Config = config

	return nil
}

// SourceConfig is an interface that all source connector configurations implement.
// It provides a way to identify and work with different source connector types.
type SourceConfig interface {
	isSourceConfig()
}

type sourceconfig struct{}

func (s sourceconfig) isSourceConfig() {}

// AzureSourceConnectorConfig represents the configuration for an Azure Blob Storage source connector.
// It supports authentication via connection string, account key, or SAS token.
type AzureSourceConnectorConfig struct {
	sourceconfig

	RemoteURL        string  `json:"remote_url"`
	AccountName      *string `json:"account_name,omitempty"`
	AccountKey       *string `json:"account_key,omitempty"`
	ConnectionString *string `json:"connection_string,omitempty"`
	SASToken         *string `json:"sas_token,omitempty"`
	Recursive        bool    `json:"recursive"`
}

// BoxSourceConnectorConfig represents the configuration for a Box source connector.
// It contains Box app configuration and file access settings.
type BoxSourceConnectorConfig struct {
	sourceconfig

	BoxAppConfig string `json:"box_app_config"`
	Recursive    bool   `json:"recursive"`
}

// ConfluenceSourceConnectorConfig represents the configuration for a Confluence source connector.
// It contains authentication details and content extraction settings.
type ConfluenceSourceConnectorConfig struct {
	sourceconfig

	URL                       string   `json:"url"`
	Username                  string   `json:"username"`
	Password                  *string  `json:"password,omitempty"`
	APIToken                  *string  `json:"api_token,omitempty"`
	Token                     *string  `json:"token,omitempty"`
	Cloud                     bool     `json:"cloud"`
	ExtractImages             *bool    `json:"extract_images,omitempty"`
	ExtractFiles              *bool    `json:"extract_files,omitempty"`
	MaxNumOfSpaces            int      `json:"max_num_of_spaces"`
	MaxNumOfDocsFromEachSpace int      `json:"max_num_of_docs_from_each_space"`
	Spaces                    []string `json:"spaces"`
}

// CouchbaseSourceConnectorConfig represents the configuration for a Couchbase source connector.
// It contains connection details, bucket information, and authentication credentials.
type CouchbaseSourceConnectorConfig struct {
	sourceconfig

	Bucket           string  `json:"bucket"`
	ConnectionString string  `json:"connection_string"`
	Scope            *string `json:"scope,omitempty"`
	Collection       *string `json:"collection,omitempty"`
	BatchSize        int     `json:"batch_size"`
	Username         string  `json:"username"`
	Password         string  `json:"password"`
	CollectionID     string  `json:"collection_id"`
}

// JiraSourceConnectorConfig represents the configuration for a Jira source connector.
// It contains authentication details and project/issue filtering settings.
type JiraSourceConnectorConfig struct {
	sourceconfig

	URL                 string   `json:"url"`
	Username            string   `json:"username"`
	Password            *string  `json:"password,omitempty"`
	Token               *string  `json:"token,omitempty"`
	Cloud               *bool    `json:"cloud,omitempty"`
	Projects            []string `json:"projects,omitempty"`
	Boards              []string `json:"boards,omitempty"`
	Issues              []string `json:"issues,omitempty"`
	StatusFilters       []string `json:"status_filters,omitempty"`
	DownloadAttachments *bool    `json:"download_attachments,omitempty"`
}

// PostgresSourceConnectorConfig represents the configuration for a PostgreSQL source connector.
// It contains database connection details and table configuration.
type PostgresSourceConnectorConfig struct {
	sourceconfig

	Host      string   `json:"host"`
	Database  string   `json:"database"`
	Port      int      `json:"port"`
	Username  string   `json:"username"`
	Password  string   `json:"password"`
	TableName string   `json:"table_name"`
	BatchSize int      `json:"batch_size"`
	IDColumn  string   `json:"id_column"`
	Fields    []string `json:"fields"`
}

// S3SourceConnectorConfig represents the configuration for an Amazon S3 source connector.
// It supports both AWS S3 and S3-compatible storage services for ingesting data.
type S3SourceConnectorConfig struct {
	sourceconfig

	RemoteURL   string  `json:"remote_url"`
	Anonymous   bool    `json:"anonymous"`
	Key         *string `json:"key,omitempty"`
	Secret      *string `json:"secret,omitempty"`
	Token       *string `json:"token,omitempty"`
	EndpointURL *string `json:"endpoint_url,omitempty"`
	Recursive   bool    `json:"recursive"`
}

// SharePointSourceConnectorConfig represents the configuration for a SharePoint source connector.
// It contains Microsoft Graph API authentication and site access details.
type SharePointSourceConnectorConfig struct {
	sourceconfig

	Site         string  `json:"site"`
	Tenant       string  `json:"tenant"`
	AuthorityURL *string `json:"authority_url,omitempty"`
	UserPName    string  `json:"user_pname"`
	ClientID     string  `json:"client_id"`
	ClientCred   string  `json:"client_cred"`
	Recursive    bool    `json:"recursive"`
	Path         *string `json:"path,omitempty"`
}

// SnowflakeSourceConnectorConfig represents the configuration for a Snowflake source connector.
// It contains account details, authentication, and table configuration.
type SnowflakeSourceConnectorConfig struct {
	sourceconfig

	Account   string   `json:"account"`
	Role      string   `json:"role"`
	User      string   `json:"user"`
	Password  string   `json:"password"`
	Host      string   `json:"host"`
	Port      *int     `json:"port,omitempty"`
	Database  string   `json:"database"`
	Schema    *string  `json:"schema,omitempty"`
	TableName *string  `json:"table_name,omitempty"`
	BatchSize *int     `json:"batch_size,omitempty"`
	IDColumn  *string  `json:"id_column,omitempty"`
	Fields    []string `json:"fields,omitempty"`
}

// DropboxSourceConnectorConfig represents the configuration for a Dropbox source connector.
// It contains access token and file path configuration.
type DropboxSourceConnectorConfig struct {
	sourceconfig

	Token     string `json:"token"`
	RemoteURL string `json:"remote_url"`
	Recursive bool   `json:"recursive"`
}

// GCSSourceConnectorConfig represents the configuration for a Google Cloud Storage source connector.
// It contains the remote URL and service account key for authentication.
type GCSSourceConnectorConfig struct {
	sourceconfig

	RemoteURL         string `json:"remote_url"`
	ServiceAccountKey string `json:"service_account_key"`
	Recursive         bool   `json:"recursive"`
}

// GoogleDriveSourceConnectorConfig represents the configuration for a Google Drive source connector.
// It contains drive ID, service account key, and file filtering settings.
type GoogleDriveSourceConnectorConfig struct {
	sourceconfig

	DriveID           string   `json:"drive_id"`
	ServiceAccountKey string   `json:"service_account_key"`
	Extensions        []string `json:"extensions,omitempty"`
	Recursive         bool     `json:"recursive"`
}

// KafkaCloudSourceConnectorConfig represents the configuration for a Kafka Cloud source connector.
// It contains broker details, topic information, and authentication credentials.
type KafkaCloudSourceConnectorConfig struct {
	sourceconfig

	BootstrapServers     string  `json:"bootstrap_servers"`
	Port                 int     `json:"port"`
	GroupID              *string `json:"group_id,omitempty"`
	Topic                string  `json:"topic"`
	KafkaAPIKey          string  `json:"kafka_api_key"`
	Secret               string  `json:"secret"`
	NumMessagesToConsume int     `json:"num_messages_to_consume"`
}

// OneDriveSourceConnectorConfig represents the configuration for a OneDrive source connector.
// It contains Microsoft Graph API authentication and file access settings.
type OneDriveSourceConnectorConfig struct {
	sourceconfig

	ClientID     string `json:"client_id"`
	UserPName    string `json:"user_pname"`
	Tenant       string `json:"tenant"`
	AuthorityURL string `json:"authority_url"`
	ClientCred   string `json:"client_cred"`
	Recursive    bool   `json:"recursive"`
	Path         string `json:"path"`
}

// OutlookSourceConnectorConfig represents the configuration for an Outlook source connector.
// It contains Microsoft Graph API authentication and email folder settings.
type OutlookSourceConnectorConfig struct {
	sourceconfig

	AuthorityURL   *string  `json:"authority_url,omitempty"`
	Tenant         *string  `json:"tenant,omitempty"`
	ClientID       string   `json:"client_id"`
	ClientCred     string   `json:"client_cred"`
	OutlookFolders []string `json:"outlook_folders,omitempty"`
	Recursive      bool     `json:"recursive"`
	UserEmail      string   `json:"user_email"`
}

// SalesforceSourceConnectorConfig represents the configuration for a Salesforce source connector.
// It contains authentication details and data category filtering.
type SalesforceSourceConnectorConfig struct {
	sourceconfig

	Username    string   `json:"username"`
	ConsumerKey string   `json:"consumer_key"`
	PrivateKey  string   `json:"private_key"`
	Categories  []string `json:"categories"`
}

// SlackSourceConnectorConfig represents the configuration for a Slack source connector.
// It contains channel selection, date range filtering, and authentication token.
type SlackSourceConnectorConfig struct {
	sourceconfig

	Channels  []string `json:"channels"`
	StartDate *string  `json:"start_date,omitempty"`
	EndDate   *string  `json:"end_date,omitempty"`
	Token     string   `json:"token"`
}

// ZendeskSourceConnectorConfig represents the configuration for a Zendesk source connector.
// It contains subdomain, authentication, and item type filtering.
type ZendeskSourceConnectorConfig struct {
	sourceconfig

	Subdomain string  `json:"subdomain"`
	Email     string  `json:"email"`
	APIToken  string  `json:"api_token"`
	ItemType  *string `json:"item_type,omitempty"`
	BatchSize *int    `json:"batch_size,omitempty"`
}
