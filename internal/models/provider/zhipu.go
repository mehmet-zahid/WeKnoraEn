package provider

import (
	"fmt"

	"github.com/Tencent/WeKnora/internal/types"
)

const (
	// ZhipuChatBaseURL default BaseURL for Zhipu AI Chat
	ZhipuChatBaseURL = "https://open.bigmodel.cn/api/paas/v4"
	// ZhipuEmbeddingBaseURL default BaseURL for Zhipu AI Embedding
	ZhipuEmbeddingBaseURL = "https://open.bigmodel.cn/api/paas/v4/embeddings"
	// ZhipuRerankBaseURL default BaseURL for Zhipu AI Rerank
	ZhipuRerankBaseURL = "https://open.bigmodel.cn/api/paas/v4/rerank"
)

// ZhipuProvider implements Zhipu AI Provider interface
type ZhipuProvider struct{}

func init() {
	Register(&ZhipuProvider{})
}

// Info returns Zhipu AI provider metadata
func (p *ZhipuProvider) Info() ProviderInfo {
	return ProviderInfo{
		Name:        ProviderZhipu,
		DisplayName: "Zhipu BigModel",
		Description: "glm-4.7, embedding-3, rerank, etc.",
		DefaultURLs: map[types.ModelType]string{
			types.ModelTypeKnowledgeQA: ZhipuChatBaseURL,
			types.ModelTypeEmbedding:   ZhipuEmbeddingBaseURL,
			types.ModelTypeRerank:      ZhipuRerankBaseURL,
			types.ModelTypeVLLM:        ZhipuChatBaseURL,
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

// ValidateConfig validates Zhipu AI provider configuration
func (p *ZhipuProvider) ValidateConfig(config *Config) error {
	if config.APIKey == "" {
		return fmt.Errorf("API key is required for Zhipu AI")
	}
	if config.ModelName == "" {
		return fmt.Errorf("model name is required")
	}
	return nil
}
