package model

import (
	"time"
	"github.com/google/uuid"
)

// RevenueAnalytics represents revenue analytics data
type RevenueAnalytics struct {
	Period          string                 `json:"period"` // daily, weekly, monthly, yearly
	StartDate       time.Time              `json:"start_date"`
	EndDate         time.Time              `json:"end_date"`
	TotalRevenue    float64                `json:"total_revenue"`
	GrossRevenue    float64                `json:"gross_revenue"`
	NetRevenue      float64                `json:"net_revenue"`
	RefundAmount    float64                `json:"refund_amount"`
	TotalSales      int                    `json:"total_sales"`
	UniqueBuyers    int                    `json:"unique_buyers"`
	AvgOrderValue   float64                `json:"avg_order_value"`
	ConversionRate  float64                `json:"conversion_rate"`
	RevenueGrowth   float64                `json:"revenue_growth"`
	DailyBreakdown  []DailyRevenueMetric   `json:"daily_breakdown"`
	CourseBreakdown []CourseRevenueMetric  `json:"course_breakdown"`
	TopCourses      []CourseRevenueMetric  `json:"top_courses"`
}

// DailyRevenueMetric represents daily revenue data
type DailyRevenueMetric struct {
	Date         time.Time `json:"date"`
	Revenue      float64   `json:"revenue"`
	Sales        int       `json:"sales"`
	Refunds      float64   `json:"refunds"`
	NetRevenue   float64   `json:"net_revenue"`
	NewCustomers int       `json:"new_customers"`
}

// CourseRevenueMetric represents revenue data per course
type CourseRevenueMetric struct {
	CourseID        uuid.UUID `json:"course_id"`
	CourseTitle     string    `json:"course_title"`
	TotalRevenue    float64   `json:"total_revenue"`
	TotalSales      int       `json:"total_sales"`
	RefundAmount    float64   `json:"refund_amount"`
	NetRevenue      float64   `json:"net_revenue"`
	AvgSalePrice    float64   `json:"avg_sale_price"`
	ConversionRate  float64   `json:"conversion_rate"`
	RevenueGrowth   float64   `json:"revenue_growth"`
}

// EngagementAnalytics represents student engagement analytics
type EngagementAnalytics struct {
	Period                  string                    `json:"period"`
	StartDate               time.Time                 `json:"start_date"`
	EndDate                 time.Time                 `json:"end_date"`
	TotalStudents           int                       `json:"total_students"`
	ActiveStudents          int                       `json:"active_students"`
	EngagementRate          float64                   `json:"engagement_rate"`
	AvgWatchTime            int                       `json:"avg_watch_time"`
	CompletionRate          float64                   `json:"completion_rate"`
	DropOffRate             float64                   `json:"drop_off_rate"`
	DiscussionParticipation float64                   `json:"discussion_participation"`
	QuizParticipation       float64                   `json:"quiz_participation"`
	AIInteractions          int                       `json:"ai_interactions"`
	DailyEngagement         []DailyEngagementMetric   `json:"daily_engagement"`
	CourseEngagement        []CourseEngagementMetric  `json:"course_engagement"`
	VideoEngagement         []VideoEngagementMetric   `json:"video_engagement"`
}

// DailyEngagementMetric represents daily engagement data
type DailyEngagementMetric struct {
	Date                time.Time `json:"date"`
	ActiveStudents      int       `json:"active_students"`
	TotalWatchTime      int       `json:"total_watch_time"`
	VideoCompletions    int       `json:"video_completions"`
	QuizCompletions     int       `json:"quiz_completions"`
	ForumPosts          int       `json:"forum_posts"`
	AIInteractions      int       `json:"ai_interactions"`
	EngagementScore     float64   `json:"engagement_score"`
}

