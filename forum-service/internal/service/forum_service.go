package service

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"forum-service/internal/model"
	"forum-service/internal/repository"

	"github.com/google/uuid"
)

type ForumService struct {
	forumRepo      *repository.ForumRepository
	mentionPattern *regexp.Regexp
}

func NewForumService(forumRepo *repository.ForumRepository) *ForumService {
	return &ForumService{
		forumRepo:      forumRepo,
		mentionPattern: regexp.MustCompile(`@(\w+)`),
	}
}

// extractMentions extracts all @username mentions from content
func (s *ForumService) extractMentions(content string) []string {
	matches := s.mentionPattern.FindAllStringSubmatch(content, -1)
	var usernames []string
	seen := make(map[string]bool)
	
	for _, match := range matches {
		if len(match) > 1 {
			username := match[1]
			if !seen[username] {
				usernames = append(usernames, username)
				seen[username] = true
			}
		}
	}
	
	return usernames
}

// processMentions handles mention creation and notifications
func (s *ForumService) processMentions(ctx context.Context, postID uuid.UUID, content string, mentionerID uuid.UUID, topicID uuid.UUID) error {
	usernames := s.extractMentions(content)
	
	for _, username := range usernames {
		// Get user ID by username
		mentionedUserID, err := s.forumRepo.GetUserIDByUsername(ctx, username)
		if err != nil || mentionedUserID == nil {
			// User not found, skip this mention
			continue
		}
		
		// Don't mention yourself
		if *mentionedUserID == mentionerID {
			continue
		}
		
		// Create mention record
		mention := &model.Mention{
			ID:              uuid.New(),
			PostID:          postID,
			MentionedUserID: *mentionedUserID,
			MentionerUserID: mentionerID,
			IsRead:          false,
			CreatedAt:       time.Now(),
		}
		
		if err := s.forumRepo.CreateMention(ctx, mention); err != nil {
			// Log error but continue with other mentions
			fmt.Printf("Failed to create mention for %s: %v\n", username, err)
			continue
		}
		
		// Create notification
		refType := "post"
		notification := &model.Notification{
			ID:            uuid.New(),
			UserID:        *mentionedUserID,
			Type:          "mention",
			Title:         "You were mentioned in a forum post",
			Message:       fmt.Sprintf("@%s mentioned you in a post", username),
			ReferenceID:   &postID,
			ReferenceType: &refType,
			IsRead:        false,
			CreatedAt:     time.Now(),
		}
		
		if err := s.forumRepo.CreateNotification(ctx, notification); err != nil {
			// Log error but continue
			fmt.Printf("Failed to create notification for %s: %v\n", username, err)
		}
	}
	
	return nil
}

// Topic operations
func (s *ForumService) CreateTopic(ctx context.Context, req *model.CreateTopicRequest, userID uuid.UUID) (*model.Topic, error) {
	if req.Title == "" {
		return nil, fmt.Errorf("title is required")
	}
	if req.Description == "" {
		return nil, fmt.Errorf("description is required")
	}
	if req.Category == "" {
		return nil, fmt.Errorf("category is required")
	}

	topic := &model.Topic{
		ID:           uuid.New(),
		CourseID:     req.CourseID,
		CreatedByID:  userID,
		Title:        req.Title,
		Description:  req.Description,
		Category:     req.Category,
		Tags:         req.Tags,
		IsSticky:     false,
		IsLocked:     false,
		Status:       "pending", // All topics start as pending
		ViewCount:    0,
		PostCount:    0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.forumRepo.CreateTopic(ctx, topic); err != nil {
		return nil, fmt.Errorf("failed to create topic: %w", err)
	}

	return topic, nil
}

func (s *ForumService) GetTopic(ctx context.Context, topicID uuid.UUID, userID *uuid.UUID) (*model.TopicResponse, error) {
	topic, err := s.forumRepo.GetTopicByID(ctx, topicID)
	if err != nil {
		return nil, fmt.Errorf("failed to get topic: %w", err)
	}

	// Increment view count
	go s.forumRepo.IncrementTopicViewCount(context.Background(), topicID)

	response := s.topicToResponse(ctx, topic)

	return response, nil
}

func (s *ForumService) ListTopics(ctx context.Context, options *model.ListTopicsOptions, userID *uuid.UUID) ([]*model.TopicResponse, int, error) {
	topics, total, err := s.forumRepo.ListTopics(ctx, options)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list topics: %w", err)
	}

	var responses []*model.TopicResponse
	for _, topic := range topics {
		response := s.topicToResponse(ctx, topic)
		responses = append(responses, response)
	}

	return responses, total, nil
}

