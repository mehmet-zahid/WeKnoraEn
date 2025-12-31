# WeKnora Development Guide

## Quick Development Mode (Recommended)

If you need to frequently modify `app` or `frontend` code, **you don't need to rebuild Docker images every time**. You can use local development mode.

### Method 1: Using Make Commands (Recommended)

#### 1. Start Infrastructure Services

```bash
make dev-start
```

This will start the following Docker containers:
- PostgreSQL (Database)
- Redis (Cache)
- MinIO (Object Storage)
- Neo4j (Graph Database)
- DocReader (Document Reading Service)
- Jaeger (Distributed Tracing)

#### 2. Start Backend Application (New Terminal)

```bash
make dev-app
```

This will run the Go application locally. After modifying code, press Ctrl+C to stop, then restart.

#### 3. Start Frontend (New Terminal)

```bash
make dev-frontend
```

This will start the Vite development server with hot reload support. Code changes will automatically refresh.

#### 4. Check Service Status

```bash
make dev-status
```

#### 5. Stop All Services

```bash
make dev-stop
```

### Method 2: Using Script Commands

If you prefer to use scripts directly:

```bash
# Start infrastructure
./scripts/dev.sh start

# Start backend (new terminal)
./scripts/dev.sh app

# Start frontend (new terminal)
./scripts/dev.sh frontend

# View logs
./scripts/dev.sh logs

# Stop all services
./scripts/dev.sh stop
```

## Access Addresses

### Development Environment

- **Frontend Development Server**: http://localhost:5173
- **Backend API**: http://localhost:8080
- **PostgreSQL**: localhost:5432
- **Redis**: localhost:6379
- **MinIO Console**: http://localhost:9001
- **Neo4j Browser**: http://localhost:7474
- **Jaeger UI**: http://localhost:16686

## Development Workflow Comparison

### ❌ Old Method (Slow)

```bash
# After each code modification:
sh scripts/build_images.sh -p      # Rebuild images (very slow)
sh scripts/start_all.sh --no-pull  # Restart containers
```

**Time Required**: 2-5 minutes per modification

### ✅ New Method (Fast)

```bash
# First startup (only once):
make dev-start

# Run in two separate terminals:
make dev-app       # After modifying Go code, Ctrl+C and restart (seconds)
make dev-frontend  # Frontend code changes auto-reload (no restart needed)
```

**Time Required**:
- First startup: 1-2 minutes
- Subsequent backend modifications: 5-10 seconds (restart Go app)
- Subsequent frontend modifications: Real-time hot reload

## Using Air for Backend Hot Reload (Optional)

If you want the backend to automatically restart after code changes, you can install `air`:

### 1. Install Air

```bash
go install github.com/cosmtrek/air@latest
```

### 2. Create Configuration File

Create `.air.toml` in the project root:

```toml
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  args_bin = []
  bin = "./tmp/main"
  cmd = "go build -o ./tmp/main ./cmd/server"
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "testdata", "frontend", "migrations"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html", "yaml"]
  include_file = []
  kill_delay = "0s"
  log = "build-errors.log"
  poll = false
  poll_interval = 0
  rerun = false
  rerun_delay = 500
  send_interrupt = false
  stop_on_error = false

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  main_only = false
  time = false

[misc]
  clean_on_exit = false

[screen]
  clear_on_rebuild = false
  keep_scroll = true
```

### 3. Start with Air

```bash
# In project root directory
air
```

Now Go code changes will automatically recompile and restart!

## Other Development Tips

### Only Modify Frontend

If you only modify the frontend, you just need:

```bash
cd frontend
npm run dev
```

The frontend will connect to the backend API at http://localhost:8080.

### Only Modify Backend

If you only modify the backend, you just need:

```bash
# Start infrastructure
make dev-start

# Run backend
make dev-app
```

### Debug Mode

#### Backend Debugging

Use VS Code or GoLand's debugging features, configured to connect to the locally running Go application.

VS Code configuration example (`.vscode/launch.json`):

```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch Server",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}/cmd/server",
            "env": {
                "DB_HOST": "localhost",
                "DOCREADER_ADDR": "localhost:50051",
                "MINIO_ENDPOINT": "localhost:9000",
                "REDIS_ADDR": "localhost:6379",
                "OTEL_EXPORTER_OTLP_ENDPOINT": "localhost:4317",
                "NEO4J_URI": "bolt://localhost:7687"
            },
            "args": []
        }
    ]
}
```

#### Frontend Debugging

Use browser developer tools. Vite provides source maps.

## Production Deployment

When you complete development and need to deploy, then build images:

```bash
# Build all images
sh scripts/build_images.sh

# Or build specific images only
sh scripts/build_images.sh -p  # Only build backend
sh scripts/build_images.sh -f  # Only build frontend

# Start production environment
sh scripts/start_all.sh
```

## Common Issues

### Q: Error connecting to database when starting dev-app

A: Make sure you've run `make dev-start` first and wait for all services to start (approximately 30 seconds).

### Q: CORS error when frontend accesses API

A: Check the frontend proxy configuration and ensure `vite.config.ts` has the correct proxy settings.

### Q: DocReader service needs to be rebuilt?

A: DocReader still uses Docker images. If you need to modify it, rebuild:

```bash
sh scripts/build_images.sh -d
make dev-restart
```

## Summary

- **Daily Development**: Use `make dev-*` commands for rapid iteration
- **Integration Testing**: Use `sh scripts/start_all.sh --no-pull` to test the complete environment
- **Production Deployment**: Use `sh scripts/build_images.sh` + `sh scripts/start_all.sh`

