package taskstorage_test

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/yigithankarabulut/distributed-mail-queue-service/internal/storage/taskstorage"
	"github.com/yigithankarabulut/distributed-mail-queue-service/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

func Test_taskStorage_Insert(t *testing.T) {
	mockDb, mock, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, _ := gorm.Open(dialector, &gorm.Config{})
	{
		tc := "Case 1: Valid Case And Success"
		mock.ExpectBegin()
		mock.ExpectQuery("INSERT INTO \"mail_task_queues\"").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()
		mock.ExpectClose()
		storage := taskstorage.New(taskstorage.WithTaskDB(db))
		_, err := storage.Insert(context.Background(), model.MailTaskQueue{})
		if err != nil {
			t.Errorf("TaskStorage.Insert() %s error = %v, want nil", tc, err)
		}
	}
	{
		tc := "Case 2: Invalid Case And Error"
		mock.ExpectBegin()
		mock.ExpectQuery("INSERT INTO \"mail_task_queues\"").
			WillReturnError(gorm.ErrInvalidData)
		mock.ExpectRollback()
		mock.ExpectClose()
		storage := taskstorage.New(taskstorage.WithTaskDB(db))
		_, err := storage.Insert(context.Background(), model.MailTaskQueue{})
		if err == nil {
			t.Errorf("TaskStorage.Insert() %s error = nil, want error", tc)
		}
	}
}

