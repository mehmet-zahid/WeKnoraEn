# Event System Usage Examples

## Integrating Event System in Chat Pipeline

### 1. Setting Up Event Bus During Service Initialization

```go
// internal/container/container.go or main.go

import (
    "github.com/Tencent/WeKnora/internal/event"
)

func InitializeEventSystem() {
    // Get global event bus
    bus := event.GetGlobalEventBus()
    
    // Register monitoring handler
    event.NewMonitoringHandler(bus)
    
    // Register analytics handler
    event.NewAnalyticsHandler(bus)
    
    // Or register custom handler
    bus.On(event.EventQueryReceived, func(ctx context.Context, e event.Event) error {
        // Custom processing logic
        return nil
    })
}
```

### 2. Sending Events in Query Processing Service

#### Example: Adding Events in search.go

```go
// internal/application/service/chat_pipline/search.go

import (
    "github.com/Tencent/WeKnora/internal/event"
    "time"
)

func (p *PluginSearch) OnEvent(
    ctx context.Context,
    eventType types.EventType,
    chatManage *types.ChatManage,
    next func() *PluginError,
) *PluginError {
    // Send retrieval start event
    startTime := time.Now()
    event.Emit(ctx, event.NewEvent(event.EventRetrievalStart, event.RetrievalData{
        Query:           chatManage.ProcessedQuery,
        KnowledgeBaseID: chatManage.KnowledgeBaseID,
        TopK:            chatManage.EmbeddingTopK,
        RetrievalType:   "vector",
    }).WithSessionID(chatManage.SessionID))
    
    // Execute retrieval logic
    results, err := p.performSearch(ctx, chatManage)
    if err != nil {
        // Send error event
        event.Emit(ctx, event.NewEvent(event.EventError, event.ErrorData{
            Error:     err.Error(),
            Stage:     "retrieval",
            SessionID: chatManage.SessionID,
            Query:     chatManage.ProcessedQuery,
        }).WithSessionID(chatManage.SessionID))
        return ErrSearch.WithError(err)
    }
    
    // Send retrieval complete event
    event.Emit(ctx, event.NewEvent(event.EventRetrievalComplete, event.RetrievalData{
        Query:           chatManage.ProcessedQuery,
        KnowledgeBaseID: chatManage.KnowledgeBaseID,
        TopK:            chatManage.EmbeddingTopK,
        RetrievalType:   "vector",
        ResultCount:     len(results),
        Duration:        time.Since(startTime).Milliseconds(),
        Results:         results,
    }).WithSessionID(chatManage.SessionID))
    
    chatManage.SearchResult = results
    return next()
}
```

#### Example: Adding Events in rewrite.go

```go
// internal/application/service/chat_pipline/rewrite.go

func (p *PluginRewriteQuery) OnEvent(
    ctx context.Context,
    eventType types.EventType,
    chatManage *types.ChatManage,
    next func() *PluginError,
) *PluginError {
    // Send rewrite start event
    event.Emit(ctx, event.NewEvent(event.EventQueryRewrite, event.QueryData{
        OriginalQuery: chatManage.Query,
        SessionID:     chatManage.SessionID,
    }).WithSessionID(chatManage.SessionID))
    
    // Execute query rewriting
    rewrittenQuery, err := p.rewriteQuery(ctx, chatManage)
    if err != nil {
        return ErrRewrite.WithError(err)
    }
    
    // Send rewrite complete event
    event.Emit(ctx, event.NewEvent(event.EventQueryRewritten, event.QueryData{
        OriginalQuery:  chatManage.Query,
        RewrittenQuery: rewrittenQuery,
        SessionID:      chatManage.SessionID,
    }).WithSessionID(chatManage.SessionID))
    
    chatManage.RewriteQuery = rewrittenQuery
    return next()
}
```

#### Example: Adding Events in rerank.go

