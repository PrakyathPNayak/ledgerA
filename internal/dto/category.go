package dto

import (
	"ledgerA/internal/model"

	"github.com/google/uuid"
)

// CreateCategoryRequest contains category creation payload.
type CreateCategoryRequest struct {
	Name string `json:"name" binding:"required" validate:"required,min=1,max=100"`
}

// UpdateCategoryRequest contains category update payload.
type UpdateCategoryRequest struct {
	Name string `json:"name" binding:"required" validate:"required,min=1,max=100"`
}

// CategoryResponse contains category data sent to clients.
type CategoryResponse struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Name      string    `json:"name"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

// CategoryWithSubsResponse contains category with nested subcategories.
type CategoryWithSubsResponse struct {
	CategoryResponse
	Subcategories []SubcategoryResponse `json:"subcategories"`
}

// ToModel converts CreateCategoryRequest into model.Category.
func (req CreateCategoryRequest) ToModel(userID uuid.UUID) model.Category {
	return model.Category{UserID: userID, Name: req.Name}
}

// NewCategoryResponse creates CategoryResponse from model.Category.
func NewCategoryResponse(category model.Category) CategoryResponse {
	return CategoryResponse{
		ID:        category.ID,
		UserID:    category.UserID,
		Name:      category.Name,
		CreatedAt: category.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: category.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
