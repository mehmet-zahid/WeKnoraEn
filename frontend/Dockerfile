# Build stage
FROM oven/bun:alpine as build-stage

WORKDIR /app

# Set environment variables
ENV NODE_OPTIONS="--max-old-space-size=4096"
ENV VITE_IS_DOCKER=true

# File size limit (MB), passed via build argument
ARG MAX_FILE_SIZE_MB=50
ENV VITE_MAX_FILE_SIZE_MB=${MAX_FILE_SIZE_MB}

# Copy dependency files
COPY package*.json ./
COPY packages/xlsx-0.20.2.tgz ./packages/xlsx-0.20.2.tgz

# Install dependencies
RUN bun install

# Copy project files
COPY . .

# Build application
RUN bun run build

# Production stage
FROM nginx:stable-alpine as production-stage

# Copy build artifacts to nginx service directory
COPY --from=build-stage /app/dist /usr/share/nginx/html

# Copy nginx configuration template file
COPY nginx.conf /etc/nginx/templates/default.conf.template

# Set default environment variable (MB)
ENV MAX_FILE_SIZE_MB=50

# Expose port
EXPOSE 80

# Convert MAX_FILE_SIZE_MB to MAX_FILE_SIZE with unit on startup, then replace in nginx config
CMD ["/bin/sh", "-c", "export MAX_FILE_SIZE=${MAX_FILE_SIZE_MB}M && envsubst '${MAX_FILE_SIZE}' < /etc/nginx/templates/default.conf.template > /etc/nginx/conf.d/default.conf && nginx -g 'daemon off;'"] 