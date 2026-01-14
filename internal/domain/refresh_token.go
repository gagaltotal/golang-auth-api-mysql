package domain

import (
	"time"

	"gorm.io/gorm"
)

type RefreshToken struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"index;not null" json:"user_id"`
	User      User           `gorm:"foreignKey:UserID" json:"-"`
	Token     string         `gorm:"uniqueIndex;size:500;not null" json:"token"`
	ExpiresAt time.Time      `gorm:"not null" json:"expires_at"`
	Revoked   bool           `gorm:"default:false" json:"revoked"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
