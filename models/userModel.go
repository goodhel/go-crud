package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email      string       `gorm:"unique" json:"email"`
	Password   string       `json:"password"`
	Name       string       `json:"name"`
	Avatar     *string      `json:"avatar"`
	Role       []Role       `gorm:"many2many:user_roles;" json:"roles"`
	FileUpload []FileUpload `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
