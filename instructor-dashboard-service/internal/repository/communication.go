package repository

import (
	"database/sql"
	"fmt"
	"time"

	"instructor-dashboard-service/internal/model"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type CommunicationRepository struct {
	db *sql.DB
}

func NewCommunicationRepository(db *sql.DB) *CommunicationRepository {
	return &CommunicationRepository{db: db}
}

// SendBroadcast sends a broadcast message to multiple students
func (r *CommunicationRepository) SendBroadcast(instructorID uuid.UUID, req *model.SendBroadcastRequest) error {
	// Get target student IDs based on course IDs if provided
	var studentIDs []uuid.UUID

	if len(req.CourseIDs) > 0 {
		// Get students enrolled in specified courses
		query := `
			SELECT DISTINCT e.user_id
			FROM enrollments e
			JOIN courses c ON e.course_id = c.id
			WHERE c.creator_id = $1 AND c.id = ANY($2)
		`
		rows, err := r.db.Query(query, instructorID, pq.Array(req.CourseIDs))
		if err != nil {
			return fmt.Errorf("failed to get students from courses: %w", err)
		}
		defer rows.Close()

		for rows.Next() {
			var studentID uuid.UUID
			if err := rows.Scan(&studentID); err != nil {
				return fmt.Errorf("failed to scan student ID: %w", err)
			}
			studentIDs = append(studentIDs, studentID)
		}
	} else if len(req.StudentIDs) > 0 {
		// Use provided student IDs
		studentIDs = req.StudentIDs
	} else {
		// Get all students enrolled in instructor's courses
		query := `
			SELECT DISTINCT e.user_id
			FROM enrollments e
			JOIN courses c ON e.course_id = c.id
			WHERE c.creator_id = $1
		`
		rows, err := r.db.Query(query, instructorID)
		if err != nil {
			return fmt.Errorf("failed to get all students: %w", err)
		}
		defer rows.Close()

		for rows.Next() {
			var studentID uuid.UUID
			if err := rows.Scan(&studentID); err != nil {
				return fmt.Errorf("failed to scan student ID: %w", err)
			}
			studentIDs = append(studentIDs, studentID)
		}
	}

	// Create communication records for each student
	for _, studentID := range studentIDs {
		// For course-specific messages, use the first course ID
		var courseID *uuid.UUID
		if len(req.CourseIDs) > 0 {
			courseID = &req.CourseIDs[0]
		}

		communication := &model.Communication{
			ID:                uuid.New(),
			InstructorID:      instructorID,
			StudentID:         studentID,
			CourseID:          *courseID,
			CommunicationType: req.CommunicationType,
			Subject:           req.Subject,
			Message:           req.Message,
			Status:            "scheduled",
			ScheduledAt:       req.ScheduledAt,
			CreatedAt:         time.Now(),
		}

		// If no scheduled time, mark as sent immediately
		if req.ScheduledAt == nil {
			communication.Status = "sent"
			now := time.Now()
			communication.SentAt = &now
		}

		err := r.CreateCommunication(communication)
		if err != nil {
			return fmt.Errorf("failed to create communication for student %s: %w", studentID.String(), err)
		}
	}

	return nil
}

// CreateCommunication creates a new communication record
func (r *CommunicationRepository) CreateCommunication(comm *model.Communication) error {
	query := `
		INSERT INTO instructor_student_communications (
			id, instructor_id, student_id, course_id, communication_type,
			subject, message, status, scheduled_at, sent_at, created_at, metadata
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	_, err := r.db.Exec(query,
		comm.ID,
		comm.InstructorID,
		comm.StudentID,
		comm.CourseID,
		comm.CommunicationType,
		comm.Subject,
		comm.Message,
		comm.Status,
		comm.ScheduledAt,
		comm.SentAt,
		comm.CreatedAt,
		comm.Metadata,
	)

	if err != nil {
		return fmt.Errorf("failed to create communication: %w", err)
	}

	return nil
}

// GetCommunicationHistory retrieves communication history for an instructor
func (r *CommunicationRepository) GetCommunicationHistory(instructorID uuid.UUID, filter *model.Filter, pagination *model.Pagination) ([]model.Communication, error) {
	communications := []model.Communication{}

	whereClause := "WHERE isc.instructor_id = $1"
	args := []interface{}{instructorID}
	argIndex := 2

	if filter.StartDate != nil {
		whereClause += fmt.Sprintf(" AND isc.created_at >= $%d", argIndex)
		args = append(args, *filter.StartDate)
		argIndex++
	}

	if filter.EndDate != nil {
		whereClause += fmt.Sprintf(" AND isc.created_at <= $%d", argIndex)
		args = append(args, *filter.EndDate)
		argIndex++
	}

	if filter.Status != "" {
		whereClause += fmt.Sprintf(" AND isc.status = $%d", argIndex)
		args = append(args, filter.Status)
		argIndex++
	}

	if filter.CourseID != nil {
		whereClause += fmt.Sprintf(" AND isc.course_id = $%d", argIndex)
		args = append(args, *filter.CourseID)
		argIndex++
	}

	offset := (pagination.Page - 1) * pagination.PageSize

	query := fmt.Sprintf(`
		SELECT
			isc.id, isc.instructor_id, isc.student_id, isc.course_id,
			isc.communication_type, isc.subject, isc.message, isc.status,
			isc.scheduled_at, isc.sent_at, isc.read_at, isc.replied_at,
			isc.metadata, isc.created_at,
			u.email as student_email, u.first_name as student_first_name,
			u.last_name as student_last_name, c.title as course_title
		FROM instructor_student_communications isc
		JOIN users u ON isc.student_id = u.id
		LEFT JOIN courses c ON isc.course_id = c.id
		%s
		ORDER BY isc.created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)

	args = append(args, pagination.PageSize, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get communication history: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var comm model.Communication
		var courseTitle sql.NullString
		err := rows.Scan(
			&comm.ID,
			&comm.InstructorID,
			&comm.StudentID,
			&comm.CourseID,
			&comm.CommunicationType,
			&comm.Subject,
			&comm.Message,
			&comm.Status,
			&comm.ScheduledAt,
			&comm.SentAt,
			&comm.ReadAt,
			&comm.RepliedAt,
			&comm.Metadata,
			&comm.CreatedAt,
			&comm.StudentEmail,
			&comm.StudentFirstName,
			&comm.StudentLastName,
			&courseTitle,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan communication: %w", err)
		}

		if courseTitle.Valid {
			comm.CourseTitle = courseTitle.String
		}

		communications = append(communications, comm)
	}

	return communications, nil
}

