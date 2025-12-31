#!/bin/bash
# This script is used to start/stop Ollama and docker-compose services on demand

# Set colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No color

# Get project root directory (parent directory of script location)
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$( cd "$SCRIPT_DIR/.." && pwd )"

# Version information
VERSION="1.0.1" # Version update
SCRIPT_NAME=$(basename "$0")

# Show help information
show_help() {
    printf "%b\n" "${GREEN}WeKnora Startup Script v${VERSION}${NC}"
    printf "%b\n" "${GREEN}Usage:${NC} $0 [options]"
    echo "Options:"
    echo "  -h, --help     Show help information"
    echo "  -o, --ollama   Start Ollama service"
    echo "  -d, --docker    Start Docker container services"
    echo "  -a, --all      Start all services (default)"
    echo "  -s, --stop     Stop all services"
    echo "  -c, --check    Check environment and diagnose issues"
    echo "  -r, --restart  Rebuild and restart specified container"
    echo "  -l, --list     List all running containers"
    echo "  -p, --pull     Pull latest Docker images"
    echo "  --no-pull      Start without pulling images (default will pull)"
    echo "  -v, --version  Show version information"
    exit 0
}

# Show version information
show_version() {
    printf "%b\n" "${GREEN}WeKnora Startup Script v${VERSION}${NC}"
    exit 0
}

# Log functions
log_info() {
    printf "%b\n" "${BLUE}[INFO]${NC} $1"
}

log_warning() {
    printf "%b\n" "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    printf "%b\n" "${RED}[ERROR]${NC} $1"
}

log_success() {
    printf "%b\n" "${GREEN}[SUCCESS]${NC} $1"
}

# Select available Docker Compose command (prefer docker compose, fallback to docker-compose)
DOCKER_COMPOSE_BIN=""
DOCKER_COMPOSE_SUBCMD=""

detect_compose_cmd() {
	# Prefer Docker Compose plugin
	if docker compose version &> /dev/null; then
		DOCKER_COMPOSE_BIN="docker"
		DOCKER_COMPOSE_SUBCMD="compose"
		return 0
	fi

	# Fallback to docker-compose (v1)
	if command -v docker-compose &> /dev/null; then
		if docker-compose version &> /dev/null; then
			DOCKER_COMPOSE_BIN="docker-compose"
			DOCKER_COMPOSE_SUBCMD=""
			return 0
		fi
	fi

	# Both unavailable
	return 1
}

# Check and create .env file
check_env_file() {
    log_info "Checking environment variable configuration..."
    if [ ! -f "$PROJECT_ROOT/.env" ]; then
        log_warning ".env file does not exist, creating from template"
        if [ -f "$PROJECT_ROOT/.env.example" ]; then
            cp "$PROJECT_ROOT/.env.example" "$PROJECT_ROOT/.env"
            log_success "Created .env file from .env.example"
        else
            log_error ".env.example template file not found, cannot create .env file"
            return 1
        fi
    else
        log_info ".env file already exists"
    fi
    
    # Check if required environment variables are set
    source "$PROJECT_ROOT/.env"
    local missing_vars=()
    
    # Check basic variables
    if [ -z "$DB_DRIVER" ]; then missing_vars+=("DB_DRIVER"); fi
    if [ -z "$STORAGE_TYPE" ]; then missing_vars+=("STORAGE_TYPE"); fi
    
    return 0
}

