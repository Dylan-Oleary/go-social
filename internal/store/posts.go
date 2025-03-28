package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type Post struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	UserID    int64     `json:"user_id"`
	Tags      []string  `json:"tags"`
	Version   int       `json:"version"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Comments  []Comment `json:"comments"`
	User      User      `json:"user"`
}

type PostWithMetadata struct {
	Post
	CommentCount int `json:"comments_count"`
}

type PostStore struct {
	db *sql.DB
}

func (s *PostStore) Create(ctx context.Context, post *Post) error {
	query := `
        INSERT INTO posts (content, title, user_id, tags)
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at, updated_at
    `

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Content,
		post.Title,
		post.UserID,
		pq.Array(post.Tags),
	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostStore) GetByID(ctx context.Context, id int64) (*Post, error) {
	query := `
        SELECT id, user_id, content, title, tags, created_at, updated_at, version
        FROM posts 
        WHERE id = $1
    `

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var post Post

	err := s.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&post.ID,
		&post.UserID,
		&post.Content,
		&post.Title,
		pq.Array(&post.Tags),
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &post, nil
}

func (s *PostStore) GetUserFeed(ctx context.Context, userId int64) ([]PostWithMetadata, error) {
	query := `
        SELECT
            p.id, p.user_id, p.title, p.content, p.created_at, p.version, p.tags,
            u.username,
            COUNT(c.id) as comments_count
        FROM posts p
        LEFT JOIN comments c ON c.post_id = p.id
        LEFT JOIN users u ON u.id = p.user_id
        LEFT JOIN followers f ON f.follower_id = $1 AND f.user_id = p.user_id
        WHERE f.follower_id = $1 OR p.user_id = $1
        GROUP BY p.id, u.username
        ORDER BY p.created_at DESC;
    `

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var feed []PostWithMetadata
	for rows.Next() {
		var post PostWithMetadata
		err := rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Title,
			&post.Content,
			&post.CreatedAt,
			&post.Version,
			pq.Array(&post.Tags),
			&post.User.Username,
			&post.CommentCount,
		)

		if err != nil {
			return nil, err
		}

		feed = append(feed, post)
	}

	return feed, nil
}

func (s *PostStore) DeleteByID(ctx context.Context, id int64) error {
	query := "DELETE FROM posts p WHERE p.id = $1"

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *PostStore) Update(ctx context.Context, p *Post) error {
	query := `
        UPDATE posts p
        SET title = $2, content = $3, version = p.version + 1
        WHERE p.id = $1
        AND p.version = $4
        RETURNING p.version
    `

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		p.ID,
		p.Title,
		p.Content,
		p.Version,
	).Scan(&p.Version)
	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			return ErrNotFound
		default:
			return err
		}
	}

	return nil
}