// CourseEngagementMetric represents engagement data per course
type CourseEngagementMetric struct {
	CourseID            uuid.UUID `json:"course_id"`
	CourseTitle         string    `json:"course_title"`
	TotalStudents       int       `json:"total_students"`
	ActiveStudents      int       `json:"active_students"`
	CompletionRate      float64   `json:"completion_rate"`
	AvgProgress         float64   `json:"avg_progress"`
	AvgWatchTime        int       `json:"avg_watch_time"`
	EngagementScore     float64   `json:"engagement_score"`
	DropOffRate         float64   `json:"drop_off_rate"`
	StudentSatisfaction float64   `json:"student_satisfaction"`
}

// VideoEngagementMetric represents engagement data per video
type VideoEngagementMetric struct {
	VideoID           uuid.UUID `json:"video_id"`
	VideoTitle        string    `json:"video_title"`
	CourseID          uuid.UUID `json:"course_id"`
	CourseTitle       string    `json:"course_title"`
	TotalViews        int       `json:"total_views"`
	UniqueViewers     int       `json:"unique_viewers"`
	AvgWatchTime      int       `json:"avg_watch_time"`
	CompletionRate    float64   `json:"completion_rate"`
	ReplayRate        float64   `json:"replay_rate"`
	DropOffPoints     []DropOffPoint `json:"drop_off_points"`
	EngagementScore   float64   `json:"engagement_score"`
	AIQuestions       int       `json:"ai_questions"`
	BookmarksCreated  int       `json:"bookmarks_created"`
}

// DropOffPoint represents points where students commonly drop off
type DropOffPoint struct {
	TimestampSeconds int     `json:"timestamp_seconds"`
	DropOffRate      float64 `json:"drop_off_rate"`
	StudentCount     int     `json:"student_count"`
}

// StudentAnalytics represents student analytics data
type StudentAnalytics struct {
	Period                string                `json:"period"`
	StartDate             time.Time             `json:"start_date"`
	EndDate               time.Time             `json:"end_date"`
	TotalStudents         int                   `json:"total_students"`
	NewStudents           int                   `json:"new_students"`
	ActiveStudents        int                   `json:"active_students"`
	RetentionRate         float64               `json:"retention_rate"`
	ChurnRate             float64               `json:"churn_rate"`
	AvgStudentLifetime    int                   `json:"avg_student_lifetime"`
	AvgCoursesPerStudent  float64               `json:"avg_courses_per_student"`
	StudentSatisfaction   float64               `json:"student_satisfaction"`
	StudentDemographics   StudentDemographics   `json:"student_demographics"`
	StudentProgress       []StudentProgress     `json:"student_progress"`
	TopStudents           []TopStudent          `json:"top_students"`
	AtRiskStudents        []AtRiskStudent       `json:"at_risk_students"`
}

// StudentDemographics represents demographic data
type StudentDemographics struct {
	CountryDistribution   []CountryMetric    `json:"country_distribution"`
	AgeDistribution       []AgeGroupMetric   `json:"age_distribution"`
	DeviceDistribution    []DeviceMetric     `json:"device_distribution"`
	LanguageDistribution  []LanguageMetric   `json:"language_distribution"`
}

// CountryMetric represents country-based metrics
type CountryMetric struct {
	Country      string  `json:"country"`
	StudentCount int     `json:"student_count"`
	Percentage   float64 `json:"percentage"`
}

// AgeGroupMetric represents age group metrics
type AgeGroupMetric struct {
	AgeGroup     string  `json:"age_group"`
	StudentCount int     `json:"student_count"`
	Percentage   float64 `json:"percentage"`
}

// DeviceMetric represents device usage metrics
type DeviceMetric struct {
	DeviceType   string  `json:"device_type"`
	StudentCount int     `json:"student_count"`
	Percentage   float64 `json:"percentage"`
}

// LanguageMetric represents language preference metrics
type LanguageMetric struct {
	Language     string  `json:"language"`
	StudentCount int     `json:"student_count"`
	Percentage   float64 `json:"percentage"`
}

