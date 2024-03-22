package model

import "gorm.io/gorm"

// MailQueue struct
type MailQueue struct {
	gorm.Model
	UserID         uint
	User           User   `gorm:"foreignKey:UserID"`
	RecipientEmail string `gorm:"not null"`
	Subject        string
	Body           string
	Status         string
	TryCount       int
}
