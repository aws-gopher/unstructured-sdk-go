package unstructured

import (
	"encoding/json"
	"fmt"
)

// EmbedderSubtype is a type that represents an embedder subtype.
type EmbedderSubtype string

// EmbedderSubtype constants.
const (
	EmbedderSubtypeAzureOpenAI EmbedderSubtype = "azure_openai"
	EmbedderSubtypeBedrock     EmbedderSubtype = "bedrock"
	EmbedderSubtypeTogetherAI  EmbedderSubtype = "togetherai"
	EmbedderSubtypeVoyageAI    EmbedderSubtype = "voyageai"
)

// EmbedderModel is a type that represents an embedder model.
type EmbedderModel string

// EmbedderModel constants for Azure OpenAI.
const (
	EmbedderModelAzureOpenAITextEmbedding3Small EmbedderModel = "text-embedding-3-small"
	EmbedderModelAzureOpenAITextEmbedding3Large EmbedderModel = "text-embedding-3-large"
	EmbedderModelAzureOpenAITextEmbeddingAda002 EmbedderModel = "text-embedding-ada-002"
)

// EmbedderModel constants for Bedrock.
const (
	EmbedderModelBedrockTitanEmbedTextV2        EmbedderModel = "amazon.titan-embed-text-v2:0"
	EmbedderModelBedrockTitanEmbedTextV1        EmbedderModel = "amazon.titan-embed-text-v1"
	EmbedderModelBedrockTitanEmbedImageV1       EmbedderModel = "amazon.titan-embed-image-v1"
	EmbedderModelBedrockCohereEmbedEnglish      EmbedderModel = "cohere.embed-english-v3"
	EmbedderModelBedrockCohereEmbedMultilingual EmbedderModel = "cohere.embed-multilingual-v3"
)

// EmbedderModel constants for TogetherAI.
const (
	EmbedderModelTogetherAIM2Bert80M32kRetrieval EmbedderModel = "togethercomputer/m2-bert-80M-32k-retrieval"
)

// EmbedderModel constants for VoyageAI.
const (
	EmbedderModelVoyageAI3           EmbedderModel = "voyage-3"
	EmbedderModelVoyageAI3Large      EmbedderModel = "voyage-3-large"
	EmbedderModelVoyageAI3Lite       EmbedderModel = "voyage-3-lite"
	EmbedderModelVoyageAICode3       EmbedderModel = "voyage-code-3"
	EmbedderModelVoyageAIFinance2    EmbedderModel = "voyage-finance-2"
	EmbedderModelVoyageAILaw2        EmbedderModel = "voyage-law-2"
	EmbedderModelVoyageAICode2       EmbedderModel = "voyage-code-2"
	EmbedderModelVoyageAIMultimodal3 EmbedderModel = "voyage-multimodal-3"
)

// Embedder represents an embedding node in a workflow.
type Embedder struct {
	ID        string          `json:"-"`
	Name      string          `json:"-"`
	Subtype   EmbedderSubtype `json:"-"`
	ModelName EmbedderModel   `json:"model_name"`
}

var _ WorkflowNode = new(Embedder)

// isNode implements the WorkflowNode interface.
func (e Embedder) isNode() {}

// MarshalJSON implements the json.Marshaler interface.
func (e Embedder) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(struct {
		ModelName EmbedderModel `json:"model_name"`
	}{
		ModelName: e.ModelName,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal embedder settings: %w", err)
	}

	header, err := json.Marshal(header{
		ID:       e.ID,
		Name:     e.Name,
		Type:     nodeTypeEmbed,
		Subtype:  string(e.Subtype),
		Settings: json.RawMessage(data),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal embedder header: %w", err)
	}

	return header, nil
}

// ValidateModel validates that the model is compatible with the subtype.
func (e *Embedder) ValidateModel() error {
	switch e.Subtype {
	case EmbedderSubtypeAzureOpenAI:
		switch e.ModelName {
		case EmbedderModelAzureOpenAITextEmbedding3Small,
			EmbedderModelAzureOpenAITextEmbedding3Large,
			EmbedderModelAzureOpenAITextEmbeddingAda002:
			return nil
		default:
			return fmt.Errorf("invalid model %s for Azure OpenAI embedder", e.ModelName)
		}

	case EmbedderSubtypeBedrock:
		switch e.ModelName {
		case EmbedderModelBedrockTitanEmbedTextV2,
			EmbedderModelBedrockTitanEmbedTextV1,
			EmbedderModelBedrockTitanEmbedImageV1,
			EmbedderModelBedrockCohereEmbedEnglish,
			EmbedderModelBedrockCohereEmbedMultilingual:
			return nil
		default:
			return fmt.Errorf("invalid model %s for Bedrock embedder", e.ModelName)
		}

	case EmbedderSubtypeTogetherAI:
		switch e.ModelName {
		case EmbedderModelTogetherAIM2Bert80M32kRetrieval:
			return nil
		default:
			return fmt.Errorf("invalid model %s for TogetherAI embedder", e.ModelName)
		}

	case EmbedderSubtypeVoyageAI:
		switch e.ModelName {
		case EmbedderModelVoyageAI3,
			EmbedderModelVoyageAI3Large,
			EmbedderModelVoyageAI3Lite,
			EmbedderModelVoyageAICode3,
			EmbedderModelVoyageAIFinance2,
			EmbedderModelVoyageAILaw2,
			EmbedderModelVoyageAICode2,
			EmbedderModelVoyageAIMultimodal3:
			return nil
		default:
			return fmt.Errorf("invalid model %s for VoyageAI embedder", e.ModelName)
		}

	default:
		return fmt.Errorf("unknown embedder subtype: %s", e.Subtype)
	}
}

func unmarshalEmbedder(header header) (WorkflowNode, error) {
	embedder := &Embedder{
		ID:      header.ID,
		Name:    header.Name,
		Subtype: EmbedderSubtype(header.Subtype),
	}

	if err := json.Unmarshal(header.Settings, embedder); err != nil {
		return nil, fmt.Errorf("failed to unmarshal embedder node: %w", err)
	}

	if err := embedder.ValidateModel(); err != nil {
		return nil, fmt.Errorf("invalid embedder configuration: %w", err)
	}

	return embedder, nil
}
