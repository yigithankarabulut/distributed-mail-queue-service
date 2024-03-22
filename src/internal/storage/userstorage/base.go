package userstorage

import (
	"context"
	"github.com/yigithankarabulut/distributed-mail-queue-service/src/model"
	"gorm.io/gorm"
)

// UserStorer is an interface for storing users.
type UserStorer interface {
	Insert(ctx context.Context, user model.User, tx ...*gorm.DB) error
	GetByID(ctx context.Context, id uint) (model.User, error)
	GetByEmail(ctx context.Context, email string) (model.User, error)
	Update(ctx context.Context, user model.User, tx ...*gorm.DB) error
	CreateTx() *gorm.DB
	CommitTx(tx *gorm.DB)
	RollbackTx(tx *gorm.DB)
	SetTx(tx ...*gorm.DB) *gorm.DB
}

// userStorage is a storage for users.
type userStorage struct {
	db *gorm.DB
}

// Option is a type for user storage options.
type Option func(*userStorage)

// WithUserDB sets the database for user storage.
func WithUserDB(db *gorm.DB) Option {
	return func(u *userStorage) {
		u.db = db
	}
}

// New creates a new user storage.
func New(opts ...Option) UserStorer {
	u := &userStorage{}
	for _, opt := range opts {
		opt(u)
	}
	return u
}
