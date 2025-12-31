package config

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/Tencent/WeKnora/internal/types"
	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// Config represents the application's overall configuration
type Config struct {
	Conversation    *ConversationConfig    `yaml:"conversation"     json:"conversation"`
	Server          *ServerConfig          `yaml:"server"           json:"server"`
	KnowledgeBase   *KnowledgeBaseConfig   `yaml:"knowledge_base"   json:"knowledge_base"`
	Tenant          *TenantConfig          `yaml:"tenant"           json:"tenant"`
	Models          []ModelConfig          `yaml:"models"           json:"models"`
	VectorDatabase  *VectorDatabaseConfig  `yaml:"vector_database"  json:"vector_database"`
	DocReader       *DocReaderConfig       `yaml:"docreader"        json:"docreader"`
	StreamManager   *StreamManagerConfig   `yaml:"stream_manager"   json:"stream_manager"`
	ExtractManager  *ExtractManagerConfig  `yaml:"extract"          json:"extract"`
	WebSearch       *WebSearchConfig       `yaml:"web_search"       json:"web_search"`
	PromptTemplates *PromptTemplatesConfig `yaml:"prompt_templates" json:"prompt_templates"`
}

type DocReaderConfig struct {
	Addr string `yaml:"addr" json:"addr"`
}

type VectorDatabaseConfig struct {
	Driver string `yaml:"driver" json:"driver"`
}

// ConversationConfig represents conversation service configuration
type ConversationConfig struct {
	MaxRounds                  int            `yaml:"max_rounds"                    json:"max_rounds"`
	KeywordThreshold           float64        `yaml:"keyword_threshold"             json:"keyword_threshold"`
	EmbeddingTopK              int            `yaml:"embedding_top_k"               json:"embedding_top_k"`
	VectorThreshold            float64        `yaml:"vector_threshold"              json:"vector_threshold"`
	RerankTopK                 int            `yaml:"rerank_top_k"                  json:"rerank_top_k"`
	RerankThreshold            float64        `yaml:"rerank_threshold"              json:"rerank_threshold"`
	FallbackStrategy           string         `yaml:"fallback_strategy"             json:"fallback_strategy"`
	FallbackResponse           string         `yaml:"fallback_response"             json:"fallback_response"`
	FallbackPrompt             string         `yaml:"fallback_prompt"               json:"fallback_prompt"`
	EnableRewrite              bool           `yaml:"enable_rewrite"                json:"enable_rewrite"`
	EnableQueryExpansion       bool           `yaml:"enable_query_expansion"        json:"enable_query_expansion"`
	EnableRerank               bool           `yaml:"enable_rerank"                 json:"enable_rerank"`
	Summary                    *SummaryConfig `yaml:"summary"                       json:"summary"`
	GenerateSessionTitlePrompt string         `yaml:"generate_session_title_prompt" json:"generate_session_title_prompt"`
	GenerateSummaryPrompt      string         `yaml:"generate_summary_prompt"       json:"generate_summary_prompt"`
	RewritePromptSystem        string         `yaml:"rewrite_prompt_system"         json:"rewrite_prompt_system"`
	RewritePromptUser          string         `yaml:"rewrite_prompt_user"           json:"rewrite_prompt_user"`
	SimplifyQueryPrompt        string         `yaml:"simplify_query_prompt"         json:"simplify_query_prompt"`
	SimplifyQueryPromptUser    string         `yaml:"simplify_query_prompt_user"    json:"simplify_query_prompt_user"`
	ExtractEntitiesPrompt      string         `yaml:"extract_entities_prompt"       json:"extract_entities_prompt"`
	ExtractRelationshipsPrompt string         `yaml:"extract_relationships_prompt"  json:"extract_relationships_prompt"`
	// GenerateQuestionsPrompt is used to generate questions for document chunks to improve recall
	GenerateQuestionsPrompt string `yaml:"generate_questions_prompt" json:"generate_questions_prompt"`
}

// SummaryConfig represents summary configuration
type SummaryConfig struct {
	MaxTokens           int     `yaml:"max_tokens"            json:"max_tokens"`
	RepeatPenalty       float64 `yaml:"repeat_penalty"        json:"repeat_penalty"`
	TopK                int     `yaml:"top_k"                 json:"top_k"`
	TopP                float64 `yaml:"top_p"                 json:"top_p"`
	FrequencyPenalty    float64 `yaml:"frequency_penalty"     json:"frequency_penalty"`
	PresencePenalty     float64 `yaml:"presence_penalty"      json:"presence_penalty"`
	Prompt              string  `yaml:"prompt"                json:"prompt"`
	ContextTemplate     string  `yaml:"context_template"      json:"context_template"`
	Temperature         float64 `yaml:"temperature"           json:"temperature"`
	Seed                int     `yaml:"seed"                  json:"seed"`
	MaxCompletionTokens int     `yaml:"max_completion_tokens" json:"max_completion_tokens"`
	NoMatchPrefix       string  `yaml:"no_match_prefix"       json:"no_match_prefix"`
	Thinking            *bool   `yaml:"thinking"              json:"thinking"`
}

