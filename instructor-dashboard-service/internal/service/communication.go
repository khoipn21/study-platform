package service

import (
	"fmt"
	"instructor-dashboard-service/internal/model"
	"instructor-dashboard-service/internal/repository"

	"github.com/google/uuid"
)

type CommunicationService struct {
	communicationRepo *repository.CommunicationRepository
}

func NewCommunicationService(communicationRepo *repository.CommunicationRepository) *CommunicationService {
	return &CommunicationService{
		communicationRepo: communicationRepo,
	}
}

// SendBroadcast sends a broadcast message to students
func (s *CommunicationService) SendBroadcast(instructorID uuid.UUID, req *model.SendBroadcastRequest) error {
	// Validate request
	if req.Subject == "" || req.Message == "" {
		return fmt.Errorf("subject and message are required")
	}

	if req.CommunicationType == "" {
		req.CommunicationType = "bulk_announcement"
	}

	// Validate communication type
	validTypes := map[string]bool{
		"welcome_message":        true,
		"milestone_congratulation": true,
		"progress_reminder":      true,
		"course_update":          true,
		"direct_message":         true,
		"bulk_announcement":      true,
	}

	if !validTypes[req.CommunicationType] {
		return fmt.Errorf("invalid communication type: %s", req.CommunicationType)
	}

	err := s.communicationRepo.SendBroadcast(instructorID, req)
	if err != nil {
		return fmt.Errorf("failed to send broadcast: %w", err)
	}

	return nil
}

// GetCommunicationHistory retrieves communication history
func (s *CommunicationService) GetCommunicationHistory(instructorID uuid.UUID, filter *model.Filter, pagination *model.Pagination) ([]model.Communication, *model.Pagination, error) {
	// Get total count for pagination
	totalCount, err := s.communicationRepo.GetCommunicationHistoryCount(instructorID, filter)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get communication count: %w", err)
	}

	// Calculate pagination
	pagination.Total = totalCount
	pagination.Pages = (totalCount + pagination.PageSize - 1) / pagination.PageSize

	// Get communications
	communications, err := s.communicationRepo.GetCommunicationHistory(instructorID, filter, pagination)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get communication history: %w", err)
	}

	return communications, pagination, nil
}

// SetupAutomatedMessages creates automated message rules
func (s *CommunicationService) SetupAutomatedMessages(instructorID uuid.UUID, rules []model.AutomatedMessageRule) error {
	// Validate rules
	for i, rule := range rules {
		if rule.RuleName == "" {
			return fmt.Errorf("rule %d: rule name is required", i)
		}

		if rule.MessageTemplate == "" {
			return fmt.Errorf("rule %d: message template is required", i)
		}

		if rule.Subject == "" {
			return fmt.Errorf("rule %d: subject is required", i)
		}

		// Validate trigger type
		validTriggers := map[string]bool{
			"enrollment":  true,
			"completion":  true,
			"inactivity":  true,
			"milestone":   true,
			"progress":    true,
			"deadline":    true,
		}

		if !validTriggers[rule.TriggerType] {
			return fmt.Errorf("rule %d: invalid trigger type: %s", i, rule.TriggerType)
		}

		// Set instructor ID
		rules[i].InstructorID = instructorID
	}

	err := s.communicationRepo.SetupAutomatedMessages(instructorID, rules)
	if err != nil {
		return fmt.Errorf("failed to setup automated messages: %w", err)
	}

	return nil
}

// GetAutomatedMessageRules retrieves automated message rules
func (s *CommunicationService) GetAutomatedMessageRules(instructorID uuid.UUID) ([]model.AutomatedMessageRule, error) {
	rules, err := s.communicationRepo.GetAutomatedMessageRules(instructorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get automated message rules: %w", err)
	}

	return rules, nil
}

// SendPersonalMessage sends a personal message to a specific student
func (s *CommunicationService) SendPersonalMessage(instructorID, studentID uuid.UUID, req *model.PersonalMessageRequest) error {
	// Validate request
	if req.Subject == "" || req.Message == "" {
		return fmt.Errorf("subject and message are required")
	}

	// Create communication record
	var courseID uuid.UUID
	if req.CourseID != nil {
		courseID = *req.CourseID
	}

	communication := &model.Communication{
		ID:                uuid.New(),
		InstructorID:      instructorID,
		StudentID:         studentID,
		CourseID:          courseID,
		CommunicationType: "direct_message",
		Subject:           req.Subject,
		Message:           req.Message,
		Status:            "sent",
		Metadata:          req.Metadata,
	}

	err := s.communicationRepo.CreateCommunication(communication)
	if err != nil {
		return fmt.Errorf("failed to send personal message: %w", err)
	}

	return nil
}

// MarkMessageRead marks a message as read
func (s *CommunicationService) MarkMessageRead(instructorID, communicationID uuid.UUID) error {
	// TODO: Verify the communication belongs to the instructor
	err := s.communicationRepo.MarkCommunicationRead(communicationID)
	if err != nil {
		return fmt.Errorf("failed to mark message as read: %w", err)
	}

	return nil
}

// UpdateCommunicationStatus updates communication status
func (s *CommunicationService) UpdateCommunicationStatus(instructorID, communicationID uuid.UUID, status string) error {
	// Validate status
	validStatuses := map[string]bool{
		"draft":     true,
		"scheduled": true,
		"sent":      true,
		"delivered": true,
		"failed":    true,
		"cancelled": true,
	}

	if !validStatuses[status] {
		return fmt.Errorf("invalid status: %s", status)
	}

	// TODO: Verify the communication belongs to the instructor
	err := s.communicationRepo.UpdateCommunicationStatus(communicationID, status)
	if err != nil {
		return fmt.Errorf("failed to update communication status: %w", err)
	}

	return nil
}

// GetCommunicationTemplates retrieves message templates
func (s *CommunicationService) GetCommunicationTemplates(instructorID uuid.UUID) ([]model.MessageTemplate, error) {
	// This would retrieve predefined message templates
	// For now, return some default templates
	templates := []model.MessageTemplate{
		{
			ID:       uuid.New().String(),
			Name:     "Welcome Message",
			Subject:  "Welcome to {{course_name}}!",
			Template: "Hi {{student_name}},\n\nWelcome to {{course_name}}! I'm excited to have you in the course.\n\nBest regards,\n{{instructor_name}}",
			Category: "welcome",
		},
		{
			ID:       uuid.New().String(),
			Name:     "Progress Reminder",
			Subject:  "Keep up the great work in {{course_name}}!",
			Template: "Hi {{student_name}},\n\nI noticed you've made great progress in {{course_name}}. Keep it up!\n\nBest regards,\n{{instructor_name}}",
			Category: "progress",
		},
		{
			ID:       uuid.New().String(),
			Name:     "Course Completion",
			Subject:  "Congratulations on completing {{course_name}}!",
			Template: "Hi {{student_name}},\n\nCongratulations on completing {{course_name}}! You've done excellent work.\n\nBest regards,\n{{instructor_name}}",
			Category: "completion",
		},
	}

	return templates, nil
}

// GetNotifications retrieves notifications for an instructor
func (s *CommunicationService) GetNotifications(instructorID uuid.UUID, limit int, status string) ([]model.Notification, error) {
	notifications, err := s.communicationRepo.GetNotifications(instructorID, limit, status)
	if err != nil {
		return nil, fmt.Errorf("failed to get notifications: %w", err)
	}

	return notifications, nil
}