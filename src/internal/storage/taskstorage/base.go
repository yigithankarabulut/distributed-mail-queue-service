package taskstorage

import "gorm.io/gorm"

// TaskStorer is an interface for storing mail tasks
type TaskStorer interface {
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
