# Study Platform - Microservices Architecture

A comprehensive online learning platform built with Go microservices architecture, featuring authentication, course management, progress tracking, and a robust API gateway.

## 🏗️ Architecture Overview

The Study Platform consists of multiple microservices that communicate via gRPC, with an HTTP API Gateway serving as the unified entry point:

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Web Client    │    │  Mobile App     │    │   Admin Panel   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                    ┌─────────────────┐
                    │   API Gateway   │ ← Rate Limiting, Auth, Circuit Breaker
                    │   (Port 8080)   │
                    └─────────────────┘
                                 │
        ┌────────────────────────┼────────────────────────┐
        │                       │                        │
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  Auth Service   │    │ Course Service  │    │Progress Service │
│   (Port 8081)   │    │   (Port 8082)   │    │   (Port 8083)   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
        │                       │                        │
        └───────────────────────┼────────────────────────┘
                                │
                    ┌─────────────────┐
                    │   PostgreSQL    │
                    │   (Port 2345)   │
                    └─────────────────┘
```

## 🚀 Quick Start

### Prerequisites

- **Docker & Docker Compose**: Required for containerized deployment
- **Go 1.23+**: For local development
- **Git**: For cloning the repository

### 1. Clone the Repository

```bash
git clone <repository-url>
cd Study-Platform
```

### 2. Start All Services

```bash
# Start all services in detached mode
docker-compose up -d

# Check service status
docker-compose ps
```

### 3. Verify Services are Running

```bash
# Check API Gateway health
curl http://localhost:8080/api/v1/health

# Check circuit breaker status
curl http://localhost:8080/api/v1/health/circuit-breakers
```

### 4. Access the API Documentation

Open your browser and navigate to:
- **Swagger UI**: http://localhost:8080/api/v1/docs
- **OpenAPI Spec**: http://localhost:8080/api/v1/docs/openapi.jsonV

## 📋 Services Overview

### 🔐 Auth Service (Port 8081)
**Handles user authentication and authorization**

**Features:**
- User registration and login
- JWT token generation and validation
- Role-based access control (Student, Instructor, Admin)
- OAuth 2.0 integration (Google, GitHub, Facebook)
- Password hashing with bcrypt
- User role management

**gRPC Methods:**
- `Register` - Register new user
- `Login` - User authentication
- `ValidateToken` - JWT token validation
- `GetUserRoles` - Get user roles
- `AssignRole` - Assign role to user
- `RemoveRole` - Remove role from user
- `GetOAuthURL` - Get OAuth authorization URL
- `OAuthCallback` - Handle OAuth callback
- `LinkOAuthAccount` - Link OAuth account
- `UnlinkOAuthAccount` - Unlink OAuth account
- `GetLinkedAccounts` - Get linked OAuth accounts

### 📚 Course Service (Port 8082)
**Manages courses and lectures**

**Features:**
- Course CRUD operations
- Lecture management with ordering
- Advanced search and filtering
- Course categorization and tagging
- Enrollment system
- Rating and review system
- Instructor course management

**gRPC Methods:**
- `CreateCourse` - Create new course
- `GetCourse` - Get course details
- `UpdateCourse` - Update course information
- `DeleteCourse` - Delete course
- `ListCourses` - List courses with pagination
- `SearchCourses` - Search courses with filters
- `CreateLecture` - Create new lecture
- `GetLecture` - Get lecture details
- `UpdateLecture` - Update lecture
- `DeleteLecture` - Delete lecture
- `ListLectures` - List course lectures
- `EnrollInCourse` - Enroll user in course
- `GetEnrollments` - Get user enrollments

### 📊 Progress Service (Port 8083)
**Tracks user progress and analytics**

**Features:**
- Progress tracking for lectures and courses
- Enrollment management
- Course completion tracking
- User analytics and insights
- Learning path recommendations
- Time tracking for lectures

**gRPC Methods:**
- `UpdateProgress` - Update user progress
- `GetProgress` - Get progress for specific content
- `GetUserProgress` - Get user's overall progress
- `GetLectureProgress` - Get lecture progress
- `GetCourseCompletion` - Get course completion status
- `MarkLectureComplete` - Mark lecture as complete
- `CreateEnrollment` - Create enrollment
- `GetEnrollment` - Get enrollment details
- `ListEnrollments` - List user enrollments
- `GetUserAnalytics` - Get user analytics
- `GetCourseAnalytics` - Get course analytics
- `GetLearningPath` - Get learning path recommendations

### 🌐 API Gateway (Port 8080)
**Unified HTTP API entry point**

**Features:**
- RESTful API endpoints
- JWT authentication middleware
- Rate limiting (100 req/sec, 200 burst)
- Circuit breaker pattern
- Automatic retry mechanisms
- Request/response logging
- CORS support
- OpenAPI documentation

**Key Endpoints:**
- `GET /api/v1/health` - Health check
- `GET /api/v1/health/circuit-breakers` - Circuit breaker status
- `GET /api/v1/docs` - Swagger UI
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login
- `GET /api/v1/courses` - List courses
- `POST /api/v1/courses` - Create course (auth required)
- `GET /api/v1/enrollments` - List enrollments (auth required)
- `POST /api/v1/enrollments` - Create enrollment (auth required)

## 🔧 Configuration

### Environment Variables

Each service can be configured using environment variables:

**Database Configuration:**
```bash
DATABASE_URL=postgres://admin:password@postgres:5432/studyplatform?sslmode=disable
```

**Service Ports:**
```bash
# Auth Service
GRPC_PORT=8080

