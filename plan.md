# Online Learning Platform - Implementation Plan

## 1. Project Overview

The Online Learning Platform is a microservices-based application providing courses, video lectures, progress tracking, and an AI-powered chatbot to assist learners. The system now also includes forum discussions for community learning and payment processing for premium courses. The platform is designed to handle high traffic with low latency while maintaining security and scalability.

### 1.1 Architecture Overview

- **Microservices Architecture**: 7 core microservices with gRPC/HTTP communication
- **API Gateway**: Single entry point for clients
- **Frontend**: Web application (separate from this plan)
- **Databases**: PostgreSQL for persistent storage
- **Object Storage**: Bucket service for video files
- **Containerization**: Docker for deployment
- **CI/CD**: GitHub Actions

## 2. Microservices Breakdown

### 2.1 Auth Service

- User registration and authentication
- Role-based access control
- JWT token generation and validation
- Protocol: gRPC
- Database: PostgreSQL (users table with role enum)

### 2.2 Course Service

- Course management (CRUD operations)
- Lecture organization
- Protocol: gRPC
- Database: PostgreSQL (courses, lectures tables)

### 2.3 Progress Service

- User progress tracking
- Course completion status
- Protocol: gRPC
- Database: PostgreSQL (progress, enrollment tables)

### 2.4 Video Service

- Video processing and metadata management
- Video streaming control
- Access management
- Protocol: HTTP
- Database: PostgreSQL (video metadata)
- Integration: Bucket Service for storage

### 2.5 Bucket Service

- Video file storage and retrieval
- Multi-format video processing
- Thumbnail generation
- Secure access control via signed URLs
- Protocol: HTTP
- Storage: Object storage system (MinIO, AWS S3, or similar)
- Integration: CDN for content delivery

### 2.6 Chatbot Service

- AI-powered learning assistant
- Course recommendations
- Chat history storage
- Protocols: HTTP (REST API), WebSocket (real-time chat)
- Database: PostgreSQL (chat_history table)
- Integration: OpenAI API, Gemini API, or Llama API

### 2.7 Forum Service

- Discussion topics and posts
- Q&A functionality
- Protocol: HTTP
- Database: PostgreSQL (forum_topics, forum_posts tables)

### 2.8 Payment Service

- Process course purchases
- Manage payment methods
- Handle subscriptions
- Protocol: HTTP
- Database: PostgreSQL (payment_methods, transactions, subscriptions tables)
- Integration: Payment gateway (e.g., Stripe, PayPal)

## 3. Database Schema

### 3.1 Users Table

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(100) NOT NULL,
    role user_role NOT NULL DEFAULT 'student',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

### 3.2 User Role Enum

```sql
CREATE TYPE user_role AS ENUM ('admin', 'instructor', 'student');
```

### 3.3 Courses Table

