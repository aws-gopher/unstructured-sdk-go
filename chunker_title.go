package unstructured

import (
	"encoding/json"
	"fmt"
)

// ChunkerTitle is a node that chunks text by character.
type ChunkerTitle struct {
	ID                  string `json:"-"`
	Name                string `json:"-"`
	APIURL              string `json:"unstructured_api_url,omitempty"`
	APIKey              string `json:"unstructured_api_key,omitempty"`
	CombineTextUnderN   int    `json:"combine_text_under_n_chars,omitempty"`
	IncludeOrigElements bool   `json:"include_orig_elements,omitempty"`
	NewAfterNChars      int    `json:"new_after_n_chars,omitempty"`
	MaxCharacters       int    `json:"max_characters,omitempty"`
	Overlap             int    `json:"overlap,omitempty"`
	OverlapAll          bool   `json:"overlap_all"`
}

var _ WorkflowNode = new(ChunkerTitle)

// isNode implements the WorkflowNode interface.
func (c ChunkerTitle) isNode() {}

// MarshalJSON implements the json.Marshaler interface.
func (c ChunkerTitle) MarshalJSON() ([]byte, error) {
	type alias ChunkerTitle

	data, err := json.Marshal(struct {
		alias
		ContextualChunkingStrategy string `json:"contextual_chunking_strategy"`
	}{
		alias:                      alias(c),
		ContextualChunkingStrategy: "v1",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal chunker title: %w", err)
	}

	headerData, err := json.Marshal(header{
		ID:       c.ID,
		Name:     c.Name,
		Type:     nodeTypeChunk,
		Subtype:  string(ChunkerSubtypeTitle),
		Settings: json.RawMessage(data),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal chunker title header: %w", err)
	}

	return headerData, nil
}
