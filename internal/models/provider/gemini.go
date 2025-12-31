package provider

import (
	"fmt"

	"github.com/Tencent/WeKnora/internal/types"
)

const (
	// GeminiBaseURL Google Gemini API BaseURL
	GeminiBaseURL = "https://generativelanguage.googleapis.com/v1beta"
	// GeminiOpenAICompatBaseURL Gemini OpenAI compatible mode BaseURL
	GeminiOpenAICompatBaseURL = "https://generativelanguage.googleapis.com/v1beta/openai"
)

// GeminiProvider implements Google Gemini Provider interface
type GeminiProvider struct{}

func init() {
	Register(&GeminiProvider{})
}

// Info returns Gemini provider metadata
func (p *GeminiProvider) Info() ProviderInfo {
	return ProviderInfo{
		Name:        ProviderGemini,
		DisplayName: "Google Gemini",
		Description: "gemini-3-flash-preview, gemini-2.5-pro, etc.",
		DefaultURLs: map[types.ModelType]string{
			types.ModelTypeKnowledgeQA: GeminiOpenAICompatBaseURL,
		},
		ModelTypes: []types.ModelType{
			types.ModelTypeKnowledgeQA,
		},
		RequiresAuth: true,
	}
}

// ValidateConfig validates Gemini provider configuration
func (p *GeminiProvider) ValidateConfig(config *Config) error {
	if config.APIKey == "" {
		return fmt.Errorf("API key is required for Google Gemini provider")
	}
	if config.ModelName == "" {
		return fmt.Errorf("model name is required")
	}
	return nil
}
