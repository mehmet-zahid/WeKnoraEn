package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Tencent/WeKnora/internal/errors"
	"github.com/Tencent/WeKnora/internal/logger"
	"github.com/Tencent/WeKnora/internal/types"
	"github.com/Tencent/WeKnora/internal/types/interfaces"
	secutils "github.com/Tencent/WeKnora/internal/utils"
)

// TagHandler handles knowledge base tag operations.
type TagHandler struct {
	tagService interfaces.KnowledgeTagService
}

// DeleteTagRequest represents the request body for deleting a tag
type DeleteTagRequest struct {
	ExcludeIDs []string `json:"exclude_ids"` // Chunk IDs to exclude from deletion
}

// NewTagHandler creates a new TagHandler.
func NewTagHandler(tagService interfaces.KnowledgeTagService) *TagHandler {
	return &TagHandler{tagService: tagService}
}

// ListTags godoc
// @Summary      Get Tag List
// @Description  Get all tags and statistics under knowledge base
// @Tags         Tag Management
// @Accept       json
// @Produce      json
// @Param        id         path      string  true   "Knowledge Base ID"
// @Param        page       query     int     false  "Page number"
// @Param        page_size  query     int     false  "Page size"
// @Param        keyword    query     string  false  "Keyword search"
// @Success      200        {object}  map[string]interface{}  "Tag list"
// @Failure      400        {object}  errors.AppError         "Invalid request parameters"
// @Security     Bearer
// @Security     ApiKeyAuth
// @Router       /knowledge-bases/{id}/tags [get]
func (h *TagHandler) ListTags(c *gin.Context) {
	ctx := c.Request.Context()
	kbID := secutils.SanitizeForLog(c.Param("id"))

	var page types.Pagination
	if err := c.ShouldBindQuery(&page); err != nil {
		logger.Error(ctx, "Failed to bind pagination query", err)
		c.Error(errors.NewBadRequestError("Invalid pagination parameters").WithDetails(err.Error()))
		return
	}

	keyword := secutils.SanitizeForLog(c.Query("keyword"))

	tags, err := h.tagService.ListTags(ctx, kbID, &page, keyword)
	if err != nil {
		logger.ErrorWithFields(ctx, err, nil)
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    tags,
	})
}

type createTagRequest struct {
	Name      string `json:"name"       binding:"required"`
	Color     string `json:"color"`
	SortOrder int    `json:"sort_order"`
}

// CreateTag godoc
// @Summary      Create Tag
// @Description  Create new tag under knowledge base
// @Tags         Tag Management
// @Accept       json
// @Produce      json
// @Param        id       path      string  true  "Knowledge Base ID"
// @Param        request  body      object{name=string,color=string,sort_order=int}  true  "Tag information"
// @Success      200      {object}  map[string]interface{}  "Created tag"
// @Failure      400      {object}  errors.AppError         "Invalid request parameters"
// @Security     Bearer
// @Security     ApiKeyAuth
// @Router       /knowledge-bases/{id}/tags [post]
func (h *TagHandler) CreateTag(c *gin.Context) {
	ctx := c.Request.Context()
	kbID := secutils.SanitizeForLog(c.Param("id"))

	var req createTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error(ctx, "Failed to bind create tag payload", err)
		c.Error(errors.NewBadRequestError("Invalid request parameters").WithDetails(err.Error()))
		return
	}

	tag, err := h.tagService.CreateTag(ctx, kbID,
		secutils.SanitizeForLog(req.Name), secutils.SanitizeForLog(req.Color), req.SortOrder)
	if err != nil {
		logger.ErrorWithFields(ctx, err, map[string]interface{}{
			"kb_id": kbID,
		})
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    tag,
	})
}

type updateTagRequest struct {
	Name      *string `json:"name"`
	Color     *string `json:"color"`
	SortOrder *int    `json:"sort_order"`
}

// UpdateTag godoc
// @Summary      Update Tag
// @Description  Update tag information
// @Tags         Tag Management
// @Accept       json
// @Produce      json
// @Param        id       path      string  true  "Knowledge Base ID"
// @Param        tag_id   path      string  true  "Tag ID"
// @Param        request  body      object  true  "Tag update information"
// @Success      200      {object}  map[string]interface{}  "Updated tag"
// @Failure      400      {object}  errors.AppError         "Invalid request parameters"
// @Security     Bearer
// @Security     ApiKeyAuth
// @Router       /knowledge-bases/{id}/tags/{tag_id} [put]
func (h *TagHandler) UpdateTag(c *gin.Context) {
	ctx := c.Request.Context()

	tagID := secutils.SanitizeForLog(c.Param("tag_id"))
	var req updateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error(ctx, "Failed to bind update tag payload", err)
		c.Error(errors.NewBadRequestError("Invalid request parameters").WithDetails(err.Error()))
		return
	}

	tag, err := h.tagService.UpdateTag(ctx, tagID, req.Name, req.Color, req.SortOrder)
	if err != nil {
		logger.ErrorWithFields(ctx, err, map[string]interface{}{
			"tag_id": tagID,
		})
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    tag,
	})
}

// DeleteTag godoc
// @Summary      Delete Tag
// @Description  Delete tag, can use force=true to force delete referenced tags, content_only=true to only delete tag content while keeping the tag itself
// @Tags         Tag Management
// @Accept       json
// @Produce      json
// @Param        id            path      string              true   "Knowledge Base ID"
// @Param        tag_id        path      string              true   "Tag ID"
// @Param        force         query     bool                false  "Force delete"
// @Param        content_only  query     bool                false  "Delete content only, keep tag"
// @Param        body          body      DeleteTagRequest    false  "Delete options"
// @Success      200           {object}  map[string]interface{}  "Delete successful"
// @Failure      400           {object}  errors.AppError         "Invalid request parameters"
// @Security     Bearer
// @Security     ApiKeyAuth
// @Router       /knowledge-bases/{id}/tags/{tag_id} [delete]
func (h *TagHandler) DeleteTag(c *gin.Context) {
	ctx := c.Request.Context()
	tagID := secutils.SanitizeForLog(c.Param("tag_id"))

	force := c.Query("force") == "true"
	contentOnly := c.Query("content_only") == "true"

	var req DeleteTagRequest
	// Ignore bind error since body is optional
	_ = c.ShouldBindJSON(&req)

	if err := h.tagService.DeleteTag(ctx, tagID, force, contentOnly, req.ExcludeIDs); err != nil {
		logger.ErrorWithFields(ctx, err, map[string]interface{}{
			"tag_id": tagID,
		})
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

// NOTE: TagHandler currently exposes CRUD for tags and statistics.
// Knowledge / Chunk tagging is handled via dedicated knowledge and FAQ APIs.
