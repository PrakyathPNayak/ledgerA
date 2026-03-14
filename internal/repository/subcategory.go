package repository

import (
	"context"
	"fmt"
	"ledgerA/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type subcategoryRepository struct {
	db *gorm.DB
}

// NewSubcategoryRepository creates a new SubcategoryRepository.
func NewSubcategoryRepository(db *gorm.DB) SubcategoryRepository {
	return &subcategoryRepository{db: db}
}

func (r *subcategoryRepository) ListByCategoryID(ctx context.Context, userID uuid.UUID, categoryID uuid.UUID) ([]model.Subcategory, int64, error) {
	var subs []model.Subcategory
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Subcategory{}).Where("user_id = ? AND category_id = ?", userID, categoryID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("subcategoryRepo.ListByCategoryID.Count: %w", err)
	}
	if err := query.Order("name asc").Find(&subs).Error; err != nil {
		return nil, 0, fmt.Errorf("subcategoryRepo.ListByCategoryID.Find: %w", err)
	}
	return subs, total, nil
}

func (r *subcategoryRepository) FindByID(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*model.Subcategory, error) {
	var sub model.Subcategory
	if err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&sub).Error; err != nil {
		return nil, fmt.Errorf("subcategoryRepo.FindByID: %w", err)
	}
	return &sub, nil
}

func (r *subcategoryRepository) FindByNameForCategory(ctx context.Context, userID uuid.UUID, categoryID uuid.UUID, name string) (*model.Subcategory, error) {
	var sub model.Subcategory
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND category_id = ? AND name = ?", userID, categoryID, name).
		First(&sub).Error; err != nil {
		return nil, fmt.Errorf("subcategoryRepo.FindByNameForCategory: %w", err)
	}
	return &sub, nil
}

func (r *subcategoryRepository) Create(ctx context.Context, subcategory *model.Subcategory) error {
	if err := r.db.WithContext(ctx).Create(subcategory).Error; err != nil {
		return fmt.Errorf("subcategoryRepo.Create: %w", err)
	}
	return nil
}

func (r *subcategoryRepository) Update(ctx context.Context, subcategory *model.Subcategory) error {
	if err := r.db.WithContext(ctx).Save(subcategory).Error; err != nil {
		return fmt.Errorf("subcategoryRepo.Update: %w", err)
	}
	return nil
}

func (r *subcategoryRepository) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&model.Subcategory{}).Error; err != nil {
		return fmt.Errorf("subcategoryRepo.Delete: %w", err)
	}
	return nil
}
