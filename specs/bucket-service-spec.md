# Bucket Service Specification

## Overview
The Bucket Service is responsible for file storage management using AWS S3 (or S3-compatible storage). It handles file uploads, downloads, and metadata management for course materials, user avatars, and video thumbnails.

## Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   API Gateway   │───▶│  Bucket Service │───▶│  AWS S3 / MinIO │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              │
                              ▼
                       ┌─────────────────┐
                       │   PostgreSQL    │
                       │   (Metadata)    │
                       └─────────────────┘
```

## Technology Stack
- **Language**: Go 1.21+
- **Storage**: AWS S3 / MinIO (S3-compatible)
- **Database**: PostgreSQL (file metadata)
- **Authentication**: JWT tokens from Auth Service
- **Image Processing**: Go imaging libraries
- **Container**: Docker

## Features

### Core Features
1. **File Upload**
   - Multiple file format support
   - File size validation
   - MIME type validation
   - Progress tracking
   - Resumable uploads

2. **File Management**
   - File metadata storage
   - File versioning
   - Soft deletion
   - File organization by buckets

3. **Security**
   - Pre-signed URLs
   - Access control
   - File encryption at rest
   - Secure file serving

4. **Image Processing**
   - Thumbnail generation
   - Image resizing
   - Format conversion
   - Optimization

## Database Schema

```sql
-- File metadata table
CREATE TABLE files (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    filename VARCHAR(255) NOT NULL,
    original_filename VARCHAR(255) NOT NULL,
    content_type VARCHAR(100) NOT NULL,
    size_bytes BIGINT NOT NULL,
    bucket_name VARCHAR(100) NOT NULL,
    object_key VARCHAR(500) NOT NULL,
    upload_user_id UUID NOT NULL,
    is_public BOOLEAN DEFAULT FALSE,
    metadata JSONB,
    checksum VARCHAR(64),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP NULL
);

-- File access permissions
CREATE TABLE file_permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    file_id UUID REFERENCES files(id) ON DELETE CASCADE,
    user_id UUID,
    permission_type VARCHAR(20) NOT NULL, -- 'read', 'write', 'delete'
    granted_by UUID NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Upload sessions for resumable uploads
CREATE TABLE upload_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    upload_id VARCHAR(255) NOT NULL, -- S3 multipart upload ID
    filename VARCHAR(255) NOT NULL,
    content_type VARCHAR(100) NOT NULL,
    total_size BIGINT NOT NULL,
    uploaded_size BIGINT DEFAULT 0,
    bucket_name VARCHAR(100) NOT NULL,
    object_key VARCHAR(500) NOT NULL,
    user_id UUID NOT NULL,
    status VARCHAR(20) DEFAULT 'active', -- 'active', 'completed', 'aborted'
    created_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP NOT NULL
);

-- Indexes
CREATE INDEX idx_files_user_id ON files(upload_user_id);
CREATE INDEX idx_files_bucket_key ON files(bucket_name, object_key);
CREATE INDEX idx_files_created_at ON files(created_at);
CREATE INDEX idx_file_permissions_file_user ON file_permissions(file_id, user_id);
CREATE INDEX idx_upload_sessions_user_status ON upload_sessions(user_id, status);
```

## API Endpoints

### File Upload
```http
POST /api/files/upload
Content-Type: multipart/form-data

Form fields:
- file: binary file data
- bucket: bucket name (optional, defaults to 'general')
- is_public: boolean (optional, defaults to false)
- metadata: JSON object (optional)

Response:
{
    "file_id": "uuid",
    "filename": "string",
    "size": 1024,
    "content_type": "image/jpeg",
    "url": "https://bucket.s3.amazonaws.com/path/to/file",
    "thumbnail_url": "https://bucket.s3.amazonaws.com/path/to/thumbnail" // for images
}
```

### Resumable Upload Start
```http
POST /api/files/upload/start
Content-Type: application/json

{
    "filename": "video.mp4",
    "content_type": "video/mp4",
    "size": 1073741824,
    "bucket": "videos"
}

Response:
{
    "session_id": "uuid",
    "upload_id": "s3-multipart-upload-id",
    "presigned_urls": [
        {
            "part_number": 1,
            "url": "https://presigned-url-1"
        }
    ]
}
```

### File Download
```http
GET /api/files/{file_id}

Response:
- Redirect to S3 signed URL for private files
- Direct S3 URL for public files
```

### File Metadata
```http
GET /api/files/{file_id}/metadata

Response:
{
    "id": "uuid",
    "filename": "document.pdf",
    "original_filename": "My Document.pdf",
    "content_type": "application/pdf",
    "size": 1024000,
    "is_public": false,
    "created_at": "2024-01-01T12:00:00Z",
    "metadata": {}
}
```

### List User Files
```http
GET /api/files?bucket=videos&page=1&limit=20&sort=created_at&order=desc

Response:
{
    "files": [...],
    "pagination": {
        "page": 1,
        "limit": 20,
        "total": 100,
        "total_pages": 5
    }
}
```

## Configuration

### Environment Variables
```env
# Service Configuration
BUCKET_SERVICE_PORT=8084
BUCKET_SERVICE_HOST=0.0.0.0