```sql
CREATE TABLE courses (
    id UUID PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    description TEXT,
    creator_id UUID REFERENCES users(id),
    thumbnail_url TEXT,
    price DECIMAL(10,2) DEFAULT 0.00,
    is_free BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

### 3.4 Lectures Table

```sql
CREATE TABLE lectures (
    id UUID PRIMARY KEY,
    course_id UUID REFERENCES courses(id),
    title VARCHAR(100) NOT NULL,
    description TEXT,
    video_url TEXT,
    duration INT,
    sequence_order INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

### 3.5 Videos Table

```sql
CREATE TABLE videos (
    id UUID PRIMARY KEY,
    lecture_id UUID REFERENCES lectures(id),
    title VARCHAR(100) NOT NULL,
    description TEXT,
    bucket_path TEXT NOT NULL,
    formats JSONB,
    duration INT,
    thumbnail_url TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

### 3.6 Progress Table

```sql
CREATE TABLE progress (
    user_id UUID REFERENCES users(id),
    lecture_id UUID REFERENCES lectures(id),
    watched_duration INT DEFAULT 0,
    completed BOOLEAN DEFAULT FALSE,
    last_watched_at TIMESTAMP,
    PRIMARY KEY (user_id, lecture_id)
);
```

### 3.7 Chat History Table

```sql
CREATE TABLE chat_history (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    message TEXT NOT NULL,
    is_user BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

### 3.8 Forum Topics Table

```sql
CREATE TABLE forum_topics (
    id UUID PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    course_id UUID REFERENCES courses(id),
    creator_id UUID REFERENCES users(id),
    is_pinned BOOLEAN DEFAULT FALSE,
    is_locked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

### 3.9 Forum Posts Table

```sql
CREATE TABLE forum_posts (
    id UUID PRIMARY KEY,
    topic_id UUID REFERENCES forum_topics(id),
    user_id UUID REFERENCES users(id),
    content TEXT NOT NULL,
    is_solution BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

### 3.10 Payment Methods Table

```sql
CREATE TABLE payment_methods (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    provider VARCHAR(50) NOT NULL,
    token VARCHAR(255) NOT NULL,
    card_last_four VARCHAR(4),
    card_expiry VARCHAR(7),
    is_default BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

### 3.11 Transactions Table

```sql
CREATE TABLE transactions (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    course_id UUID REFERENCES courses(id),
    payment_method_id UUID REFERENCES payment_methods(id),
    amount DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) NOT NULL,
    status VARCHAR(20) NOT NULL,
    transaction_reference VARCHAR(100) UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

### 3.12 Subscriptions Table

```sql
CREATE TABLE subscriptions (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    payment_method_id UUID REFERENCES payment_methods(id),
    plan_name VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL,
    billing_period VARCHAR(20) NOT NULL,
    next_billing_date TIMESTAMP,
    price DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

### 3.13 Enrollment Table

```sql
CREATE TABLE enrollment (
    user_id UUID REFERENCES users(id),
    course_id UUID REFERENCES courses(id),
    transaction_id UUID REFERENCES transactions(id),
    enrolled_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE,
    PRIMARY KEY (user_id, course_id)
);
```

## 4. API Definitions

### 4.1 Auth Service gRPC

```protobuf
service AuthService {
    rpc Register(RegisterRequest) returns (RegisterResponse);
    rpc Login(LoginRequest) returns (LoginResponse);
    rpc ValidateToken(ValidateTokenRequest) returns (ValidateTokenResponse);
    rpc GetUserRoles(GetUserRolesRequest) returns (GetUserRolesResponse);
    rpc AssignRole(AssignRoleRequest) returns (AssignRoleResponse);
    rpc RemoveRole(RemoveRoleRequest) returns (RemoveRoleResponse);
}
```

### 4.2 Course Service gRPC

```protobuf
service CourseService {
    rpc CreateCourse(CreateCourseRequest) returns (CourseResponse);
    rpc GetCourse(GetCourseRequest) returns (CourseResponse);
    rpc ListCourses(ListCoursesRequest) returns (ListCoursesResponse);
    rpc UpdateCourse(UpdateCourseRequest) returns (CourseResponse);
    rpc DeleteCourse(DeleteCourseRequest) returns (DeleteCourseResponse);
    rpc CreateLecture(CreateLectureRequest) returns (LectureResponse);
    rpc GetLecture(GetLectureRequest) returns (LectureResponse);
    rpc ListLectures(ListLecturesRequest) returns (ListLecturesResponse);
    rpc UpdateLecture(UpdateLectureRequest) returns (LectureResponse);
    rpc DeleteLecture(DeleteLectureRequest) returns (DeleteLectureResponse);
}
```

### 4.3 Progress Service gRPC

```protobuf
service ProgressService {
    rpc UpdateProgress(UpdateProgressRequest) returns (UpdateProgressResponse);
    rpc GetUserProgress(GetUserProgressRequest) returns (GetUserProgressResponse);
    rpc GetCourseProgress(GetCourseProgressRequest) returns (GetCourseProgressResponse);
    rpc EnrollUser(EnrollUserRequest) returns (EnrollUserResponse);
    rpc GetUserEnrollments(GetUserEnrollmentsRequest) returns (GetUserEnrollmentsResponse);
    rpc GetCourseEnrollments(GetCourseEnrollmentsRequest) returns (GetCourseEnrollmentsResponse);
}
```

### 4.4 Video Service HTTP

```
GET /api/videos/{video_id}/metadata
GET /api/videos/{video_id}/stream
POST /api/videos/upload
PUT /api/videos/{video_id}
DELETE /api/videos/{video_id}
```

### 4.5 Bucket Service HTTP

```
GET /api/storage/videos/{video_id}
GET /api/storage/videos/{video_id}/formats/{format}
GET /api/storage/videos/{video_id}/thumbnail
POST /api/storage/videos/upload
PUT /api/storage/videos/{video_id}
DELETE /api/storage/videos/{video_id}
```

### 4.6 Chatbot Service HTTP/WebSocket

```
POST /api/chat/message
GET /api/chat/history/{user_id}
WebSocket: /ws/chat/{user_id}
```

### 4.7 Forum Service HTTP

```
GET /api/forums/courses/{course_id}/topics
POST /api/forums/topics
GET /api/forums/topics/{topic_id}
POST /api/forums/topics/{topic_id}/posts
PUT /api/forums/topics/{topic_id}
PUT /api/forums/posts/{post_id}
DELETE /api/forums/posts/{post_id}
PUT /api/forums/posts/{post_id}/solution
```

### 4.8 Payment Service HTTP

```
POST /api/payments/methods
GET /api/payments/methods
DELETE /api/payments/methods/{method_id}
POST /api/payments/purchase/course/{course_id}
POST /api/payments/subscribe
GET /api/payments/transactions
GET /api/payments/subscriptions
PUT /api/payments/subscriptions/{subscription_id}
DELETE /api/payments/subscriptions/{subscription_id}
```

## 5. Implementation Plan

### 5.1 Phase 1: Core Infrastructure (Week 1-2)

- Set up Go project structure for microservices
- Initialize Docker configuration
- Set up PostgreSQL database
- Implement Auth Service with role-based access
- Create CI/CD pipeline with GitHub Actions

### 5.2 Phase 2: Course and Progress Services (Week 3-4)

- Implement Course Service
- Implement Progress Service
- Develop inter-service communication
- Create API Gateway
- Implement enrollment functionality

### 5.3 Phase 3: Video and Storage Services (Week 5-6)

- Implement Bucket Service with object storage
- Set up video processing capabilities
- Implement Video Service with metadata management
- Integrate with CDN for content delivery
- Implement secure access control

### 5.4 Phase 4: Chatbot and Forum Services (Week 7-8)

- Implement Chatbot Service with AI integration
- Implement WebSocket for real-time chat
- Store chat history in PostgreSQL
- Implement Forum Service for course discussions
- Develop topic and post management

### 5.5 Phase 5: Payment Service (Week 9-10)

- Implement Payment Service
- Integrate with payment gateway
- Set up subscription management
- Implement purchase flow
- Secure payment processing

### 5.6 Phase 6: Testing and Deployment (Week 11-12)

- Write unit and integration tests
- Perform load testing
- Optimize for performance
- Complete Docker configuration
- Documentation

## 6. Project Structure

```
.
├── api-gateway/
│   └── main.go
├── auth-service/
│   ├── cmd/
│   │   └── main.go
│   ├── internal/
│   │   ├── handler/
│   │   ├── model/
│   │   ├── repository/
│   │   └── service/
│   └── proto/
├── course-service/
│   ├── cmd/
│   │   └── main.go
│   ├── internal/
│   │   ├── handler/
│   │   ├── model/
│   │   ├── repository/
│   │   └── service/
│   └── proto/
├── progress-service/
│   ├── cmd/
│   │   └── main.go
│   ├── internal/
│   │   ├── handler/
│   │   ├── model/
│   │   ├── repository/
│   │   └── service/
│   └── proto/
├── video-service/
│   ├── cmd/
│   │   └── main.go
│   └── internal/
│       ├── handler/
│       ├── model/
│       ├── repository/
│       └── service/
├── bucket-service/
│   ├── cmd/
│   │   └── main.go
│   └── internal/
│       ├── handler/
│       ├── model/
│       ├── storage/
│       └── service/
├── chatbot-service/
│   ├── cmd/
│   │   └── main.go
│   └── internal/
│       ├── handler/
│       ├── model/
│       ├── repository/
│       └── service/
├── forum-service/
│   ├── cmd/
│   │   └── main.go
│   └── internal/
│       ├── handler/
│       ├── model/
│       ├── repository/
│       └── service/
├── payment-service/
│   ├── cmd/
│   │   └── main.go
│   └── internal/
│       ├── handler/
│       ├── model/
│       ├── repository/
│       └── service/
├── proto/
│   ├── auth.proto
│   ├── course.proto
│   └── progress.proto
├── pkg/
│   ├── logger/
│   ├── database/
│   └── middleware/
├── docker-compose.yml
├── Dockerfile.auth
├── Dockerfile.course
├── Dockerfile.progress
├── Dockerfile.video
├── Dockerfile.bucket
├── Dockerfile.chatbot
├── Dockerfile.forum
├── Dockerfile.payment
└── Dockerfile.gateway
```

## 7. Technical Considerations

### 7.1 Performance Optimization

- Connection pooling for database access
- Caching strategies for frequently accessed data
- Efficient video streaming with adaptive bitrate
- Chunked video transfers for better performance
- Proper indexing of database tables
- CDN for static assets, videos, and course thumbnails
- Parallel video processing for different formats

### 7.2 Security Measures

- JWT token-based authentication
- Role-based access control
- Password hashing with bcrypt
- Rate limiting for API endpoints
- Input validation and sanitization
- PCI compliance for payment processing
- Data encryption for sensitive information
- Signed URLs for secure video access
- Bucket access control policies

### 7.3 Scalability Approach

- Horizontal scaling of microservices
- Database replication for read-heavy workloads
- Load balancing with API Gateway
- Health checks and auto-recovery
- Sharding for high-traffic services
- Stateless design for easier scaling
- Autoscaling bucket storage based on demand

### 7.4 Error Handling

- Consistent error responses across services
- Circuit breaker pattern for inter-service communication
- Logging and monitoring integration
- Retry mechanisms for transient failures
- Graceful degradation of non-critical features
- Dead letter queues for failed processing jobs

## 8. Testing Strategy

### 8.1 Unit Testing

- Test individual components in isolation
- Mock external dependencies
- Aim for high code coverage

### 8.2 Integration Testing

- Test service-to-service communication
- Test database interactions
- Validate API contracts
- Test video upload and streaming flows

### 8.3 Load Testing

- Simulate high traffic scenarios
- Measure response times under load
- Identify bottlenecks
- Test payment processing under concurrent transactions
- Test video streaming with multiple concurrent users

### 8.4 End-to-End Testing

- Test complete user flows
- Validate business requirements
- Ensure proper error handling
- Test forum and payment features
- Test video upload, processing, and playback

## 9. Deployment Strategy

### 9.1 Local Development

- Docker Compose for local environment
- Development database with sample data
- Hot reloading for faster development
- Mock payment gateway for testing
- Local MinIO instance for bucket storage testing

### 9.2 Production Deployment

- Docker containers for all services
- Kubernetes (optional for future scaling)
- Database migration strategy
- Backup and recovery procedures
- Secure payment processing environment
- S3-compatible object storage in production

## 10. Future Enhancements

### 10.1 Technical Improvements

- Implement GraphQL API
- Add Redis for caching
- Implement event-driven architecture
- Add monitoring and alerting
- Enhance payment analytics
- Video content analysis and search
- Automated video captioning

### 10.2 Feature Enhancements

- Social learning features
- Advanced forum features (likes, badges)
- Recommendations engine
- Mobile application support
- Advanced analytics
- Multiple payment options
- Affiliate marketing system
- Live streaming capabilities
- Interactive video features
