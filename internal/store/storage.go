package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrConflict          = errors.New("resource already exists")
	ErrNotFound          = errors.New("resource not found")
	QueryTimeoutDuration = time.Second * 5
)

type Storage struct {
	Comments interface {
		Create(ctx context.Context, c *Comment) error
		GetByPostID(ctx context.Context, postId int64) ([]Comment, error)
	}
	Followers interface {
		Follow(ctx context.Context, userToFollowId int64, followerUserId int64) error
		Unfollow(ctx context.Context, userToUnfollowId int64, followerUserId int64) error
	}
	Posts interface {
		Create(ctx context.Context, p *Post) error
		DeleteByID(ctz context.Context, id int64) error
		GetByID(ctx context.Context, id int64) (*Post, error)
		GetUserFeed(ctx context.Context, userId int64, fq PaginationFeedQuery) ([]PostWithMetadata, error)
		Update(ctx context.Context, p *Post) error
	}
	Users interface {
		Activate(ctx context.Context, token string) error
		Create(ctx context.Context, u *User, tx *sql.Tx) error
		CreateAndInvite(ctz context.Context, u *User, token string, invitationExp time.Duration) error
		Delete(ctx context.Context, userId int64) error
		GetByID(ctx context.Context, id int64) (*User, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Comments:  &CommentStore{db},
		Followers: &FollowersStore{db},
		Posts:     &PostStore{db},
		Users:     &UserStore{db},
	}
}

func withTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
