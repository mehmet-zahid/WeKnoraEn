## Introduction
WeKnora is an enterprise-grade RAG framework ready for production deployment, implementing intelligent document understanding and retrieval capabilities. The system adopts a modular design, separating document understanding, vector storage, inference, and other functionalities.

![arc](./images/arc.png)

---

## Pipeline
WeKnora processes documents through multiple steps: Insert -> Knowledge Extraction -> Indexing -> Retrieval -> Generation. The entire workflow supports multiple retrieval methods.

![](./images/pipeline2.jpeg)

Using a user-uploaded accommodation invoice PDF file as an example, here's a detailed explanation of the data flow:

### 1. Request Reception and Initialization
+ **Request Identification**: The system receives a request and assigns it a unique `request_id=Lkq0OGLYu2fV` to track the entire processing flow.
+ **Tenant and Session Validation**:
    - The system first validates tenant information (ID: 1, Name: Default Tenant).
    - Then it begins processing a Knowledge Base Q&A (Knowledge QA) request, which belongs to session `1f241340-ae75-40a5-8731-9a3a82e34fdd`.
+ **User Question**: The user's original question is: "**What is the room type for the stay?**"
+ **Message Creation**: The system creates message records for both the user's question and the upcoming answer, with IDs `703ddf09-...` and `6f057649-...` respectively.

### 2. Knowledge Base Q&A Flow Initiation
The system formally calls the knowledge base Q&A service and defines a complete processing pipeline (Pipeline) to be executed in order, containing the following 9 events:
`[rewrite_query, preprocess_query, chunk_search, chunk_rerank, chunk_merge, filter_top_k, into_chat_message, chat_completion_stream, stream_filter]`

---

### 3. Event Execution Details
#### Event 1: `rewrite_query` - Query Rewriting
+ **Purpose**: To make retrieval more precise, the system needs to combine context to understand the user's true intent.
+ **Operations**:
    1. The system retrieves the most recent 20 historical messages from the current session (actually retrieved 8) as context.
    2. Calls a local large language model named `deepseek-r1:7b`.
    3. The model analyzes the chat history to identify the questioner as "Liwx" and rewrites the original question "What is the room type for the stay?" more specifically.
+ **Result**: The question is successfully rewritten to: "**What is the room type for Liwx's stay?**"

#### Event 2: `preprocess_query` - Query Preprocessing
+ **Purpose**: Tokenize the rewritten question and convert it into a keyword sequence suitable for search engine processing.
+ **Operations**: Tokenizes the rewritten question.
+ **Result**: Generates a keyword string: "`need rewrite user question stay room type according provide information stay person Liwx choose room type twin bed room therefore rewrite after complete question is Liwx this time stay room type`"

#### Event 3: `chunk_search` - Knowledge Chunk Retrieval
This is the core **Retrieval** step. The system performs two hybrid searches.

+ **First Search (using rewritten complete sentence)**:
    - **Vector Retrieval**:
        1. Loads the embedding model `bge-m3:latest` to convert the question into a 1024-dimensional vector.
        2. Performs vector similarity search in PostgreSQL database, finding 2 relevant knowledge chunks with IDs `e3bf6599-...` and `3989c6ce-...`.
    - **Keyword Retrieval**:
        1. Simultaneously, the system also performs keyword search.
        2. Finds the same 2 knowledge chunks mentioned above.
    - **Result Merging**: The 4 results found by both methods (actually 2 duplicates) are deduplicated, resulting in 2 unique knowledge chunks.
+ **Second Search (using preprocessed keyword sequence)**:
    - The system repeats the **vector retrieval** and **keyword retrieval** processes using the tokenized keywords.
    - Eventually obtains the same 2 knowledge chunks.
+ **Final Result**: After two searches and result merging, the system identifies 2 most relevant knowledge chunks and extracts their content for answer generation.

#### Event 4: `chunk_rerank` - Result Reranking
+ **Purpose**: Use a more powerful model to perform finer sorting of initially retrieved results to improve final answer quality.
+ **Operations**: The log shows `Rerank model ID is empty, skipping reranking`. This means the system configured a reranking step but didn't specify a specific reranking model, so **this step is skipped**.

