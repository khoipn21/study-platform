package model

import (
	"time"

	"github.com/google/uuid"
)

// UserProgress represents progress tracking for a user on a specific lecture
type UserProgress struct {
	ID                 uuid.UUID `json:"id" db:"id"`
	UserID             uuid.UUID `json:"user_id" db:"user_id"`
	CourseID           uuid.UUID `json:"course_id" db:"course_id"`
	LectureID          uuid.UUID `json:"lecture_id" db:"lecture_id"`
	ProgressPercentage float64   `json:"progress_percentage" db:"progress_percentage"`
	WatchTimeSeconds   int32     `json:"watch_time_seconds" db:"watch_time_seconds"`
	IsCompleted        bool      `json:"is_completed" db:"is_completed"`
	LastAccessed       time.Time `json:"last_accessed" db:"last_accessed"`
	CompletedAt        *time.Time `json:"completed_at" db:"completed_at"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

// Enrollment represents a user's enrollment in a course
type Enrollment struct {
	ID                   uuid.UUID `json:"id" db:"id"`
	UserID               uuid.UUID `json:"user_id" db:"user_id"`
	CourseID             uuid.UUID `json:"course_id" db:"course_id"`
	Status               string    `json:"status" db:"status"`
	ProgressPercentage   float64   `json:"progress_percentage" db:"progress_percentage"`
	CompletedLectures    int32     `json:"completed_lectures" db:"completed_lectures"`
	TotalLectures        int32     `json:"total_lectures" db:"total_lectures"`
	TotalWatchTimeSeconds int32     `json:"total_watch_time_seconds" db:"total_watch_time_seconds"`
	EnrolledAt           time.Time `json:"enrolled_at" db:"enrolled_at"`
	CompletedAt          *time.Time `json:"completed_at" db:"completed_at"`
	LastAccessed         *time.Time `json:"last_accessed" db:"last_accessed"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`
}

// LectureProgress represents progress information for a lecture
type LectureProgress struct {
	LectureID          uuid.UUID  `json:"lecture_id" db:"lecture_id"`
	Title              string     `json:"title" db:"title"`
	OrderNumber        int32      `json:"order_number" db:"order_number"`
	ProgressPercentage float64    `json:"progress_percentage" db:"progress_percentage"`
	WatchTimeSeconds   int32      `json:"watch_time_seconds" db:"watch_time_seconds"`
	IsCompleted        bool       `json:"is_completed" db:"is_completed"`
	LastAccessed       *time.Time `json:"last_accessed" db:"last_accessed"`
	CompletedAt        *time.Time `json:"completed_at" db:"completed_at"`
}

// CourseCompletion represents overall course completion status
type CourseCompletion struct {
	CourseID             uuid.UUID         `json:"course_id" db:"course_id"`
	CourseTitle          string            `json:"course_title" db:"course_title"`
	UserID               uuid.UUID         `json:"user_id" db:"user_id"`
	CompletionPercentage float64           `json:"completion_percentage" db:"completion_percentage"`
	CompletedLectures    int32             `json:"completed_lectures" db:"completed_lectures"`
	TotalLectures        int32             `json:"total_lectures" db:"total_lectures"`
	TotalWatchTimeSeconds int32             `json:"total_watch_time_seconds" db:"total_watch_time_seconds"`
	LectureProgress      []*LectureProgress `json:"lecture_progress"`
	StartedAt            *time.Time        `json:"started_at" db:"started_at"`
	CompletedAt          *time.Time        `json:"completed_at" db:"completed_at"`
	LastAccessed         *time.Time        `json:"last_accessed" db:"last_accessed"`
}

// UserAnalytics represents user learning analytics
type UserAnalytics struct {
	UserID                     uuid.UUID  `json:"user_id" db:"user_id"`
	TotalCoursesEnrolled       int32      `json:"total_courses_enrolled" db:"total_courses_enrolled"`
	TotalCoursesCompleted      int32      `json:"total_courses_completed" db:"total_courses_completed"`
	TotalLecturesCompleted     int32      `json:"total_lectures_completed" db:"total_lectures_completed"`
	TotalWatchTimeSeconds      int32      `json:"total_watch_time_seconds" db:"total_watch_time_seconds"`
	AverageProgressPercentage  float64    `json:"average_progress_percentage" db:"average_progress_percentage"`
	CoursesInProgress          int32      `json:"courses_in_progress" db:"courses_in_progress"`
	MostActiveDay              string     `json:"most_active_day" db:"most_active_day"`
	StreakDays                 int32      `json:"streak_days" db:"streak_days"`
	LastActivity               *time.Time `json:"last_activity" db:"last_activity"`
}

// CourseAnalytics represents course performance analytics
type CourseAnalytics struct {
	CourseID                  uuid.UUID `json:"course_id" db:"course_id"`
	TotalEnrollments          int32     `json:"total_enrollments" db:"total_enrollments"`
	TotalCompletions          int32     `json:"total_completions" db:"total_completions"`
	CompletionRate            float64   `json:"completion_rate" db:"completion_rate"`
	AverageProgressPercentage float64   `json:"average_progress_percentage" db:"average_progress_percentage"`
	ActiveStudents            int32     `json:"active_students" db:"active_students"`
	TotalWatchTimeSeconds     int32     `json:"total_watch_time_seconds" db:"total_watch_time_seconds"`
	AverageWatchTimePerStudent float64   `json:"average_watch_time_per_student" db:"average_watch_time_per_student"`
	MostPopularLecture        string    `json:"most_popular_lecture" db:"most_popular_lecture"`
	DropoutRate               float64   `json:"dropout_rate" db:"dropout_rate"`
}

// InstructorAnalytics represents instructor performance analytics
type InstructorAnalytics struct {
	InstructorID           uuid.UUID `json:"instructor_id" db:"instructor_id"`
	TotalCourses           int32     `json:"total_courses" db:"total_courses"`
	TotalStudents          int32     `json:"total_students" db:"total_students"`
	TotalCompletions       int32     `json:"total_completions" db:"total_completions"`
	AverageCompletionRate  float64   `json:"average_completion_rate" db:"average_completion_rate"`
	TotalWatchTimeSeconds  int32     `json:"total_watch_time_seconds" db:"total_watch_time_seconds"`
	AverageCourseRating    float64   `json:"average_course_rating" db:"average_course_rating"`
	ActiveCourses          int32     `json:"active_courses" db:"active_courses"`
	BestPerformingCourse   string    `json:"best_performing_course" db:"best_performing_course"`
	TotalRevenueCents      int32     `json:"total_revenue_cents" db:"total_revenue_cents"`
}