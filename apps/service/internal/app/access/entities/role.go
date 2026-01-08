package entities

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID          uint64         `gorm:"primaryKey;column:id" json:"id"`
	Name        string         `gorm:"column:name;not null" json:"name"`
	Description string         `gorm:"column:description;not null" json:"description"`
	ShopID      uint64         `gorm:"column:shop_id;not null" json:"shop_id"`
	CreatedAt   time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
}

func (Role) TableName() string {
	return "roles"
}
