# WeKnora - Docker Olmadan Çalıştırma Rehberi

## 📍 Main.go Dosyası Konumu
```
cmd/server/main.go
```

## 🚀 Hızlı Başlangıç

### 1. Gereksinimler
- Go 1.24+ yüklü olmalı
- PostgreSQL çalışıyor olmalı (localhost:5432)
- Redis çalışıyor olmalı (localhost:6379)
- DocReader servisi çalışıyor olmalı (localhost:50051) - opsiyonel

### 2. Bağımlılıkları Yükle
```bash
go mod download
```

### 3. Environment Variables Ayarla
`.env` dosyası proje kök dizininde olmalı. Gerekli değişkenler:
- `DB_DRIVER`, `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`
- `REDIS_ADDR`
- `STORAGE_TYPE`
- `DOCREADER_ADDR` (opsiyonel)

### 4. Çalıştırma Yöntemleri

#### Yöntem 1: go run ile direkt çalıştırma
```bash
go run cmd/server/main.go
```

#### Yöntem 2: Makefile kullanarak
```bash
make run
```

#### Yöntem 3: Önce build edip sonra çalıştırma
```bash
make build
./WeKnora
```

#### Yöntem 4: Development script kullanarak (önerilen)
```bash
# Önce bağımlılıkları başlat (PostgreSQL, Redis vb.)
make dev-start

# Sonra backend'i çalıştır
make dev-app
```

## 🔧 Gelişmiş Kullanım

### Environment Variables ile Çalıştırma
```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=weknora
export REDIS_ADDR=localhost:6379
export STORAGE_TYPE=local

go run cmd/server/main.go
```

### Air ile Hot Reload (Geliştirme için)
```bash
# Air kurulumu
go install github.com/air-verse/air@latest

# Air ile çalıştırma (kod değişikliklerinde otomatik restart)
air
```

## 📝 Notlar

1. **Config Dosyası**: `config/config.yaml` dosyası otomatik olarak yüklenir
2. **Port**: Varsayılan olarak `8080` portunda çalışır
3. **Database Migration**: İlk çalıştırmada veritabanı migration'ları otomatik çalışabilir
4. **DocReader**: Eğer docreader servisi yoksa, bazı özellikler çalışmayabilir

## 🐛 Sorun Giderme

### "config file not found" hatası
- `config/config.yaml` dosyasının mevcut olduğundan emin olun

### Database bağlantı hatası
- PostgreSQL'in çalıştığından emin olun
- `.env` dosyasındaki database bilgilerini kontrol edin

### Redis bağlantı hatası
- Redis'in çalıştığından emin olun
- `REDIS_ADDR` environment variable'ını kontrol edin