func (s *ForumService) UpdateTopic(ctx context.Context, topicID uuid.UUID, req *model.UpdateTopicRequest, userID uuid.UUID, userRole string) (*model.Topic, error) {
	topic, err := s.forumRepo.GetTopicByID(ctx, topicID)
	if err != nil {
		return nil, fmt.Errorf("topic not found: %w", err)
	}

	// Check permissions
	if topic.CreatedByID != userID && userRole != "admin" && userRole != "instructor" {
		return nil, fmt.Errorf("access denied")
	}

	// Update fields
	if req.Title != nil {
		topic.Title = *req.Title
	}
	if req.Description != nil {
		topic.Description = *req.Description
	}
	if req.Category != nil {
		topic.Category = *req.Category
	}
	if req.Tags != nil {
		topic.Tags = *req.Tags
	}
	topic.UpdatedAt = time.Now()

	if err := s.forumRepo.UpdateTopic(ctx, topic); err != nil {
		return nil, fmt.Errorf("failed to update topic: %w", err)
	}

	return topic, nil
}

func (s *ForumService) DeleteTopic(ctx context.Context, topicID uuid.UUID, userID uuid.UUID, userRole string) error {
	topic, err := s.forumRepo.GetTopicByID(ctx, topicID)
	if err != nil {
		return fmt.Errorf("topic not found: %w", err)
	}

	// Check permissions
	if topic.CreatedByID != userID && userRole != "admin" && userRole != "instructor" {
		return fmt.Errorf("access denied")
	}

	if err := s.forumRepo.DeleteTopic(ctx, topicID); err != nil {
		return fmt.Errorf("failed to delete topic: %w", err)
	}

	return nil
}

func (s *ForumService) ToggleTopicSticky(ctx context.Context, topicID uuid.UUID, userRole string) (*model.Topic, error) {
	if userRole != "admin" && userRole != "instructor" {
		return nil, fmt.Errorf("access denied: only moderators can pin topics")
	}

	topic, err := s.forumRepo.GetTopicByID(ctx, topicID)
	if err != nil {
		return nil, fmt.Errorf("topic not found: %w", err)
	}

	topic.IsSticky = !topic.IsSticky
	topic.UpdatedAt = time.Now()

	if err := s.forumRepo.UpdateTopic(ctx, topic); err != nil {
		return nil, fmt.Errorf("failed to update topic: %w", err)
	}

	return topic, nil
}

func (s *ForumService) ToggleTopicLock(ctx context.Context, topicID uuid.UUID, userRole string) (*model.Topic, error) {
	if userRole != "admin" && userRole != "instructor" {
		return nil, fmt.Errorf("access denied: only moderators can lock topics")
	}

	topic, err := s.forumRepo.GetTopicByID(ctx, topicID)
	if err != nil {
		return nil, fmt.Errorf("topic not found: %w", err)
	}

	topic.IsLocked = !topic.IsLocked
	topic.UpdatedAt = time.Now()

	if err := s.forumRepo.UpdateTopic(ctx, topic); err != nil {
		return nil, fmt.Errorf("failed to update topic: %w", err)
	}

	return topic, nil
}

