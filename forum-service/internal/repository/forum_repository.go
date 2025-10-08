package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"forum-service/internal/model"

	"github.com/google/uuid"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type ForumRepository struct {
	db *sql.DB
}

func NewForumRepository(db *sql.DB) *ForumRepository {
	return &ForumRepository{db: db}
}

// Topic operations
func (r *ForumRepository) CreateTopic(ctx context.Context, topic *model.Topic) error {
	query := `
		INSERT INTO forum_topics (id, course_id, created_by_id, title, description, category, tags, is_sticky, is_locked, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err := r.db.ExecContext(ctx, query,
		topic.ID,
		topic.CourseID,
		topic.CreatedByID,
		topic.Title,
		topic.Description,
		topic.Category,
		pq.Array(topic.Tags),
		topic.IsSticky,
		topic.IsLocked,
		topic.CreatedAt,
		topic.UpdatedAt,
	)

	return err
}

func (r *ForumRepository) GetTopicByID(ctx context.Context, topicID uuid.UUID) (*model.Topic, error) {
	query := `
		SELECT id, course_id, created_by_id, title, description, category, tags, is_sticky, is_locked, 
		       view_count, post_count, last_post_at, last_post_by_id, created_at, updated_at
		FROM forum_topics
		WHERE id = $1 AND deleted_at IS NULL`

	topic := &model.Topic{}
	var courseID sql.NullString
	var lastPostAt sql.NullTime
	var lastPostByID sql.NullString

	err := r.db.QueryRowContext(ctx, query, topicID).Scan(
		&topic.ID,
		&courseID,
		&topic.CreatedByID,
		&topic.Title,
		&topic.Description,
		&topic.Category,
		pq.Array(&topic.Tags),
		&topic.IsSticky,
		&topic.IsLocked,
		&topic.ViewCount,
		&topic.PostCount,
		&lastPostAt,
		&lastPostByID,
		&topic.CreatedAt,
		&topic.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if courseID.Valid {
		courseUUID, err := uuid.Parse(courseID.String)
		if err == nil {
			topic.CourseID = &courseUUID
		}
	}

	if lastPostAt.Valid {
		topic.LastPostAt = &lastPostAt.Time
	}

	if lastPostByID.Valid {
		lastPostByUUID, err := uuid.Parse(lastPostByID.String)
		if err == nil {
			topic.LastPostByID = &lastPostByUUID
		}
	}

	return topic, nil
}

func (r *ForumRepository) ListTopics(ctx context.Context, options *model.ListTopicsOptions) ([]*model.Topic, int, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	// Base query
	baseQuery := `
		FROM forum_topics ft
		WHERE ft.deleted_at IS NULL`

	// Add filters
	if options.CourseID != nil {
		conditions = append(conditions, fmt.Sprintf("ft.course_id = $%d", argIndex))
		args = append(args, *options.CourseID)
		argIndex++
	}

	if options.Category != "" {
		conditions = append(conditions, fmt.Sprintf("ft.category = $%d", argIndex))
		args = append(args, options.Category)
		argIndex++
	}

	if len(options.Tags) > 0 {
		conditions = append(conditions, fmt.Sprintf("ft.tags && $%d", argIndex))
		args = append(args, pq.Array(options.Tags))
		argIndex++
	}

	if options.Search != "" {
		searchPattern := "%" + options.Search + "%"
		conditions = append(conditions, fmt.Sprintf("(ft.title ILIKE $%d OR ft.description ILIKE $%d)", argIndex, argIndex))
		args = append(args, searchPattern)
		argIndex++
	}

	// Build WHERE clause
	whereClause := baseQuery
	if len(conditions) > 0 {
		whereClause += " AND " + strings.Join(conditions, " AND ")
	}

	// Count query
	countQuery := "SELECT COUNT(*) " + whereClause
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Main query with sorting and pagination
	selectQuery := `
		SELECT ft.id, ft.course_id, ft.created_by_id, ft.title, ft.description, ft.category, ft.tags,
		       ft.is_sticky, ft.is_locked, ft.view_count, ft.post_count, ft.last_post_at, ft.last_post_by_id,
		       ft.created_at, ft.updated_at ` + whereClause

	// Add sorting
	sortBy := "created_at"
	if options.SortBy != "" {
		switch options.SortBy {
		case "updated_at", "post_count", "view_count":
			sortBy = options.SortBy
		}
	}

	sortOrder := "DESC"
	if options.SortOrder == "asc" {
		sortOrder = "ASC"
	}

	selectQuery += fmt.Sprintf(" ORDER BY ft.is_sticky DESC, ft.%s %s", sortBy, sortOrder)

	// Add pagination
	limit := 20
	if options.Limit > 0 && options.Limit <= 100 {
		limit = options.Limit
	}

	offset := 0
	if options.Page > 0 {
		offset = (options.Page - 1) * limit
	}

	selectQuery += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)

	// Execute query
	rows, err := r.db.QueryContext(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var topics []*model.Topic
	for rows.Next() {
		topic := &model.Topic{}
		var courseID sql.NullString
		var lastPostAt sql.NullTime
		var lastPostByID sql.NullString

		err := rows.Scan(
			&topic.ID,
			&courseID,
			&topic.CreatedByID,
			&topic.Title,
			&topic.Description,
			&topic.Category,
			pq.Array(&topic.Tags),
			&topic.IsSticky,
			&topic.IsLocked,
			&topic.ViewCount,
			&topic.PostCount,
			&lastPostAt,
			&lastPostByID,
			&topic.CreatedAt,
			&topic.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		if courseID.Valid {
			courseUUID, err := uuid.Parse(courseID.String)
			if err == nil {
				topic.CourseID = &courseUUID
			}
		}

		if lastPostAt.Valid {
			topic.LastPostAt = &lastPostAt.Time
		}

		if lastPostByID.Valid {
			lastPostByUUID, err := uuid.Parse(lastPostByID.String)
			if err == nil {
				topic.LastPostByID = &lastPostByUUID
			}
		}

		topics = append(topics, topic)
	}

	return topics, total, nil
}

func (r *ForumRepository) UpdateTopic(ctx context.Context, topic *model.Topic) error {
	query := `
		UPDATE forum_topics
		SET title = $2, description = $3, category = $4, tags = $5, is_sticky = $6, is_locked = $7, updated_at = $8
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query,
		topic.ID,
		topic.Title,
		topic.Description,
		topic.Category,
		pq.Array(topic.Tags),
		topic.IsSticky,
		topic.IsLocked,
		topic.UpdatedAt,
	)

	return err
}

