package taskstorage

import (
	"context"
	"github.com/yigithankarabulut/distributed-mail-queue-service/model"
	"gorm.io/gorm"
)

func (s *taskStorage) Insert(ctx context.Context, task model.MailTaskQueue, tx ...*gorm.DB) (model.MailTaskQueue, error) {
	db := s.db
	if len(tx) > 0 {
		db = tx[0]
	}
	if err := db.Create(&task).Error; err != nil {
		return task, err
	}
	return task, nil
}

func (s *taskStorage) GetByID(ctx context.Context, id uint) (model.MailTaskQueue, error) {
	var task model.MailTaskQueue
	if err := s.db.Where("id = ?", id).First(&task).Error; err != nil {
		return task, err
	}
	return task, nil
}

func (s *taskStorage) GetAll(ctx context.Context, userID uint) ([]model.MailTaskQueue, error) {
	var tasks []model.MailTaskQueue
	if err := s.db.Where("user_id = ?", userID).Find(&tasks).Error; err != nil {
		return tasks, err
	}
	return tasks, nil
}

func (s *taskStorage) GetAllByStatus(ctx context.Context, state int) ([]model.MailTaskQueue, error) {
	var tasks []model.MailTaskQueue
	if err := s.db.Where("status = ?", state).Find(&tasks).Error; err != nil {
		return tasks, err
	}
	return tasks, nil
}

func (s *taskStorage) GetAllByStatusWithUserID(ctx context.Context, state int, userID uint) ([]model.MailTaskQueue, error) {
	var tasks []model.MailTaskQueue
	if err := s.db.Where("status = ? AND user_id = ?", state, userID).Find(&tasks).Error; err != nil {
		return tasks, err
	}
	return tasks, nil
}

func (s *taskStorage) Update(ctx context.Context, task model.MailTaskQueue, tx ...*gorm.DB) error {
	db := s.db
	if len(tx) > 0 {
		db = tx[0]
	}
	if err := db.Save(&task).Error; err != nil {
		return err
	}
	return nil
}

func (s *taskStorage) Delete(ctx context.Context, id uint) error {
	if err := s.db.Where("id = ?", id).Delete(&model.MailTaskQueue{}).Error; err != nil {
		return err
	}
	return nil
}

func (s *taskStorage) CreateTx() *gorm.DB {
	return s.db.Begin()
}

func (s *taskStorage) CommitTx(tx *gorm.DB) {
	tx.Commit()
}

func (s *taskStorage) RollbackTx(tx *gorm.DB) {
	tx.Rollback()
}

func (s *taskStorage) SetTx(tx ...*gorm.DB) *gorm.DB {
	if len(tx) > 0 {
		s.db = tx[0]
	}
	return s.db
}
