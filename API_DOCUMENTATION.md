# Study Platform API Documentation

## Base URL
```
http://localhost:8080/api/v1
```

## Authentication

Most endpoints require JWT authentication. Include the token in the Authorization header:
```
Authorization: Bearer YOUR_JWT_TOKEN
```

## Response Format

### Success Response
```json
{
  "success": true,
  "message": "Operation successful",
  "data": {
    // Response data
  }
}
```

### Error Response
```json
{
  "success": false,
  "message": "Error description",
  "error": "detailed_error_code"
}
```

## Endpoints

### Health & Monitoring

#### GET /health
Check API Gateway health status.

**Response:**
```json
{
  "status": "healthy",
  "service": "api-gateway",
  "version": "1.0.0"
}
```

#### GET /health/circuit-breakers
Check circuit breaker status for all services.

**Response:**
```json
{
  "status": "healthy",
  "circuit_breakers": {
    "auth-service": {"state": "closed"},
    "course-service": {"state": "closed"},
    "progress-service": {"state": "closed"}
  }
}
```

### Authentication

#### POST /auth/register
Register a new user account.

**Request Body:**
```json
{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "username": "johndoe",
      "email": "john@example.com",
      "role": "student"
    }
  }
}
```

#### POST /auth/login
Authenticate user and return JWT token.

**Request Body:**
```json
{
  "email": "john@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "username": "johndoe",
      "email": "john@example.com",
      "role": "student"
    }
  }
}
```

#### POST /auth/validate
Validate JWT token.

**Request Body:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response:**
```json
{
  "success": true,
  "message": "Token is valid",
  "data": {
    "user_id": "123e4567-e89b-12d3-a456-426614174000",
    "role": "student"
  }
}
```

#### GET /auth/profile
Get user profile information. **Requires authentication.**

**Response:**
```json
{
  "success": true,
  "message": "Profile retrieved successfully",
  "data": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "roles": ["student"]
  }
}
```

#### GET /auth/oauth/{provider}/url
Get OAuth authorization URL for specified provider.

**Parameters:**
- `provider`: `google`, `github`, or `facebook`

**Response:**
```json
{
  "success": true,
  "message": "OAuth URL generated successfully",
  "data": {
    "url": "https://accounts.google.com/oauth/authorize?..."
  }
}
```

#### GET /auth/oauth/{provider}/callback
Handle OAuth callback from provider.

**Parameters:**
- `provider`: `google`, `github`, or `facebook`
- `code`: Authorization code from OAuth provider
- `state`: State parameter for security

### Courses

#### GET /courses
List all courses with pagination.

**Query Parameters:**
- `page`: Page number (default: 1)
- `page_size`: Items per page (default: 10, max: 100)
- `category`: Filter by category
- `level`: Filter by level (`beginner`, `intermediate`, `advanced`)
- `status`: Filter by status (`draft`, `published`, `archived`)

**Response:**
```json
{
  "success": true,
  "message": "Courses retrieved successfully",
  "data": {
    "courses": [
      {
        "id": "123e4567-e89b-12d3-a456-426614174000",
        "title": "Introduction to Go Programming",
        "description": "Learn the fundamentals of Go programming language",
        "instructor_id": "123e4567-e89b-12d3-a456-426614174001",
        "instructor_name": "John Instructor",
        "category": "Programming",
        "level": "beginner",
        "price": 99.99,
        "currency": "USD",
        "status": "published",
        "rating": 4.5,
        "rating_count": 120,
        "enrollment_count": 500,
        "duration_minutes": 480,
        "tags": ["go", "programming", "backend"],
        "thumbnail_url": "https://example.com/thumbnail.jpg",
        "created_at": "2023-01-01T00:00:00Z",
        "updated_at": "2023-01-01T00:00:00Z"
      }
    ],
    "total": 50,
    "page": 1,
    "page_size": 10,
    "total_pages": 5
  }
}
```

#### GET /courses/search
Search courses with advanced filters.

**Query Parameters:**
- `q`: Search query
- `category`: Filter by category
- `level`: Filter by level
- `min_price`: Minimum price
- `max_price`: Maximum price
- `min_rating`: Minimum rating
- `tags`: Comma-separated tags

**Response:** Same as GET /courses

#### GET /courses/{id}
Get detailed information about a specific course.

**Response:**
```json
{
  "success": true,
  "message": "Course retrieved successfully",
  "data": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "title": "Introduction to Go Programming",
    "description": "Learn the fundamentals of Go programming language",
    "instructor_id": "123e4567-e89b-12d3-a456-426614174001",
    "instructor_name": "John Instructor",
    "category": "Programming",
    "level": "beginner",
    "price": 99.99,
    "currency": "USD",
    "status": "published",
    "rating": 4.5,
    "rating_count": 120,
    "enrollment_count": 500,
    "duration_minutes": 480,
    "tags": ["go", "programming", "backend"],
    "thumbnail_url": "https://example.com/thumbnail.jpg",
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z"
  }
}
```

