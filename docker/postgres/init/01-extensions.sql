-- PostgreSQL Extension'larını otomatik yükle
-- Bu script container ilk başlatıldığında otomatik çalışır

-- pgvector extension (vector search için)
CREATE EXTENSION IF NOT EXISTS vector;

-- pg_trgm extension (text similarity için)
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- pg_search extension (ParadeDB'ye özel, standart PostgreSQL'de yok)
-- Eğer ParadeDB kullanmıyorsanız bu satırı yorum satırı yapın
-- BEGIN
--     CREATE EXTENSION IF NOT EXISTS pg_search;
-- EXCEPTION WHEN OTHERS THEN
--     RAISE NOTICE 'pg_search extension not available (ParadeDB-specific), skipping...';
-- END;

-- uuid-ossp extension (UUID generation için)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