// StudentProgress represents individual student progress
type StudentProgress struct {
	StudentID        uuid.UUID `json:"student_id"`
	StudentName      string    `json:"student_name"`
	StudentEmail     string    `json:"student_email"`
	EnrollmentDate   time.Time `json:"enrollment_date"`
	LastActivityAt   time.Time `json:"last_activity_at"`
	CoursesEnrolled  int       `json:"courses_enrolled"`
	CoursesCompleted int       `json:"courses_completed"`
	OverallProgress  float64   `json:"overall_progress"`
	TotalWatchTime   int       `json:"total_watch_time"`
	EngagementScore  float64   `json:"engagement_score"`
	Status           string    `json:"status"`
}

// TopStudent represents high-performing students
type TopStudent struct {
	StudentID         uuid.UUID `json:"student_id"`
	StudentName       string    `json:"student_name"`
	StudentEmail      string    `json:"student_email"`
	CoursesCompleted  int       `json:"courses_completed"`
	TotalWatchTime    int       `json:"total_watch_time"`
	EngagementScore   float64   `json:"engagement_score"`
	ForumPosts        int       `json:"forum_posts"`
	HelpfulRatings    int       `json:"helpful_ratings"`
	CertificatesEarned int      `json:"certificates_earned"`
}

// AtRiskStudent represents students who might drop out
type AtRiskStudent struct {
	StudentID         uuid.UUID `json:"student_id"`
	StudentName       string    `json:"student_name"`
	StudentEmail      string    `json:"student_email"`
	LastActivityAt    time.Time `json:"last_activity_at"`
	ProgressRate      float64   `json:"progress_rate"`
	EngagementScore   float64   `json:"engagement_score"`
	DaysInactive      int       `json:"days_inactive"`
	RiskScore         float64   `json:"risk_score"`
	RiskFactors       []string  `json:"risk_factors"`
	RecommendedActions []string `json:"recommended_actions"`
}

// VideoAnalytics represents comprehensive video analytics
type VideoAnalytics struct {
	Period            string                `json:"period"`
	StartDate         time.Time             `json:"start_date"`
	EndDate           time.Time             `json:"end_date"`
	TotalVideos       int                   `json:"total_videos"`
	TotalViews        int                   `json:"total_views"`
	TotalWatchTime    int                   `json:"total_watch_time"`
	AvgWatchTime      int                   `json:"avg_watch_time"`
	AvgCompletionRate float64               `json:"avg_completion_rate"`
	TopVideos         []VideoPerformance    `json:"top_videos"`
	UnderperformingVideos []VideoPerformance `json:"underperforming_videos"`
	EngagementHeatmap []EngagementHeatmap   `json:"engagement_heatmap"`
	ViewingPatterns   ViewingPatterns       `json:"viewing_patterns"`
}

// VideoPerformance represents individual video performance
type VideoPerformance struct {
	VideoID           uuid.UUID `json:"video_id"`
	VideoTitle        string    `json:"video_title"`
	CourseID          uuid.UUID `json:"course_id"`
	CourseTitle       string    `json:"course_title"`
	Duration          int       `json:"duration"`
	TotalViews        int       `json:"total_views"`
	UniqueViewers     int       `json:"unique_viewers"`
	TotalWatchTime    int       `json:"total_watch_time"`
	AvgWatchTime      int       `json:"avg_watch_time"`
	CompletionRate    float64   `json:"completion_rate"`
	ReplayRate        float64   `json:"replay_rate"`
	EngagementScore   float64   `json:"engagement_score"`
	ThumbnailClickRate float64  `json:"thumbnail_click_rate"`
	DropOffRate       float64   `json:"drop_off_rate"`
	QualityIssues     int       `json:"quality_issues"`
	BufferingEvents   int       `json:"buffering_events"`
}

