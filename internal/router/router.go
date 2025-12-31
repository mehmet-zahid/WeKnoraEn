package router

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/dig"

	"github.com/Tencent/WeKnora/internal/config"
	"github.com/Tencent/WeKnora/internal/handler"
	"github.com/Tencent/WeKnora/internal/handler/session"
	"github.com/Tencent/WeKnora/internal/middleware"
	"github.com/Tencent/WeKnora/internal/types/interfaces"

	_ "github.com/Tencent/WeKnora/docs" // swagger docs
)

// RouterParams represents router parameters
type RouterParams struct {
	dig.In

	Config                *config.Config
	UserService           interfaces.UserService
	KBService             interfaces.KnowledgeBaseService
	KnowledgeService      interfaces.KnowledgeService
	ChunkService          interfaces.ChunkService
	SessionService        interfaces.SessionService
	MessageService        interfaces.MessageService
	ModelService          interfaces.ModelService
	EvaluationService     interfaces.EvaluationService
	KBHandler             *handler.KnowledgeBaseHandler
	KnowledgeHandler      *handler.KnowledgeHandler
	TenantHandler         *handler.TenantHandler
	TenantService         interfaces.TenantService
	ChunkHandler          *handler.ChunkHandler
	SessionHandler        *session.Handler
	MessageHandler        *handler.MessageHandler
	ModelHandler          *handler.ModelHandler
	EvaluationHandler     *handler.EvaluationHandler
	AuthHandler           *handler.AuthHandler
	InitializationHandler *handler.InitializationHandler
	SystemHandler         *handler.SystemHandler
	MCPServiceHandler     *handler.MCPServiceHandler
	WebSearchHandler      *handler.WebSearchHandler
	FAQHandler            *handler.FAQHandler
	TagHandler            *handler.TagHandler
	CustomAgentHandler    *handler.CustomAgentHandler
}

// NewRouter creates a new router
func NewRouter(params RouterParams) *gin.Engine {
	r := gin.New()

	// CORS middleware should be placed first
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-API-Key", "X-Request-ID"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Basic middleware (no authentication required)
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.ErrorHandler())

	// Health check (no authentication required)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Swagger API documentation (only enabled in non-production environments)
	// Determined by GIN_MODE environment variable: Swagger disabled in release mode
	if gin.Mode() != gin.ReleaseMode {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
			ginSwagger.DefaultModelsExpandDepth(-1), // Default collapse Models
			ginSwagger.DocExpansion("list"),         // Expansion mode: "list"(expand tags), "full"(expand all), "none"(collapse all)
			ginSwagger.DeepLinking(true),            // Enable deep linking
			ginSwagger.PersistAuthorization(true),   // Persist authorization information
		))
	}

	// Authentication middleware
	r.Use(middleware.Auth(params.TenantService, params.UserService, params.Config))

	// Add OpenTelemetry tracing middleware
	r.Use(middleware.TracingMiddleware())

	// API routes requiring authentication
	v1 := r.Group("/api/v1")
	{
		RegisterAuthRoutes(v1, params.AuthHandler)
		RegisterTenantRoutes(v1, params.TenantHandler)
		RegisterKnowledgeBaseRoutes(v1, params.KBHandler)
		RegisterKnowledgeTagRoutes(v1, params.TagHandler)
		RegisterKnowledgeRoutes(v1, params.KnowledgeHandler)
		RegisterFAQRoutes(v1, params.FAQHandler)
		RegisterChunkRoutes(v1, params.ChunkHandler)
		RegisterSessionRoutes(v1, params.SessionHandler)
		RegisterChatRoutes(v1, params.SessionHandler)
		RegisterMessageRoutes(v1, params.MessageHandler)
		RegisterModelRoutes(v1, params.ModelHandler)
		RegisterEvaluationRoutes(v1, params.EvaluationHandler)
		RegisterInitializationRoutes(v1, params.InitializationHandler)
		RegisterSystemRoutes(v1, params.SystemHandler)
		RegisterMCPServiceRoutes(v1, params.MCPServiceHandler)
		RegisterWebSearchRoutes(v1, params.WebSearchHandler)
		RegisterCustomAgentRoutes(v1, params.CustomAgentHandler)
	}

	return r
}

