package userstorage_test

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/yigithankarabulut/distributed-mail-queue-service/internal/storage/userstorage"
	"github.com/yigithankarabulut/distributed-mail-queue-service/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

func Test_userStorage_Insert(t *testing.T) {
	mockDb, mock, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, _ := gorm.Open(dialector, &gorm.Config{})
	{
		tc := "Case 1: User is valid and there is no error"
		mock.ExpectBegin()
		mock.ExpectQuery("INSERT INTO \"users\"").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()
		mock.ExpectClose()
		storage := userstorage.New(userstorage.WithUserDB(db))
		err := storage.Insert(context.Background(), model.User{
			Email:        "test",
			SmtpPassword: "test",
			SmtpUsername: "test",
			SmtpPort:     1,
			SmtpHost:     "test",
			Password:     "test",
		})
		t.Run(tc, func(t *testing.T) {
			if err != nil {
				t.Errorf("%s: Expected err to be nil but got %v", tc, err)
			}
		})
	}
	{
		tc := "Case 2: User is invalid and there is an error"
		mock.ExpectBegin()
		mock.ExpectQuery("INSERT INTO \"users\"").
			WillReturnError(gorm.ErrInvalidData)
		mock.ExpectRollback()
		mock.ExpectClose()
		storage := userstorage.New(userstorage.WithUserDB(db))
		err := storage.Insert(context.Background(), model.User{
			Email:        "test",
			SmtpPassword: "test",
		})
		t.Run(tc, func(t *testing.T) {
			if err == nil {
				t.Errorf("%s: Expected err to be not nil but got nil", tc)
			}
		})
	}
}

func Test_userStorage_GetByEmail(t *testing.T) {
	mockDb, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, _ := gorm.Open(dialector, &gorm.Config{})
	{
		tc := "Case 1: User found and there is no error"
		mock.ExpectQuery("SELECT * FROM \"users\" WHERE email = $1 AND \"users\".\"deleted_at\" IS NULL ORDER BY \"users\".\"id\" LIMIT $2").
			WithArgs("test", 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email", "smtp_password", "smtp_username", "smtp_port", "smtp_host", "password"}).
				AddRow(1, "test", "test", "test", 1, "test", "test"))
		storage := userstorage.New(userstorage.WithUserDB(db))
		_, err := storage.GetByEmail(context.Background(), "test")
		t.Run(tc, func(t *testing.T) {
			if err != nil {
				t.Errorf("%s: Expected err to be nil but got %v", tc, err)
			}
		})
	}
	{
		tc := "Case 2: User not found and there is an error"
		mock.ExpectQuery("SELECT * FROM \"users\" WHERE email = $1 AND \"users\".\"deleted_at\" IS NULL ORDER BY \"users\".\"id\" LIMIT $2").
			WithArgs("test", 1).
			WillReturnError(gorm.ErrRecordNotFound)
		storage := userstorage.New(userstorage.WithUserDB(db))
		_, err := storage.GetByEmail(context.Background(), "test")
		t.Run(tc, func(t *testing.T) {
			if err == nil {
				t.Errorf("%s: Expected err to be not nil but got nil", tc)
			}
		})
	}
}

func Test_userStorage_GetByID(t *testing.T) {
	mockDb, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, _ := gorm.Open(dialector, &gorm.Config{})
	{
		tc := "Case 1: User found and there is no error"
		mock.ExpectQuery("SELECT * FROM \"users\" WHERE id = $1 AND \"users\".\"deleted_at\" IS NULL ORDER BY \"users\".\"id\" LIMIT $2").
			WithArgs(1, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email", "smtp_password", "smtp_username", "smtp_port", "smtp_host", "password"}).
				AddRow(1, "test", "test", "test", 1, "test", "test"))
		storage := userstorage.New(userstorage.WithUserDB(db))
		_, err := storage.GetByID(context.Background(), 1)
		t.Run(tc, func(t *testing.T) {
			if err != nil {
				t.Errorf("%s: Expected err to be nil but got %v", tc, err)
			}
		})
	}
	{
		tc := "Case 2: User not found and there is an error"
		mock.ExpectQuery("SELECT * FROM \"users\" WHERE id = $1 AND \"users\".\"deleted_at\" IS NULL ORDER BY \"users\".\"id\" LIMIT $2").
			WithArgs(1, 1).
			WillReturnError(gorm.ErrRecordNotFound)
		storage := userstorage.New(userstorage.WithUserDB(db))
		_, err := storage.GetByID(context.Background(), 1)
		t.Run(tc, func(t *testing.T) {
			if err == nil {
				t.Errorf("%s: Expected err to be not nil but got nil", tc)
			}
		})
	}
}

