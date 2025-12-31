// Package types defines data structures and types used throughout the system
// These types are shared across different service modules to ensure data consistency
package types

import (
	"time"

	"gorm.io/gorm"
)

// ChunkType defines different types of Chunk
type ChunkType = string

const (
	// ChunkTypeText represents ordinary text Chunk
	ChunkTypeText ChunkType = "text"
	// ChunkTypeImageOCR represents image OCR text Chunk
	ChunkTypeImageOCR ChunkType = "image_ocr"
	// ChunkTypeImageCaption represents image caption Chunk
	ChunkTypeImageCaption ChunkType = "image_caption"
	// ChunkTypeSummary represents summary type Chunk
	ChunkTypeSummary = "summary"
	// ChunkTypeEntity represents entity type Chunk
	ChunkTypeEntity ChunkType = "entity"
	// ChunkTypeRelationship represents relationship type Chunk
	ChunkTypeRelationship ChunkType = "relationship"
	// ChunkTypeFAQ represents FAQ entry Chunk
	ChunkTypeFAQ ChunkType = "faq"
	// ChunkTypeWebSearch represents Web search result Chunk
	ChunkTypeWebSearch ChunkType = "web_search"
	// ChunkTypeTableSummary represents data table summary Chunk
	ChunkTypeTableSummary ChunkType = "table_summary"
	// ChunkTypeTableColumn represents data table column description Chunk
	ChunkTypeTableColumn ChunkType = "table_column"
)

// ChunkStatus defines different states of Chunk
type ChunkStatus int

const (
	ChunkStatusDefault ChunkStatus = 0
	// ChunkStatusStored represents stored Chunk
	ChunkStatusStored ChunkStatus = 1
	// ChunkStatusIndexed represents indexed Chunk
	ChunkStatusIndexed ChunkStatus = 2
)

// ChunkFlags defines Chunk flag bits for managing multiple boolean states
type ChunkFlags int

const (
	// ChunkFlagRecommended represents recommended state (1 << 0 = 1)
	// When this flag is set, the Chunk can be recommended to users
	ChunkFlagRecommended ChunkFlags = 1 << 0
	// More flags can be extended in the future:
	// ChunkFlagPinned ChunkFlags = 1 << 1  // Pinned
	// ChunkFlagHot    ChunkFlags = 1 << 2  // Hot
)

// HasFlag checks if specified flag is set
func (f ChunkFlags) HasFlag(flag ChunkFlags) bool {
	return f&flag != 0
}

// SetFlag sets specified flag
func (f ChunkFlags) SetFlag(flag ChunkFlags) ChunkFlags {
	return f | flag
}

// ClearFlag clears specified flag
func (f ChunkFlags) ClearFlag(flag ChunkFlags) ChunkFlags {
	return f &^ flag
}

// ToggleFlag toggles specified flag
func (f ChunkFlags) ToggleFlag(flag ChunkFlags) ChunkFlags {
	return f ^ flag
}

// ImageInfo represents image information associated with Chunk
type ImageInfo struct {
	// Image URL (COS)
	URL string `json:"url"          gorm:"type:text"`
	// Original image URL
	OriginalURL string `json:"original_url" gorm:"type:text"`
	// Start position of image in text
	StartPos int `json:"start_pos"`
	// End position of image in text
	EndPos int `json:"end_pos"`
	// Image caption
	Caption string `json:"caption"`
	// Image OCR text
	OCRText string `json:"ocr_text"`
}

// Chunk represents a document chunk
// Chunks are meaningful text segments extracted from original documents
// and are the basic units of knowledge base retrieval
// Each chunk contains a portion of the original content
// and maintains its positional relationship with the original text
// Chunks can be independently embedded as vectors and retrieved, supporting precise content localization
type Chunk struct {
	// Unique identifier of the chunk, using UUID format
	ID string `json:"id"                       gorm:"type:varchar(36);primaryKey"`
	// Tenant ID, used for multi-tenant isolation
	TenantID uint64 `json:"tenant_id"`
	// ID of the parent knowledge, associated with the Knowledge model
	KnowledgeID string `json:"knowledge_id"`
	// ID of the knowledge base, for quick location
	KnowledgeBaseID string `json:"knowledge_base_id"`
	// Optional tag ID for categorization within a knowledge base (used for FAQ)
	TagID string `json:"tag_id"                   gorm:"type:varchar(36);index"`
	// Actual text content of the chunk
	Content string `json:"content"`
	// Index position of the chunk in the original document
	ChunkIndex int `json:"chunk_index"`
	// Whether the chunk is enabled, can be used to temporarily disable certain chunks
	IsEnabled bool `json:"is_enabled"               gorm:"default:true"`
	// Flags 存储多个布尔状态的位标志（如推荐状态等）
	// 默认值为 ChunkFlagRecommended (1)，表示默认可推荐
	Flags ChunkFlags `json:"flags"                    gorm:"default:1"`
	// Status of the chunk
	Status int `json:"status"                   gorm:"default:0"`
	// Starting character position in the original text
	StartAt int `json:"start_at"`
	// Ending character position in the original text
	EndAt int `json:"end_at"`
	// Previous chunk ID
	PreChunkID string `json:"pre_chunk_id"`
	// Next chunk ID
	NextChunkID string `json:"next_chunk_id"`
	// Chunk 类型，用于区分不同类型的 Chunk
	ChunkType ChunkType `json:"chunk_type"               gorm:"type:varchar(20);default:'text'"`
	// 父 Chunk ID，用于关联图片 Chunk 和原始文本 Chunk
	ParentChunkID string `json:"parent_chunk_id"          gorm:"type:varchar(36);index"`
	// 关系 Chunk ID，用于关联关系 Chunk 和原始文本 Chunk
	RelationChunks JSON `json:"relation_chunks"          gorm:"type:json"`
	// 间接关系 Chunk ID，用于关联间接关系 Chunk 和原始文本 Chunk
	IndirectRelationChunks JSON `json:"indirect_relation_chunks" gorm:"type:json"`
	// Metadata 存储 chunk 级别的扩展信息，例如 FAQ 元数据
	Metadata JSON `json:"metadata"                 gorm:"type:json"`
	// ContentHash 存储内容的 hash 值，用于快速匹配（主要用于 FAQ）
	ContentHash string `json:"content_hash"             gorm:"type:varchar(64);index"`
	// 图片信息，存储为 JSON
	ImageInfo string `json:"image_info"               gorm:"type:text"`
	// Chunk creation time
	CreatedAt time.Time `json:"created_at"`
	// Chunk last update time
	UpdatedAt time.Time `json:"updated_at"`
	// Soft delete marker, supports data recovery
	DeletedAt gorm.DeletedAt `json:"deleted_at"               gorm:"index"`
}
