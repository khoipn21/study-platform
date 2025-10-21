package model

import (
	"time"

	"github.com/google/uuid"
)

type Topic struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	CourseID      *uuid.UUID `json:"course_id,omitempty" db:"course_id"`
	CreatedByID   uuid.UUID  `json:"created_by_id" db:"created_by_id"`
	Title         string     `json:"title" db:"title"`
	Description   string     `json:"description" db:"description"`
	Category      string     `json:"category" db:"category"`
	Tags          []string   `json:"tags" db:"tags"`
	IsSticky      bool       `json:"is_sticky" db:"is_sticky"`
	IsLocked      bool       `json:"is_locked" db:"is_locked"`
	Status        string     `json:"status" db:"status"`
	PinOrder      *int       `json:"pin_order,omitempty" db:"pin_order"`
	ViewCount     int        `json:"view_count" db:"view_count"`
	PostCount     int        `json:"post_count" db:"post_count"`
	LastPostAt    *time.Time `json:"last_post_at,omitempty" db:"last_post_at"`
	LastPostByID  *uuid.UUID `json:"last_post_by_id,omitempty" db:"last_post_by_id"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

type Post struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	TopicID     uuid.UUID  `json:"topic_id" db:"topic_id"`
	AuthorID    uuid.UUID  `json:"author_id" db:"author_id"`
	ParentID    *uuid.UUID `json:"parent_id,omitempty" db:"parent_id"`
	Content     string     `json:"content" db:"content"`
	Status      string     `json:"status" db:"status"`
	PinOrder    *int       `json:"pin_order,omitempty" db:"pin_order"`
	IsEdited    bool       `json:"is_edited" db:"is_edited"`
	EditedAt    *time.Time `json:"edited_at,omitempty" db:"edited_at"`
	UpVotes     int        `json:"up_votes" db:"up_votes"`
	DownVotes   int        `json:"down_votes" db:"down_votes"`
	IsAnswer    bool       `json:"is_answer" db:"is_answer"`
	IsPinned    bool       `json:"is_pinned" db:"is_pinned"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

type Vote struct {
	ID       uuid.UUID `json:"id" db:"id"`
	PostID   uuid.UUID `json:"post_id" db:"post_id"`
	UserID   uuid.UUID `json:"user_id" db:"user_id"`
	VoteType VoteType  `json:"vote_type" db:"vote_type"`
	VotedAt  time.Time `json:"voted_at" db:"voted_at"`
}

type VoteType string

const (
	VoteTypeUp   VoteType = "up"
	VoteTypeDown VoteType = "down"
)

type TopicCategory struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Color       string    `json:"color" db:"color"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type Moderation struct {
	ID          uuid.UUID      `json:"id" db:"id"`
	PostID      *uuid.UUID     `json:"post_id,omitempty" db:"post_id"`
	TopicID     *uuid.UUID     `json:"topic_id,omitempty" db:"topic_id"`
	ModeratorID uuid.UUID      `json:"moderator_id" db:"moderator_id"`
	Action      ModerationAction `json:"action" db:"action"`
	Reason      string         `json:"reason" db:"reason"`
	CreatedAt   time.Time      `json:"created_at" db:"created_at"`
}

type ModerationAction string

const (
	ModerationActionHide   ModerationAction = "hide"
	ModerationActionDelete ModerationAction = "delete"
	ModerationActionLock   ModerationAction = "lock"
	ModerationActionSticky ModerationAction = "sticky"
	ModerationActionPin    ModerationAction = "pin"
)

// Request/Response DTOs
type CreateTopicRequest struct {
	CourseID    *uuid.UUID `json:"course_id,omitempty"`
	Title       string     `json:"title" binding:"required,min=5,max=200"`
	Description string     `json:"description" binding:"required,min=10"`
	Category    string     `json:"category" binding:"required"`
	Tags        []string   `json:"tags,omitempty"`
}

type UpdateTopicRequest struct {
	Title       *string   `json:"title,omitempty"`
	Description *string   `json:"description,omitempty"`
	Category    *string   `json:"category,omitempty"`
	Tags        *[]string `json:"tags,omitempty"`
}

type CreatePostRequest struct {
	TopicID  uuid.UUID  `json:"topic_id" binding:"required"`
	ParentID *uuid.UUID `json:"parent_id,omitempty"`
	Content  string     `json:"content" binding:"required,min=5"`
}

type UpdatePostRequest struct {
	Content *string `json:"content,omitempty"`
}

type VoteRequest struct {
	PostID   uuid.UUID `json:"post_id" binding:"required"`
	VoteType VoteType  `json:"vote_type" binding:"required,oneof=up down"`
}

type TopicResponse struct {
	*Topic
	CreatedBy    *UserInfo `json:"created_by,omitempty"`
	LastPostBy   *UserInfo `json:"last_post_by,omitempty"`
	UserVote     *VoteType `json:"user_vote,omitempty"`
	IsSubscribed bool      `json:"is_subscribed"`
}

type PostResponse struct {
	*Post
	Author     *UserInfo      `json:"author,omitempty"`
	Children   []*PostResponse `json:"children,omitempty"`
	UserVote   *VoteType      `json:"user_vote,omitempty"`
	VoteTotal  int            `json:"vote_total"`
}

type UserInfo struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Avatar   *string   `json:"avatar,omitempty"`
	Role     string    `json:"role"`
}

type ForumStats struct {
	TotalTopics int `json:"total_topics"`
	TotalPosts  int `json:"total_posts"`
	TotalUsers  int `json:"total_users"`
}

type SearchFilters struct {
	CourseID   *uuid.UUID `json:"course_id,omitempty"`
	Category   string     `json:"category,omitempty"`
	Tags       []string   `json:"tags,omitempty"`
	AuthorID   *uuid.UUID `json:"author_id,omitempty"`
	HasAnswers *bool      `json:"has_answers,omitempty"`
	DateFrom   *time.Time `json:"date_from,omitempty"`
	DateTo     *time.Time `json:"date_to,omitempty"`
}

type ListTopicsOptions struct {
	CourseID  *uuid.UUID `json:"course_id,omitempty"`
	Category  string     `json:"category,omitempty"`
	Tags      []string   `json:"tags,omitempty"`
	Search    string     `json:"search,omitempty"`
	Status    string     `json:"status,omitempty"` // pending, approved, rejected
	SortBy    string     `json:"sort_by,omitempty"` // created_at, updated_at, post_count, view_count
	SortOrder string     `json:"sort_order,omitempty"` // asc, desc
	Page      int        `json:"page,omitempty"`
	Limit     int        `json:"limit,omitempty"`
}

type ListPostsOptions struct {
	TopicID   uuid.UUID `json:"topic_id"`
	ParentID  *uuid.UUID `json:"parent_id,omitempty"`
	Status    string    `json:"status,omitempty"` // pending, approved, rejected
	SortBy    string    `json:"sort_by,omitempty"` // created_at, votes, is_answer
	SortOrder string    `json:"sort_order,omitempty"` // asc, desc
	Page      int       `json:"page,omitempty"`
	Limit     int       `json:"limit,omitempty"`
}