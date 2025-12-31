package chatpipline

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/Tencent/WeKnora/internal/types"
	"github.com/Tencent/WeKnora/internal/utils"
)

// PluginIntoChatMessage handles the transformation of search results into chat messages
type PluginIntoChatMessage struct{}

// NewPluginIntoChatMessage creates and registers a new PluginIntoChatMessage instance
func NewPluginIntoChatMessage(eventManager *EventManager) *PluginIntoChatMessage {
	res := &PluginIntoChatMessage{}
	eventManager.Register(res)
	return res
}

// ActivationEvents returns the event types this plugin handles
func (p *PluginIntoChatMessage) ActivationEvents() []types.EventType {
	return []types.EventType{types.INTO_CHAT_MESSAGE}
}

// OnEvent processes the INTO_CHAT_MESSAGE event to format chat message content
func (p *PluginIntoChatMessage) OnEvent(ctx context.Context,
	eventType types.EventType, chatManage *types.ChatManage, next func() *PluginError,
) *PluginError {
	pipelineInfo(ctx, "IntoChatMessage", "input", map[string]interface{}{
		"session_id":       chatManage.SessionID,
		"merge_result_cnt": len(chatManage.MergeResult),
		"template_len":     len(chatManage.SummaryConfig.ContextTemplate),
	})

	// Separate FAQ and document results when FAQ priority is enabled
	var faqResults, docResults []*types.SearchResult
	var hasHighConfidenceFAQ bool

	if chatManage.FAQPriorityEnabled {
		for _, result := range chatManage.MergeResult {
			if result.ChunkType == string(types.ChunkTypeFAQ) {
				faqResults = append(faqResults, result)
				// Check if this FAQ has high confidence (above direct answer threshold)
				if result.Score >= chatManage.FAQDirectAnswerThreshold && !hasHighConfidenceFAQ {
					hasHighConfidenceFAQ = true
					pipelineInfo(ctx, "IntoChatMessage", "high_confidence_faq", map[string]interface{}{
						"chunk_id":  result.ID,
						"score":     fmt.Sprintf("%.4f", result.Score),
						"threshold": chatManage.FAQDirectAnswerThreshold,
					})
				}
			} else {
				docResults = append(docResults, result)
			}
		}
		pipelineInfo(ctx, "IntoChatMessage", "faq_separation", map[string]interface{}{
			"faq_count":           len(faqResults),
			"doc_count":           len(docResults),
			"has_high_confidence": hasHighConfidenceFAQ,
		})
	}

	// 验证用户查询的安全性
	safeQuery, isValid := utils.ValidateInput(chatManage.Query)
	if !isValid {
		pipelineWarn(ctx, "IntoChatMessage", "invalid_query", map[string]interface{}{
			"session_id": chatManage.SessionID,
		})
		return ErrTemplateExecute.WithError(fmt.Errorf("User query contains invalid content"))
	}

	// Prepare weekday names
	weekdayName := []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}

	var contextsBuilder strings.Builder

	// Build contexts string based on FAQ priority strategy
	if chatManage.FAQPriorityEnabled && len(faqResults) > 0 {
		// Build structured context with FAQ prioritization
		contextsBuilder.WriteString("### Source 1: Standard Q&A Library (FAQ)\n")
		contextsBuilder.WriteString("[High Confidence - Please prioritize]\n")
		for i, result := range faqResults {
			passage := getEnrichedPassageForChat(ctx, result)
			if hasHighConfidenceFAQ && i == 0 {
				contextsBuilder.WriteString(fmt.Sprintf("[FAQ-%d] ⭐ Exact Match: %s\n", i+1, passage))
			} else {
				contextsBuilder.WriteString(fmt.Sprintf("[FAQ-%d] %s\n", i+1, passage))
			}
		}

		if len(docResults) > 0 {
			contextsBuilder.WriteString("\n### Source 2: Reference Documents\n")
			contextsBuilder.WriteString("[Supplementary Material - Only refer when FAQ cannot answer]\n")
			for i, result := range docResults {
				passage := getEnrichedPassageForChat(ctx, result)
				contextsBuilder.WriteString(fmt.Sprintf("[DOC-%d] %s\n", i+1, passage))
			}
		}
	} else {
		// Original behavior: simple numbered list
		passages := make([]string, len(chatManage.MergeResult))
		for i, result := range chatManage.MergeResult {
			passages[i] = getEnrichedPassageForChat(ctx, result)
		}
		for i, passage := range passages {
			if i > 0 {
				contextsBuilder.WriteString("\n\n")
			}
			contextsBuilder.WriteString(fmt.Sprintf("[%d] %s", i+1, passage))
		}
	}

	// Replace placeholders in context template
	userContent := chatManage.SummaryConfig.ContextTemplate
	userContent = strings.ReplaceAll(userContent, "{{query}}", safeQuery)
	userContent = strings.ReplaceAll(userContent, "{{contexts}}", contextsBuilder.String())
	userContent = strings.ReplaceAll(userContent, "{{current_time}}", time.Now().Format("2006-01-02 15:04:05"))
	userContent = strings.ReplaceAll(userContent, "{{current_week}}", weekdayName[time.Now().Weekday()])

	// Set formatted content back to chat management
	chatManage.UserContent = userContent
	pipelineInfo(ctx, "IntoChatMessage", "output", map[string]interface{}{
		"session_id":       chatManage.SessionID,
		"user_content_len": len(chatManage.UserContent),
		"faq_priority":     chatManage.FAQPriorityEnabled,
	})
	return next()
}