# Course Service  
GRPC_PORT=8080

# Progress Service
GRPC_PORT=8080

# API Gateway
HTTP_PORT=8080
AUTH_SERVICE_URL=auth-service:8080
COURSE_SERVICE_URL=course-service:8080
PROGRESS_SERVICE_URL=progress-service:8080
```

**Authentication:**
```bash
JWT_SECRET=your-secret-key-here
```

### Docker Compose Services

The `docker-compose.yml` includes:
- **postgres**: PostgreSQL database
- **redis**: Redis cache
- **minio**: MinIO object storage
- **auth-service**: Authentication service
- **course-service**: Course management service
- **progress-service**: Progress tracking service
- **api-gateway**: HTTP API Gateway

## 🧪 Testing the API

### 1. User Registration

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "password123"
  }'
```

### 2. User Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

### 3. List Courses

```bash
curl -X GET http://localhost:8080/api/v1/courses
```

### 4. Create Course (Authentication Required)

```bash
curl -X POST http://localhost:8080/api/v1/courses \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "title": "Introduction to Go Programming",
    "description": "Learn the fundamentals of Go programming language",
    "category": "Programming",
    "level": "beginner",
    "price": 99.99,
    "currency": "USD",
    "tags": ["go", "programming", "backend"]
  }'
```

### 5. Enroll in Course

```bash
curl -X POST http://localhost:8080/api/v1/enrollments \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "course_id": "COURSE_ID"
  }'
```

### 6. Get User Enrollments

```bash
curl -X GET http://localhost:8080/api/v1/enrollments \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## 📊 Monitoring & Health Checks

### Service Health

```bash
# API Gateway health
curl http://localhost:8080/api/v1/health

# Circuit breaker status
curl http://localhost:8080/api/v1/health/circuit-breakers
```

### View Logs

```bash
# View all service logs
docker-compose logs

# View specific service logs
docker-compose logs auth-service
docker-compose logs course-service
docker-compose logs progress-service
docker-compose logs api-gateway

# Follow logs in real-time
docker-compose logs -f api-gateway
```

### Check Service Status

```bash
# Check running services
docker-compose ps

# Check resource usage
docker stats
```

## 🛠️ Development

### Local Development Setup

1. **Install Dependencies**
```bash
go mod download
```

2. **Start Database**
```bash
docker-compose up -d postgres redis minio
```

3. **Run Services Individually**
```bash
# Auth Service
cd auth-service
go run cmd/main.go

# Course Service
cd course-service
go run cmd/main.go

# Progress Service
cd progress-service
go run cmd/main.go