// RegisterChunkRoutes registers chunk-related routes
func RegisterChunkRoutes(r *gin.RouterGroup, handler *handler.ChunkHandler) {
	// Chunk route group
	chunks := r.Group("/chunks")
	{
		// Get chunk list
		chunks.GET("/:knowledge_id", handler.ListKnowledgeChunks)
		// Get single chunk by chunk_id (knowledge_id not required)
		chunks.GET("/by-id/:id", handler.GetChunkByIDOnly)
		// Delete chunk
		chunks.DELETE("/:knowledge_id/:id", handler.DeleteChunk)
		// Delete all chunks under knowledge
		chunks.DELETE("/:knowledge_id", handler.DeleteChunksByKnowledgeID)
		// Update chunk information
		chunks.PUT("/:knowledge_id/:id", handler.UpdateChunk)
		// Delete single generated question (by question ID)
		chunks.DELETE("/by-id/:id/questions", handler.DeleteGeneratedQuestion)
	}
}

// RegisterKnowledgeRoutes registers knowledge-related routes
func RegisterKnowledgeRoutes(r *gin.RouterGroup, handler *handler.KnowledgeHandler) {
	// Knowledge routes under knowledge base
	kb := r.Group("/knowledge-bases/:id/knowledge")
	{
		// Create knowledge from file
		kb.POST("/file", handler.CreateKnowledgeFromFile)
		// Create knowledge from URL
		kb.POST("/url", handler.CreateKnowledgeFromURL)
		// Manual Markdown entry
		kb.POST("/manual", handler.CreateManualKnowledge)
		// Get knowledge list under knowledge base
		kb.GET("", handler.ListKnowledge)
	}

	// Knowledge route group
	k := r.Group("/knowledge")
	{
		// Batch get knowledge
		k.GET("/batch", handler.GetKnowledgeBatch)
		// Get knowledge details
		k.GET("/:id", handler.GetKnowledge)
		// 删除知识
		k.DELETE("/:id", handler.DeleteKnowledge)
		// 更新知识
		k.PUT("/:id", handler.UpdateKnowledge)
		// 更新手工 Markdown 知识
		k.PUT("/manual/:id", handler.UpdateManualKnowledge)
		// 获取知识文件
		k.GET("/:id/download", handler.DownloadKnowledgeFile)
		// 更新图像分块信息
		k.PUT("/image/:id/:chunk_id", handler.UpdateImageInfo)
		// 批量更新知识标签
		k.PUT("/tags", handler.UpdateKnowledgeTagBatch)
		// 搜索知识
		k.GET("/search", handler.SearchKnowledge)
	}
}

// RegisterFAQRoutes 注册 FAQ 相关路由
func RegisterFAQRoutes(r *gin.RouterGroup, handler *handler.FAQHandler) {
	if handler == nil {
		return
	}
	faq := r.Group("/knowledge-bases/:id/faq")
	{
		faq.GET("/entries", handler.ListEntries)
		faq.GET("/entries/export", handler.ExportEntries)
		faq.GET("/entries/:entry_id", handler.GetEntry)
		faq.POST("/entries", handler.UpsertEntries)
		faq.POST("/entry", handler.CreateEntry)
		faq.PUT("/entries/:entry_id", handler.UpdateEntry)
		// Unified batch update API - supports is_enabled, is_recommended, tag_id
		faq.PUT("/entries/fields", handler.UpdateEntryFieldsBatch)
		faq.PUT("/entries/tags", handler.UpdateEntryTagBatch)
		faq.DELETE("/entries", handler.DeleteEntries)
		faq.POST("/search", handler.SearchFAQ)
	}
	// FAQ import progress route (outside of knowledge-base scope)
	faqImport := r.Group("/faq/import")
	{
		faqImport.GET("/progress/:task_id", handler.GetImportProgress)
	}
}

// RegisterKnowledgeBaseRoutes 注册知识库相关的路由
func RegisterKnowledgeBaseRoutes(r *gin.RouterGroup, handler *handler.KnowledgeBaseHandler) {
	// 知识库路由组
	kb := r.Group("/knowledge-bases")
	{
		// 创建知识库
		kb.POST("", handler.CreateKnowledgeBase)
		// 获取知识库列表
		kb.GET("", handler.ListKnowledgeBases)
		// 获取知识库详情
		kb.GET("/:id", handler.GetKnowledgeBase)
		// 更新知识库
		kb.PUT("/:id", handler.UpdateKnowledgeBase)
		// 删除知识库
		kb.DELETE("/:id", handler.DeleteKnowledgeBase)
		// 混合搜索
		kb.GET("/:id/hybrid-search", handler.HybridSearch)
		// 拷贝知识库
		kb.POST("/copy", handler.CopyKnowledgeBase)
		// 获取知识库复制进度
		kb.GET("/copy/progress/:task_id", handler.GetKBCloneProgress)
	}
}

