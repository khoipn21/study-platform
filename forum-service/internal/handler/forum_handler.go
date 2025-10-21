package handler

import (
	"net/http"
	"strconv"

	"forum-service/internal/model"
	"forum-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ForumHandler struct {
	forumService *service.ForumService
}

func NewForumHandler(forumService *service.ForumService) *ForumHandler {
	return &ForumHandler{
		forumService: forumService,
	}
}

// Topic handlers
func (h *ForumHandler) CreateTopic(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDStr, ok := userIDInterface.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID format"})
		return
	}

	var req model.CreateTopicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	topic, err := h.forumService.CreateTopic(c.Request.Context(), &req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, topic)
}

func (h *ForumHandler) GetTopic(c *gin.Context) {
	topicIDStr := c.Param("topicId")
	topicID, err := uuid.Parse(topicIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid topic ID"})
		return
	}

	var userID *uuid.UUID
	if userIDInterface, exists := c.Get("user_id"); exists {
		if userIDStr, ok := userIDInterface.(string); ok {
			if parsedUserID, err := uuid.Parse(userIDStr); err == nil {
				userID = &parsedUserID
			}
		}
	}

	topic, err := h.forumService.GetTopic(c.Request.Context(), topicID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, topic)
}

func (h *ForumHandler) ListTopics(c *gin.Context) {
	var userID *uuid.UUID
	if userIDInterface, exists := c.Get("user_id"); exists {
		if userIDStr, ok := userIDInterface.(string); ok {
			if parsedUserID, err := uuid.Parse(userIDStr); err == nil {
				userID = &parsedUserID
			}
		}
	}

	options := &model.ListTopicsOptions{
		Search:    c.Query("search"),
		Category:  c.Query("category"),
		SortBy:    c.Query("sort_by"),
		SortOrder: c.Query("sort_order"),
		Page:      1,
		Limit:     20,
	}

	// Parse course ID
	if courseIDStr := c.Query("course_id"); courseIDStr != "" {
		if courseID, err := uuid.Parse(courseIDStr); err == nil {
			options.CourseID = &courseID
		}
	}

	// Parse tags
	if tagsStr := c.Query("tags"); tagsStr != "" {
		// Simple comma-separated tags for now
		// In production, you might want better parsing
		options.Tags = []string{tagsStr}
	}

	// Parse pagination
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			options.Page = page
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 && limit <= 100 {
			options.Limit = limit
		}
	}

	topics, total, err := h.forumService.ListTopics(c.Request.Context(), options, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"topics": topics,
		"total":  total,
		"page":   options.Page,
		"limit":  options.Limit,
	})
}

func (h *ForumHandler) UpdateTopic(c *gin.Context) {
	topicIDStr := c.Param("topicId")
	topicID, err := uuid.Parse(topicIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid topic ID"})
		return
	}

	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDStr, ok := userIDInterface.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID format"})
		return
	}

	userRole, _ := c.Get("user_role")
	userRoleStr, _ := userRole.(string)

	var req model.UpdateTopicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	topic, err := h.forumService.UpdateTopic(c.Request.Context(), topicID, &req, userID, userRoleStr)
	if err != nil {
		if err.Error() == "access denied" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, topic)
}

func (h *ForumHandler) DeleteTopic(c *gin.Context) {
	topicIDStr := c.Param("topicId")
	topicID, err := uuid.Parse(topicIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid topic ID"})
		return
	}

	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDStr, ok := userIDInterface.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID format"})
		return
	}

	userRole, _ := c.Get("user_role")
	userRoleStr, _ := userRole.(string)

	err = h.forumService.DeleteTopic(c.Request.Context(), topicID, userID, userRoleStr)
	if err != nil {
		if err.Error() == "access denied" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Topic deleted successfully"})
}

