package model

import (
	"gorm.io/gorm"
)

type UserData struct {
	gorm.Model
	Name           string
	Email          string `gorm:"unique"`
	Password       string
	Phone_no       string
	Is_spam        bool
	Reported_count int
}
