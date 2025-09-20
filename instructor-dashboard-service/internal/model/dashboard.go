package model

import (
	"time"
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

// JSONB represents a JSONB field for PostgreSQL
type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = make(JSONB)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot scan %T into JSONB", value)
	}

	return json.Unmarshal(bytes, j)
}

// DashboardOverview represents the main dashboard overview data
type DashboardOverview struct {
	InstructorID       uuid.UUID `json:"instructor_id"`
	TotalRevenue       float64   `json:"total_revenue"`
	MonthlyRevenue     float64   `json:"monthly_revenue"`
	TotalStudents      int       `json:"total_students"`
	ActiveStudents     int       `json:"active_students"`
	TotalCourses       int       `json:"total_courses"`
	PublishedCourses   int       `json:"published_courses"`
	AvgCourseRating    float64   `json:"avg_course_rating"`
	CompletionRate     float64   `json:"completion_rate"`
	EngagementScore    float64   `json:"engagement_score"`
	RecentActivity     []ActivityItem `json:"recent_activity"`
	QuickStats         QuickStats     `json:"quick_stats"`
	TopPerformingCourse Course        `json:"top_performing_course,omitempty"`
}

// QuickStats represents quick statistics for the dashboard
type QuickStats struct {
	NewEnrollmentsToday    int     `json:"new_enrollments_today"`
	RevenueToday          float64 `json:"revenue_today"`
	VideoWatchTimeToday   int     `json:"video_watch_time_today"`
	NewReviewsToday       int     `json:"new_reviews_today"`
	ForumPostsToday       int     `json:"forum_posts_today"`
	AIInteractionsToday   int     `json:"ai_interactions_today"`
}

// ActivityItem represents a recent activity item
type ActivityItem struct {
	ID          uuid.UUID   `json:"id"`
	Type        string      `json:"type"` // enrollment, completion, review, etc.
	Description string      `json:"description"`
	UserID      *uuid.UUID  `json:"user_id,omitempty"`
	UserName    string      `json:"user_name,omitempty"`
	CourseID    *uuid.UUID  `json:"course_id,omitempty"`
	CourseName  string      `json:"course_name,omitempty"`
	Timestamp   time.Time   `json:"timestamp"`
	Metadata    JSONB       `json:"metadata,omitempty"`
}

// InstructorDashboardSettings represents dashboard configuration
type InstructorDashboardSettings struct {
	ID                      uuid.UUID `json:"id" db:"id"`
	InstructorID           uuid.UUID `json:"instructor_id" db:"instructor_id"`
	DashboardLayout        JSONB     `json:"dashboard_layout" db:"dashboard_layout"`
	NotificationPreferences JSONB     `json:"notification_preferences" db:"notification_preferences"`
	AnalyticsPreferences   JSONB     `json:"analytics_preferences" db:"analytics_preferences"`
	DefaultCourseSettings  JSONB     `json:"default_course_settings" db:"default_course_settings"`
	AIAssistanceEnabled    bool      `json:"ai_assistance_enabled" db:"ai_assistance_enabled"`
	CreatedAt              time.Time `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time `json:"updated_at" db:"updated_at"`
}

// Course represents a course with instructor-specific data
type Course struct {
	ID                    uuid.UUID  `json:"id" db:"id"`
	Title                string     `json:"title" db:"title"`
	Description          string     `json:"description" db:"description"`
	CreatorID            uuid.UUID  `json:"creator_id" db:"creator_id"`
	Status               string     `json:"status" db:"status"`
	IsPaid               bool       `json:"is_paid" db:"is_paid"`
	Price                *float64   `json:"price" db:"price"`
	Currency             string     `json:"currency" db:"currency"`
	AverageRating        float64    `json:"average_rating" db:"average_rating"`
	TotalEnrollments     int        `json:"total_enrollments" db:"total_enrollments"`
	ActiveStudents       int        `json:"active_students"`
	CompletionRate       float64    `json:"completion_rate"`
	TotalRevenue         float64    `json:"total_revenue"`
	MonthlyRevenue       float64    `json:"monthly_revenue"`
	EngagementScore      float64    `json:"engagement_score"`
	LastActivityAt       *time.Time `json:"last_activity_at"`
	InstructorNotes      string     `json:"instructor_notes" db:"instructor_notes"`
	MarketingDescription string     `json:"marketing_description" db:"marketing_description"`
	CreatedAt            time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at" db:"updated_at"`
}

// Student represents a student with instructor-specific data
type Student struct {
	ID                 uuid.UUID `json:"id" db:"id"`
	Email              string    `json:"email" db:"email"`
	FirstName          string    `json:"first_name" db:"first_name"`
	LastName           string    `json:"last_name" db:"last_name"`
	ProfilePictureURL  string    `json:"profile_picture_url" db:"profile_picture_url"`
	EnrollmentDate     time.Time `json:"enrollment_date"`
	LastActivityAt     time.Time `json:"last_activity_at"`
	TotalCoursesEnrolled int     `json:"total_courses_enrolled"`
	CompletedCourses   int       `json:"completed_courses"`
	TotalWatchTime     int       `json:"total_watch_time"`
	EngagementScore    float64   `json:"engagement_score"`
	TotalSpent         float64   `json:"total_spent"`
	Status             string    `json:"status"` // active, inactive, suspended
}

