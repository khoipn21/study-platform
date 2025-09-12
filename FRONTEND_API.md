# Study Platform — Frontend API Guide (via API Gateway)

Base URL

- Development: `http://localhost:8080/api/v1`

Auth & Headers

- Authentication: Bearer JWT. Add `Authorization: Bearer <token>` on protected routes.
- Some proxied services (chatbot, forum, payment, bucket, video) also accept `X-User-ID`/`X-User-Role` when requests are forwarded internally by the gateway.
- Content-Type: JSON unless noted (file upload endpoints use `multipart/form-data`).

Response Shapes

- Gateway-handled routes (Auth, Courses, Progress, Enrollments, Analytics) wrap responses as:
  - `{ success: boolean, message: string, data?: any, error?: string }`
- Proxied routes (Video, Files, Chatbot, Forum, Payment) return the upstream service JSON as-is.

Rate Limits & CORS

- Global rate limiting is enabled in the gateway (burst-friendly). CORS is enabled in dev to allow cross-origin requests.

Health & Docs

- `GET /health` → `{ status: "healthy", service: "api-gateway", version: "1.0.0" }`
- `GET /health/circuit-breakers` → current circuit breaker states.
- `GET /docs` → Swagger UI. `GET /docs/openapi.json` → OpenAPI spec (partial, hand-written).

Auth

- `POST /auth/register` (public)
  - Body: `{ username: string, email: string, password: string, role?: "student"|"instructor"|"admin" }` (role defaults to `student`)
  - Data: `{ user: { id, username, email, role }, token }`
- `POST /auth/login` (public)
  - Body: `{ email: string, password: string }`
  - Data: `{ user: { id, username, email, role }, token }`
- `POST /auth/validate` (public)
  - Body: `{ token: string }`
  - Data: `{ valid: boolean, user?: { id, username, email, role } }`
- `GET /auth/oauth/{provider}/url` (public) → Data: `{ url: string }`
- `GET /auth/oauth/{provider}/callback?code=...&state=...` (public)
  - Data: `{ user: {...}, token }`
- `GET /auth/profile` (auth)
  - Data: `{ id: string, roles: string[] }`

Courses & Lectures

- `GET /courses` (public)
  - Query: `page?` (int, default 1), `page_size?` (int, default 10), `category?`, `level?`, `status?`, `instructor_id?`
  - Data: `{ courses: Course[], total, page, page_size, total_pages }`
- `GET /courses/search` (public)
  - Query: `q`, `page?`, `page_size?`, `category?`, `level?`, `min_price?`, `max_price?`, `min_rating?`
  - Data: same shape as list
- `GET /courses/{id}` (public)
  - Data: `Course`
- `GET /courses/{course_id}/lectures` (public)
  - Query: `page?`, `page_size?`, `status?`
  - Data: `{ lectures: Lecture[], total, page, page_size, total_pages }`
- `GET /courses/lectures/{id}` (public)
  - Data: `Lecture`
- `POST /courses` (auth)
  - Body: `{ title, description, instructor_id, category, level, price, currency, thumbnail_url?, tags?[] }`
  - Data: `Course`
- `PUT /courses/{id}` (auth)
  - Body: `{ title?, description?, category?, level?, price?, currency?, thumbnail_url?, status?, tags?[] }`
  - Data: `Course`
- `DELETE /courses/{id}` (auth)
  - Data: `null`
- `POST /courses/lectures` (auth)
  - Body: `{ course_id, title, description?, order_number?, duration_minutes?, video_url?, video_id?, is_free? }`
  - Data: `Lecture`
- `POST /courses/{course_id}/enroll` (auth)
  - Data: `Enrollment`

Progress

- `POST /progress/update` (auth)
  - Body: `{ course_id, lecture_id, progress_percentage, watch_time_seconds, is_completed }`
  - Data: `UserProgress`
- `GET /progress/courses/{course_id}/lectures/{lecture_id}` (auth) → Data: `UserProgress`
- `GET /progress/courses/{course_id}` (auth)
  - Data: `{ progress: UserProgress[], overall_progress_percentage }`
- `GET /progress/lectures/{course_id}` (auth)
  - Data: `{ lecture_progress: LectureProgress[] }`
- `GET /progress/courses/{course_id}/completion` (auth) → Data: `CourseCompletion`
- `POST /progress/lectures/complete` (auth)
  - Body: `{ course_id, lecture_id, watch_time_seconds }`
  - Data: `{ progress: UserProgress, course_completed: boolean }`

Enrollments

- `POST /enrollments` (auth)
  - Body: `{ course_id }`
  - Data: `Enrollment`
- `GET /enrollments` (auth)
  - Query: `page?`, `page_size?`, `status?`
  - Data: `{ enrollments: Enrollment[], total, page, page_size, total_pages }`
