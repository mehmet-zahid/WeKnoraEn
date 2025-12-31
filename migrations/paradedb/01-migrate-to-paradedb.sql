-- Migration script: Migrate from PostgreSQL to ParadeDB
-- Note: Please ensure you have backed up your data before executing this script

-- 1. Export data (execute in PostgreSQL)
-- pg_dump -U postgres -h localhost -p 5432 -d your_database > backup.sql

-- 2. Import data (execute in ParadeDB)
-- psql -U postgres -h localhost -p 5432 -d your_database < backup.sql

-- 3. Verify data


-- Insert some sample data
-- INSERT INTO tenants (id, name, description, status, api_key)
-- VALUES 
--     (1, 'Demo Tenant', 'This is a demo tenant for testing', 'active', 'sk-00000001abcdefg123456')
-- ON CONFLICT DO NOTHING;

-- SELECT setval('tenants_id_seq', (SELECT MAX(id) FROM tenants));


-- -- Create knowledge base
-- INSERT INTO knowledge_bases (id, name, description, tenant_id, chunking_config, image_processing_config, embedding_model_id)
-- VALUES 
--     ('kb-00000001', 'Default Knowledge Base', 'Default knowledge base for testing', 1, '{"chunk_size": 512, "chunk_overlap": 50, "separators": ["\n\n", "\n", "。"], "keep_separator": true}', '{"enable_multimodal": false, "model_id": ""}', 'model-embedding-00000001'),
--     ('kb-00000002', 'Test Knowledge Base', 'Test knowledge base for development', 1, '{"chunk_size": 512, "chunk_overlap": 50, "separators": ["\n\n", "\n", "。"], "keep_separator": true}', '{"enable_multimodal": false, "model_id": ""}', 'model-embedding-00000001'),
--     ('kb-00000003', 'Test Knowledge Base 2', 'Test knowledge base for development 2', 1, '{"chunk_size": 512, "chunk_overlap": 50, "separators": ["\n\n", "\n", "。"], "keep_separator": true}', '{"enable_multimodal": false, "model_id": ""}', 'model-embedding-00000001')
-- ON CONFLICT DO NOTHING;


SELECT COUNT(*) FROM tenants;
SELECT COUNT(*) FROM models;
SELECT COUNT(*) FROM knowledge_bases;
SELECT COUNT(*) FROM knowledges;


-- Test Chinese full-text search

-- Create documents table
CREATE TABLE chinese_documents (
    id SERIAL PRIMARY KEY,
    title TEXT,
    content TEXT,
    published_date DATE
);

-- Create BM25 index on the table, using Chinese tokenizer to support Chinese text
CREATE INDEX idx_documents_bm25 ON chinese_documents
USING bm25 (id, content)
WITH (
    key_field = 'id',
    text_fields = '{
        "content": {
          "tokenizer": {"type": "chinese_lindera"}
        }
    }'
);

INSERT INTO chinese_documents (title, content, published_date)
VALUES 
('Development of Artificial Intelligence', 'Artificial intelligence technology is developing rapidly, affecting all aspects of our lives. Large language models are an important recent breakthrough.', '2023-01-15'),
('Machine Learning Fundamentals', 'Machine learning is an important branch of artificial intelligence, including methods such as supervised learning, unsupervised learning, and reinforcement learning.', '2023-02-20'),
('Deep Learning Applications', 'Deep learning has wide applications in image recognition, natural language processing, and speech recognition.', '2023-03-10'),
('Natural Language Processing Technology', 'Natural language processing allows computers to understand, interpret, and generate human language, and is one of the core technologies of artificial intelligence.', '2023-04-05'),
('Introduction to Computer Vision', 'Computer vision enables machines to "see" and understand the visual world, with wide applications in security, healthcare, and other fields.', '2023-05-12');

INSERT INTO chinese_documents (title, content, published_date)
VALUES 
('hello world', 'hello world', '2023-05-12');