# Install Ollama (different methods for different platforms)
install_ollama() {
    # Check if it's a remote service
    get_ollama_base_url
    
    if [ $IS_REMOTE -eq 1 ]; then
        log_info "Remote Ollama service configuration detected, no need to install Ollama locally"
        return 0
    fi

    log_info "Local Ollama not installed, installing..."
    
    OS=$(uname)
    if [ "$OS" = "Darwin" ]; then
        # Mac installation method
        log_info "Mac system detected, using brew to install Ollama..."
        if ! command -v brew &> /dev/null; then
            # Install via package
            log_info "Homebrew not installed, using direct download method..."
            curl -fsSL https://ollama.com/download/Ollama-darwin.zip -o ollama.zip
            unzip ollama.zip
            mv ollama /usr/local/bin
            rm ollama.zip
        else
            brew install ollama
        fi
    else
        # Linux installation method
        log_info "Linux system detected, using installation script..."
        curl -fsSL https://ollama.com/install.sh | sh
    fi
    
    if [ $? -eq 0 ]; then
        log_success "Local Ollama installation completed"
        return 0
    else
        log_error "Local Ollama installation failed"
        return 1
    fi
}

# Get Ollama base URL, check if it's a remote service
get_ollama_base_url() {

    check_env_file

    # Get Ollama base URL from environment variable
    OLLAMA_URL=${OLLAMA_BASE_URL:-"http://host.docker.internal:11434"}
    # Extract host part
    OLLAMA_HOST=$(echo "$OLLAMA_URL" | sed -E 's|^https?://||' | sed -E 's|:[0-9]+$||' | sed -E 's|/.*$||')
    # Extract port part
    OLLAMA_PORT=$(echo "$OLLAMA_URL" | grep -oE ':[0-9]+' | grep -oE '[0-9]+' || echo "11434")
    # Check if it's localhost or 127.0.0.1
    IS_REMOTE=0
    if [ "$OLLAMA_HOST" = "localhost" ] || [ "$OLLAMA_HOST" = "127.0.0.1" ] || [ "$OLLAMA_HOST" = "host.docker.internal" ]; then
        IS_REMOTE=0  # Local service
    else
        IS_REMOTE=1  # Remote service
    fi
}

# Start Ollama service
start_ollama() {
    log_info "Checking Ollama service..."
    # Extract host and port
    get_ollama_base_url
    log_info "Ollama service address: $OLLAMA_URL"
    
    if [ $IS_REMOTE -eq 1 ]; then
        log_info "Remote Ollama service detected, will use remote service directly, no local installation or startup"
        # Check if remote service is available
        if curl -s "$OLLAMA_URL/api/tags" &> /dev/null; then
            log_success "Remote Ollama service is accessible"
            return 0
        else
            log_warning "Remote Ollama service is not accessible, please confirm the service address is correct and started"
            return 1
        fi
    fi
    
    # Following is for local service handling
    # Check if Ollama is installed
    if ! command -v ollama &> /dev/null; then
        install_ollama
        if [ $? -ne 0 ]; then
            return 1
        fi
    fi

    # Check if Ollama service is running
    if curl -s "http://localhost:$OLLAMA_PORT/api/tags" &> /dev/null; then
        log_success "Local Ollama service is already running, port: $OLLAMA_PORT"
    else
        log_info "Starting local Ollama service..."
        # Note: Official recommendation is to use systemctl or launchctl to manage services, direct background running is only for temporary scenarios
        systemctl restart ollama || (ollama serve > /dev/null 2>&1 < /dev/null &)
        
        # Wait for service to start
        MAX_RETRIES=30
        COUNT=0
        while [ $COUNT -lt $MAX_RETRIES ]; do
            if curl -s "http://localhost:$OLLAMA_PORT/api/tags" &> /dev/null; then
                log_success "Local Ollama service started successfully, port: $OLLAMA_PORT"
                break
            fi
            echo -ne "Waiting for Ollama service to start... ($COUNT/$MAX_RETRIES)\r"
            sleep 1
            COUNT=$((COUNT + 1))
        done
        echo "" # New line
        
        if [ $COUNT -eq $MAX_RETRIES ]; then
            log_error "Local Ollama service startup failed"
            return 1
        fi
    fi

    log_success "Local Ollama service address: http://localhost:$OLLAMA_PORT"
    return 0
}

