package service

import (
	"context"
	"fmt"
	"time"

	"forum-service/internal/model"
	"forum-service/internal/repository"

	"github.com/google/uuid"
)

type ForumService struct {
	forumRepo *repository.ForumRepository
}

func NewForumService(forumRepo *repository.ForumRepository) *ForumService {
	return &ForumService{
		forumRepo: forumRepo,
	}
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

	response := &model.TopicResponse{
		Topic: topic,
	}

	// TODO: Add user info, vote info, subscription info
	// This would typically involve calling other services (auth service for user info)

	return response, nil
}

func (s *ForumService) ListTopics(ctx context.Context, options *model.ListTopicsOptions, userID *uuid.UUID) ([]*model.TopicResponse, int, error) {
	topics, total, err := s.forumRepo.ListTopics(ctx, options)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list topics: %w", err)
	}

	var responses []*model.TopicResponse
	for _, topic := range topics {
		response := &model.TopicResponse{
			Topic: topic,
		}
		// TODO: Add user info, vote info, subscription info
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
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.forumRepo.CreatePost(ctx, post); err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	return post, nil
}

func (s *ForumService) GetPost(ctx context.Context, postID uuid.UUID, userID *uuid.UUID) (*model.PostResponse, error) {
	post, err := s.forumRepo.GetPostByID(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	response := &model.PostResponse{
		Post:      post,
		VoteTotal: post.UpVotes - post.DownVotes,
	}

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
		response := &model.PostResponse{
			Post:      post,
			VoteTotal: post.UpVotes - post.DownVotes,
		}

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