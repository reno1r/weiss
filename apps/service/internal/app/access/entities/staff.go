package entities

import (
	"time"

	shopentities "github.com/reno1r/weiss/apps/service/internal/app/shop/entities"
	userentities "github.com/reno1r/weiss/apps/service/internal/app/user/entities"
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

	User *userentities.User `gorm:"foreignKey:UserID" json:"user"`
	Role *Role              `gorm:"foreignKey:RoleID" json:"role"`
	Shop *shopentities.Shop `gorm:"foreignKey:ShopID" json:"shop"`
}

func (Staff) TableName() string {
	return "staffs"
}
