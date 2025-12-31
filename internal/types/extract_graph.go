package types

const (
	TypeChunkExtract       = "chunk:extract"
	TypeDocumentProcess    = "document:process"    // Document processing task
	TypeFAQImport          = "faq:import"          // FAQ import task
	TypeQuestionGeneration = "question:generation" // Question generation task
	TypeSummaryGeneration  = "summary:generation"  // Summary generation task
	TypeKBClone            = "kb:clone"            // Knowledge base clone task
	TypeIndexDelete        = "index:delete"        // Index deletion task
	TypeKBDelete           = "kb:delete"           // Knowledge base deletion task
	TypeDataTableSummary   = "datatable:summary"   // Data table summary task
)

// ExtractChunkPayload represents the extract chunk task payload
type ExtractChunkPayload struct {
	TenantID uint64 `json:"tenant_id"`
	ChunkID  string `json:"chunk_id"`
	ModelID  string `json:"model_id"`
}

// DocumentProcessPayload represents the document process task payload
type DocumentProcessPayload struct {
	RequestId                string   `json:"request_id"`
	TenantID                 uint64   `json:"tenant_id"`
	KnowledgeID              string   `json:"knowledge_id"`
	KnowledgeBaseID          string   `json:"knowledge_base_id"`
	FilePath                 string   `json:"file_path,omitempty"` // File path (used for file import)
	FileName                 string   `json:"file_name,omitempty"` // File name (used for file import)
	FileType                 string   `json:"file_type,omitempty"` // File type (used for file import)
	URL                      string   `json:"url,omitempty"`       // URL (used for URL import)
	Passages                 []string `json:"passages,omitempty"`  // Text passages (used for text import)
	EnableMultimodel         bool     `json:"enable_multimodel"`
	EnableQuestionGeneration bool     `json:"enable_question_generation"` // Whether to enable question generation
	QuestionCount            int      `json:"question_count,omitempty"`   // Number of questions to generate per chunk
}

// FAQImportPayload represents the FAQ import task payload
type FAQImportPayload struct {
	TenantID    uint64            `json:"tenant_id"`
	TaskID      string            `json:"task_id"`
	KBID        string            `json:"kb_id"`
	KnowledgeID string            `json:"knowledge_id"`
	Entries     []FAQEntryPayload `json:"entries"`
	Mode        string            `json:"mode"`
}

// QuestionGenerationPayload represents the question generation task payload
type QuestionGenerationPayload struct {
	TenantID        uint64 `json:"tenant_id"`
	KnowledgeBaseID string `json:"knowledge_base_id"`
	KnowledgeID     string `json:"knowledge_id"`
	QuestionCount   int    `json:"question_count"`
}

// SummaryGenerationPayload represents the summary generation task payload
type SummaryGenerationPayload struct {
	TenantID        uint64 `json:"tenant_id"`
	KnowledgeBaseID string `json:"knowledge_base_id"`
	KnowledgeID     string `json:"knowledge_id"`
}

// KBClonePayload represents the knowledge base clone task payload
type KBClonePayload struct {
	TenantID uint64 `json:"tenant_id"`
	TaskID   string `json:"task_id"`
	SourceID string `json:"source_id"`
	TargetID string `json:"target_id"`
}

// IndexDeletePayload represents the index delete task payload
type IndexDeletePayload struct {
	TenantID         uint64                  `json:"tenant_id"`
	KnowledgeBaseID  string                  `json:"knowledge_base_id"`
	EmbeddingModelID string                  `json:"embedding_model_id"`
	KBType           string                  `json:"kb_type"`
	ChunkIDs         []string                `json:"chunk_ids"`
	EffectiveEngines []RetrieverEngineParams `json:"effective_engines"`
}

// KBDeletePayload represents the knowledge base delete task payload
type KBDeletePayload struct {
	TenantID         uint64                  `json:"tenant_id"`
	KnowledgeBaseID  string                  `json:"knowledge_base_id"`
	EffectiveEngines []RetrieverEngineParams `json:"effective_engines"`
}

// KBCloneTaskStatus represents the status of a knowledge base clone task
type KBCloneTaskStatus string

const (
	KBCloneStatusPending    KBCloneTaskStatus = "pending"
	KBCloneStatusProcessing KBCloneTaskStatus = "processing"
	KBCloneStatusCompleted  KBCloneTaskStatus = "completed"
	KBCloneStatusFailed     KBCloneTaskStatus = "failed"
)

// KBCloneProgress represents the progress of a knowledge base clone task
type KBCloneProgress struct {
	TaskID    string            `json:"task_id"`
	SourceID  string            `json:"source_id"`
	TargetID  string            `json:"target_id"`
	Status    KBCloneTaskStatus `json:"status"`
	Progress  int               `json:"progress"`   // 0-100
	Total     int               `json:"total"`      // Total knowledge count
	Processed int               `json:"processed"`  // Processed count
	Message   string            `json:"message"`    // Status message
	Error     string            `json:"error"`      // Error information
	CreatedAt int64             `json:"created_at"` // Task creation time
	UpdatedAt int64             `json:"updated_at"` // Last update time
}

// ChunkContext represents chunk content with surrounding context
type ChunkContext struct {
	ChunkID     string `json:"chunk_id"`
	Content     string `json:"content"`
	PrevContent string `json:"prev_content,omitempty"` // Previous chunk content for context
	NextContent string `json:"next_content,omitempty"` // Next chunk content for context
}

// PromptTemplateStructured represents the prompt template structured
type PromptTemplateStructured struct {
	Description string      `json:"description"`
	Tags        []string    `json:"tags"`
	Examples    []GraphData `json:"examples"`
}

type GraphNode struct {
	Name       string   `json:"name,omitempty"`
	Chunks     []string `json:"chunks,omitempty"`
	Attributes []string `json:"attributes,omitempty"`
}

// GraphRelation represents the relation of the graph
type GraphRelation struct {
	Node1 string `json:"node1,omitempty"`
	Node2 string `json:"node2,omitempty"`
	Type  string `json:"type,omitempty"`
}

type GraphData struct {
	Text     string           `json:"text,omitempty"`
	Node     []*GraphNode     `json:"node,omitempty"`
	Relation []*GraphRelation `json:"relation,omitempty"`
}

// NameSpace represents the name space of the knowledge base and knowledge
type NameSpace struct {
	KnowledgeBase string `json:"knowledge_base"`
	Knowledge     string `json:"knowledge"`
}

// Labels returns the labels of the name space
func (n NameSpace) Labels() []string {
	res := make([]string, 0)
	if n.KnowledgeBase != "" {
		res = append(res, n.KnowledgeBase)
	}
	if n.Knowledge != "" {
		res = append(res, n.Knowledge)
	}
	return res
}