# Database Configuration
DB_HOST=postgres
DB_PORT=5432
DB_USER=bucket_user
DB_PASSWORD=bucket_password
DB_NAME=bucket_service

# S3 Configuration
S3_ENDPOINT=https://s3.amazonaws.com
S3_REGION=us-east-1
S3_ACCESS_KEY_ID=your-access-key
S3_SECRET_ACCESS_KEY=your-secret-key
S3_BUCKET_PREFIX=study-platform

# MinIO Configuration (alternative)
MINIO_ENDPOINT=http://minio:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin

# File Upload Limits
MAX_FILE_SIZE=5368709120  # 5GB
MAX_UPLOAD_PARTS=1000
CHUNK_SIZE=5242880        # 5MB

# Security
JWT_SECRET=your-jwt-secret
CORS_ORIGINS=http://localhost:3000,https://yourdomain.com

# Image Processing
THUMBNAIL_SIZES=150x150,300x300,600x600
IMAGE_QUALITY=85
```

## Service Implementation

### Directory Structure
```
bucket-service/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── handler/
│   │   ├── files.go
│   │   ├── upload.go
│   │   └── health.go
│   ├── service/
│   │   ├── bucket.go
│   │   ├── s3.go
│   │   └── image.go
│   ├── repository/
│   │   └── file_repository.go
│   ├── middleware/
│   │   ├── auth.go
│   │   ├── cors.go
│   │   └── logging.go
│   └── model/
│       └── file.go
├── migrations/
│   ├── 001_create_files_table.up.sql
│   ├── 002_create_file_permissions_table.up.sql
│   └── 003_create_upload_sessions_table.up.sql
├── Dockerfile
├── docker-compose.yml
└── go.mod
```

## Implementation Steps

### Phase 1: Basic Setup
1. **Project Structure**
   ```bash
   mkdir -p bucket-service/{cmd/server,internal/{config,handler,service,repository,middleware,model},migrations}
   cd bucket-service && go mod init bucket-service
   ```

2. **Dependencies**
   ```bash
   go get github.com/aws/aws-sdk-go-v2/service/s3
   go get github.com/aws/aws-sdk-go-v2/config
   go get github.com/gin-gonic/gin
   go get github.com/lib/pq
   go get gorm.io/gorm
   go get github.com/golang-migrate/migrate/v4
   go get github.com/disintegration/imaging
   ```

3. **Database Migration Setup**
   - Create migration files for file metadata
   - Set up GORM models
   - Configure database connection

### Phase 2: S3 Integration
1. **S3 Client Setup**
   - Configure AWS SDK
   - Implement bucket operations
   - Set up credential management

2. **File Operations**
   - Upload single files
   - Multipart upload for large files
   - Generate signed URLs
   - File deletion

### Phase 3: Image Processing
1. **Thumbnail Generation**
   - Resize images on upload
   - Multiple thumbnail sizes
   - Format optimization

2. **Image Optimization**
   - Compress images
   - Convert formats
   - Progressive JPEG support

### Phase 4: Security & Permissions
1. **Access Control**
   - File permission system
   - User-based access
   - Public/private file handling

2. **Security Features**
   - File type validation
   - Virus scanning integration
   - Rate limiting

### Phase 5: Advanced Features
1. **Resumable Uploads**
   - Multipart upload sessions
   - Progress tracking
   - Resume capability

2. **CDN Integration**
   - CloudFront integration
   - Cache invalidation
   - Bandwidth optimization

## Docker Configuration

### Dockerfile
```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o bucket-service cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

COPY --from=builder /app/bucket-service .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8084
CMD ["./bucket-service"]
```

### Docker Compose Service
```yaml
bucket-service:
  build:
    context: ./bucket-service
    dockerfile: Dockerfile
  ports:
    - "8084:8084"
  environment:
    - DB_HOST=postgres
    - DB_PORT=5432
    - DB_USER=bucket_user
    - DB_PASSWORD=bucket_password
    - DB_NAME=bucket_service
    - S3_ENDPOINT=https://s3.amazonaws.com
    - S3_REGION=us-east-1
    - S3_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID}
    - S3_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY}
  depends_on:
    - postgres
  volumes:
    - ./bucket-service/migrations:/app/migrations
  networks:
    - study-platform
```

## Testing Strategy

### Unit Tests
- S3 client operations
- Image processing functions
- File validation logic
- Permission checking

### Integration Tests
- End-to-end upload/download
- Database operations
- S3 integration
- Authentication flow

### Load Testing
- Concurrent upload testing
- Large file handling
- Memory usage optimization
- Connection pooling

## Monitoring & Logging

### Metrics to Track
- Upload/download rates
- File storage usage
- Error rates
- Response times
- Active upload sessions

### Health Checks
- Database connectivity
- S3 connectivity
- Disk space usage
- Memory usage

## Security Considerations

### File Security
- Virus scanning
- File type validation
- Content scanning
- Malware detection

### Access Security
- Signed URLs with expiration
- IP-based restrictions
- Rate limiting per user
- Audit logging

## Deployment Checklist

- [ ] Configure S3 bucket and IAM roles
- [ ] Set up database with proper indexes
- [ ] Configure environment variables
- [ ] Set up monitoring and alerting
- [ ] Configure backup strategy
- [ ] Test disaster recovery
- [ ] Performance benchmarking
- [ ] Security audit