// CourseOptimizationSuggestion represents AI-generated suggestions
type CourseOptimizationSuggestion struct {
	ID                   uuid.UUID `json:"id" db:"id"`
	CourseID            uuid.UUID `json:"course_id" db:"course_id"`
	InstructorID        uuid.UUID `json:"instructor_id" db:"instructor_id"`
	SuggestionType      string    `json:"suggestion_type" db:"suggestion_type"`
	Title               string    `json:"title" db:"title"`
	Description         string    `json:"description" db:"description"`
	PriorityScore       float64   `json:"priority_score" db:"priority_score"`
	ExpectedImpact      JSONB     `json:"expected_impact" db:"expected_impact"`
	ImplementationEffort string   `json:"implementation_effort" db:"implementation_effort"`
	Status              string    `json:"status" db:"status"`
	ImplementedAt       *time.Time `json:"implemented_at" db:"implemented_at"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
}

// InstructorPerformanceMetrics represents daily performance metrics
type InstructorPerformanceMetrics struct {
	ID                      uuid.UUID `json:"id" db:"id"`
	InstructorID           uuid.UUID `json:"instructor_id" db:"instructor_id"`
	MetricDate             time.Time `json:"metric_date" db:"metric_date"`
	TotalRevenue           float64   `json:"total_revenue" db:"total_revenue"`
	TotalEnrollments       int       `json:"total_enrollments" db:"total_enrollments"`
	AvgCourseRating        float64   `json:"avg_course_rating" db:"avg_course_rating"`
	TotalStudents          int       `json:"total_students" db:"total_students"`
	CourseCompletionRate   float64   `json:"course_completion_rate" db:"course_completion_rate"`
	StudentSatisfactionScore float64 `json:"student_satisfaction_score" db:"student_satisfaction_score"`
	ContentEngagementRate  float64   `json:"content_engagement_rate" db:"content_engagement_rate"`
	MetricsData           JSONB     `json:"metrics_data" db:"metrics_data"`
	CalculatedAt          time.Time `json:"calculated_at" db:"calculated_at"`
}

// TeamMember represents instructor team member
type TeamMember struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	InstructorID   uuid.UUID  `json:"instructor_id" db:"instructor_id"`
	TeamMemberID   uuid.UUID  `json:"team_member_id" db:"team_member_id"`
	Role           string     `json:"role" db:"role"`
	Permissions    JSONB      `json:"permissions" db:"permissions"`
	CourseAccess   []uuid.UUID `json:"course_access"`
	InvitedAt      time.Time  `json:"invited_at" db:"invited_at"`
	JoinedAt       *time.Time `json:"joined_at" db:"joined_at"`
	Status         string     `json:"status" db:"status"`
	InvitationToken string    `json:"invitation_token,omitempty" db:"invitation_token"`
	InvitedBy      uuid.UUID  `json:"invited_by" db:"invited_by"`

	// Joined fields from users table
	Email         string `json:"email"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	ProfilePictureURL string `json:"profile_picture_url"`
}

// Communication represents instructor-student communication
type Communication struct {
	ID                uuid.UUID  `json:"id" db:"id"`
	InstructorID      uuid.UUID  `json:"instructor_id" db:"instructor_id"`
	StudentID         uuid.UUID  `json:"student_id" db:"student_id"`
	CourseID          uuid.UUID  `json:"course_id" db:"course_id"`
	CommunicationType string     `json:"communication_type" db:"communication_type"`
	Subject           string     `json:"subject" db:"subject"`
	Message           string     `json:"message" db:"message"`
	Status            string     `json:"status" db:"status"`
	ScheduledAt       *time.Time `json:"scheduled_at" db:"scheduled_at"`
	SentAt            *time.Time `json:"sent_at" db:"sent_at"`
	ReadAt            *time.Time `json:"read_at" db:"read_at"`
	RepliedAt         *time.Time `json:"replied_at" db:"replied_at"`
	Metadata          JSONB      `json:"metadata" db:"metadata"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`

	// Joined fields
	StudentEmail     string `json:"student_email"`
	StudentFirstName string `json:"student_first_name"`
	StudentLastName  string `json:"student_last_name"`
	CourseTitle      string `json:"course_title"`
}

