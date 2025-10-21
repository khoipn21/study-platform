# Swagger Documentation Status

## Current Status: 🟡 In Progress

### ✅ Documented Endpoints (5 total)

#### Authentication (4 endpoints)
- ✅ POST `/auth/register` - Register new user
- ✅ POST `/auth/login` - User login  
- ✅ POST `/auth/validate` - Validate JWT token
- ✅ GET `/auth/profile` - Get user profile

#### Courses (1 endpoint)
- ✅ GET `/courses` - List all courses

### 🔄 Next Priority Endpoints

#### Courses (High Priority)
- 🔲 GET `/courses/{id}` - Get course details
- 🔲 POST `/courses` - Create course
- 🔲 PUT `/courses/{id}` - Update course
- 🔲 DELETE `/courses/{id}` - Delete course
- 🔲 GET `/courses/search` - Search courses
- 🔲 GET `/courses/{course_id}/lectures` - List lectures
- 🔲 POST `/courses/{course_id}/lectures` - Create lecture
- 🔲 POST `/courses/{course_id}/enroll` - Enroll in course

#### Progress & Enrollments
- 🔲 POST `/progress/update` - Update progress
- 🔲 GET `/progress/courses/{course_id}` - Get course progress
- 🔲 POST `/progress/lectures/complete` - Mark lecture complete
- 🔲 GET `/enrollments` - List user enrollments
- 🔲 POST `/enrollments` - Create enrollment

#### Payments
- 🔲 POST `/payments/stripe/payment-intents` - Create payment
- 🔲 GET `/payments/transactions` - List transactions
- 🔲 POST `/payments/purchase/course/{course_id}` - Purchase course

#### Videos
- 🔲 POST `/videos/upload-url` - Get upload URL
- 🔲 POST `/videos/upload` - Upload video
- 🔲 GET `/videos/{video_id}` - Get video details
- 🔲 DELETE `/videos/{video_id}/delete` - Delete video

#### Files
- 🔲 POST `/files/upload` - Upload file
- 🔲 GET `/files` - List files
- 🔲 GET `/files/{fileId}` - Download file
- 🔲 DELETE `/files/{fileId}` - Delete file

#### Forum
- 🔲 GET `/forum/topics` - List topics
- 🔲 POST `/forum/topics` - Create topic
- 🔲 GET `/forum/topics/{topicId}` - Get topic
- 🔲 POST `/forum/posts` - Create post
- 🔲 GET `/forum/topics/{topicId}/posts` - List posts

#### Notes
- 🔲 POST `/notes/courses/{course_id}/lectures/{lecture_id}` - Create note
- 🔲 GET `/notes/courses/{course_id}/lectures/{lecture_id}` - Get lecture notes
- 🔲 PUT `/notes/{note_id}` - Update note
- 🔲 DELETE `/notes/{note_id}` - Delete note

### 📊 Progress

```
Total Endpoints: ~80
Documented: 5 (6%)
In Progress: 0
Pending: 75
```

### 🎯 Goals

**Phase 1 (Current):**
- ✅ Authentication endpoints
- ✅ Basic course listing
- 🔄 Core course operations

**Phase 2:**
- Course management (CRUD)
- Enrollment & Progress tracking
- Payment integration

**Phase 3:**
- Video management
- File operations
- Forum & community

**Phase 4:**
- Notes & annotations
- Chatbot integration
- Instructor dashboard
- Student dashboard

### 🔧 How to Add Annotations

1. **Find the handler method** in `api-gateway/internal/handler/`
2. **Add godoc comments** before the function:

```go
// MethodName godoc
// @Summary      Short description
// @Description  Longer description
// @Tags         Category
// @Accept       json
// @Produce      json
// @Param        name type dataType required "description"
// @Success      200 {object} APIResponse "Success message"
// @Failure      400 {object} APIResponse "Error message"
// @Security     BearerAuth
// @Router       /path [method]
func (h *Handler) MethodName(w http.ResponseWriter, r *http.Request) {
    // ...
}
```

3. **Regenerate docs:**
```bash
cd api-gateway
swag init -g cmd/main.go -o docs --parseDependency --parseInternal
```

4. **Test locally:**
```bash
# Start API gateway
go run cmd/main.go

# Visit Swagger UI
open http://localhost:8080/swagger/index.html
```

### 📚 Resources

- Swagger UI: https://study.khoipn.id.vn/swagger/index.html
- OpenAPI Spec: https://study.khoipn.id.vn/api/v1/swagger/doc.json
- Swaggo Docs: https://github.com/swaggo/swag
- OpenAPI 3.0 Spec: https://swagger.io/specification/

### 🤝 Contributing

To help complete the Swagger documentation:

1. Pick an endpoint from "Next Priority" list
2. Add annotations following the template
3. Regenerate docs with `swag init`
4. Test in Swagger UI
5. Commit and push

---

**Last Updated:** 2025-10-21  
**Maintainer:** Development Team