func Test_taskStorage_GetByID(t *testing.T) {
	mockDb, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, _ := gorm.Open(dialector, &gorm.Config{})
	{
		tc := "Case 1: Valid Case And Success"
		mock.ExpectQuery("SELECT * FROM \"mail_task_queues\" WHERE id = $1 AND \"mail_task_queues\".\"deleted_at\" IS NULL ORDER BY \"mail_task_queues\".\"id\" LIMIT $2").
			WithArgs(1, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		storage := taskstorage.New(taskstorage.WithTaskDB(db))
		_, err := storage.GetByID(context.Background(), 1)
		if err != nil {
			t.Errorf("TaskStorage.GetByID() %s error = %v, want nil", tc, err)
		}
	}
	{
		tc := "Case 2: Invalid Case And Error"
		mock.ExpectQuery("SELECT * FROM \"mail_task_queues\" WHERE id = $1 AND \"mail_task_queues\".\"deleted_at\" IS NULL ORDER BY \"mail_task_queues\".\"id\" LIMIT $2").
			WithArgs(1, 1).
			WillReturnError(gorm.ErrInvalidData)
		storage := taskstorage.New(taskstorage.WithTaskDB(db))
		_, err := storage.GetByID(context.Background(), 1)
		if err == nil {
			t.Errorf("TaskStorage.GetByID() %s error = nil, want error", tc)
		}
	}
	{
		tc := "Case 3: Invalid Case And Not Found"
		mock.ExpectQuery("SELECT * FROM \"mail_task_queues\" WHERE id = $1 AND \"mail_task_queues\".\"deleted_at\" IS NULL ORDER BY \"mail_task_queues\".\"id\" LIMIT $2").
			WithArgs(1, 1).
			WillReturnError(gorm.ErrRecordNotFound)
		storage := taskstorage.New(taskstorage.WithTaskDB(db))
		_, err := storage.GetByID(context.Background(), 1)
		if err == nil {
			t.Errorf("TaskStorage.GetByID() %s error = nil, want error", tc)
		}
	}
}

func Test_taskStorage_GetAll(t *testing.T) {
	mockDb, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, _ := gorm.Open(dialector, &gorm.Config{})
	{
		tc := "Case 1: Valid Case And Success"
		mock.ExpectQuery("SELECT * FROM \"mail_task_queues\" WHERE user_id = $1 AND \"mail_task_queues\".\"deleted_at\" IS NULL").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		storage := taskstorage.New(taskstorage.WithTaskDB(db))
		_, err := storage.GetAll(context.Background(), 1)
		if err != nil {
			t.Errorf("TaskStorage.GetAll() %s error = %v, want nil", tc, err)
		}
	}
	{
		tc := "Case 2: Invalid Case And Error"
		mock.ExpectQuery("SELECT * FROM \"mail_task_queues\" WHERE user_id = $1 AND \"mail_task_queues\".\"deleted_at\" IS NULL").
			WithArgs(1).
			WillReturnError(gorm.ErrInvalidData)
		storage := taskstorage.New(taskstorage.WithTaskDB(db))
		_, err := storage.GetAll(context.Background(), 1)
		if err == nil {
			t.Errorf("TaskStorage.GetAll() %s error = nil, want error", tc)
		}
	}
}

func Test_taskStorage_GetAllByUnprocessedTasks(t *testing.T) {
	mockDb, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, _ := gorm.Open(dialector, &gorm.Config{})
	{
		tc := "Case 1: Valid Case And Success"
		mock.ExpectQuery("SELECT * FROM \"mail_task_queues\" WHERE (status = $1 AND updated_at < NOW() - INTERVAL '5 minutes') AND \"mail_task_queues\".\"deleted_at\" IS NULL").
			WithArgs(0).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		storage := taskstorage.New(taskstorage.WithTaskDB(db))
		_, err := storage.GetAllByUnprocessedTasks(context.Background())
		if err != nil {
			t.Errorf("TaskStorage.GetAllByUnprocessedTasks() %s error = %v, want nil", tc, err)
		}
	}
	{
		tc := "Case 2: Wrong Status Value And Error"
		mock.ExpectQuery("SELECT * FROM \"mail_task_queues\" WHERE (status = $1 AND updated_at < NOW() - INTERVAL '5 minutes') AND \"mail_task_queues\".\"deleted_at\" IS NULL").
			WithArgs(1).
			WillReturnError(gorm.ErrInvalidData)
		storage := taskstorage.New(taskstorage.WithTaskDB(db))
		_, err := storage.GetAllByUnprocessedTasks(context.Background())
		if err == nil {
			t.Errorf("TaskStorage.GetAllByUnprocessedTasks() %s error = nil, want error", tc)
		}
	}
}

func Test_taskStorage_GetAllByStatusWithUserID(t *testing.T) {
	mockDb, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, _ := gorm.Open(dialector, &gorm.Config{})
	{
		tc := "Case 1: Valid Case And Success"
		mock.ExpectQuery("SELECT * FROM \"mail_task_queues\" WHERE (status = $1 AND user_id = $2) AND \"mail_task_queues\".\"deleted_at\" IS NULL").
			WithArgs(0, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		storage := taskstorage.New(taskstorage.WithTaskDB(db))
		_, err := storage.GetAllByStatusWithUserID(context.Background(), 0, 1)
		if err != nil {
			t.Errorf("TaskStorage.GetAllByStatusWithUserID() %s error = %v, want nil", tc, err)
		}
	}
	{
		tc := "Case 2: Invalid Case And Error"
		mock.ExpectQuery("SELECT * FROM \"mail_task_queues\" WHERE (status = $1 AND user_id = $2) AND \"mail_task_queues\".\"deleted_at\" IS NULL").
			WithArgs(1, 1).
			WillReturnError(gorm.ErrInvalidData)
		storage := taskstorage.New(taskstorage.WithTaskDB(db))
		_, err := storage.GetAllByStatusWithUserID(context.Background(), 1, 1)
		if err == nil {
			t.Errorf("TaskStorage.GetAllByStatusWithUserID() %s error = nil, want error", tc)
		}
	}
}

func Test_taskStorage_Update(t *testing.T) {
	mockDb, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, _ := gorm.Open(dialector, &gorm.Config{})
	{
		tc := "Case 1: Valid Case And Success"
		mock.ExpectBegin()
		mock.ExpectQuery("INSERT INTO \"mail_task_queues\" (\"created_at\",\"updated_at\",\"deleted_at\",\"user_id\",\"status\",\"try_count\",\"recipient_email\",\"subject\",\"body\",\"scheduled_at\") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING \"id\"").
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()
		mock.ExpectClose()
		storage := taskstorage.New(taskstorage.WithTaskDB(db))
		err := storage.Update(context.Background(), model.MailTaskQueue{
			UserID:   1,
			TryCount: 1,
		})
		if err != nil {
			t.Errorf("TaskStorage.Update() %s error = %v, want nil", tc, err)
		}
	}
	{
		tc := "Case 2: Invalid Case And Error"
		mock.ExpectBegin()
		mock.ExpectQuery("INSERT INTO \"mail_task_queues\" (\"created_at\",\"updated_at\",\"deleted_at\",\"user_id\",\"status\",\"try_count\",\"recipient_email\",\"subject\",\"body\",\"scheduled_at\") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING \"id\"").
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnError(gorm.ErrInvalidData)
		mock.ExpectRollback()
		mock.ExpectClose()
		storage := taskstorage.New(taskstorage.WithTaskDB(db))
		err := storage.Update(context.Background(), model.MailTaskQueue{
			UserID:   1,
			TryCount: 1,
		})
		if err == nil {
			t.Errorf("TaskStorage.Update() %s error = nil, want error", tc)
		}
	}
}

func Test_taskStorage_Delete(t *testing.T) {
	mockDb, mock, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, _ := gorm.Open(dialector, &gorm.Config{})
	{
		tc := "Case 1: Valid Case And Success"
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE \"mail_task_queues\" SET \"deleted_at\"=\\$1 WHERE id = \\$2 AND \"mail_task_queues\".\"deleted_at\" IS NULL").
			WithArgs(sqlmock.AnyArg(), 1).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()
		mock.ExpectClose()
		storage := taskstorage.New(taskstorage.WithTaskDB(db))
		err := storage.Delete(context.Background(), 1)
		if err != nil {
			t.Errorf("TaskStorage.Delete() %s error = %v, want nil", tc, err)
		}
	}
	{
		tc := "Case 2: Invalid Case And Error"
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE \"mail_task_queues\" SET \"deleted_at\"=\\$1 WHERE id = \\$2 AND \"mail_task_queues\".\"deleted_at\" IS NULL").
			WithArgs(sqlmock.AnyArg(), 1).
			WillReturnError(gorm.ErrInvalidData)
		mock.ExpectRollback()
		mock.ExpectClose()
		storage := taskstorage.New(taskstorage.WithTaskDB(db))
		err := storage.Delete(context.Background(), 1)
		if err == nil {
			t.Errorf("TaskStorage.Delete() %s error = nil, want error", tc)
		}
	}
}