```go
// internal/application/service/chat_pipline/rerank.go

func (p *PluginRerank) OnEvent(
    ctx context.Context,
    eventType types.EventType,
    chatManage *types.ChatManage,
    next func() *PluginError,
) *PluginError {
    // Send rerank start event
    startTime := time.Now()
    inputCount := len(chatManage.SearchResult)
    
    event.Emit(ctx, event.NewEvent(event.EventRerankStart, event.RerankData{
        Query:      chatManage.ProcessedQuery,
        InputCount: inputCount,
        ModelID:    chatManage.RerankModelID,
    }).WithSessionID(chatManage.SessionID))
    
    // Execute reranking
    rerankResults, err := p.performRerank(ctx, chatManage)
    if err != nil {
        return ErrRerank.WithError(err)
    }
    
    // Send rerank complete event
    event.Emit(ctx, event.NewEvent(event.EventRerankComplete, event.RerankData{
        Query:       chatManage.ProcessedQuery,
        InputCount:  inputCount,
        OutputCount: len(rerankResults),
        ModelID:     chatManage.RerankModelID,
        Duration:    time.Since(startTime).Milliseconds(),
        Results:     rerankResults,
    }).WithSessionID(chatManage.SessionID))
    
    chatManage.RerankResult = rerankResults
    return next()
}
```

#### Example: Adding Events in chat_completion.go

```go
// internal/application/service/chat_pipline/chat_completion.go

func (p *PluginChatCompletion) OnEvent(
    ctx context.Context,
    eventType types.EventType,
    chatManage *types.ChatManage,
    next func() *PluginError,
) *PluginError {
    // Send chat start event
    startTime := time.Now()
    event.Emit(ctx, event.NewEvent(event.EventChatStart, event.ChatData{
        Query:    chatManage.Query,
        ModelID:  chatManage.ChatModelID,
        IsStream: false,
    }).WithSessionID(chatManage.SessionID))
    
    // Prepare model and messages
    chatModel, opt, err := prepareChatModel(ctx, p.modelService, chatManage)
    if err != nil {
        return ErrGetChatModel.WithError(err)
    }
    
    chatMessages := prepareMessagesWithHistory(chatManage)
    
    // Call model
    chatResponse, err := chatModel.Chat(ctx, chatMessages, opt)
    if err != nil {
        event.Emit(ctx, event.NewEvent(event.EventError, event.ErrorData{
            Error:     err.Error(),
            Stage:     "chat_completion",
            SessionID: chatManage.SessionID,
            Query:     chatManage.Query,
        }).WithSessionID(chatManage.SessionID))
        return ErrModelCall.WithError(err)
    }
    
    // Send chat complete event
    event.Emit(ctx, event.NewEvent(event.EventChatComplete, event.ChatData{
        Query:      chatManage.Query,
        ModelID:    chatManage.ChatModelID,
        Response:   chatResponse.Content,
        TokenCount: chatResponse.TokenCount,
        Duration:   time.Since(startTime).Milliseconds(),
        IsStream:   false,
    }).WithSessionID(chatManage.SessionID))
    
    chatManage.ChatResponse = chatResponse
    return next()
}
```

### 3. Sending Request Received Event in Handler Layer

```go
// internal/handler/message.go

func (h *MessageHandler) SendMessage(c *gin.Context) {
    ctx := c.Request.Context()
    
    // Parse request
    var req types.SendMessageRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // Send query received event
    event.Emit(ctx, event.NewEvent(event.EventQueryReceived, event.QueryData{
        OriginalQuery: req.Content,
        SessionID:     req.SessionID,
        UserID:        c.GetString("user_id"),
    }).WithSessionID(req.SessionID).WithRequestID(c.GetString("request_id")))
    
    // Process message...
}
```

### 4. Custom Monitoring Handler

