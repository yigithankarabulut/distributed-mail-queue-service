package model

import "gorm.io/gorm"

// User struct
type User struct {
	gorm.Model
	Email        string `gorm:"unique"`
	Password     string `gorm:"not null"`
	SmtpHost     string `gorm:"not null"`
	SmtpPort     int    `gorm:"not null"`
	SmtpUsername string `gorm:"not null"`
	SmtpPassword string `gorm:"not null"`
}
