package entities

import "gorm.io/gorm"

type RefreshToken struct {
	gorm.Model
	Token     string `gorm:"index;not null" json:"token"`
	ExpiresAt int64  `gorm:"not null" json:"expiresAt"`
	UserID    uint   `gorm:"not null" json:"userId"`
}
