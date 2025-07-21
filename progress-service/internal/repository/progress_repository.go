package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/study-platform/progress-service/internal/model"
	"github.com/study-platform/pkg/database"
)

type ProgressRepository struct {
	db *database.DB
}

func NewProgressRepository(db *database.DB) *ProgressRepository {
	return &ProgressRepository{db: db}
}

// Progress tracking methods
func (r *ProgressRepository) CreateProgress(progress *model.UserProgress) error {
	query := `
		INSERT INTO progress (id, user_id, course_id, lecture_id, progress_percentage, watch_time_seconds, is_completed, last_accessed, completed_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := r.db.Exec(query, progress.ID, progress.UserID, progress.CourseID, progress.LectureID, 
		progress.ProgressPercentage, progress.WatchTimeSeconds, progress.IsCompleted, progress.LastAccessed, 
		progress.CompletedAt, progress.CreatedAt, progress.UpdatedAt)
	return err
}

func (r *ProgressRepository) UpdateProgress(progress *model.UserProgress) error {
	query := `
		UPDATE progress 
		SET progress_percentage = $4, watch_time_seconds = $5, is_completed = $6, last_accessed = $7, completed_at = $8, updated_at = $9
		WHERE user_id = $1 AND course_id = $2 AND lecture_id = $3
	`
	_, err := r.db.Exec(query, progress.UserID, progress.CourseID, progress.LectureID, 
		progress.ProgressPercentage, progress.WatchTimeSeconds, progress.IsCompleted, progress.LastAccessed, 
		progress.CompletedAt, progress.UpdatedAt)
	return err
}

func (r *ProgressRepository) GetProgress(userID, courseID, lectureID uuid.UUID) (*model.UserProgress, error) {
	query := `
		SELECT id, user_id, course_id, lecture_id, progress_percentage, watch_time_seconds, is_completed, last_accessed, completed_at, created_at, updated_at
		FROM progress
		WHERE user_id = $1 AND course_id = $2 AND lecture_id = $3
	`
	row := r.db.QueryRow(query, userID, courseID, lectureID)
	
	progress := &model.UserProgress{}
	err := row.Scan(&progress.ID, &progress.UserID, &progress.CourseID, &progress.LectureID, 
		&progress.ProgressPercentage, &progress.WatchTimeSeconds, &progress.IsCompleted, 
		&progress.LastAccessed, &progress.CompletedAt, &progress.CreatedAt, &progress.UpdatedAt)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("progress not found")
		}
		return nil, err
	}
	
	return progress, nil
}

func (r *ProgressRepository) GetUserProgress(userID, courseID uuid.UUID) ([]*model.UserProgress, error) {
	query := `
		SELECT id, user_id, course_id, lecture_id, progress_percentage, watch_time_seconds, is_completed, last_accessed, completed_at, created_at, updated_at
		FROM progress
		WHERE user_id = $1 AND course_id = $2
		ORDER BY created_at ASC
	`
	rows, err := r.db.Query(query, userID, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var progressList []*model.UserProgress
	for rows.Next() {
		progress := &model.UserProgress{}
		err := rows.Scan(&progress.ID, &progress.UserID, &progress.CourseID, &progress.LectureID, 
			&progress.ProgressPercentage, &progress.WatchTimeSeconds, &progress.IsCompleted, 
			&progress.LastAccessed, &progress.CompletedAt, &progress.CreatedAt, &progress.UpdatedAt)
		if err != nil {
			return nil, err
		}
		progressList = append(progressList, progress)
	}
	
	return progressList, nil
}

func (r *ProgressRepository) GetCourseProgress(courseID uuid.UUID, page, pageSize int) ([]*model.UserProgress, int, error) {
	// Count total records
	countQuery := `SELECT COUNT(*) FROM progress WHERE course_id = $1`
	var totalCount int
	err := r.db.QueryRow(countQuery, courseID).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	query := `
		SELECT id, user_id, course_id, lecture_id, progress_percentage, watch_time_seconds, is_completed, last_accessed, completed_at, created_at, updated_at
		FROM progress
		WHERE course_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(query, courseID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var progressList []*model.UserProgress
	for rows.Next() {
		progress := &model.UserProgress{}
		err := rows.Scan(&progress.ID, &progress.UserID, &progress.CourseID, &progress.LectureID, 
			&progress.ProgressPercentage, &progress.WatchTimeSeconds, &progress.IsCompleted, 
			&progress.LastAccessed, &progress.CompletedAt, &progress.CreatedAt, &progress.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		progressList = append(progressList, progress)
	}
	
	return progressList, totalCount, nil
}

// Enrollment methods
func (r *ProgressRepository) CreateEnrollment(enrollment *model.Enrollment) error {
	query := `
		INSERT INTO enrollments (id, user_id, course_id, status, progress_percentage, completed_lectures, total_lectures, total_watch_time_seconds, enrolled_at, completed_at, last_accessed, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`
	_, err := r.db.Exec(query, enrollment.ID, enrollment.UserID, enrollment.CourseID, enrollment.Status, 
		enrollment.ProgressPercentage, enrollment.CompletedLectures, enrollment.TotalLectures, 
		enrollment.TotalWatchTimeSeconds, enrollment.EnrolledAt, enrollment.CompletedAt, 
		enrollment.LastAccessed, enrollment.CreatedAt, enrollment.UpdatedAt)
	return err
}

func (r *ProgressRepository) GetEnrollment(userID, courseID uuid.UUID) (*model.Enrollment, error) {
	query := `
		SELECT id, user_id, course_id, status, progress_percentage, completed_lectures, total_lectures, total_watch_time_seconds, enrolled_at, completed_at, last_accessed, created_at, updated_at
		FROM enrollments
		WHERE user_id = $1 AND course_id = $2
	`
	row := r.db.QueryRow(query, userID, courseID)
	
	enrollment := &model.Enrollment{}
	err := row.Scan(&enrollment.ID, &enrollment.UserID, &enrollment.CourseID, &enrollment.Status, 
		&enrollment.ProgressPercentage, &enrollment.CompletedLectures, &enrollment.TotalLectures, 
		&enrollment.TotalWatchTimeSeconds, &enrollment.EnrolledAt, &enrollment.CompletedAt, 
		&enrollment.LastAccessed, &enrollment.CreatedAt, &enrollment.UpdatedAt)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("enrollment not found")
		}
		return nil, err
	}
	
	return enrollment, nil
}