// EngagementHeatmap represents engagement data at specific timestamps
type EngagementHeatmap struct {
	VideoID          uuid.UUID `json:"video_id"`
	TimestampSeconds int       `json:"timestamp_seconds"`
	EngagementScore  float64   `json:"engagement_score"`
	ViewerCount      int       `json:"viewer_count"`
	DropOffCount     int       `json:"drop_off_count"`
	ReplayCount      int       `json:"replay_count"`
	AIQuestionCount  int       `json:"ai_question_count"`
	BookmarkCount    int       `json:"bookmark_count"`
}

// ViewingPatterns represents common viewing behavior patterns
type ViewingPatterns struct {
	PeakViewingHours    []HourlyMetric    `json:"peak_viewing_hours"`
	DayOfWeekPatterns   []DayMetric       `json:"day_of_week_patterns"`
	DeviceUsagePatterns []DevicePattern   `json:"device_usage_patterns"`
	QualityPreferences  []QualityPreference `json:"quality_preferences"`
}

// HourlyMetric represents hourly viewing data
type HourlyMetric struct {
	Hour       int     `json:"hour"`
	ViewCount  int     `json:"view_count"`
	WatchTime  int     `json:"watch_time"`
	Percentage float64 `json:"percentage"`
}

// DayMetric represents daily viewing patterns
type DayMetric struct {
	DayOfWeek  string  `json:"day_of_week"`
	ViewCount  int     `json:"view_count"`
	WatchTime  int     `json:"watch_time"`
	Percentage float64 `json:"percentage"`
}

// DevicePattern represents device usage patterns
type DevicePattern struct {
	DeviceType string  `json:"device_type"`
	ViewCount  int     `json:"view_count"`
	WatchTime  int     `json:"watch_time"`
	Percentage float64 `json:"percentage"`
}

// QualityPreference represents video quality preferences
type QualityPreference struct {
	Quality    string  `json:"quality"`
	ViewCount  int     `json:"view_count"`
	Percentage float64 `json:"percentage"`
}

// AnalyticsRequest represents common analytics request parameters
type AnalyticsRequest struct {
	Period    string     `json:"period" form:"period"`         // daily, weekly, monthly, yearly
	StartDate *time.Time `json:"start_date" form:"start_date"`
	EndDate   *time.Time `json:"end_date" form:"end_date"`
	CourseID  *uuid.UUID `json:"course_id" form:"course_id"`
	VideoID   *uuid.UUID `json:"video_id" form:"video_id"`
	Timezone  string     `json:"timezone" form:"timezone"`
}

// AnalyticsFilter represents filtering options for analytics
type AnalyticsFilter struct {
	CourseIDs    []uuid.UUID `json:"course_ids" form:"course_ids"`
	StudentIDs   []uuid.UUID `json:"student_ids" form:"student_ids"`
	Countries    []string    `json:"countries" form:"countries"`
	DeviceTypes  []string    `json:"device_types" form:"device_types"`
	MinEngagement *float64   `json:"min_engagement" form:"min_engagement"`
	MaxEngagement *float64   `json:"max_engagement" form:"max_engagement"`
	MinRevenue   *float64    `json:"min_revenue" form:"min_revenue"`
	MaxRevenue   *float64    `json:"max_revenue" form:"max_revenue"`
}

// AnalyticsSummary represents a high-level analytics summary
type AnalyticsSummary struct {
	Period              string    `json:"period"`
	StartDate           time.Time `json:"start_date"`
	EndDate             time.Time `json:"end_date"`
	TotalRevenue        float64   `json:"total_revenue"`
	RevenueGrowth       float64   `json:"revenue_growth"`
	TotalStudents       int       `json:"total_students"`
	ActiveStudents      int       `json:"active_students"`
	EngagementRate      float64   `json:"engagement_rate"`
	CompletionRate      float64   `json:"completion_rate"`
	StudentSatisfaction float64   `json:"student_satisfaction"`
	ConversionRate      float64   `json:"conversion_rate"`
}