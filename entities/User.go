package entities

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username      string         `gorm:"not null" json:"username"`
	Email         string         `gorm:"unique_index;not null" json:"email"`
	Password      string         `gorm:"not null" json:"password"`
	RefreshTokens []RefreshToken `gorm:"not null" json:"refreshTokens"`
}
