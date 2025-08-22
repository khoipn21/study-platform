# Study Platform - Unique Port Architecture Documentation

## Overview

This document outlines the **unique port architecture** for the Study Platform microservices. Each service has been assigned a **unique internal port** to prevent conflicts, with only the API Gateway exposed externally.

## Architecture Principles

1. **🔧 Unique Internal Ports**: Each service uses a **unique port** internally within Docker containers
2. **🌐 Single External Entry Point**: Only API Gateway is exposed externally on port **8080**
3. **🔒 Internal Service Communication**: All inter-service communication happens through Docker's internal network
4. **📝 Consistent Environment Variables**: Standardized environment variable naming for port configuration
5. **🚫 No Port Conflicts**: Each service has its own dedicated port range

## Port Allocation Summary

### **Port Allocation Plan**

| Service | Internal Port | External Port | Environment Variable | Protocol | Status |
|---------|---------------|---------------|---------------------|----------|--------|
| **API Gateway** | **8080** | **8080** | `HTTP_PORT` | HTTP | ✅ **ONLY EXTERNAL ACCESS** |
| **Auth Service** | **8081** | - | `GRPC_PORT` | gRPC | ❌ Internal only |
| **Course Service** | **8082** | - | `GRPC_PORT` | gRPC | ❌ Internal only |
| **Progress Service** | **8083** | - | `GRPC_PORT` | gRPC | ❌ Internal only |
| **Video Service** | **8084** | - | `VIDEO_SERVICE_PORT` | HTTP/WebSocket | ❌ Internal only |
| **Bucket Service** | **8085** | - | `BUCKET_SERVICE_PORT` | HTTP | ❌ Internal only |
| **Chatbot Service** | **8086** | - | `CHATBOT_PORT` | HTTP/WebSocket | ❌ Internal only |
| **Forum Service** | **8087** | - | `FORUM_PORT` | HTTP | ❌ Internal only |
| **Payment Service** | **8088** | - | `PAYMENT_PORT` | HTTP | ❌ Internal only |

### **Infrastructure Services**
| Service | Internal Port | External Port | Purpose |
|---------|---------------|---------------|---------|
| PostgreSQL | 5432 | 2345 | Database |
| Redis | 6379 | 6379 | Cache/Message Queue |
| MinIO | 9000/9001 | 9000/9001 | Object Storage |

## Internal Service Communication

### API Gateway Service URLs

The API Gateway routes to all services using these internal URLs:

```yaml
AUTH_SERVICE_URL: auth-service:8081
COURSE_SERVICE_URL: course-service:8082
PROGRESS_SERVICE_URL: progress-service:8083
VIDEO_SERVICE_URL: http://video-service:8084
BUCKET_SERVICE_URL: http://bucket-service:8085
CHATBOT_SERVICE_URL: http://chatbot-service:8086
FORUM_SERVICE_URL: http://forum-service:8087
PAYMENT_SERVICE_URL: http://payment-service:8088
```

### Service Communication Matrix

| From Service | To Service | URL | Protocol |
|-------------|-----------|-----|----------|
| API Gateway | Auth Service | `auth-service:8081` | gRPC |
| API Gateway | Course Service | `course-service:8082` | gRPC |
| API Gateway | Progress Service | `progress-service:8083` | gRPC |
| API Gateway | Video Service | `http://video-service:8084` | HTTP |
| API Gateway | Bucket Service | `http://bucket-service:8085` | HTTP |
| API Gateway | Chatbot Service | `http://chatbot-service:8086` | HTTP |
| API Gateway | Forum Service | `http://forum-service:8087` | HTTP |
| API Gateway | Payment Service | `http://payment-service:8088` | HTTP |
| Payment Service | Progress Service | `progress-service:8083` | gRPC |

## Service Details

### 1. API Gateway (Port 8080)
- **Type**: HTTP Server (Gorilla Mux)
- **Internal Port**: 8080
- **External Port**: 8080 (EXPOSED)
- **Environment Variable**: `HTTP_PORT=8080`
- **Dockerfile**: `EXPOSE 8080`
- **Role**: Routes requests to all internal services
- **Health Check**: Available at `/api/v1/health`

### 2. Auth Service (Port 8081)
- **Type**: gRPC Server
- **Internal Port**: 8081
- **Environment Variable**: `GRPC_PORT=8081`
- **Dockerfile**: `EXPOSE 8081`
- **Health Check**: `grpc_health_probe -addr=:8081`
- **Features**: JWT authentication, OAuth 2.0, user management

### 3. Course Service (Port 8082)
- **Type**: gRPC Server
- **Internal Port**: 8082
- **Environment Variable**: `GRPC_PORT=8082`
- **Dockerfile**: `EXPOSE 8082`
- **Health Check**: `grpc_health_probe -addr=:8082`
- **Features**: Course CRUD, enrollment, ratings, search

### 4. Progress Service (Port 8083)
- **Type**: gRPC Server
- **Internal Port**: 8083
- **Environment Variable**: `GRPC_PORT=8083`
- **Dockerfile**: `EXPOSE 8083`
- **Health Check**: `grpc_health_probe -addr=:8083`
- **Features**: Learning progress tracking, analytics

