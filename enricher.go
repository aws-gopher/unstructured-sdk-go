package unstructured

import (
	"encoding/json"
	"fmt"
)

// Enricher is a node that enriches text.
type Enricher struct {
	ID                string         `json:"-"`
	Name              string         `json:"-"`
	Subtype           EnrichmentType `json:"-"`
	NERPromptOverride string         `json:"prompt_interface_overrides,omitempty"`
}

// EnrichmentType is a type that represents an enrichment type.
type EnrichmentType string

// EnrichmentType constants.
const (
	EnrichmentTypeImageOpenAI      EnrichmentType = "openai_image_description"
	EnrichmentTypeTableOpenAI      EnrichmentType = "openai_table_description"
	EnrichmentTypeTable2HTMLOpenAI EnrichmentType = "openai_table2html"
	EnrichmentTypeNEROpenAI        EnrichmentType = "openai_ner"

	EnrichmentTypeImageAnthropic EnrichmentType = "anthropic_image_description"
	EnrichmentTypeTableAnthropic EnrichmentType = "anthropic_table_description"
	EnrichmentTypeNERAnthropic   EnrichmentType = "anthropic_ner"

	EnrichmentTypeImageBedrock EnrichmentType = "bedrock_image_description"
	EnrichmentTypeTableBedrock EnrichmentType = "bedrock_table_description"
)

var _ WorkflowNode = new(Enricher)

func (e Enricher) isNode() {}

// MarshalJSON implements the json.Marshaler interface.
func (e Enricher) MarshalJSON() ([]byte, error) {
	var settings json.RawMessage

	if e.NERPromptOverride != "" && (e.Subtype == EnrichmentTypeNERAnthropic || e.Subtype == EnrichmentTypeNEROpenAI) {
		nested := struct {
			PromptOverride struct {
				Prompt struct {
					User string `json:"user"`
				} `json:"prompt"`
			} `json:"prompt_interface_overrides"`
		}{}
		nested.PromptOverride.Prompt.User = e.NERPromptOverride

		data, err := json.Marshal(nested)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal enricher nested settings: %w", err)
		}

		settings = json.RawMessage(data)
	}

	headerData, err := json.Marshal(header{
		ID:       e.ID,
		Name:     e.Name,
		Type:     nodeTypeEnrich,
		Subtype:  string(e.Subtype),
		Settings: settings,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal enricher header: %w", err)
	}

	return headerData, nil
}

func unmarshalEnricher(header header) (WorkflowNode, error) {
	enricher := &Enricher{
		ID:   header.ID,
		Name: header.Name,
	}

	if err := json.Unmarshal(header.Settings, enricher); err != nil {
		return nil, fmt.Errorf("failed to unmarshal enricher: %w", err)
	}

	return enricher, nil
}
