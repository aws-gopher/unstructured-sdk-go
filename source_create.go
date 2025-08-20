package unstructured

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// CreateSourceRequest represents a request to create a new source connector.
// It contains the name and configuration for the source.
type CreateSourceRequest struct {
	Name   string
	Config SourceConfigInput
}

// SourceConfigInput is an interface that all source connector configurations must implement.
// It provides a way to identify the type of source connector and marshal its configuration.
type SourceConfigInput interface {
	isSourceConfigInput()
	Type() string
}

// CreateSource creates a new source connector with the specified configuration.
// It returns the created source connector with its assigned ID and metadata.
func (c *Client) CreateSource(ctx context.Context, in CreateSourceRequest) (*Source, error) {
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
		return nil, fmt.Errorf("failed to marshal source request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx,
		http.MethodPost,
		c.endpoint.JoinPath("/sources/").String(),
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	var source Source
	if err := c.do(req, &source); err != nil {
		return nil, fmt.Errorf("failed to create source: %w", err)
	}

	return &source, nil
}

type sourceconfiginput struct{}

func (s sourceconfiginput) isSourceConfigInput() {}

// AzureSourceConnectorConfigInput represents the configuration for an Azure Blob Storage source connector.
// It supports authentication via connection string, account key, or SAS token.
type AzureSourceConnectorConfigInput struct {
	sourceconfiginput

	RemoteURL        string  `json:"remote_url"`
	AccountName      *string `json:"account_name,omitempty"`
	AccountKey       *string `json:"account_key,omitempty"`
	ConnectionString *string `json:"connection_string,omitempty"`
	SASToken         *string `json:"sas_token,omitempty"`
	Recursive        *bool   `json:"recursive,omitempty"`
}

// Type always returns the connector type identifier for Azure Blob Storage: "azure".
func (c AzureSourceConnectorConfigInput) Type() string { return ConnectorTypeAzure }

// BoxSourceConnectorConfigInput represents the configuration for a Box source connector.
// It contains Box app configuration and file access settings.
type BoxSourceConnectorConfigInput struct {
	sourceconfiginput

	BoxAppConfig string `json:"box_app_config"`
	RemoteURL    string `json:"remote_url"`
	Recursive    *bool  `json:"recursive,omitempty"`
}

// Type always returns the connector type identifier for Box: "box".
func (c BoxSourceConnectorConfigInput) Type() string { return ConnectorTypeBox }

// ConfluenceSourceConnectorConfigInput represents the configuration for a Confluence source connector.
// It contains authentication details and content extraction settings.
type ConfluenceSourceConnectorConfigInput struct {
	sourceconfiginput

	URL                       string   `json:"url"`
	Username                  string   `json:"username"`
	Password                  *string  `json:"password,omitempty"`
	APIToken                  *string  `json:"api_token,omitempty"`
	Token                     *string  `json:"token,omitempty"`
	Cloud                     *bool    `json:"cloud,omitempty"`
	ExtractImages             *bool    `json:"extract_images,omitempty"`
	ExtractFiles              *bool    `json:"extract_files,omitempty"`
	MaxNumOfSpaces            *int     `json:"max_num_of_spaces,omitempty"`
	MaxNumOfDocsFromEachSpace *int     `json:"max_num_of_docs_from_each_space,omitempty"`
	Spaces                    []string `json:"spaces,omitempty"`
}

// Type always returns the connector type identifier for Confluence: "confluence".
func (c ConfluenceSourceConnectorConfigInput) Type() string { return ConnectorTypeConfluence }

// CouchbaseSourceConnectorConfigInput represents the configuration for a Couchbase source connector.
// It contains connection details, bucket information, and authentication credentials.
type CouchbaseSourceConnectorConfigInput struct {
	sourceconfiginput

	Bucket           string  `json:"bucket"`
	ConnectionString string  `json:"connection_string"`
	Scope            *string `json:"scope,omitempty"`
	Collection       *string `json:"collection,omitempty"`
	BatchSize        int     `json:"batch_size"`
	Username         string  `json:"username"`
	Password         string  `json:"password"`
	CollectionID     string  `json:"collection_id"`
}

// Type always returns the connector type identifier for Couchbase: "couchbase".
func (c CouchbaseSourceConnectorConfigInput) Type() string { return ConnectorTypeCouchbase }

// DropboxSourceConnectorConfigInput represents the configuration for a Dropbox source connector.
// It contains access token and file path configuration.
type DropboxSourceConnectorConfigInput struct {
	sourceconfiginput

	Token     string `json:"token"`
	RemoteURL string `json:"remote_url"`
	Recursive *bool  `json:"recursive,omitempty"`
}

// Type always returns the connector type identifier for Dropbox: "dropbox".
func (c DropboxSourceConnectorConfigInput) Type() string { return ConnectorTypeDropbox }

// GCSSourceConnectorConfigInput represents the configuration for a Google Cloud Storage source connector.
// It contains the remote URL and service account key for authentication.
type GCSSourceConnectorConfigInput struct {
	sourceconfiginput

	RemoteURL         string `json:"remote_url"`
	ServiceAccountKey string `json:"service_account_key"`
	Recursive         *bool  `json:"recursive,omitempty"`
}

// Type always returns the connector type identifier for Google Cloud Storage: "gcs".
func (c GCSSourceConnectorConfigInput) Type() string { return ConnectorTypeGCS }

// GoogleDriveSourceConnectorConfigInput represents the configuration for a Google Drive source connector.
// It contains drive ID, service account key, and file filtering settings.
type GoogleDriveSourceConnectorConfigInput struct {
	sourceconfiginput

	DriveID           string   `json:"drive_id"`
	ServiceAccountKey *string  `json:"service_account_key,omitempty"`
	Extensions        []string `json:"extensions,omitempty"`
	Recursive         *bool    `json:"recursive,omitempty"`
}

// Type always returns the connector type identifier for Google Drive: "google_drive".
func (c GoogleDriveSourceConnectorConfigInput) Type() string { return ConnectorTypeGoogleDrive }

// KafkaCloudSourceConnectorConfigInput represents the configuration for a Kafka Cloud source connector.
// It contains broker details, topic information, and authentication credentials.
type KafkaCloudSourceConnectorConfigInput struct {
	sourceconfiginput

	BootstrapServers     string  `json:"bootstrap_servers"`
	Port                 *int    `json:"port,omitempty"`
	GroupID              *string `json:"group_id,omitempty"`
	Topic                string  `json:"topic"`
	KafkaAPIKey          string  `json:"kafka_api_key"`
	Secret               string  `json:"secret"`
	NumMessagesToConsume *int    `json:"num_messages_to_consume,omitempty"`
}

// Type always returns the connector type identifier for Kafka Cloud: "kafka-cloud".
func (c KafkaCloudSourceConnectorConfigInput) Type() string { return ConnectorTypeKafkaCloud }

// OneDriveSourceConnectorConfigInput represents the configuration for a OneDrive source connector.
// It contains Microsoft Graph API authentication and file access settings.
type OneDriveSourceConnectorConfigInput struct {
	sourceconfiginput

	ClientID     string `json:"client_id"`
	UserPName    string `json:"user_pname"`
	Tenant       string `json:"tenant"`
	AuthorityURL string `json:"authority_url"`
	ClientCred   string `json:"client_cred"`
	Recursive    *bool  `json:"recursive,omitempty"`
	Path         string `json:"path"`
}

// Type always returns the connector type identifier for OneDrive: "onedrive".
func (c OneDriveSourceConnectorConfigInput) Type() string { return ConnectorTypeOneDrive }

// OutlookSourceConnectorConfigInput represents the configuration for an Outlook source connector.
// It contains Microsoft Graph API authentication and email folder settings.
type OutlookSourceConnectorConfigInput struct {
	sourceconfiginput

	AuthorityURL   *string  `json:"authority_url,omitempty"`
	Tenant         *string  `json:"tenant,omitempty"`
	ClientID       string   `json:"client_id"`
	ClientCred     string   `json:"client_cred"`
	OutlookFolders []string `json:"outlook_folders,omitempty"`
	Recursive      *bool    `json:"recursive,omitempty"`
	UserEmail      string   `json:"user_email"`
}

// Type always returns the connector type identifier for Outlook: "outlook".
func (c OutlookSourceConnectorConfigInput) Type() string { return ConnectorTypeOutlook }

// PostgresSourceConnectorConfigInput represents the configuration for a PostgreSQL source connector.
// It contains database connection details and table configuration.
type PostgresSourceConnectorConfigInput struct {
	sourceconfiginput

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

// Type always returns the connector type identifier for PostgreSQL: "postgres".
func (c PostgresSourceConnectorConfigInput) Type() string { return ConnectorTypePostgres }

// S3SourceConnectorConfigInput represents the configuration for an Amazon S3 source connector.
// It supports both AWS S3 and S3-compatible storage services.
type S3SourceConnectorConfigInput struct {
	sourceconfiginput

	RemoteURL   string  `json:"remote_url"`
	Anonymous   *bool   `json:"anonymous,omitempty"`
	Key         *string `json:"key,omitempty"`
	Secret      *string `json:"secret,omitempty"`
	Token       *string `json:"token,omitempty"`
	EndpointURL *string `json:"endpoint_url,omitempty"`
	Recursive   *bool   `json:"recursive,omitempty"`
}

// Type always returns the connector type identifier for S3: "s3".
func (c S3SourceConnectorConfigInput) Type() string { return ConnectorTypeS3 }

// SalesforceSourceConnectorConfigInput represents the configuration for a Salesforce source connector.
// It contains authentication details and data category filtering.
type SalesforceSourceConnectorConfigInput struct {
	sourceconfiginput

	Username    string   `json:"username"`
	ConsumerKey string   `json:"consumer_key"`
	PrivateKey  string   `json:"private_key"`
	Categories  []string `json:"categories"`
}

// Type always returns the connector type identifier for Salesforce: "salesforce".
func (c SalesforceSourceConnectorConfigInput) Type() string { return ConnectorTypeSalesforce }

// SharePointSourceConnectorConfigInput represents the configuration for a SharePoint source connector.
// It contains Microsoft Graph API authentication and site access details.
type SharePointSourceConnectorConfigInput struct {
	sourceconfiginput

	Site         string  `json:"site"`
	Tenant       string  `json:"tenant"`
	AuthorityURL *string `json:"authority_url,omitempty"`
	UserPName    string  `json:"user_pname"`
	ClientID     string  `json:"client_id"`
	ClientCred   string  `json:"client_cred"`
	Recursive    *bool   `json:"recursive,omitempty"`
	Path         *string `json:"path,omitempty"`
}

// Type always returns the connector type identifier for SharePoint: "sharepoint".
func (c SharePointSourceConnectorConfigInput) Type() string { return ConnectorTypeSharePoint }

// SlackSourceConnectorConfigInput represents the configuration for a Slack source connector.
// It contains channel selection, date range filtering, and authentication token.
type SlackSourceConnectorConfigInput struct {
	sourceconfiginput

	Channels  []string `json:"channels"`
	StartDate *string  `json:"start_date,omitempty"`
	EndDate   *string  `json:"end_date,omitempty"`
	Token     string   `json:"token"`
}

// Type always returns the connector type identifier for Slack: "slack".
func (c SlackSourceConnectorConfigInput) Type() string { return ConnectorTypeSlack }

// SnowflakeSourceConnectorConfigInput represents the configuration for a Snowflake source connector.
// It contains account details, authentication, and table configuration.
type SnowflakeSourceConnectorConfigInput struct {
	sourceconfiginput

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

// Type always returns the connector type identifier for Snowflake: "snowflake".
func (c SnowflakeSourceConnectorConfigInput) Type() string { return ConnectorTypeSnowflake }

// JiraSourceConnectorConfigInput represents the configuration for a Jira source connector.
// It contains authentication details and project/issue filtering settings.
type JiraSourceConnectorConfigInput struct {
	sourceconfiginput

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

// Type always returns the connector type identifier for Jira: "jira".
func (c JiraSourceConnectorConfigInput) Type() string { return ConnectorTypeJira }

// ZendeskSourceConnectorConfigInput represents the configuration for a Zendesk source connector.
// It contains subdomain, authentication, and item type filtering.
type ZendeskSourceConnectorConfigInput struct {
	sourceconfiginput

	Subdomain string  `json:"subdomain"`
	Email     string  `json:"email"`
	APIToken  string  `json:"api_token"`
	ItemType  *string `json:"item_type,omitempty"`
	BatchSize *int    `json:"batch_size,omitempty"`
}

// Type always returns the connector type identifier for Zendesk: "zendesk".
func (c ZendeskSourceConnectorConfigInput) Type() string { return ConnectorTypeZendesk }