#### Event 5: `chunk_merge` - Chunk Merging
+ **Purpose**: Merge adjacent or related knowledge chunks in content to form more complete context.
+ **Operations**: The system analyzes the 2 retrieved chunks and attempts to merge them. According to the log, after processing, they remain 2 independent chunks but are sorted by relevance score.

#### Event 6: `filter_top_k` - Top-K Filtering
+ **Purpose**: Keep only the most relevant K results to prevent too much irrelevant information from interfering with the language model.
+ **Operations**: The system is configured to keep the top 5 (Top-K = 5) most relevant chunks. Since there are only 2 chunks currently, they all pass this filter.

#### Event 7 & 8: `into_chat_message` & `chat_completion_stream` - Answer Generation
This is the **Generation** step.

+ **Purpose**: Generate natural and fluent answers based on retrieved information.
+ **Operations**:
    1. The system integrates the content of the 2 retrieved knowledge chunks, the user's original question, and chat history together to form a complete prompt.
    2. Calls the `deepseek-r1:7b` large language model again and requests answer generation in **streaming** mode. Streaming output enables a typewriter effect, improving user experience.

#### Event 9: `stream_filter` - Streaming Output Filtering
+ **Purpose**: Post-process the real-time text stream generated by the model, filtering out unwanted special markers or content.
+ **Operations**:
    - The system sets up a filter to remove internal markers that the model may generate during reasoning, such as `<think>` and `</think>`.
    - The log shows the first token chunk output by the model is `<think> According`, and the filter successfully intercepts and removes the `<think>` marker, passing only "According" and subsequent content.

### 4. Completion and Response
+ **Send References**: While generating the answer, the system sends the 2 knowledge chunks used as evidence as "reference content" to the frontend for user verification.
+ **Update Message**: After the model finishes generating all content, the system updates the complete answer to the previously created message record (ID: `6f057649-...`).
+ **Request End**: The server returns a `200` success status code, marking the end of the complete flow from question to answer.

### Summary
This log completely records a typical RAG flow: The system precisely understands user intent through **query rewriting** and **preprocessing**, then uses **vector and keyword hybrid retrieval** to find relevant information from the knowledge base. Although **reranking** was skipped, **merging** and **filtering** were still performed. Finally, the retrieved knowledge is used as context for the large language model to **generate** fluent and accurate answers, with **streaming filtering** ensuring output purity.

## Document Parsing and Chunking
The code implements an independent microservice communicating via gRPC, specifically responsible for deep parsing, chunking, and multimodal information extraction of document content. It is the core executor of the "asynchronous processing" stage.

### **Overall Architecture**
This is a Python-based gRPC service whose core responsibility is to receive files (or URLs) and parse them into structured text chunks (Chunks) ready for subsequent processing (such as vectorization).

+ `server.py`: The service entry point and network layer. It is responsible for starting a multi-process, multi-threaded gRPC server, receiving requests from the Go backend, and returning parsing results.
+ `parser.py`: Implements the **Facade pattern** in design patterns. It provides a unified `Parser` class that shields the complexity of internal multiple specific parsers (such as PDF, DOCX, Markdown, etc.). External callers (`server.py`) only need to interact with this `Parser` class.
+ `base_parser.py`: The base class for parsers, defining core logic and abstract methods shared by all specific parsers. This is the "brain" of the entire parsing process, containing the most complex text chunking, image processing, OCR, and image caption generation functionalities.

---

### **Detailed Workflow**
When the Go backend starts an asynchronous task, it carries file content and configuration information and initiates a gRPC call to this Python service. The following is the complete processing flow:

#### **第一步：请求接收与分发 (**`server.py`** & **`parser.py`**)
1. **gRPC Service Entry (**`server.py: serve`**)**:
    - The service starts via the `serve()` function. It starts a **multi-process, multi-threaded** server based on environment variables (`GRPC_WORKER_PROCESSES`, `GRPC_MAX_WORKERS`) to fully utilize CPU resources and improve concurrent processing capabilities.
    - Each worker process listens on a specified port (e.g., 50051), ready to receive requests.
