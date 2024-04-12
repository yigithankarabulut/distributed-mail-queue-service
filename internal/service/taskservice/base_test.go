package taskservice_test

import (
	"context"
	"github.com/yigithankarabulut/distributed-mail-queue-service/model"
	"gorm.io/gorm"
)

type mockTaskStorer struct {
	errInsert                   error
	errGetByID                  error
	errGetAll                   error
	errGetAllByUnprocessedTasks error
	errGetAllByStatusWithUserID error
	errUpdate                   error
	errDelete                   error
	taskModelArr                []model.MailTaskQueue
}

func (m *mockTaskStorer) Insert(ctx context.Context, task model.MailTaskQueue, tx ...*gorm.DB) (model.MailTaskQueue, error) {
	return task, m.errInsert
}

func (m *mockTaskStorer) GetByID(ctx context.Context, id uint) (model.MailTaskQueue, error) {
	return model.MailTaskQueue{}, m.errGetByID
}

func (m *mockTaskStorer) GetAll(ctx context.Context, userID uint) ([]model.MailTaskQueue, error) {
	return nil, m.errGetAll
}

func (m *mockTaskStorer) GetAllByUnprocessedTasks(ctx context.Context) ([]model.MailTaskQueue, error) {
	return m.taskModelArr, m.errGetAllByUnprocessedTasks
}

func (m *mockTaskStorer) GetAllByStatusWithUserID(ctx context.Context, state int, userID uint) ([]model.MailTaskQueue, error) {
	return nil, m.errGetAllByStatusWithUserID
}

func (m *mockTaskStorer) Update(ctx context.Context, task model.MailTaskQueue, tx ...*gorm.DB) error {
	return m.errUpdate
}

func (m *mockTaskStorer) Delete(ctx context.Context, id uint) error {
	return m.errDelete
}

type mockUserStorer struct {
	errInsert     error
	errGetByID    error
	errGetByEmail error
	errUpdate     error
	errDelete     error
}

func (m *mockUserStorer) Insert(ctx context.Context, user model.User, tx ...*gorm.DB) error {
	return m.errInsert
}

func (m *mockUserStorer) GetByID(ctx context.Context, id uint) (model.User, error) {
	return model.User{}, m.errGetByID
}

func (m *mockUserStorer) GetByEmail(ctx context.Context, email string) (model.User, error) {
	return model.User{}, m.errGetByEmail
}

func (m *mockUserStorer) Update(ctx context.Context, user model.User, tx ...*gorm.DB) error {
	return m.errUpdate
}

func (m *mockUserStorer) Delete(ctx context.Context, id uint) error {
	return m.errDelete
}

func (m *mockUserStorer) CreateTx() *gorm.DB {
	return nil
}

func (m *mockUserStorer) CommitTx(tx *gorm.DB) {

}

func (m *mockUserStorer) RollbackTx(tx *gorm.DB) {

}

func (m *mockUserStorer) SetTx(tx ...*gorm.DB) *gorm.DB {
	return nil
}

type mockTaskQueue struct {
	errPublishTask   error
	errSubscribeTask error
	errStartConsume  <-chan error
}

func (m *mockTaskQueue) PublishTask(ctx context.Context, task interface{}) error {
	return m.errPublishTask
}

func (m *mockTaskQueue) SubscribeTask(ctx context.Context, consumerID int) error {
	return m.errSubscribeTask
}

func (m *mockTaskQueue) StartConsume(ctx context.Context) <-chan error {
	return m.errStartConsume
}