func Test_userStorage_Update(t *testing.T) {
	mockDb, mock, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, _ := gorm.Open(dialector, &gorm.Config{})
	{
		tc := "Case 1: User is valid and there is no error"
		mock.ExpectBegin()
		mock.ExpectQuery("INSERT INTO \"users\"").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()
		mock.ExpectClose()
		storage := userstorage.New(userstorage.WithUserDB(db))
		err := storage.Update(context.Background(), model.User{
			Email:        "test",
			SmtpPassword: "test",
			SmtpUsername: "test",
			SmtpPort:     1,
			SmtpHost:     "test",
			Password:     "test",
		})
		t.Run(tc, func(t *testing.T) {
			if err != nil {
				t.Errorf("%s: Expected err to be nil but got %v", tc, err)
			}
		})
	}
	{
		tc := "Case 2: User is invalid and there is an error"
		mock.ExpectBegin()
		mock.ExpectQuery("INSERT INTO \"users\"").
			WillReturnError(gorm.ErrInvalidData)
		mock.ExpectRollback()
		mock.ExpectClose()
		storage := userstorage.New(userstorage.WithUserDB(db))
		err := storage.Update(context.Background(), model.User{
			SmtpPassword: "test",
		})
		t.Run(tc, func(t *testing.T) {
			if err == nil {
				t.Errorf("%s: Expected err to be not nil but got nil", tc)
			}
		})
	}
}

func Test_userStorage_Delete(t *testing.T) {
	mockDb, mock, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, _ := gorm.Open(dialector, &gorm.Config{})
	{
		tc := "Case 1: User is found and there is no error"
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE \"users\" SET \"deleted_at\"=").
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
		mock.ExpectClose()
		storage := userstorage.New(userstorage.WithUserDB(db))
		err := storage.Delete(context.Background(), 1)
		t.Run(tc, func(t *testing.T) {
			if err != nil {
				t.Errorf("%s: Expected err to be nil but got %v", tc, err)
			}
		})
	}
	{
		tc := "Case 2: User is not found and there is an error"
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE \"users\" SET \"deleted_at\"=").
			WillReturnError(gorm.ErrRecordNotFound)
		mock.ExpectRollback()
		mock.ExpectClose()
		storage := userstorage.New(userstorage.WithUserDB(db))
		err := storage.Delete(context.Background(), 0)
		t.Run(tc, func(t *testing.T) {
			if err == nil {
				t.Errorf("%s: Expected err to be not nil but got nil", tc)
			}
		})
	}
}

func Test_userStorage_CommitTx(t *testing.T) {
	mockDb, mock, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, _ := gorm.Open(dialector, &gorm.Config{})
	tc := "Case 1: Commit transaction"
	storage := userstorage.New(userstorage.WithUserDB(db))
	mock.ExpectBegin()
	storage.CommitTx(db.Begin())
	t.Run(tc, func(t *testing.T) {
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("%s: Expected all expectations to be met but got %v", tc, err)
		}
	})
}

func Test_userStorage_CreateTx(t *testing.T) {
	mockDb, mock, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, _ := gorm.Open(dialector, &gorm.Config{})
	tc := "Case 1: Create transaction"
	storage := userstorage.New(userstorage.WithUserDB(db))
	mock.ExpectBegin()
	storage.CreateTx()
	t.Run(tc, func(t *testing.T) {
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("%s: Expected all expectations to be met but got %v", tc, err)
		}
	})
}

func Test_userStorage_RollbackTx(t *testing.T) {
	mockDb, mock, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, _ := gorm.Open(dialector, &gorm.Config{})
	tc := "Case 1: Rollback transaction"
	storage := userstorage.New(userstorage.WithUserDB(db))
	mock.ExpectBegin()
	storage.RollbackTx(db.Begin())
	t.Run(tc, func(t *testing.T) {
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("%s: Expected all expectations to be met but got %v", tc, err)
		}
	})
}

func Test_userStorage_SetTx(t *testing.T) {
	mockDb, mock, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, _ := gorm.Open(dialector, &gorm.Config{})
	tc := "Case 1: Set transaction"
	storage := userstorage.New(userstorage.WithUserDB(db))
	mock.ExpectBegin()
	storage.SetTx(db.Begin())
	t.Run(tc, func(t *testing.T) {
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("%s: Expected all expectations to be met but got %v", tc, err)
		}
	})
}
