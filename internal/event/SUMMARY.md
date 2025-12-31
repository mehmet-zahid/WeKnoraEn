# WeKnora Event System Summary

## Overview

Successfully created a complete event emission and listening mechanism for the WeKnora project, supporting event processing for various steps in the user query processing flow.

## Core Features

### ✅ Implemented Features

1. **Event Bus (EventBus)**
   - `Emit(ctx, event)` - Emit event
   - `On(eventType, handler)` - Register event listener
   - `Off(eventType)` - Remove event listener
   - `EmitAndWait(ctx, event)` - Emit event and wait for all handlers to complete
   - Synchronous/Asynchronous modes

2. **Event Types**
   - Query processing events (receive, validate, preprocess, rewrite)
   - Retrieval events (start, vector retrieval, keyword retrieval, entity retrieval, complete)
   - Rerank events (start, complete)
   - Merge events (start, complete)
   - Chat generation events (start, complete, streaming output)
   - Error events

3. **Event Data Structures**
   - `QueryData` - Query data
   - `RetrievalData` - Retrieval data
   - `RerankData` - Rerank data
   - `MergeData` - Merge data
   - `ChatData` - Chat data
   - `ErrorData` - Error data

4. **Middleware Support**
   - `WithLogging()` - Logging middleware
   - `WithTiming()` - Timing middleware
   - `WithRecovery()` - Error recovery middleware
   - `Chain()` - Middleware composition

5. **Global Event Bus**
   - Singleton global event bus
   - Global convenience functions (`On`, `Emit`, `EmitAndWait`, etc.)

6. **Examples and Tests**
   - Complete unit tests
   - Performance benchmarks
   - Complete usage examples
   - Real-world scenario demonstrations

## File Structure

```
internal/event/
├── event.go                    # Core event bus implementation
├── event_data.go              # Event data structure definitions
├── middleware.go              # Middleware implementation
├── global.go                  # Global event bus
├── integration_example.go     # Integration examples (monitoring, analytics handlers)
├── example_test.go            # Tests and examples
├── demo/
│   └── main.go               # Complete RAG flow demonstration
├── README.md                 # Detailed documentation
├── usage_example.md          # Usage example documentation
└── SUMMARY.md                # This document
```

## Performance Metrics

- **Event Emission Performance**: ~9 nanoseconds/event (benchmark)
- **Concurrency Safety**: Uses `sync.RWMutex` to ensure thread safety
- **Memory Overhead**: Very low, only stores event handler function references

## Use Cases

### 1. Monitoring and Metrics Collection

```go
bus.On(event.EventRetrievalComplete, func(ctx context.Context, e event.Event) error {
    data := e.Data.(event.RetrievalData)
    // Send to Prometheus or other monitoring systems
    metricsCollector.RecordRetrievalDuration(data.Duration)
    return nil
})
```

### 2. Logging

```go
bus.On(event.EventQueryRewritten, func(ctx context.Context, e event.Event) error {
    data := e.Data.(event.QueryData)
    logger.Infof(ctx, "Query rewritten: %s -> %s", 
        data.OriginalQuery, data.RewrittenQuery)
    return nil
})
```

### 3. User Behavior Analytics

```go
bus.On(event.EventQueryReceived, func(ctx context.Context, e event.Event) error {
    data := e.Data.(event.QueryData)
    // Send to analytics platform
    analytics.TrackQuery(data.UserID, data.OriginalQuery)
    return nil
})
```

### 4. Error Tracking

```go
bus.On(event.EventError, func(ctx context.Context, e event.Event) error {
    data := e.Data.(event.ErrorData)
    // Send to error tracking system
    sentry.CaptureException(data.Error)
    return nil
})
```

## Integration Steps

### Step 1: Initialize Event System

At application startup (e.g., `main.go` or `container.go`):

```go
import "github.com/Tencent/WeKnora/internal/event"

func Initialize() {
    // Get global event bus
    bus := event.GetGlobalEventBus()
    
    // Setup monitoring and analytics
    event.NewMonitoringHandler(bus)
    event.NewAnalyticsHandler(bus)
}
```

