package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/study-platform/api-gateway/internal/handler"
	"github.com/study-platform/api-gateway/internal/middleware"
)

type Router struct {
	authHandler                *handler.AuthHandler
	courseHandler              *handler.CourseHandler
	progressHandler            *handler.ProgressHandler
	videoHandler               *handler.VideoHandler
	bucketHandler              *handler.BucketHandler
	chatbotHandler             *handler.ChatbotHandler
	forumHandler               *handler.ForumHandler
	paymentHandler             *handler.PaymentHandler
	lemonSqueezyHandler        *handler.LemonSqueezyHandler
	instructorDashboardHandler *handler.InstructorDashboardHandler
	studentDashboardHandler    *handler.StudentDashboardHandler
	docsHandler                *handler.DocsHandler
	courseAccessHandler        *handler.CourseAccessHandler
	progressTrackingHandler    *handler.ProgressTrackingHandler
	notesHandler               *handler.NotesHandler
	authMiddleware             *middleware.AuthMiddleware
	loggingMiddleware          *middleware.LoggingMiddleware
	rateLimitMiddleware        *middleware.RateLimitMiddleware
	circuitBreakerManager      *middleware.CircuitBreakerManager
	securityMiddleware         *middleware.SecurityMiddleware
}

func NewRouter(
	authHandler *handler.AuthHandler,
	courseHandler *handler.CourseHandler,
	progressHandler *handler.ProgressHandler,
	videoHandler *handler.VideoHandler,
	bucketHandler *handler.BucketHandler,
	chatbotHandler *handler.ChatbotHandler,
	forumHandler *handler.ForumHandler,
	paymentHandler *handler.PaymentHandler,
	lemonSqueezyHandler *handler.LemonSqueezyHandler,
	instructorDashboardHandler *handler.InstructorDashboardHandler,
	studentDashboardHandler *handler.StudentDashboardHandler,
	docsHandler *handler.DocsHandler,
	courseAccessHandler *handler.CourseAccessHandler,
	progressTrackingHandler *handler.ProgressTrackingHandler,
	notesHandler *handler.NotesHandler,
	authMiddleware *middleware.AuthMiddleware,
	loggingMiddleware *middleware.LoggingMiddleware,
	rateLimitMiddleware *middleware.RateLimitMiddleware,
	circuitBreakerManager *middleware.CircuitBreakerManager,
	securityMiddleware *middleware.SecurityMiddleware,
) *Router {
	return &Router{
		authHandler:                authHandler,
		courseHandler:              courseHandler,
		progressHandler:            progressHandler,
		videoHandler:               videoHandler,
		bucketHandler:              bucketHandler,
		chatbotHandler:             chatbotHandler,
		forumHandler:               forumHandler,
		paymentHandler:             paymentHandler,
		lemonSqueezyHandler:        lemonSqueezyHandler,
		instructorDashboardHandler: instructorDashboardHandler,
		studentDashboardHandler:    studentDashboardHandler,
		docsHandler:                docsHandler,
		courseAccessHandler:        courseAccessHandler,
		progressTrackingHandler:    progressTrackingHandler,
		notesHandler:               notesHandler,
		authMiddleware:             authMiddleware,
		loggingMiddleware:          loggingMiddleware,
		rateLimitMiddleware:        rateLimitMiddleware,
		circuitBreakerManager:      circuitBreakerManager,
		securityMiddleware:         securityMiddleware,
	}
}

