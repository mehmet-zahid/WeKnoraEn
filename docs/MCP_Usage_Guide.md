## MCP Usage Guide

### Overview
- MCP (Model Context Protocol) allows WeKnora to securely connect to external tools or data sources, extending the capabilities that Agents can call during reasoning.
- All services are centrally managed in the frontend `Settings > MCP Services` (`frontend/src/views/settings/McpSettings.vue`), eliminating the need to manually modify configuration files.
- Each service includes name, transport method (SSE / HTTP Streamable / Stdio), connection address or command, authentication information, and advanced timeout and retry policies.

### Entry Point and Interface
- Open the `Settings -> MCP Services` option in the left menu of the console to view all MCP services under the current tenant.
- The list allows you to quickly enable/disable services, view descriptions, and perform "Test / Edit / Delete" operations through the right menu.
- The "Add Service" button opens the `McpServiceDialog` for creating or modifying services.

### Common Operations

1. **Create New Service**
   - Click "Add Service", fill in the name and description, and select the transport method.
   - SSE / HTTP Streamable requires an accessible service URL; Stdio requires configuring `uvx`/`npx` commands and parameters, with optional environment variables.
   - Fill in API Key, Bearer Token, timeout and retry policies as needed. After saving, the service will appear in the list.

2. **Enable/Disable Service**
   - Toggle the enable status in the list switch. The system will immediately call the backend `updateMCPService`. If it fails, the status will be automatically rolled back and a prompt will be displayed.

3. **Connection Test**
   - Select "Test" from the more menu. The frontend will call `/api/v1/mcp-services/{id}/test` and display the `McpTestResult`.
   - On success, it will show the list of available tools (including input schemas) and resources; on failure, it will display error information to help troubleshoot network or authentication issues.

4. **Edit / Delete**
   - "Edit" will load the existing configuration. Modify and save.
   - "Delete" requires confirmation in a popup. After completion, the list will automatically refresh.

### Usage Recommendations
- **Transport Method Selection**: Prefer SSE for streaming experience; switch to standard HTTP Streamable when compatibility is needed; use Stdio for local debugging or offline environments and start the MCP Server on the same machine.
- **Authentication Management**: Store API Keys / Tokens in "Authentication Configuration". For production environments, it is recommended to create separate minimum-privilege keys and rotate them regularly.
- **Retry Policy**: For public networks or third-party services, appropriately increase `retry_count` and `retry_delay` to avoid intermittent timeouts causing Agent interruptions.

