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

// FAQHandler handles FAQ knowledge base operations.
type FAQHandler struct {
	knowledgeService interfaces.KnowledgeService
}

// NewFAQHandler creates a new FAQ handler
func NewFAQHandler(knowledgeService interfaces.KnowledgeService) *FAQHandler {
	return &FAQHandler{knowledgeService: knowledgeService}
}

// ListEntries godoc
// @Summary      Get FAQ Entry List
// @Description  Get FAQ entry list under knowledge base with pagination and filtering support
// @Tags         FAQ Management
// @Accept       json
// @Produce      json
// @Param        id           path      string  true   "Knowledge Base ID"
// @Param        page         query     int     false  "Page number"
// @Param        page_size    query     int     false  "Page size"
// @Param        tag_id       query     string  false  "Tag ID filter"
// @Param        keyword      query     string  false  "Keyword search"
// @Param        search_field query     string  false  "Search field: standard_question(standard question), similar_questions(similar questions), answers(answers), default search all"
// @Param        sort_order   query     string  false  "Sort order: asc(sort by update time ascending), default sort by update time descending"
// @Success      200        {object}  map[string]interface{}  "FAQ list"
// @Failure      400        {object}  errors.AppError         "Invalid request parameters"
// @Security     Bearer
// @Security     ApiKeyAuth
// @Router       /knowledge-bases/{id}/faq/entries [get]
func (h *FAQHandler) ListEntries(c *gin.Context) {
	ctx := c.Request.Context()
	var page types.Pagination
	if err := c.ShouldBindQuery(&page); err != nil {
		logger.Error(ctx, "Failed to bind pagination query", err)
		c.Error(errors.NewBadRequestError("Invalid pagination parameters").WithDetails(err.Error()))
		return
	}

	tagID := secutils.SanitizeForLog(c.Query("tag_id"))
	keyword := secutils.SanitizeForLog(c.Query("keyword"))
	searchField := secutils.SanitizeForLog(c.Query("search_field"))
	sortOrder := secutils.SanitizeForLog(c.Query("sort_order"))

	result, err := h.knowledgeService.ListFAQEntries(ctx, secutils.SanitizeForLog(c.Param("id")), &page, tagID, keyword, searchField, sortOrder)
	if err != nil {
		logger.ErrorWithFields(ctx, err, nil)
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// UpsertEntries godoc
// @Summary      Batch Upsert FAQ Entries
// @Description  Asynchronously batch update or insert FAQ entries
// @Tags         FAQ Management
// @Accept       json
// @Produce      json
// @Param        id       path      string                    true  "Knowledge Base ID"
// @Param        request  body      types.FAQBatchUpsertPayload  true  "Batch operation request"
// @Success      200      {object}  map[string]interface{}    "Task ID"
// @Failure      400      {object}  errors.AppError           "Invalid request parameters"
// @Security     Bearer
// @Security     ApiKeyAuth
// @Router       /knowledge-bases/{id}/faq/entries [post]
func (h *FAQHandler) UpsertEntries(c *gin.Context) {
	ctx := c.Request.Context()
	var req types.FAQBatchUpsertPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error(ctx, "Failed to bind FAQ upsert payload", err)
		c.Error(errors.NewBadRequestError("Invalid request parameters").WithDetails(err.Error()))
		return
	}

	taskID, err := h.knowledgeService.UpsertFAQEntries(ctx, secutils.SanitizeForLog(c.Param("id")), &req)
	if err != nil {
		logger.ErrorWithFields(ctx, err, nil)
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"task_id": taskID,
		},
	})
}

// CreateEntry godoc
// @Summary      Create Single FAQ Entry
// @Description  Synchronously create a single FAQ entry
// @Tags         FAQ Management
// @Accept       json
// @Produce      json
// @Param        id       path      string                true  "Knowledge Base ID"
// @Param        request  body      types.FAQEntryPayload true  "FAQ entry"
// @Success      200      {object}  map[string]interface{}  "Created FAQ entry"
// @Failure      400      {object}  errors.AppError         "Invalid request parameters"
// @Security     Bearer
// @Security     ApiKeyAuth
// @Router       /knowledge-bases/{id}/faq/entry [post]
func (h *FAQHandler) CreateEntry(c *gin.Context) {
	ctx := c.Request.Context()
	var req types.FAQEntryPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error(ctx, "Failed to bind FAQ entry payload", err)
		c.Error(errors.NewBadRequestError("Invalid request parameters").WithDetails(err.Error()))
		return
	}

	entry, err := h.knowledgeService.CreateFAQEntry(ctx, secutils.SanitizeForLog(c.Param("id")), &req)
	if err != nil {
		logger.ErrorWithFields(ctx, err, nil)
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    entry,
	})
}

