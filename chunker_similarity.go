package unstructured

import (
	"encoding/json"
	"fmt"
)

// ChunkerSimilarity is a node that chunks text by character.
type ChunkerSimilarity struct {
	ID                  string `json:"-"`
	Name                string `json:"-"`
	APIURL              string `json:"unstructured_api_url,omitempty"`
	APIKey              string `json:"unstructured_api_key,omitempty"`
	IncludeOrigElements bool   `json:"include_orig_elements,omitempty"`
	NewAfterNChars      int    `json:"new_after_n_chars,omitempty"`
	MaxCharacters       int    `json:"max_characters,omitempty"`
	Overlap             int    `json:"overlap,omitempty"`
	OverlapAll          bool   `json:"overlap_all"`
}

var _ WorkflowNode = new(ChunkerSimilarity)

// isNode implements the WorkflowNode interface.
func (c ChunkerSimilarity) isNode() {}

// MarshalJSON implements the json.Marshaler interface.
func (c ChunkerSimilarity) MarshalJSON() ([]byte, error) {
	type alias ChunkerSimilarity

	data, err := json.Marshal(struct {
		alias
		ContextualChunkingStrategy string `json:"contextual_chunking_strategy"`
	}{
		alias:                      alias(c),
		ContextualChunkingStrategy: "v1",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal chunker similarity: %w", err)
	}

	headerData, err := json.Marshal(header{
		ID:       c.ID,
		Name:     c.Name,
		Type:     nodeTypeChunk,
		Subtype:  string(ChunkerSubtypeSimilarity),
		Settings: json.RawMessage(data),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal chunker similarity header: %w", err)
	}

	return headerData, nil
}