func (h *ForumHandler) ToggleTopicSticky(c *gin.Context) {
	topicIDStr := c.Param("topicId")
	topicID, err := uuid.Parse(topicIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid topic ID"})
		return
	}

	userRole, exists := c.Get("user_role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userRoleStr, _ := userRole.(string)

	topic, err := h.forumService.ToggleTopicSticky(c.Request.Context(), topicID, userRoleStr)
	if err != nil {
		if err.Error() == "access denied: only moderators can pin topics" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, topic)
}

func (h *ForumHandler) ToggleTopicLock(c *gin.Context) {
	topicIDStr := c.Param("topicId")
	topicID, err := uuid.Parse(topicIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid topic ID"})
		return
	}

	userRole, exists := c.Get("user_role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userRoleStr, _ := userRole.(string)

	topic, err := h.forumService.ToggleTopicLock(c.Request.Context(), topicID, userRoleStr)
	if err != nil {
		if err.Error() == "access denied: only moderators can lock topics" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, topic)
}

// Post handlers
func (h *ForumHandler) CreatePost(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDStr, ok := userIDInterface.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID format"})
		return
	}

	var req model.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post, err := h.forumService.CreatePost(c.Request.Context(), &req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, post)
}

func (h *ForumHandler) GetPost(c *gin.Context) {
	postIDStr := c.Param("postId")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	var userID *uuid.UUID
	if userIDInterface, exists := c.Get("user_id"); exists {
		if userIDStr, ok := userIDInterface.(string); ok {
			if parsedUserID, err := uuid.Parse(userIDStr); err == nil {
				userID = &parsedUserID
			}
		}
	}

	post, err := h.forumService.GetPost(c.Request.Context(), postID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, post)
}

func (h *ForumHandler) ListPosts(c *gin.Context) {
	topicIDStr := c.Param("topicId")
	topicID, err := uuid.Parse(topicIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid topic ID"})
		return
	}

	var userID *uuid.UUID
	if userIDInterface, exists := c.Get("user_id"); exists {
		if userIDStr, ok := userIDInterface.(string); ok {
			if parsedUserID, err := uuid.Parse(userIDStr); err == nil {
				userID = &parsedUserID
			}
		}
	}

	options := &model.ListPostsOptions{
		TopicID:   topicID,
		SortBy:    c.Query("sort_by"),
		SortOrder: c.Query("sort_order"),
		Page:      1,
		Limit:     50,
	}

	// Parse parent ID for nested replies
	if parentIDStr := c.Query("parent_id"); parentIDStr != "" {
		if parentID, err := uuid.Parse(parentIDStr); err == nil {
			options.ParentID = &parentID
		}
	}

	// Parse pagination
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			options.Page = page
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 && limit <= 100 {
			options.Limit = limit
		}
	}

	posts, total, err := h.forumService.ListPosts(c.Request.Context(), options, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
		"total": total,
		"page":  options.Page,
		"limit": options.Limit,
	})
}

func (h *ForumHandler) UpdatePost(c *gin.Context) {
	postIDStr := c.Param("postId")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDStr, ok := userIDInterface.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID format"})
		return
	}

	userRole, _ := c.Get("user_role")
	userRoleStr, _ := userRole.(string)

	var req model.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post, err := h.forumService.UpdatePost(c.Request.Context(), postID, &req, userID, userRoleStr)
	if err != nil {
		if err.Error() == "access denied" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, post)
}

func (h *ForumHandler) DeletePost(c *gin.Context) {
	postIDStr := c.Param("postId")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDStr, ok := userIDInterface.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID format"})
		return
	}

	userRole, _ := c.Get("user_role")
	userRoleStr, _ := userRole.(string)

	err = h.forumService.DeletePost(c.Request.Context(), postID, userID, userRoleStr)
	if err != nil {
		if err.Error() == "access denied" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}

func (h *ForumHandler) MarkPostAsAnswer(c *gin.Context) {
	postIDStr := c.Param("postId")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDStr, ok := userIDInterface.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID format"})
		return
	}

	userRole, _ := c.Get("user_role")
	userRoleStr, _ := userRole.(string)

	post, err := h.forumService.MarkPostAsAnswer(c.Request.Context(), postID, userID, userRoleStr)
	if err != nil {
		if err.Error() == "access denied" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, post)
}