2. **Request Processing (**`server.py: ReadFromFile`**)**:
    - When the Go backend initiates a `ReadFromFile` request, one of the worker processes receives the request.
    - The method first parses parameters in the request, including:
        * `file_name`, `file_type`, `file_content`: Basic file information and binary content.
        * `read_config`: A complex object containing all parsing configurations, such as `chunk_size` (chunk size), `chunk_overlap` (overlap size), `enable_multimodal` (whether to enable multimodal processing), `storage_config` (object storage configuration), `vlm_config` (visual language model configuration), etc.
    - It integrates these configurations into a `ChunkingConfig` data object.
    - The most critical step is calling `self.parser.parse_file(...)`, handing the parsing task to the `Parser` facade class.
3. **Parser Selection (**`parser.py: Parser.parse_file`**)**:
    - After the `Parser` class receives the task, it first calls the `get_parser(file_type)` method.
    - This method looks up the corresponding specific parser class (e.g., `PDFParser`) in a dictionary `self.parsers` based on the file type (e.g., `'pdf'`).
    - After finding it, it **instantiates** this `PDFParser` class and passes all configuration information including `ChunkingConfig` to the constructor.

#### **Step 2: Core Parsing and Chunking (**`base_parser.py`**)**
This touches the core of the entire flow: **How to ensure contextual integrity and original order of information**.

According to the `base_parser.py` code, **the text, tables, and images in the finally split Chunks are saved in the order they appear in the original document**.

This order is guaranteed mainly thanks to several well-designed methods in `BaseParser` working together. Let's trace this flow in detail.

The guarantee of order can be divided into four stages:

1. **Stage 1: Unified Text Stream Creation (**`pdf_parser.py`**)**:
    - In the `parse_into_text` method, the code processes the PDF **page by page**.
    - Within each page, it concatenates all content into a long string (`page_content_parts`) according to a certain logic (first extract non-table text, then append tables, finally append image placeholders).
    - **Key Point**: Although at this stage, the concatenation order of text, tables, and image placeholders may not be 100% precise to the character level, it ensures that **content from the same page stays together** and roughly follows the top-to-bottom reading order.
    - Finally, content from all pages is connected by `"\n\n--- Page Break ---\n\n"`, forming a **single, ordered text stream (**`final_text`**) containing all information (text, Markdown tables, image placeholders)**.
2. **Stage 2: Atomization and Protection (**`_split_into_units`**)**:
    - This single `final_text` is passed to the `_split_into_units` method.
    - This method is **key to ensuring structural integrity**. It uses regular expressions to identify **entire Markdown tables** and **entire Markdown image placeholders** as **indivisible atomic units**.
    - It splits these atomic units (tables, images) and the ordinary text blocks between them into a list (`units`) according to their **original order** of appearance in `final_text`.
    - **Result**: We now have a list, for example `['some text', '![...](...)', 'other text', '|...|...|\n|---|---|\n...', 'more text']`. The element order in this list **exactly matches their order in the original document**.
3. **Stage 3: Sequential Chunking (**`chunk_text`**)**:
    - The `chunk_text` method receives this **ordered **`units`** list**.
    - Its working mechanism is very simple and direct: it **sequentially** traverses each unit (`unit`) in this list.
    - It **adds these units in order** to a temporary `current_chunk` list until the chunk length approaches the `chunk_size` limit.
    - After a chunk is full, it is saved, then a new chunk starts (possibly with overlap from the previous chunk).
    - **Key Point**: Because `chunk_text` **strictly processes according to the order of the **`units`** list**, it never disrupts the relative order between tables, text, and images. A table that appears first in the document will necessarily appear in a Chunk with a lower sequence number.
4. **Stage 4: Image Information Attachment (**`process_chunks_images`**)**:
    - After text chunks are split, the `process_chunks_images` method is called.
    - It processes **every** already generated Chunk.
    - Within each Chunk, it finds image placeholders and then performs AI processing.
    - Finally, it attaches the processed image information (including permanent URL, OCR text, image description, etc.) to **that Chunk's own** `.images` attribute.
    - **Key Point**: This process **does not change the order of Chunks or the content of their **`.content`** attributes**. It only attaches additional information to already existing, correctly ordered Chunks.

#### **Step 3: Multimodal Processing (if enabled) (**`base_parser.py`**)**
If `enable_multimodal` is `True`, after text chunking is complete, it enters the most complex multimodal processing stage.

