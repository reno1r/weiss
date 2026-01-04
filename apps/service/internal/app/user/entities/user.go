package entities

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint64         `gorm:"primaryKey;column:id" json:"id"`
	FullName  string         `gorm:"column:full_name;not null" json:"full_name"`
	Phone     string         `gorm:"column:phone;not null;uniqueIndex:idx_users_phone" json:"phone"`
	Email     string         `gorm:"column:email;not null;uniqueIndex:idx_users_email" json:"email"`
	Password  string         `gorm:"column:password;not null" json:"password,omitempty"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
}

func (User) TableName() string {
	return "users"
}
