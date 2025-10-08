package repository

import (
	"database/sql"
	"fmt"
	"time"
	"strings"

	"instructor-dashboard-service/internal/model"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type DashboardRepository struct {
	db *sql.DB
}

func NewDashboardRepository(db *sql.DB) *DashboardRepository {
	return &DashboardRepository{db: db}
}

// GetDashboardOverview retrieves the main dashboard overview data
func (r *DashboardRepository) GetDashboardOverview(instructorID uuid.UUID) (*model.DashboardOverview, error) {
	overview := &model.DashboardOverview{
		InstructorID: instructorID,
	}

	// Get basic metrics
	query := `
		SELECT
			COALESCE(SUM(CASE WHEN t.status = 'completed' THEN t.amount ELSE 0 END), 0) as total_revenue,
			COALESCE(SUM(CASE WHEN t.status = 'completed' AND t.created_at >= CURRENT_DATE - INTERVAL '30 days' THEN t.amount ELSE 0 END), 0) as monthly_revenue,
			COUNT(DISTINCT e.user_id) as total_students,
			COUNT(DISTINCT CASE WHEN p.last_accessed_at >= CURRENT_DATE - INTERVAL '7 days' THEN e.user_id END) as active_students,
			COUNT(DISTINCT c.id) as total_courses,
			COUNT(DISTINCT CASE WHEN c.status = 'published' THEN c.id END) as published_courses,
			COALESCE(AVG(c.rating), 0) as avg_course_rating,
			COALESCE(AVG(p.completion_percentage), 0) as completion_rate
		FROM courses c
		LEFT JOIN enrollments e ON c.id = e.course_id
		LEFT JOIN transactions t ON c.id = t.course_id
		LEFT JOIN progress p ON c.id = p.course_id
		WHERE c.instructor_id = $1 AND c.deleted_at IS NULL
	`

	err := r.db.QueryRow(query, instructorID).Scan(
		&overview.TotalRevenue,
		&overview.MonthlyRevenue,
		&overview.TotalStudents,
		&overview.ActiveStudents,
		&overview.TotalCourses,
		&overview.PublishedCourses,
		&overview.AvgCourseRating,
		&overview.CompletionRate,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get overview metrics: %w", err)
	}

	// Get engagement score
	engagementQuery := `
		SELECT COALESCE(AVG(v.engagement_score), 0)
		FROM videos v
		JOIN courses c ON v.course_id = c.id
		WHERE c.instructor_id = $1 AND c.deleted_at IS NULL
	`
	err = r.db.QueryRow(engagementQuery, instructorID).Scan(&overview.EngagementScore)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get engagement score: %w", err)
	}

	// Get quick stats
	quickStats, err := r.getQuickStats(instructorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get quick stats: %w", err)
	}
	overview.QuickStats = *quickStats

	// Get recent activity
	activities, err := r.getRecentActivity(instructorID, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent activity: %w", err)
	}
	overview.RecentActivity = activities

	// Get top performing course
	topCourse, err := r.getTopPerformingCourse(instructorID)
	if err == nil {
		overview.TopPerformingCourse = *topCourse
	}

	return overview, nil
}

