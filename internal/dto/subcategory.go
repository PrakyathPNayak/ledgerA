package dto

import (
	"ledgerA/internal/model"

	"github.com/google/uuid"
)

// CreateSubcategoryRequest contains subcategory creation payload.
type CreateSubcategoryRequest struct {
	Name string `json:"name" binding:"required" validate:"required,min=1,max=100"`
}

// UpdateSubcategoryRequest contains subcategory update payload.
type UpdateSubcategoryRequest struct {
	Name string `json:"name" binding:"required" validate:"required,min=1,max=100"`
}

// SubcategoryResponse contains subcategory data sent to clients.
type SubcategoryResponse struct {
	ID         uuid.UUID `json:"id"`
	CategoryID uuid.UUID `json:"category_id"`
	UserID     uuid.UUID `json:"user_id"`
	Name       string    `json:"name"`
	CreatedAt  string    `json:"created_at"`
	UpdatedAt  string    `json:"updated_at"`
}

// ToModel converts CreateSubcategoryRequest into model.Subcategory.
func (req CreateSubcategoryRequest) ToModel(userID uuid.UUID, categoryID uuid.UUID) model.Subcategory {
	return model.Subcategory{UserID: userID, CategoryID: categoryID, Name: req.Name}
}

// NewSubcategoryResponse creates SubcategoryResponse from model.Subcategory.
func NewSubcategoryResponse(sub model.Subcategory) SubcategoryResponse {
	return SubcategoryResponse{
		ID:         sub.ID,
		CategoryID: sub.CategoryID,
		UserID:     sub.UserID,
		Name:       sub.Name,
		CreatedAt:  sub.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  sub.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