// GetCommunicationHistoryCount gets total count for pagination
func (r *CommunicationRepository) GetCommunicationHistoryCount(instructorID uuid.UUID, filter *model.Filter) (int, error) {
	whereClause := "WHERE instructor_id = $1"
	args := []interface{}{instructorID}
	argIndex := 2

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

	if filter.Status != "" {
		whereClause += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, filter.Status)
		argIndex++
	}

	if filter.CourseID != nil {
		whereClause += fmt.Sprintf(" AND course_id = $%d", argIndex)
		args = append(args, *filter.CourseID)
		argIndex++
	}

	query := fmt.Sprintf("SELECT COUNT(*) FROM instructor_student_communications %s", whereClause)

	var count int
	err := r.db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get communication count: %w", err)
	}

	return count, nil
}

// SetupAutomatedMessages creates automated message rules
func (r *CommunicationRepository) SetupAutomatedMessages(instructorID uuid.UUID, rules []model.AutomatedMessageRule) error {
	// This would be implemented to store automated message rules
	// For now, just return success
	return nil
}

// GetAutomatedMessageRules retrieves automated message rules
func (r *CommunicationRepository) GetAutomatedMessageRules(instructorID uuid.UUID) ([]model.AutomatedMessageRule, error) {
	// This would be implemented to get automated message rules
	// For now, return empty results
	rules := []model.AutomatedMessageRule{}
	return rules, nil
}

// UpdateCommunicationStatus updates the status of a communication
func (r *CommunicationRepository) UpdateCommunicationStatus(communicationID uuid.UUID, status string) error {
	query := `
		UPDATE instructor_student_communications
		SET status = $1, updated_at = NOW()
		WHERE id = $2
	`

	_, err := r.db.Exec(query, status, communicationID)
	if err != nil {
		return fmt.Errorf("failed to update communication status: %w", err)
	}

	return nil
}

// MarkCommunicationRead marks a communication as read
func (r *CommunicationRepository) MarkCommunicationRead(communicationID uuid.UUID) error {
	query := `
		UPDATE instructor_student_communications
		SET read_at = NOW()
		WHERE id = $1 AND read_at IS NULL
	`

	_, err := r.db.Exec(query, communicationID)
	if err != nil {
		return fmt.Errorf("failed to mark communication as read: %w", err)
	}

	return nil
}

