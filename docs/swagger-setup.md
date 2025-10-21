# Swagger/OpenAPI Documentation Setup

## Overview
This document describes how to access and maintain Swagger/OpenAPI documentation for all services.

## Accessing Swagger UI

### Production
- **API Gateway**: https://study.khoipn.id.vn/api/v1/docs
- **OpenAPI Spec**: https://study.khoipn.id.vn/api/v1/docs/openapi.json

### Local Development
- **API Gateway**: http://localhost:8080/api/v1/docs
- **Auth Service**: http://localhost:8081/swagger/index.html
- **Course Service**: http://localhost:8082/swagger/index.html
- **Progress Service**: http://localhost:8083/swagger/index.html
- **Video Service**: http://localhost:8084/swagger/index.html
- **Bucket Service**: http://localhost:8085/swagger/index.html
- **Chatbot Service**: http://localhost:8086/swagger/index.html
- **Forum Service**: http://localhost:8087/swagger/index.html
- **Payment Service**: http://localhost:8088/swagger/index.html
- **Instructor Dashboard**: http://localhost:8089/swagger/index.html

## Regenerating Documentation

When you add or modify API endpoints, regenerate the Swagger docs:

```bash
# From project root
./scripts/generate-swagger.sh

# Or individually for each service:
cd api-gateway && swag init -g cmd/main.go -o docs
cd auth-service && swag init -g cmd/main.go -o docs
cd course-service && swag init -g cmd/main.go -o docs
# ... etc
```

## Adding Swagger Annotations

Example for a new endpoint:

```go
// @Summary      Create a new course
// @Description  Creates a new course with the provided details
// @Tags         courses
// @Accept       json
// @Produce      json
// @Param        course  body      CreateCourseRequest  true  "Course data"
// @Success      201     {object}  Course
// @Failure      400     {object}  ErrorResponse
// @Failure      401     {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /courses [post]
func (h *CourseHandler) CreateCourse(w http.ResponseWriter, r *http.Request) {
    // implementation
}
```

## Security Configuration

All authenticated endpoints use Bearer token authentication:

```yaml
securityDefinitions:
  BearerAuth:
    type: apiKey
    name: Authorization
    in: header
    description: "Bearer {token}"
```

## Best Practices

1. **Keep docs up-to-date**: Run `./scripts/generate-swagger.sh` before committing
2. **Use tags**: Group related endpoints using @Tags
3. **Document errors**: Include all possible error responses
4. **Add examples**: Use @Param and @Success with examples
5. **Version APIs**: Use semantic versioning in @version

## CI/CD Integration

Swagger docs are automatically generated and validated in the CI/CD pipeline:
- Pre-build: Verify docs are up-to-date
- Post-build: Regenerate if needed
- Deployment: Serve via API Gateway
