package provider

import (
	"fmt"

	"github.com/Tencent/WeKnora/internal/types"
)

const (
	SiliconFlowBaseURL = "https://api.siliconflow.cn/v1"
)

// SiliconFlowProvider implements SiliconFlow Provider interface
type SiliconFlowProvider struct{}

func init() {
	Register(&SiliconFlowProvider{})
}

// Info returns SiliconFlow provider metadata
func (p *SiliconFlowProvider) Info() ProviderInfo {
	return ProviderInfo{
		Name:        ProviderSiliconFlow,
		DisplayName: "SiliconFlow",
		Description: "deepseek-ai/DeepSeek-V3.1, etc.",
		DefaultURLs: map[types.ModelType]string{
			types.ModelTypeKnowledgeQA: SiliconFlowBaseURL,
			types.ModelTypeEmbedding:   SiliconFlowBaseURL,
			types.ModelTypeRerank:      SiliconFlowBaseURL,
			types.ModelTypeVLLM:        SiliconFlowBaseURL,
		},
		ModelTypes: []types.ModelType{
			types.ModelTypeKnowledgeQA,
			types.ModelTypeEmbedding,
			types.ModelTypeRerank,
			types.ModelTypeVLLM,
		},
		RequiresAuth: true,
	}
}

// ValidateConfig validates SiliconFlow provider configuration
func (p *SiliconFlowProvider) ValidateConfig(config *Config) error {
	if config.APIKey == "" {
		return fmt.Errorf("API key is required for SiliconFlow provider")
	}
	return nil
}