// getEnrichedPassageForChat merges Content and ImageInfo text content for chat message preparation
func getEnrichedPassageForChat(ctx context.Context, result *types.SearchResult) string {
	// If there's no image information, return content directly
	if result.Content == "" && result.ImageInfo == "" {
		return ""
	}

	// If there's only content, no image information
	if result.ImageInfo == "" {
		return result.Content
	}

	// Process image information and merge with content
	return enrichContentWithImageInfo(ctx, result.Content, result.ImageInfo)
}

// Regular expression for matching Markdown image links
var markdownImageRegex = regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+)\)`)

// enrichContentWithImageInfo merges image information with text content
func enrichContentWithImageInfo(ctx context.Context, content string, imageInfoJSON string) string {
	// Parse ImageInfo
	var imageInfos []types.ImageInfo
	err := json.Unmarshal([]byte(imageInfoJSON), &imageInfos)
	if err != nil {
		pipelineWarn(ctx, "IntoChatMessage", "image_parse_error", map[string]interface{}{
			"error": err.Error(),
		})
		return content
	}

	if len(imageInfos) == 0 {
		return content
	}

	// Create mapping from image URL to information
	imageInfoMap := make(map[string]*types.ImageInfo)
	for i := range imageInfos {
		if imageInfos[i].URL != "" {
			imageInfoMap[imageInfos[i].URL] = &imageInfos[i]
		}
		// 同时检查原始URL
		if imageInfos[i].OriginalURL != "" {
			imageInfoMap[imageInfos[i].OriginalURL] = &imageInfos[i]
		}
	}

	// 查找内容中的所有Markdown图片链接
	matches := markdownImageRegex.FindAllStringSubmatch(content, -1)

	// 用于存储已处理的图片URL
	processedURLs := make(map[string]bool)

	pipelineInfo(ctx, "IntoChatMessage", "image_markdown_links", map[string]interface{}{
		"match_count": len(matches),
	})

	// 替换每个图片链接，添加描述和OCR文本
	for _, match := range matches {
		if len(match) < 3 {
			continue
		}

		// 提取图片URL，忽略alt文本
		imgURL := match[2]

		// 标记该URL已处理
		processedURLs[imgURL] = true

		// 查找匹配的图片信息
		imgInfo, found := imageInfoMap[imgURL]

		// 如果找到匹配的图片信息，添加描述和OCR文本
		if found && imgInfo != nil {
			replacement := match[0] + "\n"
			if imgInfo.Caption != "" {
				replacement += fmt.Sprintf("图片描述: %s\n", imgInfo.Caption)
			}
			if imgInfo.OCRText != "" {
				replacement += fmt.Sprintf("图片文本: %s\n", imgInfo.OCRText)
			}
			content = strings.Replace(content, match[0], replacement, 1)
		}
	}

	// 处理未在内容中找到但存在于ImageInfo中的图片
	var additionalImageTexts []string
	for _, imgInfo := range imageInfos {
		// 如果图片URL已经处理过，跳过
		if processedURLs[imgInfo.URL] || processedURLs[imgInfo.OriginalURL] {
			continue
		}

		var imgTexts []string
		if imgInfo.Caption != "" {
			imgTexts = append(imgTexts, fmt.Sprintf("图片 %s 的描述信息: %s", imgInfo.URL, imgInfo.Caption))
		}
		if imgInfo.OCRText != "" {
			imgTexts = append(imgTexts, fmt.Sprintf("图片 %s 的文本: %s", imgInfo.URL, imgInfo.OCRText))
		}

		if len(imgTexts) > 0 {
			additionalImageTexts = append(additionalImageTexts, imgTexts...)
		}
	}

	// 如果有额外的图片信息，添加到内容末尾
	if len(additionalImageTexts) > 0 {
		if content != "" {
			content += "\n\n"
		}
		content += "附加图片信息:\n" + strings.Join(additionalImageTexts, "\n")
	}

	pipelineInfo(ctx, "IntoChatMessage", "image_enrich_summary", map[string]interface{}{
		"markdown_images": len(matches),
		"additional_imgs": len(additionalImageTexts),
	})

	return content
}
