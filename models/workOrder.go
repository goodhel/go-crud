package models

import (
	"time"

	"gorm.io/gorm"
)

type DWorkOrder struct {
	ID           uint           `gorm:"primarykey"`
	FileUploadID uint           `json:"file_upload_id"`
	NoWorkOrder  string         `gorm:"size:100;" json:"no_work_order"`
	Customer     string         `gorm:"size: 150;" json:"customer"`
	PartNumber   string         `gorm:"size: 200;" json:"part_number"`
	PartName     string         `gorm:"size: 200;" json:"part_name"`
	Qty          int            `json:"qty"`
	TotalOrder   int            `json:"total_order"`
	TotalBox     int            `json:"total_box"`
	ProdDate     *time.Time     `json:"prod_date"`
	DelivDate    *time.Time     `json:"deliv_date"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
