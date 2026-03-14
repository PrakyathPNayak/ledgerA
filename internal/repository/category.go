package repository

import (
	"context"
	"fmt"
	"ledgerA/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type categoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository creates a new CategoryRepository.
func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]model.Category, int64, error) {
	var categories []model.Category
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Category{}).Where("user_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("categoryRepo.ListByUserID.Count: %w", err)
	}
	if err := query.Order("name asc").Find(&categories).Error; err != nil {
		return nil, 0, fmt.Errorf("categoryRepo.ListByUserID.Find: %w", err)
	}
	return categories, total, nil
}

func (r *categoryRepository) FindByID(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*model.Category, error) {
	var category model.Category
	if err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&category).Error; err != nil {
		return nil, fmt.Errorf("categoryRepo.FindByID: %w", err)
	}
	return &category, nil
}

func (r *categoryRepository) FindByName(ctx context.Context, userID uuid.UUID, name string) (*model.Category, error) {
	var category model.Category
	if err := r.db.WithContext(ctx).Where("user_id = ? AND name = ?", userID, name).First(&category).Error; err != nil {
		return nil, fmt.Errorf("categoryRepo.FindByName: %w", err)
	}
	return &category, nil
}

func (r *categoryRepository) Create(ctx context.Context, category *model.Category) error {
	if err := r.db.WithContext(ctx).Create(category).Error; err != nil {
		return fmt.Errorf("categoryRepo.Create: %w", err)
	}
	return nil
}

func (r *categoryRepository) Update(ctx context.Context, category *model.Category) error {
	if err := r.db.WithContext(ctx).Save(category).Error; err != nil {
		return fmt.Errorf("categoryRepo.Update: %w", err)
	}
	return nil
}

func (r *categoryRepository) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&model.Category{}).Error; err != nil {
		return fmt.Errorf("categoryRepo.Delete: %w", err)
	}
	return nil
}
