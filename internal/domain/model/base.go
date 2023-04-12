package model

import (
	"time"
)

// BaseModel base model
type BaseModel struct {
	ID        uint       `gorm:"id primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"deleted_at" json:"deleted_at"`
}
