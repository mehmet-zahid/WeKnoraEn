# WeKnora Knowledge Graph

## Quick Start

- Configure related environment variables in `.env`
    - Enable Neo4j: `NEO4J_ENABLE=true`
    - Neo4j URI: `NEO4J_URI=bolt://neo4j:7687`
    - Neo4j Username: `NEO4J_USERNAME=neo4j`
    - Neo4j Password: `NEO4J_PASSWORD=password`

- Start Neo4j
```bash
docker-compose --profile neo4j up -d
```

- Enable entity and relationship extraction in the knowledge base settings page, and configure related content according to the prompts

## Generate Graph

After uploading any document, the system will automatically extract entities and relationships and generate the corresponding knowledge graph.

![Knowledge Graph Example](./images/graph3.png)

## View Graph

Log in to `http://localhost:7474` and execute `match (n) return (n)` to view the generated knowledge graph.

During conversations, the system will automatically query the knowledge graph and retrieve relevant knowledge.
