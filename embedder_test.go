package unstructured

import (
	"encoding/json"
	"testing"
)

func TestEmbedder_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		embedder *Embedder
		wantErr  bool
	}{
		{
			name: "Azure OpenAI embedder",
			embedder: &Embedder{
				ID:        "test-id",
				Name:      "Test Embedder",
				Subtype:   EmbedderSubtypeAzureOpenAI,
				ModelName: EmbedderModelAzureOpenAITextEmbedding3Small,
			},
			wantErr: false,
		},
		{
			name: "Bedrock embedder",
			embedder: &Embedder{
				ID:        "test-id",
				Name:      "Test Embedder",
				Subtype:   EmbedderSubtypeBedrock,
				ModelName: EmbedderModelBedrockTitanEmbedTextV2,
			},
			wantErr: false,
		},
		{
			name: "TogetherAI embedder",
			embedder: &Embedder{
				ID:        "test-id",
				Name:      "Test Embedder",
				Subtype:   EmbedderSubtypeTogetherAI,
				ModelName: EmbedderModelTogetherAIM2Bert80M32kRetrieval,
			},
			wantErr: false,
		},
		{
			name: "VoyageAI embedder",
			embedder: &Embedder{
				ID:        "test-id",
				Name:      "Test Embedder",
				Subtype:   EmbedderSubtypeVoyageAI,
				ModelName: EmbedderModelVoyageAI3,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.embedder)
			if (err != nil) != tt.wantErr {
				t.Errorf("Embedder.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify the JSON structure
				var result map[string]interface{}
				if err := json.Unmarshal(data, &result); err != nil {
					t.Errorf("Failed to unmarshal result: %v", err)
					return
				}

				// Check required fields
				if result["type"] != "embed" {
					t.Errorf("Expected type 'embed', got %v", result["type"])
				}

				if result["subtype"] != string(tt.embedder.Subtype) {
					t.Errorf("Expected subtype %s, got %v", tt.embedder.Subtype, result["subtype"])
				}

				// Check settings
				settings, ok := result["settings"].(map[string]interface{})
				if !ok {
					t.Errorf("Settings not found or not an object")
					return
				}

				if settings["model_name"] != string(tt.embedder.ModelName) {
					t.Errorf("Expected model_name %s, got %v", tt.embedder.ModelName, settings["model_name"])
				}
			}
		})
	}
}

func TestEmbedder_ValidateModel(t *testing.T) {
	tests := []struct {
		name     string
		embedder *Embedder
		wantErr  bool
	}{
		{
			name: "Valid Azure OpenAI model",
			embedder: &Embedder{
				Subtype:   EmbedderSubtypeAzureOpenAI,
				ModelName: EmbedderModelAzureOpenAITextEmbedding3Small,
			},
			wantErr: false,
		},
		{
			name: "Invalid Azure OpenAI model",
			embedder: &Embedder{
				Subtype:   EmbedderSubtypeAzureOpenAI,
				ModelName: "invalid-model",
			},
			wantErr: true,
		},
		{
			name: "Valid Bedrock model",
			embedder: &Embedder{
				Subtype:   EmbedderSubtypeBedrock,
				ModelName: EmbedderModelBedrockTitanEmbedTextV2,
			},
			wantErr: false,
		},
		{
			name: "Invalid Bedrock model",
			embedder: &Embedder{
				Subtype:   EmbedderSubtypeBedrock,
				ModelName: "invalid-model",
			},
			wantErr: true,
		},
		{
			name: "Valid TogetherAI model",
			embedder: &Embedder{
				Subtype:   EmbedderSubtypeTogetherAI,
				ModelName: EmbedderModelTogetherAIM2Bert80M32kRetrieval,
			},
			wantErr: false,
		},
		{
			name: "Invalid TogetherAI model",
			embedder: &Embedder{
				Subtype:   EmbedderSubtypeTogetherAI,
				ModelName: "invalid-model",
			},
			wantErr: true,
		},
		{
			name: "Valid VoyageAI model",
			embedder: &Embedder{
				Subtype:   EmbedderSubtypeVoyageAI,
				ModelName: EmbedderModelVoyageAI3,
			},
			wantErr: false,
		},
		{
			name: "Invalid VoyageAI model",
			embedder: &Embedder{
				Subtype:   EmbedderSubtypeVoyageAI,
				ModelName: "invalid-model",
			},
			wantErr: true,
		},
		{
			name: "Unknown subtype",
			embedder: &Embedder{
				Subtype:   "unknown",
				ModelName: EmbedderModelVoyageAI3,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.embedder.ValidateModel()
			if (err != nil) != tt.wantErr {
				t.Errorf("Embedder.ValidateModel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
