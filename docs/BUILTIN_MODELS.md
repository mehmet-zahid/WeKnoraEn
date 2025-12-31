# Built-in Model Management Guide

## Overview

Built-in models are system-level model configurations visible to all tenants, but sensitive information is hidden and cannot be edited or deleted. Built-in models are typically used to provide system default model configurations, ensuring all tenants can use unified model services.

## Built-in Model Features

- **Visible to All Tenants**: Built-in models are visible to all tenants without separate configuration
- **Security Protection**: Sensitive information (API Key, Base URL) of built-in models is hidden and cannot be viewed in detail
- **Read-only Protection**: Built-in models cannot be edited or deleted, only set as default models
- **Unified Management**: Maintained by system administrators to ensure configuration consistency and security

## How to Add Built-in Models

Built-in models need to be inserted directly through the database. The following are the steps to add built-in models:

### 1. Prepare Model Data

First, ensure you have the model configuration information to be set as a built-in model, including:
- Model name (name)
- Model type (type): `KnowledgeQA`, `Embedding`, `Rerank`, or `VLLM`
- Model source (source): `local` or `remote`
- Model parameters (parameters): including base_url, api_key, provider, etc.
- Tenant ID (tenant_id): It is recommended to use a tenant ID less than 10000 to avoid conflicts

**Supported Providers**: `generic` (custom), `openai`, `aliyun`, `zhipu`, `volcengine`, `hunyuan`, `deepseek`, `minimax`, `mimo`, `siliconflow`, `jina`, `openrouter`, `gemini`

### 2. Execute SQL Insert Statement

Use the following SQL statement to insert built-in models:

```sql
-- Example: Insert an LLM built-in model
INSERT INTO models (
    id,
    tenant_id,
    name,
    type,
    source,
    description,
    parameters,
    is_default,
    status,
    is_builtin
) VALUES (
    'builtin-llm-001',                    -- Use fixed ID, recommended to use builtin- prefix
    10000,                                -- Tenant ID (use first tenant)
    'GPT-4',                              -- Model name
    'KnowledgeQA',                        -- Model type
    'remote',                             -- Model source
    'Built-in LLM Model',                 -- Description
    '{"base_url": "https://api.openai.com/v1", "api_key": "sk-xxx", "provider": "openai"}'::jsonb,  -- Parameters (JSON format)
    false,                                -- Is default
    'active',                             -- Status
    true                                  -- Mark as built-in model
) ON CONFLICT (id) DO NOTHING;

-- Example: Insert an Embedding built-in model
INSERT INTO models (
    id,
    tenant_id,
    name,
    type,
    source,
    description,
    parameters,
    is_default,
    status,
    is_builtin
) VALUES (
    'builtin-embedding-001',
    10000,
    'text-embedding-ada-002',
    'Embedding',
    'remote',
    'Built-in Embedding Model',
    '{"base_url": "https://api.openai.com/v1", "api_key": "sk-xxx", "provider": "openai", "embedding_parameters": {"dimension": 1536, "truncate_prompt_tokens": 0}}'::jsonb,
    false,
    'active',
    true
) ON CONFLICT (id) DO NOTHING;

-- Example: Insert a ReRank built-in model
INSERT INTO models (
    id,
    tenant_id,
    name,
    type,
    source,
    description,
    parameters,
    is_default,
    status,
    is_builtin
) VALUES (
    'builtin-rerank-001',
    10000,
    'bge-reranker-base',
    'Rerank',
    'remote',
    'Built-in ReRank Model',
    '{"base_url": "https://api.jina.ai/v1", "api_key": "jina-xxx", "provider": "jina"}'::jsonb,
    false,
    'active',
    true
) ON CONFLICT (id) DO NOTHING;

-- Example: Insert a VLLM built-in model
INSERT INTO models (
    id,
    tenant_id,
    name,
    type,
    source,
    description,
    parameters,
    is_default,
    status,
    is_builtin
) VALUES (
    'builtin-vllm-001',
    10000,
    'gpt-4-vision',
    'VLLM',
    'remote',
    'Built-in VLLM Model',
    '{"base_url": "https://dashscope.aliyuncs.com/compatible-mode/v1", "api_key": "sk-xxx", "provider": "aliyun"}'::jsonb,
    false,
    'active',
    true
) ON CONFLICT (id) DO NOTHING;
```

### 3. Verify Insertion Results

Execute the following SQL query to verify if the built-in model was successfully inserted:

```sql
SELECT id, name, type, is_builtin, status 
FROM models 
WHERE is_builtin = true
ORDER BY type, created_at;
```

## Notes

1. **ID Naming Convention**: It is recommended to use the format `builtin-{type}-{number}`, for example `builtin-llm-001`, `builtin-embedding-001`
2. **Tenant ID**: Built-in models can belong to any tenant, but it is recommended to use the first tenant ID (usually 10000)
3. **Parameter Format**: The `parameters` field must be valid JSON format
4. **Idempotency**: Use `ON CONFLICT (id) DO NOTHING` to ensure repeated execution does not cause errors
5. **Security**: The API Key and Base URL of built-in models will be automatically hidden in the frontend, but the original data in the database still exists. Please keep database access permissions secure.

## Set Existing Model as Built-in Model

If you already have a model and want to set it as a built-in model, you can use an UPDATE statement:

```sql
UPDATE models 
SET is_builtin = true 
WHERE id = 'model_id' AND name = 'model_name';
```

## Remove Built-in Model

If you need to remove the built-in model flag (restore to normal model), execute:

```sql
UPDATE models 
SET is_builtin = false 
WHERE id = 'model_id';
```

Note: After removing the built-in model flag, the model will be restored to a normal model and can be edited and deleted.
