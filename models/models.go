package models

import "gorm.io/gorm"

type User struct {
	gorm.Model `json:"-"`
	ID         uint   `json:"id" gorm:"autoIncrement"`
	Username   string `json:"username" gorm:"uniqueIndex"`
	Password   string `json:"-"`
}
