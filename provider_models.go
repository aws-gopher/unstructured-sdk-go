package unstructured

// Provider represents an AI model provider.
type Provider string

// Provider constants.
const (
	ProviderAuto      Provider = "auto"
	ProviderAnthropic Provider = "anthropic"
	ProviderOpenAI    Provider = "openai"
	ProviderBedrock   Provider = "bedrock"
)

// Model represents an AI model identifier.
type Model string

// Model constants.
const (
	ModelGPT4o                 Model = "gpt-4o"
	ModelGPT4oMini             Model = "gpt-4o-mini"
	ModelClaude35Sonnet        Model = "claude-3-5-sonnet-20241022"
	ModelClaude37Sonnet        Model = "claude-3-7-sonnet-20250219"
	ModelBedrockNovaLite       Model = "us.amazon.nova-lite-v1:0"
	ModelBedrockNovaPro        Model = "us.amazon.nova-pro-v1:0"
	ModelBedrockClaude3Opus    Model = "us.anthropic.claude-3-opus-20240229-v1:0"
	ModelBedrockClaude3Haiku   Model = "us.anthropic.claude-3-haiku-20240307-v1:0"
	ModelBedrockClaude3Sonnet  Model = "us.anthropic.claude-3-sonnet-20240229-v1:0"
	ModelBedrockClaude35Sonnet Model = "us.anthropic.claude-3-5-sonnet-20241022-v2:0"
	ModelBedrockLlama3211B     Model = "us.meta.llama3-2-11b-instruct-v1:0"
	ModelBedrockLlama3290B     Model = "us.meta.llama3-2-90b-instruct-v1:0"
)

func init() { var _ = providerModels }

var providerModels = map[Provider][]Model{
	ProviderOpenAI: {
		ModelGPT4o,
		ModelGPT4oMini,
	},
	ProviderAnthropic: {
		ModelClaude35Sonnet,
		ModelClaude37Sonnet,
	},
	ProviderBedrock: {
		ModelBedrockNovaLite,
		ModelBedrockNovaPro,
		ModelBedrockClaude3Opus,
		ModelBedrockClaude3Haiku,
		ModelBedrockClaude3Sonnet,
		ModelBedrockClaude35Sonnet,
		ModelBedrockLlama3211B,
		ModelBedrockLlama3290B,
	},
}
