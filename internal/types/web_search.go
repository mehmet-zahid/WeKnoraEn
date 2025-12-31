package types

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// WebSearchConfig represents the web search configuration for a tenant
type WebSearchConfig struct {
	Provider          string   `json:"provider"`           // Search engine provider ID
	APIKey            string   `json:"api_key"`            // API key (if required)
	MaxResults        int      `json:"max_results"`        // Maximum number of search results
	IncludeDate       bool     `json:"include_date"`       // Whether to include date
	CompressionMethod string   `json:"compression_method"` // Compression method: none, summary, extract, rag
	Blacklist         []string `json:"blacklist"`          // Blacklist rule list
	// RAG compression related configuration
	EmbeddingModelID   string `json:"embedding_model_id,omitempty"`  // Embedding model ID (for RAG compression)
	EmbeddingDimension int    `json:"embedding_dimension,omitempty"` // Embedding dimension (for RAG compression)
	RerankModelID      string `json:"rerank_model_id,omitempty"`     // Rerank model ID (for RAG compression)
	DocumentFragments  int    `json:"document_fragments,omitempty"`  // Number of document fragments (for RAG compression)
}

// Value implements driver.Valuer interface for WebSearchConfig
func (c WebSearchConfig) Value() (driver.Value, error) {
	return json.Marshal(c)
}

// Scan implements sql.Scanner interface for WebSearchConfig
func (c *WebSearchConfig) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(b, c)
}

// WebSearchResult represents a single web search result
type WebSearchResult struct {
	Title       string     `json:"title"`                  // Search result title
	URL         string     `json:"url"`                    // Result URL
	Snippet     string     `json:"snippet"`                // Snippet
	Content     string     `json:"content"`                // Full content (optional, requires additional fetching)
	Source      string     `json:"source"`                 // Source (e.g., duckduckgo, etc.)
	PublishedAt *time.Time `json:"published_at,omitempty"` // Publication time (if available)
}

// WebSearchProviderInfo represents information about a web search provider
type WebSearchProviderInfo struct {
	ID             string `json:"id"`                // Provider ID
	Name           string `json:"name"`              // Provider name
	Free           bool   `json:"free"`              // Whether it's free
	RequiresAPIKey bool   `json:"requires_api_key"`  // Whether API key is required
	Description    string `json:"description"`       // Description
	APIURL         string `json:"api_url,omitempty"` // API URL (optional)
}