### Step 2: Emit Events at Various Processing Stages

Add event emission in various plugins of the query processing flow:

```go
// In search.go
event.Emit(ctx, event.NewEvent(event.EventRetrievalStart, event.RetrievalData{
    Query:           chatManage.ProcessedQuery,
    KnowledgeBaseID: chatManage.KnowledgeBaseID,
    TopK:            chatManage.EmbeddingTopK,
}).WithSessionID(chatManage.SessionID))

// In rerank.go
event.Emit(ctx, event.NewEvent(event.EventRerankComplete, event.RerankData{
    Query:       chatManage.ProcessedQuery,
    InputCount:  len(chatManage.SearchResult),
    OutputCount: len(rerankResults),
    Duration:    time.Since(startTime).Milliseconds(),
}).WithSessionID(chatManage.SessionID))
```

### Step 3: Register Custom Event Handlers

Register custom handlers as needed:

```go
event.On(event.EventQueryRewritten, func(ctx context.Context, e event.Event) error {
    // Custom processing logic
    return nil
})
```

## Advantages

1. **Low Coupling**: Event emitters and listeners are completely decoupled, easy to maintain and extend
2. **High Performance**: Very low performance overhead (~9 nanoseconds/event)
3. **Flexibility**: Supports synchronous/asynchronous, single/multiple listeners
4. **Extensibility**: Easy to add new event types and handlers
5. **Type Safety**: Predefined event data structures
6. **Middleware Support**: Easy to add cross-cutting concerns (logging, timing, error handling, etc.)
7. **Test Friendly**: Easy to verify event behavior in tests

## Test Results

✅ All unit tests passed
✅ Performance tests passed (~9 nanoseconds/event)
✅ Async processing tests passed
✅ Multiple handler tests passed
✅ Complete flow demonstration successful

## Future Recommendations

### Optional Enhancements

1. **Event Persistence**: Save critical events to database or message queue
2. **Event Replay**: Support event replay for debugging or analysis
3. **Event Filtering**: Support more complex event filtering and routing
4. **Priority Queue**: Support event priority processing
5. **Distributed Events**: Support cross-service events via message queue

### Integration Recommendations

1. **Monitoring Integration**: Integrate Prometheus for metrics collection
2. **Logging Integration**: Unified structured logging
3. **Tracing Integration**: Integrate with existing tracing system
4. **Alerting Integration**: Event-based alerting mechanism

## Example Output

Running `go run ./internal/event/demo/main.go` shows complete RAG flow event output:

```
Step 1: Query Received
[MONITOR] Query received - Session: session-xxx, Query: What is RAG technology?
[ANALYTICS] Query tracked - User: user-123, Session: session-xxx

Step 2: Query Rewriting
[MONITOR] Query rewrite started
[MONITOR] Query rewritten - Original: What is RAG technology?, Rewritten: Retrieval Augmented Generation technology...
[CUSTOM] Query Transformation: ...

Step 3: Vector Retrieval
[MONITOR] Retrieval started - Type: vector, TopK: 20
[MONITOR] Retrieval completed - Results: 18, Duration: 301ms
[CUSTOM] Retrieval Efficiency: Rate: 90.00%

Step 4: Result Reranking
[MONITOR] Rerank started - Input: 18
[MONITOR] Rerank completed - Output: 5, Duration: 201ms
[CUSTOM] Rerank Statistics: Reduction: 72.22%

Step 5: Chat Completion
[MONITOR] Chat generation started
[MONITOR] Chat generation completed - Tokens: 256, Duration: 801ms
[ANALYTICS] Chat metrics - Model: gpt-4, Tokens: 256
```

## Summary

The event system has been fully implemented and tested, and can be immediately integrated into the WeKnora project for monitoring, logging, analytics, and debugging various stages of the query processing flow. The system design is simple, performs excellently, and is easy to use and extend.
