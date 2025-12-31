# Common Questions

## 1. How to View Logs?
```bash
docker compose logs -f app docreader postgres
```

## 2. How to Start and Stop Services?
```bash
# Start services
./scripts/start_all.sh

# Stop services
./scripts/start_all.sh --stop

# Clear database
./scripts/start_all.sh --stop && make clean-db
```

## 3. Unable to Upload Documents After Service Startup?

This is usually caused by Embedding models and conversation models not being correctly configured. Follow these steps to troubleshoot:

1. Check if the model information in the `.env` configuration is complete. If using ollama to access local models, ensure the local ollama service is running normally, and the following environment variables in `.env` need to be correctly set:
```bash
# LLM Model
INIT_LLM_MODEL_NAME=your_llm_model
# Embedding Model
INIT_EMBEDDING_MODEL_NAME=your_embedding_model
# Embedding model vector dimension
INIT_EMBEDDING_MODEL_DIMENSION=your_embedding_model_dimension
# Embedding model ID, usually a string
INIT_EMBEDDING_MODEL_ID=your_embedding_model_id
```

If accessing models through remote API, you need to additionally provide the corresponding `BASE_URL` and `API_KEY`:
```bash
# LLM model access address
INIT_LLM_MODEL_BASE_URL=your_llm_model_base_url
# LLM model API key, can be set if authentication is required
INIT_LLM_MODEL_API_KEY=your_llm_model_api_key
# Embedding model access address
INIT_EMBEDDING_MODEL_BASE_URL=your_embedding_model_base_url
# Embedding model API key, can be set if authentication is required
INIT_EMBEDDING_MODEL_API_KEY=your_embedding_model_api_key
```

When reranking functionality is needed, you need to additionally configure the Rerank model. The specific configuration is as follows:
```bash
# Rerank model name to use
INIT_RERANK_MODEL_NAME=your_rerank_model_name
# Rerank model access address
INIT_RERANK_MODEL_BASE_URL=your_rerank_model_base_url
# Rerank model API key, can be set if authentication is required
INIT_RERANK_MODEL_API_KEY=your_rerank_model_api_key
```

2. Check the main service logs to see if there are any `ERROR` log outputs.

## 4. How to Enable Multimodal Functionality?
1. Ensure the following configuration in `.env` is correctly set:
```bash
# VLM_MODEL_NAME Multimodal model name to use
VLM_MODEL_NAME=your_vlm_model_name

# VLM_MODEL_BASE_URL Multimodal model access address to use
VLM_MODEL_BASE_URL=your_vlm_model_base_url

# VLM_MODEL_API_KEY Multimodal model API key to use
VLM_MODEL_API_KEY=your_vlm_model_api_key
```
Note: Multimodal large models currently only support remote API access, so you need to provide `VLM_MODEL_BASE_URL` and `VLM_MODEL_API_KEY`.

2. Parsed files need to be uploaded to COS. Ensure the `COS` information in `.env` is correctly set:
```bash
# Tencent Cloud COS access key ID
COS_SECRET_ID=your_cos_secret_id

# Tencent Cloud COS secret key
COS_SECRET_KEY=your_cos_secret_key

# Tencent Cloud COS region, e.g., ap-guangzhou
COS_REGION=your_cos_region

# Tencent Cloud COS bucket name
COS_BUCKET_NAME=your_cos_bucket_name

# Tencent Cloud COS application ID
COS_APP_ID=your_cos_app_id

# Tencent Cloud COS path prefix for storing files
COS_PATH_PREFIX=your_cos_path_prefix
```
Important: Make sure to set file permissions in COS to **public read**, otherwise the document parsing module cannot parse files normally.

3. Check the document parsing module logs to see if OCR and Caption are correctly parsed and printed.

## 5. How to Use Data Analysis Functionality?

Before using the data analysis functionality, ensure the agent has configured relevant tools:

1. **Intelligent Reasoning**: Need to check the following two tools in the tool configuration:
   - View Data Metadata
   - Data Analysis

2. **Quick Q&A Agent**: No need to manually select tools, can directly perform simple data query operations.

### Notes and Usage Guidelines

1. **Supported File Formats**
   - Currently only supports **CSV** (`.csv`) and **Excel** (`.xlsx`, `.xls`) format files.
   - For complex Excel files, if reading fails, it is recommended to convert them to standard CSV format and re-upload.

2. **Query Restrictions**
   - Only supports **read-only queries**, including `SELECT`, `SHOW`, `DESCRIBE`, `EXPLAIN`, `PRAGMA` and other statements.
   - Prohibits executing any data modification operations, such as `INSERT`, `UPDATE`, `DELETE`, `CREATE`, `DROP`, etc.

## P.S.
If the above methods do not solve the problem, please describe your issue in an issue and provide necessary log information to help us troubleshoot.