func (r *ForumRepository) DeleteTopic(ctx context.Context, topicID uuid.UUID) error {
	query := `
		UPDATE forum_topics
		SET deleted_at = $2
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, topicID, time.Now())
	return err
}

func (r *ForumRepository) IncrementTopicViewCount(ctx context.Context, topicID uuid.UUID) error {
	query := `
		UPDATE forum_topics
		SET view_count = view_count + 1
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, topicID)
	return err
}

// Post operations
func (r *ForumRepository) CreatePost(ctx context.Context, post *model.Post) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert post
	query := `
		INSERT INTO forum_posts (id, topic_id, author_id, parent_id, content, is_edited, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err = tx.ExecContext(ctx, query,
		post.ID,
		post.TopicID,
		post.AuthorID,
		post.ParentID,
		post.Content,
		post.IsEdited,
		post.CreatedAt,
		post.UpdatedAt,
	)
	if err != nil {
		return err
	}

	// Update topic stats
	updateTopicQuery := `
		UPDATE forum_topics
		SET post_count = post_count + 1, last_post_at = $2, last_post_by_id = $3, updated_at = $4
		WHERE id = $1`

	_, err = tx.ExecContext(ctx, updateTopicQuery,
		post.TopicID,
		post.CreatedAt,
		post.AuthorID,
		post.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *ForumRepository) GetPostByID(ctx context.Context, postID uuid.UUID) (*model.Post, error) {
	query := `
		SELECT id, topic_id, author_id, parent_id, content, is_edited, edited_at, up_votes, down_votes,
		       is_answer, is_pinned, created_at, updated_at
		FROM forum_posts
		WHERE id = $1 AND deleted_at IS NULL`

	post := &model.Post{}
	var parentID sql.NullString
	var editedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, postID).Scan(
		&post.ID,
		&post.TopicID,
		&post.AuthorID,
		&parentID,
		&post.Content,
		&post.IsEdited,
		&editedAt,
		&post.UpVotes,
		&post.DownVotes,
		&post.IsAnswer,
		&post.IsPinned,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if parentID.Valid {
		parentUUID, err := uuid.Parse(parentID.String)
		if err == nil {
			post.ParentID = &parentUUID
		}
	}

	if editedAt.Valid {
		post.EditedAt = &editedAt.Time
	}

	return post, nil
}

func (r *ForumRepository) ListPosts(ctx context.Context, options *model.ListPostsOptions) ([]*model.Post, int, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	// Base query
	baseQuery := `
		FROM forum_posts fp
		WHERE fp.deleted_at IS NULL AND fp.topic_id = $1`
	args = append(args, options.TopicID)
	argIndex++

	// Add parent filter
	if options.ParentID != nil {
		conditions = append(conditions, fmt.Sprintf("fp.parent_id = $%d", argIndex))
		args = append(args, *options.ParentID)
		argIndex++
	} else {
		conditions = append(conditions, "fp.parent_id IS NULL")
	}

	// Build WHERE clause
	whereClause := baseQuery
	if len(conditions) > 0 {
		whereClause += " AND " + strings.Join(conditions, " AND ")
	}

	// Count query
	countQuery := "SELECT COUNT(*) " + whereClause
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Main query with sorting and pagination
	selectQuery := `
		SELECT fp.id, fp.topic_id, fp.author_id, fp.parent_id, fp.content, fp.is_edited, fp.edited_at,
		       fp.up_votes, fp.down_votes, fp.is_answer, fp.is_pinned, fp.created_at, fp.updated_at ` + whereClause

	// Add sorting
	sortBy := "created_at"
	if options.SortBy != "" {
		switch options.SortBy {
		case "votes", "is_answer":
			sortBy = options.SortBy
		}
	}

	if sortBy == "votes" {
		sortBy = "(up_votes - down_votes)"
	}

	sortOrder := "ASC"
	if options.SortOrder == "desc" {
		sortOrder = "DESC"
	}

	selectQuery += fmt.Sprintf(" ORDER BY fp.is_pinned DESC, fp.is_answer DESC, %s %s", sortBy, sortOrder)

	// Add pagination
	limit := 50
	if options.Limit > 0 && options.Limit <= 100 {
		limit = options.Limit
	}

	offset := 0
	if options.Page > 0 {
		offset = (options.Page - 1) * limit
	}

	selectQuery += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)

	// Execute query
	rows, err := r.db.QueryContext(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var posts []*model.Post
	for rows.Next() {
		post := &model.Post{}
		var parentID sql.NullString
		var editedAt sql.NullTime

		err := rows.Scan(
			&post.ID,
			&post.TopicID,
			&post.AuthorID,
			&parentID,
			&post.Content,
			&post.IsEdited,
			&editedAt,
			&post.UpVotes,
			&post.DownVotes,
			&post.IsAnswer,
			&post.IsPinned,
			&post.CreatedAt,
			&post.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		if parentID.Valid {
			parentUUID, err := uuid.Parse(parentID.String)
			if err == nil {
				post.ParentID = &parentUUID
			}
		}

		if editedAt.Valid {
			post.EditedAt = &editedAt.Time
		}

		posts = append(posts, post)
	}

	return posts, total, nil
}

func (r *ForumRepository) UpdatePost(ctx context.Context, post *model.Post) error {
	query := `
		UPDATE forum_posts
		SET content = $2, is_edited = $3, edited_at = $4, is_answer = $5, is_pinned = $6, updated_at = $7
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query,
		post.ID,
		post.Content,
		post.IsEdited,
		post.EditedAt,
		post.IsAnswer,
		post.IsPinned,
		post.UpdatedAt,
	)

	return err
}

