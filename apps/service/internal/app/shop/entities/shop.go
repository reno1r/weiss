package entities

import (
	"time"

	"gorm.io/gorm"
)

type Shop struct {
	ID          uint64         `gorm:"primaryKey;column:id" json:"id"`
	Name        string         `gorm:"column:name;not null" json:"name"`
	Description string         `gorm:"column:description;not null" json:"description"`
	Address     string         `gorm:"column:address;not null" json:"address"`
	Phone       string         `gorm:"column:phone;not null" json:"phone"`
	Email       string         `gorm:"column:email;not null" json:"email"`
	Website     string         `gorm:"column:website;not null" json:"website"`
	Logo        string         `gorm:"column:logo;not null" json:"logo"`
	CreatedAt   time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
}

func (Shop) TableName() string {
	return "shops"
}
