package model

import "github.com/google/uuid"

// Budget represents a spending limit for a category in a given period.
type Budget struct {
	BaseModel
	UserID     uuid.UUID `gorm:"type:uuid;not null"`
	CategoryID uuid.UUID `gorm:"type:uuid;not null"`
	Amount     float64   `gorm:"type:numeric(20,4);not null"`
	Period     string    `gorm:"type:text;not null"` // monthly, yearly
	IsActive   bool      `gorm:"type:boolean;not null;default:true"`
}