```go
// internal/monitoring/event_monitor.go

package monitoring

import (
    "context"
    "github.com/Tencent/WeKnora/internal/event"
    "github.com/prometheus/client_golang/prometheus"
)

var (
    retrievalDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "retrieval_duration_milliseconds",
            Help: "Duration of retrieval operations",
        },
        []string{"knowledge_base_id", "retrieval_type"},
    )
    
    rerankDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "rerank_duration_milliseconds",
            Help: "Duration of rerank operations",
        },
        []string{"model_id"},
    )
)

func init() {
    prometheus.MustRegister(retrievalDuration)
    prometheus.MustRegister(rerankDuration)
}

func SetupEventMonitoring() {
    bus := event.GetGlobalEventBus()
    
    // Monitor retrieval performance
    bus.On(event.EventRetrievalComplete, func(ctx context.Context, e event.Event) error {
        data := e.Data.(event.RetrievalData)
        retrievalDuration.WithLabelValues(
            data.KnowledgeBaseID,
            data.RetrievalType,
        ).Observe(float64(data.Duration))
        return nil
    })
    
    // Monitor rerank performance
    bus.On(event.EventRerankComplete, func(ctx context.Context, e event.Event) error {
        data := e.Data.(event.RerankData)
        rerankDuration.WithLabelValues(data.ModelID).Observe(float64(data.Duration))
        return nil
    })
}
```

### 5. Logging Handler

```go
// internal/logging/event_logger.go

package logging

import (
    "context"
    "encoding/json"
    "github.com/Tencent/WeKnora/internal/event"
    "github.com/Tencent/WeKnora/internal/logger"
)

func SetupEventLogging() {
    bus := event.GetGlobalEventBus()
    
    // Structured logging for all events
    logHandler := event.ApplyMiddleware(
        func(ctx context.Context, e event.Event) error {
            data, _ := json.Marshal(e.Data)
            logger.Infof(ctx, "Event: type=%s, session=%s, request=%s, data=%s",
                e.Type, e.SessionID, e.RequestID, string(data))
            return nil
        },
        event.WithTiming(),
    )
    
    // Register to all key events
    bus.On(event.EventQueryReceived, logHandler)
    bus.On(event.EventQueryRewritten, logHandler)
    bus.On(event.EventRetrievalComplete, logHandler)
    bus.On(event.EventRerankComplete, logHandler)
    bus.On(event.EventChatComplete, logHandler)
    bus.On(event.EventError, logHandler)
}
```

### 6. Complete Initialization Flow

```go
// cmd/server/main.go or internal/container/container.go

func Initialize() {
    // 1. Initialize event system
    eventBus := event.GetGlobalEventBus()
    
    // 2. Setup monitoring
    event.NewMonitoringHandler(eventBus)
    
    // 3. Setup analytics
    event.NewAnalyticsHandler(eventBus)
    
    // 4. Setup Prometheus monitoring (if needed)
    // monitoring.SetupEventMonitoring()
    
    // 5. Setup structured logging (if needed)
    // logging.SetupEventLogging()
    
    // 6. Other initialization...
}
```

## Testing Event System

```go
// Use independent event bus in tests
func TestMyService(t *testing.T) {
    ctx := context.Background()
    
    // Create test-specific event bus
    testBus := event.NewEventBus()
    
    // Register test listener
    var receivedEvents []event.Event
    testBus.On(event.EventQueryReceived, func(ctx context.Context, e event.Event) error {
        receivedEvents = append(receivedEvents, e)
        return nil
    })
    
    // Execute test...
    testBus.Emit(ctx, event.NewEvent(event.EventQueryReceived, event.QueryData{
        OriginalQuery: "test",
    }))
    
    // Verify events
    if len(receivedEvents) != 1 {
        t.Errorf("Expected 1 event, got %d", len(receivedEvents))
    }
}
```

## Async Processing Example

```go
// For events that don't affect the main flow, use async mode
func SetupAsyncAnalytics() {
    asyncBus := event.NewAsyncEventBus()
    
    asyncBus.On(event.EventQueryReceived, func(ctx context.Context, e event.Event) error {
        // Async send to analytics platform, don't block main flow
        // sendToAnalyticsPlatform(e)
        return nil
    })
    
    // Use async bus to send events
    // asyncBus.Emit(ctx, event)
}
```

## Performance Optimization Recommendations

1. **Avoid using synchronous event bus on critical path**: For monitoring, logging, etc. that don't affect business logic, use async mode
2. **Use middleware wisely**: Only use middleware where needed, avoid unnecessary overhead
3. **Control event data size**: Avoid passing large amounts of data in events, especially in async mode
4. **Use dedicated listeners**: Don't do too much in one listener, maintain single responsibility