# Stop Ollama service
stop_ollama() {
    log_info "Stopping Ollama service..."
    
    # Check if it's a remote service
    get_ollama_base_url
    
    if [ $IS_REMOTE -eq 1 ]; then
        log_info "Remote Ollama service detected, no need to stop locally"
        return 0
    fi
    
    # Check if Ollama is installed
    if ! command -v ollama &> /dev/null; then
        log_info "Local Ollama not installed, no need to stop"
        return 0
    fi
    
    # Find and terminate Ollama process
    if pgrep -x "ollama" > /dev/null; then
        # Prefer systemctl
        if command -v systemctl &> /dev/null; then
            sudo systemctl stop ollama
        else
            pkill -f "ollama serve"
        fi
        log_success "Local Ollama service stopped"
    else
        log_info "Local Ollama service is not running"
    fi
    
    return 0
}

# Check if Docker is installed
check_docker() {
    log_info "Checking Docker environment..."
    
    if ! command -v docker &> /dev/null; then
        log_error "Docker is not installed, please install Docker first"
        return 1
    fi
    
	# Check and select available Docker Compose command
	if detect_compose_cmd; then
		if [ "$DOCKER_COMPOSE_BIN" = "docker" ]; then
			log_info "Docker Compose plugin detected (docker compose)"
		else
			log_info "docker-compose (v1) detected"
		fi
	else
		log_error "Docker Compose not detected (neither docker compose nor docker-compose found). Please install one of them."
		return 1
	fi
    
    # Check Docker service running status
    if ! docker info &> /dev/null; then
        log_error "Docker service is not running, please start Docker service"
        return 1
    fi
    
    log_success "Docker environment check passed"
    return 0
}

check_platform() {
     # Detect current system platform
    log_info "Detecting system platform information..."
    if [ "$(uname -m)" = "x86_64" ]; then
        export PLATFORM="linux/amd64"
    elif [ "$(uname -m)" = "aarch64" ] || [ "$(uname -m)" = "arm64" ]; then
        export PLATFORM="linux/arm64"
    else
        log_warning "Unrecognized platform type: $(uname -m), will use default platform linux/amd64"
        export PLATFORM="linux/amd64"
    fi
    log_info "Current platform: $PLATFORM"
}

# Start Docker containers
start_docker() {
    log_info "Starting Docker containers..."
    
    # Check Docker environment
    check_docker
    if [ $? -ne 0 ]; then
        return 1
    fi
    
    # Check .env file
    check_env_file
    
    # Read .env file
    source "$PROJECT_ROOT/.env"
    storage_type=${STORAGE_TYPE:-local}
    
    check_platform
    
    # Change to project root directory before executing docker-compose commands
    cd "$PROJECT_ROOT"
    
    # Start basic services
    log_info "Starting core service containers..."
	# Start using detected Compose command
	if [ "$NO_PULL" = true ]; then
		# Don't pull images, use local images
		log_info "Skipping image pull, using local images..."
		PLATFORM=$PLATFORM "$DOCKER_COMPOSE_BIN" $DOCKER_COMPOSE_SUBCMD up --build -d
	else
		# Pull latest images
		log_info "Pulling latest images..."
		PLATFORM=$PLATFORM "$DOCKER_COMPOSE_BIN" $DOCKER_COMPOSE_SUBCMD up --pull always -d
	fi
    if [ $? -ne 0 ]; then
        log_error "Docker container startup failed"
        return 1
    fi
    
    log_success "All Docker containers started successfully"
    
    # Show container status
    log_info "Current container status:"
	"$DOCKER_COMPOSE_BIN" $DOCKER_COMPOSE_SUBCMD ps
    
    return 0
}

