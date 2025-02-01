package models

import "gorm.io/gorm"

type Staff struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Hospital string `gorm:"not null"`
}
