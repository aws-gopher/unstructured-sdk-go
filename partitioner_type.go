package unstructured

import (
	"encoding/json"
	"fmt"
)

// OutputFormat represents the output format for document processing.
type OutputFormat string

// Output format constants.
const (
	OutputFormatHTML OutputFormat = "text/html"
	OutputFormatJSON OutputFormat = "application/json"
)

// Partitioner strategy constants.
const (
	PartitionerStrategyAuto  = "auto"
	PartitionerStrategyVLM   = "vlm"
	PartitionerStrategyHiRes = "hi_res"
	PartitionerStrategyFast  = "fast"
)

func unmarshalPartitioner(header header) (WorkflowNode, error) {
	var partitioner WorkflowNode

	switch header.Subtype {
	case PartitionerStrategyAuto:
		partitioner = &PartitionerAuto{
			ID:   header.ID,
			Name: header.Name,
		}

	case PartitionerStrategyVLM:
		partitioner = &PartitionerVLM{
			ID:   header.ID,
			Name: header.Name,
		}

	case PartitionerStrategyHiRes:
		partitioner = &PartitionerHiRes{
			ID:   header.ID,
			Name: header.Name,
		}

	case PartitionerStrategyFast:
		partitioner = &PartitionerFast{
			ID:   header.ID,
			Name: header.Name,
		}

	default:
		return nil, fmt.Errorf("unknown partitioner strategy: %s", header.Subtype)
	}

	if err := json.Unmarshal(header.Settings, partitioner); err != nil {
		return nil, fmt.Errorf("failed to unmarshal partitioner node: %w", err)
	}

	return partitioner, nil
}