// MarkCommunicationReplied marks a communication as replied
func (r *CommunicationRepository) MarkCommunicationReplied(communicationID uuid.UUID) error {
	query := `
		UPDATE instructor_student_communications
		SET replied_at = NOW()
		WHERE id = $1 AND replied_at IS NULL
	`

	_, err := r.db.Exec(query, communicationID)
	if err != nil {
		return fmt.Errorf("failed to mark communication as replied: %w", err)
	}

	return nil
}

// GetNotifications retrieves notifications for an instructor
func (r *CommunicationRepository) GetNotifications(instructorID uuid.UUID, limit int, status string) ([]model.Notification, error) {
	notifications := []model.Notification{}

	whereClause := "WHERE instructor_id = $1"
	args := []interface{}{instructorID}
	argIndex := 2

	// Add status filter if provided
	if status == "read" {
		whereClause += fmt.Sprintf(" AND is_read = true")
	} else if status == "unread" {
		whereClause += fmt.Sprintf(" AND is_read = false")
	}

	query := fmt.Sprintf(`
		SELECT
			id, instructor_id, type, title, message, is_read, priority,
			action_url, metadata, created_at, read_at
		FROM instructor_notifications
		%s
		ORDER BY created_at DESC
		LIMIT $%d
	`, whereClause, argIndex)

	args = append(args, limit)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		// If the table doesn't exist, return mock data
		if err.Error() == `pq: relation "instructor_notifications" does not exist` {
			return r.getMockNotifications(instructorID, limit), nil
		}
		return nil, fmt.Errorf("failed to get notifications: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var notif model.Notification
		var actionURL sql.NullString
		err := rows.Scan(
			&notif.ID,
			&notif.InstructorID,
			&notif.Type,
			&notif.Title,
			&notif.Message,
			&notif.IsRead,
			&notif.Priority,
			&actionURL,
			&notif.Metadata,
			&notif.CreatedAt,
			&notif.ReadAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan notification: %w", err)
		}

		if actionURL.Valid {
			notif.ActionURL = &actionURL.String
		}

		notifications = append(notifications, notif)
	}

	return notifications, nil
}

// getMockNotifications returns mock notification data
func (r *CommunicationRepository) getMockNotifications(instructorID uuid.UUID, limit int) []model.Notification {
	now := time.Now()
	mockNotifications := []model.Notification{
		{
			ID:           uuid.New(),
			InstructorID: instructorID,
			Type:         "new_enrollment",
			Title:        "New Student Enrollment",
			Message:      "John Doe just enrolled in your course 'Advanced React Development'",
			IsRead:       false,
			Priority:     "medium",
			CreatedAt:    now.Add(-2 * time.Hour),
		},
		{
			ID:           uuid.New(),
			InstructorID: instructorID,
			Type:         "course_completion",
			Title:        "Course Completion",
			Message:      "Sarah Smith completed your course 'JavaScript Fundamentals' with a 95% score",
			IsRead:       false,
			Priority:     "high",
			CreatedAt:    now.Add(-4 * time.Hour),
		},
		{
			ID:           uuid.New(),
			InstructorID: instructorID,
			Type:         "new_review",
			Title:        "New 5-Star Review",
			Message:      "You received a new 5-star review for 'Python for Beginners'",
			IsRead:       true,
			Priority:     "medium",
			CreatedAt:    now.Add(-6 * time.Hour),
			ReadAt:       &[]time.Time{now.Add(-5 * time.Hour)}[0],
		},
		{
			ID:           uuid.New(),
			InstructorID: instructorID,
			Type:         "revenue_milestone",
			Title:        "Revenue Milestone Reached",
			Message:      "Congratulations! You've reached $5,000 in total course revenue",
			IsRead:       true,
			Priority:     "high",
			CreatedAt:    now.Add(-12 * time.Hour),
			ReadAt:       &[]time.Time{now.Add(-11 * time.Hour)}[0],
		},
		{
			ID:           uuid.New(),
			InstructorID: instructorID,
			Type:         "system_update",
			Title:        "Platform Update",
			Message:      "New analytics features are now available in your instructor dashboard",
			IsRead:       false,
			Priority:     "low",
			CreatedAt:    now.Add(-24 * time.Hour),
		},
	}

	// Apply limit
	if limit > 0 && limit < len(mockNotifications) {
		mockNotifications = mockNotifications[:limit]
	}

	return mockNotifications
}