func (r *ProgressRepository) ListEnrollments(userID uuid.UUID, status string, page, pageSize int) ([]*model.Enrollment, int, error) {
	// Build query with optional status filter
	countQuery := `SELECT COUNT(*) FROM enrollments WHERE user_id = $1`
	query := `
		SELECT id, user_id, course_id, status, progress_percentage, completed_lectures, total_lectures, total_watch_time_seconds, enrolled_at, completed_at, last_accessed, created_at, updated_at
		FROM enrollments
		WHERE user_id = $1
	`
	
	var countArgs []interface{}
	var queryArgs []interface{}
	countArgs = append(countArgs, userID)
	queryArgs = append(queryArgs, userID)
	
	if status != "" {
		countQuery += " AND status = $2"
		query += " AND status = $2"
		countArgs = append(countArgs, status)
		queryArgs = append(queryArgs, status)
	}
	
	// Count total records
	var totalCount int
	err := r.db.QueryRow(countQuery, countArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// Add pagination
	offset := (page - 1) * pageSize
	query += " ORDER BY enrolled_at DESC LIMIT $" + fmt.Sprintf("%d", len(queryArgs)+1) + " OFFSET $" + fmt.Sprintf("%d", len(queryArgs)+2)
	queryArgs = append(queryArgs, pageSize, offset)

	rows, err := r.db.Query(query, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var enrollments []*model.Enrollment
	for rows.Next() {
		enrollment := &model.Enrollment{}
		err := rows.Scan(&enrollment.ID, &enrollment.UserID, &enrollment.CourseID, &enrollment.Status, 
			&enrollment.ProgressPercentage, &enrollment.CompletedLectures, &enrollment.TotalLectures, 
			&enrollment.TotalWatchTimeSeconds, &enrollment.EnrolledAt, &enrollment.CompletedAt, 
			&enrollment.LastAccessed, &enrollment.CreatedAt, &enrollment.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		enrollments = append(enrollments, enrollment)
	}
	
	return enrollments, totalCount, nil
}

func (r *ProgressRepository) UpdateEnrollmentStatus(userID, courseID uuid.UUID, status string) error {
	now := time.Now()
	query := `
		UPDATE enrollments 
		SET status = $3, updated_at = $4
		WHERE user_id = $1 AND course_id = $2
	`
	_, err := r.db.Exec(query, userID, courseID, status, now)
	return err
}

func (r *ProgressRepository) UpdateEnrollmentProgress(userID, courseID uuid.UUID, progressPercentage float64, completedLectures, totalWatchTime int32) error {
	now := time.Now()
	var completedAt *time.Time
	if progressPercentage >= 100.0 {
		completedAt = &now
	}
	
	query := `
		UPDATE enrollments 
		SET progress_percentage = $3, completed_lectures = $4, total_watch_time_seconds = $5, completed_at = $6, last_accessed = $7, updated_at = $8
		WHERE user_id = $1 AND course_id = $2
	`
	_, err := r.db.Exec(query, userID, courseID, progressPercentage, completedLectures, totalWatchTime, completedAt, now, now)
	return err
}

// Lecture progress methods
func (r *ProgressRepository) GetLectureProgress(userID, courseID uuid.UUID) ([]*model.LectureProgress, error) {
	query := `
		SELECT 
			l.id as lecture_id,
			l.title,
			l.order_number,
			COALESCE(p.progress_percentage, 0) as progress_percentage,
			COALESCE(p.watch_time_seconds, 0) as watch_time_seconds,
			COALESCE(p.is_completed, false) as is_completed,
			p.last_accessed,
			p.completed_at
		FROM lectures l
		LEFT JOIN progress p ON l.id = p.lecture_id AND p.user_id = $1
		WHERE l.course_id = $2
		ORDER BY l.order_number ASC
	`
	rows, err := r.db.Query(query, userID, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lectureProgress []*model.LectureProgress
	for rows.Next() {
		lp := &model.LectureProgress{}
		err := rows.Scan(&lp.LectureID, &lp.Title, &lp.OrderNumber, &lp.ProgressPercentage, 
			&lp.WatchTimeSeconds, &lp.IsCompleted, &lp.LastAccessed, &lp.CompletedAt)
		if err != nil {
			return nil, err
		}
		lectureProgress = append(lectureProgress, lp)
	}
	
	return lectureProgress, nil
}

func (r *ProgressRepository) GetCourseCompletion(userID, courseID uuid.UUID) (*model.CourseCompletion, error) {
	// Get course info and overall statistics
	query := `
		SELECT 
			c.id as course_id,
			c.title as course_title,
			COUNT(l.id) as total_lectures,
			COUNT(CASE WHEN p.is_completed = true THEN 1 END) as completed_lectures,
			COALESCE(SUM(p.watch_time_seconds), 0) as total_watch_time,
			COALESCE(AVG(p.progress_percentage), 0) as avg_progress,
			MIN(p.created_at) as started_at,
			MAX(CASE WHEN p.is_completed = true THEN p.completed_at END) as completed_at,
			MAX(p.last_accessed) as last_accessed
		FROM courses c
		LEFT JOIN lectures l ON c.id = l.course_id
		LEFT JOIN progress p ON l.id = p.lecture_id AND p.user_id = $1
		WHERE c.id = $2
		GROUP BY c.id, c.title
	`
	row := r.db.QueryRow(query, userID, courseID)
	
	completion := &model.CourseCompletion{UserID: userID}
	var totalLectures, completedLectures int32
	var totalWatchTime int32
	var avgProgress float64
	
	err := row.Scan(&completion.CourseID, &completion.CourseTitle, &totalLectures, &completedLectures, 
		&totalWatchTime, &avgProgress, &completion.StartedAt, &completion.CompletedAt, &completion.LastAccessed)
	if err != nil {
		return nil, err
	}
	
	completion.TotalLectures = totalLectures
	completion.CompletedLectures = completedLectures
	completion.TotalWatchTimeSeconds = totalWatchTime
	if totalLectures > 0 {
		completion.CompletionPercentage = (float64(completedLectures) / float64(totalLectures)) * 100.0
	}
	
	// Get lecture progress details
	lectureProgress, err := r.GetLectureProgress(userID, courseID)
	if err != nil {
		return nil, err
	}
	completion.LectureProgress = lectureProgress
	
	return completion, nil
}

// Analytics methods
func (r *ProgressRepository) GetUserAnalytics(userID uuid.UUID) (*model.UserAnalytics, error) {
	query := `
		SELECT 
			$1 as user_id,
			COUNT(DISTINCT e.course_id) as total_courses_enrolled,
			COUNT(DISTINCT CASE WHEN e.status = 'completed' THEN e.course_id END) as total_courses_completed,
			COUNT(CASE WHEN p.is_completed = true THEN 1 END) as total_lectures_completed,
			COALESCE(SUM(p.watch_time_seconds), 0) as total_watch_time,
			COALESCE(AVG(e.progress_percentage), 0) as avg_progress,
			COUNT(DISTINCT CASE WHEN e.status = 'enrolled' THEN e.course_id END) as courses_in_progress,
			COALESCE(MAX(p.last_accessed), MAX(e.last_accessed)) as last_activity
		FROM enrollments e
		LEFT JOIN progress p ON e.user_id = p.user_id AND e.course_id = p.course_id
		WHERE e.user_id = $1
	`
	row := r.db.QueryRow(query, userID)
	
	analytics := &model.UserAnalytics{}
	err := row.Scan(&analytics.UserID, &analytics.TotalCoursesEnrolled, &analytics.TotalCoursesCompleted, 
		&analytics.TotalLecturesCompleted, &analytics.TotalWatchTimeSeconds, &analytics.AverageProgressPercentage, 
		&analytics.CoursesInProgress, &analytics.LastActivity)
	if err != nil {
		return nil, err
	}
	
	// Set default values for fields not calculated in the query
	analytics.MostActiveDay = "Monday" // Placeholder
	analytics.StreakDays = 0           // Placeholder
	
	return analytics, nil
}

func (r *ProgressRepository) GetCourseAnalytics(courseID uuid.UUID) (*model.CourseAnalytics, error) {
	query := `
		SELECT 
			$1 as course_id,
			COUNT(DISTINCT e.user_id) as total_enrollments,
			COUNT(DISTINCT CASE WHEN e.status = 'completed' THEN e.user_id END) as total_completions,
			CASE 
				WHEN COUNT(DISTINCT e.user_id) > 0 THEN 
					(COUNT(DISTINCT CASE WHEN e.status = 'completed' THEN e.user_id END)::float / COUNT(DISTINCT e.user_id)::float) * 100.0
				ELSE 0
			END as completion_rate,
			COALESCE(AVG(e.progress_percentage), 0) as avg_progress,
			COUNT(DISTINCT CASE WHEN e.status = 'enrolled' THEN e.user_id END) as active_students,
			COALESCE(SUM(p.watch_time_seconds), 0) as total_watch_time,
			CASE 
				WHEN COUNT(DISTINCT e.user_id) > 0 THEN 
					COALESCE(SUM(p.watch_time_seconds), 0)::float / COUNT(DISTINCT e.user_id)::float
				ELSE 0
			END as avg_watch_time_per_student
		FROM enrollments e
		LEFT JOIN progress p ON e.user_id = p.user_id AND e.course_id = p.course_id
		WHERE e.course_id = $1
	`
	row := r.db.QueryRow(query, courseID)
	
	analytics := &model.CourseAnalytics{}
	err := row.Scan(&analytics.CourseID, &analytics.TotalEnrollments, &analytics.TotalCompletions, 
		&analytics.CompletionRate, &analytics.AverageProgressPercentage, &analytics.ActiveStudents, 
		&analytics.TotalWatchTimeSeconds, &analytics.AverageWatchTimePerStudent)
	if err != nil {
		return nil, err
	}
	
	// Set default values for fields not calculated in the query
	analytics.MostPopularLecture = "N/A" // Placeholder
	analytics.DropoutRate = 0.0          // Placeholder
	
	return analytics, nil
}

func (r *ProgressRepository) GetInstructorAnalytics(instructorID uuid.UUID) (*model.InstructorAnalytics, error) {
	query := `
		SELECT 
			$1 as instructor_id,
			COUNT(DISTINCT c.id) as total_courses,
			COUNT(DISTINCT e.user_id) as total_students,
			COUNT(DISTINCT CASE WHEN e.status = 'completed' THEN e.user_id END) as total_completions,
			CASE 
				WHEN COUNT(DISTINCT e.user_id) > 0 THEN 
					(COUNT(DISTINCT CASE WHEN e.status = 'completed' THEN e.user_id END)::float / COUNT(DISTINCT e.user_id)::float) * 100.0
				ELSE 0
			END as avg_completion_rate,
			COALESCE(SUM(p.watch_time_seconds), 0) as total_watch_time,
			COALESCE(AVG(c.rating), 0) as avg_course_rating,
			COUNT(DISTINCT CASE WHEN c.status = 'published' THEN c.id END) as active_courses
		FROM courses c
		LEFT JOIN enrollments e ON c.id = e.course_id
		LEFT JOIN progress p ON e.user_id = p.user_id AND e.course_id = p.course_id
		WHERE c.instructor_id = $1
	`
	row := r.db.QueryRow(query, instructorID)
	
	analytics := &model.InstructorAnalytics{}
	err := row.Scan(&analytics.InstructorID, &analytics.TotalCourses, &analytics.TotalStudents, 
		&analytics.TotalCompletions, &analytics.AverageCompletionRate, &analytics.TotalWatchTimeSeconds, 
		&analytics.AverageCourseRating, &analytics.ActiveCourses)
	if err != nil {
		return nil, err
	}
	
	// Set default values for fields not calculated in the query
	analytics.BestPerformingCourse = "N/A" // Placeholder
	analytics.TotalRevenueCents = 0        // Placeholder
	
	return analytics, nil
}