1. **Concurrent Task Launch (**`BaseParser.process_chunks_images`**)**:
    - This method uses `asyncio` (Python's asynchronous I/O framework) to **concurrently process all images in text chunks**, greatly improving efficiency.
    - It creates an asynchronous task `process_chunk_images_async` for each `Chunk`.
2. **Processing Images in a Single Chunk (**`BaseParser.process_chunk_images_async`**)**:
    - **Extract Image References**: First, use regular expressions `extract_images_from_chunk` to find all image references in the current chunk's text (e.g., `![alt text](image.png)`).
    - **Image Persistence**: For each found image, concurrently call `download_and_upload_image`. This function is responsible for:
        * Retrieving image data from its original location (which may be inside a PDF, local path, or remote URL).
        * **Uploading the image to the configured object storage (COS/MinIO)**. This step is crucial, converting temporary, unstable image references into persistent, publicly accessible URLs.
        * Returning the persistent URL and image object (PIL Image).
    - **Concurrent AI Processing**: Collect all successfully uploaded images and call `process_multiple_images`.
        * This method internally uses `asyncio.Semaphore` to limit concurrency (e.g., processing at most 5 images simultaneously) to prevent excessive memory consumption or triggering model API rate limits.
        * For each image, it calls `process_image_async`.
3. **Processing a Single Image (**`BaseParser.process_image_async`**)**:
    - **OCR**: Calls `perform_ocr`, which uses an OCR engine (such as `PaddleOCR`) to recognize all text in the image.
    - **Image Caption**: Calls `get_image_caption`, which sends image data (converted to Base64) to the configured visual language model (VLM) to generate a natural language description of the image content.
    - This method returns `(ocr_text, caption, persistent_url)`.
4. **Result Aggregation**:
    - After all image processing is complete, structured information including persistent URLs, OCR text, and image descriptions is attached to the corresponding `Chunk` object's `.images` field.

#### **Step 4: Return Results (**`server.py`**)**
1. **Data Conversion (**`server.py: _convert_chunk_to_proto`**)**:
    - After `parser.parse_file` completes execution, it returns a list (`ParseResult`) containing all processed `Chunk` objects.
    - The `ReadFromFile` method receives this result and calls `_convert_chunk_to_proto` to convert Python `Chunk` objects (including their internal image information) into gRPC-defined Protobuf message format.
2. **Response Return**:
    - Finally, the gRPC server sends the `ReadResponse` message containing all chunks and multimodal information back to the caller—the Go backend service.

At this point, the Go backend has obtained structured, information-rich document data and can proceed with the next steps of vectorization and index storage.


## Deployment
Supports local deployment via Docker images and provides interface services through API ports.

## Performance and Monitoring
WeKnora includes rich monitoring and testing components:

+ **Distributed Tracing**: Integrates Jaeger to track the complete execution path of requests in the service architecture. Essentially, Jaeger is a technology that helps users "see" the complete lifecycle of requests in distributed systems.
+ **Health Monitoring**: Monitors service health status.
+ **Scalability**: Through containerized deployment, multiple services can meet large-scale concurrent requests.

## Q&A
### Question 1: What is the purpose of executing two hybrid searches in the retrieval process? And what are the differences between the first and second searches?
This is an excellent observation. The system performs two hybrid searches to **maximize retrieval accuracy and recall rate**. Essentially, it's a combination method of **Query Expansion and multi-strategy retrieval**.

#### Purpose
By searching with two different forms of queries (rewritten sentence vs. tokenized keyword sequence), the system can combine the advantages of both query methods:

+ **Semantic Retrieval Depth**: Using complete sentences for search can better leverage the vector model's (such as `bge-m3`) ability to understand the overall meaning of sentences, finding knowledge chunks that are semantically closest.
+ **Keyword Retrieval Breadth**: Using tokenized keywords for search ensures that even if the knowledge chunk's expression differs from the original question, as long as it contains core keywords, it has a chance to be matched. This is especially effective for traditional keyword matching algorithms (such as BM25).

Simply put, it's **asking the same question in two different ways**, then aggregating results from both sides to ensure the most relevant knowledge is not missed.

#### Differences Between the Two Searches
Their core difference lies in the **input query text (Query Text)**:

1. **First Hybrid Search**
    - **Input**: Uses a **grammatically complete natural language question** generated after the `rewrite_query` event.
    - **Log Evidence**:

```plain
INFO [2025-08-29 09:46:36.896] [request_id=Lkq0OGLYu2fV] knowledgebase.go:266[HybridSearch] | Hybrid search parameters, knowledge base ID: kb-00000001, query text: The user question that needs to be rewritten is: "What is the room type for the stay?". According to the provided information, the room type chosen by the guest Liwx is a twin bed room. Therefore, the complete rewritten question is: "What is the room type for Liwx's stay?"
```

2. **Second Hybrid Search**
    - **Input**: Uses a **space-separated keyword sequence** generated after the `preprocess_query` event.
    - **Log Evidence**:

```plain
INFO [2025-08-29 09:46:37.257] [request_id=Lkq0OGLYu2fV] knowledgebase.go:266[HybridSearch] | Hybrid search parameters, knowledge base ID: kb-00000001, query text: need rewrite user question stay room type according provide information stay person Liwx choose room type twin bed room therefore rewrite after complete question is Liwx this time stay room type
```

Finally, the system deduplicates and merges the results from these two searches (the log shows 2 results found each time, still 2 after deduplication), obtaining a more reliable knowledge set for subsequent answer generation.



### Question 2: Reranker Model Analysis
Rerankers are very advanced technology in the RAG field, with significant differences in working principles and applicable scenarios.

Simply put, they represent an evolution from "**specialized discriminative models**" to "**using large language models (LLM) for discrimination**" to "**deeply mining LLM internal information for discrimination**".

The following are their detailed differences:

#### 1. Normal Reranker (Conventional Reranker / Cross-Encoder)
This is the most classic and mainstream reranking method.

+ **Model Type**: **Sequence Classification Model**. Essentially a **Cross-Encoder**, usually based on bidirectional encoder architectures like BERT, RoBERTa. `BAAI/bge-reranker-base/large/v2-m3` all belong to this category.
+ **Working Principle**:
    1. It concatenates the **Query** and **document to be ranked (Passage)** into a single input sequence, for example: `[CLS] what is panda? [SEP] The giant panda is a bear species endemic to China. [SEP]`.
    2. This concatenated sequence is fully fed into the model. The model's internal self-attention mechanism can simultaneously analyze every word in both the query and document, calculating **fine-grained interaction relationships** between them.
    3. The model finally outputs a **single score (Logit)**, which directly represents the relevance between the query and document. The higher the score, the stronger the relevance.
+ **Key Characteristics**:
    - **Advantages**: Because queries and documents undergo full, deep interaction within the model, its **accuracy is usually very high**, making it the gold standard for measuring Reranker performance.
    - **Disadvantages**: **Relatively slow**. Because it must independently execute a complete, costly computation for **every "query-document" pair**. If initial retrieval returns 100 documents, it needs to run 100 times.



#### 2. LLM-based Reranker
This method creatively leverages the capabilities of general-purpose large language models (LLM) for reranking.

+ **Model Type**: **Causal Language Model**, i.e., the LLMs we commonly refer to like GPT, Llama, Gemma used for text generation. `BAAI/bge-reranker-v2-gemma` is a typical example.
+ **Working Principle**:
    1. It **does not directly output a score**, but rather **transforms the reranking task into a Q&A or text generation task**.
    2. It organizes input through a carefully designed **Prompt**, for example: `"Given a query A and a passage B, determine whether the passage contains an answer to the query by providing a prediction of either 'Yes' or 'No'. A: {query} B: {passage}"`.
    3. It feeds this complete Prompt to the LLM, then **observes the probability of the LLM generating the word 'Yes' at the end**.
    4. This **probability of generating 'Yes' (or its Logit value) is used as the relevance score**. If the model is very confident the answer is "Yes", it means it believes document B contains the answer to query A, i.e., high relevance.
+ **Key Characteristics**:
    - **Advantages**: Can leverage LLM's powerful **semantic understanding, reasoning, and world knowledge**. For complex queries requiring deep understanding and reasoning to judge relevance, the effect may be better.
    - **Disadvantages**: Computational overhead can be very large (depending on LLM size), and performance **highly depends on Prompt design**.



#### 3. LLM-based Layerwise Reranker
This is the "enhanced version" of the second method, a more cutting-edge and complex exploratory technology.

+ **Model Type**: Also a **Causal Language Model**, for example `BAAI/bge-reranker-v2-minicpm-layerwise`.
+ **Working Principle**:
    1. The input part is exactly the same as the second method, also using a "Yes/No" Prompt.
    2. The core difference lies in the **score extraction method**. It no longer relies solely on the LLM's **final layer** output (i.e., the final prediction result).
    3. It believes that during the LLM's layer-by-layer information processing, network layers (Layer) at different depths may capture semantic relevance information at different levels. Therefore, it extracts prediction Logits for the word "Yes" from **multiple intermediate layers of the model**.
    4. The `cutoff_layers=[28]` parameter in the code tells the model: "Please give me the output of layer 28". Finally, you get one or more scores from different network layers, which can be averaged or combined in other ways to form a more robust final relevance judgment.
+ **Key Characteristics**:
    - **Advantages**: Theoretically can obtain **richer, more comprehensive relevance signals**, potentially achieving higher accuracy than looking at only the final layer. It's currently a method for exploring performance limits.
    - **Disadvantages**: **Highest complexity**, requires specific modifications to the model to extract intermediate layer information (the `trust_remote_code=True` in the code is a signal), and computational overhead is also very large.

#### Summary Comparison
| Feature | 1. Normal Reranker | 2. LLM-based Reranker | 3. LLM-based Layerwise Reranker |
| :--- | :--- | :--- | :--- |
| **Underlying Model** | Cross-Encoder (e.g., BERT) | Causal Language Model (e.g., Gemma) | Causal Language Model (e.g., MiniCPM) |
| **Working Principle** | Computes deep interaction between Query and Passage, directly outputs relevance score | Transforms ranking task into "Yes/No" prediction, uses "Yes" probability as score | Similar to 2, but extracts "Yes" probability from multiple intermediate layers of LLM |
| **Output** | Single relevance score | Single relevance score (from final layer) | Multiple relevance scores (from different layers) |
| **Advantages** | **Best balance of speed and accuracy**, mature and stable | Leverages LLM reasoning capabilities, handles complex problems | Theoretically highest accuracy, richer signals |
| **Disadvantages** | Slower than vector retrieval | Large computational overhead, depends on Prompt design | **Highest complexity**, largest computational overhead |
| **Recommended Scenarios** | **Preferred for most production environments**, good effect, easy to deploy | Scenarios with extreme requirements for answer quality and sufficient computing resources | Academic research or scenarios pursuing SOTA (State-of-the-art) performance |

#### Usage Recommendations
1. **Starting Stage**: Strongly recommend starting with **`Normal Reranker`**, for example `BAAI/bge-reranker-v2-m3`. It's one of the best-performing models currently, can significantly improve your RAG system performance, and is relatively easy to integrate and deploy.
2. **Advanced Exploration**: If you find that conventional Rerankers perform poorly on very subtle queries or queries requiring complex reasoning, and you have sufficient GPU resources, you can try `LLM-based Reranker`.
3. **Cutting-edge Research**: `Layerwise Reranker` is more suitable for researchers or experts who want to squeeze out the last bit of performance on specific tasks.


### Question 3: How is knowledge (with reranking) after coarse or fine filtering assembled and sent to the large model?
This part mainly involves designing prompts, typical instruction details. The core task is to answer user questions based on context. When assembling context, the following must be specified:
- **Key Constraint**: Must strictly answer according to the provided documents, prohibit using your own knowledge to answer
- **Unknown Situation Handling**: If there is insufficient information in the documents to answer the question, please inform "Based on the available materials, I cannot answer this question"
- **Citation Requirement**: When answering, if you reference content from a document, please add the document number at the end of the sentence

---

## Manual Knowledge Online Editing

The platform's knowledge base page has added dual entry points for "Upload Document / Online Editing", supporting direct writing and maintenance of Markdown knowledge in the browser:

- **Draft Mode**: Used to temporarily store content. Drafts do not participate in retrieval.
- **Publish Operation**: Automatically triggers vectorization and index building.
- **Published Markdown Knowledge**: Can be opened again for editing and republishing.
- **Add to Knowledge Base Tool**: Provided at the end of assistant answers on the conversation page, allowing one-click import of current Q&A into the editor for confirmation and saving.