# Stop Docker containers
stop_docker() {
    log_info "Stopping Docker containers..."
    
    # Check Docker environment
    check_docker
    if [ $? -ne 0 ]; then
        # Try to stop even if check fails, just in case
        log_warning "Docker environment check failed, will still attempt to stop containers..."
    fi
    
    # Change to project root directory before executing docker-compose commands
    cd "$PROJECT_ROOT"
    
    # Stop all containers
	"$DOCKER_COMPOSE_BIN" $DOCKER_COMPOSE_SUBCMD down --remove-orphans
    if [ $? -ne 0 ]; then
        log_error "Docker container stop failed"
        return 1
    fi
    
    log_success "All Docker containers stopped"
    return 0
}

# List all running containers
list_containers() {
    log_info "Listing all running containers..."
    
    # Check Docker environment
    check_docker
    if [ $? -ne 0 ]; then
        return 1
    fi
    
    # Change to project root directory before executing docker-compose commands
    cd "$PROJECT_ROOT"
    
    # List all containers
    printf "%b\n" "${BLUE}Currently running containers:${NC}"
	"$DOCKER_COMPOSE_BIN" $DOCKER_COMPOSE_SUBCMD ps --services | sort
    
    return 0
}

# Pull latest Docker images
pull_images() {
    log_info "Pulling latest Docker images..."
    
    # Check Docker environment
    check_docker
    if [ $? -ne 0 ]; then
        return 1
    fi
    
    # Check .env file
    check_env_file
    
    # Read .env file
    source "$PROJECT_ROOT/.env"
    storage_type=${STORAGE_TYPE:-local}
    
    check_platform
    
    # Change to project root directory before executing docker-compose commands
    cd "$PROJECT_ROOT"
    
    # Pull all images
    log_info "Pulling latest images for all services..."
	PLATFORM=$PLATFORM "$DOCKER_COMPOSE_BIN" $DOCKER_COMPOSE_SUBCMD pull
    if [ $? -ne 0 ]; then
        log_error "Image pull failed"
        return 1
    fi
    
    log_success "All images pulled to latest version successfully"
    
    # Show pulled image information
    log_info "Pulled images:"
    docker images --format "table {{.Repository}}\t{{.Tag}}\t{{.CreatedAt}}\t{{.Size}}" | head -10
    
    return 0
}

# Restart specified container
restart_container() {
    local container_name="$1"
    
    if [ -z "$container_name" ]; then
        log_error "Container name not specified"
        echo "Available containers:"
        list_containers
        return 1
    fi
    
    log_info "Rebuilding and restarting container: $container_name"
    
    # Check Docker environment
    check_docker
    if [ $? -ne 0 ]; then
        return 1
    fi
    
    check_platform
    
    # Change to project root directory before executing docker-compose commands
    cd "$PROJECT_ROOT"
    
    # Check if container exists
	if ! "$DOCKER_COMPOSE_BIN" $DOCKER_COMPOSE_SUBCMD ps --services | grep -q "^$container_name$"; then
        log_error "Container '$container_name' does not exist or is not running"
        echo "Available containers:"
        list_containers
        return 1
    fi
    
    # Build and restart container
    log_info "Rebuilding container '$container_name'..."
	PLATFORM=$PLATFORM "$DOCKER_COMPOSE_BIN" $DOCKER_COMPOSE_SUBCMD build "$container_name"
    if [ $? -ne 0 ]; then
        log_error "Container '$container_name' build failed"
        return 1
    fi
    
    log_info "Restarting container '$container_name'..."
	PLATFORM=$PLATFORM "$DOCKER_COMPOSE_BIN" $DOCKER_COMPOSE_SUBCMD up -d --no-deps "$container_name"
    if [ $? -ne 0 ]; then
        log_error "Container '$container_name' restart failed"
        return 1
    fi
    
    log_success "Container '$container_name' rebuilt and restarted successfully"
    return 0
}

