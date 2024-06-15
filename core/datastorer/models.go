package datastorer

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Login           string `json:"login" gorm:"unique; not null"`
	Username        string `json:"username" gorm:"not null"`
	Password        string `json:"password" gorm:"not null"`
	PermissionLevel int    `json:"permission_level" gorm:"not null; default:1"`
}

func (User) TableName() string {
	return "users"
}

type Device struct {
	gorm.Model
	DisplayName string `json:"display_name" gorm:"not null; default:''"`
	ClientID    string `json:"client_id" gorm:"unique; not null"`
	Password    string `json:"password" gorm:"not null"`
	Module      string `json:"plugin" gorm:"not null"`
}

func (Device) TableName() string {
	return "devices"
}