func (rt *Router) SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	// CORS middleware will be applied globally
	r.Use(middleware.CORSMiddleware)

	// Apply other middleware
	r.Use(middleware.SetJSONContentType)
	r.Use(rt.rateLimitMiddleware.RateLimit)
	r.Use(rt.loggingMiddleware.LogRequest)

	// API version prefix
	api := r.PathPrefix("/api/v1").Subrouter()

	// Create a general routes subrouter for non-specific endpoints
	generalRoutes := api.PathPrefix("/").Subrouter()

	// Health check endpoint
	generalRoutes.HandleFunc("/health", rt.healthCheck).Methods("GET")
	generalRoutes.HandleFunc("/health/circuit-breakers", rt.circuitBreakerStatus).Methods("GET")

	// Documentation endpoints
	generalRoutes.HandleFunc("/docs/openapi.json", rt.docsHandler.GetAPISpec).Methods("GET")
	generalRoutes.HandleFunc("/docs", rt.docsHandler.GetSwaggerUI).Methods("GET")

	// Debug endpoint for testing
	generalRoutes.HandleFunc("/debug/files", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Debug endpoint works", "path": "` + r.URL.Path + `"}`))
	}).Methods("GET")

	// Debug dashboard endpoint without auth
	generalRoutes.HandleFunc("/debug/dashboard", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Dashboard debug endpoint works", "path": "` + r.URL.Path + `"}`))
	}).Methods("GET")

	// Direct file routes for testing
	generalRoutes.HandleFunc("/files", rt.bucketHandler.ListFiles).Methods("GET")
	generalRoutes.HandleFunc("/files/upload", rt.bucketHandler.UploadFile).Methods("POST")

	// Auth routes (no authentication required)
	authRoutes := api.PathPrefix("/auth").Subrouter()
	authRoutes.HandleFunc("/register", rt.authHandler.Register).Methods("POST")
	authRoutes.HandleFunc("/login", rt.authHandler.Login).Methods("POST")
	authRoutes.HandleFunc("/validate", rt.authHandler.ValidateToken).Methods("POST")
	authRoutes.HandleFunc("/oauth/{provider}/url", rt.authHandler.GetOAuthURL).Methods("GET")
	authRoutes.HandleFunc("/oauth/{provider}/callback", rt.authHandler.OAuthCallback).Methods("GET")

	// Protected auth routes (authentication required)
	protectedAuthRoutes := api.PathPrefix("/auth").Subrouter()
	protectedAuthRoutes.Use(rt.authMiddleware.RequireAuth)
	protectedAuthRoutes.HandleFunc("/profile", rt.authHandler.GetProfile).Methods("GET")

	// Course routes (some public, some protected)
	courseRoutes := api.PathPrefix("/courses").Subrouter()

	// Public course routes
	courseRoutes.HandleFunc("", rt.courseHandler.ListCourses).Methods("GET")
	courseRoutes.HandleFunc("/search", rt.courseHandler.SearchCourses).Methods("GET")
	courseRoutes.HandleFunc("/{id}", rt.courseHandler.GetCourse).Methods("GET")
	courseRoutes.HandleFunc("/{course_id}/lectures", rt.courseHandler.ListLectures).Methods("GET")
	courseRoutes.HandleFunc("/lectures/{id}", rt.courseHandler.GetLecture).Methods("GET")

	// Protected course routes (authentication required)
	protectedCourseRoutes := api.PathPrefix("/courses").Subrouter()
	protectedCourseRoutes.Use(rt.authMiddleware.RequireAuth)
	protectedCourseRoutes.HandleFunc("/upload", rt.courseHandler.CreateCourseWithThumbnail).Methods("POST")
	protectedCourseRoutes.HandleFunc("", rt.courseHandler.CreateCourse).Methods("POST")
	protectedCourseRoutes.HandleFunc("/{id}", rt.courseHandler.UpdateCourse).Methods("PUT")
	protectedCourseRoutes.HandleFunc("/{id}/upload", rt.courseHandler.UpdateCourseWithThumbnail).Methods("PUT")
	protectedCourseRoutes.HandleFunc("/{id}", rt.courseHandler.DeleteCourse).Methods("DELETE")
	protectedCourseRoutes.HandleFunc("/{course_id}/lectures", rt.courseHandler.CreateLecture).Methods("POST")
	protectedCourseRoutes.HandleFunc("/{course_id}/enroll", rt.courseHandler.EnrollInCourse).Methods("POST")

	// Lecture Resource signed URL routes (require authentication)
	lectureResourceRoutes := api.PathPrefix("/lecture-resources").Subrouter()
	lectureResourceRoutes.Use(rt.authMiddleware.RequireAuth)
	lectureResourceRoutes.HandleFunc("/{resource_id}/download-url", rt.courseHandler.GetLectureResourceDownloadURL).Methods("GET")
	lectureResourceRoutes.HandleFunc("/{resource_id}/preview-url", rt.courseHandler.GetLectureResourcePreviewURL).Methods("GET")

	// Notes routes (all require authentication)
	notesRoutes := api.PathPrefix("/notes").Subrouter()
	notesRoutes.Use(rt.authMiddleware.RequireAuth)

	// Notes CRUD operations
	notesRoutes.HandleFunc("/courses/{course_id}/lectures/{lecture_id}", rt.notesHandler.GetNotesByLecture).Methods("GET")
	notesRoutes.HandleFunc("/courses/{course_id}/lectures/{lecture_id}", rt.notesHandler.CreateNote).Methods("POST")
	notesRoutes.HandleFunc("/courses/{course_id}", rt.notesHandler.GetNotesByCourse).Methods("GET")
	notesRoutes.HandleFunc("/{note_id}", rt.notesHandler.GetNote).Methods("GET")
	notesRoutes.HandleFunc("/{note_id}", rt.notesHandler.UpdateNote).Methods("PUT")
	notesRoutes.HandleFunc("/{note_id}", rt.notesHandler.DeleteNote).Methods("DELETE")

	// Progress routes (all require authentication)
	progressRoutes := api.PathPrefix("/progress").Subrouter()
	progressRoutes.Use(rt.authMiddleware.RequireAuth)

	// Progress tracking
	progressRoutes.HandleFunc("/update", rt.progressHandler.UpdateProgress).Methods("POST")
	progressRoutes.HandleFunc("/courses/{course_id}/lectures/{lecture_id}", rt.progressHandler.GetProgress).Methods("GET")
	progressRoutes.HandleFunc("/courses/{course_id}", rt.progressHandler.GetUserProgress).Methods("GET")
	progressRoutes.HandleFunc("/lectures/{course_id}", rt.progressHandler.GetLectureProgress).Methods("GET")
	progressRoutes.HandleFunc("/courses/{course_id}/completion", rt.progressHandler.GetCourseCompletion).Methods("GET")
	progressRoutes.HandleFunc("/lectures/complete", rt.progressHandler.MarkLectureComplete).Methods("POST")

	// Enrollment routes
	enrollmentRoutes := api.PathPrefix("/enrollments").Subrouter()
	enrollmentRoutes.Use(rt.authMiddleware.RequireAuth)
	enrollmentRoutes.HandleFunc("", rt.progressHandler.CreateEnrollment).Methods("POST")
	enrollmentRoutes.HandleFunc("", rt.progressHandler.ListEnrollments).Methods("GET")
	enrollmentRoutes.HandleFunc("/courses/{course_id}", rt.progressHandler.GetEnrollment).Methods("GET")

	// Analytics routes
	analyticsRoutes := api.PathPrefix("/analytics").Subrouter()
	analyticsRoutes.Use(rt.authMiddleware.RequireAuth)
	analyticsRoutes.HandleFunc("/user", rt.progressHandler.GetUserAnalytics).Methods("GET")

	// Course Access Control Routes (new endpoints for enrollment verification)
	courseAccessRoutes := api.PathPrefix("/courses").Subrouter()
	courseAccessRoutes.Use(rt.authMiddleware.RequireAuth)

	// Course access endpoints
	courseAccessRoutes.HandleFunc("/{courseId}/access", rt.courseAccessHandler.CheckCourseAccess).Methods("GET")
	courseAccessRoutes.HandleFunc("/{courseId}/lectures/enrolled", rt.courseAccessHandler.GetEnrolledCourseLectures).Methods("GET")

	// My enrolled courses endpoint
	myCoursesRoutes := api.PathPrefix("/enrollments").Subrouter()
	myCoursesRoutes.Use(rt.authMiddleware.RequireAuth)
	myCoursesRoutes.HandleFunc("/my-courses", rt.courseAccessHandler.GetMyEnrolledCourses).Methods("GET")

	// Lecture Access and Video Streaming Routes
	lectureRoutes := api.PathPrefix("/lectures").Subrouter()
	lectureRoutes.Use(rt.authMiddleware.RequireAuth)

	// Video streaming with enrollment verification
	lectureRoutes.HandleFunc("/{lectureId}/stream", rt.videoHandler.GetLectureStreamURL).Methods("GET")

	// Progress Tracking Routes (enhanced)
	progressTrackingRoutes := api.PathPrefix("/progress").Subrouter()
	progressTrackingRoutes.Use(rt.authMiddleware.RequireAuth)

	// Progress tracking endpoints
	progressTrackingRoutes.HandleFunc("/track", rt.progressTrackingHandler.TrackProgress).Methods("POST")
	progressTrackingRoutes.HandleFunc("/complete", rt.progressTrackingHandler.MarkLectureComplete).Methods("POST")
	progressTrackingRoutes.HandleFunc("/courses/{courseId}/user", rt.progressTrackingHandler.GetUserProgress).Methods("GET")
	progressTrackingRoutes.HandleFunc("/lectures/{courseId}/{lectureId}", rt.progressTrackingHandler.GetLectureProgress).Methods("GET")

	// File management routes (protected - require authentication)
	protectedFileRoutes := api.PathPrefix("/files").Subrouter()
	protectedFileRoutes.Use(rt.authMiddleware.RequireAuth)

	// File operations
	protectedFileRoutes.HandleFunc("/upload", rt.bucketHandler.UploadFile).Methods("POST")
	protectedFileRoutes.HandleFunc("", rt.bucketHandler.ListFiles).Methods("GET")
	protectedFileRoutes.HandleFunc("/{fileId}", rt.bucketHandler.DownloadFile).Methods("GET")
	protectedFileRoutes.HandleFunc("/{fileId}", rt.bucketHandler.DeleteFile).Methods("DELETE")
	protectedFileRoutes.HandleFunc("/{fileId}/metadata", rt.bucketHandler.GetFileMetadata).Methods("GET")

	// Multipart upload operations
	protectedFileRoutes.HandleFunc("/upload/start", rt.bucketHandler.StartMultipartUpload).Methods("POST")
	protectedFileRoutes.HandleFunc("/upload/{sessionId}/complete", rt.bucketHandler.CompleteMultipartUpload).Methods("POST")
	protectedFileRoutes.HandleFunc("/upload/{sessionId}", rt.bucketHandler.AbortMultipartUpload).Methods("DELETE")
	protectedFileRoutes.HandleFunc("/upload/{sessionId}/progress", rt.bucketHandler.GetUploadProgress).Methods("GET")

	// Student Dashboard routes (require authentication)
	dashboardRoutes := api.PathPrefix("/dashboard").Subrouter()
	// dashboardRoutes.Use(rt.authMiddleware.RequireAuth) // Temporarily disabled for debugging

	// Student dashboard endpoints
	if rt.studentDashboardHandler == nil {
		panic("studentDashboardHandler is nil!")
	}
	dashboardRoutes.HandleFunc("/user/{userId}", rt.studentDashboardHandler.GetUserDashboard).Methods("GET")
	dashboardRoutes.HandleFunc("/goals/{goalId}", rt.studentDashboardHandler.UpdateStudyGoal).Methods("PATCH")
	dashboardRoutes.HandleFunc("/goals", rt.studentDashboardHandler.CreateStudyGoal).Methods("POST")
	dashboardRoutes.HandleFunc("/stats/weekly/{userId}", rt.studentDashboardHandler.GetWeeklyStats).Methods("GET")

	// Debug test endpoint
	dashboardRoutes.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Student dashboard routes are working"}`))
	}).Methods("GET")

	// Admin routes (require admin role)
	adminRoutes := api.PathPrefix("/admin").Subrouter()
	adminRoutes.Use(rt.authMiddleware.RequireAuth)
	adminRoutes.Use(rt.authMiddleware.RequireAdmin)
	// Add admin-specific routes here as needed

	// Instructor routes (require instructor role)
	instructorRoutes := api.PathPrefix("/instructor").Subrouter()
	instructorRoutes.Use(rt.authMiddleware.RequireAuth)
	instructorRoutes.Use(rt.authMiddleware.RequireInstructor)

	// Instructor Dashboard Overview
	instructorRoutes.HandleFunc("/dashboard/overview", rt.instructorDashboardHandler.GetDashboardOverview).Methods("GET")

	// Instructor Course Management
	instructorRoutes.HandleFunc("/courses", rt.instructorDashboardHandler.GetInstructorCourses).Methods("GET")
	instructorRoutes.HandleFunc("/courses", rt.instructorDashboardHandler.CreateCourse).Methods("POST")
	instructorRoutes.HandleFunc("/courses/{id}/bulk-operations", rt.instructorDashboardHandler.BulkCourseOperations).Methods("POST")

	// Instructor Analytics
	instructorRoutes.HandleFunc("/analytics/revenue", rt.instructorDashboardHandler.GetRevenueAnalytics).Methods("GET")
	instructorRoutes.HandleFunc("/analytics/engagement", rt.instructorDashboardHandler.GetEngagementAnalytics).Methods("GET")

	// Instructor Student Management
	instructorRoutes.HandleFunc("/students", rt.instructorDashboardHandler.GetStudents).Methods("GET")

	// Instructor Communication
	instructorRoutes.HandleFunc("/communication/broadcast", rt.instructorDashboardHandler.BroadcastCommunication).Methods("POST")

	// Instructor Video Analytics
	instructorRoutes.HandleFunc("/videos/analytics", rt.instructorDashboardHandler.GetVideoAnalytics).Methods("GET")

	// Instructor Team Management
	instructorRoutes.HandleFunc("/team", rt.instructorDashboardHandler.GetTeamMembers).Methods("GET")
	instructorRoutes.HandleFunc("/team/invite", rt.instructorDashboardHandler.InviteTeamMember).Methods("POST")
	instructorRoutes.HandleFunc("/team/{id}", rt.instructorDashboardHandler.UpdateTeamMember).Methods("PUT")
	instructorRoutes.HandleFunc("/team/{id}", rt.instructorDashboardHandler.RemoveTeamMember).Methods("DELETE")

	// Additional Instructor Routes
	instructorRoutes.HandleFunc("/dashboard/settings", rt.instructorDashboardHandler.UpdateDashboardSettings).Methods("PUT")
	instructorRoutes.HandleFunc("/courses/{id}/analytics", rt.instructorDashboardHandler.GetCourseAnalytics).Methods("GET")
	instructorRoutes.HandleFunc("/courses/{id}", rt.instructorDashboardHandler.GetCourse).Methods("GET")
	instructorRoutes.HandleFunc("/courses/{id}", rt.instructorDashboardHandler.UpdateCourse).Methods("PUT")
	instructorRoutes.HandleFunc("/courses/{id}", rt.instructorDashboardHandler.DeleteCourse).Methods("DELETE")
	instructorRoutes.HandleFunc("/students/{id}", rt.instructorDashboardHandler.GetStudentDetails).Methods("GET")
	instructorRoutes.HandleFunc("/analytics/students", rt.instructorDashboardHandler.GetStudentAnalytics).Methods("GET")
	instructorRoutes.HandleFunc("/videos/{id}/engagement", rt.instructorDashboardHandler.GetVideoEngagement).Methods("GET")
	instructorRoutes.HandleFunc("/communication/history", rt.instructorDashboardHandler.GetCommunicationHistory).Methods("GET")
	instructorRoutes.HandleFunc("/communication/automated", rt.instructorDashboardHandler.SetupAutomatedMessages).Methods("POST")
	instructorRoutes.HandleFunc("/suggestions", rt.instructorDashboardHandler.GetAISuggestions).Methods("GET")
	instructorRoutes.HandleFunc("/suggestions/{id}/implement", rt.instructorDashboardHandler.ImplementSuggestion).Methods("POST")
	instructorRoutes.HandleFunc("/notifications", rt.instructorDashboardHandler.GetNotifications).Methods("GET")
	instructorRoutes.HandleFunc("/notifications/settings", rt.instructorDashboardHandler.GetNotificationSettings).Methods("GET")
	instructorRoutes.HandleFunc("/notifications/settings", rt.instructorDashboardHandler.UpdateNotificationSettings).Methods("PUT")

	// Video Upload and Management
	instructorRoutes.HandleFunc("/videos/upload", rt.instructorDashboardHandler.UploadVideo).Methods("POST")
	instructorRoutes.HandleFunc("/videos/status/{lectureId}", rt.instructorDashboardHandler.GetVideoStatus).Methods("GET")

	// Instructor Dashboard Health Check
	instructorRoutes.HandleFunc("/dashboard/health", rt.instructorDashboardHandler.HealthCheck).Methods("GET")

	// Chatbot routes
	chatRoutes := api.PathPrefix("/chat").Subrouter()
	chatRoutes.Use(rt.authMiddleware.RequireAuth)

	// Chat session management
	chatRoutes.HandleFunc("/sessions", rt.chatbotHandler.CreateSession).Methods("POST")
	chatRoutes.HandleFunc("/sessions", rt.chatbotHandler.GetUserSessions).Methods("GET")
	chatRoutes.HandleFunc("/sessions/{sessionId}", rt.chatbotHandler.GetSession).Methods("GET")
	chatRoutes.HandleFunc("/sessions/{sessionId}", rt.chatbotHandler.UpdateSession).Methods("PUT")
	chatRoutes.HandleFunc("/sessions/{sessionId}", rt.chatbotHandler.DeleteSession).Methods("DELETE")

	// Chat messaging
	chatRoutes.HandleFunc("/message", rt.chatbotHandler.SendMessage).Methods("POST")
	chatRoutes.HandleFunc("/sessions/{sessionId}/messages", rt.chatbotHandler.GetMessages).Methods("GET")

	// WebSocket endpoint for real-time chat
	chatRoutes.HandleFunc("/ws", rt.chatbotHandler.HandleWebSocket).Methods("GET")

	// Forum routes
	forumRoutes := api.PathPrefix("/forum").Subrouter()

	// Public forum routes (can work without auth but benefit from it)
	forumRoutes.HandleFunc("/topics", rt.forumHandler.ListTopics).Methods("GET")
	forumRoutes.HandleFunc("/topics/{topicId}", rt.forumHandler.GetTopic).Methods("GET")
	forumRoutes.HandleFunc("/topics/{topicId}/posts", rt.forumHandler.ListPosts).Methods("GET")
	forumRoutes.HandleFunc("/posts/{postId}", rt.forumHandler.GetPost).Methods("GET")
	forumRoutes.HandleFunc("/search", rt.forumHandler.SearchTopics).Methods("GET")

	// Course-specific forum routes
	forumRoutes.HandleFunc("/courses/{courseId}/topics", rt.forumHandler.ListCourseTopics).Methods("GET")

	// Protected forum routes (require authentication)
	protectedForumRoutes := api.PathPrefix("/forum").Subrouter()
	protectedForumRoutes.Use(rt.authMiddleware.RequireAuth)

	// Topic management
	protectedForumRoutes.HandleFunc("/topics", rt.forumHandler.CreateTopic).Methods("POST")
	protectedForumRoutes.HandleFunc("/topics/{topicId}", rt.forumHandler.UpdateTopic).Methods("PUT")
	protectedForumRoutes.HandleFunc("/topics/{topicId}", rt.forumHandler.DeleteTopic).Methods("DELETE")
	protectedForumRoutes.HandleFunc("/topics/{topicId}/sticky", rt.forumHandler.ToggleTopicSticky).Methods("PUT")
	protectedForumRoutes.HandleFunc("/topics/{topicId}/lock", rt.forumHandler.ToggleTopicLock).Methods("PUT")

	// Post management
	protectedForumRoutes.HandleFunc("/posts", rt.forumHandler.CreatePost).Methods("POST")
	protectedForumRoutes.HandleFunc("/posts/{postId}", rt.forumHandler.UpdatePost).Methods("PUT")
	protectedForumRoutes.HandleFunc("/posts/{postId}", rt.forumHandler.DeletePost).Methods("DELETE")
	protectedForumRoutes.HandleFunc("/posts/{postId}/answer", rt.forumHandler.MarkPostAsAnswer).Methods("PUT")
	protectedForumRoutes.HandleFunc("/posts/{postId}/pin", rt.forumHandler.TogglePostPin).Methods("PUT")

	// Voting
	protectedForumRoutes.HandleFunc("/votes", rt.forumHandler.VotePost).Methods("POST")
	protectedForumRoutes.HandleFunc("/posts/{postId}/vote", rt.forumHandler.RemoveVote).Methods("DELETE")

	// Approval routes (instructor/admin only)
	protectedForumRoutes.HandleFunc("/pending/topics", rt.forumHandler.GetPendingTopics).Methods("GET")
	protectedForumRoutes.HandleFunc("/topics/{topicId}/approve", rt.forumHandler.ApproveTopic).Methods("PUT")
	protectedForumRoutes.HandleFunc("/topics/{topicId}/reject", rt.forumHandler.RejectTopic).Methods("PUT")
	protectedForumRoutes.HandleFunc("/posts/{postId}/approve", rt.forumHandler.ApprovePost).Methods("PUT")
	protectedForumRoutes.HandleFunc("/posts/{postId}/reject", rt.forumHandler.RejectPost).Methods("PUT")

	// Pin order management (instructor/admin only)
	protectedForumRoutes.HandleFunc("/topics/{topicId}/pin-order", rt.forumHandler.SetTopicPinOrder).Methods("PUT")
	protectedForumRoutes.HandleFunc("/posts/{postId}/pin-order", rt.forumHandler.SetPostPinOrder).Methods("PUT")

	// Payment routes (all require authentication)
	paymentRoutes := api.PathPrefix("/payments").Subrouter()
	paymentRoutes.Use(rt.authMiddleware.RequireAuth)

	// Payment methods
	paymentRoutes.HandleFunc("/methods", rt.paymentHandler.CreatePaymentMethod).Methods("POST")
	paymentRoutes.HandleFunc("/methods", rt.paymentHandler.GetPaymentMethods).Methods("GET")
	paymentRoutes.HandleFunc("/methods/{methodId}", rt.paymentHandler.UpdatePaymentMethod).Methods("PUT")
	paymentRoutes.HandleFunc("/methods/{methodId}", rt.paymentHandler.DeletePaymentMethod).Methods("DELETE")
	paymentRoutes.HandleFunc("/methods/{methodId}/default", rt.paymentHandler.SetDefaultPaymentMethod).Methods("PUT")

	// Course purchase
	paymentRoutes.HandleFunc("/purchase/course/{courseId}", rt.paymentHandler.PurchaseCourse).Methods("POST")
	paymentRoutes.HandleFunc("/validate", rt.paymentHandler.ValidatePayment).Methods("POST")

	// Transactions
	paymentRoutes.HandleFunc("/transactions", rt.paymentHandler.GetTransactions).Methods("GET")
	paymentRoutes.HandleFunc("/transactions/{transactionId}", rt.paymentHandler.GetTransaction).Methods("GET")
	paymentRoutes.HandleFunc("/transactions/{transactionId}/refund", rt.paymentHandler.RefundTransaction).Methods("POST")

	// Subscriptions
	paymentRoutes.HandleFunc("/subscriptions", rt.paymentHandler.CreateSubscription).Methods("POST")
	paymentRoutes.HandleFunc("/subscriptions", rt.paymentHandler.GetSubscriptions).Methods("GET")
	paymentRoutes.HandleFunc("/subscriptions/{subscriptionId}", rt.paymentHandler.UpdateSubscription).Methods("PUT")
	paymentRoutes.HandleFunc("/subscriptions/{subscriptionId}", rt.paymentHandler.CancelSubscription).Methods("DELETE")

	// Lemon Squeezy routes (require authentication)
	paymentRoutes.HandleFunc("/lemonsqueezy/checkout/course/{course_id}", rt.paymentHandler.CreateLemonSqueezyCheckout).Methods("POST")
	paymentRoutes.HandleFunc("/lemonsqueezy/verify/{order_id}", rt.paymentHandler.VerifyLemonSqueezyPayment).Methods("POST")
	paymentRoutes.HandleFunc("/lemonsqueezy/products", rt.paymentHandler.GetLemonSqueezyProducts).Methods("GET")
	paymentRoutes.HandleFunc("/lemonsqueezy/variants", rt.paymentHandler.GetLemonSqueezyVariants).Methods("GET")

	// Stripe routes (match frontend paths: /payments/stripe/*)
	// Note: CORS handled by global middleware to avoid duplicates
	stripeRoutes := api.PathPrefix("/payments/stripe").Subrouter()

	// Public Stripe routes (no auth required) - order matters!
	stripeRoutes.HandleFunc("/config", rt.paymentHandler.GetStripeConfig).Methods("GET")

	// Protected Stripe routes - use Handle() with middleware
	stripeRoutes.Handle("/payment-intents", rt.authMiddleware.RequireAuth(http.HandlerFunc(rt.paymentHandler.CreateStripePaymentIntent))).Methods("POST")
	stripeRoutes.Handle("/payment-intents/{payment_intent_id}", rt.authMiddleware.RequireAuth(http.HandlerFunc(rt.paymentHandler.GetStripePaymentIntent))).Methods("GET")
	stripeRoutes.Handle("/payment-intents/confirm", rt.authMiddleware.RequireAuth(http.HandlerFunc(rt.paymentHandler.ConfirmStripePaymentIntent))).Methods("POST")
	stripeRoutes.Handle("/transactions", rt.authMiddleware.RequireAuth(http.HandlerFunc(rt.paymentHandler.ListStripeTransactions))).Methods("GET")

	// Additional frontend endpoints
	stripeRoutes.Handle("/customers", rt.authMiddleware.RequireAuth(http.HandlerFunc(rt.paymentHandler.CreateStripePaymentIntent))).Methods("POST")
	stripeRoutes.Handle("/customers/{customer_id}", rt.authMiddleware.RequireAuth(http.HandlerFunc(rt.paymentHandler.GetStripePaymentIntent))).Methods("GET", "PUT")
	stripeRoutes.Handle("/products", rt.authMiddleware.RequireAuth(http.HandlerFunc(rt.paymentHandler.CreateStripePaymentIntent))).Methods("POST", "GET")
	stripeRoutes.Handle("/products/{product_id}", rt.authMiddleware.RequireAuth(http.HandlerFunc(rt.paymentHandler.GetStripePaymentIntent))).Methods("GET")
	stripeRoutes.Handle("/prices", rt.authMiddleware.RequireAuth(http.HandlerFunc(rt.paymentHandler.CreateStripePaymentIntent))).Methods("POST", "GET")
	stripeRoutes.Handle("/prices/{price_id}", rt.authMiddleware.RequireAuth(http.HandlerFunc(rt.paymentHandler.GetStripePaymentIntent))).Methods("GET")
	stripeRoutes.Handle("/subscriptions", rt.authMiddleware.RequireAuth(http.HandlerFunc(rt.paymentHandler.CreateStripePaymentIntent))).Methods("POST", "GET")
	stripeRoutes.Handle("/subscriptions/{subscription_id}", rt.authMiddleware.RequireAuth(http.HandlerFunc(rt.paymentHandler.GetStripePaymentIntent))).Methods("GET", "PUT", "DELETE")
	stripeRoutes.Handle("/purchase/course/{course_id}", rt.authMiddleware.RequireAuth(http.HandlerFunc(rt.paymentHandler.PurchaseCourse))).Methods("POST")
	stripeRoutes.Handle("/webhooks/verify", rt.authMiddleware.RequireAuth(http.HandlerFunc(rt.paymentHandler.CreateStripePaymentIntent))).Methods("POST")

	// Stripe webhook (no authentication required)
	api.HandleFunc("/payments/stripe/webhook", rt.paymentHandler.HandleStripeWebhook).Methods("POST")

	// Lemon Squeezy webhook (no authentication required)
	api.HandleFunc("/lemonsqueezy/webhook", rt.paymentHandler.HandleLemonSqueezyWebhook).Methods("POST")

	// LemonSqueezy Course Integration Routes (protected)
	lemonSqueezyRoutes := api.PathPrefix("/lemonsqueezy").Subrouter()
	lemonSqueezyRoutes.Use(rt.authMiddleware.RequireAuth)

	// Course-specific LemonSqueezy operations
	lemonSqueezyRoutes.HandleFunc("/courses/{courseId}/checkout", rt.lemonSqueezyHandler.CreateCourseCheckout).Methods("POST")
	lemonSqueezyRoutes.HandleFunc("/courses/{courseId}/link-product", rt.lemonSqueezyHandler.LinkCourseToLemonSqueezyProduct).Methods("POST")
	lemonSqueezyRoutes.HandleFunc("/courses/{courseId}/unlink-product", rt.lemonSqueezyHandler.UnlinkCourseFromLemonSqueezy).Methods("DELETE")
	lemonSqueezyRoutes.HandleFunc("/checkouts/{checkoutId}/status", rt.lemonSqueezyHandler.GetCourseCheckoutStatus).Methods("GET")
	lemonSqueezyRoutes.HandleFunc("/courses/products", rt.lemonSqueezyHandler.ListCourseProducts).Methods("GET")

	// Additional LemonSqueezy webhook endpoint for course-specific events
	api.HandleFunc("/lemonsqueezy/webhook/courses", rt.lemonSqueezyHandler.HandleLemonSqueezyWebhook).Methods("POST")

	// Video routes - single subrouter with CORS, selective auth per endpoint
	videoRoutes := api.PathPrefix("/videos").Subrouter()

	// Catch-all OPTIONS handler for all video routes - place this FIRST
	videoRoutes.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// CORS middleware already applied above, just return success
		w.WriteHeader(http.StatusNoContent)
	})

	// Public video routes (no auth required)
	videoRoutes.HandleFunc("/search", rt.videoHandler.SearchVideos).Methods("GET")
	videoRoutes.HandleFunc("/webhooks/cloudflare", rt.videoHandler.CloudflareWebhook).Methods("POST")
	videoRoutes.HandleFunc("/course/{course_id}", rt.videoHandler.ListCourseVideos).Methods("GET")

	// Protected video routes (require auth) - create handlers that include auth check
	videoRoutes.Handle("/upload-url", rt.authMiddleware.RequireAuth(http.HandlerFunc(rt.videoHandler.GetUploadURL))).Methods("POST")
	videoRoutes.Handle("/upload", rt.authMiddleware.RequireAuth(http.HandlerFunc(rt.videoHandler.UploadVideo))).Methods("POST")
	videoRoutes.Handle("/user/{user_id}", rt.authMiddleware.RequireAuth(http.HandlerFunc(rt.videoHandler.ListUserVideos))).Methods("GET")

	// Parameterized routes LAST (they catch everything) - public first
	videoRoutes.HandleFunc("/{video_id}", rt.videoHandler.GetVideo).Methods("GET")

	// Parameterized protected routes (with auth middleware)
	videoRoutes.Handle("/{video_id}/update", rt.authMiddleware.RequireAuth(http.HandlerFunc(rt.videoHandler.UpdateVideo))).Methods("PUT")
	videoRoutes.Handle("/{video_id}/delete", rt.authMiddleware.RequireAuth(http.HandlerFunc(rt.videoHandler.DeleteVideo))).Methods("DELETE")

	// Session management (protected)
	videoRoutes.Handle("/{video_id}/sessions", rt.authMiddleware.RequireAuth(http.HandlerFunc(rt.videoHandler.CreateViewingSession))).Methods("POST")
	videoRoutes.Handle("/sessions/{session_id}/progress", rt.authMiddleware.RequireAuth(http.HandlerFunc(rt.videoHandler.UpdateSessionProgress))).Methods("PUT")
	videoRoutes.Handle("/sessions/{session_id}/network", rt.authMiddleware.RequireAuth(http.HandlerFunc(rt.videoHandler.UpdateNetworkStatus))).Methods("POST")

	// Analytics (protected)
	videoRoutes.Handle("/{video_id}/analytics", rt.authMiddleware.RequireAuth(http.HandlerFunc(rt.videoHandler.GetVideoAnalytics))).Methods("GET")

	// Manual status updates (protected)
	videoRoutes.Handle("/{video_id}/status", rt.authMiddleware.RequireAuth(http.HandlerFunc(rt.videoHandler.UpdateVideoStatus))).Methods("PUT")

	// WebSocket endpoints (protected)
	videoRoutes.Handle("/ws/stats", rt.authMiddleware.RequireAuth(http.HandlerFunc(rt.videoHandler.WebSocketStats))).Methods("GET")
	videoRoutes.Handle("/ws/broadcast", rt.authMiddleware.RequireAuth(http.HandlerFunc(rt.videoHandler.BroadcastMessage))).Methods("POST")
	videoRoutes.Handle("/ws/session/{session_id}/send", rt.authMiddleware.RequireAuth(http.HandlerFunc(rt.videoHandler.SendToSession))).Methods("POST")
	videoRoutes.Handle("/ws/session/{session_id}", rt.authMiddleware.RequireAuth(http.HandlerFunc(rt.videoHandler.GetSessionInfo))).Methods("GET")
	videoRoutes.Handle("/ws/{session_id}", rt.authMiddleware.RequireAuth(http.HandlerFunc(rt.videoHandler.WebSocketProxy))).Methods("GET")

	return r
}

func (rt *Router) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Simple JSON encoding
	jsonStr := `{"status":"healthy","service":"api-gateway","version":"1.0.0"}`
	w.Write([]byte(jsonStr))
}

func (rt *Router) circuitBreakerStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	status := rt.circuitBreakerManager.GetStatus()

	// Simple JSON encoding for circuit breaker status
	jsonStr := `{"status":"healthy","circuit_breakers":{`
	first := true
	for service, state := range status {
		if !first {
			jsonStr += ","
		}
		jsonStr += `"` + service + `":` + `{"state":"` + state.(map[string]interface{})["state"].(string) + `"}`
		first = false
	}
	jsonStr += `}}`

	w.Write([]byte(jsonStr))
}
