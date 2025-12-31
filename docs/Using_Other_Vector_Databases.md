### How to Integrate New Vector Databases

This document provides a complete guide for adding new vector database support to the WeKnora project. By implementing standardized interfaces and following a structured process, developers can efficiently integrate custom vector databases.

### Integration Process

#### 1. Implement Basic Retrieval Engine Interface

First, you need to implement the `RetrieveEngine` interface in the `interfaces` package, which defines the core capabilities of the retrieval engine:

```go
type RetrieveEngine interface {
    // Returns the type identifier of the retrieval engine
    EngineType() types.RetrieverEngineType

    // Executes retrieval operation, returns matching results
    Retrieve(ctx context.Context, params types.RetrieveParams) ([]*types.RetrieveResult, error)

    // Returns the list of retrieval types supported by this engine
    Support() []types.RetrieverType
}
```

#### 2. Implement Storage Layer Interface

Implement the `RetrieveEngineRepository` interface to extend the basic retrieval engine capabilities and add index management functionality:

```go
type RetrieveEngineRepository interface {
    // Save a single index information
    Save(ctx context.Context, indexInfo *types.IndexInfo, params map[string]any) error
    
    // Batch save multiple index information
    BatchSave(ctx context.Context, indexInfoList []*types.IndexInfo, params map[string]any) error
    
    // Estimate the storage space required for indexes
    EstimateStorageSize(ctx context.Context, indexInfoList []*types.IndexInfo, params map[string]any) int64
    
    // Delete indexes by chunk ID list
    DeleteByChunkIDList(ctx context.Context, indexIDList []string, dimension int) error
    
    // Copy index data to avoid recalculating embedding vectors
    CopyIndices(
        ctx context.Context,
        sourceKnowledgeBaseID string,
        sourceToTargetKBIDMap map[string]string,
        sourceToTargetChunkIDMap map[string]string,
        targetKnowledgeBaseID string,
        dimension int,
    ) error
    
    // Delete indexes by knowledge ID list
    DeleteByKnowledgeIDList(ctx context.Context, knowledgeIDList []string, dimension int) error
    
    // Inherits RetrieveEngine interface
    RetrieveEngine
}
```

#### 3. Implement Service Layer Interface

Create a service that implements the `RetrieveEngineService` interface, responsible for handling the business logic of index creation and management:

```go
type RetrieveEngineService interface {
    // Create a single index
    Index(ctx context.Context,
        embedder embedding.Embedder,
        indexInfo *types.IndexInfo,
        retrieverTypes []types.RetrieverType,
    ) error

    // Batch create indexes
    BatchIndex(ctx context.Context,
        embedder embedding.Embedder,
        indexInfoList []*types.IndexInfo,
        retrieverTypes []types.RetrieverType,
    ) error

    // Estimate index storage space
    EstimateStorageSize(ctx context.Context,
        embedder embedding.Embedder,
        indexInfoList []*types.IndexInfo,
        retrieverTypes []types.RetrieverType,
    ) int64
    
    // Copy index data
    CopyIndices(
        ctx context.Context,
        sourceKnowledgeBaseID string,
        sourceToTargetKBIDMap map[string]string,
        sourceToTargetChunkIDMap map[string]string,
        targetKnowledgeBaseID string,
        dimension int,
    ) error

    // Delete indexes
    DeleteByChunkIDList(ctx context.Context, indexIDList []string, dimension int) error
    DeleteByKnowledgeIDList(ctx context.Context, knowledgeIDList []string, dimension int) error

    // Inherits RetrieveEngine interface
    RetrieveEngine
}
```

#### 4. Add Environment Variable Configuration

Add the necessary connection parameters for the new database in the environment configuration:

```
# Add the new database driver name to RETRIEVE_DRIVER (multiple drivers separated by commas)
RETRIEVE_DRIVER=postgres,elasticsearch_v8,your_database

# Connection parameters for the new database
YOUR_DATABASE_ADDR=your_database_host:port
YOUR_DATABASE_USERNAME=username
YOUR_DATABASE_PASSWORD=password
# Other necessary connection parameters...
```

#### 5. Register Retrieval Engine

Add the initialization and registration logic for the new database in the `initRetrieveEngineRegistry` function in the `internal/container/container.go` file:

```go
func initRetrieveEngineRegistry(db *gorm.DB, cfg *config.Config) (interfaces.RetrieveEngineRegistry, error) {
    registry := retriever.NewRetrieveEngineRegistry()
    retrieveDriver := strings.Split(os.Getenv("RETRIEVE_DRIVER"), ",")
    log := logger.GetLogger(context.Background())

    // Existing PostgreSQL and Elasticsearch initialization code...
    
    // Add initialization code for new vector database
    if slices.Contains(retrieveDriver, "your_database") {
        // Initialize database client
        client, err := your_database.NewClient(your_database.Config{
            Addresses: []string{os.Getenv("YOUR_DATABASE_ADDR")},
            Username:  os.Getenv("YOUR_DATABASE_USERNAME"),
            Password:  os.Getenv("YOUR_DATABASE_PASSWORD"),
            // Other connection parameters...
        })
        
        if err != nil {
            log.Errorf("Create your_database client failed: %v", err)
        } else {
            // Create retrieval engine repository
            yourDatabaseRepo := your_database.NewYourDatabaseRepository(client, cfg)
            
            // Register retrieval engine
            if err := registry.Register(
                retriever.NewKVHybridRetrieveEngine(
                    yourDatabaseRepo, types.YourDatabaseRetrieverEngineType,
                ),
            ); err != nil {
                log.Errorf("Register your_database retrieve engine failed: %v", err)
            } else {
                log.Infof("Register your_database retrieve engine success")
            }
        }
    }

    return registry, nil
}
```

#### 6. Define Retrieval Engine Type Constants

Add new retrieval engine type constants in the `internal/types/retriever.go` file:

```go
// RetrieverEngineType defines retrieval engine types
const (
    ElasticsearchRetrieverEngineType RetrieverEngineType = "elasticsearch"
    PostgresRetrieverEngineType      RetrieverEngineType = "postgres"
    YourDatabaseRetrieverEngineType  RetrieverEngineType = "your_database" // Add new database type
)
```

## Reference Implementation Examples

It is recommended to refer to the existing PostgreSQL and Elasticsearch implementations as development templates. These implementations are located in the following directories:

- PostgreSQL: `internal/application/repository/retriever/postgres/`
- ElasticsearchV7: `internal/application/repository/retriever/elasticsearch/v7/`
- ElasticsearchV8: `internal/application/repository/retriever/elasticsearch/v8/`

By following the above steps and referring to existing implementations, you can successfully integrate new vector databases into the WeKnora system, expanding its vector retrieval capabilities.