// getQuickStats retrieves today's quick statistics
func (r *DashboardRepository) getQuickStats(instructorID uuid.UUID) (*model.QuickStats, error) {
	stats := &model.QuickStats{}

	query := `
		SELECT
			COUNT(DISTINCT CASE WHEN e.created_at >= CURRENT_DATE THEN e.user_id END) as new_enrollments_today,
			COALESCE(SUM(CASE WHEN t.status = 'completed' AND t.created_at >= CURRENT_DATE THEN t.amount ELSE 0 END), 0) as revenue_today,
			COALESCE(SUM(CASE WHEN vs.created_at >= CURRENT_DATE THEN vs.total_watch_time_seconds ELSE 0 END), 0) as video_watch_time_today,
			COUNT(DISTINCT CASE WHEN cr.created_at >= CURRENT_DATE THEN cr.id END) as new_reviews_today
		FROM courses c
		LEFT JOIN enrollments e ON c.id = e.course_id
		LEFT JOIN transactions t ON c.id = t.course_id
		LEFT JOIN videos v ON c.id = v.course_id
		LEFT JOIN viewing_sessions vs ON v.id = vs.video_id
		LEFT JOIN course_reviews cr ON c.id = cr.course_id
		WHERE c.instructor_id = $1 AND c.deleted_at IS NULL
	`

	err := r.db.QueryRow(query, instructorID).Scan(
		&stats.NewEnrollmentsToday,
		&stats.RevenueToday,
		&stats.VideoWatchTimeToday,
		&stats.NewReviewsToday,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get quick stats: %w", err)
	}

	// Get forum posts and AI interactions (simplified queries)
	forumQuery := `
		SELECT COUNT(DISTINCT fp.id)
		FROM forum_posts fp
		JOIN forum_topics ft ON fp.topic_id = ft.id
		JOIN courses c ON ft.course_id = c.id
		WHERE c.instructor_id = $1 AND c.deleted_at IS NULL AND fp.created_at >= CURRENT_DATE
	`
	r.db.QueryRow(forumQuery, instructorID).Scan(&stats.ForumPostsToday)

	// AI interactions would need to be tracked separately
	stats.AIInteractionsToday = 0 // Placeholder

	return stats, nil
}

// getRecentActivity retrieves recent activity items
func (r *DashboardRepository) getRecentActivity(instructorID uuid.UUID, limit int) ([]model.ActivityItem, error) {
	activities := []model.ActivityItem{}

	// Get recent enrollments
	enrollmentQuery := `
		SELECT
			e.id,
			'enrollment' as type,
			'New student enrolled in ' || c.title as description,
			e.user_id,
			u.first_name || ' ' || u.last_name as user_name,
			c.id as course_id,
			c.title as course_name,
			e.created_at
		FROM enrollments e
		JOIN courses c ON e.course_id = c.id
		JOIN users u ON e.user_id = u.id
		WHERE c.instructor_id = $1 AND c.deleted_at IS NULL
		ORDER BY e.created_at DESC
		LIMIT $2
	`

	rows, err := r.db.Query(enrollmentQuery, instructorID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent enrollments: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var activity model.ActivityItem
		err := rows.Scan(
			&activity.ID,
			&activity.Type,
			&activity.Description,
			&activity.UserID,
			&activity.UserName,
			&activity.CourseID,
			&activity.CourseName,
			&activity.Timestamp,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan activity: %w", err)
		}
		activities = append(activities, activity)
	}

	return activities, nil
}

// getTopPerformingCourse retrieves the top performing course
func (r *DashboardRepository) getTopPerformingCourse(instructorID uuid.UUID) (*model.Course, error) {
	course := &model.Course{}

	query := `
		SELECT
			c.id,
			c.title,
			c.description,
			c.instructor_id,
			c.status,
			c.is_paid,
			c.price,
			c.currency,
			COALESCE(c.rating, 0) as average_rating,
			c.enrollment_count,
			COALESCE(SUM(CASE WHEN t.status = 'completed' THEN t.amount ELSE 0 END), 0) as total_revenue,
			c.created_at,
			c.updated_at
		FROM courses c
		LEFT JOIN transactions t ON c.id = t.course_id
		WHERE c.instructor_id = $1 AND c.deleted_at IS NULL
		GROUP BY c.id, c.title, c.description, c.instructor_id, c.status, c.is_paid, c.price, c.currency, c.rating, c.enrollment_count, c.created_at, c.updated_at
		ORDER BY total_revenue DESC
		LIMIT 1
	`

	err := r.db.QueryRow(query, instructorID).Scan(
		&course.ID,
		&course.Title,
		&course.Description,
		&course.InstructorID,
		&course.Status,
		&course.Price, // Skip is_paid as it's computed
		&course.Price,
		&course.Currency,
		&course.AverageRating,
		&course.TotalEnrollments,
		&course.TotalRevenue,
		&course.CreatedAt,
		&course.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get top performing course: %w", err)
	}

	return course, nil
}

