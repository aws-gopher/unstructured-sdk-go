package unstructured

import (
	"encoding/json"
	"errors"
	"fmt"
)

// WorkflowNodes is a slice of WorkflowNode.
type WorkflowNodes []WorkflowNode

// ValidateNodeOrder validates the order of nodes in a workflow.
func (w WorkflowNodes) ValidateNodeOrder() (err error) {
	if len(w) == 0 {
		return errors.New("first node must be a partitioner")
	}

	// you have to partition.
	switch w[0].(type) {
	case *PartitionerAuto, *PartitionerVLM, *PartitionerHiRes, *PartitionerFast:
		// good
	default:
		err = errors.Join(err, errors.New("first node must be a partitioner"))
	}

	last := nodeTypePartition

	var (
		didEnrichTable bool
		didEnrichNER   bool
		didEnrichImage bool
	)

	for i, node := range w[1:] {
		switch node := node.(type) {
		case *PartitionerAuto, *PartitionerVLM, *PartitionerHiRes, *PartitionerFast:
			err = errors.Join(err, errors.New("only the first node may be a partitioner"))

		case *ChunkerCharacter, *ChunkerTitle, *ChunkerPage, *ChunkerSimilarity:
			// you can chunk after you partition.
			if last != nodeTypePartition && last != nodeTypeEnrich {
				err = errors.Join(err, fmt.Errorf("%s must be after %s or %s", nodeTypeChunk, nodeTypePartition, nodeTypeEnrich))
			}

			last = nodeTypeChunk

		case *Embedder:
			// you can embed after you chunk.
			if last != nodeTypeChunk {
				err = errors.Join(err, fmt.Errorf("%s must be after %s", nodeTypeEmbed, nodeTypeChunk))
			}

			last = nodeTypeEmbed

		case *Enricher:
			// you can enrich before you chunk...
			if i == len(w[1:])-1 {
				err = errors.Join(err, fmt.Errorf("%s must not be the last node", nodeTypeEnrich))
			}

			// and after you partition or enrich.
			if last != nodeTypePartition && last != nodeTypeEnrich {
				err = errors.Join(err, fmt.Errorf("%s must be after %s or %s", nodeTypeEnrich, nodeTypePartition, nodeTypeEnrich))
			}

			// you can only have one image enrichment.
			if node.isImage() && didEnrichImage {
				err = errors.Join(err, errors.New("only one image enrichment is allowed"))
			}

			didEnrichImage = node.isImage()

			// you can only have one table enrichment.
			if node.isTable() && didEnrichTable {
				err = errors.Join(err, errors.New("only one table enrichment is allowed"))
			}

			didEnrichTable = node.isTable()

			// you can only have one NER enrichment.
			if node.isNER() && didEnrichNER {
				err = errors.Join(err, errors.New("only one NER enrichment is allowed"))
			}

			didEnrichNER = node.isNER()

			last = nodeTypeEnrich

		default:
			err = errors.Join(err, fmt.Errorf("invalid node type %T at index %d", node, i+1))
		}
	}

	return err
}

// MarshalJSON implements the json.Marshaler interface.
func (w WorkflowNodes) MarshalJSON() ([]byte, error) {
	nodes := make([]json.RawMessage, len(w))

	for i, node := range w {
		msg, err := json.Marshal(node)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal workflow node: %w", err)
		}

		nodes[i] = msg
	}

	headerData, err := json.Marshal(nodes)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal workflow nodes: %w", err)
	}

	return headerData, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (w *WorkflowNodes) UnmarshalJSON(data []byte) error {
	var nodes []json.RawMessage
	if err := json.Unmarshal(data, &nodes); err != nil {
		return fmt.Errorf("failed to unmarshal workflow nodes: %w", err)
	}

	if cap(*w) < len(nodes) {
		*w = make(WorkflowNodes, 0, len(nodes))
	}

	for _, node := range nodes {
		val, err := unmarshalNode(node)
		if err != nil {
			return fmt.Errorf("failed to unmarshal workflow node: %w", err)
		}

		*w = append(*w, val)
	}

	return nil
}

// WorkflowNode is a node in a workflow.
type WorkflowNode interface {
	json.Marshaler
	isNode()
}

type header struct {
	ID       string          `json:"id,omitempty"`
	Name     string          `json:"name"`
	Type     string          `json:"type"`
	Subtype  string          `json:"subtype"`
	Settings json.RawMessage `json:"settings"`
}

const (
	nodeTypePartition = "partition"
	nodeTypeEnrich    = "prompter"
	nodeTypeChunk     = "chunk"
	nodeTypeEmbed     = "embed"
)

func unmarshalNode(data []byte) (WorkflowNode, error) {
	var header header
	if err := json.Unmarshal(data, &header); err != nil {
		return nil, fmt.Errorf("failed to unmarshal workflow node: %w", err)
	}

	switch header.Type {
	case nodeTypePartition:
		return unmarshalPartitioner(header)

	case nodeTypeChunk:
		return unmarshalChunker(header)

	case nodeTypeEmbed:
		return unmarshalEmbedder(header)

	case nodeTypeEnrich:
		return unmarshalEnricher(header)
	}

	return nil, fmt.Errorf("unknown node type: %s", header.Type)
}