### 5. Video Service (Port 8084)
- **Type**: HTTP Server (Gin) + WebSocket
- **Internal Port**: 8084
- **Environment Variable**: `VIDEO_SERVICE_PORT=8084`
- **Dockerfile**: `EXPOSE 8084`
- **Health Check**: `wget http://localhost:8084/health`
- **Features**: Cloudflare Stream integration, adaptive streaming, WebSocket

### 6. Bucket Service (Port 8085)
- **Type**: HTTP Server (Gin)
- **Internal Port**: 8085
- **Environment Variable**: `BUCKET_SERVICE_PORT=8085`
- **Dockerfile**: `EXPOSE 8085`
- **Health Check**: `wget http://localhost:8085/health`
- **Features**: AWS S3 integration, file upload/download

### 7. Chatbot Service (Port 8086)
- **Type**: HTTP Server (Gin) + WebSocket
- **Internal Port**: 8086
- **Environment Variable**: `CHATBOT_PORT=8086`
- **Dockerfile**: `EXPOSE 8086`
- **Health Check**: `wget http://localhost:8086/health`
- **Features**: OpenAI integration, real-time chat

### 8. Forum Service (Port 8087)
- **Type**: HTTP Server (Gin)
- **Internal Port**: 8087
- **Environment Variable**: `FORUM_PORT=8087`
- **Dockerfile**: `EXPOSE 8087`
- **Health Check**: `wget http://localhost:8087/health`
- **Features**: Discussion forums, voting system

### 9. Payment Service (Port 8088)
- **Type**: HTTP Server (Gin)
- **Internal Port**: 8088
- **Environment Variable**: `PAYMENT_PORT=8088`
- **Dockerfile**: `EXPOSE 8088`
- **Health Check**: `wget http://localhost:8088/health`
- **Features**: Stripe/PayPal integration, subscription management

## Docker Configuration

### Docker Compose Port Mapping

```yaml
services:
  # Only API Gateway exposed externally
  api-gateway:
    ports:
      - "8080:8080"  # ONLY EXTERNAL ACCESS POINT
  
  # All other services: internal only
  auth-service:
    ports: []
    environment:
      GRPC_PORT: 8081
  
  course-service:
    ports: []
    environment:
      GRPC_PORT: 8082
  
  progress-service:
    ports: []
    environment:
      GRPC_PORT: 8083
  
  video-service:
    ports: []
    environment:
      VIDEO_SERVICE_PORT: 8084
  
  bucket-service:
    ports: []
    environment:
      BUCKET_SERVICE_PORT: 8085
  
  chatbot-service:
    ports: []
    environment:
      CHATBOT_PORT: 8086
  
  forum-service:
    ports: []
    environment:
      FORUM_PORT: 8087
  
  payment-service:
    ports: []
    environment:
      PAYMENT_PORT: 8088
```

### Health Check Configuration

Each service has its own health check configuration:

#### gRPC Services (Auth, Course, Progress)
```yaml
healthcheck:
  test: ["CMD", "grpc_health_probe", "-addr=:808X"]
  interval: 30s
  timeout: 10s
  retries: 3
```

#### HTTP Services (Video, Bucket, Chatbot, Forum, Payment)
```yaml
healthcheck:
  test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:808X/health"]
  interval: 30s
  timeout: 10s
  retries: 3
```

## API Access Patterns

### External Client Access
```bash
# All external requests go through API Gateway on port 8080
curl http://localhost:8080/api/v1/auth/login
curl http://localhost:8080/api/v1/courses
curl http://localhost:8080/api/v1/videos/search
curl http://localhost:8080/api/v1/files/upload
curl http://localhost:8080/api/v1/forum/topics
curl http://localhost:8080/api/v1/payments/methods
```

### Development Testing
```bash
# Health check through API Gateway
curl http://localhost:8080/api/v1/health

# Service-specific health checks (internal container access only)
docker-compose exec video-service wget http://localhost:8084/health
docker-compose exec bucket-service wget http://localhost:8085/health
```

## Security Benefits

1. **🔒 Reduced Attack Surface**: Only one port exposed externally
2. **🛡️ Centralized Security**: All security policies applied at API Gateway
3. **🔍 Unified Monitoring**: Single entry point for all traffic monitoring
4. **🚫 Service Isolation**: Internal services cannot be accessed directly from outside
5. **📝 Consistent Logging**: All requests logged through API Gateway
6. **⚡ Port Conflict Prevention**: Unique ports eliminate configuration conflicts

## Configuration Files Updated

### Service Configurations
- `/auth-service/cmd/main.go`: `GRPC_PORT=8081`
- `/course-service/cmd/main.go`: `GRPC_PORT=8082`
- `/progress-service/cmd/main.go`: `GRPC_PORT=8083`
- `/video-service/internal/config/config.go`: `VIDEO_SERVICE_PORT=8084`
- `/bucket-service/internal/config/config.go`: `BUCKET_SERVICE_PORT=8085`
- `/chatbot-service/internal/config/config.go`: `CHATBOT_PORT=8086`
- `/forum-service/internal/config/config.go`: `FORUM_PORT=8087`
- `/payment-service/internal/config/config.go`: `PAYMENT_PORT=8088`

