package models

import "time"

type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

type Todo struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Title       string     `gorm:"size:200;not null" json:"title" validate:"required,min=3,max=200"`
	Description string     `gorm:"type:text" json:"description" validate:"max=2000"`
	Completed   bool       `gorm:"default:false" json:"completed"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	Priority    Priority   `gorm:"size:10;default:medium" json:"priority" validate:"oneof=low medium high"`
	OwnerID     uint       `json:"owner_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
