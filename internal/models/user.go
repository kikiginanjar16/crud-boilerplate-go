package models

import "time"

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"size:120;not null" json:"name" validate:"required,min=2,max=120"`
	Email        string    `gorm:"size:180;uniqueIndex;not null" json:"email" validate:"required,email"`
	PasswordHash string    `gorm:"size:255;not null" json:"-"`
	Role         Role      `gorm:"size:20;default:user" json:"role" validate:"oneof=admin user"`
	AvatarURL    string    `gorm:"size:255" json:"avatar_url"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