#### POST /courses
Create a new course. **Requires authentication.**

**Request Body:**
```json
{
  "title": "Introduction to Go Programming",
  "description": "Learn the fundamentals of Go programming language",
  "category": "Programming",
  "level": "beginner",
  "price": 99.99,
  "currency": "USD",
  "tags": ["go", "programming", "backend"],
  "thumbnail_url": "https://example.com/thumbnail.jpg"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Course created successfully",
  "data": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "title": "Introduction to Go Programming",
    "description": "Learn the fundamentals of Go programming language",
    "instructor_id": "123e4567-e89b-12d3-a456-426614174001",
    "category": "Programming",
    "level": "beginner",
    "price": 99.99,
    "currency": "USD",
    "status": "draft",
    "tags": ["go", "programming", "backend"],
    "thumbnail_url": "https://example.com/thumbnail.jpg",
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z"
  }
}
```

#### PUT /courses/{id}
Update an existing course. **Requires authentication.**

**Request Body:** Same as POST /courses

#### DELETE /courses/{id}
Delete a course. **Requires authentication.**

**Response:**
```json
{
  "success": true,
  "message": "Course deleted successfully"
}
```

#### POST /courses/{course_id}/enroll
Enroll in a course. **Requires authentication.**

**Response:**
```json
{
  "success": true,
  "message": "Enrollment created successfully",
  "data": {
    "enrollment_id": "123e4567-e89b-12d3-a456-426614174000",
    "course_id": "123e4567-e89b-12d3-a456-426614174001",
    "user_id": "123e4567-e89b-12d3-a456-426614174002",
    "enrolled_at": "2023-01-01T00:00:00Z",
    "status": "enrolled"
  }
}
```

### Lectures

#### GET /courses/{course_id}/lectures
List all lectures for a course.

**Response:**
```json
{
  "success": true,
  "message": "Lectures retrieved successfully",
  "data": {
    "lectures": [
      {
        "id": "123e4567-e89b-12d3-a456-426614174000",
        "title": "Introduction to Variables",
        "description": "Learn about variables in Go",
        "course_id": "123e4567-e89b-12d3-a456-426614174001",
        "order": 1,
        "duration_minutes": 15,
        "video_url": "https://example.com/video.mp4",
        "content": "Lecture content...",
        "is_preview": true,
        "created_at": "2023-01-01T00:00:00Z",
        "updated_at": "2023-01-01T00:00:00Z"
      }
    ]
  }
}
```

#### GET /lectures/{id}
Get detailed information about a specific lecture.

#### POST /lectures
Create a new lecture. **Requires authentication.**

**Request Body:**
```json
{
  "title": "Introduction to Variables",
  "description": "Learn about variables in Go",
  "course_id": "123e4567-e89b-12d3-a456-426614174001",
  "order": 1,
  "duration_minutes": 15,
  "video_url": "https://example.com/video.mp4",
  "content": "Lecture content...",
  "is_preview": true
}
```

### Enrollments

#### GET /enrollments
List user's enrollments. **Requires authentication.**

**Query Parameters:**
- `page`: Page number (default: 1)
- `page_size`: Items per page (default: 10)
- `status`: Filter by status (`enrolled`, `completed`, `dropped`)

**Response:**
```json
{
  "success": true,
  "message": "Enrollments retrieved successfully",
  "data": {
    "enrollments": [
      {
        "id": "123e4567-e89b-12d3-a456-426614174000",
        "course_id": "123e4567-e89b-12d3-a456-426614174001",
        "user_id": "123e4567-e89b-12d3-a456-426614174002",
        "enrolled_at": "2023-01-01T00:00:00Z",
        "status": "enrolled",
        "progress_percentage": 75.5,
        "completed_lectures": 12,
        "total_lectures": 16,
        "total_watch_time_seconds": 3600,
        "created_at": "2023-01-01T00:00:00Z",
        "updated_at": "2023-01-01T00:00:00Z"
      }
    ],
    "total": 5,
    "page": 1,
    "page_size": 10,
    "total_pages": 1
  }
}
```

#### POST /enrollments
Create a new enrollment. **Requires authentication.**

**Request Body:**
```json
{
  "course_id": "123e4567-e89b-12d3-a456-426614174001"
}
```

#### GET /enrollments/courses/{course_id}
Get enrollment details for a specific course. **Requires authentication.**

### Progress Tracking

#### POST /progress/update
Update learning progress. **Requires authentication.**

