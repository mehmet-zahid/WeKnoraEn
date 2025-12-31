package provider

import (
	"fmt"

	"github.com/Tencent/WeKnora/internal/types"
)

// GenericProvider implements generic OpenAI compatible Provider interface
type GenericProvider struct{}

func init() {
	Register(&GenericProvider{})
}

// Info returns generic provider metadata
func (p *GenericProvider) Info() ProviderInfo {
	return ProviderInfo{
		Name:        ProviderGeneric,
		DisplayName: "Custom (OpenAI Format Compatible)",
		Description: "Generic API endpoint",
		DefaultURLs: map[types.ModelType]string{}, // User needs to configure manually
		ModelTypes: []types.ModelType{
			types.ModelTypeKnowledgeQA,
			types.ModelTypeEmbedding,
			types.ModelTypeRerank,
			types.ModelTypeVLLM,
		},
		RequiresAuth: false, // 可能需要也可能不需要
	}
}

// ValidateConfig validates generic provider configuration
func (p *GenericProvider) ValidateConfig(config *Config) error {
	if config.BaseURL == "" {
		return fmt.Errorf("base URL is required for generic provider")
	}
	if config.ModelName == "" {
		return fmt.Errorf("model name is required")
	}
	return nil
}