// NotificationSetting represents instructor notification preferences
type NotificationSetting struct {
	ID               uuid.UUID  `json:"id" db:"id"`
	InstructorID     uuid.UUID  `json:"instructor_id" db:"instructor_id"`
	NotificationType string     `json:"notification_type" db:"notification_type"`
	Enabled          bool       `json:"enabled" db:"enabled"`
	DeliveryMethod   string     `json:"delivery_method" db:"delivery_method"`
	ThresholdValue   *float64   `json:"threshold_value" db:"threshold_value"`
	Frequency        string     `json:"frequency" db:"frequency"`
	LastSentAt       *time.Time `json:"last_sent_at" db:"last_sent_at"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
}

// Request/Response Models

// UpdateDashboardSettingsRequest represents the request to update dashboard settings
type UpdateDashboardSettingsRequest struct {
	DashboardLayout        JSONB `json:"dashboard_layout,omitempty"`
	NotificationPreferences JSONB `json:"notification_preferences,omitempty"`
	AnalyticsPreferences   JSONB `json:"analytics_preferences,omitempty"`
	DefaultCourseSettings  JSONB `json:"default_course_settings,omitempty"`
	AIAssistanceEnabled    *bool `json:"ai_assistance_enabled,omitempty"`
}

// BulkCourseOperationRequest represents bulk operations on courses
type BulkCourseOperationRequest struct {
	Operation string      `json:"operation"` // publish, unpublish, delete, update_price
	CourseIDs []uuid.UUID `json:"course_ids"`
	Parameters JSONB      `json:"parameters,omitempty"`
}

// SendBroadcastRequest represents a broadcast message request
type SendBroadcastRequest struct {
	CourseIDs         []uuid.UUID `json:"course_ids,omitempty"`
	StudentIDs        []uuid.UUID `json:"student_ids,omitempty"`
	Subject           string      `json:"subject"`
	Message           string      `json:"message"`
	CommunicationType string      `json:"communication_type"`
	ScheduledAt       *time.Time  `json:"scheduled_at,omitempty"`
}

// InviteTeamMemberRequest represents team member invitation
type InviteTeamMemberRequest struct {
	Email        string      `json:"email"`
	Role         string      `json:"role"`
	Permissions  JSONB       `json:"permissions,omitempty"`
	CourseAccess []uuid.UUID `json:"course_access,omitempty"`
}

// UpdateTeamMemberRequest represents team member update
type UpdateTeamMemberRequest struct {
	Role         string      `json:"role,omitempty"`
	Permissions  JSONB       `json:"permissions,omitempty"`
	CourseAccess []uuid.UUID `json:"course_access,omitempty"`
	Status       string      `json:"status,omitempty"`
}

// AutomatedMessageRule represents rules for automated messages
type AutomatedMessageRule struct {
	ID                uuid.UUID  `json:"id" db:"id"`
	InstructorID      uuid.UUID  `json:"instructor_id" db:"instructor_id"`
	RuleName          string     `json:"rule_name" db:"rule_name"`
	TriggerType       string     `json:"trigger_type" db:"trigger_type"` // enrollment, completion, inactivity, milestone
	TriggerConditions JSONB      `json:"trigger_conditions" db:"trigger_conditions"`
	MessageTemplate   string     `json:"message_template" db:"message_template"`
	Subject           string     `json:"subject" db:"subject"`
	DelayHours        int        `json:"delay_hours" db:"delay_hours"`
	IsActive          bool       `json:"is_active" db:"is_active"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
}

// Pagination represents pagination parameters
type Pagination struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"page_size" form:"page_size"`
	Total    int `json:"total"`
	Pages    int `json:"pages"`
}

// Filter represents common filter parameters
type Filter struct {
	StartDate *time.Time `json:"start_date" form:"start_date"`
	EndDate   *time.Time `json:"end_date" form:"end_date"`
	Status    string     `json:"status" form:"status"`
	CourseID  *uuid.UUID `json:"course_id" form:"course_id"`
}

// MessageTemplate represents a reusable message template
type MessageTemplate struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Subject  string `json:"subject"`
	Template string `json:"template"`
	Category string `json:"category"`
}

// PersonalMessageRequest represents a personal message request
type PersonalMessageRequest struct {
	CourseID *uuid.UUID `json:"course_id,omitempty"`
	Subject  string     `json:"subject" binding:"required"`
	Message  string     `json:"message" binding:"required"`
	Metadata JSONB      `json:"metadata,omitempty"`
}

// Notification represents a notification for an instructor
type Notification struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	InstructorID uuid.UUID `json:"instructor_id" db:"instructor_id"`
	Type        string     `json:"type" db:"type"` // new_enrollment, course_completion, review, revenue_milestone, etc.
	Title       string     `json:"title" db:"title"`
	Message     string     `json:"message" db:"message"`
	IsRead      bool       `json:"is_read" db:"is_read"`
	Priority    string     `json:"priority" db:"priority"` // low, medium, high, urgent
	ActionURL   *string    `json:"action_url" db:"action_url"`
	Metadata    JSONB      `json:"metadata" db:"metadata"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	ReadAt      *time.Time `json:"read_at" db:"read_at"`
}