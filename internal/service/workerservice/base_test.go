package workerservice_test

import (
	"context"
	"github.com/yigithankarabulut/distributed-mail-queue-service/internal/service/mailservice"
	"github.com/yigithankarabulut/distributed-mail-queue-service/model"
	"gopkg.in/gomail.v2"
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
	taskModel                   model.MailTaskQueue
}

func (m *mockTaskStorer) Insert(ctx context.Context, task model.MailTaskQueue, tx ...*gorm.DB) (model.MailTaskQueue, error) {
	return m.taskModel, m.errInsert
}

func (m *mockTaskStorer) GetByID(ctx context.Context, id uint) (model.MailTaskQueue, error) {
	return m.taskModel, m.errGetByID
}

func (m *mockTaskStorer) GetAll(ctx context.Context, userID uint) ([]model.MailTaskQueue, error) {
	return m.taskModelArr, m.errGetAll
}

func (m *mockTaskStorer) GetAllByUnprocessedTasks(ctx context.Context) ([]model.MailTaskQueue, error) {
	return m.taskModelArr, m.errGetAllByUnprocessedTasks
}

func (m *mockTaskStorer) GetAllByStatusWithUserID(ctx context.Context, state int, userID uint) ([]model.MailTaskQueue, error) {
	return m.taskModelArr, m.errGetAllByStatusWithUserID
}

func (m *mockTaskStorer) Update(ctx context.Context, task model.MailTaskQueue, tx ...*gorm.DB) error {
	return m.errUpdate
}

func (m *mockTaskStorer) Delete(ctx context.Context, id uint) error {
	return m.errDelete
}

type mockMailService struct {
	errAddTask  error
	errSendMail error
}

func (m *mockMailService) AddTask(task model.MailTaskQueue) error {
	return m.errAddTask
}

func (m *mockMailService) NewDialer() *gomail.Dialer {
	return nil
}

func (m *mockMailService) NewMessage() *gomail.Message {
	return nil
}

func (m *mockMailService) SendMail(d mailservice.Dialer, me *gomail.Message) error {
	return m.errSendMail
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
