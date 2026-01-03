#!/bin/bash
# WeKnora local çalıştırma scripti

cd "$(dirname "$0")"

# .env dosyasını yükle
if [ -f ".env" ]; then
    set -a
    source .env
    set +a
    echo "✓ .env dosyası yüklendi"
else
    echo "✗ .env dosyası bulunamadı!"
    exit 1
fi

# Localhost adreslerini ayarla (Docker container adreslerini override et)
export DB_HOST=localhost
export REDIS_ADDR=localhost:6379
export DOCREADER_ADDR=localhost:50051

echo "✓ Environment variable'lar ayarlandı"
echo "  DB_HOST: $DB_HOST"
echo "  DB_PORT: $DB_PORT"
echo "  DB_NAME: $DB_NAME"
echo ""

# Go uygulamasını çalıştır
echo "🚀 WeKnora başlatılıyor..."
go run cmd/server/main.go