- `GET /enrollments/courses/{course_id}` (auth) → Data: `Enrollment`

Analytics

- `GET /analytics/user` (auth) → Data: `UserAnalytics`

Files (Bucket Service via Gateway)

- `POST /files/upload` (auth, multipart/form-data)
  - Form fields:
    - `file` (required, file)
    - `bucket` (optional, default `general`)
    - `is_public` (optional, `true|false`, default `false`)
    - `metadata` (optional, string)
  - Response: `{ file_id, filename, size, content_type, url, thumbnail_url? }`
- `GET /files` (auth)
  - Query: `page?`, `limit?`, `bucket?`, `content_type?`, `search?`, `is_public?`, `sort?`, `order?`
  - Response: `{ files: File[], total, page, limit, total_pages }`
- `GET /files/{fileId}` (auth)
  - Redirects to a presigned download URL if permitted.
- `DELETE /files/{fileId}` (auth)
  - Response: `{ message: "File deleted successfully" }`
- `GET /files/{fileId}/metadata` (auth) → `File`
- Multipart uploads (auth):
  - `POST /files/upload/start` Body: `{ filename, content_type, size, bucket? }` → `{ session_id, upload_id, presigned_urls[] }`
  - `POST /files/upload/{sessionId}/complete` Body: provider-specific parts manifest → finalization JSON
  - `DELETE /files/upload/{sessionId}` → abort session
  - `GET /files/upload/{sessionId}/progress` → progress JSON

Video (Video Service via Gateway)

- Public
  - `GET /videos/search?q=...&limit?&offset?` → `{ videos, query, limit, offset }`
  - `GET /videos/{video_id}` (optional auth for private videos) → `Video`
  - `GET /videos/course/{course_id}?limit?&offset?` → `{ videos, course_id, limit, offset }`
  - `POST /videos/webhooks/cloudflare` (webhook)
- Auth required
  - `POST /videos/upload` Body: `{ title, description?, course_id?, lecture_id?, visibility? }` → Upload handoff `{ video_id, cloudflare_uid, upload_url, ... }`
  - `PUT /videos/{video_id}/update` Body: `{ title?, description?, visibility? }` → `Video`
  - `DELETE /videos/{video_id}/delete` → `{ message }`
  - `GET /videos/user/{user_id}?limit?&offset?` → `{ videos, limit, offset }`
  - `POST /videos/{video_id}/sessions` Body: `{ device_info: { user_agent?, screen_resolution?, connection_type? } }` → session `{ session_id, stream_url, websocket_url, ... }`
  - `PUT /videos/sessions/{session_id}/progress` Body: `{ current_time, quality? }` → `{ message }`
  - `POST /videos/sessions/{session_id}/network` Body: `{ bandwidth_mbps, latency_ms, packet_loss, connection_type, buffer_health, current_time, current_quality }` → `{ recommended_quality, quality_score, should_preload, buffer_target }`
  - `GET /videos/{video_id}/analytics?period=7d` → analytics JSON (mock data in dev)
  - WebSocket helpers (not proxied yet in gateway): stats/broadcast/session endpoints exposed but WS proxy is TODO.

Chatbot (Chatbot Service via Gateway, auth)

- Sessions
  - `POST /chat/sessions` Body: `{ title?, course_id? }` (creates an empty session) → `ChatSession`
  - `GET /chat/sessions?limit?&offset?` → `{ sessions, limit, offset }`
  - `GET /chat/sessions/{sessionId}` → `ChatSession` (+recent messages)
  - `PUT /chat/sessions/{sessionId}` Body: `{ title?, is_active? }` → `ChatSession`
  - `DELETE /chat/sessions/{sessionId}` → `{ message }`
- Messages
  - `POST /chat/message` Body: `{ session_id?, course_id?, message, context? }`
    - If `session_id` omitted, a new session is created.
    - Response: `{ session_id, message_id, role: "assistant", content, tokens_used, created_at }`
  - `GET /chat/sessions/{sessionId}/messages?limit?&offset?` → `{ messages, limit, offset }`
- Realtime
  - `GET /chat/ws` (not implemented through gateway; WS available on service directly in dev).

Forum (Forum Service via Gateway)

- Public
  - `GET /forum/topics?search?&category?&tags?&sort_by?&sort_order?&course_id?&page?&limit?` → `{ topics, total, page, limit }`
  - `GET /forum/topics/{topicId}` → `TopicResponse`
  - `GET /forum/topics/{topicId}/posts?parent_id?&sort_by?&sort_order?&page?&limit?` → `{ posts, total, page, limit }`
  - `GET /forum/posts/{postId}` → `PostResponse`
  - `GET /forum/search?q=...&course_id?&category?` → `{ topics, total, query }`
