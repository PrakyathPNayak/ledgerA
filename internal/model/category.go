package model

import "github.com/google/uuid"

// Category represents a top-level grouping for transactions.
type Category struct {
	BaseModel
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_categories_user_name,where:deleted_at IS NULL"`
	Name   string    `gorm:"type:text;not null;uniqueIndex:idx_categories_user_name,where:deleted_at IS NULL"`
}