func (h *ForumHandler) TogglePostPin(c *gin.Context) {
	postIDStr := c.Param("postId")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	userRole, exists := c.Get("user_role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userRoleStr, _ := userRole.(string)

	post, err := h.forumService.TogglePostPin(c.Request.Context(), postID, userRoleStr)
	if err != nil {
		if err.Error() == "access denied: only moderators can pin posts" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, post)
}

// Voting handlers
func (h *ForumHandler) VotePost(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDStr, ok := userIDInterface.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID format"})
		return
	}

	var req model.VoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.forumService.VotePost(c.Request.Context(), &req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Vote recorded successfully"})
}

func (h *ForumHandler) RemoveVote(c *gin.Context) {
	postIDStr := c.Param("postId")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDStr, ok := userIDInterface.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID format"})
		return
	}

	err = h.forumService.RemoveVote(c.Request.Context(), postID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Vote removed successfully"})
}

// Search handler
func (h *ForumHandler) SearchTopics(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	var userID *uuid.UUID
	if userIDInterface, exists := c.Get("user_id"); exists {
		if userIDStr, ok := userIDInterface.(string); ok {
			if parsedUserID, err := uuid.Parse(userIDStr); err == nil {
				userID = &parsedUserID
			}
		}
	}

	filters := &model.SearchFilters{}

	// Parse course ID filter
	if courseIDStr := c.Query("course_id"); courseIDStr != "" {
		if courseID, err := uuid.Parse(courseIDStr); err == nil {
			filters.CourseID = &courseID
		}
	}

	// Parse category filter
	if category := c.Query("category"); category != "" {
		filters.Category = category
	}

	topics, total, err := h.forumService.SearchTopics(c.Request.Context(), query, filters, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"topics": topics,
		"total":  total,
		"query":  query,
	})
}
// Approval handlers
func (h *ForumHandler) GetPendingTopics(c *gin.Context) {
	userRole, _ := c.Get("user_role")
	userRoleStr, _ := userRole.(string)
	
	// Only instructors and admins can view pending items
	if userRoleStr != "instructor" && userRoleStr != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}
	
	var courseID *uuid.UUID
	if courseIDStr := c.Query("course_id"); courseIDStr != "" {
		cid, err := uuid.Parse(courseIDStr)
		if err == nil {
			courseID = &cid
		}
	}
	
	topics, err := h.forumService.GetPendingTopics(c.Request.Context(), courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"topics": topics})
}

func (h *ForumHandler) ApproveTopic(c *gin.Context) {
	topicIDStr := c.Param("topicId")
	topicID, err := uuid.Parse(topicIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid topic ID"})
		return
	}
	
	userIDInterface, _ := c.Get("user_id")
	userIDStr, _ := userIDInterface.(string)
	userID, _ := uuid.Parse(userIDStr)
	
	userRole, _ := c.Get("user_role")
	userRoleStr, _ := userRole.(string)
	
	err = h.forumService.ApproveTopic(c.Request.Context(), topicID, userID, userRoleStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Topic approved successfully"})
}

func (h *ForumHandler) RejectTopic(c *gin.Context) {
	topicIDStr := c.Param("topicId")
	topicID, err := uuid.Parse(topicIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid topic ID"})
		return
	}
	
	userIDInterface, _ := c.Get("user_id")
	userIDStr, _ := userIDInterface.(string)
	userID, _ := uuid.Parse(userIDStr)
	
	userRole, _ := c.Get("user_role")
	userRoleStr, _ := userRole.(string)
	
	err = h.forumService.RejectTopic(c.Request.Context(), topicID, userID, userRoleStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Topic rejected successfully"})
}

func (h *ForumHandler) ApprovePost(c *gin.Context) {
	postIDStr := c.Param("postId")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}
	
	userIDInterface, _ := c.Get("user_id")
	userIDStr, _ := userIDInterface.(string)
	userID, _ := uuid.Parse(userIDStr)
	
	userRole, _ := c.Get("user_role")
	userRoleStr, _ := userRole.(string)
	
	err = h.forumService.ApprovePost(c.Request.Context(), postID, userID, userRoleStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Post approved successfully"})
}

func (h *ForumHandler) RejectPost(c *gin.Context) {
	postIDStr := c.Param("postId")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}
	
	userIDInterface, _ := c.Get("user_id")
	userIDStr, _ := userIDInterface.(string)
	userID, _ := uuid.Parse(userIDStr)
	
	userRole, _ := c.Get("user_role")
	userRoleStr, _ := userRole.(string)
	
	err = h.forumService.RejectPost(c.Request.Context(), postID, userID, userRoleStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Post rejected successfully"})
}

func (h *ForumHandler) SetTopicPinOrder(c *gin.Context) {
	topicIDStr := c.Param("topicId")
	topicID, err := uuid.Parse(topicIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid topic ID"})
		return
	}
	
	var req struct {
		PinOrder *int `json:"pin_order"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	userRole, _ := c.Get("user_role")
	userRoleStr, _ := userRole.(string)
	
	err = h.forumService.SetTopicPinOrder(c.Request.Context(), topicID, req.PinOrder, userRoleStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Pin order updated successfully"})
}

func (h *ForumHandler) SetPostPinOrder(c *gin.Context) {
	postIDStr := c.Param("postId")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}
	
	var req struct {
		PinOrder *int `json:"pin_order"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	userRole, _ := c.Get("user_role")
	userRoleStr, _ := userRole.(string)
	
	err = h.forumService.SetPostPinOrder(c.Request.Context(), postID, req.PinOrder, userRoleStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Pin order updated successfully"})
}
