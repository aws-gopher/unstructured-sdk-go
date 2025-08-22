package unstructured

import (
	"encoding/json"
	"fmt"
	"time"
)

// sourceConfigFactories maps source type strings to factory functions
// that create new instances of the appropriate concrete source config type.
var sourceConfigFactories = map[string]func() SourceConfig{
	ConnectorTypeAzure:             func() SourceConfig { return new(AzureSourceConnectorConfig) },
	ConnectorTypeBox:               func() SourceConfig { return new(BoxSourceConnectorConfig) },
	ConnectorTypeConfluence:        func() SourceConfig { return new(ConfluenceSourceConnectorConfig) },
	ConnectorTypeCouchbase:         func() SourceConfig { return new(CouchbaseConnectorConfig) },
	ConnectorTypeDatabricksVolumes: func() SourceConfig { return new(DatabricksVolumesConnectorConfig) },
	ConnectorTypeDropbox:           func() SourceConfig { return new(DropboxSourceConnectorConfig) },
	ConnectorTypeElasticsearch:     func() SourceConfig { return new(ElasticsearchConnectorConfig) },
	ConnectorTypeGCS:               func() SourceConfig { return new(GCSConnectorConfig) },
	ConnectorTypeGoogleDrive:       func() SourceConfig { return new(GoogleDriveSourceConnectorConfig) },
	ConnectorTypeJira:              func() SourceConfig { return new(JiraSourceConnectorConfig) },
	ConnectorTypeKafkaCloud:        func() SourceConfig { return new(KafkaCloudConnectorConfig) },
	ConnectorTypeMongoDB:           func() SourceConfig { return new(MongoDBConnectorConfig) },
	ConnectorTypeOneDrive:          func() SourceConfig { return new(OneDriveConnectorConfig) },
	ConnectorTypeOutlook:           func() SourceConfig { return new(OutlookSourceConnectorConfig) },
	ConnectorTypePostgres:          func() SourceConfig { return new(PostgresConnectorConfig) },
	ConnectorTypeS3:                func() SourceConfig { return new(S3ConnectorConfig) },
	ConnectorTypeSalesforce:        func() SourceConfig { return new(SalesforceSourceConnectorConfig) },
	ConnectorTypeSharePoint:        func() SourceConfig { return new(SharePointSourceConnectorConfig) },
	ConnectorTypeSlack:             func() SourceConfig { return new(SlackSourceConnectorConfig) },
	ConnectorTypeSnowflake:         func() SourceConfig { return new(SnowflakeConnectorConfig) },
	ConnectorTypeZendesk:           func() SourceConfig { return new(ZendeskSourceConnectorConfig) },
}

// Source represents a source connector that ingests files or data from various locations.
// It contains metadata about the connector and its configuration.
type Source struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	CreatedAt time.Time    `json:"created_at,omitzero"`
	UpdatedAt time.Time    `json:"updated_at,omitzero"`
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
	Type() string
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
	Recursive        *bool   `json:"recursive,omitempty"`
}

var _ SourceConfig = (*AzureSourceConnectorConfig)(nil)

// Type always returns the connector type identifier for Azure: "azure".
func (c AzureSourceConnectorConfig) Type() string { return ConnectorTypeAzure }

// BoxSourceConnectorConfig represents the configuration for a Box source connector.
// It contains Box app configuration and file access settings.
type BoxSourceConnectorConfig struct {
	sourceconfig

	BoxAppConfig string `json:"box_app_config"`
	RemoteURL    string `json:"remote_url"`
	Recursive    *bool  `json:"recursive,omitempty"`
}

var _ SourceConfig = (*BoxSourceConnectorConfig)(nil)

// Type always returns the connector type identifier for Box: "box".
func (c BoxSourceConnectorConfig) Type() string { return ConnectorTypeBox }

// ConfluenceSourceConnectorConfig represents the configuration for a Confluence source connector.
// It contains authentication details and content extraction settings.
type ConfluenceSourceConnectorConfig struct {
	sourceconfig

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

var _ SourceConfig = (*ConfluenceSourceConnectorConfig)(nil)

// Type always returns the connector type identifier for Confluence: "confluence".
func (c ConfluenceSourceConnectorConfig) Type() string { return ConnectorTypeConfluence }

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

var _ SourceConfig = (*JiraSourceConnectorConfig)(nil)

// Type always returns the connector type identifier for Jira: "jira".
func (c JiraSourceConnectorConfig) Type() string { return ConnectorTypeJira }

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
	Recursive    *bool   `json:"recursive,omitempty"`
	Path         *string `json:"path,omitempty"`
}

