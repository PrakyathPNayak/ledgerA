package service

import (
	"context"
	"encoding/json"
	"fmt"
	"ledgerA/internal/dto"
	"ledgerA/internal/model"
	"ledgerA/internal/repository"
	"strings"

	"github.com/google/uuid"
)

type categoryService struct {
	categoryRepo repository.CategoryRepository
	auditRepo    repository.AuditRepository
}

// NewCategoryService creates a new CategoryService.
func NewCategoryService(categoryRepo repository.CategoryRepository, auditRepo repository.AuditRepository) CategoryService {
	return &categoryService{categoryRepo: categoryRepo, auditRepo: auditRepo}
}

func (s *categoryService) List(ctx context.Context, userID uuid.UUID) ([]model.Category, int64, error) {
	items, total, err := s.categoryRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("categoryService.List: %w", err)
	}
	return items, total, nil
}

func (s *categoryService) Create(ctx context.Context, userID uuid.UUID, req dto.CreateCategoryRequest) (*model.Category, error) {
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return nil, fmt.Errorf("categoryService.Create: empty category name")
	}

	item := req.ToModel(userID)
	if err := s.categoryRepo.Create(ctx, &item); err != nil {
		if isDuplicateKeyError(err) {
			existing, findErr := s.categoryRepo.FindByName(ctx, userID, req.Name)
			if findErr == nil {
				return existing, nil
			}
		}
		return nil, fmt.Errorf("categoryService.Create: %w", err)
	}
	if err := s.writeAudit(ctx, userID, item.ID, "create", nil, item); err != nil {
		return nil, fmt.Errorf("categoryService.Create.Audit: %w", err)
	}
	return &item, nil
}

func (s *categoryService) Update(ctx context.Context, userID uuid.UUID, id uuid.UUID, req dto.UpdateCategoryRequest) (*model.Category, error) {
	item, err := s.categoryRepo.FindByID(ctx, userID, id)
	if err != nil {
		return nil, fmt.Errorf("categoryService.Update.Find: %w", err)
	}
	before := *item
	item.Name = req.Name
	if err := s.categoryRepo.Update(ctx, item); err != nil {
		return nil, fmt.Errorf("categoryService.Update.Save: %w", err)
	}
	if err := s.writeAudit(ctx, userID, item.ID, "update", before, item); err != nil {
		return nil, fmt.Errorf("categoryService.Update.Audit: %w", err)
	}
	return item, nil
}

func (s *categoryService) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	item, err := s.categoryRepo.FindByID(ctx, userID, id)
	if err != nil {
		return fmt.Errorf("categoryService.Delete.Find: %w", err)
	}
	if err := s.categoryRepo.Delete(ctx, userID, id); err != nil {
		return fmt.Errorf("categoryService.Delete.Delete: %w", err)
	}
	if err := s.writeAudit(ctx, userID, id, "delete", item, nil); err != nil {
		return fmt.Errorf("categoryService.Delete.Audit: %w", err)
	}
	return nil
}

func (s *categoryService) writeAudit(ctx context.Context, userID uuid.UUID, entityID uuid.UUID, action string, before any, after any) error {
	diffPayload := map[string]any{"before": before, "after": after}
	diffBytes, err := json.Marshal(diffPayload)
	if err != nil {
		return fmt.Errorf("categoryService.writeAudit.Marshal: %w", err)
	}
	entry := model.AuditLog{UserID: userID, EntityType: "category", EntityID: entityID, Action: action, Diff: string(diffBytes)}
	if err := s.auditRepo.Create(ctx, &entry); err != nil {
		return fmt.Errorf("categoryService.writeAudit.Create: %w", err)
	}
	return nil
}
