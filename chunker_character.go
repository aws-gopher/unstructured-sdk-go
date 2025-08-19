package unstructured

import (
	"encoding/json"
	"fmt"
)

// ChunkerCharacter is a node that chunks text by character.
type ChunkerCharacter struct {
	ID                         string           `json:"-"`
	Name                       string           `json:"-"`
	APIURL                     string           `json:"unstructured_api_url,omitempty"`
	APIKey                     string           `json:"unstructured_api_key,omitempty"`
	IncludeOrigElements        bool             `json:"include_orig_elements,omitempty"`
	NewAfterNChars             int              `json:"new_after_n_chars,omitempty"`
	MaxCharacters              int              `json:"max_characters,omitempty"`
	Overlap                    int              `json:"overlap,omitempty"`
	OverlapAll                 bool             `json:"overlap_all"`
	ContextualChunkingStrategy ChunkingStrategy `json:"contextual_chunking_strategy,omitempty"`
}

// ChunkingStrategy is a strategy for contextual chunking.
type ChunkingStrategy string

// ChunkingStrategyV1 is a strategy for contextual chunking.
const ChunkingStrategyV1 = "v1"

var _ WorkflowNode = new(ChunkerCharacter)

// isNode implements the WorkflowNode interface.
func (c ChunkerCharacter) isNode() {}

// MarshalJSON implements the json.Marshaler interface.
func (c ChunkerCharacter) MarshalJSON() ([]byte, error) {
	type alias ChunkerCharacter

	data, err := json.Marshal(alias(c))
	if err != nil {
		return nil, fmt.Errorf("failed to marshal chunker character: %w", err)
	}

	headerData, err := json.Marshal(header{
		ID:       c.ID,
		Name:     c.Name,
		Type:     nodeTypeChunk,
		Subtype:  string(ChunkerSubtypeCharacter),
		Settings: json.RawMessage(data),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal chunker character header: %w", err)
	}

	return headerData, nil
}