- Auth required
  - `POST /forum/topics` Body: `{ title, description, category, tags?[], course_id? }` → `Topic`
  - `PUT /forum/topics/{topicId}` Body: `{ title?, description?, category?, tags?[] }` → `Topic`
  - `DELETE /forum/topics/{topicId}` → `{ message }`
  - `PUT /forum/topics/{topicId}/sticky` → `Topic` (moderator/admin)
  - `PUT /forum/topics/{topicId}/lock` → `Topic` (moderator/admin)
  - `POST /forum/posts` Body: `{ topic_id, parent_id?, content }` → `Post`
  - `PUT /forum/posts/{postId}` Body: `{ content? }` → `Post`
  - `DELETE /forum/posts/{postId}` → `{ message }`
  - `POST /forum/votes` Body: `{ post_id, vote_type: "up"|"down" }` → `{ message }`
  - `DELETE /forum/posts/{postId}/vote` → `{ message }`

Payments (Payment Service via Gateway, auth)

- Payment Methods
  - `POST /payments/methods` Body: `{ provider, token, card_last_four?, card_expiry?, is_default? }` → `PaymentMethod`
  - `GET /payments/methods` → `{ payment_methods: PaymentMethod[] }`
  - `PUT /payments/methods/{methodId}` Body: `{ card_last_four?, card_expiry?, is_default? }` → `PaymentMethod`
  - `DELETE /payments/methods/{methodId}` → `{ message }`
  - `PUT /payments/methods/{methodId}/default` → `PaymentMethod`
- Purchase
  - `POST /payments/purchase/course/{courseId}` Body: `{ payment_method_id, amount, currency? }` → `Transaction`
  - `POST /payments/validate` Body: `{ transaction_id, provider_token }` → `Transaction`
- Transactions
  - `GET /payments/transactions?limit?&offset?` → `{ transactions[], limit, offset }`
  - `GET /payments/transactions/{transactionId}` → `Transaction`
  - `POST /payments/transactions/{transactionId}/refund` Body: `{ amount?, reason? }` → `Transaction`
- Subscriptions
  - `POST /payments/subscriptions` Body: `{ payment_method_id, plan_name, billing_period, price }` → `Subscription`
  - `GET /payments/subscriptions` → `{ subscriptions: Subscription[] }`
  - `PUT /payments/subscriptions/{subscriptionId}` Body: `{ payment_method_id?, status? }` → `Subscription`
  - `DELETE /payments/subscriptions/{subscriptionId}` → `{ message }`

Types (Selected)

- Course: `{ id, title, description, instructor_id, instructor_name, category, level, price, currency, thumbnail_url, status, duration_minutes, enrollment_count, rating, rating_count, tags[], created_at, updated_at }`
- Lecture: `{ id, course_id, title, description, order_number, duration_minutes, video_url, video_id, status, is_free, created_at, updated_at }`
- Enrollment: `{ id, user_id, course_id, status, progress_percentage, completed_lectures, total_lectures, total_watch_time_seconds, enrolled_at, last_accessed?, completed_at?, created_at, updated_at }`
- UserProgress: `{ id, user_id, course_id, lecture_id, progress_percentage, watch_time_seconds, is_completed, last_accessed, created_at, updated_at, completed_at? }`
- LectureProgress: `{ lecture_id, title, order_number, progress_percentage, watch_time_seconds, is_completed, last_accessed?, completed_at? }`
- CourseCompletion: `{ course_id, course_title, user_id, completion_percentage, completed_lectures, total_lectures, total_watch_time_seconds, lecture_progress[], started_at?, completed_at?, last_accessed? }`
- PaymentMethod: `{ id, user_id, provider, token, card_last_four?, card_expiry?, is_default, created_at, updated_at }`
- Transaction: `{ id, user_id, course_id?, payment_method_id?, amount, currency, status, transaction_reference?, created_at, updated_at }`
- Subscription: `{ id, user_id, payment_method_id?, plan_name, status, billing_period, next_billing_date?, price, created_at, updated_at }`
- Video: `{ id, cloudflare_uid, title, description, duration_seconds?, file_size_bytes?, upload_user_id, course_id?, lecture_id?, status, visibility, thumbnail_url, stream_url, preview_url, metadata, created_at, updated_at }`
- ChatSession: `{ id, user_id, course_id?, title, is_active, created_at, updated_at }`
- TopicResponse/PostResponse per forum model.

Notes & Gotchas

- Video WebSocket endpoints are exposed by the video service; the gateway’s WS proxy is not yet implemented. For local dev, hit the video service directly if you need WS.
- File downloads return an HTTP redirect to a signed URL; follow redirects in the client or let the browser handle them.
- Some services perform role checks; instructor/admin-only operations are guarded by gateway middleware.
- Payment routes are proxied and expect `X-User-ID` to be forwarded; the gateway injects this from JWT.

Changelog

- 2025-09-11: Initial consolidated gateway-based API guide generated from code.

