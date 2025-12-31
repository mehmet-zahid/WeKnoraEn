package tools

// Tool names constants
const (
	ToolThinking            = "thinking"
	ToolTodoWrite           = "todo_write"
	ToolGrepChunks          = "grep_chunks"
	ToolKnowledgeSearch     = "knowledge_search"
	ToolListKnowledgeChunks = "list_knowledge_chunks"
	ToolQueryKnowledgeGraph = "query_knowledge_graph"
	ToolGetDocumentInfo     = "get_document_info"
	ToolDatabaseQuery       = "database_query"
	ToolDataAnalysis        = "data_analysis"
	ToolDataSchema          = "data_schema"
	ToolWebSearch           = "web_search"
	ToolWebFetch            = "web_fetch"
)

// AvailableTool defines a simple tool metadata used by settings APIs.
type AvailableTool struct {
	Name        string `json:"name"`
	Label       string `json:"label"`
	Description string `json:"description"`
}

// AvailableToolDefinitions returns the list of tools exposed to the UI.
// Keep this in sync with registered tools in this package.
func AvailableToolDefinitions() []AvailableTool {
	return []AvailableTool{
		{Name: ToolThinking, Label: "Thinking", Description: "Dynamic and reflective problem-solving thinking tool"},
		{Name: ToolTodoWrite, Label: "Plan Creation", Description: "Create structured research plans"},
		{Name: ToolGrepChunks, Label: "Keyword Search", Description: "Quickly locate documents and chunks containing specific keywords"},
		{Name: ToolKnowledgeSearch, Label: "Semantic Search", Description: "Understand questions and find semantically relevant content"},
		{Name: ToolListKnowledgeChunks, Label: "View Document Chunks", Description: "Get complete chunk content of documents"},
		{Name: ToolQueryKnowledgeGraph, Label: "Query Knowledge Graph", Description: "Query relationships from knowledge graph"},
		{Name: ToolGetDocumentInfo, Label: "Get Document Info", Description: "View document metadata"},
		{Name: ToolDatabaseQuery, Label: "Query Database", Description: "Query information from database"},
		{Name: ToolDataAnalysis, Label: "Data Analysis", Description: "Understand data files and perform data analysis"},
		{Name: ToolDataSchema, Label: "View Data Metadata", Description: "Get metadata of table files"},
	}
}

// DefaultAllowedTools returns the default allowed tools list.
func DefaultAllowedTools() []string {
	return []string{
		ToolThinking,
		ToolTodoWrite,
		ToolKnowledgeSearch,
		ToolGrepChunks,
		ToolListKnowledgeChunks,
		ToolQueryKnowledgeGraph,
		ToolGetDocumentInfo,
		ToolDatabaseQuery,
		ToolDataAnalysis,
		ToolDataSchema,
	}
}
