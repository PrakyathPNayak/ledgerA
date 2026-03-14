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

type subcategoryService struct {
	subRepo      repository.SubcategoryRepository
	categoryRepo repository.CategoryRepository
	auditRepo    repository.AuditRepository
}

// NewSubcategoryService creates a new SubcategoryService.
func NewSubcategoryService(subRepo repository.SubcategoryRepository, categoryRepo repository.CategoryRepository, auditRepo repository.AuditRepository) SubcategoryService {
	return &subcategoryService{subRepo: subRepo, categoryRepo: categoryRepo, auditRepo: auditRepo}
}

func (s *subcategoryService) ListByCategory(ctx context.Context, userID uuid.UUID, categoryID uuid.UUID) ([]model.Subcategory, int64, error) {
	items, total, err := s.subRepo.ListByCategoryID(ctx, userID, categoryID)
	if err != nil {
		return nil, 0, fmt.Errorf("subcategoryService.ListByCategory: %w", err)
	}
	return items, total, nil
}

func (s *subcategoryService) Create(ctx context.Context, userID uuid.UUID, categoryID uuid.UUID, req dto.CreateSubcategoryRequest) (*model.Subcategory, error) {
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return nil, fmt.Errorf("subcategoryService.Create: empty subcategory name")
	}

	item := req.ToModel(userID, categoryID)
	if err := s.subRepo.Create(ctx, &item); err != nil {
		if isDuplicateKeyError(err) {
			existing, findErr := s.subRepo.FindByNameForCategory(ctx, userID, categoryID, req.Name)
			if findErr == nil {
				return existing, nil
			}
		}
		return nil, fmt.Errorf("subcategoryService.Create: %w", err)
	}
	if err := s.writeAudit(ctx, userID, item.ID, "create", nil, item); err != nil {
		return nil, fmt.Errorf("subcategoryService.Create.Audit: %w", err)
	}
	return &item, nil
}

func (s *subcategoryService) Update(ctx context.Context, userID uuid.UUID, id uuid.UUID, req dto.UpdateSubcategoryRequest) (*model.Subcategory, error) {
	item, err := s.subRepo.FindByID(ctx, userID, id)
	if err != nil {
		return nil, fmt.Errorf("subcategoryService.Update.Find: %w", err)
	}
	before := *item
	item.Name = req.Name
	if err := s.subRepo.Update(ctx, item); err != nil {
		return nil, fmt.Errorf("subcategoryService.Update.Save: %w", err)
	}
	if err := s.writeAudit(ctx, userID, item.ID, "update", before, item); err != nil {
		return nil, fmt.Errorf("subcategoryService.Update.Audit: %w", err)
	}
	return item, nil
}

func (s *subcategoryService) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	item, err := s.subRepo.FindByID(ctx, userID, id)
	if err != nil {
		return fmt.Errorf("subcategoryService.Delete.Find: %w", err)
	}
	if err := s.subRepo.Delete(ctx, userID, id); err != nil {
		return fmt.Errorf("subcategoryService.Delete.Delete: %w", err)
	}
	if err := s.writeAudit(ctx, userID, id, "delete", item, nil); err != nil {
		return fmt.Errorf("subcategoryService.Delete.Audit: %w", err)
	}
	return nil
}

func (s *subcategoryService) GetOrCreateFromCategory(ctx context.Context, userID uuid.UUID, categoryID uuid.UUID, categoryName string) (*model.Subcategory, error) {
	item, err := s.subRepo.FindByNameForCategory(ctx, userID, categoryID, categoryName)
	if err == nil {
		return item, nil
	}

	newSub := model.Subcategory{UserID: userID, CategoryID: categoryID, Name: categoryName}
	if err := s.subRepo.Create(ctx, &newSub); err != nil {
		return nil, fmt.Errorf("subcategoryService.GetOrCreateFromCategory.Create: %w", err)
	}
	if err := s.writeAudit(ctx, userID, newSub.ID, "create", nil, newSub); err != nil {
		return nil, fmt.Errorf("subcategoryService.GetOrCreateFromCategory.Audit: %w", err)
	}
	return &newSub, nil
}

func (s *subcategoryService) writeAudit(ctx context.Context, userID uuid.UUID, entityID uuid.UUID, action string, before any, after any) error {
	diffPayload := map[string]any{"before": before, "after": after}
	diffBytes, err := json.Marshal(diffPayload)
	if err != nil {
		return fmt.Errorf("subcategoryService.writeAudit.Marshal: %w", err)
	}
	entry := model.AuditLog{UserID: userID, EntityType: "subcategory", EntityID: entityID, Action: action, Diff: string(diffBytes)}
	if err := s.auditRepo.Create(ctx, &entry); err != nil {
		return fmt.Errorf("subcategoryService.writeAudit.Create: %w", err)
	}
	return nil
}
