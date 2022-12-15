package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `gorm:"unique" json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Role     []Role `gorm:"many2many:user_roles;" json:"roles"`
}
