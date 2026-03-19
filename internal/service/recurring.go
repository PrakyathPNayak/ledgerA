package service

import (
	"context"
	"encoding/json"
	"fmt"
	"ledgerA/internal/dto"
	"ledgerA/internal/model"
	"ledgerA/internal/repository"
	"time"

	"github.com/google/uuid"
)

type recurringService struct {
	repo           repository.RecurringRepository
	subcategoryRepo repository.SubcategoryRepository
	categoryRepo   repository.CategoryRepository
	txnService     TransactionService
	auditRepo      repository.AuditRepository
}

// NewRecurringService creates a new RecurringService.
func NewRecurringService(repo repository.RecurringRepository, subcategoryRepo repository.SubcategoryRepository, categoryRepo repository.CategoryRepository, txnService TransactionService, auditRepo repository.AuditRepository) RecurringService {
	return &recurringService{repo: repo, subcategoryRepo: subcategoryRepo, categoryRepo: categoryRepo, txnService: txnService, auditRepo: auditRepo}
}

func (s *recurringService) List(ctx context.Context, userID uuid.UUID) ([]model.RecurringTransaction, int64, error) {
	return s.repo.ListByUserID(ctx, userID)
}

func (s *recurringService) Create(ctx context.Context, userID uuid.UUID, req dto.CreateRecurringRequest) (*model.RecurringTransaction, error) {
	startDate, err := dto.ParseTransactionDate(req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start_date: %w", err)
	}

	subcategoryID := uuid.Nil
	if req.SubcategoryID != nil {
		subcategoryID = *req.SubcategoryID
	} else {
		cat, err := s.categoryRepo.FindByID(ctx, userID, req.CategoryID)
		if err != nil {
			return nil, fmt.Errorf("category lookup: %w", err)
		}
		sub, err := s.subcategoryRepo.FindByNameForCategory(ctx, userID, req.CategoryID, cat.Name)
		if err != nil {
			sub = &model.Subcategory{UserID: userID, CategoryID: req.CategoryID, Name: cat.Name}
			if err := s.subcategoryRepo.Create(ctx, sub); err != nil {
				return nil, fmt.Errorf("auto-create subcategory: %w", err)
			}
		}
		subcategoryID = sub.ID
	}

	item := &model.RecurringTransaction{
		UserID:        userID,
		AccountID:     req.AccountID,
		CategoryID:    req.CategoryID,
		SubcategoryID: subcategoryID,
		Name:          req.Name,
		Amount:        req.Amount,
		Notes:         req.Notes,
		Frequency:     req.Frequency,
		StartDate:     startDate,
		NextDueDate:   startDate,
		IsActive:      true,
	}

	if req.EndDate != nil {
		endDate, err := dto.ParseTransactionDate(*req.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end_date: %w", err)
		}
		item.EndDate = &endDate
	}

	if err := s.repo.Create(ctx, item); err != nil {
		return nil, err
	}

	s.writeAudit(ctx, userID, item.ID, "create", nil, item)

	return item, nil
}

func (s *recurringService) writeAudit(ctx context.Context, userID uuid.UUID, entityID uuid.UUID, action string, before any, after any) {
	diffPayload := map[string]any{"before": before, "after": after}
	diffBytes, _ := json.Marshal(diffPayload)
	_ = s.auditRepo.Create(ctx, &model.AuditLog{
		UserID: userID, EntityType: "recurring_transaction", EntityID: entityID, Action: action,
		Diff: string(diffBytes),
	})
}

func (s *recurringService) Update(ctx context.Context, userID uuid.UUID, id uuid.UUID, req dto.UpdateRecurringRequest) (*model.RecurringTransaction, error) {
	item, err := s.repo.FindByID(ctx, userID, id)
	if err != nil {
		return nil, err
	}

	before := *item

	if req.AccountID != nil {
		item.AccountID = *req.AccountID
	}
	if req.CategoryID != nil {
		item.CategoryID = *req.CategoryID
	}
	if req.SubcategoryID != nil {
		item.SubcategoryID = *req.SubcategoryID
	}
	if req.Name != nil {
		item.Name = *req.Name
	}
	if req.Amount != nil {
		item.Amount = *req.Amount
	}
	if req.Notes != nil {
		item.Notes = req.Notes
	}
	if req.Frequency != nil {
		item.Frequency = *req.Frequency
	}
	if req.EndDate != nil {
		endDate, err := dto.ParseTransactionDate(*req.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end_date: %w", err)
		}
		item.EndDate = &endDate
	}
	if req.IsActive != nil {
		item.IsActive = *req.IsActive
	}

	if err := s.repo.Update(ctx, item); err != nil {
		return nil, err
	}

	s.writeAudit(ctx, userID, item.ID, "update", before, item)

	return item, nil
}

func (s *recurringService) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	item, err := s.repo.FindByID(ctx, userID, id)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, userID, id); err != nil {
		return err
	}

	s.writeAudit(ctx, userID, id, "delete", item, nil)

	return nil
}

func (s *recurringService) ProcessDue(ctx context.Context) (int, error) {
	now := time.Now()
	items, err := s.repo.FindDue(ctx, now)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, item := range items {
		txnReq := dto.CreateTransactionRequest{
			AccountID:     item.AccountID,
			CategoryID:    item.CategoryID,
			SubcategoryID: &item.SubcategoryID,
			Name:          item.Name,
			Amount:        item.Amount,
			TransactionDate: item.NextDueDate.Format("2006-01-02"),
			Notes:         item.Notes,
		}

		if _, err := s.txnService.Create(ctx, item.UserID, txnReq); err != nil {
			continue
		}

		now := time.Now()
		item.LastExecutedAt = &now
		item.NextDueDate = dto.ComputeNextDueDate(item.NextDueDate, item.Frequency)

		if item.EndDate != nil && item.NextDueDate.After(*item.EndDate) {
			item.IsActive = false
		}

		_ = s.repo.Update(ctx, &item)
		count++
	}

	return count, nil
}
