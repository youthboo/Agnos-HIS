package models

import (
	"time"
	"gorm.io/gorm"
)

type Patient struct {
	gorm.Model
	FirstNameTH  string    `gorm:"not null"`
	MiddleNameTH string
	LastNameTH   string    `gorm:"not null"`
	FirstNameEN  string
	MiddleNameEN string
	LastNameEN   string
	DateOfBirth  time.Time `gorm:"not null"` 
	PatientHN    *string   `gorm:"unique"`
	NationalID   *string    `gorm:"unique"`
	PassportID   *string   `gorm:"unique"`
	PhoneNumber  string    `gorm:"not null"`
	Email        string    `gorm:"unique"`
	Gender       string    `gorm:"not null;check:gender IN ('M', 'F')"` 
	Hospital     string    `gorm:"not null"`
}