func (r *ForumRepository) DeletePost(ctx context.Context, postID uuid.UUID) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Get post info to update topic stats
	var topicID uuid.UUID
	err = tx.QueryRowContext(ctx, "SELECT topic_id FROM forum_posts WHERE id = $1", postID).Scan(&topicID)
	if err != nil {
		return err
	}

	// Soft delete post
	query := `
		UPDATE forum_posts
		SET deleted_at = $2
		WHERE id = $1`

	_, err = tx.ExecContext(ctx, query, postID, time.Now())
	if err != nil {
		return err
	}

	// Update topic post count
	updateTopicQuery := `
		UPDATE forum_topics
		SET post_count = post_count - 1
		WHERE id = $1`

	_, err = tx.ExecContext(ctx, updateTopicQuery, topicID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// Voting operations
func (r *ForumRepository) CreateOrUpdateVote(ctx context.Context, vote *model.Vote) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Check if vote exists
	var existingVoteType sql.NullString
	err = tx.QueryRowContext(ctx,
		"SELECT vote_type FROM forum_votes WHERE post_id = $1 AND user_id = $2",
		vote.PostID, vote.UserID).Scan(&existingVoteType)

	if err != nil && err != sql.ErrNoRows {
		return err
	}

	var oldVoteType *model.VoteType
	if existingVoteType.Valid {
		vt := model.VoteType(existingVoteType.String)
		oldVoteType = &vt
	}

	// Insert or update vote
	if oldVoteType == nil {
		// Insert new vote
		_, err = tx.ExecContext(ctx,
			"INSERT INTO forum_votes (id, post_id, user_id, vote_type, voted_at) VALUES ($1, $2, $3, $4, $5)",
			vote.ID, vote.PostID, vote.UserID, vote.VoteType, vote.VotedAt)
	} else {
		// Update existing vote
		_, err = tx.ExecContext(ctx,
			"UPDATE forum_votes SET vote_type = $3, voted_at = $4 WHERE post_id = $1 AND user_id = $2",
			vote.PostID, vote.UserID, vote.VoteType, vote.VotedAt)
	}

	if err != nil {
		return err
	}

	// Update post vote counts
	var upDelta, downDelta int
	if oldVoteType == nil {
		// New vote
		if vote.VoteType == model.VoteTypeUp {
			upDelta = 1
		} else {
			downDelta = 1
		}
	} else {
		// Vote change
		if *oldVoteType == model.VoteTypeUp && vote.VoteType == model.VoteTypeDown {
			upDelta = -1
			downDelta = 1
		} else if *oldVoteType == model.VoteTypeDown && vote.VoteType == model.VoteTypeUp {
			upDelta = 1
			downDelta = -1
		}
	}

	if upDelta != 0 || downDelta != 0 {
		_, err = tx.ExecContext(ctx,
			"UPDATE forum_posts SET up_votes = up_votes + $2, down_votes = down_votes + $3 WHERE id = $1",
			vote.PostID, upDelta, downDelta)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *ForumRepository) DeleteVote(ctx context.Context, postID, userID uuid.UUID) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Get existing vote type
	var voteType string
	err = tx.QueryRowContext(ctx,
		"SELECT vote_type FROM forum_votes WHERE post_id = $1 AND user_id = $2",
		postID, userID).Scan(&voteType)

	if err != nil {
		return err // Vote doesn't exist or other error
	}

	// Delete vote
	_, err = tx.ExecContext(ctx,
		"DELETE FROM forum_votes WHERE post_id = $1 AND user_id = $2",
		postID, userID)
	if err != nil {
		return err
	}

	// Update post vote counts
	var upDelta, downDelta int
	if voteType == string(model.VoteTypeUp) {
		upDelta = -1
	} else {
		downDelta = -1
	}

	_, err = tx.ExecContext(ctx,
		"UPDATE forum_posts SET up_votes = up_votes + $2, down_votes = down_votes + $3 WHERE id = $1",
		postID, upDelta, downDelta)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *ForumRepository) GetUserVote(ctx context.Context, postID, userID uuid.UUID) (*model.VoteType, error) {
	var voteType string
	err := r.db.QueryRowContext(ctx,
		"SELECT vote_type FROM forum_votes WHERE post_id = $1 AND user_id = $2",
		postID, userID).Scan(&voteType)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No vote found
		}
		return nil, err
	}

	vt := model.VoteType(voteType)
	return &vt, nil
}