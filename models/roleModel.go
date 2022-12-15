package models

type Role struct {
	ID   uint   `gorm:"primarykey" json:"id"`
	Name string `json:"name"`
	User []User `gorm:"many2many:user_roles;" json:"users"`
}
