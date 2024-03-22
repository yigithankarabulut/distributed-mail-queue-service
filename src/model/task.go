package model

import "gorm.io/gorm"

// MailTaskQueue is a struct that represent the mail task queue table in the database.
type MailTaskQueue struct {
	gorm.Model
	UserID         uint
	User           User   `gorm:"foreignKey:UserID"`
	RecipientEmail string `gorm:"not null"`
	Subject        string
	Body           string
	Status         string
	TryCount       int // max 5
}
