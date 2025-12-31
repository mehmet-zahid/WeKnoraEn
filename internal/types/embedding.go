package types

// SourceType represents the type of content source
type SourceType int

const (
	ChunkSourceType   SourceType = iota // Source is a text chunk
	PassageSourceType                   // Source is a passage
	SummarySourceType                   // Source is a summary
)

// MatchType represents the type of matching algorithm
type MatchType int

const (
	MatchTypeEmbedding MatchType = iota
	MatchTypeKeywords
	MatchTypeNearByChunk
	MatchTypeHistory
	MatchTypeParentChunk   // Parent Chunk match type
	MatchTypeRelationChunk // Relation Chunk match type
	MatchTypeGraph
	MatchTypeWebSearch    // Web search match type
	MatchTypeDirectLoad   // Direct load match type
	MatchTypeDataAnalysis // Data analysis match type
)

// IndexInfo contains information about indexed content
type IndexInfo struct {
	ID              string     // Unique identifier
	Content         string     // Content text
	SourceID        string     // ID of the source document
	SourceType      SourceType // Type of the source
	ChunkID         string     // ID of the text chunk
	KnowledgeID     string     // ID of the knowledge
	KnowledgeBaseID string     // ID of the knowledge base
	KnowledgeType   string     // Type of the knowledge (e.g., "faq", "manual")
	TagID           string     // Tag ID for categorization (used for FAQ priority filtering)
	IsEnabled       bool       // Whether the chunk is enabled for retrieval
}