var _ SourceConfig = (*SharePointSourceConnectorConfig)(nil)

// Type always returns the connector type identifier for SharePoint: "sharepoint".
func (c SharePointSourceConnectorConfig) Type() string { return ConnectorTypeSharePoint }

// DropboxSourceConnectorConfig represents the configuration for a Dropbox source connector.
// It contains access token and file path configuration.
type DropboxSourceConnectorConfig struct {
	sourceconfig

	Token     string `json:"token"`
	RemoteURL string `json:"remote_url"`
	Recursive *bool  `json:"recursive,omitempty"`
}

var _ SourceConfig = (*DropboxSourceConnectorConfig)(nil)

// Type always returns the connector type identifier for Dropbox: "dropbox".
func (c DropboxSourceConnectorConfig) Type() string { return ConnectorTypeDropbox }

// GoogleDriveSourceConnectorConfig represents the configuration for a Google Drive source connector.
// It contains drive ID, service account key, and file filtering settings.
type GoogleDriveSourceConnectorConfig struct {
	sourceconfig

	DriveID           string   `json:"drive_id"`
	ServiceAccountKey *string  `json:"service_account_key,omitempty"`
	Extensions        []string `json:"extensions,omitempty"`
	Recursive         *bool    `json:"recursive,omitempty"`
}

var _ SourceConfig = (*GoogleDriveSourceConnectorConfig)(nil)

// Type always returns the connector type identifier for Google Drive: "google_drive".
func (c GoogleDriveSourceConnectorConfig) Type() string { return ConnectorTypeGoogleDrive }

// OutlookSourceConnectorConfig represents the configuration for an Outlook source connector.
// It contains Microsoft Graph API authentication and email folder settings.
type OutlookSourceConnectorConfig struct {
	sourceconfig

	AuthorityURL   *string  `json:"authority_url,omitempty"`
	Tenant         *string  `json:"tenant,omitempty"`
	ClientID       string   `json:"client_id"`
	ClientCred     string   `json:"client_cred"`
	OutlookFolders []string `json:"outlook_folders,omitempty"`
	Recursive      *bool    `json:"recursive,omitempty"`
	UserEmail      string   `json:"user_email"`
}

var _ SourceConfig = (*OutlookSourceConnectorConfig)(nil)

// Type always returns the connector type identifier for Outlook: "outlook".
func (c OutlookSourceConnectorConfig) Type() string { return ConnectorTypeOutlook }

// SalesforceSourceConnectorConfig represents the configuration for a Salesforce source connector.
// It contains authentication details and data category filtering.
type SalesforceSourceConnectorConfig struct {
	sourceconfig

	Username    string   `json:"username"`
	ConsumerKey string   `json:"consumer_key"`
	PrivateKey  string   `json:"private_key"`
	Categories  []string `json:"categories"`
}

var _ SourceConfig = (*SalesforceSourceConnectorConfig)(nil)

// Type always returns the connector type identifier for Salesforce: "salesforce".
func (c SalesforceSourceConnectorConfig) Type() string { return ConnectorTypeSalesforce }

// SlackSourceConnectorConfig represents the configuration for a Slack source connector.
// It contains channel selection, date range filtering, and authentication token.
type SlackSourceConnectorConfig struct {
	sourceconfig

	Channels  []string `json:"channels"`
	StartDate *string  `json:"start_date,omitempty"`
	EndDate   *string  `json:"end_date,omitempty"`
	Token     string   `json:"token"`
}

var _ SourceConfig = (*SlackSourceConnectorConfig)(nil)

// Type always returns the connector type identifier for Slack: "slack".
func (c SlackSourceConnectorConfig) Type() string { return ConnectorTypeSlack }

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

var _ SourceConfig = (*ZendeskSourceConnectorConfig)(nil)

// Type always returns the connector type identifier for Zendesk: "zendesk".
func (c ZendeskSourceConnectorConfig) Type() string { return ConnectorTypeZendesk }
