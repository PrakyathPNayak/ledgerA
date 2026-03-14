package model

import "github.com/google/uuid"

// Subcategory represents a secondary grouping under a Category.
type Subcategory struct {
	BaseModel
	CategoryID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_subcategories_category_name,where:deleted_at IS NULL"`
	UserID     uuid.UUID `gorm:"type:uuid;not null"`
	Name       string    `gorm:"type:text;not null;uniqueIndex:idx_subcategories_category_name,where:deleted_at IS NULL"`
}
