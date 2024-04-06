package model

import (
	"gorm.io/gorm"
	"time"
)

// MailTaskQueue is a struct that represent the mail task queue table in the database.
type MailTaskQueue struct {
	gorm.Model
	UserID         uint
	User           User   `gorm:"foreignKey:UserID"`
	Status         int    `gorm:"default:0"`
	TryCount       int    `gorm:"default:0"`
	RecipientEmail string `gorm:"not null"`
	Subject        string
	Body           string
	ScheduledAt    time.Time
}