**Request Body:**
```json
{
  "course_id": "123e4567-e89b-12d3-a456-426614174001",
  "lecture_id": "123e4567-e89b-12d3-a456-426614174002",
  "progress_percentage": 85.5,
  "watch_time_seconds": 1200,
  "completed": false
}
```

#### GET /progress/courses/{course_id}/lectures/{lecture_id}
Get progress for a specific lecture. **Requires authentication.**

**Response:**
```json
{
  "success": true,
  "message": "Progress retrieved successfully",
  "data": {
    "course_id": "123e4567-e89b-12d3-a456-426614174001",
    "lecture_id": "123e4567-e89b-12d3-a456-426614174002",
    "user_id": "123e4567-e89b-12d3-a456-426614174003",
    "progress_percentage": 85.5,
    "watch_time_seconds": 1200,
    "completed": false,
    "last_accessed": "2023-01-01T00:00:00Z"
  }
}
```

#### GET /progress/courses/{course_id}
Get overall progress for a course. **Requires authentication.**

#### GET /progress/lectures/{course_id}
Get progress for all lectures in a course. **Requires authentication.**

#### GET /progress/courses/{course_id}/completion
Get course completion status. **Requires authentication.**

#### POST /progress/lectures/complete
Mark a lecture as complete. **Requires authentication.**

**Request Body:**
```json
{
  "course_id": "123e4567-e89b-12d3-a456-426614174001",
  "lecture_id": "123e4567-e89b-12d3-a456-426614174002"
}
```

### Analytics

#### GET /analytics/user
Get user analytics. **Requires authentication.**

**Response:**
```json
{
  "success": true,
  "message": "Analytics retrieved successfully",
  "data": {
    "total_courses": 5,
    "completed_courses": 2,
    "total_watch_time_hours": 24.5,
    "average_progress": 68.2,
    "learning_streak_days": 7,
    "favorite_categories": ["Programming", "Data Science"]
  }
}
```

## Error Codes

| HTTP Status | Error Code | Description |
|-------------|------------|-------------|
| 400 | `invalid_request` | Invalid request format or parameters |
| 401 | `unauthorized` | Missing or invalid authentication |
| 403 | `forbidden` | Insufficient permissions |
| 404 | `not_found` | Resource not found |
| 409 | `conflict` | Resource already exists |
| 422 | `validation_error` | Input validation failed |
| 429 | `rate_limit_exceeded` | Too many requests |
| 500 | `internal_error` | Internal server error |

## Rate Limiting

The API implements rate limiting to prevent abuse:
- **Limit**: 100 requests per second per IP address
- **Burst**: 200 requests maximum burst capacity
- **Headers**: Rate limit information is included in response headers:
  - `X-RateLimit-Limit`: Maximum requests per window
  - `X-RateLimit-Remaining`: Remaining requests in current window
  - `X-RateLimit-Reset`: Time when rate limit resets

## Circuit Breaker

The API Gateway implements circuit breaker pattern for service resilience:
- **States**: `closed` (healthy), `open` (failing), `half-open` (testing)
- **Thresholds**: 5 failures to open circuit
- **Timeout**: 60 seconds before attempting recovery
- **Monitoring**: Check status at `/health/circuit-breakers`

## Pagination

List endpoints support pagination:
- `page`: Page number (starts from 1)
- `page_size`: Items per page (default: 10, max: 100)
- Response includes:
  - `total`: Total number of items
  - `page`: Current page number
  - `page_size`: Items per page
  - `total_pages`: Total number of pages

## Filtering and Searching

Many endpoints support filtering and searching:
- Use query parameters for filtering
- Multiple filters can be combined
- Search supports partial matches
- Use comma-separated values for multiple values

## Authentication Flow

1. **Register**: `POST /auth/register`
2. **Login**: `POST /auth/login` (returns JWT token)
3. **Use Token**: Include `Authorization: Bearer TOKEN` header
4. **Validate**: `POST /auth/validate` (optional, for token verification)

## OAuth Flow

1. **Get URL**: `GET /auth/oauth/{provider}/url`
2. **Redirect**: User authorizes with OAuth provider
3. **Callback**: `GET /auth/oauth/{provider}/callback`
4. **Login**: Returns JWT token for authenticated user

## Testing Examples

### Using curl

```bash
# Register new user
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"password123"}'

# Login
TOKEN=$(curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}' \
  | jq -r '.data.token')

# List courses
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/courses

# Create course
curl -X POST http://localhost:8080/api/v1/courses \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"title":"Test Course","description":"Test description","category":"Programming","level":"beginner","price":99.99}'
```

### Using Postman

1. Import the OpenAPI specification from `/docs/openapi.json`
2. Set up environment variables for base URL and tokens
3. Use the collection to test all endpoints

## WebSocket Support

Currently not implemented. Future versions will include:
- Real-time progress updates
- Live chat during lectures
- Notifications for course updates