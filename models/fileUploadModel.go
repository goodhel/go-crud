package models

import (
	"time"

	"gorm.io/gorm"
)

type FileUpload struct {
	ID           uint           `gorm:"primarykey"`
	Name         string         `gorm:"size:100;" json:"name"`
	OriginalName string         `gorm:"size:150;" json:"original_name"`
	Mime         string         `gorm:"size:75;" json:"mime"`
	Path         string         `gorm:"size:200;" json:"path"`
	Extension    string         `gorm:"size:20;" json:"extension"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
