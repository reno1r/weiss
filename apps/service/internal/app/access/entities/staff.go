package entities

import (
	"time"

	"gorm.io/gorm"
)

type Staff struct {
	ID        uint64         `gorm:"primaryKey;column:id" json:"id"`
	UserID    uint64         `gorm:"column:user_id;not null" json:"user_id"`
	RoleID    uint64         `gorm:"column:role_id;not null" json:"role_id"`
	ShopID    uint64         `gorm:"column:shop_id;not null" json:"shop_id"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
}

func (Staff) TableName() string {
	return "staffs"
}
