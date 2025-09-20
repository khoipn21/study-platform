package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type StudentDashboardHandler struct{}

func NewStudentDashboardHandler() *StudentDashboardHandler {
	return &StudentDashboardHandler{}
}

// DashboardStats represents student dashboard statistics
type DashboardStats struct {
	TotalCourses      int `json:"totalCourses"`
	CompletedCourses  int `json:"completedCourses"`
	InProgressCourses int `json:"inProgressCourses"`
	TotalWatchTime    int `json:"totalWatchTime"` // in minutes
	StreakDays        int `json:"streakDays"`
	CertificatesEarned int `json:"certificatesEarned"`
	ForumPosts        int `json:"forumPosts"`
	ChatSessions      int `json:"chatSessions"`
}

// StudyGoal represents a learning goal
type StudyGoal struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Target      int    `json:"target"`
	Current     int    `json:"current"`
	Unit        string `json:"unit"`
	Deadline    string `json:"deadline,omitempty"`
	IsCompleted bool   `json:"isCompleted"`
	CreatedAt   string `json:"createdAt"`
}

// WeeklyStat represents weekly learning statistics
type WeeklyStat struct {
	Date               string `json:"date"`
	Minutes            int    `json:"minutes"`
	LecturesCompleted  int    `json:"lecturesCompleted"`
	CoursesStarted     int    `json:"coursesStarted"`
}

// GetUserDashboard handles GET /api/v1/dashboard/user/{userId}
func (h *StudentDashboardHandler) GetUserDashboard(w http.ResponseWriter, r *http.Request) {
	log.Printf("DEBUG: StudentDashboardHandler.GetUserDashboard called for path: %s", r.URL.Path)
	vars := mux.Vars(r)
	userID := vars["userId"]
	log.Printf("DEBUG: Extracted userID: %s", userID)

	if userID == "" {
		log.Printf("DEBUG: UserID is empty, returning 400")
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// For now, return mock data - in production this would fetch from database
	dashboardData := map[string]interface{}{
		"stats": DashboardStats{
			TotalCourses:       8,
			CompletedCourses:   3,
			InProgressCourses:  2,
			TotalWatchTime:     1247,
			StreakDays:         12,
			CertificatesEarned: 3,
			ForumPosts:         15,
			ChatSessions:       28,
		},
		"recentProgress": []map[string]interface{}{
			{
				"courseId":     "react-mastery",
				"courseTitle":  "React Mastery: Build Modern Web Applications",
				"instructor":   "Sarah Wilson",
				"progress":     68,
				"lastWatched":  "2024-01-15T14:30:00Z",
				"difficulty":   "intermediate",
				"category":     "Web Development",
			},
		},
		"recentActivity": []map[string]interface{}{
			{
				"id":        "activity_1",
				"type":      "lecture_completed",
				"title":     "Completed: React State Management",
				"timestamp": "2024-01-15T14:30:00Z",
				"courseId":  "react-mastery",
			},
		},
		"achievements": []map[string]interface{}{
			{
				"id":          "first_course",
				"title":       "First Steps",
				"description": "Complete your first course",
				"category":    "milestone",
				"points":      100,
				"unlockedAt":  "2024-01-10T12:00:00Z",
			},
		},
		"studyGoals": []StudyGoal{
			{
				ID:          "daily_30min",
				Title:       "Daily Learning",
				Description: "Study for at least 30 minutes every day",
				Type:        "daily",
				Target:      30,
				Current:     22,
				Unit:        "minutes",
				IsCompleted: false,
				CreatedAt:   "2024-01-01T00:00:00Z",
			},
		},
		"learningPaths": []map[string]interface{}{
			{
				"id":           "frontend_developer",
				"title":        "Frontend Developer Path",
				"description":  "Master modern frontend development",
				"progress":     60,
				"difficulty":   "intermediate",
				"category":     "Web Development",
			},
		},
		"upcomingDeadlines": []map[string]interface{}{
			{
				"id":       "goal_deadline_1",
				"title":    "React Mastery Weekly Goal",
				"type":     "goal",
				"dueDate":  "2024-01-21T23:59:59Z",
				"priority": "medium",
			},
		},
		"recommendations": []map[string]interface{}{
			{
				"id":           "rec_1",
				"type":         "course",
				"title":        "Advanced React Patterns",
				"reason":       "Based on your progress in React Mastery",
				"rating":       4.8,
				"studentsCount": 15420,
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dashboardData)
}

// UpdateStudyGoal handles PATCH /api/v1/dashboard/goals/{goalId}
func (h *StudentDashboardHandler) UpdateStudyGoal(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	goalID := vars["goalId"]

	if goalID == "" {
		http.Error(w, "Goal ID is required", http.StatusBadRequest)
		return
	}

	var updateData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// In production, this would update the goal in the database
	response := map[string]interface{}{
		"id":      goalID,
		"message": "Goal updated successfully",
		"updated": updateData,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// CreateStudyGoal handles POST /api/v1/dashboard/goals
func (h *StudentDashboardHandler) CreateStudyGoal(w http.ResponseWriter, r *http.Request) {
	var goal StudyGoal
	if err := json.NewDecoder(r.Body).Decode(&goal); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Basic validation
	if goal.Title == "" || goal.Type == "" || goal.Target <= 0 {
		http.Error(w, "Title, type, and target are required", http.StatusBadRequest)
		return
	}

	// Set default values
	goal.ID = "goal_" + strconv.FormatInt(time.Now().Unix(), 10)
	goal.Current = 0
	goal.IsCompleted = false
	goal.CreatedAt = time.Now().Format(time.RFC3339)

	// In production, this would save to database
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(goal)
}

// GetWeeklyStats handles GET /api/v1/dashboard/stats/weekly/{userId}
func (h *StudentDashboardHandler) GetWeeklyStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Generate mock weekly stats for the past 7 days
	today := time.Now()
	var stats []WeeklyStat

	for i := 6; i >= 0; i-- {
		date := today.AddDate(0, 0, -i)
		stats = append(stats, WeeklyStat{
			Date:              date.Format("2006-01-02"),
			Minutes:           30 + (i * 15), // Mock increasing study time
			LecturesCompleted: i % 3,         // Mock varying lecture completion
			CoursesStarted:    func() int { if i == 6 { return 1 }; return 0 }(),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stats)
}