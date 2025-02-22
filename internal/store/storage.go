package store

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrNotFound = errors.New("resource not found")
)

type Storage struct {
	Comments interface {
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
