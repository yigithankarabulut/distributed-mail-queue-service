package userstorage

import (
	"context"
	"github.com/yigithankarabulut/distributed-mail-queue-service/model"
	"gorm.io/gorm"
)

func (s *userStorage) Insert(ctx context.Context, user model.User, tx ...*gorm.DB) error {
	db := s.db
	if len(tx) > 0 {
		db = tx[0]
	}
	if err := db.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func (s *userStorage) GetByEmail(ctx context.Context, email string) (model.User, error) {
	var user model.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (s *userStorage) GetByID(ctx context.Context, id uint) (model.User, error) {
	var user model.User

	if err := s.db.Where("id = ?", id).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (s *userStorage) Update(ctx context.Context, user model.User, tx ...*gorm.DB) error {
	db := s.db
	if len(tx) > 0 {
		db = tx[0]
	}
	if err := db.Save(&user).Error; err != nil {
		return err
	}
	return nil
}

func (s *userStorage) Delete(ctx context.Context, id uint) error {
	if err := s.db.Where("id = ?", id).Delete(&model.User{}).Error; err != nil {
		return err
	}
	return nil
}

func (s *userStorage) CreateTx() *gorm.DB {
	return s.db.Begin()
}

func (s *userStorage) CommitTx(tx *gorm.DB) {
	tx.Commit()
}

func (s *userStorage) RollbackTx(tx *gorm.DB) {
	tx.Rollback()
}

func (s *userStorage) SetTx(tx ...*gorm.DB) *gorm.DB {
	if len(tx) > 0 {
		s.db = tx[0]
	}
	return s.db
}