# Check system environment
check_environment() {
    log_info "Starting environment check..."
    
    # Check operating system
    OS=$(uname)
    log_info "Operating system: $OS"
    
    # Check Docker
    check_docker
    
    # Check .env file
    check_env_file
    
    get_ollama_base_url
    
    if [ $IS_REMOTE -eq 1 ]; then
        log_info "Remote Ollama service configuration detected"
        if curl -s "$OLLAMA_URL/api/tags" &> /dev/null; then
            version=$(curl -s "$OLLAMA_URL/api/tags" | grep -o '"version":"[^"]*"' | cut -d'"' -f4)
            log_success "Remote Ollama service is accessible, version: $version"
        else
            log_warning "Remote Ollama service is not accessible, please confirm the service address is correct and started"
        fi
    else
        if command -v ollama &> /dev/null; then
            log_success "Local Ollama is installed"
            if curl -s "http://localhost:$OLLAMA_PORT/api/tags" &> /dev/null; then
                version=$(curl -s "http://localhost:$OLLAMA_PORT/api/tags" | grep -o '"version":"[^"]*"' | cut -d'"' -f4)
                log_success "Local Ollama service is running, version: $version"
            else
                log_warning "Local Ollama is installed but service is not running"
            fi
        else
            log_warning "Local Ollama is not installed"
        fi
    fi
    
    # Check disk space
    log_info "Checking disk space..."
    df -h | grep -E "(Filesystem|/$)"
    
    # Check memory
    log_info "Checking memory usage..."
    if [ "$OS" = "Darwin" ]; then
        vm_stat | perl -ne '/page size of (\d+)/ and $size=$1; /Pages free:\s*(\d+)/ and print "Free Memory: ", $1 * $size / 1048576, " MB\n"'
    else
        free -h | grep -E "(total|Mem:)"
    fi
    
    # Check CPU
    log_info "CPU information:"
    if [ "$OS" = "Darwin" ]; then
        sysctl -n machdep.cpu.brand_string
        echo "CPU cores: $(sysctl -n hw.ncpu)"
    else
        grep "model name" /proc/cpuinfo | head -1
        echo "CPU cores: $(nproc)"
    fi
    
    # Check container status
    log_info "Checking container status..."
    if docker info &> /dev/null; then
        docker ps -a
    else
        log_warning "Unable to get container status, Docker may not be running"
    fi
    
    log_success "Environment check completed"
    return 0
}

# Parse command line arguments
START_OLLAMA=false
START_DOCKER=false
STOP_SERVICES=false
CHECK_ENVIRONMENT=false
LIST_CONTAINERS=false
RESTART_CONTAINER=false
PULL_IMAGES=false
NO_PULL=false
CONTAINER_NAME=""

# Default to start all services when no arguments provided
if [ $# -eq 0 ]; then
    START_OLLAMA=true
    START_DOCKER=true
fi

while [ "$1" != "" ]; do
    case $1 in
        -h | --help )       show_help
                            ;;
        -o | --ollama )     START_OLLAMA=true
                            ;;
        -d | --docker )     START_DOCKER=true
                            ;;
        -a | --all )        START_OLLAMA=true
                            START_DOCKER=true
                            ;;
        -s | --stop )       STOP_SERVICES=true
                            ;;
        -c | --check )      CHECK_ENVIRONMENT=true
                            ;;
        -l | --list )       LIST_CONTAINERS=true
                            ;;
        -p | --pull )       PULL_IMAGES=true
                            ;;
        --no-pull )         NO_PULL=true
                            START_OLLAMA=true
                            START_DOCKER=true
                            ;;
        -r | --restart )    RESTART_CONTAINER=true
                            CONTAINER_NAME="$2"
                            shift
                            ;;
        -v | --version )    show_version
                            ;;
        * )                 log_error "Unknown option: $1"
                            show_help
                            ;;
    esac
    shift
done

# Execute environment check
if [ "$CHECK_ENVIRONMENT" = true ]; then
    check_environment
    exit $?
fi

# List all containers
if [ "$LIST_CONTAINERS" = true ]; then
    list_containers
    exit $?
fi

