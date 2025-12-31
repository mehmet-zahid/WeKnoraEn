package provider

import (
	"fmt"
	"strings"

	"github.com/Tencent/WeKnora/internal/types"
)

const (
	// AliyunChatBaseURL default BaseURL for Alibaba Cloud DashScope Chat/Embedding
	AliyunChatBaseURL = "https://dashscope.aliyuncs.com/compatible-mode/v1"
	// AliyunRerankBaseURL default BaseURL for Alibaba Cloud DashScope Rerank
	AliyunRerankBaseURL = "https://dashscope.aliyuncs.com/api/v1/services/rerank/text-rerank/text-rerank"
)

// AliyunProvider implements Alibaba Cloud DashScope Provider interface
type AliyunProvider struct{}

func init() {
	Register(&AliyunProvider{})
}

// Info returns Alibaba Cloud provider metadata
func (p *AliyunProvider) Info() ProviderInfo {
	return ProviderInfo{
		Name:        ProviderAliyun,
		DisplayName: "Alibaba Cloud DashScope",
		Description: "qwen-plus, tongyi-embedding-vision-plus, qwen3-rerank, etc.",
		DefaultURLs: map[types.ModelType]string{
			types.ModelTypeKnowledgeQA: AliyunChatBaseURL,
			types.ModelTypeEmbedding:   AliyunChatBaseURL,
			types.ModelTypeRerank:      AliyunRerankBaseURL,
			types.ModelTypeVLLM:        AliyunChatBaseURL,
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

// ValidateConfig validates Alibaba Cloud provider configuration
func (p *AliyunProvider) ValidateConfig(config *Config) error {
	if config.APIKey == "" {
		return fmt.Errorf("API key is required for Aliyun DashScope")
	}
	if config.ModelName == "" {
		return fmt.Errorf("model name is required")
	}
	return nil
}

// IsQwen3Model checks if model name is Qwen3 model
// Qwen3 models require special handling of enable_thinking parameter
func IsQwen3Model(modelName string) bool {
	return strings.HasPrefix(modelName, "qwen3-")
}

// IsDeepSeekModel checks if model name is DeepSeek model
// DeepSeek models don't support tool_choice parameter
func IsDeepSeekModel(modelName string) bool {
	return strings.Contains(strings.ToLower(modelName), "deepseek")
}