// GetDashboardSettings retrieves dashboard settings for an instructor
func (r *DashboardRepository) GetDashboardSettings(instructorID uuid.UUID) (*model.InstructorDashboardSettings, error) {
	settings := &model.InstructorDashboardSettings{}

	query := `
		SELECT id, instructor_id, dashboard_layout, notification_preferences,
			   analytics_preferences, default_course_settings, ai_assistance_enabled,
			   created_at, updated_at
		FROM instructor_dashboard_settings
		WHERE instructor_id = $1
	`

	err := r.db.QueryRow(query, instructorID).Scan(
		&settings.ID,
		&settings.InstructorID,
		&settings.DashboardLayout,
		&settings.NotificationPreferences,
		&settings.AnalyticsPreferences,
		&settings.DefaultCourseSettings,
		&settings.AIAssistanceEnabled,
		&settings.CreatedAt,
		&settings.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get dashboard settings: %w", err)
	}

	return settings, nil
}

// UpdateDashboardSettings updates dashboard settings
func (r *DashboardRepository) UpdateDashboardSettings(instructorID uuid.UUID, req *model.UpdateDashboardSettingsRequest) error {
	updateFields := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.DashboardLayout != nil {
		updateFields = append(updateFields, fmt.Sprintf("dashboard_layout = $%d", argIndex))
		args = append(args, req.DashboardLayout)
		argIndex++
	}

	if req.NotificationPreferences != nil {
		updateFields = append(updateFields, fmt.Sprintf("notification_preferences = $%d", argIndex))
		args = append(args, req.NotificationPreferences)
		argIndex++
	}

	if req.AnalyticsPreferences != nil {
		updateFields = append(updateFields, fmt.Sprintf("analytics_preferences = $%d", argIndex))
		args = append(args, req.AnalyticsPreferences)
		argIndex++
	}

	if req.DefaultCourseSettings != nil {
		updateFields = append(updateFields, fmt.Sprintf("default_course_settings = $%d", argIndex))
		args = append(args, req.DefaultCourseSettings)
		argIndex++
	}

	if req.AIAssistanceEnabled != nil {
		updateFields = append(updateFields, fmt.Sprintf("ai_assistance_enabled = $%d", argIndex))
		args = append(args, req.AIAssistanceEnabled)
		argIndex++
	}

	if len(updateFields) == 0 {
		return fmt.Errorf("no fields to update")
	}

	updateFields = append(updateFields, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, time.Now())
	argIndex++

	args = append(args, instructorID)

	query := fmt.Sprintf(`
		UPDATE instructor_dashboard_settings
		SET %s
		WHERE instructor_id = $%d
	`, strings.Join(updateFields, ", "), argIndex)

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update dashboard settings: %w", err)
	}

	return nil
}

