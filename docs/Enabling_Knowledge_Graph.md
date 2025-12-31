# Enabling Knowledge Graph Feature Guide

This document explains how to enable and verify the Knowledge Graph (Neo4j) feature in WeKnora, helping you complete the entire process from environment setup to frontend configuration.

## Prerequisites

- WeKnora backend and frontend have been deployed.
- Docker/Docker Compose runtime environment is available.
- Accessible Neo4j service locally or remotely (recommended to use the project's built-in Docker Compose).

## Step 1: Configure Environment Variables

Add or modify the following variables in the `.env` file in the project root:

```
NEO4J_ENABLE=true
NEO4J_URI=bolt://neo4j:7687
NEO4J_USERNAME=neo4j
NEO4J_PASSWORD=your_strong_password
# Optional: NEO4J_DATABASE=neo4j
```

Notes:

- `NEO4J_ENABLE` must be set to `true` to enable knowledge graph related logic.
- `neo4j` in `NEO4J_URI` is the docker-compose service name. If using an external instance, replace with the actual address.
- If using secret management in production, ensure passwords are injected securely.

## Step 2: Start Neo4j Service

The project includes a Neo4j component that can be started directly with:

```bash
docker-compose --profile neo4j up -d
```

Common verification commands:

```bash
docker ps | grep neo4j
```

If you need custom mounts or memory settings, edit the `neo4j` service configuration in `docker-compose.yml`.

## Step 3: Restart WeKnora Services

To make the new environment variables take effect, restart the backend and frontend (examples for reference):

```bash
make stop && make start
# or
docker compose up -d --build
```

Ensure the backend logs show a successful Neo4j initialization message.

## Step 4: Enable Entity/Relationship Extraction in Frontend

1. Log in to the WeKnora frontend management page.
2. Open "Knowledge Base Settings" or create a new knowledge base.
3. Check the "Enable Entity Extraction" and "Enable Relationship Extraction" switches.
4. Fill in the required LLM, callback, or model parameters as prompted by the interface (if any).

After saving, the system will automatically trigger entity and relationship extraction tasks during document ingestion.

## Step 5: Verify Knowledge Graph

### Method 1: Neo4j Console

1. Access `http://localhost:7474` (or corresponding host/port).
2. Log in with the account and password from `.env`.
3. Execute `MATCH (n) RETURN n LIMIT 50;` to check if there are new nodes/relationships.

### Method 2: WeKnora Interface

After uploading documents in the knowledge base or conversation page, the frontend should display a graph visualization entry. During conversations, the system will automatically query the graph based on intent and return supplementary information.

## Troubleshooting

- **Cannot connect to Neo4j**: Verify network reachability, `NEO4J_URI`, username and password are correct, and check Neo4j container logs.
- **No nodes generated**: Confirm that entity/relationship extraction is enabled in the knowledge base, and uploaded documents have been parsed. Check backend logs for extraction task errors.
- **No query results**: Try executing `CALL db.schema.visualization;` in the Neo4j console to check if the schema exists. Re-import documents if necessary.

After completing the above steps, the Knowledge Graph feature is successfully enabled and can be combined with RAG and Agent workflows to improve Q&A quality.