// UpdateEntry godoc
// @Summary      Update FAQ Entry
// @Description  Update specified FAQ entry
// @Tags         FAQ Management
// @Accept       json
// @Produce      json
// @Param        id        path      string                true  "Knowledge Base ID"
// @Param        entry_id  path      string                true  "FAQ Entry ID"
// @Param        request   body      types.FAQEntryPayload true  "FAQ entry"
// @Success      200       {object}  map[string]interface{}  "Update successful"
// @Failure      400       {object}  errors.AppError         "Invalid request parameters"
// @Security     Bearer
// @Security     ApiKeyAuth
// @Router       /knowledge-bases/{id}/faq/entries/{entry_id} [put]
func (h *FAQHandler) UpdateEntry(c *gin.Context) {
	ctx := c.Request.Context()
	var req types.FAQEntryPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error(ctx, "Failed to bind FAQ entry payload", err)
		c.Error(errors.NewBadRequestError("Invalid request parameters").WithDetails(err.Error()))
		return
	}

	if err := h.knowledgeService.UpdateFAQEntry(ctx,
		secutils.SanitizeForLog(c.Param("id")), secutils.SanitizeForLog(c.Param("entry_id")), &req); err != nil {
		logger.ErrorWithFields(ctx, err, nil)
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

// UpdateEntryTagBatch godoc
// @Summary      Batch Update FAQ Tags
// @Description  Batch update tags for FAQ entries
// @Tags         FAQ Management
// @Accept       json
// @Produce      json
// @Param        id       path      string  true  "Knowledge Base ID"
// @Param        request  body      object  true  "Tag update request"
// @Success      200      {object}  map[string]interface{}  "Update successful"
// @Failure      400      {object}  errors.AppError         "Invalid request parameters"
// @Security     Bearer
// @Security     ApiKeyAuth
// @Router       /knowledge-bases/{id}/faq/entries/tags [put]
func (h *FAQHandler) UpdateEntryTagBatch(c *gin.Context) {
	ctx := c.Request.Context()
	var req faqEntryTagBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error(ctx, "Failed to bind FAQ entry tag batch payload", err)
		c.Error(errors.NewBadRequestError("Invalid request parameters").WithDetails(err.Error()))
		return
	}
	if err := h.knowledgeService.UpdateFAQEntryTagBatch(ctx,
		secutils.SanitizeForLog(c.Param("id")), req.Updates); err != nil {
		logger.ErrorWithFields(ctx, err, nil)
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

// UpdateEntryFieldsBatch godoc
// @Summary      Batch Update FAQ Fields
// @Description  Batch update multiple fields for FAQ entries (is_enabled, is_recommended, tag_id)
// @Tags         FAQ Management
// @Accept       json
// @Produce      json
// @Param        id       path      string                        true  "Knowledge Base ID"
// @Param        request  body      types.FAQEntryFieldsBatchUpdate  true  "Field update request"
// @Success      200      {object}  map[string]interface{}        "Update successful"
// @Failure      400      {object}  errors.AppError               "Invalid request parameters"
// @Security     Bearer
// @Security     ApiKeyAuth
// @Router       /knowledge-bases/{id}/faq/entries/fields [put]
func (h *FAQHandler) UpdateEntryFieldsBatch(c *gin.Context) {
	ctx := c.Request.Context()
	var req types.FAQEntryFieldsBatchUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error(ctx, "Failed to bind FAQ entry fields batch payload", err)
		c.Error(errors.NewBadRequestError("Invalid request parameters").WithDetails(err.Error()))
		return
	}
	if err := h.knowledgeService.UpdateFAQEntryFieldsBatch(ctx,
		secutils.SanitizeForLog(c.Param("id")), &req); err != nil {
		logger.ErrorWithFields(ctx, err, nil)
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

// faqDeleteRequest is a request for deleting FAQ entries in batch
type faqDeleteRequest struct {
	IDs []string `json:"ids" binding:"required,min=1,dive,required"`
}

// faqEntryTagBatchRequest is a request for updating tags for FAQ entries in batch
type faqEntryTagBatchRequest struct {
	Updates map[string]*string `json:"updates" binding:"required,min=1"`
}

// DeleteEntries godoc
// @Summary      Batch Delete FAQ Entries
// @Description  Batch delete specified FAQ entries
// @Tags         FAQ Management
// @Accept       json
// @Produce      json
// @Param        id       path      string  true  "Knowledge Base ID"
// @Param        request  body      object{ids=[]string}  true  "FAQ ID list to delete"
// @Success      200      {object}  map[string]interface{}  "Delete successful"
// @Failure      400      {object}  errors.AppError         "Invalid request parameters"
// @Security     Bearer
// @Security     ApiKeyAuth
// @Router       /knowledge-bases/{id}/faq/entries [delete]
func (h *FAQHandler) DeleteEntries(c *gin.Context) {
	ctx := c.Request.Context()
	var req faqDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf(ctx, "Failed to bind FAQ delete payload: %s", secutils.SanitizeForLog(err.Error()))
		c.Error(errors.NewBadRequestError("Invalid request parameters").WithDetails(err.Error()))
		return
	}

	if err := h.knowledgeService.DeleteFAQEntries(ctx,
		secutils.SanitizeForLog(c.Param("id")),
		secutils.SanitizeForLogArray(req.IDs)); err != nil {
		logger.ErrorWithFields(ctx, err, nil)
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

// SearchFAQ godoc
// @Summary      Search FAQ
// @Description  Search in FAQ using hybrid search, supports two-level priority tag recall: first_priority_tag_ids has highest priority, second_priority_tag_ids has second priority
// @Tags         FAQ Management
// @Accept       json
// @Produce      json
// @Param        id       path      string                true  "Knowledge Base ID"
// @Param        request  body      types.FAQSearchRequest  true  "Search request"
// @Success      200      {object}  map[string]interface{}  "Search results"
// @Failure      400      {object}  errors.AppError         "Invalid request parameters"
// @Security     Bearer
// @Security     ApiKeyAuth
// @Router       /knowledge-bases/{id}/faq/search [post]
func (h *FAQHandler) SearchFAQ(c *gin.Context) {
	ctx := c.Request.Context()
	var req types.FAQSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error(ctx, "Failed to bind FAQ search payload", err)
		c.Error(errors.NewBadRequestError("Invalid request parameters").WithDetails(err.Error()))
		return
	}
	req.QueryText = secutils.SanitizeForLog(req.QueryText)
	if req.MatchCount <= 0 {
		req.MatchCount = 10
	}
	if req.MatchCount > 200 {
		req.MatchCount = 200
	}
	entries, err := h.knowledgeService.SearchFAQEntries(ctx, secutils.SanitizeForLog(c.Param("id")), &req)
	if err != nil {
		logger.ErrorWithFields(ctx, err, nil)
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    entries,
	})
}

// ExportEntries godoc
// @Summary      Export FAQ entries
// @Description  Export all FAQ entries as CSV file
// @Tags         FAQ Management
// @Accept       json
// @Produce      text/csv
// @Param        id   path      string  true  "Knowledge base ID"
// @Success      200  {file}    file    "CSV file"
// @Failure      400  {object}  errors.AppError  "Invalid request parameters"
// @Security     Bearer
// @Security     ApiKeyAuth
// @Router       /knowledge-bases/{id}/faq/entries/export [get]
func (h *FAQHandler) ExportEntries(c *gin.Context) {
	ctx := c.Request.Context()
	kbID := secutils.SanitizeForLog(c.Param("id"))

	csvData, err := h.knowledgeService.ExportFAQEntries(ctx, kbID)
	if err != nil {
		logger.ErrorWithFields(ctx, err, nil)
		c.Error(err)
		return
	}

	// Set response headers for CSV download
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=faq_export.csv")
	// Add BOM for Excel compatibility with UTF-8
	bom := []byte{0xEF, 0xBB, 0xBF}
	c.Data(http.StatusOK, "text/csv; charset=utf-8", append(bom, csvData...))
}

// GetEntry godoc
// @Summary      Get FAQ entry details
// @Description  Get details of a single FAQ entry by ID
// @Tags         FAQ Management
// @Accept       json
// @Produce      json
// @Param        id        path      string  true  "Knowledge base ID"
// @Param        entry_id  path      string  true  "FAQ entry ID"
// @Success      200       {object}  map[string]interface{}  "FAQ entry details"
// @Failure      400       {object}  errors.AppError         "Invalid request parameters"
// @Failure      404       {object}  errors.AppError         "Entry not found"
// @Security     Bearer
// @Security     ApiKeyAuth
// @Router       /knowledge-bases/{id}/faq/entries/{entry_id} [get]
func (h *FAQHandler) GetEntry(c *gin.Context) {
	ctx := c.Request.Context()
	kbID := secutils.SanitizeForLog(c.Param("id"))
	entryID := secutils.SanitizeForLog(c.Param("entry_id"))

	entry, err := h.knowledgeService.GetFAQEntry(ctx, kbID, entryID)
	if err != nil {
		logger.ErrorWithFields(ctx, err, nil)
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    entry,
	})
}

// GetImportProgress godoc
// @Summary      Get FAQ import progress
// @Description  Get progress of FAQ import task
// @Tags         FAQ Management
// @Accept       json
// @Produce      json
// @Param        task_id  path      string  true  "Task ID"
// @Success      200      {object}  map[string]interface{}  "Import progress"
// @Failure      404      {object}  errors.AppError         "Task not found"
// @Security     Bearer
// @Security     ApiKeyAuth
// @Router       /faq/import/progress/{task_id} [get]
func (h *FAQHandler) GetImportProgress(c *gin.Context) {
	ctx := c.Request.Context()
	taskID := secutils.SanitizeForLog(c.Param("task_id"))

	progress, err := h.knowledgeService.GetFAQImportProgress(ctx, taskID)
	if err != nil {
		logger.ErrorWithFields(ctx, err, nil)
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    progress,
	})
}