// ServerConfig represents server configuration
type ServerConfig struct {
	Port            int           `yaml:"port"             json:"port"`
	Host            string        `yaml:"host"             json:"host"`
	LogPath         string        `yaml:"log_path"         json:"log_path"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" json:"shutdown_timeout" default:"30s"`
}

// KnowledgeBaseConfig represents knowledge base configuration
type KnowledgeBaseConfig struct {
	ChunkSize       int                    `yaml:"chunk_size"       json:"chunk_size"`
	ChunkOverlap    int                    `yaml:"chunk_overlap"    json:"chunk_overlap"`
	SplitMarkers    []string               `yaml:"split_markers"    json:"split_markers"`
	KeepSeparator   bool                   `yaml:"keep_separator"   json:"keep_separator"`
	ImageProcessing *ImageProcessingConfig `yaml:"image_processing" json:"image_processing"`
}

// ImageProcessingConfig represents image processing configuration
type ImageProcessingConfig struct {
	EnableMultimodal bool `yaml:"enable_multimodal" json:"enable_multimodal"`
}

// TenantConfig represents tenant configuration
type TenantConfig struct {
	DefaultSessionName        string `yaml:"default_session_name"        json:"default_session_name"`
	DefaultSessionTitle       string `yaml:"default_session_title"       json:"default_session_title"`
	DefaultSessionDescription string `yaml:"default_session_description" json:"default_session_description"`
	// EnableCrossTenantAccess enables cross-tenant access for users with permission
	EnableCrossTenantAccess bool `yaml:"enable_cross_tenant_access" json:"enable_cross_tenant_access"`
}

// PromptTemplate represents a prompt template
type PromptTemplate struct {
	ID               string `yaml:"id"                 json:"id"`
	Name             string `yaml:"name"               json:"name"`
	Description      string `yaml:"description"        json:"description"`
	Content          string `yaml:"content"            json:"content"`
	HasKnowledgeBase bool   `yaml:"has_knowledge_base" json:"has_knowledge_base,omitempty"`
	HasWebSearch     bool   `yaml:"has_web_search"     json:"has_web_search,omitempty"`
}

// PromptTemplatesConfig represents prompt templates configuration
type PromptTemplatesConfig struct {
	SystemPrompt    []PromptTemplate `yaml:"system_prompt"    json:"system_prompt"`
	ContextTemplate []PromptTemplate `yaml:"context_template" json:"context_template"`
	RewriteSystem   []PromptTemplate `yaml:"rewrite_system"   json:"rewrite_system"`
	RewriteUser     []PromptTemplate `yaml:"rewrite_user"     json:"rewrite_user"`
	Fallback        []PromptTemplate `yaml:"fallback"         json:"fallback"`
}

// ModelConfig represents model configuration
type ModelConfig struct {
	Type       string                 `yaml:"type"       json:"type"`
	Source     string                 `yaml:"source"     json:"source"`
	ModelName  string                 `yaml:"model_name" json:"model_name"`
	Parameters map[string]interface{} `yaml:"parameters" json:"parameters"`
}

// StreamManagerConfig represents stream manager configuration
type StreamManagerConfig struct {
	Type           string        `yaml:"type"            json:"type"`            // Type: "memory" or "redis"
	Redis          RedisConfig   `yaml:"redis"           json:"redis"`           // Redis configuration
	CleanupTimeout time.Duration `yaml:"cleanup_timeout" json:"cleanup_timeout"` // Cleanup timeout in seconds
}

// RedisConfig represents Redis configuration
type RedisConfig struct {
	Address  string        `yaml:"address"  json:"address"`  // Redis address
	Password string        `yaml:"password" json:"password"` // Redis password
	DB       int           `yaml:"db"       json:"db"`       // Redis database
	Prefix   string        `yaml:"prefix"   json:"prefix"`   // Key prefix
	TTL      time.Duration `yaml:"ttl"      json:"ttl"`      // Expiration time (hours)
}

// ExtractManagerConfig represents extraction manager configuration
type ExtractManagerConfig struct {
	ExtractGraph  *types.PromptTemplateStructured `yaml:"extract_graph"  json:"extract_graph"`
	ExtractEntity *types.PromptTemplateStructured `yaml:"extract_entity" json:"extract_entity"`
	FabriText     *FebriText                      `yaml:"fabri_text"     json:"fabri_text"`
}

type FebriText struct {
	WithTag   string `yaml:"with_tag"    json:"with_tag"`
	WithNoTag string `yaml:"with_no_tag" json:"with_no_tag"`
}

