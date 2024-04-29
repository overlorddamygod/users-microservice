package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID    int32  `gorm:"primaryKey;autoIncrement"`
	Name  string `binding:"required,min=3"`
	Email string `binding:"required,email"`
}
