# WeKnora API Documentation

## Table of Contents

- [Overview](#overview)
- [Basic Information](#basic-information)
- [Authentication Mechanism](#authentication-mechanism)
- [Error Handling](#error-handling)
- [API Overview](#api-overview)

## Overview

WeKnora provides a series of RESTful APIs for creating and managing knowledge bases, retrieving knowledge, and performing knowledge-based Q&A. This document describes in detail how to use these APIs.

## Basic Information

- **Base URL**: `/api/v1`
- **Response Format**: JSON
- **Authentication Method**: API Key

## Authentication Mechanism

All API requests need to include `X-API-Key` in the HTTP request header for authentication:

```
X-API-Key: your_api_key
```

For convenient issue tracking and debugging, it is recommended to add `X-Request-ID` to each request's HTTP request header:

```
X-Request-ID: unique_request_id
```

### Get API Key

After completing account registration on the web page, please go to the account information page to get your API Key.

Please keep your API Key secure and avoid leakage. The API Key represents your account identity and has full API access permissions.

## Error Handling

All APIs use standard HTTP status codes to indicate request status and return a unified error response format:

```json
{
  "success": false,
  "error": {
    "code": "error_code",
    "message": "error_message",
    "details": "error_details"
  }
}
```

## API Overview

WeKnora APIs are divided into the following categories by function:

| Category | Description | Documentation Link |
|----------|-------------|---------------------|
| Tenant Management | Create and manage tenant accounts | [tenant.md](./tenant.md) |
| Knowledge Base Management | Create, query, and manage knowledge bases | [knowledge-base.md](./knowledge-base.md) |
| Knowledge Management | Upload, retrieve, and manage knowledge content | [knowledge.md](./knowledge.md) |
| Model Management | Configure and manage various AI models | [model.md](./model.md) |
| Chunk Management | Manage chunked content of knowledge | [chunk.md](./chunk.md) |
| Tag Management | Manage tag classifications of knowledge bases | [tag.md](./tag.md) |
| FAQ Management | Manage FAQ Q&A pairs | [faq.md](./faq.md) |
| Session Management | Create and manage conversation sessions | [session.md](./session.md) |
| Knowledge Search | Search content in knowledge bases | [knowledge-search.md](./knowledge-search.md) |
| Chat Functionality | Q&A based on knowledge bases and Agents | [chat.md](./chat.md) |
| Message Management | Get and manage conversation messages | [message.md](./message.md) |
| Evaluation Functionality | Evaluate model performance | [evaluation.md](./evaluation.md) |