// GetInstructorCourses retrieves courses for an instructor with pagination
func (r *DashboardRepository) GetInstructorCourses(instructorID uuid.UUID, filter *model.Filter, pagination *model.Pagination) ([]model.Course, error) {
	courses := []model.Course{}

	whereClause := "WHERE c.instructor_id = $1 AND c.deleted_at IS NULL"
	args := []interface{}{instructorID}
	argIndex := 2

	if filter.Status != "" {
		whereClause += fmt.Sprintf(" AND c.status = $%d", argIndex)
		args = append(args, filter.Status)
		argIndex++
	}

	if filter.StartDate != nil {
		whereClause += fmt.Sprintf(" AND c.created_at >= $%d", argIndex)
		args = append(args, *filter.StartDate)
		argIndex++
	}

	if filter.EndDate != nil {
		whereClause += fmt.Sprintf(" AND c.created_at <= $%d", argIndex)
		args = append(args, *filter.EndDate)
		argIndex++
	}

	offset := (pagination.Page - 1) * pagination.PageSize

	query := fmt.Sprintf(`
		SELECT
			c.id, c.title, c.description, c.instructor_id, c.status,
			c.price, c.currency, COALESCE(c.rating, 0) as average_rating, c.enrollment_count,
			COUNT(DISTINCT e.user_id) as active_students,
			COALESCE(AVG(p.completion_percentage), 0) as completion_rate,
			COALESCE(SUM(CASE WHEN t.status = 'completed' THEN t.amount ELSE 0 END), 0) as total_revenue,
			COALESCE(SUM(CASE WHEN t.status = 'completed' AND t.created_at >= CURRENT_DATE - INTERVAL '30 days' THEN t.amount ELSE 0 END), 0) as monthly_revenue,
			COALESCE(AVG(v.engagement_score), 0) as engagement_score,
			MAX(p.last_accessed_at) as last_activity_at,
			c.created_at, c.updated_at
		FROM courses c
		LEFT JOIN enrollments e ON c.id = e.course_id
		LEFT JOIN progress p ON c.id = p.course_id
		LEFT JOIN transactions t ON c.id = t.course_id
		LEFT JOIN videos v ON c.id = v.course_id
		%s
		GROUP BY c.id, c.title, c.description, c.instructor_id, c.status,
				 c.price, c.currency, c.rating, c.enrollment_count, c.created_at, c.updated_at
		ORDER BY c.created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)

	args = append(args, pagination.PageSize, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get instructor courses: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var course model.Course
		var price sql.NullFloat64
		err := rows.Scan(
			&course.ID,
			&course.Title,
			&course.Description,
			&course.InstructorID,
			&course.Status,
			&price,
			&course.Currency,
			&course.AverageRating,
			&course.TotalEnrollments,
			&course.ActiveStudents,
			&course.CompletionRate,
			&course.TotalRevenue,
			&course.MonthlyRevenue,
			&course.EngagementScore,
			&course.LastActivityAt,
			&course.CreatedAt,
			&course.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan course: %w", err)
		}

		// Handle nullable price
		if price.Valid {
			course.Price = price.Float64
		} else {
			course.Price = 0
		}

		courses = append(courses, course)
	}

	return courses, nil
}

// GetInstructorCoursesCount gets total count for pagination
func (r *DashboardRepository) GetInstructorCoursesCount(instructorID uuid.UUID, filter *model.Filter) (int, error) {
	whereClause := "WHERE instructor_id = $1 AND deleted_at IS NULL"
	args := []interface{}{instructorID}
	argIndex := 2

	if filter.Status != "" {
		whereClause += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, filter.Status)
		argIndex++
	}

	if filter.StartDate != nil {
		whereClause += fmt.Sprintf(" AND created_at >= $%d", argIndex)
		args = append(args, *filter.StartDate)
		argIndex++
	}

	if filter.EndDate != nil {
		whereClause += fmt.Sprintf(" AND created_at <= $%d", argIndex)
		args = append(args, *filter.EndDate)
		argIndex++
	}

	query := fmt.Sprintf("SELECT COUNT(*) FROM courses %s", whereClause)

	var count int
	err := r.db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get courses count: %w", err)
	}

	return count, nil
}

// BulkCourseOperation performs bulk operations on courses
func (r *DashboardRepository) BulkCourseOperation(instructorID uuid.UUID, req *model.BulkCourseOperationRequest) error {
	// Verify all courses belong to instructor
	courseIDsStr := make([]string, len(req.CourseIDs))
	for i, id := range req.CourseIDs {
		courseIDsStr[i] = fmt.Sprintf("'%s'", id.String())
	}

	verifyQuery := fmt.Sprintf(`
		SELECT COUNT(*) FROM courses
		WHERE instructor_id = $1 AND id = ANY($2) AND deleted_at IS NULL
	`)

	var verifiedCount int
	err := r.db.QueryRow(verifyQuery, instructorID, pq.Array(req.CourseIDs)).Scan(&verifiedCount)
	if err != nil {
		return fmt.Errorf("failed to verify course ownership: %w", err)
	}

	if verifiedCount != len(req.CourseIDs) {
		return fmt.Errorf("some courses don't belong to instructor")
	}

	// Perform the operation
	switch req.Operation {
	case "publish":
		_, err = r.db.Exec(`
			UPDATE courses
			SET status = 'published', updated_at = NOW()
			WHERE instructor_id = $1 AND id = ANY($2) AND deleted_at IS NULL
		`, instructorID, pq.Array(req.CourseIDs))

	case "unpublish":
		_, err = r.db.Exec(`
			UPDATE courses
			SET status = 'draft', updated_at = NOW()
			WHERE instructor_id = $1 AND id = ANY($2) AND deleted_at IS NULL
		`, instructorID, pq.Array(req.CourseIDs))

	case "update_price":
		if price, ok := req.Parameters["price"]; ok {
			_, err = r.db.Exec(`
				UPDATE courses
				SET price = $3, updated_at = NOW()
				WHERE instructor_id = $1 AND id = ANY($2) AND deleted_at IS NULL
			`, instructorID, pq.Array(req.CourseIDs), price)
		} else {
			return fmt.Errorf("price parameter required for update_price operation")
		}

	case "delete":
		_, err = r.db.Exec(`
			UPDATE courses
			SET deleted_at = NOW(), updated_at = NOW()
			WHERE instructor_id = $1 AND id = ANY($2) AND deleted_at IS NULL
		`, instructorID, pq.Array(req.CourseIDs))

	default:
		return fmt.Errorf("unsupported operation: %s", req.Operation)
	}

	if err != nil {
		return fmt.Errorf("failed to perform bulk operation: %w", err)
	}

	return nil
}

// GetCourse retrieves a specific course for an instructor with complete details
func (r *DashboardRepository) GetCourse(instructorID, courseID uuid.UUID) (*model.Course, error) {
	// Get course details
	query := `
		SELECT
			c.id, c.title, c.description, c.instructor_id, c.category, c.level, c.price, c.currency,
			c.status, c.language, c.thumbnail_url, c.learning_outcomes, c.requirements, c.tags,
			c.duration_minutes, c.auto_approve_enrollment, c.allow_previews, c.has_certificate, c.mobile_access,
			COALESCE(c.rating, 0) as average_rating,
			COUNT(DISTINCT e.id) as total_enrollments,
			c.created_at, c.updated_at
		FROM courses c
		LEFT JOIN enrollments e ON c.id = e.course_id
		WHERE c.id = $1 AND c.instructor_id = $2 AND c.deleted_at IS NULL
		GROUP BY c.id, c.title, c.description, c.instructor_id, c.category, c.level, c.price, c.currency,
				 c.status, c.language, c.thumbnail_url, c.learning_outcomes, c.requirements, c.tags,
				 c.duration_minutes, c.auto_approve_enrollment, c.allow_previews, c.has_certificate, c.mobile_access,
				 c.rating, c.created_at, c.updated_at
	`

	course := &model.Course{}
	var price sql.NullFloat64
	err := r.db.QueryRow(query, courseID, instructorID).Scan(
		&course.ID,
		&course.Title,
		&course.Description,
		&course.InstructorID,
		&course.Category,
		&course.Level,
		&price,
		&course.Currency,
		&course.Status,
		&course.Language,
		&course.ThumbnailURL,
		pq.Array(&course.LearningOutcomes),
		pq.Array(&course.Requirements),
		pq.Array(&course.Tags),
		&course.EstimatedDurationHours,
		&course.AutoApproveEnrollment,
		&course.AllowPreviews,
		&course.HasCertificate,
		&course.MobileAccess,
		&course.AverageRating,
		&course.TotalEnrollments,
		&course.CreatedAt,
		&course.UpdatedAt,
	)

	// Handle nullable price
	if price.Valid {
		course.Price = price.Float64
	} else {
		course.Price = 0
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("course not found")
		}
		return nil, fmt.Errorf("failed to get course: %w", err)
	}

	// Convert duration from minutes to hours
	if course.EstimatedDurationHours > 0 {
		course.EstimatedDurationHours = course.EstimatedDurationHours / 60
	}

	// Set DifficultyLevel to same as Level for compatibility
	course.DifficultyLevel = course.Level

	// Get lectures for the course
	lectures, err := r.getCourseLectures(courseID)
	if err != nil {
		// Don't fail if lectures can't be loaded, just log and continue with empty array
		lectures = []model.Lecture{}
	}
	course.Lectures = lectures

	// Get resources for the course
	resources, err := r.getCourseResources(courseID)
	if err != nil {
		// Don't fail if resources can't be loaded, just log and continue with empty array
		resources = []model.Resource{}
	}
	course.Resources = resources

	return course, nil
}

// getCourseLectures retrieves all lectures for a course
func (r *DashboardRepository) getCourseLectures(courseID uuid.UUID) ([]model.Lecture, error) {
	query := `
		SELECT id, title, description, type, duration_minutes, is_free, order_number,
			   video_url, video_id, status, created_at, updated_at
		FROM lectures
		WHERE course_id = $1 AND deleted_at IS NULL
		ORDER BY order_number ASC
	`

	rows, err := r.db.Query(query, courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to query lectures: %w", err)
	}
	defer rows.Close()

	var lectures []model.Lecture
	for rows.Next() {
		var lecture model.Lecture
		var videoURL, videoID sql.NullString

		err := rows.Scan(
			&lecture.ID,
			&lecture.Title,
			&lecture.Description,
			&lecture.Type,
			&lecture.DurationMinutes,
			&lecture.IsFree,
			&lecture.OrderNumber,
			&videoURL,
			&videoID,
			&lecture.Status,
			&lecture.CreatedAt,
			&lecture.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan lecture: %w", err)
		}

		if videoURL.Valid {
			lecture.VideoURL = videoURL.String
		}
		if videoID.Valid {
			lecture.VideoID = videoID.String
		}

		lectures = append(lectures, lecture)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate lectures: %w", err)
	}

	return lectures, nil
}

// getCourseResources retrieves all resources/files for a course
func (r *DashboardRepository) getCourseResources(courseID uuid.UUID) ([]model.Resource, error) {
	// First, check if there's a course_files table or similar relationship
	// For now, we'll look for files associated with the course through metadata or a join table
	// Since we don't have a direct relationship, let's check files with course metadata
	query := `
		SELECT f.id, f.filename, f.original_filename, f.content_type, f.size_bytes,
			   CONCAT('https://s3.amazonaws.com/', f.bucket_name, '/', f.object_key) as download_url,
			   f.is_public, f.created_at
		FROM files f
		WHERE f.metadata->>'course_id' = $1 AND f.deleted_at IS NULL
		ORDER BY f.created_at DESC
	`

	rows, err := r.db.Query(query, courseID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to query course resources: %w", err)
	}
	defer rows.Close()

	var resources []model.Resource
	for rows.Next() {
		var resource model.Resource
		var id sql.NullString

		err := rows.Scan(
			&id,
			&resource.Filename,
			&resource.OriginalName,
			&resource.FileType,
			&resource.FileSize,
			&resource.DownloadURL,
			&resource.IsPublic,
			&resource.UploadedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan resource: %w", err)
		}

		if id.Valid {
			if parsedID, err := uuid.Parse(id.String); err == nil {
				resource.ID = parsedID
			}
		}

		resources = append(resources, resource)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate resources: %w", err)
	}

	return resources, nil
}

// CreateCourse creates a new course for an instructor
func (r *DashboardRepository) CreateCourse(instructorID uuid.UUID, req *model.CreateCourseRequest) (*model.Course, error) {
	courseID := uuid.New()
	now := time.Now()

	// Set default values
	status := "draft"
	if req.Status != "" {
		status = req.Status
	}

	currency := "USD"
	if req.Currency != "" {
		currency = req.Currency
	}

	query := `
		INSERT INTO courses (
			id, title, description, instructor_id, category, level, price, currency,
			thumbnail_url, status, tags, is_paid, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
		)
	`

	var price float64
	isPaid := false
	if req.Price != nil && *req.Price > 0 {
		price = *req.Price
		isPaid = true
	}

	_, err := r.db.Exec(query,
		courseID,
		req.Title,
		req.Description,
		instructorID,
		req.Category,
		req.Level,
		price,
		currency,
		req.ThumbnailURL,
		status,
		pq.Array(req.Tags),
		isPaid,
		now,
		now,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create course: %w", err)
	}

	// Return the created course
	return r.GetCourse(instructorID, courseID)
}

// UpdateCourse updates an existing course for an instructor
func (r *DashboardRepository) UpdateCourse(instructorID, courseID uuid.UUID, req *model.UpdateCourseRequest) (*model.Course, error) {
	// Check if course exists and belongs to instructor
	_, err := r.GetCourse(instructorID, courseID)
	if err != nil {
		return nil, err
	}

	// Build dynamic update query
	setParts := []string{}
	args := []interface{}{}
	argCount := 1

	if req.Title != "" {
		setParts = append(setParts, fmt.Sprintf("title = $%d", argCount))
		args = append(args, req.Title)
		argCount++
	}

	if req.Description != "" {
		setParts = append(setParts, fmt.Sprintf("description = $%d", argCount))
		args = append(args, req.Description)
		argCount++
	}

	if req.Category != "" {
		setParts = append(setParts, fmt.Sprintf("category = $%d", argCount))
		args = append(args, req.Category)
		argCount++
	}

	if req.Level != "" {
		setParts = append(setParts, fmt.Sprintf("level = $%d", argCount))
		args = append(args, req.Level)
		argCount++
	}

	if req.Price != nil {
		setParts = append(setParts, fmt.Sprintf("price = $%d", argCount))
		args = append(args, *req.Price)
		argCount++

		isPaid := *req.Price > 0
		setParts = append(setParts, fmt.Sprintf("is_paid = $%d", argCount))
		args = append(args, isPaid)
		argCount++
	}

	if req.Currency != "" {
		setParts = append(setParts, fmt.Sprintf("currency = $%d", argCount))
		args = append(args, req.Currency)
		argCount++
	}

	if req.ThumbnailURL != "" {
		setParts = append(setParts, fmt.Sprintf("thumbnail_url = $%d", argCount))
		args = append(args, req.ThumbnailURL)
		argCount++
	}

	if req.Status != "" {
		setParts = append(setParts, fmt.Sprintf("status = $%d", argCount))
		args = append(args, req.Status)
		argCount++
	}

	if req.Tags != nil {
		setParts = append(setParts, fmt.Sprintf("tags = $%d", argCount))
		args = append(args, pq.Array(req.Tags))
		argCount++
	}

	if len(setParts) == 0 {
		return r.GetCourse(instructorID, courseID)
	}

	// Add updated_at
	setParts = append(setParts, fmt.Sprintf("updated_at = $%d", argCount))
	args = append(args, time.Now())
	argCount++

	// Add WHERE clause parameters
	args = append(args, courseID, instructorID)

	query := fmt.Sprintf(`
		UPDATE courses
		SET %s
		WHERE id = $%d AND instructor_id = $%d AND deleted_at IS NULL
	`, strings.Join(setParts, ", "), argCount, argCount+1)

	_, err = r.db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update course: %w", err)
	}

	// Return the updated course
	return r.GetCourse(instructorID, courseID)
}

// DeleteCourse soft deletes a course for an instructor
func (r *DashboardRepository) DeleteCourse(instructorID, courseID uuid.UUID) error {
	// Check if course exists and belongs to instructor (without soft delete filter)
	verifyQuery := `
		SELECT id FROM courses
		WHERE id = $1 AND instructor_id = $2
	`

	var existingID uuid.UUID
	err := r.db.QueryRow(verifyQuery, courseID, instructorID).Scan(&existingID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("course not found or access denied")
		}
		return fmt.Errorf("failed to verify course ownership: %w", err)
	}

	// Perform soft delete - only update if not already deleted
	query := `
		UPDATE courses
		SET deleted_at = $1, updated_at = $1
		WHERE id = $2 AND instructor_id = $3 AND deleted_at IS NULL
	`

	_, err = r.db.Exec(query, time.Now(), courseID, instructorID)
	if err != nil {
		return fmt.Errorf("failed to delete course: %w", err)
	}

	// If no rows were affected, the course was already deleted - this is still a success
	// since the end result is the same (course is soft deleted)
	return nil
}