package unstructured

import (
	"maps"
	"strings"
	"testing"
)

func TestWorkflowNodeOrder(t *testing.T) {
	t.Parallel()

	partitioners := map[string]WorkflowNode{
		"none":            nil,
		"partition_auto":  &PartitionerAuto{},
		"partition_vlm":   &PartitionerVLM{},
		"partition_hires": &PartitionerHiRes{},
		"partition_fast":  &PartitionerFast{},
	}

	chunkers := map[string]WorkflowNode{
		"none":               nil,
		"chunker_character":  &ChunkerCharacter{},
		"chunker_title":      &ChunkerTitle{},
		"chunker_page":       &ChunkerPage{},
		"chunker_similarity": &ChunkerSimilarity{},
	}

	enrichers := map[string]*Enricher{
		"none":                       nil,
		"enricher_image_openai":      {Subtype: EnrichmentTypeImageOpenAI},
		"enricher_table_openai":      {Subtype: EnrichmentTypeTableOpenAI},
		"enricher_table2html_openai": {Subtype: EnrichmentTypeTable2HTMLOpenAI},
		"enricher_ner_openai":        {Subtype: EnrichmentTypeNEROpenAI},
		"enricher_image_anthropic":   {Subtype: EnrichmentTypeImageAnthropic},
		"enricher_table_anthropic":   {Subtype: EnrichmentTypeTableAnthropic},
		"enricher_ner_anthropic":     {Subtype: EnrichmentTypeNERAnthropic},
		"enricher_image_bedrock":     {Subtype: EnrichmentTypeImageBedrock},
		"enricher_table_bedrock":     {Subtype: EnrichmentTypeTableBedrock},
	}

	embedders := map[string]WorkflowNode{
		"none":     nil,
		"embedder": &Embedder{},
	}

	type testcase struct {
		nodes   WorkflowNodes
		wantErr bool
	}

	tests := make(map[string]testcase, len(partitioners)*len(chunkers)*len(embedders)*len(enrichers)+4)

	for partitionerName, partitioner := range partitioners {
		for chunkerName, chunker := range chunkers {
			for enricherName, enricher := range enrichers {
				for embedderName, embedder := range embedders {
					labels := []string{}

					var tc testcase

					tc.wantErr = partitioner == nil

					if partitioner != nil {
						labels = append(labels, partitionerName)

						tc.nodes = append(tc.nodes, partitioner)
					}

					if enricher != nil {
						labels = append(labels, enricherName)

						tc.nodes = append(tc.nodes, enricher)

						tc.wantErr = tc.wantErr || chunker == nil
					}

					if chunker != nil {
						labels = append(labels, chunkerName)

						tc.nodes = append(tc.nodes, chunker)
					}

					if embedder != nil {
						labels = append(labels, embedderName)

						tc.nodes = append(tc.nodes, embedder)

						tc.wantErr = tc.wantErr || chunker == nil
					}

					name := strings.Join(labels, "-")
					if name == "" {
						name = "none"
					}

					tests[name] = tc
				}
			}
		}
	}

	maps.Copy(tests, map[string]testcase{
		"wrong_order": {
			nodes:   WorkflowNodes{&PartitionerAuto{}, &ChunkerCharacter{}, &Embedder{}, &Enricher{Subtype: EnrichmentTypeImageOpenAI}},
			wantErr: true,
		},
		"double_image_enricher": {
			nodes:   WorkflowNodes{&PartitionerAuto{}, &Enricher{Subtype: EnrichmentTypeImageOpenAI}, &Enricher{Subtype: EnrichmentTypeImageAnthropic}, &ChunkerCharacter{}},
			wantErr: true,
		},
		"double_table_enricher": {
			nodes:   WorkflowNodes{&PartitionerAuto{}, &Enricher{Subtype: EnrichmentTypeTableOpenAI}, &Enricher{Subtype: EnrichmentTypeTableBedrock}, &ChunkerCharacter{}},
			wantErr: true,
		},
		"double_ner_enricher": {
			nodes:   WorkflowNodes{&PartitionerAuto{}, &Enricher{Subtype: EnrichmentTypeNEROpenAI}, &Enricher{Subtype: EnrichmentTypeNERAnthropic}, &ChunkerCharacter{}},
			wantErr: true,
		},
	})

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := test.nodes.ValidateNodeOrder()

			switch {
			case !test.wantErr && got != nil:
				t.Errorf("got\n%v\nwant nil", got)

			case test.wantErr && got == nil:
				t.Errorf("got nil, want error")
			}
		})
	}
}
