package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("resource not found")
	QueryTimeoutDuration = time.Second * 5
)

type Storage struct {
	Comments interface {
		Create(ctx context.Context, c *Comment) error
		GetByPostID(ctx context.Context, postId int64) ([]Comment, error)
	}
	Posts interface {
		Create(ctx context.Context, p *Post) error
		DeleteByID(ctz context.Context, id int64) error
		GetByID(ctx context.Context, id int64) (*Post, error)
		Update(ctx context.Context, p *Post) error
	}
	Users interface {
		Create(ctx context.Context, u *User) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Comments: &CommentStore{db},
		Posts:    &PostStore{db},
		Users:    &UserStore{db},
	}
}
