# WeKnora MCP Server

This is a Model Context Protocol (MCP) server that provides access to the WeKnora knowledge management API.

## Quick Start

> It is recommended to directly refer to [MCP Configuration Guide](./MCP_CONFIG.md), no need to perform the following operations.

### 1. Install Dependencies
```bash
pip install -r requirements.txt
```

### 2. Configure Environment Variables
```bash
# Linux/macOS
export WEKNORA_BASE_URL="http://localhost:8080/api/v1"
export WEKNORA_API_KEY="your_api_key_here"

# Windows PowerShell
$env:WEKNORA_BASE_URL="http://localhost:8080/api/v1"
$env:WEKNORA_API_KEY="your_api_key_here"

# Windows CMD
set WEKNORA_BASE_URL=http://localhost:8080/api/v1
set WEKNORA_API_KEY=your_api_key_here
```

### 3. Run the Server

**Recommended Method - Using Main Entry Point:**
```bash
python main.py
```

**Other Running Methods:**
```bash
# Using original startup script
python run_server.py

# Using convenience script
python run.py

# Directly run server module
python weknora_mcp_server.py

# Run as Python module
python -m weknora_mcp_server
```

### 4. Command Line Options
```bash
python main.py --help                 # Display help information
python main.py --check-only           # Only check environment configuration
python main.py --verbose              # Enable verbose logging
python main.py --version              # Display version information
```

## Install as Python Package

### Development Mode Installation
```bash
pip install -e .
```

After installation, you can use the command line tool:
```bash
weknora-mcp-server
# or
weknora-server
```

### Production Mode Installation
```bash
pip install .
```

### Build Distribution Package
```bash
# Using setuptools
python setup.py sdist bdist_wheel

# Using modern build tools
pip install build
python -m build
```

## Test Module

Run the test script to verify that the module is working correctly:
```bash
python test_module.py
```

## Features

This MCP server provides the following tools:

### Tenant Management
- `create_tenant` - Create a new tenant
- `list_tenants` - List all tenants

### Knowledge Base Management
- `create_knowledge_base` - Create a knowledge base
- `list_knowledge_bases` - List knowledge bases
- `get_knowledge_base` - Get knowledge base details
- `delete_knowledge_base` - Delete a knowledge base
- `hybrid_search` - Hybrid search

### Knowledge Management
- `create_knowledge_from_url` - Create knowledge from URL
- `list_knowledge` - List knowledge
- `get_knowledge` - Get knowledge details
- `delete_knowledge` - Delete knowledge

### Model Management
- `create_model` - Create a model
- `list_models` - List models
- `get_model` - Get model details

### Session Management
- `create_session` - Create a chat session
- `get_session` - Get session details
- `list_sessions` - List sessions
- `delete_session` - Delete a session

### Chat Functionality
- `chat` - Send chat messages

### Chunk Management
- `list_chunks` - List knowledge chunks
- `delete_chunk` - Delete a knowledge chunk

## Troubleshooting

If you encounter import errors, please ensure:
1. All required dependency packages are installed
2. Python version is compatible (recommended 3.10+)
3. No filename conflicts (avoid using `mcp.py` as a filename)

## Usage Example

<img width="950" height="2063" alt="118d078426f42f3d4983c13386085d7f" src="https://github.com/user-attachments/assets/09111ec8-0489-415c-969d-aa3835778e14" />