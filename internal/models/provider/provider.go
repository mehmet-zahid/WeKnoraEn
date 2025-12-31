// Package provider defines the unified interface and registry for multi-vendor model API adapters.
package provider

import (
	"fmt"
	"strings"
	"sync"

	"github.com/Tencent/WeKnora/internal/types"
)

// ProviderName represents model service provider name
type ProviderName string

const (
	// OpenAI
	ProviderOpenAI ProviderName = "openai"
	// Alibaba Cloud DashScope
	ProviderAliyun ProviderName = "aliyun"
	// Zhipu AI (GLM series)
	ProviderZhipu ProviderName = "zhipu"
	// OpenRouter
	ProviderOpenRouter ProviderName = "openrouter"
	// SiliconFlow
	ProviderSiliconFlow ProviderName = "siliconflow"
	// Jina AI (Embedding and Rerank)
	ProviderJina ProviderName = "jina"
	// Generic OpenAI compatible (custom deployment)
	ProviderGeneric ProviderName = "generic"
	// DeepSeek
	ProviderDeepSeek ProviderName = "deepseek"
	// Google Gemini
	ProviderGemini ProviderName = "gemini"
	// Volcengine Ark
	ProviderVolcengine ProviderName = "volcengine"
	// Tencent Hunyuan
	ProviderHunyuan ProviderName = "hunyuan"
	// MiniMax
	ProviderMiniMax ProviderName = "minimax"
	// Xiaomi Mimo
	ProviderMimo ProviderName = "mimo"
)

// AllProviders returns all registered provider names
func AllProviders() []ProviderName {
	return []ProviderName{
		ProviderGeneric,
		ProviderAliyun,
		ProviderZhipu,
		ProviderVolcengine,
		ProviderHunyuan,
		ProviderSiliconFlow,
		ProviderDeepSeek,
		ProviderMiniMax,
		ProviderOpenAI,
		ProviderGemini,
		ProviderOpenRouter,
		ProviderJina,
		ProviderMimo,
	}
}

// ProviderInfo contains provider metadata
type ProviderInfo struct {
	Name         ProviderName               // Provider identifier
	DisplayName  string                     // Human-readable name
	Description  string                     // Provider description
	DefaultURLs  map[types.ModelType]string // Default BaseURL by model type
	ModelTypes   []types.ModelType          // Supported model types
	RequiresAuth bool                       // Whether API key is required
	ExtraFields  []ExtraFieldConfig         // Additional configuration fields
}

// GetDefaultURL gets the default URL for specified model type
func (p ProviderInfo) GetDefaultURL(modelType types.ModelType) string {
	if url, ok := p.DefaultURLs[modelType]; ok {
		return url
	}
	// Fallback to Chat URL
	if url, ok := p.DefaultURLs[types.ModelTypeKnowledgeQA]; ok {
		return url
	}
	return ""
}

// ExtraFieldConfig defines additional configuration fields for provider
type ExtraFieldConfig struct {
	Key         string `json:"key"`
	Label       string `json:"label"`
	Type        string `json:"type"` // "string", "number", "boolean", "select"
	Required    bool   `json:"required"`
	Default     string `json:"default"`
	Placeholder string `json:"placeholder"`
	Options     []struct {
		Label string `json:"label"`
		Value string `json:"value"`
	} `json:"options,omitempty"`
}

// Config represents model provider configuration
type Config struct {
	Provider  ProviderName   `json:"provider"`
	BaseURL   string         `json:"base_url"`
	APIKey    string         `json:"api_key"`
	ModelName string         `json:"model_name"`
	ModelID   string         `json:"model_id"`
	Extra     map[string]any `json:"extra,omitempty"`
}

type Provider interface {
	// Info returns provider metadata
	Info() ProviderInfo

	// ValidateConfig validates provider configuration
	ValidateConfig(config *Config) error
}

// registry stores all registered providers
var (
	registryMu sync.RWMutex
	registry   = make(map[ProviderName]Provider)
)

// Register adds a provider to the global registry
func Register(p Provider) {
	registryMu.Lock()
	defer registryMu.Unlock()
	registry[p.Info().Name] = p
}

// Get gets provider from registry by name
func Get(name ProviderName) (Provider, bool) {
	registryMu.RLock()
	defer registryMu.RUnlock()
	p, ok := registry[name]
	return p, ok
}

// GetOrDefault gets provider from registry by name, returns default provider if not found
func GetOrDefault(name ProviderName) Provider {
	p, ok := Get(name)
	if ok {
		return p
	}
	// If not found, return default provider
	p, _ = Get(ProviderGeneric)
	return p
}

// List returns all registered providers (in order defined by AllProviders)
func List() []ProviderInfo {
	registryMu.RLock()
	defer registryMu.RUnlock()

	result := make([]ProviderInfo, 0, len(registry))
	for _, name := range AllProviders() {
		if p, ok := registry[name]; ok {
			result = append(result, p.Info())
		}
	}
	return result
}

// ListByModelType returns all providers supporting specified model type (in order defined by AllProviders)
func ListByModelType(modelType types.ModelType) []ProviderInfo {
	registryMu.RLock()
	defer registryMu.RUnlock()

	result := make([]ProviderInfo, 0)
	for _, name := range AllProviders() {
		if p, ok := registry[name]; ok {
			info := p.Info()
			for _, t := range info.ModelTypes {
				if t == modelType {
					result = append(result, info)
					break
				}
			}
		}
	}
	return result
}

// DetectProvider detects provider by BaseURL
func DetectProvider(baseURL string) ProviderName {
	switch {
	case containsAny(baseURL, "dashscope.aliyuncs.com"):
		return ProviderAliyun
	case containsAny(baseURL, "open.bigmodel.cn", "zhipu"):
		return ProviderZhipu
	case containsAny(baseURL, "openrouter.ai"):
		return ProviderOpenRouter
	case containsAny(baseURL, "siliconflow.cn"):
		return ProviderSiliconFlow
	case containsAny(baseURL, "api.jina.ai"):
		return ProviderJina
	case containsAny(baseURL, "api.openai.com"):
		return ProviderOpenAI
	case containsAny(baseURL, "api.deepseek.com"):
		return ProviderDeepSeek
	case containsAny(baseURL, "generativelanguage.googleapis.com"):
		return ProviderGemini
	case containsAny(baseURL, "volces.com", "volcengine"):
		return ProviderVolcengine
	case containsAny(baseURL, "hunyuan.cloud.tencent.com"):
		return ProviderHunyuan
	case containsAny(baseURL, "minimax.io", "minimaxi.com"):
		return ProviderMiniMax
	case containsAny(baseURL, "xiaomimimo.com"):
		return ProviderMimo
	default:
		return ProviderGeneric
	}
}

func containsAny(s string, substrs ...string) bool {
	for _, sub := range substrs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

func NewConfigFromModel(model *types.Model) (*Config, error) {
	if model == nil {
		return nil, fmt.Errorf("model is nil")
	}

	providerName := ProviderName(model.Parameters.Provider)
	if providerName == "" {
		providerName = DetectProvider(model.Parameters.BaseURL)
	}

	return &Config{
		Provider:  providerName,
		BaseURL:   model.Parameters.BaseURL,
		APIKey:    model.Parameters.APIKey,
		ModelName: model.Name,
		ModelID:   model.ID,
	}, nil
}