// RegisterKnowledgeTagRoutes 注册知识库标签相关路由
func RegisterKnowledgeTagRoutes(r *gin.RouterGroup, tagHandler *handler.TagHandler) {
	if tagHandler == nil {
		return
	}
	kbTags := r.Group("/knowledge-bases/:id/tags")
	{
		kbTags.GET("", tagHandler.ListTags)
		kbTags.POST("", tagHandler.CreateTag)
		kbTags.PUT("/:tag_id", tagHandler.UpdateTag)
		kbTags.DELETE("/:tag_id", tagHandler.DeleteTag)
	}
}

// RegisterMessageRoutes 注册消息相关的路由
func RegisterMessageRoutes(r *gin.RouterGroup, handler *handler.MessageHandler) {
	// 消息路由组
	messages := r.Group("/messages")
	{
		// 加载更早的消息，用于向上滚动加载
		messages.GET("/:session_id/load", handler.LoadMessages)
		// 删除消息
		messages.DELETE("/:session_id/:id", handler.DeleteMessage)
	}
}

// RegisterSessionRoutes 注册路由
func RegisterSessionRoutes(r *gin.RouterGroup, handler *session.Handler) {
	sessions := r.Group("/sessions")
	{
		sessions.POST("", handler.CreateSession)
		sessions.GET("/:id", handler.GetSession)
		sessions.GET("", handler.GetSessionsByTenant)
		sessions.PUT("/:id", handler.UpdateSession)
		sessions.DELETE("/:id", handler.DeleteSession)
		sessions.POST("/:session_id/generate_title", handler.GenerateTitle)
		sessions.POST("/:session_id/stop", handler.StopSession)
		// 继续接收活跃流
		sessions.GET("/continue-stream/:session_id", handler.ContinueStream)
	}
}

// RegisterChatRoutes 注册路由
func RegisterChatRoutes(r *gin.RouterGroup, handler *session.Handler) {
	knowledgeChat := r.Group("/knowledge-chat")
	{
		knowledgeChat.POST("/:session_id", handler.KnowledgeQA)
	}

	// Agent-based chat
	agentChat := r.Group("/agent-chat")
	{
		agentChat.POST("/:session_id", handler.AgentQA)
	}

	// 新增知识检索接口，不需要session_id
	knowledgeSearch := r.Group("/knowledge-search")
	{
		knowledgeSearch.POST("", handler.SearchKnowledge)
	}
}

// RegisterTenantRoutes 注册租户相关的路由
func RegisterTenantRoutes(r *gin.RouterGroup, handler *handler.TenantHandler) {
	// 添加获取所有租户的路由（需要跨租户权限）
	r.GET("/tenants/all", handler.ListAllTenants)
	// 添加搜索租户的路由（需要跨租户权限，支持分页和搜索）
	r.GET("/tenants/search", handler.SearchTenants)
	// 租户路由组
	tenantRoutes := r.Group("/tenants")
	{
		tenantRoutes.POST("", handler.CreateTenant)
		tenantRoutes.GET("/:id", handler.GetTenant)
		tenantRoutes.PUT("/:id", handler.UpdateTenant)
		tenantRoutes.DELETE("/:id", handler.DeleteTenant)
		tenantRoutes.GET("", handler.ListTenants)

		// Generic KV configuration management (tenant-level)
		// Tenant ID is obtained from authentication context
		tenantRoutes.GET("/kv/:key", handler.GetTenantKV)
		tenantRoutes.PUT("/kv/:key", handler.UpdateTenantKV)
	}
}

// RegisterModelRoutes 注册模型相关的路由
func RegisterModelRoutes(r *gin.RouterGroup, handler *handler.ModelHandler) {
	// 模型路由组
	models := r.Group("/models")
	{
		// 获取模型厂商列表
		models.GET("/providers", handler.ListModelProviders)
		// 创建模型
		models.POST("", handler.CreateModel)
		// 获取模型列表
		models.GET("", handler.ListModels)
		// 获取单个模型
		models.GET("/:id", handler.GetModel)
		// 更新模型
		models.PUT("/:id", handler.UpdateModel)
		// 删除模型
		models.DELETE("/:id", handler.DeleteModel)
	}
}

func RegisterEvaluationRoutes(r *gin.RouterGroup, handler *handler.EvaluationHandler) {
	evaluationRoutes := r.Group("/evaluation")
	{
		evaluationRoutes.POST("/", handler.Evaluation)
		evaluationRoutes.GET("/", handler.GetEvaluationResult)
	}
}

