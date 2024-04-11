package taskstorage

import (
	"context"
	"github.com/yigithankarabulut/distributed-mail-queue-service/model"
	"gorm.io/gorm"
)

// TaskStorer is an interface for storing mail tasks
type TaskStorer interface {
	Insert(ctx context.Context, task model.MailTaskQueue, tx ...*gorm.DB) (model.MailTaskQueue, error)
	GetByID(ctx context.Context, id uint) (model.MailTaskQueue, error)
	GetAll(ctx context.Context, userID uint) ([]model.MailTaskQueue, error)
	GetAllByUnprocessedTasks(ctx context.Context) ([]model.MailTaskQueue, error)
	GetAllByStatusWithUserID(ctx context.Context, state int, userID uint) ([]model.MailTaskQueue, error)
	Update(ctx context.Context, task model.MailTaskQueue, tx ...*gorm.DB) error
	Delete(ctx context.Context, id uint) error
}

// taskStorage is a storage for mail tasks
type taskStorage struct {
	db *gorm.DB
}

// Option is a type for mail task storage options
type Option func(storage *taskStorage)

// WithTaskDB sets the database for mail task storage.
func WithTaskDB(db *gorm.DB) Option {
	return func(storage *taskStorage) {
		storage.db = db
	}
}

// New creates a new mail task storage instance.
func New(opts ...Option) TaskStorer {
	storage := &taskStorage{}
	for _, opt := range opts {
		opt(storage)
	}
	return storage
}
