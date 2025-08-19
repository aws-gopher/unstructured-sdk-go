package unstructured

import (
	"encoding/json"
	"fmt"
)

// PartitionerHiRes represents a high-resolution partitioner configuration for document processing.
type PartitionerHiRes struct {
	ID                     string               `json:"-"`
	Name                   string               `json:"-"`
	PageBreaks             bool                 `json:"include_page_breaks,omitzero"`
	PDFInferTableStructure bool                 `json:"pdf_infer_table_structure,omitzero"`
	ExcludeElements        []ExcludeableElement `json:"exclude_elements,omitzero"`
	XMLKeepTags            bool                 `json:"xml_keep_tags,omitzero"`
	Encoding               Encoding             `json:"encoding,omitzero"`
	OCRLanguages           []Language           `json:"ocr_languages,omitzero"`
	ExtractImageBlockTypes []BlockType          `json:"extract_image_block_types,omitzero"`
	InferTableStructure    bool                 `json:"infer_table_structure,omitzero"`
}

var _ WorkflowNode = new(PartitionerHiRes)

// MarshalJSON implements the json.Marshaler interface for PartitionerHiRes.
func (p PartitionerHiRes) MarshalJSON() ([]byte, error) {
	type alias PartitionerHiRes

	mask := struct {
		Strategy string `json:"strategy"`
		alias
	}{
		Strategy: PartitionerStrategyHiRes,
		alias:    alias(p),
	}

	data, err := json.Marshal(mask)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal partitioner hi res: %w", err)
	}

	headerData, err := json.Marshal(header{
		ID:       p.ID,
		Name:     p.Name,
		Type:     nodeTypePartition,
		Subtype:  string(PartitionerStrategyHiRes),
		Settings: json.RawMessage(data),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal partitioner hi res header: %w", err)
	}

	return headerData, nil
}

func (p *PartitionerHiRes) isNode() {}
