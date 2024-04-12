package userservice_test

import (
	"context"
	"github.com/yigithankarabulut/distributed-mail-queue-service/internal/service/mailservice"
	"github.com/yigithankarabulut/distributed-mail-queue-service/model"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
	"time"
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

type mockUserStorer struct {
	errInsert     error
	errGetByID    error
	errGetByEmail error
	errUpdate     error
	errDelete     error
	userModel     model.User
}

func (m *mockUserStorer) Insert(ctx context.Context, user model.User, tx ...*gorm.DB) error {
	return m.errInsert
}

func (m *mockUserStorer) GetByID(ctx context.Context, id uint) (model.User, error) {
	return m.userModel, m.errGetByID
}

func (m *mockUserStorer) GetByEmail(ctx context.Context, email string) (model.User, error) {
	return m.userModel, m.errGetByEmail
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

type mockPassUtils struct {
	errHashPassword    error
	hashPasswordResult string
	errComparePassword error
}

func (m *mockPassUtils) HashPassword(password string) (string, error) {
	return m.hashPasswordResult, m.errHashPassword
}

func (m *mockPassUtils) ComparePassword(hash, password string) error {
	return m.errComparePassword
}

type mockJwtUtils struct {
	errGenerateJwtToken error
	generateJwtTokenRes string
}

func (m *mockJwtUtils) GenerateJwtToken(userID uint, expiration time.Duration) (string, error) {
	return m.generateJwtTokenRes, m.errGenerateJwtToken
}