// Post operations
func (s *ForumService) CreatePost(ctx context.Context, req *model.CreatePostRequest, userID uuid.UUID) (*model.Post, error) {
	if req.Content == "" {
		return nil, fmt.Errorf("content is required")
	}

	// Check if topic exists and is not locked
	topic, err := s.forumRepo.GetTopicByID(ctx, req.TopicID)
	if err != nil {
		return nil, fmt.Errorf("topic not found: %w", err)
	}

	if topic.IsLocked {
		return nil, fmt.Errorf("topic is locked")
	}

	// If this is a reply, check if parent post exists
	if req.ParentID != nil {
		_, err := s.forumRepo.GetPostByID(ctx, *req.ParentID)
		if err != nil {
			return nil, fmt.Errorf("parent post not found: %w", err)
		}
	}

	post := &model.Post{
		ID:        uuid.New(),
		TopicID:   req.TopicID,
		AuthorID:  userID,
		ParentID:  req.ParentID,
		Content:   req.Content,
		IsEdited:  false,
		UpVotes:   0,
		DownVotes: 0,
		IsAnswer:  false,
		IsPinned:  false,
		Status:    "pending", // All posts start as pending
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.forumRepo.CreatePost(ctx, post); err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	// Process mentions asynchronously (don't fail post creation if mentions fail)
	go s.processMentions(context.Background(), post.ID, post.Content, userID, req.TopicID)

	return post, nil
}

func (s *ForumService) GetPost(ctx context.Context, postID uuid.UUID, userID *uuid.UUID) (*model.PostResponse, error) {
	post, err := s.forumRepo.GetPostByID(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	response := s.postToResponse(ctx, post)

	// Get user vote if user is authenticated
	if userID != nil {
		userVote, err := s.forumRepo.GetUserVote(ctx, postID, *userID)
		if err == nil {
			response.UserVote = userVote
		}
	}

	// TODO: Add author info from user service
	// TODO: Add children posts if needed

	return response, nil
}

func (s *ForumService) ListPosts(ctx context.Context, options *model.ListPostsOptions, userID *uuid.UUID) ([]*model.PostResponse, int, error) {
	posts, total, err := s.forumRepo.ListPosts(ctx, options)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list posts: %w", err)
	}

	var responses []*model.PostResponse
	for _, post := range posts {
		response := s.postToResponse(ctx, post)

		// Get user vote if user is authenticated
		if userID != nil {
			userVote, err := s.forumRepo.GetUserVote(ctx, post.ID, *userID)
			if err == nil {
				response.UserVote = userVote
			}
		}

		// TODO: Add author info from user service
		responses = append(responses, response)
	}

	return responses, total, nil
}

func (s *ForumService) UpdatePost(ctx context.Context, postID uuid.UUID, req *model.UpdatePostRequest, userID uuid.UUID, userRole string) (*model.Post, error) {
	post, err := s.forumRepo.GetPostByID(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("post not found: %w", err)
	}

	// Check permissions
	if post.AuthorID != userID && userRole != "admin" && userRole != "instructor" {
		return nil, fmt.Errorf("access denied")
	}

	// Update fields
	if req.Content != nil {
		post.Content = *req.Content
		post.IsEdited = true
		now := time.Now()
		post.EditedAt = &now
	}
	post.UpdatedAt = time.Now()

	if err := s.forumRepo.UpdatePost(ctx, post); err != nil {
		return nil, fmt.Errorf("failed to update post: %w", err)
	}

	return post, nil
}

func (s *ForumService) DeletePost(ctx context.Context, postID uuid.UUID, userID uuid.UUID, userRole string) error {
	post, err := s.forumRepo.GetPostByID(ctx, postID)
	if err != nil {
		return fmt.Errorf("post not found: %w", err)
	}

	// Check permissions
	if post.AuthorID != userID && userRole != "admin" && userRole != "instructor" {
		return fmt.Errorf("access denied")
	}

	if err := s.forumRepo.DeletePost(ctx, postID); err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	return nil
}

func (s *ForumService) MarkPostAsAnswer(ctx context.Context, postID uuid.UUID, userID uuid.UUID, userRole string) (*model.Post, error) {
	post, err := s.forumRepo.GetPostByID(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("post not found: %w", err)
	}

	// Get topic to check ownership
	topic, err := s.forumRepo.GetTopicByID(ctx, post.TopicID)
	if err != nil {
		return nil, fmt.Errorf("topic not found: %w", err)
	}

	// Check permissions (topic creator, post author, or moderator)
	if topic.CreatedByID != userID && post.AuthorID != userID && userRole != "admin" && userRole != "instructor" {
		return nil, fmt.Errorf("access denied")
	}

	post.IsAnswer = !post.IsAnswer
	post.UpdatedAt = time.Now()

	if err := s.forumRepo.UpdatePost(ctx, post); err != nil {
		return nil, fmt.Errorf("failed to update post: %w", err)
	}

	return post, nil
}

func (s *ForumService) TogglePostPin(ctx context.Context, postID uuid.UUID, userRole string) (*model.Post, error) {
	if userRole != "admin" && userRole != "instructor" {
		return nil, fmt.Errorf("access denied: only moderators can pin posts")
	}

	post, err := s.forumRepo.GetPostByID(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("post not found: %w", err)
	}

	post.IsPinned = !post.IsPinned
	post.UpdatedAt = time.Now()

	if err := s.forumRepo.UpdatePost(ctx, post); err != nil {
		return nil, fmt.Errorf("failed to update post: %w", err)
	}

	return post, nil
}

// Voting operations
func (s *ForumService) VotePost(ctx context.Context, req *model.VoteRequest, userID uuid.UUID) error {
	// Check if post exists
	_, err := s.forumRepo.GetPostByID(ctx, req.PostID)
	if err != nil {
		return fmt.Errorf("post not found: %w", err)
	}

	vote := &model.Vote{
		ID:       uuid.New(),
		PostID:   req.PostID,
		UserID:   userID,
		VoteType: req.VoteType,
		VotedAt:  time.Now(),
	}

	if err := s.forumRepo.CreateOrUpdateVote(ctx, vote); err != nil {
		return fmt.Errorf("failed to vote: %w", err)
	}

	return nil
}

func (s *ForumService) RemoveVote(ctx context.Context, postID uuid.UUID, userID uuid.UUID) error {
	if err := s.forumRepo.DeleteVote(ctx, postID, userID); err != nil {
		return fmt.Errorf("failed to remove vote: %w", err)
	}

	return nil
}

// Search functionality
func (s *ForumService) SearchTopics(ctx context.Context, query string, filters *model.SearchFilters, userID *uuid.UUID) ([]*model.TopicResponse, int, error) {
	options := &model.ListTopicsOptions{
		Search:    query,
		CourseID:  filters.CourseID,
		Category:  filters.Category,
		Tags:      filters.Tags,
		SortBy:    "updated_at",
		SortOrder: "desc",
		Page:      1,
		Limit:     50,
	}

	return s.ListTopics(ctx, options, userID)
}

// Approval methods - Only instructors/admins can approve
func (s *ForumService) ApproveTopic(ctx context.Context, topicID uuid.UUID, status string, userID uuid.UUID, userRole string) error {
	// Check if user has permission to approve
	if userRole != "admin" && userRole != "instructor" {
		return fmt.Errorf("access denied: only instructors and admins can approve topics")
	}

	// Get topic to check course ownership for instructors
	topic, err := s.forumRepo.GetTopicByID(ctx, topicID)
	if err != nil {
		return fmt.Errorf("topic not found: %w", err)
	}

	// If instructor, check if they're instructor of the course (for course-specific topics)
	if userRole == "instructor" && topic.CourseID != nil {
		isInstructor, err := s.forumRepo.IsUserInstructorOfCourse(ctx, userID, *topic.CourseID)
		if err != nil {
			return fmt.Errorf("failed to check instructor status: %w", err)
		}
		if !isInstructor {
			return fmt.Errorf("access denied: you are not the instructor of this course")
		}
	}

	if status != "approved" && status != "rejected" {
		return fmt.Errorf("invalid status: must be 'approved' or 'rejected'")
	}

	return s.forumRepo.ApproveTopicStatus(ctx, topicID, status)
}

func (s *ForumService) ApprovePost(ctx context.Context, postID uuid.UUID, status string, userID uuid.UUID, userRole string) error {
	// Check if user has permission to approve
	if userRole != "admin" && userRole != "instructor" {
		return fmt.Errorf("access denied: only instructors and admins can approve posts")
	}

	// Get post and topic to check course ownership for instructors
	post, err := s.forumRepo.GetPostByID(ctx, postID)
	if err != nil {
		return fmt.Errorf("post not found: %w", err)
	}

	topic, err := s.forumRepo.GetTopicByID(ctx, post.TopicID)
	if err != nil {
		return fmt.Errorf("topic not found: %w", err)
	}

	// If instructor, check if they're instructor of the course (for course-specific posts)
	if userRole == "instructor" && topic.CourseID != nil {
		isInstructor, err := s.forumRepo.IsUserInstructorOfCourse(ctx, userID, *topic.CourseID)
		if err != nil {
			return fmt.Errorf("failed to check instructor status: %w", err)
		}
		if !isInstructor {
			return fmt.Errorf("access denied: you are not the instructor of this course")
		}
	}

	if status != "approved" && status != "rejected" {
		return fmt.Errorf("invalid status: must be 'approved' or 'rejected'")
	}

	return s.forumRepo.ApprovePostStatus(ctx, postID, status)
}

func (s *ForumService) GetPendingTopics(ctx context.Context, courseID *uuid.UUID, userID uuid.UUID, userRole string) ([]*model.TopicResponse, error) {
	// Check if user has permission to see pending topics
	if userRole != "admin" && userRole != "instructor" {
		return nil, fmt.Errorf("access denied: only instructors and admins can view pending topics")
	}

	// If instructor, they can only see pending topics for their courses or general forum
	if userRole == "instructor" && courseID != nil {
		isInstructor, err := s.forumRepo.IsUserInstructorOfCourse(ctx, userID, *courseID)
		if err != nil {
			return nil, fmt.Errorf("failed to check instructor status: %w", err)
		}
		if !isInstructor {
			return nil, fmt.Errorf("access denied: you are not the instructor of this course")
		}
	}

	topics, err := s.forumRepo.GetPendingTopics(ctx, courseID)
	if err != nil {
		return nil, err
	}

	// Convert topics to response format with author info
	var responses []*model.TopicResponse
	for _, topic := range topics {
		response := s.topicToResponse(ctx, topic)
		responses = append(responses, response)
	}

	return responses, nil
}

// Pin management methods
func (s *ForumService) SetTopicPinOrder(ctx context.Context, topicID uuid.UUID, pinOrder *int, userID uuid.UUID, userRole string) error {
	if userRole != "admin" && userRole != "instructor" {
		return fmt.Errorf("access denied: only instructors and admins can pin topics")
	}

	// Get topic to check course ownership for instructors
	topic, err := s.forumRepo.GetTopicByID(ctx, topicID)
	if err != nil {
		return fmt.Errorf("topic not found: %w", err)
	}

	// If instructor, check if they're instructor of the course
	if userRole == "instructor" && topic.CourseID != nil {
		isInstructor, err := s.forumRepo.IsUserInstructorOfCourse(ctx, userID, *topic.CourseID)
		if err != nil {
			return fmt.Errorf("failed to check instructor status: %w", err)
		}
		if !isInstructor {
			return fmt.Errorf("access denied: you are not the instructor of this course")
		}
	}

	return s.forumRepo.SetTopicPinOrder(ctx, topicID, pinOrder)
}

func (s *ForumService) SetPostPinOrder(ctx context.Context, postID uuid.UUID, pinOrder *int, userID uuid.UUID, userRole string) error {
	if userRole != "admin" && userRole != "instructor" {
		return fmt.Errorf("access denied: only instructors and admins can pin posts")
	}

	// Get post and topic to check course ownership for instructors
	post, err := s.forumRepo.GetPostByID(ctx, postID)
	if err != nil {
		return fmt.Errorf("post not found: %w", err)
	}

	topic, err := s.forumRepo.GetTopicByID(ctx, post.TopicID)
	if err != nil {
		return fmt.Errorf("topic not found: %w", err)
	}

	// If instructor, check if they're instructor of the course
	if userRole == "instructor" && topic.CourseID != nil {
		isInstructor, err := s.forumRepo.IsUserInstructorOfCourse(ctx, userID, *topic.CourseID)
		if err != nil {
			return fmt.Errorf("failed to check instructor status: %w", err)
		}
		if !isInstructor {
			return fmt.Errorf("access denied: you are not the instructor of this course")
		}
	}

	return s.forumRepo.SetPostPinOrder(ctx, postID, pinOrder)
}

// getUserInfo fetches user information from database
func (s *ForumService) getUserInfo(ctx context.Context, userID uuid.UUID) *model.UserInfo {
	// Fetch user info from database via repository
	userInfo, err := s.forumRepo.GetUserInfo(ctx, userID)
	if err != nil {
		// Return basic author info if fetch fails
		return &model.UserInfo{
			ID:       userID,
			Username: "Unknown User",
			Role:     "student",
		}
	}
	
	return userInfo
}

// Helper function to convert Topic to TopicResponse with author info
func (s *ForumService) topicToResponse(ctx context.Context, topic *model.Topic) *model.TopicResponse {
	response := &model.TopicResponse{
		ID:           topic.ID.String(),
		CreatedByID:  topic.CreatedByID.String(),
		Title:        topic.Title,
		Description:  topic.Description,
		Category:     topic.Category,
		Tags:         topic.Tags,
		IsSticky:     topic.IsSticky,
		IsLocked:     topic.IsLocked,
		Status:       topic.Status,
		PinOrder:     topic.PinOrder,
		ViewCount:    topic.ViewCount,
		PostCount:    topic.PostCount,
		CreatedAt:    topic.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    topic.UpdatedAt.Format(time.RFC3339),
		CreatedBy:    s.getUserInfo(ctx, topic.CreatedByID),
		IsSubscribed: false,
	}

	if topic.CourseID != nil {
		courseIDStr := topic.CourseID.String()
		response.CourseID = &courseIDStr
	}

	if topic.LastPostAt != nil {
		lastPostAtStr := topic.LastPostAt.Format(time.RFC3339)
		response.LastPostAt = &lastPostAtStr
	}

	if topic.LastPostByID != nil {
		lastPostByIDStr := topic.LastPostByID.String()
		response.LastPostByID = &lastPostByIDStr
	}

	return response
}

// Helper function to convert Post to PostResponse with author info
func (s *ForumService) postToResponse(ctx context.Context, post *model.Post) *model.PostResponse {
	response := &model.PostResponse{
		ID:        post.ID.String(),
		TopicID:   post.TopicID.String(),
		AuthorID:  post.AuthorID.String(),
		Content:   post.Content,
		IsEdited:  post.IsEdited,
		UpVotes:   post.UpVotes,
		DownVotes: post.DownVotes,
		IsAnswer:  post.IsAnswer,
		IsPinned:  post.IsPinned,
		Status:    post.Status,
		PinOrder:  post.PinOrder,
		CreatedAt: post.CreatedAt.Format(time.RFC3339),
		UpdatedAt: post.UpdatedAt.Format(time.RFC3339),
		Author:    s.getUserInfo(ctx, post.AuthorID),
		VoteTotal: post.UpVotes - post.DownVotes,
	}

	if post.ParentID != nil {
		parentIDStr := post.ParentID.String()
		response.ParentID = &parentIDStr
	}

	if post.EditedAt != nil {
		editedAtStr := post.EditedAt.Format(time.RFC3339)
		response.EditedAt = &editedAtStr
	}

	return response
}

// Enrollment checking for course-specific forums
func (s *ForumService) CanUserAccessCourseForum(ctx context.Context, courseID uuid.UUID, userID uuid.UUID, userRole string) (bool, error) {
	// Admins and instructors can always access
	if userRole == "admin" || userRole == "instructor" {
		return true, nil
	}

	// Check if user is instructor of the course
	isInstructor, err := s.forumRepo.IsUserInstructorOfCourse(ctx, userID, courseID)
	if err != nil {
		return false, err
	}
	if isInstructor {
		return true, nil
	}

	// Check if user is enrolled in the course
	isEnrolled, err := s.forumRepo.IsUserEnrolledInCourse(ctx, userID, courseID)
	if err != nil {
		return false, err
	}

	return isEnrolled, nil
}