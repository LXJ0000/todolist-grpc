package domain

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserName string `gorm:"unique"`
	NickName string
	Password string
}