// RegisterAuthRoutes registers authentication routes
func RegisterAuthRoutes(r *gin.RouterGroup, handler *handler.AuthHandler) {
	r.POST("/auth/register", handler.Register)
	r.POST("/auth/login", handler.Login)
	r.POST("/auth/refresh", handler.RefreshToken)
	r.GET("/auth/validate", handler.ValidateToken)
	r.POST("/auth/logout", handler.Logout)
	r.GET("/auth/me", handler.GetCurrentUser)
	r.POST("/auth/change-password", handler.ChangePassword)
}

func RegisterInitializationRoutes(r *gin.RouterGroup, handler *handler.InitializationHandler) {
	// 初始化接口
	r.GET("/initialization/config/:kbId", handler.GetCurrentConfigByKB)
	r.POST("/initialization/initialize/:kbId", handler.InitializeByKB)
	r.PUT("/initialization/config/:kbId", handler.UpdateKBConfig) // 新的简化版接口，只传模型ID

	// Ollama相关接口
	r.GET("/initialization/ollama/status", handler.CheckOllamaStatus)
	r.GET("/initialization/ollama/models", handler.ListOllamaModels)
	r.POST("/initialization/ollama/models/check", handler.CheckOllamaModels)
	r.POST("/initialization/ollama/models/download", handler.DownloadOllamaModel)
	r.GET("/initialization/ollama/download/progress/:taskId", handler.GetDownloadProgress)
	r.GET("/initialization/ollama/download/tasks", handler.ListDownloadTasks)

	// 远程API相关接口
	r.POST("/initialization/remote/check", handler.CheckRemoteModel)
	r.POST("/initialization/embedding/test", handler.TestEmbeddingModel)
	r.POST("/initialization/rerank/check", handler.CheckRerankModel)
	r.POST("/initialization/multimodal/test", handler.TestMultimodalFunction)

	r.POST("/initialization/extract/text-relation", handler.ExtractTextRelations)
	r.POST("/initialization/extract/fabri-tag", handler.FabriTag)
	r.POST("/initialization/extract/fabri-text", handler.FabriText)
}

// RegisterSystemRoutes registers system information routes
func RegisterSystemRoutes(r *gin.RouterGroup, handler *handler.SystemHandler) {
	systemRoutes := r.Group("/system")
	{
		systemRoutes.GET("/info", handler.GetSystemInfo)
		systemRoutes.GET("/minio/buckets", handler.ListMinioBuckets)
	}
}

// RegisterMCPServiceRoutes registers MCP service routes
func RegisterMCPServiceRoutes(r *gin.RouterGroup, handler *handler.MCPServiceHandler) {
	mcpServices := r.Group("/mcp-services")
	{
		// Create MCP service
		mcpServices.POST("", handler.CreateMCPService)
		// List MCP services
		mcpServices.GET("", handler.ListMCPServices)
		// Get MCP service by ID
		mcpServices.GET("/:id", handler.GetMCPService)
		// Update MCP service
		mcpServices.PUT("/:id", handler.UpdateMCPService)
		// Delete MCP service
		mcpServices.DELETE("/:id", handler.DeleteMCPService)
		// Test MCP service connection
		mcpServices.POST("/:id/test", handler.TestMCPService)
		// Get MCP service tools
		mcpServices.GET("/:id/tools", handler.GetMCPServiceTools)
		// Get MCP service resources
		mcpServices.GET("/:id/resources", handler.GetMCPServiceResources)
	}
}

// RegisterWebSearchRoutes registers web search routes
func RegisterWebSearchRoutes(r *gin.RouterGroup, webSearchHandler *handler.WebSearchHandler) {
	// Web search providers
	webSearch := r.Group("/web-search")
	{
		// Get available providers
		webSearch.GET("/providers", webSearchHandler.GetProviders)
	}
}

// RegisterCustomAgentRoutes registers custom agent routes
func RegisterCustomAgentRoutes(r *gin.RouterGroup, agentHandler *handler.CustomAgentHandler) {
	agents := r.Group("/agents")
	{
		// Get placeholder definitions (must be before /:id to avoid conflict)
		agents.GET("/placeholders", agentHandler.GetPlaceholders)
		// Create custom agent
		agents.POST("", agentHandler.CreateAgent)
		// List all agents (including built-in)
		agents.GET("", agentHandler.ListAgents)
		// Get agent by ID
		agents.GET("/:id", agentHandler.GetAgent)
		// Update agent
		agents.PUT("/:id", agentHandler.UpdateAgent)
		// Delete agent
		agents.DELETE("/:id", agentHandler.DeleteAgent)
		// Copy agent
		agents.POST("/:id/copy", agentHandler.CopyAgent)
	}
}