// LoadConfig loads configuration from config file
func LoadConfig() (*Config, error) {
	// Set config file name and paths
	viper.SetConfigName("config")         // Config file name (without extension)
	viper.SetConfigType("yaml")           // Config file type
	viper.AddConfigPath(".")              // Current directory
	viper.AddConfigPath("./config")       // config subdirectory
	viper.AddConfigPath("$HOME/.appname") // User directory
	viper.AddConfigPath("/etc/appname/")  // etc directory

	// Enable environment variable substitution
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	// Replace environment variable references in config
	configFileContent, err := os.ReadFile(viper.ConfigFileUsed())
	if err != nil {
		return nil, fmt.Errorf("error reading config file content: %w", err)
	}

	// Replace ${ENV_VAR} format environment variable references
	re := regexp.MustCompile(`\${([^}]+)}`)
	result := re.ReplaceAllStringFunc(string(configFileContent), func(match string) string {
		// Extract environment variable name (remove ${} part)
		envVar := match[2 : len(match)-1]
		// Get environment variable value, keep original if not exists
		if value := os.Getenv(envVar); value != "" {
			return value
		}
		return match
	})

	// Use processed config content
	viper.ReadConfig(strings.NewReader(result))

	// Parse config into struct
	var cfg Config
	if err := viper.Unmarshal(&cfg, func(dc *mapstructure.DecoderConfig) {
		dc.TagName = "yaml"
	}); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}
	fmt.Printf("Using configuration file: %s\n", viper.ConfigFileUsed())

	// Load prompt templates (from directory or config file)
	configDir := filepath.Dir(viper.ConfigFileUsed())
	promptTemplates, err := loadPromptTemplates(configDir)
	if err != nil {
		fmt.Printf("Warning: failed to load prompt templates from directory: %v\n", err)
		// If directory loading fails, use templates from config file (if any)
	} else if promptTemplates != nil {
		cfg.PromptTemplates = promptTemplates
	}

	return &cfg, nil
}

// promptTemplateFile is used to parse template files
type promptTemplateFile struct {
	Templates []PromptTemplate `yaml:"templates"`
}

// loadPromptTemplates loads prompt templates from directory
func loadPromptTemplates(configDir string) (*PromptTemplatesConfig, error) {
	templatesDir := filepath.Join(configDir, "prompt_templates")

	// Check if directory exists
	if _, err := os.Stat(templatesDir); os.IsNotExist(err) {
		return nil, nil // Directory doesn't exist, return nil to let caller use templates from config file
	}

	config := &PromptTemplatesConfig{}

	// Define template file mapping
	templateFiles := map[string]*[]PromptTemplate{
		"system_prompt.yaml":    &config.SystemPrompt,
		"context_template.yaml": &config.ContextTemplate,
		"rewrite_system.yaml":   &config.RewriteSystem,
		"rewrite_user.yaml":     &config.RewriteUser,
		"fallback.yaml":         &config.Fallback,
	}

	// Load each template file
	for filename, target := range templateFiles {
		filePath := filepath.Join(templatesDir, filename)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			continue // File doesn't exist, skip
		}

		data, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read %s: %w", filename, err)
		}

		var file promptTemplateFile
		if err := yaml.Unmarshal(data, &file); err != nil {
			return nil, fmt.Errorf("failed to parse %s: %w", filename, err)
		}

		*target = file.Templates
	}

	return config, nil
}

// WebSearchConfig represents the web search configuration
type WebSearchConfig struct {
	Providers []WebSearchProviderConfig `yaml:"providers" json:"providers"`
	Default   WebSearchDefaultConfig    `yaml:"default"   json:"default"`
	Timeout   int                       `yaml:"timeout"   json:"timeout"` // Timeout in seconds
}

// WebSearchProviderConfig represents configuration for a web search provider
type WebSearchProviderConfig struct {
	ID             string `yaml:"id"                    json:"id"`
	Name           string `yaml:"name"                  json:"name"`
	Free           bool   `yaml:"free"                  json:"free"`
	RequiresAPIKey bool   `yaml:"requires_api_key"      json:"requires_api_key"`
	Description    string `yaml:"description,omitempty" json:"description,omitempty"`
	APIURL         string `yaml:"api_url,omitempty"     json:"api_url,omitempty"`
}

// WebSearchDefaultConfig represents the default web search configuration
type WebSearchDefaultConfig struct {
	Provider          string   `yaml:"provider"           json:"provider"`
	MaxResults        int      `yaml:"max_results"        json:"max_results"`
	IncludeDate       bool     `yaml:"include_date"       json:"include_date"`
	CompressionMethod string   `yaml:"compression_method" json:"compression_method"`
	Blacklist         []string `yaml:"blacklist"          json:"blacklist"`
}
