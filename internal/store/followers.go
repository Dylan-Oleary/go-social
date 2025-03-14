package store

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type Follower struct {
	UserID     int64  `json:"user_id"`
	FollowerID string `json:"follower_id"`
	CreatedAt  string `json:"created_at"`
}

type FollowersStore struct {
	db *sql.DB
}

func (s *FollowersStore) Follow(ctx context.Context, userToFollowId int64, followerUserId int64) error {
	query := "INSERT INTO followers (user_id, follower_id) VALUES ($1, $2)"

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, userToFollowId, followerUserId)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return ErrConflict
		}
	}

	return nil
}

func (s *FollowersStore) Unfollow(ctx context.Context, userToUnfollowId int64, followerUserId int64) error {
	query := `
        DELETE FROM followers f
        WHERE f.user_id = $1
        AND f.follower_id = $2
    `

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, userToUnfollowId, followerUserId)
	return err
}
