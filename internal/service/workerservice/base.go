package workerservice

import (
	"github.com/yigithankarabulut/distributed-mail-queue-service/internal/storage/taskqueue"
	"github.com/yigithankarabulut/distributed-mail-queue-service/model"
	"gorm.io/gorm"
)

type IWorker interface {
	TriggerWorker()
}

type worker struct {
	id        uint32
	db        *gorm.DB
	taskqueue taskqueue.TaskQueue
	ch        chan model.MailTaskQueue
}

type Option func(*worker)

func WithDB(db *gorm.DB) Option {
	return func(w *worker) {
		w.db = db
	}
}

func WithTaskQueue(rds taskqueue.TaskQueue) Option {
	return func(w *worker) {
		w.taskqueue = rds
	}
}

func WithChannel(ch chan model.MailTaskQueue) Option {
	return func(w *worker) {
		w.ch = ch
	}
}

func New(id int, opts ...Option) IWorker {
	w := &worker{id: uint32(id)}
	for _, opt := range opts {
		opt(w)
	}
	return w
}