### Docker Files
- `/Dockerfile.auth`: `EXPOSE 8081`
- `/Dockerfile.course`: `EXPOSE 8082`
- `/Dockerfile.progress`: `EXPOSE 8083`
- `/video-service/Dockerfile`: `EXPOSE 8084`
- `/bucket-service/Dockerfile`: `EXPOSE 8085`
- `/chatbot-service/Dockerfile`: `EXPOSE 8086`
- `/forum-service/Dockerfile`: `EXPOSE 8087`
- `/Dockerfile.payment`: `EXPOSE 8088`

### API Gateway Configuration
- `/api-gateway/cmd/main.go`: Updated service URL defaults
- `/docker-compose.yml`: Complete port mapping configuration

## Development Guidelines

### Adding New Services
1. **Port Assignment**: Use next available port in sequence (8089, 8090, etc.)
2. **Dockerfile**: Include `EXPOSE 808X`
3. **Environment**: Use `{SERVICE_NAME}_PORT=808X`
4. **Health Check**: Implement `/health` endpoint on assigned port
5. **API Gateway**: Add route configuration and service URL
6. **Docker Compose**: Set `ports: []` (no external ports)

### Service Communication Best Practices
```go
// ✅ Correct: Use Docker service names with unique ports
baseURL := "http://video-service:8084"
grpcConn := "auth-service:8081"

// ❌ Incorrect: Don't use localhost or default ports
baseURL := "http://localhost:8080"
grpcConn := "localhost:8080"
```

## Troubleshooting

### Port Conflict Resolution
1. ✅ **No Conflicts**: Each service has unique internal port
2. ✅ **Clear Separation**: gRPC services (8081-8083), HTTP services (8084-8088)
3. ✅ **Easy Identification**: Port number matches service order

### Service Communication Issues
1. Use Docker service names, not `localhost`
2. Use correct unique port for each service
3. Check Docker network connectivity: `docker-compose exec api-gateway ping video-service`
4. Verify environment variables: `docker-compose exec video-service env | grep PORT`

### Health Check Failures
1. Ensure health endpoints are available on correct unique ports
2. Check service logs: `docker-compose logs [service-name]`
3. Verify port configuration in service code and Dockerfile

### Testing Internal Ports
```bash
# Check if services are listening on correct ports
docker-compose exec auth-service netstat -tlnp | grep 8081
docker-compose exec course-service netstat -tlnp | grep 8082
docker-compose exec progress-service netstat -tlnp | grep 8083
docker-compose exec video-service netstat -tlnp | grep 8084
docker-compose exec bucket-service netstat -tlnp | grep 8085
docker-compose exec chatbot-service netstat -tlnp | grep 8086
docker-compose exec forum-service netstat -tlnp | grep 8087
docker-compose exec payment-service netstat -tlnp | grep 8088
```

## Migration from Previous Architecture

### Before: Unified Port Architecture (ALL SERVICES on 8080)
- **Issue**: All services used port 8080 internally
- **Problem**: Potential conflicts and confusion
- **Solution**: Assigned unique ports to each service

### After: Unique Port Architecture (EACH SERVICE has UNIQUE PORT)
- **Auth Service**: `8080` → `8081`
- **Course Service**: `8080` → `8082`
- **Progress Service**: `8080` → `8083`
- **Video Service**: `8080` → `8084`
- **Bucket Service**: `8080` → `8085`
- **Chatbot Service**: `8080` → `8086`
- **Forum Service**: `8080` → `8087`
- **Payment Service**: `8080` → `8088`
- **API Gateway**: Remains `8080` (external access)

## Testing Commands

```bash
# Verify external access works
curl http://localhost:8080/api/v1/health

# Test individual services through API Gateway
curl http://localhost:8080/api/v1/auth/login -d '{"email":"test@example.com","password":"password"}'
curl http://localhost:8080/api/v1/courses
curl http://localhost:8080/api/v1/videos/search?q=test
curl http://localhost:8080/api/v1/forum/topics
curl http://localhost:8080/api/v1/payments/methods

# Verify unique internal ports
docker-compose exec auth-service netstat -tlnp | grep 8081
docker-compose exec video-service netstat -tlnp | grep 8084
docker-compose exec payment-service netstat -tlnp | grep 8088

# Check service environment variables
docker-compose exec auth-service env | grep GRPC_PORT
docker-compose exec video-service env | grep VIDEO_SERVICE_PORT
docker-compose exec payment-service env | grep PAYMENT_PORT
```

---

**Last Updated**: August 2024  
**Architecture**: Microservices with unique internal ports  
**External Access**: API Gateway only (localhost:8080)  
**Port Range**: 8081-8088 (Internal Services), 8080 (External Gateway)  
**Security**: Single entry point, isolated services