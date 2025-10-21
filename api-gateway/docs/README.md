# API Gateway Swagger Documentation

This directory contains auto-generated Swagger/OpenAPI documentation for the Study Platform API Gateway.

## Generated Files

- `docs.go` - Go code for embedding Swagger specs
- `swagger.json` - OpenAPI 3.0 specification in JSON format
- `swagger.yaml` - OpenAPI 3.0 specification in YAML format

## Accessing Documentation

### Production
- **Swagger UI**: https://study.khoipn.id.vn/swagger/index.html
- **OpenAPI JSON**: https://study.khoipn.id.vn/api/v1/swagger/doc.json
- **Legacy Docs**: https://study.khoipn.id.vn/api/v1/docs

### Local Development
- **Swagger UI**: http://localhost:8080/swagger/index.html
- **OpenAPI JSON**: http://localhost:8080/api/v1/swagger/doc.json

## Regenerating Documentation

When you modify API endpoints or add new routes, regenerate the docs:

```bash
# From api-gateway directory
swag init -g cmd/main.go -o docs --parseDependency --parseInternal

# Or from project root
./scripts/generate-swagger.sh
```

## Adding Documentation to Endpoints

Use Swagger annotations in your handler files:

```go
// @Summary      Get user profile
// @Description  Retrieves the authenticated user's profile
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  UserProfileResponse
// @Failure      401  {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /auth/profile [get]
func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
    // Implementation
}
```

## Authentication

All protected endpoints require Bearer token authentication:

```bash
curl -H "Authorization: Bearer YOUR_TOKEN" https://study.khoipn.id.vn/api/v1/auth/profile
```

## Testing with Swagger UI

1. Open Swagger UI in your browser
2. Click "Authorize" button
3. Enter your JWT token: `Bearer YOUR_TOKEN`
4. Test API endpoints directly from the UI

## CI/CD Integration

Swagger docs are automatically regenerated and validated in the CI/CD pipeline before deployment.
