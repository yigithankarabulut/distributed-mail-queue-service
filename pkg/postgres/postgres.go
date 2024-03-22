package postgres

import (
	"database/sql"
	"fmt"
	"github.com/yigithankarabulut/distributed-mail-queue-service/config"
	"github.com/yigithankarabulut/distributed-mail-queue-service/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

// ConnectPQ connects to the postgres database and returns the connection instance.
func ConnectPQ(config config.Database) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", config.Host, config.Port, config.User, config.Name, config.Pass)
	sqlDB, err := sql.Open("pgx", dsn)
	if DB, err = gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{}); err != nil {
		return nil, err
	}
	if ping, err := DB.DB(); err != nil || ping.Ping() != nil {
		return nil, err
	}
	if config.Migrate {
		err = AutoMigrate()
		if err != nil {
			return nil, err
		}
	}
	log.Printf("Connected to postgres at %s:%s", config.Host, config.Port)
	return DB, nil
}

// AutoMigrate migrates the models to the database.
func AutoMigrate() error {
	err := DB.AutoMigrate(
		&model.User{},
		&model.MailTaskQueue{},
	)
	if err != nil {
		return err
	}
	return nil
}