# API Gateway
cd api-gateway
go run cmd/main.go
```

### Database Migrations

Database migrations are automatically applied when services start. Migration files are located in the `migrations/` directory.

### Adding New Endpoints

1. **Update Protobuf Definitions** (for gRPC services)
2. **Generate Protobuf Code**
```bash
# From service directory
protoc --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  proto/*.proto
```

3. **Implement Service Methods**
4. **Add HTTP Handlers** (for API Gateway)
5. **Update Documentation**

## 🔒 Security Features

### Authentication & Authorization
- JWT-based authentication
- Role-based access control (RBAC)
- OAuth 2.0 integration
- Secure password hashing (bcrypt)

### API Security
- Rate limiting (100 requests/second per IP)
- Request/response logging
- CORS protection
- Input validation

### Service Communication
- gRPC with TLS support
- Circuit breaker pattern
- Retry mechanisms with backoff
- Service-to-service authentication

## 📈 Performance & Reliability

### Circuit Breaker Pattern
- Automatic failure detection
- Configurable failure thresholds
- Automatic recovery mechanisms
- Health monitoring

### Retry Mechanisms
- Exponential backoff with jitter
- Configurable retry policies
- Smart retry logic for retryable errors
- Timeout handling

### Rate Limiting
- Token bucket algorithm
- Per-IP rate limiting
- Configurable limits
- Automatic cleanup

## 🚀 Production Deployment

### Docker Deployment

1. **Build Images**
```bash
docker-compose build
```

2. **Deploy Services**
```bash
docker-compose up -d
```

3. **Scale Services**
```bash
docker-compose up -d --scale auth-service=3 --scale course-service=2
```

### Environment Configuration

For production, update environment variables in `docker-compose.yml`:
- Use strong passwords and secrets
- Configure proper database URLs
- Set up SSL/TLS certificates
- Configure monitoring and logging

### Health Checks

All services include health checks:
- HTTP health endpoints
- Database connectivity checks
- Service dependency validation

## 📚 API Documentation

### Interactive Documentation
- **Swagger UI**: http://localhost:8080/api/v1/docs
- **OpenAPI Specification**: http://localhost:8080/api/v1/docs/openapi.json

### Authentication
Most endpoints require JWT authentication. Include the token in the Authorization header:
```
Authorization: Bearer YOUR_JWT_TOKEN
```

### Response Format
All API responses follow this format:
```json
{
  "success": true,
  "message": "Operation successful",
  "data": {
    // Response data
  }
}
```

### Error Handling
Error responses include:
```json
{
  "success": false,
  "message": "Error description",
  "error": "detailed_error_code"
}
```

## 🗂️ Project Structure

```
Study-Platform/
├── auth-service/              # Authentication service
│   ├── cmd/main.go           # Service entry point
│   ├── internal/             # Internal packages
│   │   ├── handler/          # gRPC handlers
│   │   ├── model/            # Data models
│   │   ├── repository/       # Data access layer
│   │   └── service/          # Business logic
│   └── proto/                # Protobuf definitions
├── course-service/            # Course management service
├── progress-service/          # Progress tracking service
├── api-gateway/              # HTTP API Gateway
│   ├── cmd/main.go           # Gateway entry point
│   ├── internal/             # Internal packages
│   │   ├── handler/          # HTTP handlers
│   │   ├── middleware/       # HTTP middleware
│   │   └── router/           # Route definitions
├── pkg/                      # Shared packages
│   ├── database/             # Database utilities
│   └── logger/               # Logging utilities
├── migrations/               # Database migrations
├── docker-compose.yml        # Docker services
└── README.md                # This file
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License.

## 🆘 Troubleshooting

### Common Issues

**Services won't start:**
```bash
# Check logs
docker-compose logs

# Restart services
docker-compose restart
```

**Database connection issues:**
```bash
# Check database health
docker-compose exec postgres pg_isready -U admin -d studyplatform
```

**Port conflicts:**
```bash
# Check port usage
lsof -i :8080
```

**Memory issues:**
```bash
# Check Docker resources
docker system df
docker system prune
```

### Getting Help

- Check the logs: `docker-compose logs [service-name]`
- Verify service health: `curl http://localhost:8080/api/v1/health`
- Review the API documentation: http://localhost:8080/api/v1/docs
- Check the implementation checklist: `specs/implementation-checklist.md`

---

**Happy Learning! 🎓**