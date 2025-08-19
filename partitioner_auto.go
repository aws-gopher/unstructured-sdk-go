package unstructured

import (
	"encoding/json"
	"fmt"
)

// PartitionerAuto is a partitioner that uses the Auto strategy.
type PartitionerAuto struct {
	ID             string       `json:"-"`
	Name           string       `json:"-"`
	Strategy       string       `json:"strategy"`
	Provider       Provider     `json:"provider,omitempty"`
	ProviderAPIKey string       `json:"provider_api_key,omitempty"`
	Model          Model        `json:"model,omitempty"`
	OutputFormat   OutputFormat `json:"output_format,omitempty"`
	Prompt         struct {
		Text string `json:"text,omitempty"`
	} `json:"prompt,omitzero"`
	FormatHTML       *bool `json:"format_html,omitzero"`
	UniqueElementIDs *bool `json:"unique_element_ids,omitzero"`
	IsDynamic        bool  `json:"is_dynamic"`
	AllowFast        bool  `json:"allow_fast"`
}

var _ WorkflowNode = new(PartitionerAuto)

// MarshalJSON implements the json.Marshaler interface.
func (p PartitionerAuto) MarshalJSON() ([]byte, error) {
	type alias PartitionerAuto

	data, err := json.Marshal(struct {
		alias
		Strategy string `json:"strategy"`
	}{
		alias:    alias(p),
		Strategy: PartitionerStrategyAuto,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal partitioner auto: %w", err)
	}

	headerData, err := json.Marshal(header{
		ID:       p.ID,
		Name:     p.Name,
		Type:     nodeTypePartition,
		Subtype:  PartitionerStrategyVLM,
		Settings: json.RawMessage(data),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal partitioner auto header: %w", err)
	}

	return headerData, nil
}

func (p *PartitionerAuto) isNode() {}