# Pull latest images
if [ "$PULL_IMAGES" = true ]; then
    pull_images
    exit $?
fi

# Restart specified container
if [ "$RESTART_CONTAINER" = true ]; then
    restart_container "$CONTAINER_NAME"
    exit $?
fi

# Execute service operations
if [ "$STOP_SERVICES" = true ]; then
    # Stop services
    stop_ollama
    OLLAMA_RESULT=$?
    
    stop_docker
    DOCKER_RESULT=$?
    
    # Show summary
    echo ""
    log_info "=== Stop Results ==="
    if [ $OLLAMA_RESULT -eq 0 ]; then
        log_success "✓ Ollama service stopped"
    else
        log_error "✗ Ollama service stop failed"
    fi
    
    if [ $DOCKER_RESULT -eq 0 ]; then
        log_success "✓ Docker containers stopped"
    else
        log_error "✗ Docker container stop failed"
    fi
    
    log_success "Service stop completed."
else
    # Start services
    OLLAMA_RESULT=1
    DOCKER_RESULT=1
    if [ "$START_OLLAMA" = true ]; then
        start_ollama
        OLLAMA_RESULT=$?
    fi
    
    if [ "$START_DOCKER" = true ]; then
        start_docker
        DOCKER_RESULT=$?
    fi
    
    # Show summary
    echo ""
    log_info "=== Startup Results ==="
    if [ "$START_OLLAMA" = true ]; then
        if [ $OLLAMA_RESULT -eq 0 ]; then
            log_success "✓ Ollama service started"
        else
            log_error "✗ Ollama service startup failed"
        fi
    fi
    
    if [ "$START_DOCKER" = true ]; then
        if [ $DOCKER_RESULT -eq 0 ]; then
            log_success "✓ Docker containers started"
        else
            log_error "✗ Docker container startup failed"
        fi
    fi
    
    if [ "$START_OLLAMA" = true ] && [ "$START_DOCKER" = true ]; then
        if [ $OLLAMA_RESULT -eq 0 ] && [ $DOCKER_RESULT -eq 0 ]; then
            log_success "All services started successfully, accessible at:"
            printf "%b\n" "${GREEN}  - Frontend: http://localhost:${FRONTEND_PORT:-80}${NC}"
            printf "%b\n" "${GREEN}  - API: http://localhost:${APP_PORT:-8080}${NC}"
            printf "%b\n" "${GREEN}  - Jaeger Tracing: http://localhost:16686${NC}"
            echo ""
            log_info "Continuously outputting container logs (Press Ctrl+C to exit logs, containers will not stop)..."
            "$DOCKER_COMPOSE_BIN" $DOCKER_COMPOSE_SUBCMD logs app docreader postgres --since=10s -f
        else
            log_error "Some services failed to start, please check logs and fix issues"
        fi
    elif [ "$START_OLLAMA" = true ] && [ $OLLAMA_RESULT -eq 0 ]; then
        log_success "Ollama service started successfully, accessible at:"
        printf "%b\n" "${GREEN}  - Ollama API: http://localhost:$OLLAMA_PORT${NC}"
    elif [ "$START_DOCKER" = true ] && [ $DOCKER_RESULT -eq 0 ]; then
        log_success "Docker containers started successfully, accessible at:"
        printf "%b\n" "${GREEN}  - Frontend: http://localhost:${FRONTEND_PORT:-80}${NC}"
        printf "%b\n" "${GREEN}  - API: http://localhost:${APP_PORT:-8080}${NC}"
        printf "%b\n" "${GREEN}  - Jaeger Tracing: http://localhost:16686${NC}"
        echo ""
        log_info "Continuously outputting container logs (Press Ctrl+C to exit logs, containers will not stop)..."
        "$DOCKER_COMPOSE_BIN" $DOCKER_COMPOSE_SUBCMD logs app docreader postgres --since=10s -f
    fi
fi

exit 0