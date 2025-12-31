package types

// PromptPlaceholder represents a placeholder that can be used in prompt templates
type PromptPlaceholder struct {
	// Name is the placeholder name (without braces), e.g., "query"
	Name string `json:"name"`
	// Label is a short label for the placeholder
	Label string `json:"label"`
	// Description explains what this placeholder represents
	Description string `json:"description"`
}

// PromptFieldType represents the type of prompt field
type PromptFieldType string

const (
	// PromptFieldSystemPrompt is for system prompts (normal mode)
	PromptFieldSystemPrompt PromptFieldType = "system_prompt"
	// PromptFieldAgentSystemPrompt is for agent mode system prompts
	PromptFieldAgentSystemPrompt PromptFieldType = "agent_system_prompt"
	// PromptFieldContextTemplate is for context templates
	PromptFieldContextTemplate PromptFieldType = "context_template"
	// PromptFieldRewriteSystemPrompt is for rewrite system prompts
	PromptFieldRewriteSystemPrompt PromptFieldType = "rewrite_system_prompt"
	// PromptFieldRewritePrompt is for rewrite user prompts
	PromptFieldRewritePrompt PromptFieldType = "rewrite_prompt"
	// PromptFieldFallbackPrompt is for fallback prompts
	PromptFieldFallbackPrompt PromptFieldType = "fallback_prompt"
)

// All available placeholders in the system
var (
	// Common placeholders
	PlaceholderQuery = PromptPlaceholder{
		Name:        "query",
		Label:       "User Question",
		Description: "User's current question or query content",
	}

	PlaceholderContexts = PromptPlaceholder{
		Name:        "contexts",
		Label:       "Retrieved Content",
		Description: "List of relevant content retrieved from knowledge base",
	}

	PlaceholderCurrentTime = PromptPlaceholder{
		Name:        "current_time",
		Label:       "Current Time",
		Description: "Current system time (format: 2006-01-02 15:04:05)",
	}

	PlaceholderCurrentWeek = PromptPlaceholder{
		Name:        "current_week",
		Label:       "Current Week",
		Description: "Current day of week (e.g., Monday, 星期一)",
	}

	// Rewrite prompt placeholders
	PlaceholderConversation = PromptPlaceholder{
		Name:        "conversation",
		Label:       "Conversation History",
		Description: "Formatted conversation history content for multi-turn conversation rewriting",
	}

	PlaceholderYesterday = PromptPlaceholder{
		Name:        "yesterday",
		Label:       "Yesterday Date",
		Description: "Yesterday's date (format: 2006-01-02)",
	}

	PlaceholderAnswer = PromptPlaceholder{
		Name:        "answer",
		Label:       "Assistant Answer",
		Description: "Assistant's answer content (for conversation history formatting)",
	}

	// Agent mode specific placeholders
	PlaceholderKnowledgeBases = PromptPlaceholder{
		Name:        "knowledge_bases",
		Label:       "Knowledge Base List",
		Description: "Auto-formatted knowledge base list, including name, description, document count, etc.",
	}

	PlaceholderWebSearchStatus = PromptPlaceholder{
		Name:        "web_search_status",
		Label:       "Web Search Status",
		Description: "Status of whether web search tool is enabled (Enabled or Disabled)",
	}
)

// PlaceholdersByField returns the available placeholders for a specific prompt field type
func PlaceholdersByField(fieldType PromptFieldType) []PromptPlaceholder {
	switch fieldType {
	case PromptFieldSystemPrompt:
		// Normal mode system prompt
		return []PromptPlaceholder{
			PlaceholderQuery,
			PlaceholderContexts,
			PlaceholderCurrentTime,
			PlaceholderCurrentWeek,
		}
	case PromptFieldAgentSystemPrompt:
		// Agent mode system prompt
		return []PromptPlaceholder{
			PlaceholderKnowledgeBases,
			PlaceholderWebSearchStatus,
			PlaceholderCurrentTime,
		}
	case PromptFieldContextTemplate:
		return []PromptPlaceholder{
			PlaceholderQuery,
			PlaceholderContexts,
			PlaceholderCurrentTime,
			PlaceholderCurrentWeek,
		}
	case PromptFieldRewriteSystemPrompt:
		// Rewrite system prompt supports same placeholders as rewrite user prompt
		return []PromptPlaceholder{
			PlaceholderQuery,
			PlaceholderConversation,
			PlaceholderCurrentTime,
			PlaceholderYesterday,
		}
	case PromptFieldRewritePrompt:
		return []PromptPlaceholder{
			PlaceholderQuery,
			PlaceholderConversation,
			PlaceholderCurrentTime,
			PlaceholderYesterday,
		}
	case PromptFieldFallbackPrompt:
		return []PromptPlaceholder{
			PlaceholderQuery,
		}
	default:
		return []PromptPlaceholder{}
	}
}

// AllPlaceholders returns all available placeholders in the system
func AllPlaceholders() []PromptPlaceholder {
	return []PromptPlaceholder{
		PlaceholderQuery,
		PlaceholderContexts,
		PlaceholderCurrentTime,
		PlaceholderCurrentWeek,
		PlaceholderConversation,
		PlaceholderYesterday,
		PlaceholderAnswer,
		PlaceholderKnowledgeBases,
		PlaceholderWebSearchStatus,
	}
}

// PlaceholderMap returns a map of field types to their available placeholders
func PlaceholderMap() map[PromptFieldType][]PromptPlaceholder {
	return map[PromptFieldType][]PromptPlaceholder{
		PromptFieldSystemPrompt:        PlaceholdersByField(PromptFieldSystemPrompt),
		PromptFieldAgentSystemPrompt:   PlaceholdersByField(PromptFieldAgentSystemPrompt),
		PromptFieldContextTemplate:     PlaceholdersByField(PromptFieldContextTemplate),
		PromptFieldRewriteSystemPrompt: PlaceholdersByField(PromptFieldRewriteSystemPrompt),
		PromptFieldRewritePrompt:       PlaceholdersByField(PromptFieldRewritePrompt),
		PromptFieldFallbackPrompt:      PlaceholdersByField(PromptFieldFallbackPrompt),
	}
}
