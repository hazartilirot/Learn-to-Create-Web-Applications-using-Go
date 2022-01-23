package models

import "gorm.io/gorm"

type Users struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null; uniqueIndex"`
}
