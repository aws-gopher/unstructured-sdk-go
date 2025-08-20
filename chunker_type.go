package unstructured

import (
	"encoding/json"
	"fmt"
)

// ChunkerSubtype is a type that represents a chunker subtype.
type ChunkerSubtype string

// ChunkerSubtype constants.
const (
	ChunkerSubtypeCharacter  ChunkerSubtype = "chunk_by_character"
	ChunkerSubtypeTitle      ChunkerSubtype = "chunk_by_title"
	ChunkerSubtypePage       ChunkerSubtype = "chunk_by_page"
	ChunkerSubtypeSimilarity ChunkerSubtype = "chunk_by_similarity"
)

func unmarshalChunker(header header) (WorkflowNode, error) {
	var chunker WorkflowNode

	switch ChunkerSubtype(header.Subtype) {
	case ChunkerSubtypeCharacter:
		chunker = &ChunkerCharacter{
			ID:   header.ID,
			Name: header.Name,
		}

	case ChunkerSubtypeTitle:
		chunker = &ChunkerTitle{
			ID:   header.ID,
			Name: header.Name,
		}

	case ChunkerSubtypePage:
		chunker = &ChunkerPage{
			ID:   header.ID,
			Name: header.Name,
		}

	case ChunkerSubtypeSimilarity:
		chunker = &ChunkerSimilarity{
			ID:   header.ID,
			Name: header.Name,
		}

	default:
		return nil, fmt.Errorf("unknown Chunker strategy: %s", header.Subtype)
	}

	if err := json.Unmarshal(header.Settings, chunker); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Chunker node: %w", err)
	}

	return chunker, nil
}
