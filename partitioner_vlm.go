package unstructured

import (
	"encoding/json"
	"fmt"
)

// PartitionerVLM is a partitioner that uses the VLM strategy.
type PartitionerVLM struct {
	ID             string       `json:"-"`
	Name           string       `json:"-"`
	Strategy       string       `json:"strategy,omitempty"`
	Provider       Provider     `json:"provider,omitempty"`
	ProviderAPIKey string       `json:"provider_api_key,omitempty"`
	Model          Model        `json:"model,omitempty"`
	OutputFormat   OutputFormat `json:"output_format,omitempty"`
	Prompt         struct {
		Text string `json:"text,omitempty"`
	} `json:"prompt,omitzero"`
	FormatHTML       *bool `json:"format_html,omitzero"`
	UniqueElementIDs *bool `json:"unique_element_ids,omitzero"`
	IsDynamic        *bool `json:"is_dynamic,omitzero"`
	AllowFast        *bool `json:"allow_fast,omitzero"`
}

var _ WorkflowNode = new(PartitionerVLM)

// MarshalJSON implements the json.Marshaler interface.
func (p PartitionerVLM) MarshalJSON() ([]byte, error) {
	type alias PartitionerVLM

	data, err := json.Marshal(alias(p))
	if err != nil {
		return nil, fmt.Errorf("failed to marshal partitioner vlm: %w", err)
	}

	headerData, err := json.Marshal(header{
		ID:       p.ID,
		Name:     p.Name,
		Type:     PartitionerStrategyVLM,
		Subtype:  string(nodeTypePartition),
		Settings: json.RawMessage(data),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal partitioner vlm header: %w", err)
	}

	return headerData, nil
}

func (p *PartitionerVLM) isNode() {}
